



package imports

import (
	"context"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/internal/gopathwalk"
)


var Debug = false




var LocalPrefix string

func localPrefixes() []string {
	if LocalPrefix != "" {
		return strings.Split(LocalPrefix, ",")
	}
	return nil
}



var importToGroup = []func(importPath string) (num int, ok bool){
	func(importPath string) (num int, ok bool) {
		for _, p := range localPrefixes() {
			if strings.HasPrefix(importPath, p) || strings.TrimSuffix(p, "/") == importPath {
				return 3, true
			}
		}
		return
	},
	func(importPath string) (num int, ok bool) {
		if strings.HasPrefix(importPath, "appengine") {
			return 2, true
		}
		return
	},
	func(importPath string) (num int, ok bool) {
		if strings.Contains(importPath, ".") {
			return 1, true
		}
		return
	},
}

func importGroup(importPath string) int {
	for _, fn := range importToGroup {
		if n, ok := fn(importPath); ok {
			return n
		}
	}
	return 0
}


type importInfo struct {
	Path  string 
	Alias string 
}


type packageInfo struct {
	Globals map[string]bool       
	Imports map[string]importInfo 
	
	
	
	Refs map[string]map[string]bool
}


var dirPackageInfo = dirPackageInfoFile


func dirPackageInfoFile(pkgName, srcDir, filename string) (*packageInfo, error) {
	considerTests := strings.HasSuffix(filename, "_test.go")

	fileBase := filepath.Base(filename)
	packageFileInfos, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}

	info := &packageInfo{
		Globals: make(map[string]bool),
		Imports: make(map[string]importInfo),
		Refs:    make(map[string]map[string]bool),
	}

	visitor := collectReferences(info.Refs)
	for _, fi := range packageFileInfos {
		if fi.Name() == fileBase || !strings.HasSuffix(fi.Name(), ".go") {
			continue
		}
		if !considerTests && strings.HasSuffix(fi.Name(), "_test.go") {
			continue
		}

		fileSet := token.NewFileSet()
		root, err := parser.ParseFile(fileSet, filepath.Join(srcDir, fi.Name()), nil, 0)
		if err != nil {
			continue
		}

		for _, decl := range root.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				info.Globals[valueSpec.Names[0].Name] = true
			}
		}

		for _, imp := range root.Imports {
			impInfo := importInfo{Path: strings.Trim(imp.Path.Value, `"`)}
			name := path.Base(impInfo.Path)
			if imp.Name != nil {
				name = strings.Trim(imp.Name.Name, `"`)
				impInfo.Alias = name
			}
			info.Imports[name] = impInfo
		}

		ast.Walk(visitor, root)
	}
	return info, nil
}



func collectReferences(refs map[string]map[string]bool) visitFn {
	var visitor visitFn
	visitor = func(node ast.Node) ast.Visitor {
		if node == nil {
			return visitor
		}
		switch v := node.(type) {
		case *ast.SelectorExpr:
			xident, ok := v.X.(*ast.Ident)
			if !ok {
				break
			}
			if xident.Obj != nil {
				
				break
			}
			pkgName := xident.Name
			r := refs[pkgName]
			if r == nil {
				r = make(map[string]bool)
				refs[pkgName] = r
			}
			if ast.IsExported(v.Sel.Name) {
				r[v.Sel.Name] = true
			}
		}
		return visitor
	}
	return visitor
}

