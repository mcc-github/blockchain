







package gcimporter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"math"
	"math/big"
	"sort"
	"strings"
)











const debugFormat = false 


const trace = false 












const exportVersion = 4










const trackAllTypes = false

type exporter struct {
	fset *token.FileSet
	out  bytes.Buffer

	
	strIndex map[string]int
	pkgIndex map[*types.Package]int
	typIndex map[types.Type]int

	
	posInfoFormat bool
	prevFile      string
	prevLine      int

	
	written int 
	indent  int 
}


type internalError string

func (e internalError) Error() string { return "gcimporter: " + string(e) }

func internalErrorf(format string, args ...interface{}) error {
	return internalError(fmt.Sprintf(format, args...))
}



func BExportData(fset *token.FileSet, pkg *types.Package) (b []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			if ierr, ok := e.(internalError); ok {
				err = ierr
				return
			}
			
			panic(e)
		}
	}()

	p := exporter{
		fset:          fset,
		strIndex:      map[string]int{"": 0}, 
		pkgIndex:      make(map[*types.Package]int),
		typIndex:      make(map[types.Type]int),
		posInfoFormat: true, 
	}

	
	
	
	
	p.rawStringln(fmt.Sprintf("version %d", exportVersion))
	var debug string
	if debugFormat {
		debug = "debug"
	}
	p.rawStringln(debug) 
	p.bool(trackAllTypes)
	p.bool(p.posInfoFormat)

	

	
	for index, typ := range predeclared {
		p.typIndex[typ] = index
	}
	if len(p.typIndex) != len(predeclared) {
		return nil, internalError("duplicate entries in type map?")
	}

	
	p.pkg(pkg, true)
	if trace {
		p.tracef("\n")
	}

	
	objcount := 0
	scope := pkg.Scope()
	for _, name := range scope.Names() {
		if !ast.IsExported(name) {
			continue
		}
		if trace {
			p.tracef("\n")
		}
		p.obj(scope.Lookup(name))
		objcount++
	}

	
	if trace {
		p.tracef("\n")
	}
	p.tag(endTag)

	
	p.int(objcount)

	if trace {
		p.tracef("\n")
	}

	

	return p.out.Bytes(), nil
}

func (p *exporter) pkg(pkg *types.Package, emptypath bool) {
	if pkg == nil {
		panic(internalError("unexpected nil pkg"))
	}

	
	if i, ok := p.pkgIndex[pkg]; ok {
		p.index('P', i)
		return
	}

	
	if trace {
		p.tracef("P%d = { ", len(p.pkgIndex))
		defer p.tracef("} ")
	}
	p.pkgIndex[pkg] = len(p.pkgIndex)

	p.tag(packageTag)
	p.string(pkg.Name())
	if emptypath {
		p.string("")
	} else {
		p.string(pkg.Path())
	}
}

func (p *exporter) obj(obj types.Object) {
	switch obj := obj.(type) {
	case *types.Const:
		p.tag(constTag)
		p.pos(obj)
		p.qualifiedName(obj)
		p.typ(obj.Type())
		p.value(obj.Val())

	case *types.TypeName:
		if obj.IsAlias() {
			p.tag(aliasTag)
			p.pos(obj)
			p.qualifiedName(obj)
		} else {
			p.tag(typeTag)
		}
		p.typ(obj.Type())

	case *types.Var:
		p.tag(varTag)
		p.pos(obj)
		p.qualifiedName(obj)
		p.typ(obj.Type())

	case *types.Func:
		p.tag(funcTag)
		p.pos(obj)
		p.qualifiedName(obj)
		sig := obj.Type().(*types.Signature)
		p.paramList(sig.Params(), sig.Variadic())
		p.paramList(sig.Results(), false)

	default:
		panic(internalErrorf("unexpected object %v (%T)", obj, obj))
	}
}

func (p *exporter) pos(obj types.Object) {
	if !p.posInfoFormat {
		return
	}

	file, line := p.fileLine(obj)
	if file == p.prevFile {
		
		
		delta := line - p.prevLine
		p.int(delta)
		if delta == 0 {
			p.int(-1) 
		}
	} else {
		
		p.int(0)
		
		
		
		
		
		
		n := commonPrefixLen(p.prevFile, file)
		p.int(n)           
		p.string(file[n:]) 
		p.prevFile = file
		p.int(line)
	}
	p.prevLine = line
}

