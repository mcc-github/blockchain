package govaluate

import (
	"errors"
	"time"
	"fmt"
)

var stageSymbolMap = map[OperatorSymbol]evaluationOperator{
	EQ:             equalStage,
	NEQ:            notEqualStage,
	GT:             gtStage,
	LT:             ltStage,
	GTE:            gteStage,
	LTE:            lteStage,
	REQ:            regexStage,
	NREQ:           notRegexStage,
	AND:            andStage,
	OR:             orStage,
	IN:             inStage,
	BITWISE_OR:     bitwiseOrStage,
	BITWISE_AND:    bitwiseAndStage,
	BITWISE_XOR:    bitwiseXORStage,
	BITWISE_LSHIFT: leftShiftStage,
	BITWISE_RSHIFT: rightShiftStage,
	PLUS:           addStage,
	MINUS:          subtractStage,
	MULTIPLY:       multiplyStage,
	DIVIDE:         divideStage,
	MODULUS:        modulusStage,
	EXPONENT:       exponentStage,
	NEGATE:         negateStage,
	INVERT:         invertStage,
	BITWISE_NOT:    bitwiseNotStage,
	TERNARY_TRUE:   ternaryIfStage,
	TERNARY_FALSE:  ternaryElseStage,
	COALESCE:       ternaryElseStage,
	SEPARATE:       separatorStage,
}


type precedent func(stream *tokenStream) (*evaluationStage, error)


type precedencePlanner struct {
	validSymbols map[string]OperatorSymbol
	validKinds   []TokenKind

	typeErrorFormat string

	next      precedent
	nextRight precedent
}

var planPrefix precedent
var planExponential precedent
var planMultiplicative precedent
var planAdditive precedent
var planBitwise precedent
var planShift precedent
var planComparator precedent
var planLogicalAnd precedent
var planLogicalOr precedent
var planTernary precedent
var planSeparator precedent

func init() {

	
	
	
	planPrefix = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    prefixSymbols,
		validKinds:      []TokenKind{PREFIX},
		typeErrorFormat: prefixErrorFormat,
		nextRight:       planFunction,
	})
	planExponential = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    exponentialSymbolsS,
		validKinds:      []TokenKind{MODIFIER},
		typeErrorFormat: modifierErrorFormat,
		next:            planFunction,
	})
	planMultiplicative = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    multiplicativeSymbols,
		validKinds:      []TokenKind{MODIFIER},
		typeErrorFormat: modifierErrorFormat,
		next:            planExponential,
	})
	planAdditive = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    additiveSymbols,
		validKinds:      []TokenKind{MODIFIER},
		typeErrorFormat: modifierErrorFormat,
		next:            planMultiplicative,
	})
	planShift = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    bitwiseShiftSymbols,
		validKinds:      []TokenKind{MODIFIER},
		typeErrorFormat: modifierErrorFormat,
		next:            planAdditive,
	})
	planBitwise = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    bitwiseSymbols,
		validKinds:      []TokenKind{MODIFIER},
		typeErrorFormat: modifierErrorFormat,
		next:            planShift,
	})
	planComparator = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    comparatorSymbols,
		validKinds:      []TokenKind{COMPARATOR},
		typeErrorFormat: comparatorErrorFormat,
		next:            planBitwise,
	})
	planLogicalAnd = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    map[string]OperatorSymbol{"&&": AND},
		validKinds:      []TokenKind{LOGICALOP},
		typeErrorFormat: logicalErrorFormat,
		next:            planComparator,
	})
	planLogicalOr = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    map[string]OperatorSymbol{"||": OR},
		validKinds:      []TokenKind{LOGICALOP},
		typeErrorFormat: logicalErrorFormat,
		next:            planLogicalAnd,
	})
	planTernary = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols:    ternarySymbols,
		validKinds:      []TokenKind{TERNARY},
		typeErrorFormat: ternaryErrorFormat,
		next:            planLogicalOr,
	})
	planSeparator = makePrecedentFromPlanner(&precedencePlanner{
		validSymbols: separatorSymbols,
		validKinds:   []TokenKind{SEPARATOR},
		next:         planTernary,
	})
}


func makePrecedentFromPlanner(planner *precedencePlanner) precedent {

	var generated precedent
	var nextRight precedent

	generated = func(stream *tokenStream) (*evaluationStage, error) {
		return planPrecedenceLevel(
			stream,
			planner.typeErrorFormat,
			planner.validSymbols,
			planner.validKinds,
			nextRight,
			planner.next,
		)
	}

	if planner.nextRight != nil {
		nextRight = planner.nextRight
	} else {
		nextRight = generated
	}

	return generated
}


func planStages(tokens []ExpressionToken) (*evaluationStage, error) {

	stream := newTokenStream(tokens)

	stage, err := planTokens(stream)
	if err != nil {
		return nil, err
	}

	
	
	reorderStages(stage)

	stage = elideLiterals(stage)
	return stage, nil
}

