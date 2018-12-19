



package cgo
















































import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)




func ProcessFiles(bp *build.Package, fset *token.FileSet, DisplayPath func(path string) string, mode parser.Mode) ([]*ast.File, error) {
	tmpdir, err := ioutil.TempDir("", strings.Replace(bp.ImportPath, "/", "_", -1)+"_C")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpdir)

	pkgdir := bp.Dir
	if DisplayPath != nil {
		pkgdir = DisplayPath(pkgdir)
	}

	cgoFiles, cgoDisplayFiles, err := Run(bp, pkgdir, tmpdir, false)
	if err != nil {
		return nil, err
	}
	var files []*ast.File
	for i := range cgoFiles {
		rd, err := os.Open(cgoFiles[i])
		if err != nil {
			return nil, err
		}
		display := filepath.Join(bp.Dir, cgoDisplayFiles[i])
		f, err := parser.ParseFile(fset, display, rd, mode)
		rd.Close()
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

var cgoRe = regexp.MustCompile(`[/\\:]`)














func Run(bp *build.Package, pkgdir, tmpdir string, useabs bool) (files, displayFiles []string, err error) {
	cgoCPPFLAGS, _, _, _ := cflags(bp, true)
	_, cgoexeCFLAGS, _, _ := cflags(bp, false)

	if len(bp.CgoPkgConfig) > 0 {
		pcCFLAGS, err := pkgConfigFlags(bp)
		if err != nil {
			return nil, nil, err
		}
		cgoCPPFLAGS = append(cgoCPPFLAGS, pcCFLAGS...)
	}

	
	cgoCPPFLAGS = append(cgoCPPFLAGS, "-I", tmpdir)

	
	files = append(files, filepath.Join(tmpdir, "_cgo_gotypes.go"))
	displayFiles = append(displayFiles, "C")
	for _, fn := range bp.CgoFiles {
		
		f := cgoRe.ReplaceAllString(fn[:len(fn)-len("go")], "_")
		files = append(files, filepath.Join(tmpdir, f+"cgo1.go"))
		displayFiles = append(displayFiles, fn)
	}

	var cgoflags []string
	if bp.Goroot && bp.ImportPath == "runtime/cgo" {
		cgoflags = append(cgoflags, "-import_runtime_cgo=false")
	}
	if bp.Goroot && bp.ImportPath == "runtime/race" || bp.ImportPath == "runtime/cgo" {
		cgoflags = append(cgoflags, "-import_syscall=false")
	}

	var cgoFiles []string = bp.CgoFiles
	if useabs {
		cgoFiles = make([]string, len(bp.CgoFiles))
		for i := range cgoFiles {
			cgoFiles[i] = filepath.Join(pkgdir, bp.CgoFiles[i])
		}
	}

	args := stringList(
		"go", "tool", "cgo", "-objdir", tmpdir, cgoflags, "--",
		cgoCPPFLAGS, cgoexeCFLAGS, cgoFiles,
	)
	if false {
		log.Printf("Running cgo for package %q: %s (dir=%s)", bp.ImportPath, args, pkgdir)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = pkgdir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, fmt.Errorf("cgo failed: %s: %s", args, err)
	}

	return files, displayFiles, nil
}




func cflags(p *build.Package, def bool) (cppflags, cflags, cxxflags, ldflags []string) {
	var defaults string
	if def {
		defaults = "-g -O2"
	}

	cppflags = stringList(envList("CGO_CPPFLAGS", ""), p.CgoCPPFLAGS)
	cflags = stringList(envList("CGO_CFLAGS", defaults), p.CgoCFLAGS)
	cxxflags = stringList(envList("CGO_CXXFLAGS", defaults), p.CgoCXXFLAGS)
	ldflags = stringList(envList("CGO_LDFLAGS", defaults), p.CgoLDFLAGS)
	return
}



func envList(key, def string) []string {
	v := os.Getenv(key)
	if v == "" {
		v = def
	}
	return strings.Fields(v)
}



func stringList(args ...interface{}) []string {
	var x []string
	for _, arg := range args {
		switch arg := arg.(type) {
		case []string:
			x = append(x, arg...)
		case string:
			x = append(x, arg)
		default:
			panic("stringList: invalid argument")
		}
	}
	return x
}
