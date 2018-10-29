










package gcimporter 

import (
	"bufio"
	"errors"
	"fmt"
	"go/build"
	"go/constant"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/scanner"
)


const debug = false

var pkgExts = [...]string{".a", ".o"}







func FindPkg(path, srcDir string) (filename, id string) {
	if path == "" {
		return
	}

	var noext string
	switch {
	default:
		
		
		if abs, err := filepath.Abs(srcDir); err == nil { 
			srcDir = abs
		}
		bp, _ := build.Import(path, srcDir, build.FindOnly|build.AllowBinary)
		if bp.PkgObj == "" {
			id = path 
			return
		}
		noext = strings.TrimSuffix(bp.PkgObj, ".a")
		id = bp.ImportPath

	case build.IsLocalImport(path):
		
		noext = filepath.Join(srcDir, path)
		id = noext

	case filepath.IsAbs(path):
		
		
		
		noext = path
		id = path
	}

	if false { 
		if path != id {
			fmt.Printf("%s -> %s\n", path, id)
		}
	}

	
	for _, ext := range pkgExts {
		filename = noext + ext
		if f, err := os.Stat(filename); err == nil && !f.IsDir() {
			return
		}
	}

	filename = "" 
	return
}













func ImportData(packages map[string]*types.Package, filename, id string, data io.Reader) (pkg *types.Package, err error) {
	
	defer func() {
		switch r := recover().(type) {
		case nil:
			
		case importError:
			err = r
		default:
			panic(r) 
		}
	}()

	var p parser
	p.init(filename, id, data, packages)
	pkg = p.parseExport()

	return
}





func Import(packages map[string]*types.Package, path, srcDir string, lookup func(path string) (io.ReadCloser, error)) (pkg *types.Package, err error) {
	var rc io.ReadCloser
	var filename, id string
	if lookup != nil {
		
		
		if path == "unsafe" {
			return types.Unsafe, nil
		}
		id = path

		
		if pkg = packages[id]; pkg != nil && pkg.Complete() {
			return
		}
		f, err := lookup(path)
		if err != nil {
			return nil, err
		}
		rc = f
	} else {
		filename, id = FindPkg(path, srcDir)
		if filename == "" {
			if path == "unsafe" {
				return types.Unsafe, nil
			}
			return nil, fmt.Errorf("can't find import: %q", id)
		}

		
		if pkg = packages[id]; pkg != nil && pkg.Complete() {
			return
		}

		
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err != nil {
				
				err = fmt.Errorf("%s: %v", filename, err)
			}
		}()
		rc = f
	}
	defer rc.Close()

	var hdr string
	buf := bufio.NewReader(rc)
	if hdr, err = FindExportData(buf); err != nil {
		return
	}

	switch hdr {
	case "$$\n":
		
		
		
		if filename == "" {
			filename = path
		}
		return ImportData(packages, filename, id, buf)

	case "$$B\n":
		var data []byte
		data, err = ioutil.ReadAll(buf)
		if err != nil {
			break
		}

		
		
		fset := token.NewFileSet()

		
		
		
		if len(data) > 0 && data[0] == 'i' {
			_, pkg, err = IImportData(fset, packages, data[1:], id)
		} else {
			_, pkg, err = BImportData(fset, packages, data, id)
		}

	default:
		err = fmt.Errorf("unknown export data header: %q", hdr)
	}

	return
}












type parser struct {
	scanner    scanner.Scanner
	tok        rune                      
	lit        string                    
	id         string                    
	sharedPkgs map[string]*types.Package 
	localPkgs  map[string]*types.Package 
}

func (p *parser) init(filename, id string, src io.Reader, packages map[string]*types.Package) {
	p.scanner.Init(src)
	p.scanner.Error = func(_ *scanner.Scanner, msg string) { p.error(msg) }
	p.scanner.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanChars | scanner.ScanStrings | scanner.ScanComments | scanner.SkipComments
	p.scanner.Whitespace = 1<<'\t' | 1<<' '
	p.scanner.Filename = filename 
	p.next()
	p.id = id
	p.sharedPkgs = packages
	if debug {
		
		for _, pkg := range packages {
			if pkg.Name() == "" {
				fmt.Printf("no package name for %s\n", pkg.Path())
			}
		}
	}
}

