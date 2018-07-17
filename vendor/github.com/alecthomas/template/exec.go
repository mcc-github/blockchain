



package template

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/alecthomas/template/parse"
)




type state struct {
	tmpl *Template
	wr   io.Writer
	node parse.Node 
	vars []variable 
}


type variable struct {
	name  string
	value reflect.Value
}


func (s *state) push(name string, value reflect.Value) {
	s.vars = append(s.vars, variable{name, value})
}


func (s *state) mark() int {
	return len(s.vars)
}


func (s *state) pop(mark int) {
	s.vars = s.vars[0:mark]
}


func (s *state) setVar(n int, value reflect.Value) {
	s.vars[len(s.vars)-n].value = value
}


func (s *state) varValue(name string) reflect.Value {
	for i := s.mark() - 1; i >= 0; i-- {
		if s.vars[i].name == name {
			return s.vars[i].value
		}
	}
	s.errorf("undefined variable: %s", name)
	return zero
}

var zero reflect.Value


func (s *state) at(node parse.Node) {
	s.node = node
}



func doublePercent(str string) string {
	if strings.Contains(str, "%") {
		str = strings.Replace(str, "%", "%%", -1)
	}
	return str
}


func (s *state) errorf(format string, args ...interface{}) {
	name := doublePercent(s.tmpl.Name())
	if s.node == nil {
		format = fmt.Sprintf("template: %s: %s", name, format)
	} else {
		location, context := s.tmpl.ErrorContext(s.node)
		format = fmt.Sprintf("template: %s: executing %q at <%s>: %s", location, name, doublePercent(context), format)
	}
	panic(fmt.Errorf(format, args...))
}



func errRecover(errp *error) {
	e := recover()
	if e != nil {
		switch err := e.(type) {
		case runtime.Error:
			panic(e)
		case error:
			*errp = err
		default:
			panic(e)
		}
	}
}







func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	tmpl := t.tmpl[name]
	if tmpl == nil {
		return fmt.Errorf("template: no template %q associated with template %q", name, t.name)
	}
	return tmpl.Execute(wr, data)
}







func (t *Template) Execute(wr io.Writer, data interface{}) (err error) {
	defer errRecover(&err)
	value := reflect.ValueOf(data)
	state := &state{
		tmpl: t,
		wr:   wr,
		vars: []variable{{"$", value}},
	}
	t.init()
	if t.Tree == nil || t.Root == nil {
		var b bytes.Buffer
		for name, tmpl := range t.tmpl {
			if tmpl.Tree == nil || tmpl.Root == nil {
				continue
			}
			if b.Len() > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%q", name)
		}
		var s string
		if b.Len() > 0 {
			s = "; defined templates are: " + b.String()
		}
		state.errorf("%q is an incomplete or empty template%s", t.Name(), s)
	}
	state.walk(value, t.Root)
	return
}



func (s *state) walk(dot reflect.Value, node parse.Node) {
	s.at(node)
	switch node := node.(type) {
	case *parse.ActionNode:
		
		
		val := s.evalPipeline(dot, node.Pipe)
		if len(node.Pipe.Decl) == 0 {
			s.printValue(node, val)
		}
	case *parse.IfNode:
		s.walkIfOrWith(parse.NodeIf, dot, node.Pipe, node.List, node.ElseList)
	case *parse.ListNode:
		for _, node := range node.Nodes {
			s.walk(dot, node)
		}
	case *parse.RangeNode:
		s.walkRange(dot, node)
	case *parse.TemplateNode:
		s.walkTemplate(dot, node)
	case *parse.TextNode:
		if _, err := s.wr.Write(node.Text); err != nil {
			s.errorf("%s", err)
		}
	case *parse.WithNode:
		s.walkIfOrWith(parse.NodeWith, dot, node.Pipe, node.List, node.ElseList)
	default:
		s.errorf("unknown node: %s", node)
	}
}



func (s *state) walkIfOrWith(typ parse.NodeType, dot reflect.Value, pipe *parse.PipeNode, list, elseList *parse.ListNode) {
	defer s.pop(s.mark())
	val := s.evalPipeline(dot, pipe)
	truth, ok := isTrue(val)
	if !ok {
		s.errorf("if/with can't use %v", val)
	}
	if truth {
		if typ == parse.NodeWith {
			s.walk(val, list)
		} else {
			s.walk(dot, list)
		}
	} else if elseList != nil {
		s.walk(dot, elseList)
	}
}



