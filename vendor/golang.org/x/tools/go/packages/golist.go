



package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/tools/internal/gopathwalk"
	"golang.org/x/tools/internal/semver"
)



type goTooOldError struct {
	error
}




func goListDriver(cfg *Config, patterns ...string) (*driverResponse, error) {
	
	var containFiles []string
	var packagesNamed []string
	restPatterns := make([]string, 0, len(patterns))
	
	
extractQueries:
	for _, pattern := range patterns {
		eqidx := strings.Index(pattern, "=")
		if eqidx < 0 {
			restPatterns = append(restPatterns, pattern)
		} else {
			query, value := pattern[:eqidx], pattern[eqidx+len("="):]
			switch query {
			case "file":
				containFiles = append(containFiles, value)
			case "pattern":
				restPatterns = append(restPatterns, value)
			case "name":
				packagesNamed = append(packagesNamed, value)
			case "": 
				restPatterns = append(restPatterns, pattern)
			default:
				for _, rune := range query {
					if rune < 'a' || rune > 'z' { 
						restPatterns = append(restPatterns, pattern)
						continue extractQueries
					}
				}
				
				return nil, fmt.Errorf("invalid query type %q in query pattern %q", query, pattern)
			}
		}
	}
	patterns = restPatterns
	
	
	restPatterns = restPatterns[:0]
	for _, pattern := range patterns {
		if strings.HasPrefix(pattern, "contains:") {
			containFile := strings.TrimPrefix(pattern, "contains:")
			containFiles = append(containFiles, containFile)
		} else {
			restPatterns = append(restPatterns, pattern)
		}
	}
	containFiles = absJoin(cfg.Dir, containFiles)

	
	var listfunc driver
	listfunc = func(cfg *Config, words ...string) (*driverResponse, error) {
		response, err := golistDriverCurrent(cfg, words...)
		if _, ok := err.(goTooOldError); ok {
			listfunc = golistDriverFallback
			return listfunc(cfg, words...)
		}
		listfunc = golistDriverCurrent
		return response, err
	}

	var response *driverResponse
	var err error

	
	if len(restPatterns) > 0 {
		response, err = listfunc(cfg, restPatterns...)
		if err != nil {
			return nil, err
		}
	} else {
		response = &driverResponse{}
	}

	if len(containFiles) == 0 && len(packagesNamed) == 0 {
		return response, nil
	}

	seenPkgs := make(map[string]*Package) 
	for _, pkg := range response.Packages {
		seenPkgs[pkg.ID] = pkg
	}
	addPkg := func(p *Package) {
		if _, ok := seenPkgs[p.ID]; ok {
			return
		}
		seenPkgs[p.ID] = p
		response.Packages = append(response.Packages, p)
	}

	containsResults, err := runContainsQueries(cfg, listfunc, addPkg, containFiles)
	if err != nil {
		return nil, err
	}
	response.Roots = append(response.Roots, containsResults...)

	namedResults, err := runNamedQueries(cfg, listfunc, addPkg, packagesNamed)
	if err != nil {
		return nil, err
	}
	response.Roots = append(response.Roots, namedResults...)

	return response, nil
}

func runContainsQueries(cfg *Config, driver driver, addPkg func(*Package), queries []string) ([]string, error) {
	var results []string
	for _, query := range queries {
		
		fdir := filepath.Dir(query)
		cfg.Dir = fdir
		dirResponse, err := driver(cfg, ".")
		if err != nil {
			return nil, err
		}
		isRoot := make(map[string]bool, len(dirResponse.Roots))
		for _, root := range dirResponse.Roots {
			isRoot[root] = true
		}
		for _, pkg := range dirResponse.Packages {
			
			
			
			
			addPkg(pkg)
			
			if !isRoot[pkg.ID] {
				continue
			}
			for _, pkgFile := range pkg.GoFiles {
				if filepath.Base(query) == filepath.Base(pkgFile) {
					results = append(results, pkg.ID)
					break
				}
			}
		}
	}
	return results, nil
}


var modCacheRegexp = regexp.MustCompile(`(.*)@([^/\\]*)(.*)`)

