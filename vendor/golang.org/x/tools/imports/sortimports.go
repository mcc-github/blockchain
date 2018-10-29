





package imports

import (
	"go/ast"
	"go/token"
	"sort"
	"strconv"
)



func sortImports(fset *token.FileSet, f *ast.File) {
	for i, d := range f.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			
			
			break
		}

		if len(d.Specs) == 0 {
			
			f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
		}

		if !d.Lparen.IsValid() {
			
			continue
		}

		
		i := 0
		specs := d.Specs[:0]
		for j, s := range d.Specs {
			if j > i && fset.Position(s.Pos()).Line > 1+fset.Position(d.Specs[j-1].End()).Line {
				
				specs = append(specs, sortSpecs(fset, f, d.Specs[i:j])...)
				i = j
			}
		}
		specs = append(specs, sortSpecs(fset, f, d.Specs[i:])...)
		d.Specs = specs

		
		if len(d.Specs) > 0 {
			lastSpec := d.Specs[len(d.Specs)-1]
			lastLine := fset.Position(lastSpec.Pos()).Line
			if rParenLine := fset.Position(d.Rparen).Line; rParenLine > lastLine+1 {
				fset.File(d.Rparen).MergeLine(rParenLine - 1)
			}
		}
	}
}

func importPath(s ast.Spec) string {
	t, err := strconv.Unquote(s.(*ast.ImportSpec).Path.Value)
	if err == nil {
		return t
	}
	return ""
}

func importName(s ast.Spec) string {
	n := s.(*ast.ImportSpec).Name
	if n == nil {
		return ""
	}
	return n.Name
}

func importComment(s ast.Spec) string {
	c := s.(*ast.ImportSpec).Comment
	if c == nil {
		return ""
	}
	return c.Text()
}


func collapse(prev, next ast.Spec) bool {
	if importPath(next) != importPath(prev) || importName(next) != importName(prev) {
		return false
	}
	return prev.(*ast.ImportSpec).Comment == nil
}

type posSpan struct {
	Start token.Pos
	End   token.Pos
}

func sortSpecs(fset *token.FileSet, f *ast.File, specs []ast.Spec) []ast.Spec {
	
	
	
	if len(specs) <= 1 {
		return specs
	}

	
	pos := make([]posSpan, len(specs))
	for i, s := range specs {
		pos[i] = posSpan{s.Pos(), s.End()}
	}

	
	
	lastLine := fset.Position(pos[len(pos)-1].End).Line
	cstart := len(f.Comments)
	cend := len(f.Comments)
	for i, g := range f.Comments {
		if g.Pos() < pos[0].Start {
			continue
		}
		if i < cstart {
			cstart = i
		}
		if fset.Position(g.End()).Line > lastLine {
			cend = i
			break
		}
	}
	comments := f.Comments[cstart:cend]

	
	importComment := map[*ast.ImportSpec][]*ast.CommentGroup{}
	specIndex := 0
	for _, g := range comments {
		for specIndex+1 < len(specs) && pos[specIndex+1].Start <= g.Pos() {
			specIndex++
		}
		s := specs[specIndex].(*ast.ImportSpec)
		importComment[s] = append(importComment[s], g)
	}

	
	
	
	
	
	sort.Sort(byImportSpec(specs))

	
	
	deduped := specs[:0]
	for i, s := range specs {
		if i == len(specs)-1 || !collapse(s, specs[i+1]) {
			deduped = append(deduped, s)
		} else {
			p := s.Pos()
			fset.File(p).MergeLine(fset.Position(p).Line)
		}
	}
	specs = deduped

	
	for i, s := range specs {
		s := s.(*ast.ImportSpec)
		if s.Name != nil {
			s.Name.NamePos = pos[i].Start
		}
		s.Path.ValuePos = pos[i].Start
		s.EndPos = pos[i].End
		nextSpecPos := pos[i].End

		for _, g := range importComment[s] {
			for _, c := range g.List {
				c.Slash = pos[i].End
				nextSpecPos = c.End()
			}
		}
		if i < len(specs)-1 {
			pos[i+1].Start = nextSpecPos
			pos[i+1].End = nextSpecPos
		}
	}

	sort.Sort(byCommentPos(comments))

	
	
	firstSpecLine := fset.Position(specs[0].Pos()).Line
	for _, s := range specs[1:] {
		p := s.Pos()
		line := fset.File(p).Line(p)
		for previousLine := line - 1; previousLine >= firstSpecLine; {
			fset.File(p).MergeLine(previousLine)
			previousLine--
		}
	}
	return specs
}

type byImportSpec []ast.Spec 

func (x byImportSpec) Len() int      { return len(x) }
func (x byImportSpec) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byImportSpec) Less(i, j int) bool {
	ipath := importPath(x[i])
	jpath := importPath(x[j])

	igroup := importGroup(ipath)
	jgroup := importGroup(jpath)
	if igroup != jgroup {
		return igroup < jgroup
	}

	if ipath != jpath {
		return ipath < jpath
	}
	iname := importName(x[i])
	jname := importName(x[j])

	if iname != jname {
		return iname < jname
	}
	return importComment(x[i]) < importComment(x[j])
}

type byCommentPos []*ast.CommentGroup

func (x byCommentPos) Len() int           { return len(x) }
func (x byCommentPos) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x byCommentPos) Less(i, j int) bool { return x[i].Pos() < x[j].Pos() }
