





package gopathwalk

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"golang.org/x/tools/internal/fastwalk"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)


type Options struct {
	Debug          bool 
	ModulesEnabled bool 
}


type RootType int

const (
	RootUnknown RootType = iota
	RootGOROOT
	RootGOPATH
	RootCurrentModule
	RootModuleCache
)


type Root struct {
	Path string
	Type RootType
}


func SrcDirsRoots() []Root {
	var roots []Root
	roots = append(roots, Root{filepath.Join(build.Default.GOROOT, "src"), RootGOROOT})
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		roots = append(roots, Root{filepath.Join(p, "src"), RootGOPATH})
	}
	return roots
}





func Walk(roots []Root, add func(root Root, dir string), opts Options) {
	for _, root := range roots {
		walkDir(root, add, opts)
	}
}

func walkDir(root Root, add func(Root, string), opts Options) {
	if _, err := os.Stat(root.Path); os.IsNotExist(err) {
		if opts.Debug {
			log.Printf("skipping nonexistant directory: %v", root.Path)
		}
		return
	}
	if opts.Debug {
		log.Printf("scanning %s", root.Path)
	}
	w := &walker{
		root: root,
		add:  add,
		opts: opts,
	}
	w.init()
	if err := fastwalk.Walk(root.Path, w.walk); err != nil {
		log.Printf("gopathwalk: scanning directory %v: %v", root.Path, err)
	}

	if opts.Debug {
		log.Printf("scanned %s", root.Path)
	}
}


type walker struct {
	root Root               
	add  func(Root, string) 
	opts Options            

	ignoredDirs []os.FileInfo 
}


func (w *walker) init() {
	var ignoredPaths []string
	if w.root.Type == RootModuleCache {
		ignoredPaths = []string{"cache"}
	}
	if !w.opts.ModulesEnabled && w.root.Type == RootGOPATH {
		ignoredPaths = w.getIgnoredDirs(w.root.Path)
		ignoredPaths = append(ignoredPaths, "v", "mod")
	}

	for _, p := range ignoredPaths {
		full := filepath.Join(w.root.Path, p)
		if fi, err := os.Stat(full); err == nil {
			w.ignoredDirs = append(w.ignoredDirs, fi)
			if w.opts.Debug {
				log.Printf("Directory added to ignore list: %s", full)
			}
		} else if w.opts.Debug {
			log.Printf("Error statting ignored directory: %v", err)
		}
	}
}




func (w *walker) getIgnoredDirs(path string) []string {
	file := filepath.Join(path, ".goimportsignore")
	slurp, err := ioutil.ReadFile(file)
	if w.opts.Debug {
		if err != nil {
			log.Print(err)
		} else {
			log.Printf("Read %s", file)
		}
	}
	if err != nil {
		return nil
	}

	var ignoredDirs []string
	bs := bufio.NewScanner(bytes.NewReader(slurp))
	for bs.Scan() {
		line := strings.TrimSpace(bs.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignoredDirs = append(ignoredDirs, line)
	}
	return ignoredDirs
}

func (w *walker) shouldSkipDir(fi os.FileInfo) bool {
	for _, ignoredDir := range w.ignoredDirs {
		if os.SameFile(fi, ignoredDir) {
			return true
		}
	}
	return false
}

func (w *walker) walk(path string, typ os.FileMode) error {
	dir := filepath.Dir(path)
	if typ.IsRegular() {
		if dir == w.root.Path {
			
			
			return fastwalk.SkipFiles
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		w.add(w.root, dir)
		return fastwalk.SkipFiles
	}
	if typ == os.ModeDir {
		base := filepath.Base(path)
		if base == "" || base[0] == '.' || base[0] == '_' ||
			base == "testdata" || (!w.opts.ModulesEnabled && base == "node_modules") {
			return filepath.SkipDir
		}
		fi, err := os.Lstat(path)
		if err == nil && w.shouldSkipDir(fi) {
			return filepath.SkipDir
		}
		return nil
	}
	if typ == os.ModeSymlink {
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".#") {
			
			return nil
		}
		fi, err := os.Lstat(path)
		if err != nil {
			
			return nil
		}
		if w.shouldTraverse(dir, fi) {
			return fastwalk.TraverseLink
		}
	}
	return nil
}




func (w *walker) shouldTraverse(dir string, fi os.FileInfo) bool {
	path := filepath.Join(dir, fi.Name())
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false
	}
	ts, err := os.Stat(target)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	if !ts.IsDir() {
		return false
	}
	if w.shouldSkipDir(ts) {
		return false
	}
	
	
	for {
		parent := filepath.Dir(path)
		if parent == path {
			
			
			return true
		}
		parentInfo, err := os.Stat(parent)
		if err != nil {
			return false
		}
		if os.SameFile(ts, parentInfo) {
			
			return false
		}
		path = parent
	}

}