func (p *exporter) fileLine(obj types.Object) (file string, line int) {
	if p.fset != nil {
		pos := p.fset.Position(obj.Pos())
		file = pos.Filename
		line = pos.Line
	}
	return
}

func commonPrefixLen(a, b string) int {
	if len(a) > len(b) {
		a, b = b, a
	}
	
	i := 0
	for i < len(a) && a[i] == b[i] {
		i++
	}
	return i
}

func (p *exporter) qualifiedName(obj types.Object) {
	p.string(obj.Name())
	p.pkg(obj.Pkg(), false)
}

func (p *exporter) typ(t types.Type) {
	if t == nil {
		panic(internalError("nil type"))
	}

	
	
	
	
	
	
	
	

	
	if i, ok := p.typIndex[t]; ok {
		p.index('T', i)
		return
	}

	
	if trackAllTypes {
		if trace {
			p.tracef("T%d = {>\n", len(p.typIndex))
			defer p.tracef("<\n} ")
		}
		p.typIndex[t] = len(p.typIndex)
	}

	switch t := t.(type) {
	case *types.Named:
		if !trackAllTypes {
			
			p.typIndex[t] = len(p.typIndex)
		}

		p.tag(namedTag)
		p.pos(t.Obj())
		p.qualifiedName(t.Obj())
		p.typ(t.Underlying())
		if !types.IsInterface(t) {
			p.assocMethods(t)
		}

	case *types.Array:
		p.tag(arrayTag)
		p.int64(t.Len())
		p.typ(t.Elem())

	case *types.Slice:
		p.tag(sliceTag)
		p.typ(t.Elem())

	case *dddSlice:
		p.tag(dddTag)
		p.typ(t.elem)

	case *types.Struct:
		p.tag(structTag)
		p.fieldList(t)

	case *types.Pointer:
		p.tag(pointerTag)
		p.typ(t.Elem())

	case *types.Signature:
		p.tag(signatureTag)
		p.paramList(t.Params(), t.Variadic())
		p.paramList(t.Results(), false)

	case *types.Interface:
		p.tag(interfaceTag)
		p.iface(t)

	case *types.Map:
		p.tag(mapTag)
		p.typ(t.Key())
		p.typ(t.Elem())

	case *types.Chan:
		p.tag(chanTag)
		p.int(int(3 - t.Dir())) 
		p.typ(t.Elem())

	default:
		panic(internalErrorf("unexpected type %T: %s", t, t))
	}
}

func (p *exporter) assocMethods(named *types.Named) {
	
	var methods []*types.Func
	for i := 0; i < named.NumMethods(); i++ {
		methods = append(methods, named.Method(i))
	}
	sort.Sort(methodsByName(methods))

	p.int(len(methods))

	if trace && methods != nil {
		p.tracef("associated methods {>\n")
	}

	for i, m := range methods {
		if trace && i > 0 {
			p.tracef("\n")
		}

		p.pos(m)
		name := m.Name()
		p.string(name)
		if !exported(name) {
			p.pkg(m.Pkg(), false)
		}

		sig := m.Type().(*types.Signature)
		p.paramList(types.NewTuple(sig.Recv()), false)
		p.paramList(sig.Params(), sig.Variadic())
		p.paramList(sig.Results(), false)
		p.int(0) 
	}

	if trace && methods != nil {
		p.tracef("<\n} ")
	}
}

type methodsByName []*types.Func

func (x methodsByName) Len() int           { return len(x) }
func (x methodsByName) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x methodsByName) Less(i, j int) bool { return x[i].Name() < x[j].Name() }

func (p *exporter) fieldList(t *types.Struct) {
	if trace && t.NumFields() > 0 {
		p.tracef("fields {>\n")
		defer p.tracef("<\n} ")
	}

	p.int(t.NumFields())
	for i := 0; i < t.NumFields(); i++ {
		if trace && i > 0 {
			p.tracef("\n")
		}
		p.field(t.Field(i))
		p.string(t.Tag(i))
	}
}