func runNamedQueries(cfg *Config, driver driver, addPkg func(*Package), queries []string) ([]string, error) {
	
	roots, modulesEnabled, err := roots(cfg)
	if err != nil {
		return nil, err
	}

	
	
	
	var matchesMu sync.Mutex
	var simpleMatches, modCacheMatches []string
	add := func(root gopathwalk.Root, dir string) {
		
		matchesMu.Lock()
		defer matchesMu.Unlock()

		path := dir[len(root.Path)+1:]
		if pathMatchesQueries(path, queries) {
			switch root.Type {
			case gopathwalk.RootModuleCache:
				modCacheMatches = append(modCacheMatches, path)
			case gopathwalk.RootCurrentModule:
				
				
				rel, err := filepath.Rel(cfg.Dir, dir)
				if err != nil {
					
					
					panic(err)
				}
				simpleMatches = append(simpleMatches, "./"+rel)
			case gopathwalk.RootGOPATH, gopathwalk.RootGOROOT:
				simpleMatches = append(simpleMatches, path)
			}
		}
	}
	gopathwalk.Walk(roots, add, gopathwalk.Options{ModulesEnabled: modulesEnabled})

	var results []string
	addResponse := func(r *driverResponse) {
		for _, pkg := range r.Packages {
			addPkg(pkg)
			for _, name := range queries {
				if pkg.Name == name {
					results = append(results, pkg.ID)
					break
				}
			}
		}
	}

	if len(simpleMatches) != 0 {
		resp, err := driver(cfg, simpleMatches...)
		if err != nil {
			return nil, err
		}
		addResponse(resp)
	}

	
	
	
	
	if len(modCacheMatches) != 0 {
		
		
		type modInfo struct {
			mod   string
			major string
		}
		mods := make(map[modInfo]string)
		var imports []string
		for _, modPath := range modCacheMatches {
			matches := modCacheRegexp.FindStringSubmatch(modPath)
			mod, ver := filepath.ToSlash(matches[1]), matches[2]
			importPath := filepath.ToSlash(filepath.Join(matches[1], matches[3]))

			major := semver.Major(ver)
			if prevVer, ok := mods[modInfo{mod, major}]; !ok || semver.Compare(ver, prevVer) > 0 {
				mods[modInfo{mod, major}] = ver
			}

			imports = append(imports, importPath)
		}

		
		var gomod bytes.Buffer
		gomod.WriteString("module modquery\nrequire (\n")
		for mod, version := range mods {
			gomod.WriteString("\t" + mod.mod + " " + version + "\n")
		}
		gomod.WriteString(")\n")

		tmpCfg := *cfg
		var err error
		tmpCfg.Dir, err = ioutil.TempDir("", "gopackages-modquery")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(tmpCfg.Dir)

		if err := ioutil.WriteFile(filepath.Join(tmpCfg.Dir, "go.mod"), gomod.Bytes(), 0777); err != nil {
			return nil, fmt.Errorf("writing go.mod for module cache query: %v", err)
		}

		
		resp, err := driver(&tmpCfg, imports...)
		if err != nil {
			return nil, fmt.Errorf("querying module cache matches: %v", err)
		}
		addResponse(resp)
	}

	return results, nil
}



func roots(cfg *Config) ([]gopathwalk.Root, bool, error) {
	stdout, err := invokeGo(cfg, "env", "GOROOT", "GOPATH", "GOMOD")
	if err != nil {
		return nil, false, err
	}

	fields := strings.Split(stdout.String(), "\n")
	if len(fields) != 4 || len(fields[3]) != 0 {
		return nil, false, fmt.Errorf("go env returned unexpected output: %q", stdout.String())
	}
	goroot, gopath, gomod := fields[0], filepath.SplitList(fields[1]), fields[2]
	modsEnabled := gomod != ""

	var roots []gopathwalk.Root
	
	roots = append(roots, gopathwalk.Root{filepath.Join(goroot, "/src"), gopathwalk.RootGOROOT})
	
	if modsEnabled {
		roots = append(roots, gopathwalk.Root{filepath.Dir(gomod), gopathwalk.RootCurrentModule})
	}
	
	for _, p := range gopath {
		if modsEnabled {
			roots = append(roots, gopathwalk.Root{filepath.Join(p, "/pkg/mod"), gopathwalk.RootModuleCache})
		} else {
			roots = append(roots, gopathwalk.Root{filepath.Join(p, "/src"), gopathwalk.RootGOPATH})
		}
	}

	return roots, modsEnabled, nil
}