func (p *parser) next() {
	p.tok = p.scanner.Scan()
	switch p.tok {
	case scanner.Ident, scanner.Int, scanner.Char, scanner.String, '·':
		p.lit = p.scanner.TokenText()
	default:
		p.lit = ""
	}
	if debug {
		fmt.Printf("%s: %q -> %q\n", scanner.TokenString(p.tok), p.scanner.TokenText(), p.lit)
	}
}

func declTypeName(pkg *types.Package, name string) *types.TypeName {
	scope := pkg.Scope()
	if obj := scope.Lookup(name); obj != nil {
		return obj.(*types.TypeName)
	}
	obj := types.NewTypeName(token.NoPos, pkg, name, nil)
	
	
	types.NewNamed(obj, nil, nil)
	scope.Insert(obj)
	return obj
}





type importError struct {
	pos scanner.Position
	err error
}

func (e importError) Error() string {
	return fmt.Sprintf("import error %s (byte offset = %d): %s", e.pos, e.pos.Offset, e.err)
}

func (p *parser) error(err interface{}) {
	if s, ok := err.(string); ok {
		err = errors.New(s)
	}
	
	panic(importError{p.scanner.Pos(), err.(error)})
}

func (p *parser) errorf(format string, args ...interface{}) {
	p.error(fmt.Sprintf(format, args...))
}

func (p *parser) expect(tok rune) string {
	lit := p.lit
	if p.tok != tok {
		p.errorf("expected %s, got %s (%s)", scanner.TokenString(tok), scanner.TokenString(p.tok), lit)
	}
	p.next()
	return lit
}

func (p *parser) expectSpecial(tok string) {
	sep := 'x' 
	i := 0
	for i < len(tok) && p.tok == rune(tok[i]) && sep > ' ' {
		sep = p.scanner.Peek() 
		p.next()
		i++
	}
	if i < len(tok) {
		p.errorf("expected %q, got %q", tok, tok[0:i])
	}
}

func (p *parser) expectKeyword(keyword string) {
	lit := p.expect(scanner.Ident)
	if lit != keyword {
		p.errorf("expected keyword %s, got %q", keyword, lit)
	}
}






func (p *parser) parsePackageId() string {
	id, err := strconv.Unquote(p.expect(scanner.String))
	if err != nil {
		p.error(err)
	}
	
	
	if id == "" {
		id = p.id
	}
	return id
}



func (p *parser) parsePackageName() string {
	return p.expect(scanner.Ident)
}


func (p *parser) parseDotIdent() string {
	ident := ""
	if p.tok != scanner.Int {
		sep := 'x' 
		for (p.tok == scanner.Ident || p.tok == scanner.Int || p.tok == '·') && sep > ' ' {
			ident += p.lit
			sep = p.scanner.Peek() 
			p.next()
		}
	}
	if ident == "" {
		p.expect(scanner.Ident) 
	}
	return ident
}



func (p *parser) parseQualifiedName() (id, name string) {
	p.expect('@')
	id = p.parsePackageId()
	p.expect('.')
	
	if p.tok == '?' {
		p.next()
	} else {
		name = p.parseDotIdent()
	}
	return
}











func (p *parser) getPkg(id, name string) *types.Package {
	
	if id == "unsafe" {
		return types.Unsafe
	}

	pkg := p.localPkgs[id]
	if pkg == nil {
		
		pkg = p.sharedPkgs[id]
		if pkg == nil {
			
			
			pkg = types.NewPackage(id, name)
			p.sharedPkgs[id] = pkg
		}
		
		if p.localPkgs == nil {
			p.localPkgs = make(map[string]*types.Package)
		}
		p.localPkgs[id] = pkg
	} else if name != "" {
		
		
		if pname := pkg.Name(); pname == "" {
			pkg.SetName(name)
		} else if pname != name {
			p.errorf("%s package name mismatch: %s (given) vs %s (expected)", id, pname, name)
		}
	}
	return pkg
}




func (p *parser) parseExportedName() (pkg *types.Package, name string) {
	id, name := p.parseQualifiedName()
	pkg = p.getPkg(id, "")
	return
}






func (p *parser) parseBasicType() types.Type {
	id := p.expect(scanner.Ident)
	obj := types.Universe.Lookup(id)
	if obj, ok := obj.(*types.TypeName); ok {
		return obj.Type()
	}
	p.errorf("not a basic type: %s", id)
	return nil
}



