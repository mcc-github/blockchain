




package astutil 

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)


func AddImport(fset *token.FileSet, f *ast.File, ipath string) (added bool) {
	return AddNamedImport(fset, f, "", ipath)
}








func AddNamedImport(fset *token.FileSet, f *ast.File, name, ipath string) (added bool) {
	if imports(f, ipath) {
		return false
	}

	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(ipath),
		},
	}
	if name != "" {
		newImport.Name = &ast.Ident{Name: name}
	}

	
	
	
	
	var (
		bestMatch  = -1         
		lastImport = -1         
		impDecl    *ast.GenDecl 
		impIndex   = -1         

		isThirdPartyPath = isThirdParty(ipath)
	)
	for i, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT {
			lastImport = i
			
			
			if declImports(gen, "C") {
				continue
			}

			
			if len(gen.Specs) == 0 && bestMatch == -1 {
				impDecl = gen
			}

			
			
			
			
			
			
			
			
			
			
			seenAnyThirdParty := false
			for j, spec := range gen.Specs {
				impspec := spec.(*ast.ImportSpec)
				p := importPath(impspec)
				n := matchLen(p, ipath)
				if n > bestMatch || (bestMatch == 0 && !seenAnyThirdParty && isThirdPartyPath) {
					bestMatch = n
					impDecl = gen
					impIndex = j
				}
				seenAnyThirdParty = seenAnyThirdParty || isThirdParty(p)
			}
		}
	}

	
	if impDecl == nil {
		impDecl = &ast.GenDecl{
			Tok: token.IMPORT,
		}
		if lastImport >= 0 {
			impDecl.TokPos = f.Decls[lastImport].End()
		} else {
			
			
			
			
			impDecl.TokPos = f.Package

			file := fset.File(f.Package)
			pkgLine := file.Line(f.Package)
			for _, c := range f.Comments {
				if file.Line(c.Pos()) > pkgLine {
					break
				}
				
				impDecl.TokPos = c.End() + 2
			}
		}
		f.Decls = append(f.Decls, nil)
		copy(f.Decls[lastImport+2:], f.Decls[lastImport+1:])
		f.Decls[lastImport+1] = impDecl
	}

	
	insertAt := 0
	if impIndex >= 0 {
		
		insertAt = impIndex + 1
	}
	impDecl.Specs = append(impDecl.Specs, nil)
	copy(impDecl.Specs[insertAt+1:], impDecl.Specs[insertAt:])
	impDecl.Specs[insertAt] = newImport
	pos := impDecl.Pos()
	if insertAt > 0 {
		
		
		if spec, ok := impDecl.Specs[insertAt-1].(*ast.ImportSpec); ok && spec.Comment != nil {
			pos = spec.Comment.End()
		} else {
			
			
			pos = impDecl.Specs[insertAt-1].Pos()
		}
	}
	if newImport.Name != nil {
		newImport.Name.NamePos = pos
	}
	newImport.Path.ValuePos = pos
	newImport.EndPos = pos

	
	if len(impDecl.Specs) == 1 {
		
		impDecl.Lparen = token.NoPos
	} else if !impDecl.Lparen.IsValid() {
		
		impDecl.Lparen = impDecl.Specs[0].Pos()
	}

	f.Imports = append(f.Imports, newImport)

	if len(f.Decls) <= 1 {
		return true
	}

	
	var first *ast.GenDecl
	for i := 0; i < len(f.Decls); i++ {
		decl := f.Decls[i]
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT || declImports(gen, "C") {
			continue
		}
		if first == nil {
			first = gen
			continue 
		}
		
		
		first.Lparen = first.Pos()
		
		for _, spec := range gen.Specs {
			spec.(*ast.ImportSpec).Path.ValuePos = first.Pos()
			first.Specs = append(first.Specs, spec)
		}
		f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
		i--
	}

	return true
}

func isThirdParty(importPath string) bool {
	
	
	return strings.Contains(importPath, ".")
}


func DeleteImport(fset *token.FileSet, f *ast.File, path string) (deleted bool) {
	return DeleteNamedImport(fset, f, "", path)
}


