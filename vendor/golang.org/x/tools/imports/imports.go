







package imports 

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)


type Options struct {
	Fragment  bool 
	AllErrors bool 

	Comments  bool 
	TabIndent bool 
	TabWidth  int  

	FormatOnly bool 
}







func Process(filename string, src []byte, opt *Options) ([]byte, error) {
	if opt == nil {
		opt = &Options{Comments: true, TabIndent: true, TabWidth: 8}
	}
	if src == nil {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		src = b
	}

	fileSet := token.NewFileSet()
	file, adjust, err := parse(fileSet, filename, src, opt)
	if err != nil {
		return nil, err
	}

	if !opt.FormatOnly {
		_, err = fixImports(fileSet, file, filename)
		if err != nil {
			return nil, err
		}
	}

	sortImports(fileSet, file)
	imps := astutil.Imports(fileSet, file)
	var spacesBefore []string 
	for _, impSection := range imps {
		
		
		
		
		lastGroup := -1
		for _, importSpec := range impSection {
			importPath, _ := strconv.Unquote(importSpec.Path.Value)
			groupNum := importGroup(importPath)
			if groupNum != lastGroup && lastGroup != -1 {
				spacesBefore = append(spacesBefore, importPath)
			}
			lastGroup = groupNum
		}

	}

	printerMode := printer.UseSpaces
	if opt.TabIndent {
		printerMode |= printer.TabIndent
	}
	printConfig := &printer.Config{Mode: printerMode, Tabwidth: opt.TabWidth}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, fileSet, file)
	if err != nil {
		return nil, err
	}
	out := buf.Bytes()
	if adjust != nil {
		out = adjust(src, out)
	}
	if len(spacesBefore) > 0 {
		out, err = addImportSpaces(bytes.NewReader(out), spacesBefore)
		if err != nil {
			return nil, err
		}
	}

	out, err = format.Source(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}



func parse(fset *token.FileSet, filename string, src []byte, opt *Options) (*ast.File, func(orig, src []byte) []byte, error) {
	parserMode := parser.Mode(0)
	if opt.Comments {
		parserMode |= parser.ParseComments
	}
	if opt.AllErrors {
		parserMode |= parser.AllErrors
	}

	
	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err == nil {
		return file, nil, nil
	}
	
	
	
	if !opt.Fragment || !strings.Contains(err.Error(), "expected 'package'") {
		return nil, nil, err
	}

	
	
	
	
	const prefix = "package main;"
	psrc := append([]byte(prefix), src...)
	file, err = parser.ParseFile(fset, filename, psrc, parserMode)
	if err == nil {
		
		
		
		psrc[len(prefix)-1] = '\n'
		fset.File(file.Package).SetLinesForContent(psrc)

		
		
		if containsMainFunc(file) {
			return file, nil, nil
		}

		adjust := func(orig, src []byte) []byte {
			
			src = src[len(prefix):]
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	
	
	
	if !strings.Contains(err.Error(), "expected declaration") {
		return nil, nil, err
	}

	
	
	
	
	
	fsrc := append(append([]byte("package p; func _() {"), src...), '}')
	file, err = parser.ParseFile(fset, filename, fsrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			
			
			src = src[len("package p\n\nfunc _() {"):]
			src = src[:len(src)-len("}\n")]
			
			
			src = bytes.Replace(src, []byte("\n\t"), []byte("\n"), -1)
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}

	
	return nil, nil, err
}



func containsMainFunc(file *ast.File) bool {
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Name.Name != "main" {
				continue
			}

			if len(f.Type.Params.List) != 0 {
				continue
			}

			if f.Type.Results != nil && len(f.Type.Results.List) != 0 {
				continue
			}

			return true
		}
	}

	return false
}

func cutSpace(b []byte) (before, middle, after []byte) {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n') {
		i++
	}
	j := len(b)
	for j > 0 && (b[j-1] == ' ' || b[j-1] == '\t' || b[j-1] == '\n') {
		j--
	}
	if i <= j {
		return b[:i], b[i:j], b[j:]
	}
	return nil, nil, b[j:]
}







func matchSpace(orig []byte, src []byte) []byte {
	before, _, after := cutSpace(orig)
	i := bytes.LastIndex(before, []byte{'\n'})
	before, indent := before[:i+1], before[i+1:]

	_, src, _ = cutSpace(src)

	var b bytes.Buffer
	b.Write(before)
	for len(src) > 0 {
		line := src
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, src = line[:i+1], line[i+1:]
		} else {
			src = nil
		}
		if len(line) > 0 && line[0] != '\n' { 
			b.Write(indent)
		}
		b.Write(line)
	}
	b.Write(after)
	return b.Bytes()
}

var impLine = regexp.MustCompile(`^\s+(?:[\w\.]+\s+)?"(.+)"`)

func addImportSpaces(r io.Reader, breaks []string) ([]byte, error) {
	var out bytes.Buffer
	in := bufio.NewReader(r)
	inImports := false
	done := false
	for {
		s, err := in.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !inImports && !done && strings.HasPrefix(s, "import") {
			inImports = true
		}
		if inImports && (strings.HasPrefix(s, "var") ||
			strings.HasPrefix(s, "func") ||
			strings.HasPrefix(s, "const") ||
			strings.HasPrefix(s, "type")) {
			done = true
			inImports = false
		}
		if inImports && len(breaks) > 0 {
			if m := impLine.FindStringSubmatch(s); m != nil {
				if m[1] == breaks[0] {
					out.WriteByte('\n')
					breaks = breaks[1:]
				}
			}
		}

		fmt.Fprint(&out, s)
	}
	return out.Bytes(), nil
}