func (p *parser) parseArrayType(parent *types.Package) types.Type {
	
	lit := p.expect(scanner.Int)
	p.expect(']')
	elem := p.parseType(parent)
	n, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		p.error(err)
	}
	return types.NewArray(elem, n)
}



func (p *parser) parseMapType(parent *types.Package) types.Type {
	p.expectKeyword("map")
	p.expect('[')
	key := p.parseType(parent)
	p.expect(']')
	elem := p.parseType(parent)
	return types.NewMap(key, elem)
}














func (p *parser) parseName(parent *types.Package, materializePkg bool) (pkg *types.Package, name string) {
	pkg = parent
	if pkg == nil {
		pkg = p.sharedPkgs[p.id]
	}
	switch p.tok {
	case scanner.Ident:
		name = p.lit
		p.next()
	case '?':
		
		p.next()
	case '@':
		
		pkg = nil
		var id string
		id, name = p.parseQualifiedName()
		if materializePkg {
			pkg = p.getPkg(id, "")
		}
	default:
		p.error("name expected")
	}
	return
}

func deref(typ types.Type) types.Type {
	if p, _ := typ.(*types.Pointer); p != nil {
		return p.Elem()
	}
	return typ
}



func (p *parser) parseField(parent *types.Package) (*types.Var, string) {
	pkg, name := p.parseName(parent, true)

	if name == "_" {
		
		
		
		
		
		
		
		const blankpkg = "<_>"
		pkg = p.getPkg(blankpkg, blankpkg)
	}

	typ := p.parseType(parent)
	anonymous := false
	if name == "" {
		
		switch typ := deref(typ).(type) {
		case *types.Basic: 
			pkg = nil 
			name = typ.Name()
		case *types.Named:
			name = typ.Obj().Name()
		default:
			p.errorf("anonymous field expected")
		}
		anonymous = true
	}
	tag := ""
	if p.tok == scanner.String {
		s := p.expect(scanner.String)
		var err error
		tag, err = strconv.Unquote(s)
		if err != nil {
			p.errorf("invalid struct tag %s: %s", s, err)
		}
	}
	return types.NewField(token.NoPos, pkg, name, typ, anonymous), tag
}




func (p *parser) parseStructType(parent *types.Package) types.Type {
	var fields []*types.Var
	var tags []string

	p.expectKeyword("struct")
	p.expect('{')
	for i := 0; p.tok != '}' && p.tok != scanner.EOF; i++ {
		if i > 0 {
			p.expect(';')
		}
		fld, tag := p.parseField(parent)
		if tag != "" && tags == nil {
			tags = make([]string, i)
		}
		if tags != nil {
			tags = append(tags, tag)
		}
		fields = append(fields, fld)
	}
	p.expect('}')

	return types.NewStruct(fields, tags)
}



func (p *parser) parseParameter() (par *types.Var, isVariadic bool) {
	_, name := p.parseName(nil, false)
	
	if i := strings.Index(name, "·"); i >= 0 {
		name = name[:i]
	}
	if p.tok == '.' {
		p.expectSpecial("...")
		isVariadic = true
	}
	typ := p.parseType(nil)
	if isVariadic {
		typ = types.NewSlice(typ)
	}
	
	if p.tok == scanner.String {
		p.next()
	}
	
	par = types.NewVar(token.NoPos, nil, name, typ)
	return
}




func (p *parser) parseParameters() (list []*types.Var, isVariadic bool) {
	p.expect('(')
	for p.tok != ')' && p.tok != scanner.EOF {
		if len(list) > 0 {
			p.expect(',')
		}
		par, variadic := p.parseParameter()
		list = append(list, par)
		if variadic {
			if isVariadic {
				p.error("... not on final argument")
			}
			isVariadic = true
		}
	}
	p.expect(')')

	return
}




func (p *parser) parseSignature(recv *types.Var) *types.Signature {
	params, isVariadic := p.parseParameters()

	
	var results []*types.Var
	if p.tok == '(' {
		var variadic bool
		results, variadic = p.parseParameters()
		if variadic {
			p.error("... not permitted on result type")
		}
	}

	return types.NewSignature(recv, types.NewTuple(params...), types.NewTuple(results...), isVariadic)
}









