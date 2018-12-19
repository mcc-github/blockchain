



package packages



import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/tools/go/gcexportdata"
)




type LoadMode int

const (
	
	
	LoadFiles LoadMode = iota

	
	
	
	LoadImports

	
	
	
	
	
	LoadTypes

	
	
	LoadSyntax

	
	
	
	
	LoadAllSyntax
)




type Config struct {
	
	Mode LoadMode

	
	
	
	
	Context context.Context

	
	
	
	Dir string

	
	
	
	
	
	
	
	
	Env []string

	
	
	BuildFlags []string

	
	
	Fset *token.FileSet

	
	
	
	
	
	
	
	
	
	
	
	
	ParseFile func(fset *token.FileSet, filename string, src []byte) (*ast.File, error)

	
	
	
	
	
	
	
	
	
	
	
	
	Tests bool

	
	
	
	
	
	
	
	Overlay map[string][]byte
}



type driver func(cfg *Config, patterns ...string) (*driverResponse, error)


type driverResponse struct {
	
	
	
	
	Roots []string `json:",omitempty"`

	
	
	
	
	
	Packages []*Package
}















func Load(cfg *Config, patterns ...string) ([]*Package, error) {
	l := newLoader(cfg)
	response, err := defaultDriver(&l.Config, patterns...)
	if err != nil {
		return nil, err
	}
	return l.refine(response.Roots, response.Packages...)
}



func defaultDriver(cfg *Config, patterns ...string) (*driverResponse, error) {
	driver := findExternalDriver(cfg)
	if driver == nil {
		driver = goListDriver
	}
	return driver(cfg, patterns...)
}


type Package struct {
	
	
	
	
	
	
	ID string

	
	Name string

	
	PkgPath string

	
	
	Errors []Error

	
	GoFiles []string

	
	
	
	CompiledGoFiles []string

	
	
	OtherFiles []string

	
	
	ExportFile string

	
	
	Imports map[string]*Package

	
	
	
	
	Types *types.Package

	
	
	Fset *token.FileSet

	
	
	IllTyped bool

	
	
	
	
	Syntax []*ast.File

	
	
	TypesInfo *types.Info
}


type Error struct {
	Pos  string 
	Msg  string
	Kind ErrorKind
}




type ErrorKind int

const (
	UnknownError ErrorKind = iota
	ListError
	ParseError
	TypeError
)

func (err Error) Error() string {
	pos := err.Pos
	if pos == "" {
		pos = "-" 
	}
	return pos + ": " + err.Msg
}






type flatPackage struct {
	ID              string
	Name            string            `json:",omitempty"`
	PkgPath         string            `json:",omitempty"`
	Errors          []Error           `json:",omitempty"`
	GoFiles         []string          `json:",omitempty"`
	CompiledGoFiles []string          `json:",omitempty"`
	OtherFiles      []string          `json:",omitempty"`
	ExportFile      string            `json:",omitempty"`
	Imports         map[string]string `json:",omitempty"`
}










func (p *Package) MarshalJSON() ([]byte, error) {
	flat := &flatPackage{
		ID:              p.ID,
		Name:            p.Name,
		PkgPath:         p.PkgPath,
		Errors:          p.Errors,
		GoFiles:         p.GoFiles,
		CompiledGoFiles: p.CompiledGoFiles,
		OtherFiles:      p.OtherFiles,
		ExportFile:      p.ExportFile,
	}
	if len(p.Imports) > 0 {
		flat.Imports = make(map[string]string, len(p.Imports))
		for path, ipkg := range p.Imports {
			flat.Imports[path] = ipkg.ID
		}
	}
	return json.Marshal(flat)
}



func (p *Package) UnmarshalJSON(b []byte) error {
	flat := &flatPackage{}
	if err := json.Unmarshal(b, &flat); err != nil {
		return err
	}
	*p = Package{
		ID:              flat.ID,
		Name:            flat.Name,
		PkgPath:         flat.PkgPath,
		Errors:          flat.Errors,
		GoFiles:         flat.GoFiles,
		CompiledGoFiles: flat.CompiledGoFiles,
		OtherFiles:      flat.OtherFiles,
		ExportFile:      flat.ExportFile,
	}
	if len(flat.Imports) > 0 {
		p.Imports = make(map[string]*Package, len(flat.Imports))
		for path, id := range flat.Imports {
			p.Imports[path] = &Package{ID: id}
		}
	}
	return nil
}