func fixImports(fset *token.FileSet, f *ast.File, filename string) (added []string, err error) {
	
	
	
	refs := make(map[string]map[string]bool)

	
	decls := make(map[string]*ast.ImportSpec)

	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	srcDir := filepath.Dir(abs)
	if Debug {
		log.Printf("fixImports(filename=%q), abs=%q, srcDir=%q ...", filename, abs, srcDir)
	}

	var packageInfo *packageInfo
	var loadedPackageInfo bool

	
	var visitor visitFn
	visitor = visitFn(func(node ast.Node) ast.Visitor {
		if node == nil {
			return visitor
		}
		switch v := node.(type) {
		case *ast.ImportSpec:
			if v.Name != nil {
				decls[v.Name.Name] = v
				break
			}
			ipath := strings.Trim(v.Path.Value, `"`)
			if ipath == "C" {
				break
			}
			local := importPathToName(ipath, srcDir)
			decls[local] = v
		case *ast.SelectorExpr:
			xident, ok := v.X.(*ast.Ident)
			if !ok {
				break
			}
			if xident.Obj != nil {
				
				break
			}
			pkgName := xident.Name
			if refs[pkgName] == nil {
				refs[pkgName] = make(map[string]bool)
			}
			if !loadedPackageInfo {
				loadedPackageInfo = true
				packageInfo, _ = dirPackageInfo(f.Name.Name, srcDir, filename)
			}
			if decls[pkgName] == nil && (packageInfo == nil || !packageInfo.Globals[pkgName]) {
				refs[pkgName][v.Sel.Name] = true
			}
		}
		return visitor
	})
	ast.Walk(visitor, f)

	
	unusedImport := map[string]string{}
	for pkg, is := range decls {
		if refs[pkg] == nil && pkg != "_" && pkg != "." {
			name := ""
			if is.Name != nil {
				name = is.Name.Name
			}
			unusedImport[strings.Trim(is.Path.Value, `"`)] = name
		}
	}
	for ipath, name := range unusedImport {
		if ipath == "C" {
			
			continue
		}
		astutil.DeleteNamedImport(fset, f, name, ipath)
	}

	for pkgName, symbols := range refs {
		if len(symbols) == 0 {
			
			delete(refs, pkgName)
		}
	}

	
	if len(refs) == 0 {
		return nil, nil
	}

	
	if !loadedPackageInfo {
		packageInfo, _ = dirPackageInfo(f.Name.Name, srcDir, filename)
	}

	
	type result struct {
		ipath string 
		name  string 
	}
	results := make(chan result, len(refs))

	ctx, cancel := context.WithCancel(context.TODO())
	var wg sync.WaitGroup
	defer func() {
		cancel()
		wg.Wait()
	}()
	var (
		firstErr     error
		firstErrOnce sync.Once
	)
	for pkgName, symbols := range refs {
		wg.Add(1)
		go func(pkgName string, symbols map[string]bool) {
			defer wg.Done()

			if packageInfo != nil {
				sibling := packageInfo.Imports[pkgName]
				if sibling.Path != "" {
					refs := packageInfo.Refs[pkgName]
					for symbol := range symbols {
						if refs[symbol] {
							results <- result{ipath: sibling.Path, name: sibling.Alias}
							return
						}
					}
				}
			}

			ipath, rename, err := findImport(ctx, pkgName, symbols, filename)
			if err != nil {
				firstErrOnce.Do(func() {
					firstErr = err
					cancel()
				})
				return
			}

			if ipath == "" {
				return 
			}

			r := result{ipath: ipath}
			if rename {
				r.name = pkgName
			}
			results <- r
			return
		}(pkgName, symbols)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.name != "" {
			astutil.AddNamedImport(fset, f, result.name, result.ipath)
		} else {
			astutil.AddImport(fset, f, result.ipath)
		}
		added = append(added, result.ipath)
	}

	if firstErr != nil {
		return nil, firstErr
	}
	return added, nil
}


var importPathToName func(importPath, srcDir string) (packageName string) = importPathToNameGoPath



func importPathToNameBasic(importPath, srcDir string) (packageName string) {
	base := path.Base(importPath)
	if strings.HasPrefix(base, "v") {
		if _, err := strconv.Atoi(base[1:]); err == nil {
			dir := path.Dir(importPath)
			if dir != "." {
				return path.Base(dir)
			}
		}
	}
	return base
}



func importPathToNameGoPath(importPath, srcDir string) (packageName string) {
	
	if pkg, ok := stdImportPackage[importPath]; ok {
		return pkg
	}

	pkgName, err := importPathToNameGoPathParse(importPath, srcDir)
	if Debug {
		log.Printf("importPathToNameGoPathParse(%q, srcDir=%q) = %q, %v", importPath, srcDir, pkgName, err)
	}
	if err == nil {
		return pkgName
	}
	return importPathToNameBasic(importPath, srcDir)
}





func importPathToNameGoPathParse(importPath, srcDir string) (packageName string, err error) {
	buildPkg, err := build.Import(importPath, srcDir, build.FindOnly)
	if err != nil {
		return "", err
	}
	d, err := os.Open(buildPkg.Dir)
	if err != nil {
		return "", err
	}
	names, err := d.Readdirnames(-1)
	d.Close()
	if err != nil {
		return "", err
	}
	sort.Strings(names) 
	var lastErr error
	var nfile int
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		nfile++
		fullFile := filepath.Join(buildPkg.Dir, name)

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fullFile, nil, parser.PackageClauseOnly)
		if err != nil {
			lastErr = err
			continue
		}
		pkgName := f.Name.Name
		if pkgName == "documentation" {
			
			
			continue
		}
		if pkgName == "main" {
			
			
			continue
		}
		return pkgName, nil
	}
	if lastErr != nil {
		return "", lastErr
	}
	return "", fmt.Errorf("no importable package found in %d Go files", nfile)
}

var stdImportPackage = map[string]string{} 

func init() {
	
	
	for _, pkg := range stdlib {
		if _, ok := stdImportPackage[pkg]; !ok {
			stdImportPackage[pkg] = path.Base(pkg)
		}
	}
}