func planTokens(stream *tokenStream) (*evaluationStage, error) {

	if !stream.hasNext() {
		return nil, nil
	}

	return planSeparator(stream)
}


func planPrecedenceLevel(
	stream *tokenStream,
	typeErrorFormat string,
	validSymbols map[string]OperatorSymbol,
	validKinds []TokenKind,
	rightPrecedent precedent,
	leftPrecedent precedent) (*evaluationStage, error) {

	var token ExpressionToken
	var symbol OperatorSymbol
	var leftStage, rightStage *evaluationStage
	var checks typeChecks
	var err error
	var keyFound bool

	if leftPrecedent != nil {

		leftStage, err = leftPrecedent(stream)
		if err != nil {
			return nil, err
		}
	}

	for stream.hasNext() {

		token = stream.next()

		if len(validKinds) > 0 {

			keyFound = false
			for _, kind := range validKinds {
				if kind == token.Kind {
					keyFound = true
					break
				}
			}

			if !keyFound {
				break
			}
		}

		if validSymbols != nil {

			if !isString(token.Value) {
				break
			}

			symbol, keyFound = validSymbols[token.Value.(string)]
			if !keyFound {
				break
			}
		}

		if rightPrecedent != nil {
			rightStage, err = rightPrecedent(stream)
			if err != nil {
				return nil, err
			}
		}

		checks = findTypeChecks(symbol)

		return &evaluationStage{

			symbol:     symbol,
			leftStage:  leftStage,
			rightStage: rightStage,
			operator:   stageSymbolMap[symbol],

			leftTypeCheck:   checks.left,
			rightTypeCheck:  checks.right,
			typeCheck:       checks.combined,
			typeErrorFormat: typeErrorFormat,
		}, nil
	}

	stream.rewind()
	return leftStage, nil
}


func planFunction(stream *tokenStream) (*evaluationStage, error) {

	var token ExpressionToken
	var rightStage *evaluationStage
	var err error

	token = stream.next()

	if token.Kind != FUNCTION {
		stream.rewind()
		return planValue(stream)
	}

	rightStage, err = planValue(stream)
	if err != nil {
		return nil, err
	}

	return &evaluationStage{

		symbol:          FUNCTIONAL,
		rightStage:      rightStage,
		operator:        makeFunctionStage(token.Value.(ExpressionFunction)),
		typeErrorFormat: "Unable to run function '%v': %v",
	}, nil
}


func planValue(stream *tokenStream) (*evaluationStage, error) {

	var token ExpressionToken
	var symbol OperatorSymbol
	var ret *evaluationStage
	var operator evaluationOperator
	var err error

	token = stream.next()

	switch token.Kind {

	case CLAUSE:

		ret, err = planTokens(stream)
		if err != nil {
			return nil, err
		}

		
		stream.next()

		
		
		
		ret = &evaluationStage {
			rightStage: ret,
			operator: noopStageRight,
			symbol: NOOP,
		}

		return ret, nil

	case CLAUSE_CLOSE:

		
		
		stream.rewind()
		return nil, nil

	case VARIABLE:
		operator = makeParameterStage(token.Value.(string))

	case NUMERIC:
		fallthrough
	case STRING:
		fallthrough
	case PATTERN:
		fallthrough
	case BOOLEAN:
		symbol = LITERAL
		operator = makeLiteralStage(token.Value)
	case TIME:
		symbol = LITERAL
		operator = makeLiteralStage(float64(token.Value.(time.Time).Unix()))

	case PREFIX:
		stream.rewind()
		return planPrefix(stream)
	}

	if operator == nil {
		errorMsg := fmt.Sprintf("Unable to plan token kind: '%s', value: '%v'", token.Kind.String(), token.Value)
		return nil, errors.New(errorMsg)
	}

	return &evaluationStage{
		symbol: symbol,
		operator: operator,
	}, nil
}


type typeChecks struct {
	left     stageTypeCheck
	right    stageTypeCheck
	combined stageCombinedTypeCheck
}


