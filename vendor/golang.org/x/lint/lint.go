






package lint 

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/gcexportdata"
)

const styleGuideBase = "https://golang.org/wiki/CodeReviewComments"


type Linter struct {
}


type Problem struct {
	Position   token.Position 
	Text       string         
	Link       string         
	Confidence float64        
	LineText   string         
	Category   string         

	
	
	ReplacementLine string
}

func (p *Problem) String() string {
	if p.Link != "" {
		return p.Text + "\n\n" + p.Link
	}
	return p.Text
}

type byPosition []Problem

func (p byPosition) Len() int      { return len(p) }
func (p byPosition) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p byPosition) Less(i, j int) bool {
	pi, pj := p[i].Position, p[j].Position

	if pi.Filename != pj.Filename {
		return pi.Filename < pj.Filename
	}
	if pi.Line != pj.Line {
		return pi.Line < pj.Line
	}
	if pi.Column != pj.Column {
		return pi.Column < pj.Column
	}

	return p[i].Text < p[j].Text
}


func (l *Linter) Lint(filename string, src []byte) ([]Problem, error) {
	return l.LintFiles(map[string][]byte{filename: src})
}



func (l *Linter) LintFiles(files map[string][]byte) ([]Problem, error) {
	pkg := &pkg{
		fset:  token.NewFileSet(),
		files: make(map[string]*file),
	}
	var pkgName string
	for filename, src := range files {
		if isGenerated(src) {
			continue 
		}
		f, err := parser.ParseFile(pkg.fset, filename, src, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if pkgName == "" {
			pkgName = f.Name.Name
		} else if f.Name.Name != pkgName {
			return nil, fmt.Errorf("%s is in package %s, not %s", filename, f.Name.Name, pkgName)
		}
		pkg.files[filename] = &file{
			pkg:      pkg,
			f:        f,
			fset:     pkg.fset,
			src:      src,
			filename: filename,
		}
	}
	if len(pkg.files) == 0 {
		return nil, nil
	}
	return pkg.lint(), nil
}

var (
	genHdr = []byte("
	genFtr = []byte(" DO NOT EDIT.")
)



func isGenerated(src []byte) bool {
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		b := sc.Bytes()
		if bytes.HasPrefix(b, genHdr) && bytes.HasSuffix(b, genFtr) && len(b) >= len(genHdr)+len(genFtr) {
			return true
		}
	}
	return false
}


type pkg struct {
	fset  *token.FileSet
	files map[string]*file

	typesPkg  *types.Package
	typesInfo *types.Info

	
	sortable map[string]bool
	
	main bool

	problems []Problem
}

func (p *pkg) lint() []Problem {
	if err := p.typeCheck(); err != nil {
		
	}

	p.scanSortable()
	p.main = p.isMain()

	for _, f := range p.files {
		f.lint()
	}

	sort.Sort(byPosition(p.problems))

	return p.problems
}


type file struct {
	pkg      *pkg
	f        *ast.File
	fset     *token.FileSet
	src      []byte
	filename string
}

func (f *file) isTest() bool { return strings.HasSuffix(f.filename, "_test.go") }

func (f *file) lint() {
	f.lintPackageComment()
	f.lintImports()
	f.lintBlankImports()
	f.lintExported()
	f.lintNames()
	f.lintVarDecls()
	f.lintElses()
	f.lintRanges()
	f.lintErrorf()
	f.lintErrors()
	f.lintErrorStrings()
	f.lintReceiverNames()
	f.lintIncDec()
	f.lintErrorReturn()
	f.lintUnexportedReturn()
	f.lintTimeNames()
	f.lintContextKeyTypes()
	f.lintContextArgs()
}

type link string
type category string




func (f *file) errorf(n ast.Node, confidence float64, args ...interface{}) *Problem {
	pos := f.fset.Position(n.Pos())
	if pos.Filename == "" {
		pos.Filename = f.filename
	}
	return f.pkg.errorfAt(pos, confidence, args...)
}

func (p *pkg) errorfAt(pos token.Position, confidence float64, args ...interface{}) *Problem {
	problem := Problem{
		Position:   pos,
		Confidence: confidence,
	}
	if pos.Filename != "" {
		
		if f, ok := p.files[pos.Filename]; ok {
			problem.LineText = srcLine(f.src, pos)
		}
	}

argLoop:
	for len(args) > 1 { 
		switch v := args[0].(type) {
		case link:
			problem.Link = string(v)
		case category:
			problem.Category = string(v)
		default:
			break argLoop
		}
		args = args[1:]
	}

	problem.Text = fmt.Sprintf(args[0].(string), args[1:]...)

	p.problems = append(p.problems, problem)
	return &p.problems[len(p.problems)-1]
}

var newImporter = func(fset *token.FileSet) types.ImporterFrom {
	return gcexportdata.NewImporter(fset, make(map[string]*types.Package))
}

func (p *pkg) typeCheck() error {
	config := &types.Config{
		
		Error:    func(error) {},
		Importer: newImporter(p.fset),
	}
	info := &types.Info{
		Types:  make(map[ast.Expr]types.TypeAndValue),
		Defs:   make(map[*ast.Ident]types.Object),
		Uses:   make(map[*ast.Ident]types.Object),
		Scopes: make(map[ast.Node]*types.Scope),
	}
	var anyFile *file
	var astFiles []*ast.File
	for _, f := range p.files {
		anyFile = f
		astFiles = append(astFiles, f.f)
	}
	pkg, err := config.Check(anyFile.f.Name.Name, p.fset, astFiles, info)
	
	
	p.typesPkg = pkg
	p.typesInfo = info
	return err
}

func (p *pkg) typeOf(expr ast.Expr) types.Type {
	if p.typesInfo == nil {
		return nil
	}
	return p.typesInfo.TypeOf(expr)
}

func (p *pkg) isNamedType(typ types.Type, importPath, name string) bool {
	n, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	tn := n.Obj()
	return tn != nil && tn.Pkg() != nil && tn.Pkg().Path() == importPath && tn.Name() == name
}


func (p *pkg) scopeOf(id *ast.Ident) *types.Scope {
	var scope *types.Scope
	if obj := p.typesInfo.ObjectOf(id); obj != nil {
		scope = obj.Parent()
	}
	if scope == p.typesPkg.Scope() {
		
		
		pos := id.Pos()
		for _, f := range p.files {
			if f.f.Pos() <= pos && pos < f.f.End() {
				scope = p.typesInfo.Scopes[f.f]
				break
			}
		}
	}
	return scope
}

func (p *pkg) scanSortable() {
	p.sortable = make(map[string]bool)

	
	const (
		Len = 1 << iota
		Less
		Swap
	)
	nmap := map[string]int{"Len": Len, "Less": Less, "Swap": Swap}
	has := make(map[string]int)
	for _, f := range p.files {
		f.walk(func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
				return true
			}
			
			recv := receiverType(fn)
			if i, ok := nmap[fn.Name.Name]; ok {
				has[recv] |= i
			}
			return false
		})
	}
	for typ, ms := range has {
		if ms == Len|Less|Swap {
			p.sortable[typ] = true
		}
	}
}