func (p *exporter) field(f *types.Var) {
	if !f.IsField() {
		panic(internalError("field expected"))
	}

	p.pos(f)
	p.fieldName(f)
	p.typ(f.Type())
}

func (p *exporter) iface(t *types.Interface) {
	
	
	p.int(0)

	n := t.NumMethods()
	if trace && n > 0 {
		p.tracef("methods {>\n")
		defer p.tracef("<\n} ")
	}
	p.int(n)
	for i := 0; i < n; i++ {
		if trace && i > 0 {
			p.tracef("\n")
		}
		p.method(t.Method(i))
	}
}

func (p *exporter) method(m *types.Func) {
	sig := m.Type().(*types.Signature)
	if sig.Recv() == nil {
		panic(internalError("method expected"))
	}

	p.pos(m)
	p.string(m.Name())
	if m.Name() != "_" && !ast.IsExported(m.Name()) {
		p.pkg(m.Pkg(), false)
	}

	
	p.paramList(sig.Params(), sig.Variadic())
	p.paramList(sig.Results(), false)
}

func (p *exporter) fieldName(f *types.Var) {
	name := f.Name()

	if f.Anonymous() {
		
		
		
		
		bname := basetypeName(f.Type())
		if name == bname {
			if ast.IsExported(name) {
				name = "" 
			} else {
				name = "?" 
			}
		} else {
			
			
			p.string("@")
		}
	}

	p.string(name)
	if name != "" && !ast.IsExported(name) {
		p.pkg(f.Pkg(), false)
	}
}

func basetypeName(typ types.Type) string {
	switch typ := deref(typ).(type) {
	case *types.Basic:
		return typ.Name()
	case *types.Named:
		return typ.Obj().Name()
	default:
		return "" 
	}
}

func (p *exporter) paramList(params *types.Tuple, variadic bool) {
	
	
	
	n := params.Len()
	if n > 0 && params.At(0).Name() == "" {
		n = -n
	}
	p.int(n)
	for i := 0; i < params.Len(); i++ {
		q := params.At(i)
		t := q.Type()
		if variadic && i == params.Len()-1 {
			t = &dddSlice{t.(*types.Slice).Elem()}
		}
		p.typ(t)
		if n > 0 {
			name := q.Name()
			p.string(name)
			if name != "_" {
				p.pkg(q.Pkg(), false)
			}
		}
		p.string("") 
	}
}

func (p *exporter) value(x constant.Value) {
	if trace {
		p.tracef("= ")
	}

	switch x.Kind() {
	case constant.Bool:
		tag := falseTag
		if constant.BoolVal(x) {
			tag = trueTag
		}
		p.tag(tag)

	case constant.Int:
		if v, exact := constant.Int64Val(x); exact {
			
			p.tag(int64Tag)
			p.int64(v)
			return
		}
		
		
		p.tag(floatTag)
		p.float(constant.ToFloat(x))

	case constant.Float:
		p.tag(floatTag)
		p.float(x)

	case constant.Complex:
		p.tag(complexTag)
		p.float(constant.Real(x))
		p.float(constant.Imag(x))

	case constant.String:
		p.tag(stringTag)
		p.string(constant.StringVal(x))

	case constant.Unknown:
		
		p.tag(unknownTag)

	default:
		panic(internalErrorf("unexpected value %v (%T)", x, x))
	}
}

func (p *exporter) float(x constant.Value) {
	if x.Kind() != constant.Float {
		panic(internalErrorf("unexpected constant %v, want float", x))
	}
	
	sign := constant.Sign(x)
	if sign == 0 {
		
		p.int(0)
		return
	}
	

	var f big.Float
	if v, exact := constant.Float64Val(x); exact {
		
		f.SetFloat64(v)
	} else if num, denom := constant.Num(x), constant.Denom(x); num.Kind() == constant.Int {
		
		r := valueToRat(num)
		f.SetRat(r.Quo(r, valueToRat(denom)))
	} else {
		
		
		f.SetFloat64(math.MaxFloat64) 
	}

	
	var m big.Float
	exp := f.MantExp(&m)

	
	
	
	m.SetMantExp(&m, int(m.MinPrec()))
	mant, acc := m.Int(nil)
	if acc != big.Exact {
		panic(internalError("internal error"))
	}

	p.int(sign)
	p.int(exp)
	p.string(string(mant.Bytes()))
}