func (p *Package) String() string { return p.ID }


type loaderPackage struct {
	*Package
	importErrors map[string]error 
	loadOnce     sync.Once
	color        uint8 
	needsrc      bool  
	needtypes    bool  
	initial      bool  
}


type loader struct {
	pkgs map[string]*loaderPackage
	Config
	exportMu sync.Mutex 
}

func newLoader(cfg *Config) *loader {
	ld := &loader{}
	if cfg != nil {
		ld.Config = *cfg
	}
	if ld.Config.Env == nil {
		ld.Config.Env = os.Environ()
	}
	if ld.Context == nil {
		ld.Context = context.Background()
	}
	if ld.Dir == "" {
		if dir, err := os.Getwd(); err == nil {
			ld.Dir = dir
		}
	}

	if ld.Mode >= LoadTypes {
		if ld.Fset == nil {
			ld.Fset = token.NewFileSet()
		}

		
		
		if ld.ParseFile == nil {
			ld.ParseFile = func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
				var isrc interface{}
				if src != nil {
					isrc = src
				}
				const mode = parser.AllErrors | parser.ParseComments
				return parser.ParseFile(fset, filename, isrc, mode)
			}
		}
	}
	return ld
}



func (ld *loader) refine(roots []string, list ...*Package) ([]*Package, error) {
	isRoot := make(map[string]bool, len(roots))
	for _, root := range roots {
		isRoot[root] = true
	}
	ld.pkgs = make(map[string]*loaderPackage)
	
	var initial []*loaderPackage
	for _, pkg := range list {
		lpkg := &loaderPackage{
			Package: pkg,
			needtypes: ld.Mode >= LoadAllSyntax ||
				ld.Mode >= LoadTypes && isRoot[pkg.ID],
			needsrc: ld.Mode >= LoadAllSyntax ||
				ld.Mode >= LoadSyntax && isRoot[pkg.ID] ||
				pkg.ExportFile == "" && pkg.PkgPath != "unsafe",
		}
		ld.pkgs[lpkg.ID] = lpkg
		if isRoot[lpkg.ID] {
			initial = append(initial, lpkg)
			lpkg.initial = true
		}
	}

	

	const (
		white = 0 
		grey  = 1 
		black = 2 
	)

	
	
	
	
	
	
	
	
	
	
	var stack []*loaderPackage
	var visit func(lpkg *loaderPackage) bool
	var srcPkgs []*loaderPackage
	visit = func(lpkg *loaderPackage) bool {
		switch lpkg.color {
		case black:
			return lpkg.needsrc
		case grey:
			panic("internal error: grey node")
		}
		lpkg.color = grey
		stack = append(stack, lpkg) 
		stubs := lpkg.Imports       
		lpkg.Imports = make(map[string]*Package, len(stubs))
		for importPath, ipkg := range stubs {
			var importErr error
			imp := ld.pkgs[ipkg.ID]
			if imp == nil {
				
				importErr = fmt.Errorf("missing package: %q", ipkg.ID)
			} else if imp.color == grey {
				importErr = fmt.Errorf("import cycle: %s", stack)
			}
			if importErr != nil {
				if lpkg.importErrors == nil {
					lpkg.importErrors = make(map[string]error)
				}
				lpkg.importErrors[importPath] = importErr
				continue
			}

			if visit(imp) {
				lpkg.needsrc = true
			}
			lpkg.Imports[importPath] = imp.Package
		}
		if lpkg.needsrc {
			srcPkgs = append(srcPkgs, lpkg)
		}
		stack = stack[:len(stack)-1] 
		lpkg.color = black

		return lpkg.needsrc
	}

	if ld.Mode < LoadImports {
		
		for _, lpkg := range initial {
			lpkg.Imports = nil
		}
	} else {
		
		for _, lpkg := range initial {
			visit(lpkg)
		}
	}
	for _, lpkg := range srcPkgs {
		
		
		for _, ipkg := range lpkg.Imports {
			imp := ld.pkgs[ipkg.ID]
			imp.needtypes = true
		}
	}
	
	
	if ld.Mode >= LoadTypes {
		var wg sync.WaitGroup
		for _, lpkg := range initial {
			wg.Add(1)
			go func(lpkg *loaderPackage) {
				ld.loadRecursive(lpkg)
				wg.Done()
			}(lpkg)
		}
		wg.Wait()
	}

	result := make([]*Package, len(initial))
	for i, lpkg := range initial {
		result[i] = lpkg.Package
	}
	return result, nil
}





