package govaluate

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)


func (this EvaluableExpression) ToSQLQuery() (string, error) {

	var stream *tokenStream
	var transactions *expressionOutputStream
	var transaction string
	var err error

	stream = newTokenStream(this.tokens)
	transactions = new(expressionOutputStream)

	for stream.hasNext() {

		transaction, err = this.findNextSQLString(stream, transactions)
		if err != nil {
			return "", err
		}

		transactions.add(transaction)
	}

	return transactions.createString(" "), nil
}

func (this EvaluableExpression) findNextSQLString(stream *tokenStream, transactions *expressionOutputStream) (string, error) {

	var token ExpressionToken
	var ret string

	token = stream.next()

	switch token.Kind {

	case STRING:
		ret = fmt.Sprintf("'%v'", token.Value)
	case PATTERN:
		ret = fmt.Sprintf("'%s'", token.Value.(*regexp.Regexp).String())
	case TIME:
		ret = fmt.Sprintf("'%s'", token.Value.(time.Time).Format(this.QueryDateFormat))

	case LOGICALOP:
		switch logicalSymbols[token.Value.(string)] {

		case AND:
			ret = "AND"
		case OR:
			ret = "OR"
		}

	case BOOLEAN:
		if token.Value.(bool) {
			ret = "1"
		} else {
			ret = "0"
		}

	case VARIABLE:
		ret = fmt.Sprintf("[%s]", token.Value.(string))

	case NUMERIC:
		ret = fmt.Sprintf("%g", token.Value.(float64))

	case COMPARATOR:
		switch comparatorSymbols[token.Value.(string)] {

		case EQ:
			ret = "="
		case NEQ:
			ret = "<>"
		case REQ:
			ret = "RLIKE"
		case NREQ:
			ret = "NOT RLIKE"
		default:
			ret = fmt.Sprintf("%s", token.Value.(string))
		}

	case TERNARY:

		switch ternarySymbols[token.Value.(string)] {

		case COALESCE:

			left := transactions.rollback()
			right, err := this.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("COALESCE(%v, %v)", left, right)
		case TERNARY_TRUE:
			fallthrough
		case TERNARY_FALSE:
			return "", errors.New("Ternary operators are unsupported in SQL output")
		}
	case PREFIX:
		switch prefixSymbols[token.Value.(string)] {

		case INVERT:
			ret = fmt.Sprintf("NOT")
		default:

			right, err := this.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("%s%s", token.Value.(string), right)
		}
	case MODIFIER:

		switch modifierSymbols[token.Value.(string)] {

		case EXPONENT:

			left := transactions.rollback()
			right, err := this.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("POW(%s, %s)", left, right)
		case MODULUS:

			left := transactions.rollback()
			right, err := this.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("MOD(%s, %s)", left, right)
		default:
			ret = fmt.Sprintf("%s", token.Value.(string))
		}
	case CLAUSE:
		ret = "("
	case CLAUSE_CLOSE:
		ret = ")"
	case SEPARATOR:
		ret = ","

	default:
		errorMsg := fmt.Sprintf("Unrecognized query token '%s' of kind '%s'", token.Value, token.Kind)
		return "", errors.New(errorMsg)
	}

	return ret, nil
}