func findTypeChecks(symbol OperatorSymbol) typeChecks {

	switch symbol {
	case GT:
		fallthrough
	case LT:
		fallthrough
	case GTE:
		fallthrough
	case LTE:
		return typeChecks{
			combined: comparatorTypeCheck,
		}
	case REQ:
		fallthrough
	case NREQ:
		return typeChecks{
			left:  isString,
			right: isRegexOrString,
		}
	case AND:
		fallthrough
	case OR:
		return typeChecks{
			left:  isBool,
			right: isBool,
		}
	case IN:
		return typeChecks{
			right: isArray,
		}
	case BITWISE_LSHIFT:
		fallthrough
	case BITWISE_RSHIFT:
		fallthrough
	case BITWISE_OR:
		fallthrough
	case BITWISE_AND:
		fallthrough
	case BITWISE_XOR:
		return typeChecks{
			left:  isFloat64,
			right: isFloat64,
		}
	case PLUS:
		return typeChecks{
			combined: additionTypeCheck,
		}
	case MINUS:
		fallthrough
	case MULTIPLY:
		fallthrough
	case DIVIDE:
		fallthrough
	case MODULUS:
		fallthrough
	case EXPONENT:
		return typeChecks{
			left:  isFloat64,
			right: isFloat64,
		}
	case NEGATE:
		return typeChecks{
			right: isFloat64,
		}
	case INVERT:
		return typeChecks{
			right: isBool,
		}
	case BITWISE_NOT:
		return typeChecks{
			right: isFloat64,
		}
	case TERNARY_TRUE:
		return typeChecks{
			left: isBool,
		}

	
	case EQ:
		fallthrough
	case NEQ:
		return typeChecks{}
	case TERNARY_FALSE:
		fallthrough
	case COALESCE:
		fallthrough
	default:
		return typeChecks{}
	}
}


func reorderStages(rootStage *evaluationStage) {

	
	var identicalPrecedences []*evaluationStage
	var currentStage, nextStage *evaluationStage
	var precedence, currentPrecedence operatorPrecedence

	nextStage = rootStage
	precedence = findOperatorPrecedenceForSymbol(rootStage.symbol)

	for nextStage != nil {

		currentStage = nextStage
		nextStage = currentStage.rightStage

		
		if currentStage.leftStage != nil {
			reorderStages(currentStage.leftStage)
		}

		currentPrecedence = findOperatorPrecedenceForSymbol(currentStage.symbol)

		if currentPrecedence == precedence {
			identicalPrecedences = append(identicalPrecedences, currentStage)
			continue
		}

		
		
		if len(identicalPrecedences) > 1 {
			mirrorStageSubtree(identicalPrecedences)
		}

		identicalPrecedences = []*evaluationStage{currentStage}
		precedence = currentPrecedence
	}

	if len(identicalPrecedences) > 1 {
		mirrorStageSubtree(identicalPrecedences)
	}
}


func mirrorStageSubtree(stages []*evaluationStage) {

	var rootStage, inverseStage, carryStage, frontStage *evaluationStage

	stagesLength := len(stages)

	
	for _, frontStage = range stages {

		carryStage = frontStage.rightStage
		frontStage.rightStage = frontStage.leftStage
		frontStage.leftStage = carryStage
	}

	
	rootStage = stages[0]
	frontStage = stages[stagesLength-1]

	carryStage = frontStage.leftStage
	frontStage.leftStage = rootStage.rightStage
	rootStage.rightStage = carryStage

	
	for i := 0; i < (stagesLength-2)/2+1; i++ {

		frontStage = stages[i+1]
		inverseStage = stages[stagesLength-i-1]

		carryStage = frontStage.rightStage
		frontStage.rightStage = inverseStage.rightStage
		inverseStage.rightStage = carryStage
	}

	
	for i := 0; i < stagesLength/2; i++ {

		frontStage = stages[i]
		inverseStage = stages[stagesLength-i-1]
		frontStage.swapWith(inverseStage)
	}
}


func elideLiterals(root *evaluationStage) *evaluationStage {

	if root.leftStage != nil {
		root.leftStage = elideLiterals(root.leftStage)
	}

	if root.rightStage != nil {
		root.rightStage = elideLiterals(root.rightStage)
	}

	return elideStage(root)
}


func elideStage(root *evaluationStage) *evaluationStage {

	var leftValue, rightValue, result interface{}
	var err error

	
	if root.rightStage == nil ||
		root.rightStage.symbol != LITERAL ||
		root.leftStage == nil ||
		root.leftStage.symbol != LITERAL {
		return root
	}

	
	switch root.symbol {
	case SEPARATE:
		fallthrough
	case IN:
		return root
	}

	
	
	leftValue, err = root.leftStage.operator(nil, nil, nil)
	if err != nil {
		return root
	}

	rightValue, err = root.rightStage.operator(nil, nil, nil)
	if err != nil {
		return root
	}

	
	err = typeCheck(root.leftTypeCheck, leftValue, root.symbol, root.typeErrorFormat)
	if err != nil {
		return root
	}

	err = typeCheck(root.rightTypeCheck, rightValue, root.symbol, root.typeErrorFormat)
	if err != nil {
		return root
	}

	if root.typeCheck != nil && !root.typeCheck(leftValue, rightValue) {
		return root
	}

	
	result, err = root.operator(leftValue, rightValue, nil)
	if err != nil {
		return root
	}

	return &evaluationStage {
		symbol: LITERAL,
		operator: makeLiteralStage(result),
	}
}