func (ld *loader) loadRecursive(lpkg *loaderPackage) {
	lpkg.loadOnce.Do(func() {
		
		var wg sync.WaitGroup
		for _, ipkg := range lpkg.Imports {
			imp := ld.pkgs[ipkg.ID]
			wg.Add(1)
			go func(imp *loaderPackage) {
				ld.loadRecursive(imp)
				wg.Done()
			}(imp)
		}
		wg.Wait()

		ld.loadPackage(lpkg)
	})
}





func (ld *loader) loadPackage(lpkg *loaderPackage) {
	if lpkg.PkgPath == "unsafe" {
		
		lpkg.Types = types.Unsafe
		lpkg.Fset = ld.Fset
		lpkg.Syntax = []*ast.File{}
		lpkg.TypesInfo = new(types.Info)
		return
	}

	
	
	
	lpkg.Types = types.NewPackage(lpkg.PkgPath, lpkg.Name)
	lpkg.Fset = ld.Fset

	
	
	
	
	
	
	if !lpkg.needtypes {
		return
	}
	if !lpkg.needsrc {
		ld.loadFromExportData(lpkg)
		return 
	}

	appendError := func(err error) {
		
		var errs []Error
		switch err := err.(type) {
		case Error:
			
			errs = append(errs, err)

		case *os.PathError:
			
			errs = append(errs, Error{
				Pos:  err.Path + ":1",
				Msg:  err.Err.Error(),
				Kind: ParseError,
			})

		case scanner.ErrorList:
			
			for _, err := range err {
				errs = append(errs, Error{
					Pos:  err.Pos.String(),
					Msg:  err.Msg,
					Kind: ParseError,
				})
			}

		case types.Error:
			
			errs = append(errs, Error{
				Pos:  err.Fset.Position(err.Pos).String(),
				Msg:  err.Msg,
				Kind: TypeError,
			})

		default:
			
			errs = append(errs, Error{
				Pos:  "-",
				Msg:  err.Error(),
				Kind: UnknownError,
			})

			
			log.Printf("internal error: error %q (%T) without position", err, err)
		}

		lpkg.Errors = append(lpkg.Errors, errs...)
	}

	files, errs := ld.parseFiles(lpkg.CompiledGoFiles)
	for _, err := range errs {
		appendError(err)
	}

	lpkg.Syntax = files

	lpkg.TypesInfo = &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Scopes:     make(map[ast.Node]*types.Scope),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	importer := importerFunc(func(path string) (*types.Package, error) {
		if path == "unsafe" {
			return types.Unsafe, nil
		}

		
		ipkg := lpkg.Imports[path]
		if ipkg == nil {
			if err := lpkg.importErrors[path]; err != nil {
				return nil, err
			}
			
			
			
			
			return nil, fmt.Errorf("no metadata for %s", path)
		}

		if ipkg.Types != nil && ipkg.Types.Complete() {
			return ipkg.Types, nil
		}
		log.Fatalf("internal error: nil Pkg importing %q from %q", path, lpkg)
		panic("unreachable")
	})

	
	
	goarch := runtime.GOARCH
	const goarchPrefix = "GOARCH="
	for _, e := range ld.Config.Env {
		if strings.HasPrefix(e, goarchPrefix) {
			goarch = e[len(goarchPrefix):]
		}
	}
	sizes := types.SizesFor("gc", goarch)

	
	tc := &types.Config{
		Importer: importer,

		
		
		
		IgnoreFuncBodies: ld.Mode < LoadAllSyntax && !lpkg.initial,

		Error: appendError,
		Sizes: sizes,
	}
	types.NewChecker(tc, ld.Fset, lpkg.Types, lpkg.TypesInfo).Files(lpkg.Syntax)

	lpkg.importErrors = nil 

	
	
	
	
	
	
	
	
	if tc.FakeImportC {
	outer:
		for _, f := range lpkg.Syntax {
			for _, imp := range f.Imports {
				if imp.Path.Value == `"C"` {
					err := types.Error{Fset: ld.Fset, Pos: imp.Pos(), Msg: `import "C" ignored`}
					appendError(err)
					break outer
				}
			}
		}
	}

	
	illTyped := len(lpkg.Errors) > 0
	if !illTyped {
		for _, imp := range lpkg.Imports {
			if imp.IllTyped {
				illTyped = true
				break
			}
		}
	}
	lpkg.IllTyped = illTyped
}