func pathMatchesQueries(path string, queries []string) bool {
	lastTwo := lastTwoComponents(path)
	for _, query := range queries {
		if strings.Contains(lastTwo, query) {
			return true
		}
		if hasHyphenOrUpperASCII(lastTwo) && !hasHyphenOrUpperASCII(query) {
			lastTwo = lowerASCIIAndRemoveHyphen(lastTwo)
			if strings.Contains(lastTwo, query) {
				return true
			}
		}
	}
	return false
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



type jsonPackage struct {
	ImportPath      string
	Dir             string
	Name            string
	Export          string
	GoFiles         []string
	CompiledGoFiles []string
	CFiles          []string
	CgoFiles        []string
	CXXFiles        []string
	MFiles          []string
	HFiles          []string
	FFiles          []string
	SFiles          []string
	SwigFiles       []string
	SwigCXXFiles    []string
	SysoFiles       []string
	Imports         []string
	ImportMap       map[string]string
	Deps            []string
	TestGoFiles     []string
	TestImports     []string
	XTestGoFiles    []string
	XTestImports    []string
	ForTest         string 
	DepOnly         bool

	Error *jsonPackageError
}

type jsonPackageError struct {
	ImportStack []string
	Pos         string
	Err         string
}

func otherFiles(p *jsonPackage) [][]string {
	return [][]string{p.CFiles, p.CXXFiles, p.MFiles, p.HFiles, p.FFiles, p.SFiles, p.SwigFiles, p.SwigCXXFiles, p.SysoFiles}
}




func golistDriverCurrent(cfg *Config, words ...string) (*driverResponse, error) {
	
	
	
	
	
	
	
	
	
	

	
	
	buf, err := invokeGo(cfg, golistargs(cfg, words)...)
	if err != nil {
		return nil, err
	}
	
	var response driverResponse
	for dec := json.NewDecoder(buf); dec.More(); {
		p := new(jsonPackage)
		if err := dec.Decode(p); err != nil {
			return nil, fmt.Errorf("JSON decoding failed: %v", err)
		}

		if p.ImportPath == "" {
			
			
			
			
			if p.Error != nil {
				return nil, Error{
					Pos: p.Error.Pos,
					Msg: p.Error.Err,
				}
			}
			return nil, fmt.Errorf("package missing import path: %+v", p)
		}

		pkg := &Package{
			Name:            p.Name,
			ID:              p.ImportPath,
			GoFiles:         absJoin(p.Dir, p.GoFiles, p.CgoFiles),
			CompiledGoFiles: absJoin(p.Dir, p.CompiledGoFiles),
			OtherFiles:      absJoin(p.Dir, otherFiles(p)...),
		}

		
		if i := strings.IndexByte(pkg.ID, ' '); i >= 0 {
			pkg.PkgPath = pkg.ID[:i]
		} else {
			pkg.PkgPath = pkg.ID
		}

		if pkg.PkgPath == "unsafe" {
			pkg.GoFiles = nil 
		}

		
		if p.Dir != "" && !filepath.IsAbs(p.Dir) {
			log.Fatalf("internal error: go list returned non-absolute Package.Dir: %s", p.Dir)
		}

		if p.Export != "" && !filepath.IsAbs(p.Export) {
			pkg.ExportFile = filepath.Join(p.Dir, p.Export)
		} else {
			pkg.ExportFile = p.Export
		}

		
		
		
		
		ids := make(map[string]bool)
		for _, id := range p.Imports {
			ids[id] = true
		}
		pkg.Imports = make(map[string]*Package)
		for path, id := range p.ImportMap {
			pkg.Imports[path] = &Package{ID: id} 
			delete(ids, id)
		}
		for id := range ids {
			if id == "C" {
				continue
			}

			pkg.Imports[id] = &Package{ID: id} 
		}
		if !p.DepOnly {
			response.Roots = append(response.Roots, pkg.ID)
		}

		
		if len(pkg.CompiledGoFiles) == 0 {
			pkg.CompiledGoFiles = pkg.GoFiles
		}

		if p.Error != nil {
			pkg.Errors = append(pkg.Errors, Error{
				Pos: p.Error.Pos,
				Msg: p.Error.Err,
			})
		}

		response.Packages = append(response.Packages, pkg)
	}

	return &response, nil
}


func absJoin(dir string, fileses ...[]string) (res []string) {
	for _, files := range fileses {
		for _, file := range files {
			if !filepath.IsAbs(file) {
				file = filepath.Join(dir, file)
			}
			res = append(res, file)
		}
	}
	return res
}

func golistargs(cfg *Config, words []string) []string {
	fullargs := []string{
		"list", "-e", "-json", "-compiled",
		fmt.Sprintf("-test=%t", cfg.Tests),
		fmt.Sprintf("-export=%t", usesExportData(cfg)),
		fmt.Sprintf("-deps=%t", cfg.Mode >= LoadImports),
	}
	fullargs = append(fullargs, cfg.BuildFlags...)
	fullargs = append(fullargs, "--")
	fullargs = append(fullargs, words...)
	return fullargs
}


func invokeGo(cfg *Config, args ...string) (*bytes.Buffer, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := exec.CommandContext(cfg.Context, "go", args...)
	
	
	
	
	
	
	cmd.Env = append(append([]string{}, cfg.Env...), "PWD="+cfg.Dir)
	cmd.Dir = cfg.Dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			
			
			
			return nil, fmt.Errorf("couldn't exec 'go %v': %s %T", args, err, err)
		}

		
		if strings.Contains(stderr.String(), "flag provided but not defined") {
			return nil, goTooOldError{fmt.Errorf("unsupported version of go: %s: %s", exitErr, stderr)}
		}

		
		
		
		
		if !usesExportData(cfg) {
			return nil, fmt.Errorf("go %v: %s: %s", args, exitErr, stderr)
		}
	}

	
	
	
	
	
	
	if len(stderr.Bytes()) != 0 && os.Getenv("GOPACKAGESPRINTGOLISTERRORS") != "" {
		fmt.Fprintf(os.Stderr, "go %v stderr: <<%s>>\n", args, stderr)
	}

	
	if false {
		fmt.Fprintf(os.Stderr, "go %v stdout: <<%s>>\n", args, stdout)
	}

	return stdout, nil
}