func (p *pkg) isMain() bool {
	for _, f := range p.files {
		if f.isMain() {
			return true
		}
	}
	return false
}

func (f *file) isMain() bool {
	if f.f.Name.Name == "main" {
		return true
	}
	return false
}






func (f *file) lintPackageComment() {
	if f.isTest() {
		return
	}

	const ref = styleGuideBase + "#package-comments"
	prefix := "Package " + f.f.Name.Name + " "

	
	
	var lastCG *ast.CommentGroup
	for _, cg := range f.f.Comments {
		if cg.Pos() > f.f.Package {
			
			break
		}
		lastCG = cg
	}
	if lastCG != nil && strings.HasPrefix(lastCG.Text(), prefix) {
		endPos := f.fset.Position(lastCG.End())
		pkgPos := f.fset.Position(f.f.Package)
		if endPos.Line+1 < pkgPos.Line {
			
			
			
			pos := token.Position{
				Filename: endPos.Filename,
				
				Line:   endPos.Line + 1,
				Column: 1,
			}
			f.pkg.errorfAt(pos, 0.9, link(ref), category("comments"), "package comment is detached; there should be no blank lines between it and the package statement")
			return
		}
	}

	if f.f.Doc == nil {
		f.errorf(f.f, 0.2, link(ref), category("comments"), "should have a package comment, unless it's in another file for this package")
		return
	}
	s := f.f.Doc.Text()
	if ts := strings.TrimLeft(s, " \t"); ts != s {
		f.errorf(f.f.Doc, 1, link(ref), category("comments"), "package comment should not have leading space")
		s = ts
	}
	
	if !f.pkg.main && !strings.HasPrefix(s, prefix) {
		f.errorf(f.f.Doc, 1, link(ref), category("comments"), `package comment should be of the form "%s..."`, prefix)
	}
}