func isTrue(val reflect.Value) (truth, ok bool) {
	if !val.IsValid() {
		
		return false, true
	}
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		truth = val.Len() > 0
	case reflect.Bool:
		truth = val.Bool()
	case reflect.Complex64, reflect.Complex128:
		truth = val.Complex() != 0
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface:
		truth = !val.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		truth = val.Int() != 0
	case reflect.Float32, reflect.Float64:
		truth = val.Float() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		truth = val.Uint() != 0
	case reflect.Struct:
		truth = true 
	default:
		return
	}
	return truth, true
}

func (s *state) walkRange(dot reflect.Value, r *parse.RangeNode) {
	s.at(r)
	defer s.pop(s.mark())
	val, _ := indirect(s.evalPipeline(dot, r.Pipe))
	
	mark := s.mark()
	oneIteration := func(index, elem reflect.Value) {
		
		if len(r.Pipe.Decl) > 0 {
			s.setVar(1, elem)
		}
		
		if len(r.Pipe.Decl) > 1 {
			s.setVar(2, index)
		}
		s.walk(elem, r.List)
		s.pop(mark)
	}
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() == 0 {
			break
		}
		for i := 0; i < val.Len(); i++ {
			oneIteration(reflect.ValueOf(i), val.Index(i))
		}
		return
	case reflect.Map:
		if val.Len() == 0 {
			break
		}
		for _, key := range sortKeys(val.MapKeys()) {
			oneIteration(key, val.MapIndex(key))
		}
		return
	case reflect.Chan:
		if val.IsNil() {
			break
		}
		i := 0
		for ; ; i++ {
			elem, ok := val.Recv()
			if !ok {
				break
			}
			oneIteration(reflect.ValueOf(i), elem)
		}
		if i == 0 {
			break
		}
		return
	case reflect.Invalid:
		break 
	default:
		s.errorf("range can't iterate over %v", val)
	}
	if r.ElseList != nil {
		s.walk(dot, r.ElseList)
	}
}

func (s *state) walkTemplate(dot reflect.Value, t *parse.TemplateNode) {
	s.at(t)
	tmpl := s.tmpl.tmpl[t.Name]
	if tmpl == nil {
		s.errorf("template %q not defined", t.Name)
	}
	
	dot = s.evalPipeline(dot, t.Pipe)
	newState := *s
	newState.tmpl = tmpl
	
	newState.vars = []variable{{"$", dot}}
	newState.walk(dot, tmpl.Root)
}









func (s *state) evalPipeline(dot reflect.Value, pipe *parse.PipeNode) (value reflect.Value) {
	if pipe == nil {
		return
	}
	s.at(pipe)
	for _, cmd := range pipe.Cmds {
		value = s.evalCommand(dot, cmd, value) 
		
		if value.Kind() == reflect.Interface && value.Type().NumMethod() == 0 {
			value = reflect.ValueOf(value.Interface()) 
		}
	}
	for _, variable := range pipe.Decl {
		s.push(variable.Ident[0], value)
	}
	return value
}

func (s *state) notAFunction(args []parse.Node, final reflect.Value) {
	if len(args) > 1 || final.IsValid() {
		s.errorf("can't give argument to non-function %s", args[0])
	}
}

func (s *state) evalCommand(dot reflect.Value, cmd *parse.CommandNode, final reflect.Value) reflect.Value {
	firstWord := cmd.Args[0]
	switch n := firstWord.(type) {
	case *parse.FieldNode:
		return s.evalFieldNode(dot, n, cmd.Args, final)
	case *parse.ChainNode:
		return s.evalChainNode(dot, n, cmd.Args, final)
	case *parse.IdentifierNode:
		
		return s.evalFunction(dot, n, cmd, cmd.Args, final)
	case *parse.PipeNode:
		
		return s.evalPipeline(dot, n)
	case *parse.VariableNode:
		return s.evalVariableNode(dot, n, cmd.Args, final)
	}
	s.at(firstWord)
	s.notAFunction(cmd.Args, final)
	switch word := firstWord.(type) {
	case *parse.BoolNode:
		return reflect.ValueOf(word.True)
	case *parse.DotNode:
		return dot
	case *parse.NilNode:
		s.errorf("nil is not a command")
	case *parse.NumberNode:
		return s.idealConstant(word)
	case *parse.StringNode:
		return reflect.ValueOf(word.Text)
	}
	s.errorf("can't evaluate command %q", firstWord)
	panic("not reached")
}