var (
	
	scanOnce sync.Once
	dirScan  map[string]*pkg 
)

type pkg struct {
	dir             string 
	importPath      string 
	importPathShort string 
}

type pkgDistance struct {
	pkg      *pkg
	distance int 
}



type byDistanceOrImportPathShortLength []pkgDistance

func (s byDistanceOrImportPathShortLength) Len() int { return len(s) }
func (s byDistanceOrImportPathShortLength) Less(i, j int) bool {
	di, dj := s[i].distance, s[j].distance
	if di == -1 {
		return false
	}
	if dj == -1 {
		return true
	}
	if di != dj {
		return di < dj
	}

	vi, vj := s[i].pkg.importPathShort, s[j].pkg.importPathShort
	if len(vi) != len(vj) {
		return len(vi) < len(vj)
	}
	return vi < vj
}
func (s byDistanceOrImportPathShortLength) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func distance(basepath, targetpath string) int {
	p, err := filepath.Rel(basepath, targetpath)
	if err != nil {
		return -1
	}
	if p == "." {
		return 0
	}
	return strings.Count(p, string(filepath.Separator)) + 1
}


func scanGoDirs() map[string]*pkg {
	result := make(map[string]*pkg)
	var mu sync.Mutex

	add := func(root gopathwalk.Root, dir string) {
		mu.Lock()
		defer mu.Unlock()

		if _, dup := result[dir]; dup {
			return
		}
		importpath := filepath.ToSlash(dir[len(root.Path)+len("/"):])
		result[dir] = &pkg{
			importPath:      importpath,
			importPathShort: VendorlessPath(importpath),
			dir:             dir,
		}
	}
	gopathwalk.Walk(gopathwalk.SrcDirsRoots(), add, gopathwalk.Options{Debug: Debug, ModulesEnabled: false})
	return result
}



func VendorlessPath(ipath string) string {
	
	if i := strings.LastIndex(ipath, "/vendor/"); i >= 0 {
		return ipath[i+len("/vendor/"):]
	}
	if strings.HasPrefix(ipath, "vendor/") {
		return ipath[len("vendor/"):]
	}
	return ipath
}



var loadExports func(ctx context.Context, expectPackage, dir string) (map[string]bool, error) = loadExportsGoPath

func loadExportsGoPath(ctx context.Context, expectPackage, dir string) (map[string]bool, error) {
	if Debug {
		log.Printf("loading exports in dir %s (seeking package %s)", dir, expectPackage)
	}
	exports := make(map[string]bool)

	buildCtx := build.Default

	
	
	
	buildCtx.ReadDir = func(dir string) (notTests []os.FileInfo, err error) {
		all, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		notTests = all[:0]
		for _, fi := range all {
			name := fi.Name()
			if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
				notTests = append(notTests, fi)
			}
		}
		return notTests, nil
	}

	files, err := buildCtx.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fset := token.NewFileSet()

	for _, fi := range files {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		match, err := buildCtx.MatchFile(dir, fi.Name())
		if err != nil || !match {
			continue
		}
		fullFile := filepath.Join(dir, fi.Name())
		f, err := parser.ParseFile(fset, fullFile, nil, 0)
		if err != nil {
			if Debug {
				log.Printf("Parsing %s: %v", fullFile, err)
			}
			return nil, err
		}
		pkgName := f.Name.Name
		if pkgName == "documentation" {
			
			
			continue
		}
		if pkgName != expectPackage {
			err := fmt.Errorf("scan of dir %v is not expected package %v (actually %v)", dir, expectPackage, pkgName)
			if Debug {
				log.Print(err)
			}
			return nil, err
		}
		for name := range f.Scope.Objects {
			if ast.IsExported(name) {
				exports[name] = true
			}
		}
	}

	if Debug {
		exportList := make([]string, 0, len(exports))
		for k := range exports {
			exportList = append(exportList, k)
		}
		sort.Strings(exportList)
		log.Printf("loaded exports in dir %v (package %v): %v", dir, expectPackage, strings.Join(exportList, ", "))
	}
	return exports, nil
}













var findImport func(ctx context.Context, pkgName string, symbols map[string]bool, filename string) (foundPkg string, rename bool, err error) = findImportGoPath