func (f *file) lintBlankImports() {
	
	if f.pkg.main || f.isTest() {
		return
	}

	
	
	for i, imp := range f.f.Imports {
		pos := f.fset.Position(imp.Pos())

		if !isBlank(imp.Name) {
			continue 
		}
		if i > 0 {
			prev := f.f.Imports[i-1]
			prevPos := f.fset.Position(prev.Pos())
			if isBlank(prev.Name) && prevPos.Line+1 == pos.Line {
				continue 
			}
		}

		
		if imp.Doc == nil && imp.Comment == nil {
			ref := ""
			f.errorf(imp, 1, link(ref), category("imports"), "a blank import should be only in a main or test package, or have a comment justifying it")
		}
	}
}


func (f *file) lintImports() {
	for i, is := range f.f.Imports {
		_ = i
		if is.Name != nil && is.Name.Name == "." && !f.isTest() {
			f.errorf(is, 1, link(styleGuideBase+"#import-dot"), category("imports"), "should not use dot imports")
		}

	}
}

const docCommentsLink = styleGuideBase + "#doc-comments"









func (f *file) lintExported() {
	if f.isTest() {
		return
	}

	var lastGen *ast.GenDecl 

	
	genDeclMissingComments := make(map[*ast.GenDecl]bool)

	f.walk(func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.GenDecl:
			if v.Tok == token.IMPORT {
				return false
			}
			
			lastGen = v
			return true
		case *ast.FuncDecl:
			f.lintFuncDoc(v)
			if v.Recv == nil {
				
				
				f.checkStutter(v.Name, "func")
			}
			
			return false
		case *ast.TypeSpec:
			
			doc := v.Doc
			if doc == nil {
				doc = lastGen.Doc
			}
			f.lintTypeDoc(v, doc)
			f.checkStutter(v.Name, "type")
			
			return false
		case *ast.ValueSpec:
			f.lintValueSpecDoc(v, lastGen, genDeclMissingComments)
			return false
		}
		return true
	})
}

var (
	allCapsRE = regexp.MustCompile(`^[A-Z0-9_]+$`)
	anyCapsRE = regexp.MustCompile(`[A-Z]`)
)




var knownNameExceptions = map[string]bool{
	"LastInsertId": true, 
	"kWh":          true,
}



func (f *file) lintNames() {
	
	if strings.Contains(f.f.Name.Name, "_") && !strings.HasSuffix(f.f.Name.Name, "_test") {
		f.errorf(f.f, 1, link("http://golang.org/doc/effective_go.html#package-names"), category("naming"), "don't use an underscore in package name")
	}
	if anyCapsRE.MatchString(f.f.Name.Name) {
		f.errorf(f.f, 1, link("http://golang.org/doc/effective_go.html#package-names"), category("mixed-caps"), "don't use MixedCaps in package name; %s should be %s", f.f.Name.Name, strings.ToLower(f.f.Name.Name))
	}

	check := func(id *ast.Ident, thing string) {
		if id.Name == "_" {
			return
		}
		if knownNameExceptions[id.Name] {
			return
		}

		
		if len(id.Name) >= 5 && allCapsRE.MatchString(id.Name) && strings.Contains(id.Name, "_") {
			f.errorf(id, 0.8, link(styleGuideBase+"#mixed-caps"), category("naming"), "don't use ALL_CAPS in Go names; use CamelCase")
			return
		}
		if len(id.Name) > 2 && id.Name[0] == 'k' && id.Name[1] >= 'A' && id.Name[1] <= 'Z' {
			should := string(id.Name[1]+'a'-'A') + id.Name[2:]
			f.errorf(id, 0.8, link(styleGuideBase+"#mixed-caps"), category("naming"), "don't use leading k in Go names; %s %s should be %s", thing, id.Name, should)
		}

		should := lintName(id.Name)
		if id.Name == should {
			return
		}

		if len(id.Name) > 2 && strings.Contains(id.Name[1:], "_") {
			f.errorf(id, 0.9, link("http://golang.org/doc/effective_go.html#mixed-caps"), category("naming"), "don't use underscores in Go names; %s %s should be %s", thing, id.Name, should)
			return
		}
		f.errorf(id, 0.8, link(styleGuideBase+"#initialisms"), category("naming"), "%s %s should be %s", thing, id.Name, should)
	}
	checkList := func(fl *ast.FieldList, thing string) {
		if fl == nil {
			return
		}
		for _, f := range fl.List {
			for _, id := range f.Names {
				check(id, thing)
			}
		}
	}
	f.walk(func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.AssignStmt:
			if v.Tok == token.ASSIGN {
				return true
			}
			for _, exp := range v.Lhs {
				if id, ok := exp.(*ast.Ident); ok {
					check(id, "var")
				}
			}
		case *ast.FuncDecl:
			if f.isTest() && (strings.HasPrefix(v.Name.Name, "Example") || strings.HasPrefix(v.Name.Name, "Test") || strings.HasPrefix(v.Name.Name, "Benchmark")) {
				return true
			}

			thing := "func"
			if v.Recv != nil {
				thing = "method"
			}

			
			
			
			if ast.IsExported(v.Name.Name) || !isCgoExported(v) {
				check(v.Name, thing)
			}

			checkList(v.Type.Params, thing+" parameter")
			checkList(v.Type.Results, thing+" result")
		case *ast.GenDecl:
			if v.Tok == token.IMPORT {
				return true
			}
			var thing string
			switch v.Tok {
			case token.CONST:
				thing = "const"
			case token.TYPE:
				thing = "type"
			case token.VAR:
				thing = "var"
			}
			for _, spec := range v.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					check(s.Name, thing)
				case *ast.ValueSpec:
					for _, id := range s.Names {
						check(id, thing)
					}
				}
			}
		case *ast.InterfaceType:
			
			
			for _, x := range v.Methods.List {
				ft, ok := x.Type.(*ast.FuncType)
				if !ok { 
					continue
				}
				checkList(ft.Params, "interface method parameter")
				checkList(ft.Results, "interface method result")
			}
		case *ast.RangeStmt:
			if v.Tok == token.ASSIGN {
				return true
			}
			if id, ok := v.Key.(*ast.Ident); ok {
				check(id, "range var")
			}
			if id, ok := v.Value.(*ast.Ident); ok {
				check(id, "range var")
			}
		case *ast.StructType:
			for _, f := range v.Fields.List {
				for _, id := range f.Names {
					check(id, "struct field")
				}
			}
		}
		return true
	})
}