func (s *state) idealConstant(constant *parse.NumberNode) reflect.Value {
	
	
	
	s.at(constant)
	switch {
	case constant.IsComplex:
		return reflect.ValueOf(constant.Complex128) 
	case constant.IsFloat && !isHexConstant(constant.Text) && strings.IndexAny(constant.Text, ".eE") >= 0:
		return reflect.ValueOf(constant.Float64)
	case constant.IsInt:
		n := int(constant.Int64)
		if int64(n) != constant.Int64 {
			s.errorf("%s overflows int", constant.Text)
		}
		return reflect.ValueOf(n)
	case constant.IsUint:
		s.errorf("%s overflows int", constant.Text)
	}
	return zero
}

func isHexConstant(s string) bool {
	return len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

func (s *state) evalFieldNode(dot reflect.Value, field *parse.FieldNode, args []parse.Node, final reflect.Value) reflect.Value {
	s.at(field)
	return s.evalFieldChain(dot, dot, field, field.Ident, args, final)
}

func (s *state) evalChainNode(dot reflect.Value, chain *parse.ChainNode, args []parse.Node, final reflect.Value) reflect.Value {
	s.at(chain)
	
	pipe := s.evalArg(dot, nil, chain.Node)
	if len(chain.Field) == 0 {
		s.errorf("internal error: no fields in evalChainNode")
	}
	return s.evalFieldChain(dot, pipe, chain, chain.Field, args, final)
}

func (s *state) evalVariableNode(dot reflect.Value, variable *parse.VariableNode, args []parse.Node, final reflect.Value) reflect.Value {
	
	s.at(variable)
	value := s.varValue(variable.Ident[0])
	if len(variable.Ident) == 1 {
		s.notAFunction(args, final)
		return value
	}
	return s.evalFieldChain(dot, value, variable, variable.Ident[1:], args, final)
}




func (s *state) evalFieldChain(dot, receiver reflect.Value, node parse.Node, ident []string, args []parse.Node, final reflect.Value) reflect.Value {
	n := len(ident)
	for i := 0; i < n-1; i++ {
		receiver = s.evalField(dot, ident[i], node, nil, zero, receiver)
	}
	
	return s.evalField(dot, ident[n-1], node, args, final, receiver)
}

func (s *state) evalFunction(dot reflect.Value, node *parse.IdentifierNode, cmd parse.Node, args []parse.Node, final reflect.Value) reflect.Value {
	s.at(node)
	name := node.Ident
	function, ok := findFunction(name, s.tmpl)
	if !ok {
		s.errorf("%q is not a defined function", name)
	}
	return s.evalCall(dot, function, cmd, name, args, final)
}




func (s *state) evalField(dot reflect.Value, fieldName string, node parse.Node, args []parse.Node, final, receiver reflect.Value) reflect.Value {
	if !receiver.IsValid() {
		return zero
	}
	typ := receiver.Type()
	receiver, _ = indirect(receiver)
	
	
	ptr := receiver
	if ptr.Kind() != reflect.Interface && ptr.CanAddr() {
		ptr = ptr.Addr()
	}
	if method := ptr.MethodByName(fieldName); method.IsValid() {
		return s.evalCall(dot, method, node, fieldName, args, final)
	}
	hasArgs := len(args) > 1 || final.IsValid()
	
	receiver, isNil := indirect(receiver)
	if isNil {
		s.errorf("nil pointer evaluating %s.%s", typ, fieldName)
	}
	switch receiver.Kind() {
	case reflect.Struct:
		tField, ok := receiver.Type().FieldByName(fieldName)
		if ok {
			field := receiver.FieldByIndex(tField.Index)
			if tField.PkgPath != "" { 
				s.errorf("%s is an unexported field of struct type %s", fieldName, typ)
			}
			
			if hasArgs {
				s.errorf("%s has arguments but cannot be invoked as function", fieldName)
			}
			return field
		}
		s.errorf("%s is not a field of struct type %s", fieldName, typ)
	case reflect.Map:
		
		nameVal := reflect.ValueOf(fieldName)
		if nameVal.Type().AssignableTo(receiver.Type().Key()) {
			if hasArgs {
				s.errorf("%s is not a method but has arguments", fieldName)
			}
			return receiver.MapIndex(nameVal)
		}
	}
	s.errorf("can't evaluate field %s in type %s", fieldName, typ)
	panic("not reached")
}

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)




