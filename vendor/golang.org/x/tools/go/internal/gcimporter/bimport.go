





package gcimporter

import (
	"encoding/binary"
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

type importer struct {
	imports    map[string]*types.Package
	data       []byte
	importpath string
	buf        []byte 
	version    int    

	
	strList       []string           
	pathList      []string           
	pkgList       []*types.Package   
	typList       []types.Type       
	interfaceList []*types.Interface 
	trackAllTypes bool

	
	posInfoFormat bool
	prevFile      string
	prevLine      int
	fake          fakeFileSet

	
	debugFormat bool
	read        int 
}





func BImportData(fset *token.FileSet, imports map[string]*types.Package, data []byte, path string) (_ int, pkg *types.Package, err error) {
	
	const currentVersion = 6
	version := -1 
	defer func() {
		if e := recover(); e != nil {
			
			if version > currentVersion {
				err = fmt.Errorf("cannot import %q (%v), export data is newer version - update tool", path, e)
			} else {
				err = fmt.Errorf("cannot import %q (%v), possibly version skew - reinstall package", path, e)
			}
		}
	}()

	p := importer{
		imports:    imports,
		data:       data,
		importpath: path,
		version:    version,
		strList:    []string{""}, 
		pathList:   []string{""}, 
		fake: fakeFileSet{
			fset:  fset,
			files: make(map[string]*token.File),
		},
	}

	
	var versionstr string
	if b := p.rawByte(); b == 'c' || b == 'd' {
		
		
		
		
		
		
		if b == 'd' {
			p.debugFormat = true
		}
		p.trackAllTypes = p.rawByte() == 'a'
		p.posInfoFormat = p.int() != 0
		versionstr = p.string()
		if versionstr == "v1" {
			version = 0
		}
	} else {
		
		
		versionstr = p.rawStringln(b)
		if s := strings.SplitN(versionstr, " ", 3); len(s) >= 2 && s[0] == "version" {
			if v, err := strconv.Atoi(s[1]); err == nil && v > 0 {
				version = v
			}
		}
	}
	p.version = version

	
	switch p.version {
	
	
	
	case currentVersion, 5, 4, 3, 2, 1:
		p.debugFormat = p.rawStringln(p.rawByte()) == "debug"
		p.trackAllTypes = p.int() != 0
		p.posInfoFormat = p.int() != 0
	case 0:
		
	default:
		errorf("unknown bexport format version %d (%q)", p.version, versionstr)
	}

	

	
	p.typList = append(p.typList, predeclared...)

	
	pkg = p.pkg()

	
	objcount := 0
	for {
		tag := p.tagOrIndex()
		if tag == endTag {
			break
		}
		p.obj(tag)
		objcount++
	}

	
	if count := p.int(); count != objcount {
		errorf("got %d objects; want %d", objcount, count)
	}

	

	
	
	for _, typ := range p.interfaceList {
		typ.Complete()
	}

	
	list := append(([]*types.Package)(nil), p.pkgList[1:]...)
	sort.Sort(byPath(list))
	pkg.SetImports(list)

	
	pkg.MarkComplete()

	return p.read, pkg, nil
}

func errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func (p *importer) pkg() *types.Package {
	
	i := p.tagOrIndex()
	if i >= 0 {
		return p.pkgList[i]
	}

	
	if i != packageTag {
		errorf("unexpected package tag %d version %d", i, p.version)
	}

	
	name := p.string()
	var path string
	if p.version >= 5 {
		path = p.path()
	} else {
		path = p.string()
	}
	if p.version >= 6 {
		p.int() 
	}

	
	if name == "" {
		errorf("empty package name in import")
	}

	
	
	if (path == "") != (len(p.pkgList) == 0) {
		errorf("package path %q for pkg index %d", path, len(p.pkgList))
	}

	
	if path == "" {
		path = p.importpath
	}
	pkg := p.imports[path]
	if pkg == nil {
		pkg = types.NewPackage(path, name)
		p.imports[path] = pkg
	} else if pkg.Name() != name {
		errorf("conflicting names %s and %s for package %q", pkg.Name(), name, path)
	}
	p.pkgList = append(p.pkgList, pkg)

	return pkg
}


func objTag(obj types.Object) int {
	switch obj.(type) {
	case *types.Const:
		return constTag
	case *types.TypeName:
		return typeTag
	case *types.Var:
		return varTag
	case *types.Func:
		return funcTag
	default:
		errorf("unexpected object: %v (%T)", obj, obj) 
		panic("unreachable")
	}
}

func sameObj(a, b types.Object) bool {
	
	
	
	return objTag(a) == objTag(b) && types.Identical(a.Type(), b.Type())
}

func (p *importer) declare(obj types.Object) {
	pkg := obj.Pkg()
	if alt := pkg.Scope().Insert(obj); alt != nil {
		
		
		
		
		
		
		
		
		
		if !sameObj(obj, alt) {
			errorf("inconsistent import:\n\t%v\npreviously imported as:\n\t%v\n", obj, alt)
		}
	}
}

func (p *importer) obj(tag int) {
	switch tag {
	case constTag:
		pos := p.pos()
		pkg, name := p.qualifiedName()
		typ := p.typ(nil, nil)
		val := p.value()
		p.declare(types.NewConst(pos, pkg, name, typ, val))

	case aliasTag:
		
		pos := p.pos()
		pkg, name := p.qualifiedName()
		typ := p.typ(nil, nil)
		p.declare(types.NewTypeName(pos, pkg, name, typ))

	case typeTag:
		p.typ(nil, nil)

	case varTag:
		pos := p.pos()
		pkg, name := p.qualifiedName()
		typ := p.typ(nil, nil)
		p.declare(types.NewVar(pos, pkg, name, typ))

	case funcTag:
		pos := p.pos()
		pkg, name := p.qualifiedName()
		params, isddd := p.paramList()
		result, _ := p.paramList()
		sig := types.NewSignature(nil, params, result, isddd)
		p.declare(types.NewFunc(pos, pkg, name, sig))

	default:
		errorf("unexpected object tag %d", tag)
	}
}

const deltaNewFile = -64 

func (p *importer) pos() token.Pos {
	if !p.posInfoFormat {
		return token.NoPos
	}

	file := p.prevFile
	line := p.prevLine
	delta := p.int()
	line += delta
	if p.version >= 5 {
		if delta == deltaNewFile {
			if n := p.int(); n >= 0 {
				
				file = p.path()
				line = n
			}
		}
	} else {
		if delta == 0 {
			if n := p.int(); n >= 0 {
				
				file = p.prevFile[:n] + p.string()
				line = p.int()
			}
		}
	}
	p.prevFile = file
	p.prevLine = line

	return p.fake.pos(file, line)
}


type fakeFileSet struct {
	fset  *token.FileSet
	files map[string]*token.File
}

func (s *fakeFileSet) pos(file string, line int) token.Pos {
	
	
	const maxlines = 64 * 1024
	f := s.files[file]
	if f == nil {
		f = s.fset.AddFile(file, -1, maxlines)
		s.files[file] = f
		
		
		fakeLinesOnce.Do(func() {
			fakeLines = make([]int, maxlines)
			for i := range fakeLines {
				fakeLines[i] = i
			}
		})
		f.SetLines(fakeLines)
	}

	if line > maxlines {
		line = 1
	}

	
	
	return f.Pos(line - 1)
}

var (
	fakeLines     []int
	fakeLinesOnce sync.Once
)

func (p *importer) qualifiedName() (pkg *types.Package, name string) {
	name = p.string()
	pkg = p.pkg()
	return
}

func (p *importer) record(t types.Type) {
	p.typList = append(p.typList, t)
}




type dddSlice struct {
	elem types.Type
}

func (t *dddSlice) Underlying() types.Type { return t }
func (t *dddSlice) String() string         { return "..." + t.elem.String() }









func (p *importer) typ(parent *types.Package, tname *types.Named) types.Type {
	
	i := p.tagOrIndex()
	if i >= 0 {
		return p.typList[i]
	}

	
	switch i {
	case namedTag:
		
		pos := p.pos()
		parent, name := p.qualifiedName()
		scope := parent.Scope()
		obj := scope.Lookup(name)

		
		if obj == nil {
			obj = types.NewTypeName(pos, parent, name, nil)
			scope.Insert(obj)
		}

		if _, ok := obj.(*types.TypeName); !ok {
			errorf("pkg = %s, name = %s => %s", parent, name, obj)
		}

		
		t0 := types.NewNamed(obj.(*types.TypeName), nil, nil)

		
		tname := obj.Type().(*types.Named) 
		p.record(tname)

		
		t0.SetUnderlying(p.typ(parent, t0))

		
		if types.IsInterface(t0) {
			return tname
		}

		
		for i := p.int(); i > 0; i-- {
			
			pos := p.pos()
			name := p.string()
			if !exported(name) {
				p.pkg()
			}

			recv, _ := p.paramList() 
			params, isddd := p.paramList()
			result, _ := p.paramList()
			p.int() 

			sig := types.NewSignature(recv.At(0), params, result, isddd)
			t0.AddMethod(types.NewFunc(pos, parent, name, sig))
		}

		return tname

	case arrayTag:
		t := new(types.Array)
		if p.trackAllTypes {
			p.record(t)
		}

		n := p.int64()
		*t = *types.NewArray(p.typ(parent, nil), n)
		return t

	case sliceTag:
		t := new(types.Slice)
		if p.trackAllTypes {
			p.record(t)
		}

		*t = *types.NewSlice(p.typ(parent, nil))
		return t

	case dddTag:
		t := new(dddSlice)
		if p.trackAllTypes {
			p.record(t)
		}

		t.elem = p.typ(parent, nil)
		return t

	case structTag:
		t := new(types.Struct)
		if p.trackAllTypes {
			p.record(t)
		}

		*t = *types.NewStruct(p.fieldList(parent))
		return t

	case pointerTag:
		t := new(types.Pointer)
		if p.trackAllTypes {
			p.record(t)
		}

		*t = *types.NewPointer(p.typ(parent, nil))
		return t

	case signatureTag:
		t := new(types.Signature)
		if p.trackAllTypes {
			p.record(t)
		}

		params, isddd := p.paramList()
		result, _ := p.paramList()
		*t = *types.NewSignature(nil, params, result, isddd)
		return t

	case interfaceTag:
		
		
		
		
		
		
		n := len(p.typList)
		if p.trackAllTypes {
			p.record(nil)
		}

		var embeddeds []types.Type
		for n := p.int(); n > 0; n-- {
			p.pos()
			embeddeds = append(embeddeds, p.typ(parent, nil))
		}

		t := newInterface(p.methodList(parent, tname), embeddeds)
		p.interfaceList = append(p.interfaceList, t)
		if p.trackAllTypes {
			p.typList[n] = t
		}
		return t

	case mapTag:
		t := new(types.Map)
		if p.trackAllTypes {
			p.record(t)
		}

		key := p.typ(parent, nil)
		val := p.typ(parent, nil)
		*t = *types.NewMap(key, val)
		return t

	case chanTag:
		t := new(types.Chan)
		if p.trackAllTypes {
			p.record(t)
		}

		dir := chanDir(p.int())
		val := p.typ(parent, nil)
		*t = *types.NewChan(dir, val)
		return t

	default:
		errorf("unexpected type tag %d", i) 
		panic("unreachable")
	}
}

func chanDir(d int) types.ChanDir {
	
	switch d {
	case 1  :
		return types.RecvOnly
	case 2  :
		return types.SendOnly
	case 3  :
		return types.SendRecv
	default:
		errorf("unexpected channel dir %d", d)
		return 0
	}
}

func (p *importer) fieldList(parent *types.Package) (fields []*types.Var, tags []string) {
	if n := p.int(); n > 0 {
		fields = make([]*types.Var, n)
		tags = make([]string, n)
		for i := range fields {
			fields[i], tags[i] = p.field(parent)
		}
	}
	return
}

func (p *importer) field(parent *types.Package) (*types.Var, string) {
	pos := p.pos()
	pkg, name, alias := p.fieldName(parent)
	typ := p.typ(parent, nil)
	tag := p.string()

	anonymous := false
	if name == "" {
		
		switch typ := deref(typ).(type) {
		case *types.Basic: 
			pkg = nil 
			name = typ.Name()
		case *types.Named:
			name = typ.Obj().Name()
		default:
			errorf("named base type expected")
		}
		anonymous = true
	} else if alias {
		
		anonymous = true
	}

	return types.NewField(pos, pkg, name, typ, anonymous), tag
}

func (p *importer) methodList(parent *types.Package, baseType *types.Named) (methods []*types.Func) {
	if n := p.int(); n > 0 {
		methods = make([]*types.Func, n)
		for i := range methods {
			methods[i] = p.method(parent, baseType)
		}
	}
	return
}

func (p *importer) method(parent *types.Package, baseType *types.Named) *types.Func {
	pos := p.pos()
	pkg, name, _ := p.fieldName(parent)
	
	
	
	
	var recv *types.Var
	if baseType != nil {
		recv = types.NewVar(token.NoPos, parent, "", baseType)
	}
	params, isddd := p.paramList()
	result, _ := p.paramList()
	sig := types.NewSignature(recv, params, result, isddd)
	return types.NewFunc(pos, pkg, name, sig)
}

func (p *importer) fieldName(parent *types.Package) (pkg *types.Package, name string, alias bool) {
	name = p.string()
	pkg = parent
	if pkg == nil {
		
		pkg = p.pkgList[0]
	}
	if p.version == 0 && name == "_" {
		
		return
	}
	switch name {
	case "":
		
	case "?":
		
		name = ""
		pkg = p.pkg()
	case "@":
		
		name = p.string()
		alias = true
		fallthrough
	default:
		if !exported(name) {
			pkg = p.pkg()
		}
	}
	return
}

func (p *importer) paramList() (*types.Tuple, bool) {
	n := p.int()
	if n == 0 {
		return nil, false
	}
	
	named := true
	if n < 0 {
		n = -n
		named = false
	}
	
	params := make([]*types.Var, n)
	isddd := false
	for i := range params {
		params[i], isddd = p.param(named)
	}
	return types.NewTuple(params...), isddd
}

func (p *importer) param(named bool) (*types.Var, bool) {
	t := p.typ(nil, nil)
	td, isddd := t.(*dddSlice)
	if isddd {
		t = types.NewSlice(td.elem)
	}

	var pkg *types.Package
	var name string
	if named {
		name = p.string()
		if name == "" {
			errorf("expected named parameter")
		}
		if name != "_" {
			pkg = p.pkg()
		}
		if i := strings.Index(name, "·"); i > 0 {
			name = name[:i] 
		}
	}

	
	p.string()

	return types.NewVar(token.NoPos, pkg, name, t), isddd
}

func exported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

func (p *importer) value() constant.Value {
	switch tag := p.tagOrIndex(); tag {
	case falseTag:
		return constant.MakeBool(false)
	case trueTag:
		return constant.MakeBool(true)
	case int64Tag:
		return constant.MakeInt64(p.int64())
	case floatTag:
		return p.float()
	case complexTag:
		re := p.float()
		im := p.float()
		return constant.BinaryOp(re, token.ADD, constant.MakeImag(im))
	case stringTag:
		return constant.MakeString(p.string())
	case unknownTag:
		return constant.MakeUnknown()
	default:
		errorf("unexpected value tag %d", tag) 
		panic("unreachable")
	}
}

func (p *importer) float() constant.Value {
	sign := p.int()
	if sign == 0 {
		return constant.MakeInt64(0)
	}

	exp := p.int()
	mant := []byte(p.string()) 

	
	for len(mant) > 0 && mant[0] == 0 {
		mant = mant[1:]
	}

	
	
	
	for i, j := 0, len(mant)-1; i < j; i, j = i+1, j-1 {
		mant[i], mant[j] = mant[j], mant[i]
	}

	
	
	exp -= len(mant) << 3
	if len(mant) > 0 {
		for msd := mant[len(mant)-1]; msd&0x80 == 0; msd <<= 1 {
			exp++
		}
	}

	x := constant.MakeFromBytes(mant)
	switch {
	case exp < 0:
		d := constant.Shift(constant.MakeInt64(1), token.SHL, uint(-exp))
		x = constant.BinaryOp(x, token.QUO, d)
	case exp > 0:
		x = constant.Shift(x, token.SHL, uint(exp))
	}

	if sign < 0 {
		x = constant.UnaryOp(token.SUB, x, 0)
	}
	return x
}




func (p *importer) tagOrIndex() int {
	if p.debugFormat {
		p.marker('t')
	}

	return int(p.rawInt64())
}

func (p *importer) int() int {
	x := p.int64()
	if int64(int(x)) != x {
		errorf("exported integer too large")
	}
	return int(x)
}

func (p *importer) int64() int64 {
	if p.debugFormat {
		p.marker('i')
	}

	return p.rawInt64()
}

func (p *importer) path() string {
	if p.debugFormat {
		p.marker('p')
	}
	
	
	i := p.rawInt64()
	if i >= 0 {
		return p.pathList[i]
	}
	
	a := make([]string, -i)
	for n := range a {
		a[n] = p.string()
	}
	s := strings.Join(a, "/")
	p.pathList = append(p.pathList, s)
	return s
}

func (p *importer) string() string {
	if p.debugFormat {
		p.marker('s')
	}
	
	
	i := p.rawInt64()
	if i >= 0 {
		return p.strList[i]
	}
	
	if n := int(-i); n <= cap(p.buf) {
		p.buf = p.buf[:n]
	} else {
		p.buf = make([]byte, n)
	}
	for i := range p.buf {
		p.buf[i] = p.rawByte()
	}
	s := string(p.buf)
	p.strList = append(p.strList, s)
	return s
}

func (p *importer) marker(want byte) {
	if got := p.rawByte(); got != want {
		errorf("incorrect marker: got %c; want %c (pos = %d)", got, want, p.read)
	}

	pos := p.read
	if n := int(p.rawInt64()); n != pos {
		errorf("incorrect position: got %d; want %d", n, pos)
	}
}


func (p *importer) rawInt64() int64 {
	i, err := binary.ReadVarint(p)
	if err != nil {
		errorf("read error: %v", err)
	}
	return i
}


func (p *importer) rawStringln(b byte) string {
	p.buf = p.buf[:0]
	for b != '\n' {
		p.buf = append(p.buf, b)
		b = p.rawByte()
	}
	return string(p.buf)
}


func (p *importer) ReadByte() (byte, error) {
	return p.rawByte(), nil
}




func (p *importer) rawByte() byte {
	b := p.data[0]
	r := 1
	if b == '|' {
		b = p.data[1]
		r = 2
		switch b {
		case 'S':
			b = '$'
		case '|':
			
		default:
			errorf("unexpected escape sequence in export data")
		}
	}
	p.data = p.data[r:]
	p.read += r
	return b

}





const (
	
	packageTag = -(iota + 1)
	constTag
	typeTag
	varTag
	funcTag
	endTag

	
	namedTag
	arrayTag
	sliceTag
	dddTag
	structTag
	pointerTag
	signatureTag
	interfaceTag
	mapTag
	chanTag

	
	falseTag
	trueTag
	int64Tag
	floatTag
	fractionTag 
	complexTag
	stringTag
	nilTag     
	unknownTag 

	
	aliasTag
)

var predeclared = []types.Type{
	
	types.Typ[types.Bool],
	types.Typ[types.Int],
	types.Typ[types.Int8],
	types.Typ[types.Int16],
	types.Typ[types.Int32],
	types.Typ[types.Int64],
	types.Typ[types.Uint],
	types.Typ[types.Uint8],
	types.Typ[types.Uint16],
	types.Typ[types.Uint32],
	types.Typ[types.Uint64],
	types.Typ[types.Uintptr],
	types.Typ[types.Float32],
	types.Typ[types.Float64],
	types.Typ[types.Complex64],
	types.Typ[types.Complex128],
	types.Typ[types.String],

	
	types.Universe.Lookup("byte").Type(),
	types.Universe.Lookup("rune").Type(),

	
	types.Universe.Lookup("error").Type(),

	
	types.Typ[types.UntypedBool],
	types.Typ[types.UntypedInt],
	types.Typ[types.UntypedRune],
	types.Typ[types.UntypedFloat],
	types.Typ[types.UntypedComplex],
	types.Typ[types.UntypedString],
	types.Typ[types.UntypedNil],

	
	types.Typ[types.UnsafePointer],

	
	types.Typ[types.Invalid], 

	
	anyType{},
}

type anyType struct{}

func (t anyType) Underlying() types.Type { return t }
func (t anyType) String() string         { return "any" }
