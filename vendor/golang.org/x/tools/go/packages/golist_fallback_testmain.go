






package packages

import (
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)






func generateTestmain(out string, testPkg, xtestPkg *Package) (extraimports, extradeps []string, err error) {
	testFuncs, err := loadTestFuncs(testPkg, xtestPkg)
	if err != nil {
		return nil, nil, err
	}
	extraimports = []string{"testing", "testing/internal/testdeps"}
	if testFuncs.TestMain == nil {
		extraimports = append(extraimports, "os")
	}
	
	
	
	extradeps = []string{
		"errors",
		"internal/cpu",
		"unsafe",
		"internal/bytealg",
		"internal/race",
		"runtime/internal/atomic",
		"runtime/internal/sys",
		"runtime",
		"sync/atomic",
		"sync",
		"io",
		"unicode",
		"unicode/utf8",
		"bytes",
		"math",
		"syscall",
		"time",
		"internal/poll",
		"internal/syscall/unix",
		"internal/testlog",
		"os",
		"math/bits",
		"strconv",
		"reflect",
		"fmt",
		"sort",
		"strings",
		"flag",
		"runtime/debug",
		"context",
		"runtime/trace",
		"testing",
		"bufio",
		"regexp/syntax",
		"regexp",
		"compress/flate",
		"encoding/binary",
		"hash",
		"hash/crc32",
		"compress/gzip",
		"path/filepath",
		"io/ioutil",
		"text/tabwriter",
		"runtime/pprof",
		"testing/internal/testdeps",
	}
	return extraimports, extradeps, writeTestmain(out, testFuncs)
}





func isTestFunc(fn *ast.FuncDecl, arg string) bool {
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 ||
		fn.Type.Params.List == nil ||
		len(fn.Type.Params.List) != 1 ||
		len(fn.Type.Params.List[0].Names) > 1 {
		return false
	}
	ptr, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	
	
	
	
	if name, ok := ptr.X.(*ast.Ident); ok && name.Name == arg {
		return true
	}
	if sel, ok := ptr.X.(*ast.SelectorExpr); ok && sel.Sel.Name == arg {
		return true
	}
	return false
}




func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { 
		return true
	}
	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}


func loadTestFuncs(ptest, pxtest *Package) (*testFuncs, error) {
	t := &testFuncs{
		TestPackage:  ptest,
		XTestPackage: pxtest,
	}
	for _, file := range ptest.GoFiles {
		if !strings.HasSuffix(file, "_test.go") {
			continue
		}
		if err := t.load(file, "_test", &t.ImportTest, &t.NeedTest); err != nil {
			return nil, err
		}
	}
	if pxtest != nil {
		for _, file := range pxtest.GoFiles {
			if err := t.load(file, "_xtest", &t.ImportXtest, &t.NeedXtest); err != nil {
				return nil, err
			}
		}
	}
	return t, nil
}


func writeTestmain(out string, t *testFuncs) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := testmainTmpl.Execute(f, t); err != nil {
		return err
	}

	return nil
}

type testFuncs struct {
	Tests        []testFunc
	Benchmarks   []testFunc
	Examples     []testFunc
	TestMain     *testFunc
	TestPackage  *Package
	XTestPackage *Package
	ImportTest   bool
	NeedTest     bool
	ImportXtest  bool
	NeedXtest    bool
}


func (t *testFuncs) Tested() string {
	return t.TestPackage.Name
}

type testFunc struct {
	Package   string 
	Name      string 
	Output    string 
	Unordered bool   
}

func (t *testFuncs) load(filename, pkg string, doImport, seen *bool) error {
	var fset = token.NewFileSet()

	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return errors.New("failed to parse test file " + filename)
	}
	for _, d := range f.Decls {
		n, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if n.Recv != nil {
			continue
		}
		name := n.Name.String()
		switch {
		case name == "TestMain":
			if isTestFunc(n, "T") {
				t.Tests = append(t.Tests, testFunc{pkg, name, "", false})
				*doImport, *seen = true, true
				continue
			}
			err := checkTestFunc(fset, n, "M")
			if err != nil {
				return err
			}
			if t.TestMain != nil {
				return errors.New("multiple definitions of TestMain")
			}
			t.TestMain = &testFunc{pkg, name, "", false}
			*doImport, *seen = true, true
		case isTest(name, "Test"):
			err := checkTestFunc(fset, n, "T")
			if err != nil {
				return err
			}
			t.Tests = append(t.Tests, testFunc{pkg, name, "", false})
			*doImport, *seen = true, true
		case isTest(name, "Benchmark"):
			err := checkTestFunc(fset, n, "B")
			if err != nil {
				return err
			}
			t.Benchmarks = append(t.Benchmarks, testFunc{pkg, name, "", false})
			*doImport, *seen = true, true
		}
	}
	ex := doc.Examples(f)
	sort.Slice(ex, func(i, j int) bool { return ex[i].Order < ex[j].Order })
	for _, e := range ex {
		*doImport = true 
		if e.Output == "" && !e.EmptyOutput {
			
			continue
		}
		t.Examples = append(t.Examples, testFunc{pkg, "Example" + e.Name, e.Output, e.Unordered})
		*seen = true
	}
	return nil
}

func checkTestFunc(fset *token.FileSet, fn *ast.FuncDecl, arg string) error {
	if !isTestFunc(fn, arg) {
		name := fn.Name.String()
		pos := fset.Position(fn.Pos())
		return fmt.Errorf("%s: wrong signature for %s, must be: func %s(%s *testing.%s)", pos, name, name, strings.ToLower(arg), arg)
	}
	return nil
}

var testmainTmpl = template.Must(template.New("main").Parse(`
package main

import (
{{if not .TestMain}}
	"os"
{{end}}
	"testing"
	"testing/internal/testdeps"

{{if .ImportTest}}
	{{if .NeedTest}}_test{{else}}_{{end}} {{.TestPackage.PkgPath | printf "%q"}}
{{end}}
{{if .ImportXtest}}
	{{if .NeedXtest}}_xtest{{else}}_{{end}} {{.XTestPackage.PkgPath | printf "%q"}}
{{end}}
)

var tests = []testing.InternalTest{
{{range .Tests}}
	{"{{.Name}}", {{.Package}}.{{.Name}}},
{{end}}
}

var benchmarks = []testing.InternalBenchmark{
{{range .Benchmarks}}
	{"{{.Name}}", {{.Package}}.{{.Name}}},
{{end}}
}

var examples = []testing.InternalExample{
{{range .Examples}}
	{"{{.Name}}", {{.Package}}.{{.Name}}, {{.Output | printf "%q"}}, {{.Unordered}}},
{{end}}
}

func init() {
	testdeps.ImportPath = {{.TestPackage.PkgPath | printf "%q"}}
}

func main() {
	m := testing.MainStart(testdeps.TestDeps{}, tests, benchmarks, examples)
{{with .TestMain}}
	{{.Package}}.{{.Name}}(m)
{{else}}
	os.Exit(m.Run())
{{end}}
}

`))