func (s *state) evalCall(dot, fun reflect.Value, node parse.Node, name string, args []parse.Node, final reflect.Value) reflect.Value {
	if args != nil {
		args = args[1:] 
	}
	typ := fun.Type()
	numIn := len(args)
	if final.IsValid() {
		numIn++
	}
	numFixed := len(args)
	if typ.IsVariadic() {
		numFixed = typ.NumIn() - 1 
		if numIn < numFixed {
			s.errorf("wrong number of args for %s: want at least %d got %d", name, typ.NumIn()-1, len(args))
		}
	} else if numIn < typ.NumIn()-1 || !typ.IsVariadic() && numIn != typ.NumIn() {
		s.errorf("wrong number of args for %s: want %d got %d", name, typ.NumIn(), len(args))
	}
	if !goodFunc(typ) {
		
		s.errorf("can't call method/function %q with %d results", name, typ.NumOut())
	}
	
	argv := make([]reflect.Value, numIn)
	
	i := 0
	for ; i < numFixed && i < len(args); i++ {
		argv[i] = s.evalArg(dot, typ.In(i), args[i])
	}
	
	if typ.IsVariadic() {
		argType := typ.In(typ.NumIn() - 1).Elem() 
		for ; i < len(args); i++ {
			argv[i] = s.evalArg(dot, argType, args[i])
		}
	}
	
	if final.IsValid() {
		t := typ.In(typ.NumIn() - 1)
		if typ.IsVariadic() {
			t = t.Elem()
		}
		argv[i] = s.validateType(final, t)
	}
	result := fun.Call(argv)
	
	if len(result) == 2 && !result[1].IsNil() {
		s.at(node)
		s.errorf("error calling %s: %s", name, result[1].Interface().(error))
	}
	return result[0]
}


func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}


func (s *state) validateType(value reflect.Value, typ reflect.Type) reflect.Value {
	if !value.IsValid() {
		if typ == nil || canBeNil(typ) {
			
			return reflect.Zero(typ)
		}
		s.errorf("invalid value; expected %s", typ)
	}
	if typ != nil && !value.Type().AssignableTo(typ) {
		if value.Kind() == reflect.Interface && !value.IsNil() {
			value = value.Elem()
			if value.Type().AssignableTo(typ) {
				return value
			}
			
		}
		
		
		
		
		switch {
		case value.Kind() == reflect.Ptr && value.Type().Elem().AssignableTo(typ):
			value = value.Elem()
			if !value.IsValid() {
				s.errorf("dereference of nil pointer of type %s", typ)
			}
		case reflect.PtrTo(value.Type()).AssignableTo(typ) && value.CanAddr():
			value = value.Addr()
		default:
			s.errorf("wrong type for value; expected %s; got %s", typ, value.Type())
		}
	}
	return value
}

func (s *state) evalArg(dot reflect.Value, typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	switch arg := n.(type) {
	case *parse.DotNode:
		return s.validateType(dot, typ)
	case *parse.NilNode:
		if canBeNil(typ) {
			return reflect.Zero(typ)
		}
		s.errorf("cannot assign nil to %s", typ)
	case *parse.FieldNode:
		return s.validateType(s.evalFieldNode(dot, arg, []parse.Node{n}, zero), typ)
	case *parse.VariableNode:
		return s.validateType(s.evalVariableNode(dot, arg, nil, zero), typ)
	case *parse.PipeNode:
		return s.validateType(s.evalPipeline(dot, arg), typ)
	case *parse.IdentifierNode:
		return s.evalFunction(dot, arg, arg, nil, zero)
	case *parse.ChainNode:
		return s.validateType(s.evalChainNode(dot, arg, nil, zero), typ)
	}
	switch typ.Kind() {
	case reflect.Bool:
		return s.evalBool(typ, n)
	case reflect.Complex64, reflect.Complex128:
		return s.evalComplex(typ, n)
	case reflect.Float32, reflect.Float64:
		return s.evalFloat(typ, n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return s.evalInteger(typ, n)
	case reflect.Interface:
		if typ.NumMethod() == 0 {
			return s.evalEmptyInterface(dot, n)
		}
	case reflect.String:
		return s.evalString(typ, n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return s.evalUnsignedInteger(typ, n)
	}
	s.errorf("can't handle %s for arg of type %s", n, typ)
	panic("not reached")
}

func (s *state) evalBool(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.BoolNode); ok {
		value := reflect.New(typ).Elem()
		value.SetBool(n.True)
		return value
	}
	s.errorf("expected bool; found %s", n)
	panic("not reached")
}

func (s *state) evalString(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.StringNode); ok {
		value := reflect.New(typ).Elem()
		value.SetString(n.Text)
		return value
	}
	s.errorf("expected string; found %s", n)
	panic("not reached")
}