func (p *parser) parseInterfaceType(parent *types.Package) types.Type {
	var methods []*types.Func

	p.expectKeyword("interface")
	p.expect('{')
	for i := 0; p.tok != '}' && p.tok != scanner.EOF; i++ {
		if i > 0 {
			p.expect(';')
		}
		pkg, name := p.parseName(parent, true)
		sig := p.parseSignature(nil)
		methods = append(methods, types.NewFunc(token.NoPos, pkg, name, sig))
	}
	p.expect('}')

	
	
	return types.NewInterface(methods, nil).Complete()
}



func (p *parser) parseChanType(parent *types.Package) types.Type {
	dir := types.SendRecv
	if p.tok == scanner.Ident {
		p.expectKeyword("chan")
		if p.tok == '<' {
			p.expectSpecial("<-")
			dir = types.SendOnly
		}
	} else {
		p.expectSpecial("<-")
		p.expectKeyword("chan")
		dir = types.RecvOnly
	}
	elem := p.parseType(parent)
	return types.NewChan(dir, elem)
}












func (p *parser) parseType(parent *types.Package) types.Type {
	switch p.tok {
	case scanner.Ident:
		switch p.lit {
		default:
			return p.parseBasicType()
		case "struct":
			return p.parseStructType(parent)
		case "func":
			
			p.next()
			return p.parseSignature(nil)
		case "interface":
			return p.parseInterfaceType(parent)
		case "map":
			return p.parseMapType(parent)
		case "chan":
			return p.parseChanType(parent)
		}
	case '@':
		
		pkg, name := p.parseExportedName()
		return declTypeName(pkg, name).Type()
	case '[':
		p.next() 
		if p.tok == ']' {
			
			p.next()
			return types.NewSlice(p.parseType(parent))
		}
		return p.parseArrayType(parent)
	case '*':
		
		p.next()
		return types.NewPointer(p.parseType(parent))
	case '<':
		return p.parseChanType(parent)
	case '(':
		
		p.next()
		typ := p.parseType(parent)
		p.expect(')')
		return typ
	}
	p.errorf("expected type, got %s (%q)", scanner.TokenString(p.tok), p.lit)
	return nil
}






func (p *parser) parseImportDecl() {
	p.expectKeyword("import")
	name := p.parsePackageName()
	p.getPkg(p.parsePackageId(), name)
}



func (p *parser) parseInt() string {
	s := ""
	switch p.tok {
	case '-':
		s = "-"
		p.next()
	case '+':
		p.next()
	}
	return s + p.expect(scanner.Int)
}



func (p *parser) parseNumber() (typ *types.Basic, val constant.Value) {
	
	mant := constant.MakeFromLiteral(p.parseInt(), token.INT, 0)
	if mant == nil {
		panic("invalid mantissa")
	}

	if p.lit == "p" {
		
		p.next()
		exp, err := strconv.ParseInt(p.parseInt(), 10, 0)
		if err != nil {
			p.error(err)
		}
		if exp < 0 {
			denom := constant.MakeInt64(1)
			denom = constant.Shift(denom, token.SHL, uint(-exp))
			typ = types.Typ[types.UntypedFloat]
			val = constant.BinaryOp(mant, token.QUO, denom)
			return
		}
		if exp > 0 {
			mant = constant.Shift(mant, token.SHL, uint(exp))
		}
		typ = types.Typ[types.UntypedFloat]
		val = mant
		return
	}

	typ = types.Typ[types.UntypedInt]
	val = mant
	return
}