func lintName(name string) (should string) {
	
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		return name
	}

	
	
	runes := []rune(name)
	w, i := 0, 0 
	for i+1 <= len(runes) {
		eow := false 
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			
			eow = true
		}
		i++
		if !eow {
			continue
		}

		
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			
			if w == 0 && unicode.IsLower(runes[w]) {
				u = strings.ToLower(u)
			}
			
			
			copy(runes[w:], []rune(u))
		} else if w > 0 && strings.ToLower(word) == word {
			
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}




var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}




func (f *file) lintTypeDoc(t *ast.TypeSpec, doc *ast.CommentGroup) {
	if !ast.IsExported(t.Name.Name) {
		return
	}
	if doc == nil {
		f.errorf(t, 1, link(docCommentsLink), category("comments"), "exported type %v should have comment or be unexported", t.Name)
		return
	}

	s := doc.Text()
	articles := [...]string{"A", "An", "The"}
	for _, a := range articles {
		if strings.HasPrefix(s, a+" ") {
			s = s[len(a)+1:]
			break
		}
	}
	if !strings.HasPrefix(s, t.Name.Name+" ") {
		f.errorf(doc, 1, link(docCommentsLink), category("comments"), `comment on exported type %v should be of the form "%v ..." (with optional leading article)`, t.Name, t.Name)
	}
}

var commonMethods = map[string]bool{
	"Error":     true,
	"Read":      true,
	"ServeHTTP": true,
	"String":    true,
	"Write":     true,
}




func (f *file) lintFuncDoc(fn *ast.FuncDecl) {
	if !ast.IsExported(fn.Name.Name) {
		
		return
	}
	kind := "function"
	name := fn.Name.Name
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		
		kind = "method"
		recv := receiverType(fn)
		if !ast.IsExported(recv) {
			
			return
		}
		if commonMethods[name] {
			return
		}
		switch name {
		case "Len", "Less", "Swap":
			if f.pkg.sortable[recv] {
				return
			}
		}
		name = recv + "." + name
	}
	if fn.Doc == nil {
		f.errorf(fn, 1, link(docCommentsLink), category("comments"), "exported %s %s should have comment or be unexported", kind, name)
		return
	}
	s := fn.Doc.Text()
	prefix := fn.Name.Name + " "
	if !strings.HasPrefix(s, prefix) {
		f.errorf(fn.Doc, 1, link(docCommentsLink), category("comments"), `comment on exported %s %s should be of the form "%s..."`, kind, name, prefix)
	}
}




