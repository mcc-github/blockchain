



package gotty

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var exp = [...]string{
	"%%",
	"%c",
	"%s",
	"%p(\\d)",
	"%P([A-z])",
	"%g([A-z])",
	"%'(.)'",
	"%{([0-9]+)}",
	"%l",
	"%\\+|%-|%\\*|%/|%m",
	"%&|%\\||%\\^",
	"%=|%>|%<",
	"%A|%O",
	"%!|%~",
	"%i",
	"%(:[\\ #\\-\\+]{0,4})?(\\d+\\.\\d+|\\d+)?[doxXs]",
	"%\\?(.*?);",
}

var regex *regexp.Regexp
var staticVar map[byte]stacker


func (term *TermInfo) Parse(attr string, params ...interface{}) (string, error) {
	
	iface, err := term.GetAttribute(attr)
	str, ok := iface.(string)
	if err != nil {
		return "", err
	}
	if !ok {
		return str, errors.New("Only string capabilities can be parsed.")
	}
	
	
	ps := &parser{}
	
	ps.dynamicVar = make(map[byte]stacker, 26)
	ps.parameters = make([]stacker, len(params))
	
	for i, x := range params {
		ps.parameters[i] = x
	}
	
	result, err := ps.walk(str)
	return result, err
}




func (term *TermInfo) ParseName(attr string, params ...interface{}) (string, error) {
	tc := GetTermcapName(attr)
	return term.Parse(tc, params)
}