func (p *parser) parseConstDecl() {
	p.expectKeyword("const")
	pkg, name := p.parseExportedName()

	var typ0 types.Type
	if p.tok != '=' {
		
		typ0 = p.parseType(nil)
	}

	p.expect('=')
	var typ types.Type
	var val constant.Value
	switch p.tok {
	case scanner.Ident:
		
		if p.lit != "true" && p.lit != "false" {
			p.error("expected true or false")
		}
		typ = types.Typ[types.UntypedBool]
		val = constant.MakeBool(p.lit == "true")
		p.next()

	case '-', scanner.Int:
		
		typ, val = p.parseNumber()

	case '(':
		
		p.next()
		if p.tok == scanner.Char {
			p.next()
			p.expect('+')
			typ = types.Typ[types.UntypedRune]
			_, val = p.parseNumber()
			p.expect(')')
			break
		}
		_, re := p.parseNumber()
		p.expect('+')
		_, im := p.parseNumber()
		p.expectKeyword("i")
		p.expect(')')
		typ = types.Typ[types.UntypedComplex]
		val = constant.BinaryOp(re, token.ADD, constant.MakeImag(im))

	case scanner.Char:
		
		typ = types.Typ[types.UntypedRune]
		val = constant.MakeFromLiteral(p.lit, token.CHAR, 0)
		p.next()

	case scanner.String:
		
		typ = types.Typ[types.UntypedString]
		val = constant.MakeFromLiteral(p.lit, token.STRING, 0)
		p.next()

	default:
		p.errorf("expected literal got %s", scanner.TokenString(p.tok))
	}

	if typ0 == nil {
		typ0 = typ
	}

	pkg.Scope().Insert(types.NewConst(token.NoPos, pkg, name, typ0, val))
}



func (p *parser) parseTypeDecl() {
	p.expectKeyword("type")
	pkg, name := p.parseExportedName()
	obj := declTypeName(pkg, name)

	
	
	
	
	
	typ := p.parseType(pkg)

	if name := obj.Type().(*types.Named); name.Underlying() == nil {
		name.SetUnderlying(typ)
	}
}



func (p *parser) parseVarDecl() {
	p.expectKeyword("var")
	pkg, name := p.parseExportedName()
	typ := p.parseType(pkg)
	pkg.Scope().Insert(types.NewVar(token.NoPos, pkg, name, typ))
}




func (p *parser) parseFunc(recv *types.Var) *types.Signature {
	sig := p.parseSignature(recv)
	if p.tok == '{' {
		p.next()
		for i := 1; i > 0; p.next() {
			switch p.tok {
			case '{':
				i++
			case '}':
				i--
			}
		}
	}
	return sig
}




func (p *parser) parseMethodDecl() {
	
	p.expect('(')
	recv, _ := p.parseParameter() 
	p.expect(')')

	
	base := deref(recv.Type()).(*types.Named)

	
	_, name := p.parseName(nil, false)
	sig := p.parseFunc(recv)

	
	pkg := base.Obj().Pkg()

	
	
	
	base.AddMethod(types.NewFunc(token.NoPos, pkg, name, sig))
}



func (p *parser) parseFuncDecl() {
	
	pkg, name := p.parseExportedName()
	typ := p.parseFunc(nil)
	pkg.Scope().Insert(types.NewFunc(token.NoPos, pkg, name, typ))
}



func (p *parser) parseDecl() {
	if p.tok == scanner.Ident {
		switch p.lit {
		case "import":
			p.parseImportDecl()
		case "const":
			p.parseConstDecl()
		case "type":
			p.parseTypeDecl()
		case "var":
			p.parseVarDecl()
		case "func":
			p.next() 
			if p.tok == '(' {
				p.parseMethodDecl()
			} else {
				p.parseFuncDecl()
			}
		}
	}
	p.expect('\n')
}







func (p *parser) parseExport() *types.Package {
	p.expectKeyword("package")
	name := p.parsePackageName()
	if p.tok == scanner.Ident && p.lit == "safe" {
		
		p.next()
	}
	p.expect('\n')

	pkg := p.getPkg(p.id, name)

	for p.tok != '$' && p.tok != scanner.EOF {
		p.parseDecl()
	}

	if ch := p.scanner.Peek(); p.tok != '$' || ch != '$' {
		
		
		p.errorf("expected '$$', got %s %c", scanner.TokenString(p.tok), ch)
	}

	if n := p.scanner.ErrorCount; n != 0 {
		p.errorf("expected no scanner errors, got %d", n)
	}

	
	var imports []*types.Package
	for id, pkg2 := range p.localPkgs {
		if pkg2.Name() == "" {
			p.errorf("%s package has no name", id)
		}
		if id == p.id {
			continue 
		}
		imports = append(imports, pkg2)
	}
	sort.Sort(byPath(imports))
	pkg.SetImports(imports)

	
	pkg.MarkComplete()

	return pkg
}

type byPath []*types.Package

func (a byPath) Len() int           { return len(a) }
func (a byPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPath) Less(i, j int) bool { return a[i].Path() < a[j].Path() }
