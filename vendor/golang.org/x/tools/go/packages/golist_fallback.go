



package packages

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/internal/cgo"
)










func golistDriverFallback(cfg *Config, words ...string) (*driverResponse, error) {
	
	
	
	
	words = cleanAbsPaths(cfg, words)

	original, deps, err := getDeps(cfg, words...)
	if err != nil {
		return nil, err
	}

	var tmpdir string 
	var needsTestVariant []struct {
		pkg, xtestPkg *Package
	}

	var response driverResponse
	allPkgs := make(map[string]bool)
	addPackage := func(p *jsonPackage) {
		id := p.ImportPath

		if allPkgs[id] {
			return
		}
		allPkgs[id] = true

		isRoot := original[id] != nil
		pkgpath := id

		if pkgpath == "unsafe" {
			p.GoFiles = nil 
		}

		importMap := func(importlist []string) map[string]*Package {
			importMap := make(map[string]*Package)
			for _, id := range importlist {

				if id == "C" {
					for _, path := range []string{"unsafe", "syscall", "runtime/cgo"} {
						if pkgpath != path && importMap[path] == nil {
							importMap[path] = &Package{ID: path}
						}
					}
					continue
				}
				importMap[vendorlessPath(id)] = &Package{ID: id}
			}
			return importMap
		}
		compiledGoFiles := absJoin(p.Dir, p.GoFiles)
		
		var cgoErrors []error
		var outdir string
		getOutdir := func() (string, error) {
			if outdir != "" {
				return outdir, nil
			}
			if tmpdir == "" {
				if tmpdir, err = ioutil.TempDir("", "gopackages"); err != nil {
					return "", err
				}
			}
			
			
			
			
			outdir = filepath.Join(tmpdir, "go-build", strings.Replace(p.ImportPath, "/", "_", -1))
			if err := os.MkdirAll(outdir, 0755); err != nil {
				outdir = ""
				return "", err
			}
			return outdir, nil
		}
		processCgo := func() bool {
			
			
			outdir, err := getOutdir()
			if err != nil {
				cgoErrors = append(cgoErrors, err)
				return false
			}
			files, _, err := runCgo(p.Dir, outdir, cfg.Env)
			if err != nil {
				cgoErrors = append(cgoErrors, err)
				return false
			}
			compiledGoFiles = append(compiledGoFiles, files...)
			return true
		}
		if len(p.CgoFiles) == 0 || !processCgo() {
			compiledGoFiles = append(compiledGoFiles, absJoin(p.Dir, p.CgoFiles)...) 
		}
		if isRoot {
			response.Roots = append(response.Roots, id)
		}
		pkg := &Package{
			ID:              id,
			Name:            p.Name,
			GoFiles:         absJoin(p.Dir, p.GoFiles, p.CgoFiles),
			CompiledGoFiles: compiledGoFiles,
			OtherFiles:      absJoin(p.Dir, otherFiles(p)...),
			PkgPath:         pkgpath,
			Imports:         importMap(p.Imports),
			
		}
		if p.Error != nil {
			pkg.Errors = append(pkg.Errors, Error{
				Pos: p.Error.Pos,
				Msg: p.Error.Err,
			})
		}
		response.Packages = append(response.Packages, pkg)
		if cfg.Tests && isRoot {
			testID := fmt.Sprintf("%s [%s.test]", id, id)
			if len(p.TestGoFiles) > 0 || len(p.XTestGoFiles) > 0 {
				response.Roots = append(response.Roots, testID)
				testPkg := &Package{
					ID:              testID,
					Name:            p.Name,
					GoFiles:         absJoin(p.Dir, p.GoFiles, p.CgoFiles, p.TestGoFiles),
					CompiledGoFiles: append(compiledGoFiles, absJoin(p.Dir, p.TestGoFiles)...),
					OtherFiles:      absJoin(p.Dir, otherFiles(p)...),
					PkgPath:         pkgpath,
					Imports:         importMap(append(p.Imports, p.TestImports...)),
					
				}
				response.Packages = append(response.Packages, testPkg)
				var xtestPkg *Package
				if len(p.XTestGoFiles) > 0 {
					xtestID := fmt.Sprintf("%s_test [%s.test]", id, id)
					response.Roots = append(response.Roots, xtestID)
					
					
					
					
					xtestPkg = &Package{
						ID:              xtestID,
						Name:            p.Name + "_test",
						GoFiles:         absJoin(p.Dir, p.XTestGoFiles),
						CompiledGoFiles: absJoin(p.Dir, p.XTestGoFiles),
						PkgPath:         pkgpath + "_test",
						Imports:         importMap(p.XTestImports),
					}
					
					
					
					needsTestVariant = append(needsTestVariant, struct{ pkg, xtestPkg *Package }{pkg, xtestPkg})
					response.Packages = append(response.Packages, xtestPkg)
				}
				
				testmainID := id + ".test"
				response.Roots = append(response.Roots, testmainID)
				imports := map[string]*Package{}
				imports[testPkg.PkgPath] = &Package{ID: testPkg.ID}
				if xtestPkg != nil {
					imports[xtestPkg.PkgPath] = &Package{ID: xtestPkg.ID}
				}
				testmainPkg := &Package{
					ID:      testmainID,
					Name:    "main",
					PkgPath: testmainID,
					Imports: imports,
				}
				response.Packages = append(response.Packages, testmainPkg)
				outdir, err := getOutdir()
				if err != nil {
					testmainPkg.Errors = append(testmainPkg.Errors, Error{
						Pos:  "-",
						Msg:  fmt.Sprintf("failed to generate testmain: %v", err),
						Kind: ListError,
					})
					return
				}
				testmain := filepath.Join(outdir, "testmain.go")
				extraimports, extradeps, err := generateTestmain(testmain, testPkg, xtestPkg)
				if err != nil {
					testmainPkg.Errors = append(testmainPkg.Errors, Error{
						Pos:  "-",
						Msg:  fmt.Sprintf("failed to generate testmain: %v", err),
						Kind: ListError,
					})
				}
				deps = append(deps, extradeps...)
				for _, imp := range extraimports { 
					imports[imp] = &Package{ID: imp}
				}
				testmainPkg.GoFiles = []string{testmain}
				testmainPkg.CompiledGoFiles = []string{testmain}
			}
		}
	}

	for _, pkg := range original {
		addPackage(pkg)
	}
	if cfg.Mode < LoadImports || len(deps) == 0 {
		return &response, nil
	}

	buf, err := invokeGo(cfg, golistArgsFallback(cfg, deps)...)
	if err != nil {
		return nil, err
	}

	
	for dec := json.NewDecoder(buf); dec.More(); {
		p := new(jsonPackage)
		if err := dec.Decode(p); err != nil {
			return nil, fmt.Errorf("JSON decoding failed: %v", err)
		}

		addPackage(p)
	}

	for _, v := range needsTestVariant {
		createTestVariants(&response, v.pkg, v.xtestPkg)
	}

	
	sort.SliceStable(response.Packages, func(i, j int) bool {
		return response.Packages[i].PkgPath < response.Packages[j].PkgPath
	})

	return &response, nil
}

