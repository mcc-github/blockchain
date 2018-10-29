package main



import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	buildContext = build.Default
	goroot       = filepath.Clean(runtime.GOROOT())
	gorootSrc    = filepath.Join(goroot, "src")
)



func importPathsNoDotExpansion(args []string) []string {
	if len(args) == 0 {
		return []string{"."}
	}
	var out []string
	for _, a := range args {
		
		
		
		if filepath.Separator == '\\' {
			a = strings.Replace(a, `\`, `/`, -1)
		}

		
		if strings.HasPrefix(a, "./") {
			a = "./" + path.Clean(a)
			if a == "./." {
				a = "."
			}
		} else {
			a = path.Clean(a)
		}
		if a == "all" || a == "std" {
			out = append(out, allPackages(a)...)
			continue
		}
		out = append(out, a)
	}
	return out
}


func importPaths(args []string) []string {
	args = importPathsNoDotExpansion(args)
	var out []string
	for _, a := range args {
		if strings.Contains(a, "...") {
			if build.IsLocalImport(a) {
				out = append(out, allPackagesInFS(a)...)
			} else {
				out = append(out, allPackages(a)...)
			}
			continue
		}
		out = append(out, a)
	}
	return out
}





func matchPattern(pattern string) func(name string) bool {
	re := regexp.QuoteMeta(pattern)
	re = strings.Replace(re, `\.\.\.`, `.*`, -1)
	
	if strings.HasSuffix(re, `/.*`) {
		re = re[:len(re)-len(`/.*`)] + `(/.*)?`
	}
	reg := regexp.MustCompile(`^` + re + `$`)
	return func(name string) bool {
		return reg.MatchString(name)
	}
}



func hasPathPrefix(s, prefix string) bool {
	switch {
	default:
		return false
	case len(s) == len(prefix):
		return s == prefix
	case len(s) > len(prefix):
		if prefix != "" && prefix[len(prefix)-1] == '/' {
			return strings.HasPrefix(s, prefix)
		}
		return s[len(prefix)] == '/' && s[:len(prefix)] == prefix
	}
}




func treeCanMatchPattern(pattern string) func(name string) bool {
	wildCard := false
	if i := strings.Index(pattern, "..."); i >= 0 {
		wildCard = true
		pattern = pattern[:i]
	}
	return func(name string) bool {
		return len(name) <= len(pattern) && hasPathPrefix(pattern, name) ||
			wildCard && strings.HasPrefix(name, pattern)
	}
}





func allPackages(pattern string) []string {
	pkgs := matchPackages(pattern)
	if len(pkgs) == 0 {
		fmt.Fprintf(os.Stderr, "warning: %q matched no packages\n", pattern)
	}
	return pkgs
}

func matchPackages(pattern string) []string {
	match := func(string) bool { return true }
	treeCanMatch := func(string) bool { return true }
	if pattern != "all" && pattern != "std" {
		match = matchPattern(pattern)
		treeCanMatch = treeCanMatchPattern(pattern)
	}

	have := map[string]bool{
		"builtin": true, 
	}
	if !buildContext.CgoEnabled {
		have["runtime/cgo"] = true 
	}
	var pkgs []string

	
	cmd := filepath.Join(goroot, "src/cmd") + string(filepath.Separator)
	filepath.Walk(cmd, func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() || path == cmd {
			return nil
		}
		name := path[len(cmd):]
		if !treeCanMatch(name) {
			return filepath.SkipDir
		}
		
		if strings.Contains(name, string(filepath.Separator)) {
			return filepath.SkipDir
		}

		
		name = "cmd/" + name
		if have[name] {
			return nil
		}
		have[name] = true
		if !match(name) {
			return nil
		}
		_, err = buildContext.ImportDir(path, 0)
		if err != nil {
			if _, noGo := err.(*build.NoGoError); !noGo {
				log.Print(err)
			}
			return nil
		}
		pkgs = append(pkgs, name)
		return nil
	})

	for _, src := range buildContext.SrcDirs() {
		if (pattern == "std" || pattern == "cmd") && src != gorootSrc {
			continue
		}
		src = filepath.Clean(src) + string(filepath.Separator)
		root := src
		if pattern == "cmd" {
			root += "cmd" + string(filepath.Separator)
		}
		filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
			if err != nil || !fi.IsDir() || path == src {
				return nil
			}

			
			_, elem := filepath.Split(path)
			if strings.HasPrefix(elem, ".") || strings.HasPrefix(elem, "_") || elem == "testdata" {
				return filepath.SkipDir
			}

			name := filepath.ToSlash(path[len(src):])
			if pattern == "std" && (strings.Contains(name, ".") || name == "cmd") {
				
				
				return filepath.SkipDir
			}
			if !treeCanMatch(name) {
				return filepath.SkipDir
			}
			if have[name] {
				return nil
			}
			have[name] = true
			if !match(name) {
				return nil
			}
			_, err = buildContext.ImportDir(path, 0)
			if err != nil {
				if _, noGo := err.(*build.NoGoError); noGo {
					return nil
				}
			}
			pkgs = append(pkgs, name)
			return nil
		})
	}
	return pkgs
}




func allPackagesInFS(pattern string) []string {
	pkgs := matchPackagesInFS(pattern)
	if len(pkgs) == 0 {
		fmt.Fprintf(os.Stderr, "warning: %q matched no packages\n", pattern)
	}
	return pkgs
}

func matchPackagesInFS(pattern string) []string {
	
	
	
	
	i := strings.Index(pattern, "...")
	dir, _ := path.Split(pattern[:i])

	
	
	
	
	prefix := ""
	if strings.HasPrefix(pattern, "./") {
		prefix = "./"
	}
	match := matchPattern(pattern)

	var pkgs []string
	filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() {
			return nil
		}
		if path == dir {
			
			
			
			
			
			
			
			
			path = filepath.Clean(path)
		}

		
		_, elem := filepath.Split(path)
		dot := strings.HasPrefix(elem, ".") && elem != "." && elem != ".."
		if dot || strings.HasPrefix(elem, "_") || elem == "testdata" {
			return filepath.SkipDir
		}

		name := prefix + filepath.ToSlash(path)
		if !match(name) {
			return nil
		}
		if _, err = build.ImportDir(path, 0); err != nil {
			if _, noGo := err.(*build.NoGoError); !noGo {
				log.Print(err)
			}
			return nil
		}
		pkgs = append(pkgs, name)
		return nil
	})
	return pkgs
}