func valueToRat(x constant.Value) *big.Rat {
	
	
	bytes := constant.Bytes(x)
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i]
	}
	return new(big.Rat).SetInt(new(big.Int).SetBytes(bytes))
}

func (p *exporter) bool(b bool) bool {
	if trace {
		p.tracef("[")
		defer p.tracef("= %v] ", b)
	}

	x := 0
	if b {
		x = 1
	}
	p.int(x)
	return b
}




func (p *exporter) index(marker byte, index int) {
	if index < 0 {
		panic(internalError("invalid index < 0"))
	}
	if debugFormat {
		p.marker('t')
	}
	if trace {
		p.tracef("%c%d ", marker, index)
	}
	p.rawInt64(int64(index))
}

func (p *exporter) tag(tag int) {
	if tag >= 0 {
		panic(internalError("invalid tag >= 0"))
	}
	if debugFormat {
		p.marker('t')
	}
	if trace {
		p.tracef("%s ", tagString[-tag])
	}
	p.rawInt64(int64(tag))
}

func (p *exporter) int(x int) {
	p.int64(int64(x))
}

func (p *exporter) int64(x int64) {
	if debugFormat {
		p.marker('i')
	}
	if trace {
		p.tracef("%d ", x)
	}
	p.rawInt64(x)
}

func (p *exporter) string(s string) {
	if debugFormat {
		p.marker('s')
	}
	if trace {
		p.tracef("%q ", s)
	}
	
	
	if i, ok := p.strIndex[s]; ok {
		p.rawInt64(int64(i))
		return
	}
	
	p.strIndex[s] = len(p.strIndex)
	p.rawInt64(-int64(len(s)))
	for i := 0; i < len(s); i++ {
		p.rawByte(s[i])
	}
}




func (p *exporter) marker(m byte) {
	p.rawByte(m)
	
	
	if false && trace {
		p.tracef("#%d ", p.written)
	}
	p.rawInt64(int64(p.written))
}


func (p *exporter) rawInt64(x int64) {
	var tmp [binary.MaxVarintLen64]byte
	n := binary.PutVarint(tmp[:], x)
	for i := 0; i < n; i++ {
		p.rawByte(tmp[i])
	}
}


func (p *exporter) rawStringln(s string) {
	for i := 0; i < len(s); i++ {
		p.rawByte(s[i])
	}
	p.rawByte('\n')
}











func (p *exporter) rawByte(b byte) {
	switch b {
	case '$':
		
		b = 'S'
		fallthrough
	case '|':
		
		p.out.WriteByte('|')
		p.written++
	}
	p.out.WriteByte(b)
	p.written++
}



func (p *exporter) tracef(format string, args ...interface{}) {
	if strings.ContainsAny(format, "<>\n") {
		var buf bytes.Buffer
		for i := 0; i < len(format); i++ {
			
			ch := format[i]
			switch ch {
			case '>':
				p.indent++
				continue
			case '<':
				p.indent--
				continue
			}
			buf.WriteByte(ch)
			if ch == '\n' {
				for j := p.indent; j > 0; j-- {
					buf.WriteString(".  ")
				}
			}
		}
		format = buf.String()
	}
	fmt.Printf(format, args...)
}



var tagString = [...]string{
	
	-packageTag: "package",

	
	-namedTag:     "named type",
	-arrayTag:     "array",
	-sliceTag:     "slice",
	-dddTag:       "ddd",
	-structTag:    "struct",
	-pointerTag:   "pointer",
	-signatureTag: "signature",
	-interfaceTag: "interface",
	-mapTag:       "map",
	-chanTag:      "chan",

	
	-falseTag:    "false",
	-trueTag:     "true",
	-int64Tag:    "int64",
	-floatTag:    "float",
	-fractionTag: "fraction",
	-complexTag:  "complex",
	-stringTag:   "string",
	-unknownTag:  "unknown",

	
	-aliasTag: "alias",
}