func createTestVariants(response *driverResponse, pkgUnderTest, xtestPkg *Package) {
	allPkgs := make(map[string]*Package)
	for _, pkg := range response.Packages {
		allPkgs[pkg.ID] = pkg
	}
	needsTestVariant := make(map[string]bool)
	needsTestVariant[pkgUnderTest.ID] = true
	var needsVariantRec func(p *Package) bool
	needsVariantRec = func(p *Package) bool {
		if needsTestVariant[p.ID] {
			return true
		}
		for _, imp := range p.Imports {
			if needsVariantRec(allPkgs[imp.ID]) {
				
				
				
				needsTestVariant[p.ID] = true
			}
		}
		if !needsTestVariant[p.ID] {
			return false
		}
		
		
		testVariant := *p
		testVariant.ID = fmt.Sprintf("%s [%s.test]", p.ID, pkgUnderTest.ID)
		testVariant.Imports = make(map[string]*Package)
		for imp, pkg := range p.Imports {
			testVariant.Imports[imp] = pkg
			if needsTestVariant[pkg.ID] {
				testVariant.Imports[imp] = &Package{ID: fmt.Sprintf("%s [%s.test]", pkg.ID, pkgUnderTest.ID)}
			}
		}
		response.Packages = append(response.Packages, &testVariant)
		return needsTestVariant[p.ID]
	}
	
	for imp, pkg := range xtestPkg.Imports {
		if allPkgs[pkg.ID] == nil {
			fmt.Printf("for %s: package %s doesn't exist\n", xtestPkg.ID, pkg.ID)
		}
		if needsVariantRec(allPkgs[pkg.ID]) {
			xtestPkg.Imports[imp] = &Package{ID: fmt.Sprintf("%s [%s.test]", pkg.ID, pkgUnderTest.ID)}
		}
	}
}