func (f *file) lintValueSpecDoc(vs *ast.ValueSpec, gd *ast.GenDecl, genDeclMissingComments map[*ast.GenDecl]bool) {
	kind := "var"
	if gd.Tok == token.CONST {
		kind = "const"
	}

	if len(vs.Names) > 1 {
		
		for _, n := range vs.Names[1:] {
			if ast.IsExported(n.Name) {
				f.errorf(vs, 1, category("comments"), "exported %s %s should have its own declaration", kind, n.Name)
				return
			}
		}
	}

	
	name := vs.Names[0].Name
	if !ast.IsExported(name) {
		return
	}

	if vs.Doc == nil && gd.Doc == nil {
		if genDeclMissingComments[gd] {
			return
		}
		block := ""
		if kind == "const" && gd.Lparen.IsValid() {
			block = " (or a comment on this block)"
		}
		f.errorf(vs, 1, link(docCommentsLink), category("comments"), "exported %s %s should have comment%s or be unexported", kind, name, block)
		genDeclMissingComments[gd] = true
		return
	}
	
	if gd.Lparen.IsValid() && gd.Doc != nil {
		return
	}
	
	
	doc := vs.Doc
	if doc == nil {
		doc = gd.Doc
	}
	prefix := name + " "
	if !strings.HasPrefix(doc.Text(), prefix) {
		f.errorf(doc, 1, link(docCommentsLink), category("comments"), `comment on exported %s %s should be of the form "%s..."`, kind, name, prefix)
	}
}

func (f *file) checkStutter(id *ast.Ident, thing string) {
	pkg, name := f.f.Name.Name, id.Name
	if !ast.IsExported(name) {
		
		return
	}
	
	
	if len(name) <= len(pkg) {
		
		
		return
	}
	if !strings.EqualFold(pkg, name[:len(pkg)]) {
		return
	}
	
	
	
	rem := name[len(pkg):]
	if next, _ := utf8.DecodeRuneInString(rem); next == '_' || unicode.IsUpper(next) {
		f.errorf(id, 0.8, link(styleGuideBase+"#package-names"), category("naming"), "%s name will be used as %s.%s by other packages, and that stutters; consider calling this %s", thing, pkg, name, rem)
	}
}



var zeroLiteral = map[string]bool{
	"false": true, 
	
	`'\x00'`: true,
	`'\000'`: true,
	
	`""`: true,
	"``": true,
	
	"0":   true,
	"0.":  true,
	"0.0": true,
	"0i":  true,
}



func (f *file) lintVarDecls() {
	var lastGen *ast.GenDecl 

	f.walk(func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.GenDecl:
			if v.Tok != token.CONST && v.Tok != token.VAR {
				return false
			}
			lastGen = v
			return true
		case *ast.ValueSpec:
			if lastGen.Tok == token.CONST {
				return false
			}
			if len(v.Names) > 1 || v.Type == nil || len(v.Values) == 0 {
				return false
			}
			rhs := v.Values[0]
			
			
			if isIdent(v.Names[0], "_") {
				return false
			}
			
			zero := false
			if lit, ok := rhs.(*ast.BasicLit); ok {
				zero = zeroLiteral[lit.Value]
			} else if isIdent(rhs, "nil") {
				zero = true
			}
			if zero {
				f.errorf(rhs, 0.9, category("zero-value"), "should drop = %s from declaration of var %s; it is the zero value", f.render(rhs), v.Names[0])
				return false
			}
			lhsTyp := f.pkg.typeOf(v.Type)
			rhsTyp := f.pkg.typeOf(rhs)

			if !validType(lhsTyp) || !validType(rhsTyp) {
				
				return false
			}

			if !types.Identical(lhsTyp, rhsTyp) {
				
				return false
			}

			
			

			
			
			
			
			if _, ok := v.Type.(*ast.InterfaceType); ok {
				return false
			}
			
			if defType, ok := f.isUntypedConst(rhs); ok && !isIdent(v.Type, defType) {
				return false
			}

			f.errorf(v.Type, 0.8, category("type-inference"), "should omit type %s from declaration of var %s; it will be inferred from the right-hand side", f.render(v.Type), v.Names[0])
			return false
		}
		return true
	})
}

func validType(T types.Type) bool {
	return T != nil &&
		T != types.Typ[types.Invalid] &&
		!strings.Contains(T.String(), "invalid type") 
}


func (f *file) lintElses() {
	
	
	
	ignore := make(map[*ast.IfStmt]bool)

	f.walk(func(node ast.Node) bool {
		ifStmt, ok := node.(*ast.IfStmt)
		if !ok || ifStmt.Else == nil {
			return true
		}
		if elseif, ok := ifStmt.Else.(*ast.IfStmt); ok {
			ignore[elseif] = true
			return true
		}
		if ignore[ifStmt] {
			return true
		}
		if _, ok := ifStmt.Else.(*ast.BlockStmt); !ok {
			
			return true
		}
		if len(ifStmt.Body.List) == 0 {
			return true
		}
		shortDecl := false 
		if ifStmt.Init != nil {
			if as, ok := ifStmt.Init.(*ast.AssignStmt); ok && as.Tok == token.DEFINE {
				shortDecl = true
			}
		}
		lastStmt := ifStmt.Body.List[len(ifStmt.Body.List)-1]
		if _, ok := lastStmt.(*ast.ReturnStmt); ok {
			extra := ""
			if shortDecl {
				extra = " (move short variable declaration to its own line if necessary)"
			}
			f.errorf(ifStmt.Else, 1, link(styleGuideBase+"#indent-error-flow"), category("indent"), "if block ends with a return statement, so drop this else and outdent its block"+extra)
		}
		return true
	})
}