func findImportGoPath(ctx context.Context, pkgName string, symbols map[string]bool, filename string) (foundPkg string, rename bool, err error) {
	pkgDir, err := filepath.Abs(filename)
	if err != nil {
		return "", false, err
	}
	pkgDir = filepath.Dir(pkgDir)

	
	
	
	if pkg, ok := findImportStdlib(pkgName, symbols); ok {
		return pkg, false, nil
	}
	if pkgName == "rand" && symbols["Read"] {
		
		
		
		
		
		
		
		return "", false, nil
	}

	
	
	
	

	
	scanOnce.Do(func() { dirScan = scanGoDirs() })

	
	var candidates []pkgDistance
	for _, pkg := range dirScan {
		if pkgIsCandidate(filename, pkgName, pkg) {
			candidates = append(candidates, pkgDistance{
				pkg:      pkg,
				distance: distance(pkgDir, pkg.dir),
			})
		}
	}

	
	
	
	
	sort.Sort(byDistanceOrImportPathShortLength(candidates))
	if Debug {
		for i, c := range candidates {
			log.Printf("%s candidate %d/%d: %v in %v", pkgName, i+1, len(candidates), c.pkg.importPathShort, c.pkg.dir)
		}
	}

	

	rescv := make([]chan *pkg, len(candidates))
	for i := range candidates {
		rescv[i] = make(chan *pkg, 1)
	}
	const maxConcurrentPackageImport = 4
	loadExportsSem := make(chan struct{}, maxConcurrentPackageImport)

	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	defer func() {
		cancel()
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i, c := range candidates {
			select {
			case loadExportsSem <- struct{}{}:
			case <-ctx.Done():
				return
			}

			wg.Add(1)
			go func(c pkgDistance, resc chan<- *pkg) {
				defer func() {
					<-loadExportsSem
					wg.Done()
				}()

				exports, err := loadExports(ctx, pkgName, c.pkg.dir)
				if err != nil {
					resc <- nil
					return
				}

				
				
				for symbol := range symbols {
					if !exports[symbol] {
						resc <- nil
						return
					}
				}
				resc <- c.pkg
			}(c, rescv[i])
		}
	}()

	for _, resc := range rescv {
		pkg := <-resc
		if pkg == nil {
			continue
		}
		
		
		needsRename := path.Base(pkg.importPath) != pkgName
		return pkg.importPathShort, needsRename, nil
	}
	return "", false, nil
}













func pkgIsCandidate(filename, pkgIdent string, pkg *pkg) bool {
	
	if !canUse(filename, pkg.dir) {
		return false
	}

	
	
	
	
	
	
	
	
	
	
	
	
	lastTwo := lastTwoComponents(pkg.importPathShort)
	if strings.Contains(lastTwo, pkgIdent) {
		return true
	}
	if hasHyphenOrUpperASCII(lastTwo) && !hasHyphenOrUpperASCII(pkgIdent) {
		lastTwo = lowerASCIIAndRemoveHyphen(lastTwo)
		if strings.Contains(lastTwo, pkgIdent) {
			return true
		}
	}

	return false
}

func hasHyphenOrUpperASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == '-' || ('A' <= b && b <= 'Z') {
			return true
		}
	}
	return false
}

func lowerASCIIAndRemoveHyphen(s string) (ret string) {
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch {
		case b == '-':
			continue
		case 'A' <= b && b <= 'Z':
			buf = append(buf, b+('a'-'A'))
		default:
			buf = append(buf, b)
		}
	}
	return string(buf)
}



func canUse(filename, dir string) bool {
	
	
	
	
	if !strings.Contains(dir, "vendor") && !strings.Contains(dir, "internal") {
		return true
	}

	dirSlash := filepath.ToSlash(dir)
	if !strings.Contains(dirSlash, "/vendor/") && !strings.Contains(dirSlash, "/internal/") && !strings.HasSuffix(dirSlash, "/internal") {
		return true
	}
	
	
	
	
	
	
	absfile, err := filepath.Abs(filename)
	if err != nil {
		return false
	}
	absdir, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absfile, absdir)
	if err != nil {
		return false
	}
	relSlash := filepath.ToSlash(rel)
	if i := strings.LastIndex(relSlash, "../"); i >= 0 {
		relSlash = relSlash[i+len("../"):]
	}
	return !strings.Contains(relSlash, "/vendor/") && !strings.Contains(relSlash, "/internal/") && !strings.HasSuffix(relSlash, "/internal")
}



func lastTwoComponents(v string) string {
	nslash := 0
	for i := len(v) - 1; i >= 0; i-- {
		if v[i] == '/' || v[i] == '\\' {
			nslash++
			if nslash == 2 {
				return v[i:]
			}
		}
	}
	return v
}

type visitFn func(node ast.Node) ast.Visitor

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	return fn(node)
}

func findImportStdlib(shortPkg string, symbols map[string]bool) (importPath string, ok bool) {
	for symbol := range symbols {
		key := shortPkg + "." + symbol
		path := stdlib[key]
		if path == "" {
			if key == "rand.Read" {
				continue
			}
			return "", false
		}
		if importPath != "" && importPath != path {
			
			return "", false
		}
		importPath = path
	}
	if importPath == "" && shortPkg == "rand" && symbols["Read"] {
		return "crypto/rand", true
	}
	return importPath, importPath != ""
}