type importerFunc func(path string) (*types.Package, error)

func (f importerFunc) Import(path string) (*types.Package, error) { return f(path) }



var ioLimit = make(chan bool, 20)








func (ld *loader) parseFiles(filenames []string) ([]*ast.File, []error) {
	var wg sync.WaitGroup
	n := len(filenames)
	parsed := make([]*ast.File, n)
	errors := make([]error, n)
	for i, file := range filenames {
		wg.Add(1)
		go func(i int, filename string) {
			ioLimit <- true 
			
			var src []byte
			for f, contents := range ld.Config.Overlay {
				if sameFile(f, filename) {
					src = contents
				}
			}
			var err error
			if src == nil {
				src, err = ioutil.ReadFile(filename)
			}
			if err != nil {
				parsed[i], errors[i] = nil, err
			} else {
				parsed[i], errors[i] = ld.ParseFile(ld.Fset, filename, src)
			}
			<-ioLimit 
			wg.Done()
		}(i, file)
	}
	wg.Wait()

	
	var o int
	for _, f := range parsed {
		if f != nil {
			parsed[o] = f
			o++
		}
	}
	parsed = parsed[:o]

	o = 0
	for _, err := range errors {
		if err != nil {
			errors[o] = err
			o++
		}
	}
	errors = errors[:o]

	return parsed, errors
}




func sameFile(x, y string) bool {
	if filepath.Base(x) == filepath.Base(y) { 
		if xi, err := os.Stat(x); err == nil {
			if yi, err := os.Stat(y); err == nil {
				return os.SameFile(xi, yi)
			}
		}
	}
	return false
}



func (ld *loader) loadFromExportData(lpkg *loaderPackage) (*types.Package, error) {
	if lpkg.PkgPath == "" {
		log.Fatalf("internal error: Package %s has no PkgPath", lpkg)
	}

	
	
	
	
	
	
	
	
	
	
	
	
	ld.exportMu.Lock()
	defer ld.exportMu.Unlock()

	if tpkg := lpkg.Types; tpkg != nil && tpkg.Complete() {
		return tpkg, nil 
	}

	lpkg.IllTyped = true 

	if lpkg.ExportFile == "" {
		
		return nil, fmt.Errorf("no export data file")
	}
	f, err := os.Open(lpkg.ExportFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	
	
	
	
	
	
	r, err := gcexportdata.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %v", lpkg.ExportFile, err)
	}

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	view := make(map[string]*types.Package) 
	seen := make(map[*loaderPackage]bool)   
	var visit func(pkgs map[string]*Package)
	visit = func(pkgs map[string]*Package) {
		for _, p := range pkgs {
			lpkg := ld.pkgs[p.ID]
			if !seen[lpkg] {
				seen[lpkg] = true
				view[lpkg.PkgPath] = lpkg.Types
				visit(lpkg.Imports)
			}
		}
	}
	visit(lpkg.Imports)

	viewLen := len(view) + 1 
	
	
	tpkg, err := gcexportdata.Read(r, ld.Fset, view, lpkg.PkgPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %v", lpkg.ExportFile, err)
	}
	if viewLen != len(view) {
		log.Fatalf("Unexpected package creation during export data loading")
	}

	lpkg.Types = tpkg
	lpkg.IllTyped = false

	return tpkg, nil
}

func usesExportData(cfg *Config) bool {
	return LoadTypes <= cfg.Mode && cfg.Mode < LoadAllSyntax
}