func cleanAbsPaths(cfg *Config, words []string) []string {
	var searchpaths []string
	var cleaned = make([]string, len(words))
	for i := range cleaned {
		cleaned[i] = words[i]
		
		
		
		if !filepath.IsAbs(cleaned[i]) || strings.HasSuffix(cleaned[i], ".go") {
			continue
		}
		
		if searchpaths == nil {
			cmd := exec.Command("go", "env", "GOPATH", "GOROOT")
			cmd.Env = cfg.Env
			out, err := cmd.Output()
			if err != nil {
				searchpaths = []string{}
				continue 
			}
			lines := strings.Split(string(out), "\n")
			if len(lines) != 3 || lines[0] == "" || lines[1] == "" || lines[2] != "" {
				continue 
			}
			
			for _, path := range filepath.SplitList(lines[0]) {
				searchpaths = append(searchpaths, filepath.Join(path, "src"))
			}
			
			searchpaths = append(searchpaths, filepath.Join(lines[1], "src"))
		}
		for _, sp := range searchpaths {
			if strings.HasPrefix(cleaned[i], sp) {
				cleaned[i] = strings.TrimPrefix(cleaned[i], sp)
				cleaned[i] = strings.TrimLeft(cleaned[i], string(filepath.Separator))
			}
		}
	}
	return cleaned
}




func vendorlessPath(ipath string) string {
	
	if i := strings.LastIndex(ipath, "/vendor/"); i >= 0 {
		return ipath[i+len("/vendor/"):]
	}
	if strings.HasPrefix(ipath, "vendor/") {
		return ipath[len("vendor/"):]
	}
	return ipath
}


func getDeps(cfg *Config, words ...string) (originalSet map[string]*jsonPackage, deps []string, err error) {
	buf, err := invokeGo(cfg, golistArgsFallback(cfg, words)...)
	if err != nil {
		return nil, nil, err
	}

	depsSet := make(map[string]bool)
	originalSet = make(map[string]*jsonPackage)
	var testImports []string

	
	for dec := json.NewDecoder(buf); dec.More(); {
		p := new(jsonPackage)
		if err := dec.Decode(p); err != nil {
			return nil, nil, fmt.Errorf("JSON decoding failed: %v", err)
		}

		originalSet[p.ImportPath] = p
		for _, dep := range p.Deps {
			depsSet[dep] = true
		}
		if cfg.Tests {
			
			pkgTestImports := append(p.TestImports, p.XTestImports...)
			for _, imp := range pkgTestImports {
				if depsSet[imp] {
					continue
				}
				depsSet[imp] = true
				testImports = append(testImports, imp)
			}
		}
	}
	
	if len(testImports) > 0 {
		buf, err = invokeGo(cfg, golistArgsFallback(cfg, testImports)...)
		if err != nil {
			return nil, nil, err
		}
		
		for dec := json.NewDecoder(buf); dec.More(); {
			p := new(jsonPackage)
			if err := dec.Decode(p); err != nil {
				return nil, nil, fmt.Errorf("JSON decoding failed: %v", err)
			}
			for _, dep := range p.Deps {
				depsSet[dep] = true
			}
		}
	}

	for orig := range originalSet {
		delete(depsSet, orig)
	}

	deps = make([]string, 0, len(depsSet))
	for dep := range depsSet {
		deps = append(deps, dep)
	}
	sort.Strings(deps) 
	return originalSet, deps, nil
}

func golistArgsFallback(cfg *Config, words []string) []string {
	fullargs := []string{"list", "-e", "-json"}
	fullargs = append(fullargs, cfg.BuildFlags...)
	fullargs = append(fullargs, "--")
	fullargs = append(fullargs, words...)
	return fullargs
}

func runCgo(pkgdir, tmpdir string, env []string) (files, displayfiles []string, err error) {
	
	
	
	
	
	bp, err := build.ImportDir(pkgdir, build.ImportMode(0))
	if err != nil {
		return nil, nil, err
	}
	for _, ev := range env {
		if v := strings.TrimPrefix(ev, "CGO_CPPFLAGS"); v != ev {
			bp.CgoCPPFLAGS = append(bp.CgoCPPFLAGS, strings.Fields(v)...)
		} else if v := strings.TrimPrefix(ev, "CGO_CFLAGS"); v != ev {
			bp.CgoCFLAGS = append(bp.CgoCFLAGS, strings.Fields(v)...)
		} else if v := strings.TrimPrefix(ev, "CGO_CXXFLAGS"); v != ev {
			bp.CgoCXXFLAGS = append(bp.CgoCXXFLAGS, strings.Fields(v)...)
		} else if v := strings.TrimPrefix(ev, "CGO_LDFLAGS"); v != ev {
			bp.CgoLDFLAGS = append(bp.CgoLDFLAGS, strings.Fields(v)...)
		}
	}
	return cgo.Run(bp, pkgdir, tmpdir, true)
}
