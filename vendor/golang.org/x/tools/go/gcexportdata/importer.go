



package gcexportdata

import (
	"fmt"
	"go/token"
	"go/types"
	"os"
)













func NewImporter(fset *token.FileSet, imports map[string]*types.Package) types.ImporterFrom {
	return importer{fset, imports}
}

type importer struct {
	fset    *token.FileSet
	imports map[string]*types.Package
}

func (imp importer) Import(importPath string) (*types.Package, error) {
	return imp.ImportFrom(importPath, "", 0)
}

func (imp importer) ImportFrom(importPath, srcDir string, mode types.ImportMode) (_ *types.Package, err error) {
	filename, path := Find(importPath, srcDir)
	if filename == "" {
		if importPath == "unsafe" {
			
			
			return types.Unsafe, nil
		}
		return nil, fmt.Errorf("can't find import: %s", importPath)
	}

	if pkg, ok := imp.imports[path]; ok && pkg.Complete() {
		return pkg, nil 
	}

	
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
		if err != nil {
			
			err = fmt.Errorf("reading export data: %s: %v", filename, err)
		}
	}()

	r, err := NewReader(f)
	if err != nil {
		return nil, err
	}

	return Read(r, imp.fset, imp.imports, path)
}