func (ps *parser) walk(attr string) (string, error) {
	
	var buf bytes.Buffer
	
	tokens := regex.FindAllStringSubmatch(attr, -1)
	if len(tokens) == 0 {
		return attr, nil
	}
	indices := regex.FindAllStringIndex(attr, -1)
	q := 0 
	
	for i := 0; i < len(attr); i++ {
		
		
		if q < len(indices) && i >= indices[q][0] && i < indices[q][1] {
			
			switch {
			case tokens[q][0][:2] == "%%":
				
				buf.WriteByte('%')
			case tokens[q][0][:2] == "%c":
				
				c, err := ps.st.pop()
				if err != nil {
					return buf.String(), err
				}
				buf.WriteByte(c.(byte))
			case tokens[q][0][:2] == "%s":
				
				str, err := ps.st.pop()
				if err != nil {
					return buf.String(), err
				}
				if _, ok := str.(string); !ok {
					return buf.String(), errors.New("Stack head is not a string")
				}
				buf.WriteString(str.(string))
			case tokens[q][0][:2] == "%p":
				
				index, err := strconv.ParseInt(tokens[q][1], 10, 8)
				index--
				if err != nil {
					return buf.String(), err
				}
				if int(index) >= len(ps.parameters) {
					return buf.String(), errors.New("Parameters index out of bound")
				}
				ps.st.push(ps.parameters[index])
			case tokens[q][0][:2] == "%P":
				
				val, err := ps.st.pop()
				if err != nil {
					return buf.String(), err
				}
				index := tokens[q][2]
				if len(index) > 1 {
					errorStr := fmt.Sprintf("%s is not a valid dynamic variables index",
						index)
					return buf.String(), errors.New(errorStr)
				}
				
				if index[0] >= 'a' && index[0] <= 'z' {
					ps.dynamicVar[index[0]] = val
				} else if index[0] >= 'A' && index[0] <= 'Z' {
					staticVar[index[0]] = val
				}
			case tokens[q][0][:2] == "%g":
				
				index := tokens[q][3]
				if len(index) > 1 {
					errorStr := fmt.Sprintf("%s is not a valid static variables index",
						index)
					return buf.String(), errors.New(errorStr)
				}
				var val stacker
				if index[0] >= 'a' && index[0] <= 'z' {
					val = ps.dynamicVar[index[0]]
				} else if index[0] >= 'A' && index[0] <= 'Z' {
					val = staticVar[index[0]]
				}
				ps.st.push(val)
			case tokens[q][0][:2] == "%'":
				
				con := tokens[q][4]
				if len(con) > 1 {
					errorStr := fmt.Sprintf("%s is not a valid character constant", con)
					return buf.String(), errors.New(errorStr)
				}
				ps.st.push(con[0])
			case tokens[q][0][:2] == "%{":
				
				con, err := strconv.ParseInt(tokens[q][5], 10, 32)
				if err != nil {
					return buf.String(), err
				}
				ps.st.push(con)
			case tokens[q][0][:2] == "%l":
				
				popStr, err := ps.st.pop()
				if err != nil {
					return buf.String(), err
				}
				if _, ok := popStr.(string); !ok {
					errStr := fmt.Sprintf("Stack head is not a string")
					return buf.String(), errors.New(errStr)
				}
				ps.st.push(len(popStr.(string)))
			case tokens[q][0][:2] == "%?":
				
				
				ifReg, _ := regexp.Compile("%\\?(.*)%t(.*)%e(.*);|%\\?(.*)%t(.*);")
				ifTokens := ifReg.FindStringSubmatch(tokens[q][0])
				var (
					ifStr string
					err   error
				)
				
				if len(ifTokens[1]) > 0 {
					ifStr, err = ps.walk(ifTokens[1])
				} else { 
					ifStr, err = ps.walk(ifTokens[4])
				}
				
				if err != nil {
					return buf.String(), err
				} else if len(ifStr) > 0 {
					
					
					return buf.String(), errors.New("If-clause cannot print statements")
				}
				var thenStr string
				
				choose, err := ps.st.pop()
				if err != nil {
					return buf.String(), err
				}
				
				if choose.(int) == 0 && len(ifTokens[1]) > 0 {
					thenStr, err = ps.walk(ifTokens[3])
				} else if choose.(int) != 0 {
					if len(ifTokens[1]) > 0 {
						thenStr, err = ps.walk(ifTokens[2])
					} else {
						thenStr, err = ps.walk(ifTokens[5])
					}
				}
				if err != nil {
					return buf.String(), err
				}
				buf.WriteString(thenStr)
			case tokens[q][0][len(tokens[q][0])-1] == 'd': 
				fallthrough
			case tokens[q][0][len(tokens[q][0])-1] == 'o': 
				fallthrough
			case tokens[q][0][len(tokens[q][0])-1] == 'x':
				fallthrough
			case tokens[q][0][len(tokens[q][0])-1] == 'X':
				fallthrough
			case tokens[q][0][len(tokens[q][0])-1] == 's':
				token := tokens[q][0]
				
				if token[1] == ':' {
					token = token[:1] + token[2:]
				}
				digit, err := ps.st.pop()
				if err != nil {
					return buf.String(), err
				}
				
				digitStr := fmt.Sprintf(token, digit.(int))
				buf.WriteString(digitStr)
			case tokens[q][0][:2] == "%i":
				
				if len(ps.parameters) < 2 {
					return buf.String(), errors.New("Not enough parameters to increment.")
				}
				val1, val2 := ps.parameters[0].(int), ps.parameters[1].(int)
				val1++
				val2++
				ps.parameters[0], ps.parameters[1] = val1, val2
			default:
				
				
				op1, err := ps.st.pop()
				if err != nil {
					return buf.String(), err
				}
				op2, err := ps.st.pop()
				if err != nil {
					return buf.String(), err
				}
				var result stacker
				switch tokens[q][0][:2] {
				case "%+":
					
					result = op2.(int) + op1.(int)
				case "%-":
					
					result = op2.(int) - op1.(int)
				case "%*":
					
					result = op2.(int) * op1.(int)
				case "%/":
					
					result = op2.(int) / op1.(int)
				case "%m":
					
					result = op2.(int) % op1.(int)
				case "%&":
					
					result = op2.(int) & op1.(int)
				case "%|":
					
					result = op2.(int) | op1.(int)
				case "%^":
					
					result = op2.(int) ^ op1.(int)
				case "%=":
					
					result = op2 == op1
				case "%>":
					
					result = op2.(int) > op1.(int)
				case "%<":
					
					result = op2.(int) < op1.(int)
				case "%A":
					
					result = op2.(bool) && op1.(bool)
				case "%O":
					
					result = op2.(bool) || op1.(bool)
				case "%!":
					
					result = !op1.(bool)
				case "%~":
					
					result = ^(op1.(int))
				}
				ps.st.push(result)
			}

			i = indices[q][1] - 1
			q++
		} else {
			
			
			j := i
			if q != len(indices) {
				for !(j >= indices[q][0] && j < indices[q][1]) {
					j++
				}
			} else {
				j = len(attr)
			}
			buf.WriteString(string(attr[i:j]))
			i = j
		}
	}
	
	return buf.String(), nil
}


func (st *stack) push(s stacker) {
	*st = append(*st, s)
}


func (st *stack) pop() (stacker, error) {
	if len(*st) == 0 {
		return nil, errors.New("Stack is empty.")
	}
	newStack := make(stack, len(*st)-1)
	val := (*st)[len(*st)-1]
	copy(newStack, (*st)[:len(*st)-1])
	*st = newStack
	return val, nil
}



func init() {
	
	expStr := strings.Join(exp[:], "|")
	regex, _ = regexp.Compile(expStr)
	
	staticVar = make(map[byte]stacker, 26)
}