func (f *file) lintRanges() {
	f.walk(func(node ast.Node) bool {
		rs, ok := node.(*ast.RangeStmt)
		if !ok {
			return true
		}

		if isIdent(rs.Key, "_") && (rs.Value == nil || isIdent(rs.Value, "_")) {
			p := f.errorf(rs.Key, 1, category("range-loop"), "should omit values from range; this loop is equivalent to `for range ...`")

			newRS := *rs 
			newRS.Value = nil
			newRS.Key = nil
			p.ReplacementLine = f.firstLineOf(&newRS, rs)

			return true
		}

		if isIdent(rs.Value, "_") {
			p := f.errorf(rs.Value, 1, category("range-loop"), "should omit 2nd value from range; this loop is equivalent to `for %s %s range ...`", f.render(rs.Key), rs.Tok)

			newRS := *rs 
			newRS.Value = nil
			p.ReplacementLine = f.firstLineOf(&newRS, rs)
		}

		return true
	})
}


func (f *file) lintErrorf() {
	f.walk(func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok || len(ce.Args) != 1 {
			return true
		}
		isErrorsNew := isPkgDot(ce.Fun, "errors", "New")
		var isTestingError bool
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if ok && se.Sel.Name == "Error" {
			if typ := f.pkg.typeOf(se.X); typ != nil {
				isTestingError = typ.String() == "*testing.T"
			}
		}
		if !isErrorsNew && !isTestingError {
			return true
		}
		if !f.imports("errors") {
			return true
		}
		arg := ce.Args[0]
		ce, ok = arg.(*ast.CallExpr)
		if !ok || !isPkgDot(ce.Fun, "fmt", "Sprintf") {
			return true
		}
		errorfPrefix := "fmt"
		if isTestingError {
			errorfPrefix = f.render(se.X)
		}
		p := f.errorf(node, 1, category("errors"), "should replace %s(fmt.Sprintf(...)) with %s.Errorf(...)", f.render(se), errorfPrefix)

		m := f.srcLineWithMatch(ce, `^(.*)`+f.render(se)+`\(fmt\.Sprintf\((.*)\)\)(.*)$`)
		if m != nil {
			p.ReplacementLine = m[1] + errorfPrefix + ".Errorf(" + m[2] + ")" + m[3]
		}

		return true
	})
}


func (f *file) lintErrors() {
	for _, decl := range f.f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			spec := spec.(*ast.ValueSpec)
			if len(spec.Names) != 1 || len(spec.Values) != 1 {
				continue
			}
			ce, ok := spec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			if !isPkgDot(ce.Fun, "errors", "New") && !isPkgDot(ce.Fun, "fmt", "Errorf") {
				continue
			}

			id := spec.Names[0]
			prefix := "err"
			if id.IsExported() {
				prefix = "Err"
			}
			if !strings.HasPrefix(id.Name, prefix) {
				f.errorf(id, 0.9, category("naming"), "error var %s should have name of the form %sFoo", id.Name, prefix)
			}
		}
	}
}

func lintErrorString(s string) (isClean bool, conf float64) {
	const basicConfidence = 0.8
	const capConfidence = basicConfidence - 0.2
	first, firstN := utf8.DecodeRuneInString(s)
	last, _ := utf8.DecodeLastRuneInString(s)
	if last == '.' || last == ':' || last == '!' || last == '\n' {
		return false, basicConfidence
	}
	if unicode.IsUpper(first) {
		
		
		if len(s) <= firstN {
			return false, capConfidence
		}
		
		if second, _ := utf8.DecodeRuneInString(s[firstN:]); !unicode.IsUpper(second) {
			return false, capConfidence
		}
	}
	return true, 0
}



func (f *file) lintErrorStrings() {
	f.walk(func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if !isPkgDot(ce.Fun, "errors", "New") && !isPkgDot(ce.Fun, "fmt", "Errorf") {
			return true
		}
		if len(ce.Args) < 1 {
			return true
		}
		str, ok := ce.Args[0].(*ast.BasicLit)
		if !ok || str.Kind != token.STRING {
			return true
		}
		s, _ := strconv.Unquote(str.Value) 
		if s == "" {
			return true
		}
		clean, conf := lintErrorString(s)
		if clean {
			return true
		}

		f.errorf(str, conf, link(styleGuideBase+"#error-strings"), category("errors"),
			"error strings should not be capitalized or end with punctuation or a newline")
		return true
	})
}