func (s *state) evalInteger(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.NumberNode); ok && n.IsInt {
		value := reflect.New(typ).Elem()
		value.SetInt(n.Int64)
		return value
	}
	s.errorf("expected integer; found %s", n)
	panic("not reached")
}

func (s *state) evalUnsignedInteger(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.NumberNode); ok && n.IsUint {
		value := reflect.New(typ).Elem()
		value.SetUint(n.Uint64)
		return value
	}
	s.errorf("expected unsigned integer; found %s", n)
	panic("not reached")
}

func (s *state) evalFloat(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.NumberNode); ok && n.IsFloat {
		value := reflect.New(typ).Elem()
		value.SetFloat(n.Float64)
		return value
	}
	s.errorf("expected float; found %s", n)
	panic("not reached")
}

func (s *state) evalComplex(typ reflect.Type, n parse.Node) reflect.Value {
	if n, ok := n.(*parse.NumberNode); ok && n.IsComplex {
		value := reflect.New(typ).Elem()
		value.SetComplex(n.Complex128)
		return value
	}
	s.errorf("expected complex; found %s", n)
	panic("not reached")
}

func (s *state) evalEmptyInterface(dot reflect.Value, n parse.Node) reflect.Value {
	s.at(n)
	switch n := n.(type) {
	case *parse.BoolNode:
		return reflect.ValueOf(n.True)
	case *parse.DotNode:
		return dot
	case *parse.FieldNode:
		return s.evalFieldNode(dot, n, nil, zero)
	case *parse.IdentifierNode:
		return s.evalFunction(dot, n, n, nil, zero)
	case *parse.NilNode:
		
		s.errorf("evalEmptyInterface: nil (can't happen)")
	case *parse.NumberNode:
		return s.idealConstant(n)
	case *parse.StringNode:
		return reflect.ValueOf(n.Text)
	case *parse.VariableNode:
		return s.evalVariableNode(dot, n, nil, zero)
	case *parse.PipeNode:
		return s.evalPipeline(dot, n)
	}
	s.errorf("can't handle assignment of %s to empty interface argument", n)
	panic("not reached")
}




func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
		if v.Kind() == reflect.Interface && v.NumMethod() > 0 {
			break
		}
	}
	return v, false
}



func (s *state) printValue(n parse.Node, v reflect.Value) {
	s.at(n)
	iface, ok := printableValue(v)
	if !ok {
		s.errorf("can't print %s of type %s", n, v.Type())
	}
	fmt.Fprint(s.wr, iface)
}



func printableValue(v reflect.Value) (interface{}, bool) {
	if v.Kind() == reflect.Ptr {
		v, _ = indirect(v) 
	}
	if !v.IsValid() {
		return "<no value>", true
	}

	if !v.Type().Implements(errorType) && !v.Type().Implements(fmtStringerType) {
		if v.CanAddr() && (reflect.PtrTo(v.Type()).Implements(errorType) || reflect.PtrTo(v.Type()).Implements(fmtStringerType)) {
			v = v.Addr()
		} else {
			switch v.Kind() {
			case reflect.Chan, reflect.Func:
				return nil, false
			}
		}
	}
	return v.Interface(), true
}



type rvs []reflect.Value

func (x rvs) Len() int      { return len(x) }
func (x rvs) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

type rvInts struct{ rvs }

func (x rvInts) Less(i, j int) bool { return x.rvs[i].Int() < x.rvs[j].Int() }

type rvUints struct{ rvs }

func (x rvUints) Less(i, j int) bool { return x.rvs[i].Uint() < x.rvs[j].Uint() }

type rvFloats struct{ rvs }

func (x rvFloats) Less(i, j int) bool { return x.rvs[i].Float() < x.rvs[j].Float() }

type rvStrings struct{ rvs }

func (x rvStrings) Less(i, j int) bool { return x.rvs[i].String() < x.rvs[j].String() }


func sortKeys(v []reflect.Value) []reflect.Value {
	if len(v) <= 1 {
		return v
	}
	switch v[0].Kind() {
	case reflect.Float32, reflect.Float64:
		sort.Sort(rvFloats{v})
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sort.Sort(rvInts{v})
	case reflect.String:
		sort.Sort(rvStrings{v})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		sort.Sort(rvUints{v})
	}
	return v
}