func DeleteNamedImport(fset *token.FileSet, f *ast.File, name, path string) (deleted bool) {
	var delspecs []*ast.ImportSpec
	var delcomments []*ast.CommentGroup

	
	for i := 0; i < len(f.Decls); i++ {
		decl := f.Decls[i]
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}
		for j := 0; j < len(gen.Specs); j++ {
			spec := gen.Specs[j]
			impspec := spec.(*ast.ImportSpec)
			if impspec.Name == nil && name != "" {
				continue
			}
			if impspec.Name != nil && impspec.Name.Name != name {
				continue
			}
			if importPath(impspec) != path {
				continue
			}

			
			
			delspecs = append(delspecs, impspec)
			deleted = true
			copy(gen.Specs[j:], gen.Specs[j+1:])
			gen.Specs = gen.Specs[:len(gen.Specs)-1]

			
			
			if len(gen.Specs) == 0 {
				copy(f.Decls[i:], f.Decls[i+1:])
				f.Decls = f.Decls[:len(f.Decls)-1]
				i--
				break
			} else if len(gen.Specs) == 1 {
				if impspec.Doc != nil {
					delcomments = append(delcomments, impspec.Doc)
				}
				if impspec.Comment != nil {
					delcomments = append(delcomments, impspec.Comment)
				}
				for _, cg := range f.Comments {
					
					if cg.End() < impspec.Pos() && fset.Position(cg.End()).Line == fset.Position(impspec.Pos()).Line {
						delcomments = append(delcomments, cg)
						break
					}
				}

				spec := gen.Specs[0].(*ast.ImportSpec)

				
				if spec.Doc != nil {
					for fset.Position(gen.TokPos).Line+1 < fset.Position(spec.Doc.Pos()).Line {
						fset.File(gen.TokPos).MergeLine(fset.Position(gen.TokPos).Line)
					}
				}
				for _, cg := range f.Comments {
					if cg.End() < spec.Pos() && fset.Position(cg.End()).Line == fset.Position(spec.Pos()).Line {
						for fset.Position(gen.TokPos).Line+1 < fset.Position(spec.Pos()).Line {
							fset.File(gen.TokPos).MergeLine(fset.Position(gen.TokPos).Line)
						}
						break
					}
				}
			}
			if j > 0 {
				lastImpspec := gen.Specs[j-1].(*ast.ImportSpec)
				lastLine := fset.Position(lastImpspec.Path.ValuePos).Line
				line := fset.Position(impspec.Path.ValuePos).Line

				
				
				if line-lastLine > 1 {
					
					
					
				} else if line != fset.File(gen.Rparen).LineCount() {
					
					fset.File(gen.Rparen).MergeLine(line)
				}
			}
			j--
		}
	}

	
	for i := 0; i < len(f.Imports); i++ {
		imp := f.Imports[i]
		for j, del := range delspecs {
			if imp == del {
				copy(f.Imports[i:], f.Imports[i+1:])
				f.Imports = f.Imports[:len(f.Imports)-1]
				copy(delspecs[j:], delspecs[j+1:])
				delspecs = delspecs[:len(delspecs)-1]
				i--
				break
			}
		}
	}

	
	for i := 0; i < len(f.Comments); i++ {
		cg := f.Comments[i]
		for j, del := range delcomments {
			if cg == del {
				copy(f.Comments[i:], f.Comments[i+1:])
				f.Comments = f.Comments[:len(f.Comments)-1]
				copy(delcomments[j:], delcomments[j+1:])
				delcomments = delcomments[:len(delcomments)-1]
				i--
				break
			}
		}
	}

	if len(delspecs) > 0 {
		panic(fmt.Sprintf("deleted specs from Decls but not Imports: %v", delspecs))
	}

	return
}


func RewriteImport(fset *token.FileSet, f *ast.File, oldPath, newPath string) (rewrote bool) {
	for _, imp := range f.Imports {
		if importPath(imp) == oldPath {
			rewrote = true
			
			
			imp.EndPos = imp.End()
			imp.Path.Value = strconv.Quote(newPath)
		}
	}
	return
}


func UsesImport(f *ast.File, path string) (used bool) {
	spec := importSpec(f, path)
	if spec == nil {
		return
	}

	name := spec.Name.String()
	switch name {
	case "<nil>":
		
		
		lastSlash := strings.LastIndex(path, "/")
		if lastSlash == -1 {
			name = path
		} else {
			name = path[lastSlash+1:]
		}
	case "_", ".":
		
		return true
	}

	ast.Walk(visitFn(func(n ast.Node) {
		sel, ok := n.(*ast.SelectorExpr)
		if ok && isTopName(sel.X, name) {
			used = true
		}
	}), f)

	return
}

type visitFn func(node ast.Node)

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	fn(node)
	return fn
}


func imports(f *ast.File, path string) bool {
	return importSpec(f, path) != nil
}



func importSpec(f *ast.File, path string) *ast.ImportSpec {
	for _, s := range f.Imports {
		if importPath(s) == path {
			return s
		}
	}
	return nil
}



func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err == nil {
		return t
	}
	return ""
}


func declImports(gen *ast.GenDecl, path string) bool {
	if gen.Tok != token.IMPORT {
		return false
	}
	for _, spec := range gen.Specs {
		impspec := spec.(*ast.ImportSpec)
		if importPath(impspec) == path {
			return true
		}
	}
	return false
}


func matchLen(x, y string) int {
	n := 0
	for i := 0; i < len(x) && i < len(y) && x[i] == y[i]; i++ {
		if x[i] == '/' {
			n++
		}
	}
	return n
}


func isTopName(n ast.Expr, name string) bool {
	id, ok := n.(*ast.Ident)
	return ok && id.Name == name && id.Obj == nil
}


func Imports(fset *token.FileSet, f *ast.File) [][]*ast.ImportSpec {
	var groups [][]*ast.ImportSpec

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			break
		}

		group := []*ast.ImportSpec{}

		var lastLine int
		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)
			pos := importSpec.Path.ValuePos
			line := fset.Position(pos).Line
			if lastLine > 0 && pos > 0 && line-lastLine > 1 {
				groups = append(groups, group)
				group = []*ast.ImportSpec{}
			}
			group = append(group, importSpec)
			lastLine = line
		}
		groups = append(groups, group)
	}

	return groups
}