func (f *file) lintReceiverNames() {
	typeReceiver := map[string]string{}
	f.walk(func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			return true
		}
		names := fn.Recv.List[0].Names
		if len(names) < 1 {
			return true
		}
		name := names[0].Name
		const ref = styleGuideBase + "#receiver-names"
		if name == "_" {
			f.errorf(n, 1, link(ref), category("naming"), `receiver name should not be an underscore, omit the name if it is unused`)
			return true
		}
		if name == "this" || name == "self" {
			f.errorf(n, 1, link(ref), category("naming"), `receiver name should be a reflection of its identity; don't use generic names such as "this" or "self"`)
			return true
		}
		recv := receiverType(fn)
		if prev, ok := typeReceiver[recv]; ok && prev != name {
			f.errorf(n, 1, link(ref), category("naming"), "receiver name %s should be consistent with previous receiver name %s for %s", name, prev, recv)
			return true
		}
		typeReceiver[recv] = name
		return true
	})
}



func (f *file) lintIncDec() {
	f.walk(func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		if len(as.Lhs) != 1 {
			return true
		}
		if !isOne(as.Rhs[0]) {
			return true
		}
		var suffix string
		switch as.Tok {
		case token.ADD_ASSIGN:
			suffix = "++"
		case token.SUB_ASSIGN:
			suffix = "--"
		default:
			return true
		}
		f.errorf(as, 0.8, category("unary-op"), "should replace %s with %s%s", f.render(as), f.render(as.Lhs[0]), suffix)
		return true
	})
}



func (f *file) lintErrorReturn() {
	f.walk(func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Type.Results == nil {
			return true
		}
		ret := fn.Type.Results.List
		if len(ret) <= 1 {
			return true
		}
		if isIdent(ret[len(ret)-1].Type, "error") {
			return true
		}
		
		
		for _, r := range ret[:len(ret)-1] {
			if isIdent(r.Type, "error") {
				f.errorf(fn, 0.9, category("arg-order"), "error should be the last type when returning multiple items")
				break 
			}
		}
		return true
	})
}



func (f *file) lintUnexportedReturn() {
	f.walk(func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if fn.Type.Results == nil {
			return false
		}
		if !fn.Name.IsExported() {
			return false
		}
		thing := "func"
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			thing = "method"
			if !ast.IsExported(receiverType(fn)) {
				
				
				return false
			}
		}
		for _, ret := range fn.Type.Results.List {
			typ := f.pkg.typeOf(ret.Type)
			if exportedType(typ) {
				continue
			}
			f.errorf(ret.Type, 0.8, category("unexported-type-in-api"),
				"exported %s %s returns unexported type %s, which can be annoying to use",
				thing, fn.Name.Name, typ)
			break 
		}
		return false
	})
}




func exportedType(typ types.Type) bool {
	switch T := typ.(type) {
	case *types.Named:
		
		return T.Obj().Pkg() == nil || T.Obj().Exported()
	case *types.Map:
		return exportedType(T.Key()) && exportedType(T.Elem())
	case interface {
		Elem() types.Type
	}: 
		return exportedType(T.Elem())
	}
	
	return true
}



var timeSuffixes = []string{
	"Sec", "Secs", "Seconds",
	"Msec", "Msecs",
	"Milli", "Millis", "Milliseconds",
	"Usec", "Usecs", "Microseconds",
	"MS", "Ms",
}

func (f *file) lintTimeNames() {
	f.walk(func(node ast.Node) bool {
		v, ok := node.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for _, name := range v.Names {
			origTyp := f.pkg.typeOf(name)
			
			
			typ := origTyp
			if pt, ok := typ.(*types.Pointer); ok {
				typ = pt.Elem()
			}
			if !f.pkg.isNamedType(typ, "time", "Duration") {
				continue
			}
			suffix := ""
			for _, suf := range timeSuffixes {
				if strings.HasSuffix(name.Name, suf) {
					suffix = suf
					break
				}
			}
			if suffix == "" {
				continue
			}
			f.errorf(v, 0.9, category("time"), "var %s is of type %v; don't use unit-specific suffix %q", name.Name, origTyp, suffix)
		}
		return true
	})
}




func (f *file) lintContextKeyTypes() {
	f.walk(func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.CallExpr:
			f.checkContextKeyType(node)
		}

		return true
	})
}



