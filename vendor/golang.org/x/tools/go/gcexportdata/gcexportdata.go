



















package gcexportdata 

import (
	"bufio"
	"bytes"
	"fmt"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"

	"golang.org/x/tools/go/internal/gcimporter"
)










func Find(importPath, srcDir string) (filename, path string) {
	return gcimporter.FindPkg(importPath, srcDir)
}




func NewReader(r io.Reader) (io.Reader, error) {
	buf := bufio.NewReader(r)
	_, err := gcimporter.FindExportData(buf)
	
	
	
	return buf, err
}













func Read(in io.Reader, fset *token.FileSet, imports map[string]*types.Package, path string) (*types.Package, error) {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("reading export data for %q: %v", path, err)
	}

	if bytes.HasPrefix(data, []byte("!<arch>")) {
		return nil, fmt.Errorf("can't read export data for %q directly from an archive file (call gcexportdata.NewReader first to extract export data)", path)
	}

	
	
	if bytes.HasPrefix(data, []byte("package ")) {
		return gcimporter.ImportData(imports, path, path, bytes.NewReader(data))
	}

	
	
	
	if len(data) > 0 && data[0] == 'i' {
		_, pkg, err := gcimporter.IImportData(fset, imports, data[1:], path)
		return pkg, err
	}

	_, pkg, err := gcimporter.BImportData(fset, imports, data, path)
	return pkg, err
}



func Write(out io.Writer, fset *token.FileSet, pkg *types.Package) error {
	b, err := gcimporter.BExportData(fset, pkg)
	if err != nil {
		return err
	}
	_, err = out.Write(b)
	return err
}