func (f *file) checkContextKeyType(x *ast.CallExpr) {
	sel, ok := x.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "context" {
		return
	}
	if sel.Sel.Name != "WithValue" {
		return
	}

	
	if len(x.Args) != 3 {
		return
	}
	key := f.pkg.typesInfo.Types[x.Args[1]]

	if ktyp, ok := key.Type.(*types.Basic); ok && ktyp.Kind() != types.Invalid {
		f.errorf(x, 1.0, category("context"), fmt.Sprintf("should not use basic type %s as key in context.WithValue", key.Type))
	}
}




func (f *file) lintContextArgs() {
	f.walk(func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || len(fn.Type.Params.List) <= 1 {
			return true
		}
		
		
		for _, arg := range fn.Type.Params.List[1:] {
			if isPkgDot(arg.Type, "context", "Context") {
				f.errorf(fn, 0.9, link("https://golang.org/pkg/context/"), category("arg-order"), "context.Context should be the first parameter of a function")
				break 
			}
		}
		return true
	})
}



func (f *file) containsComments(start, end token.Pos) bool {
	for _, cgroup := range f.f.Comments {
		comments := cgroup.List
		if comments[0].Slash >= end {
			
			return false
		}
		if comments[len(comments)-1].Slash < start {
			
			continue
		}
		for _, c := range comments {
			if start <= c.Slash && c.Slash < end && !strings.HasPrefix(c.Text, "
				return true
			}
		}
	}
	return false
}



func receiverType(fn *ast.FuncDecl) string {
	switch e := fn.Recv.List[0].Type.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		if id, ok := e.X.(*ast.Ident); ok {
			return id.Name
		}
	}
	
	return "invalid-type"
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.f)
}

func (f *file) render(x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, f.fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

func (f *file) debugRender(x interface{}) string {
	var buf bytes.Buffer
	if err := ast.Fprint(&buf, f.fset, x, nil); err != nil {
		panic(err)
	}
	return buf.String()
}



type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}



func isBlank(id *ast.Ident) bool { return id != nil && id.Name == "_" }

func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func isOne(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.INT && lit.Value == "1"
}

func isCgoExported(f *ast.FuncDecl) bool {
	if f.Recv != nil || f.Doc == nil {
		return false
	}

	cgoExport := regexp.MustCompile(fmt.Sprintf("(?m)^
	for _, c := range f.Doc.List {
		if cgoExport.MatchString(c.Text) {
			return true
		}
	}
	return false
}

var basicTypeKinds = map[types.BasicKind]string{
	types.UntypedBool:    "bool",
	types.UntypedInt:     "int",
	types.UntypedRune:    "rune",
	types.UntypedFloat:   "float64",
	types.UntypedComplex: "complex128",
	types.UntypedString:  "string",
}




func (f *file) isUntypedConst(expr ast.Expr) (defType string, ok bool) {
	
	
	exprStr := f.render(expr)
	tv, err := types.Eval(f.fset, f.pkg.typesPkg, expr.Pos(), exprStr)
	if err != nil {
		return "", false
	}
	if b, ok := tv.Type.(*types.Basic); ok {
		if dt, ok := basicTypeKinds[b.Kind()]; ok {
			return dt, true
		}
	}

	return "", false
}



func (f *file) firstLineOf(node, match ast.Node) string {
	line := f.render(node)
	if i := strings.Index(line, "\n"); i >= 0 {
		line = line[:i]
	}
	return f.indentOf(match) + line
}

func (f *file) indentOf(node ast.Node) string {
	line := srcLine(f.src, f.fset.Position(node.Pos()))
	for i, r := range line {
		switch r {
		case ' ', '\t':
		default:
			return line[:i]
		}
	}
	return line 
}

func (f *file) srcLineWithMatch(node ast.Node, pattern string) (m []string) {
	line := srcLine(f.src, f.fset.Position(node.Pos()))
	line = strings.TrimSuffix(line, "\n")
	rx := regexp.MustCompile(pattern)
	return rx.FindStringSubmatch(line)
}


func (f *file) imports(importPath string) bool {
	all := astutil.Imports(f.fset, f.f)
	for _, p := range all {
		for _, i := range p {
			uq, err := strconv.Unquote(i.Path.Value)
			if err == nil && importPath == uq {
				return true
			}
		}
	}
	return false
}


func srcLine(src []byte, p token.Position) string {
	
	lo, hi := p.Offset, p.Offset+1
	for lo > 0 && src[lo-1] != '\n' {
		lo--
	}
	for hi < len(src) && src[hi-1] != '\n' {
		hi++
	}
	return string(src[lo:hi])
}
