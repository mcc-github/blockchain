



package astutil



import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
)














































func PathEnclosingInterval(root *ast.File, start, end token.Pos) (path []ast.Node, exact bool) {
	

	
	var visit func(node ast.Node) bool
	visit = func(node ast.Node) bool {
		path = append(path, node)

		nodePos := node.Pos()
		nodeEnd := node.End()

		

		
		if start < nodePos {
			start = nodePos
		}
		if end > nodeEnd {
			end = nodeEnd
		}

		
		children := childrenOf(node)
		l := len(children)
		for i, child := range children {
			
			childPos := child.Pos()
			childEnd := child.End()

			
			augPos := childPos
			augEnd := childEnd
			if i > 0 {
				augPos = children[i-1].End() 
			}
			if i < l-1 {
				nextChildPos := children[i+1].Pos()
				
				if start >= augEnd && end <= nextChildPos {
					return false 
				}
				augEnd = nextChildPos 
			}

			
			

			
			if augPos <= start && end <= augEnd {
				_, isToken := child.(tokenNode)
				return isToken || visit(child)
			}

			
			
			
			if start < childEnd && end > augEnd {
				break
			}
		}

		
		

		
		
		
		
		if start == nodePos && end == nodeEnd {
			return true 
		}

		return false 
	}

	if start > end {
		start, end = end, start
	}

	if start < root.End() && end > root.Pos() {
		if start == end {
			end = start + 1 
		}
		exact = visit(root)

		
		for i, l := 0, len(path); i < l/2; i++ {
			path[i], path[l-1-i] = path[l-1-i], path[i]
		}
	} else {
		
		
		
		path = append(path, root)
	}

	return
}





type tokenNode struct {
	pos token.Pos
	end token.Pos
}

func (n tokenNode) Pos() token.Pos {
	return n.pos
}

func (n tokenNode) End() token.Pos {
	return n.end
}

func tok(pos token.Pos, len int) ast.Node {
	return tokenNode{pos, pos + token.Pos(len)}
}





func childrenOf(n ast.Node) []ast.Node {
	var children []ast.Node

	
	ast.Inspect(n, func(node ast.Node) bool {
		if node == n { 
			return true 
		}
		if node != nil { 
			children = append(children, node)
		}
		return false 
	})

	
	switch n := n.(type) {
	case *ast.ArrayType:
		children = append(children,
			tok(n.Lbrack, len("[")),
			tok(n.Elt.End(), len("]")))

	case *ast.AssignStmt:
		children = append(children,
			tok(n.TokPos, len(n.Tok.String())))

	case *ast.BasicLit:
		children = append(children,
			tok(n.ValuePos, len(n.Value)))

	case *ast.BinaryExpr:
		children = append(children, tok(n.OpPos, len(n.Op.String())))

	case *ast.BlockStmt:
		children = append(children,
			tok(n.Lbrace, len("{")),
			tok(n.Rbrace, len("}")))

	case *ast.BranchStmt:
		children = append(children,
			tok(n.TokPos, len(n.Tok.String())))

	case *ast.CallExpr:
		children = append(children,
			tok(n.Lparen, len("(")),
			tok(n.Rparen, len(")")))
		if n.Ellipsis != 0 {
			children = append(children, tok(n.Ellipsis, len("...")))
		}

	case *ast.CaseClause:
		if n.List == nil {
			children = append(children,
				tok(n.Case, len("default")))
		} else {
			children = append(children,
				tok(n.Case, len("case")))
		}
		children = append(children, tok(n.Colon, len(":")))

	case *ast.ChanType:
		switch n.Dir {
		case ast.RECV:
			children = append(children, tok(n.Begin, len("<-chan")))
		case ast.SEND:
			children = append(children, tok(n.Begin, len("chan<-")))
		case ast.RECV | ast.SEND:
			children = append(children, tok(n.Begin, len("chan")))
		}

	case *ast.CommClause:
		if n.Comm == nil {
			children = append(children,
				tok(n.Case, len("default")))
		} else {
			children = append(children,
				tok(n.Case, len("case")))
		}
		children = append(children, tok(n.Colon, len(":")))

	case *ast.Comment:
		

	case *ast.CommentGroup:
		

	case *ast.CompositeLit:
		children = append(children,
			tok(n.Lbrace, len("{")),
			tok(n.Rbrace, len("{")))

	case *ast.DeclStmt:
		

	case *ast.DeferStmt:
		children = append(children,
			tok(n.Defer, len("defer")))

	case *ast.Ellipsis:
		children = append(children,
			tok(n.Ellipsis, len("...")))

	case *ast.EmptyStmt:
		

	case *ast.ExprStmt:
		

	case *ast.Field:
		

	case *ast.FieldList:
		children = append(children,
			tok(n.Opening, len("(")),
			tok(n.Closing, len(")")))

	case *ast.File:
		
		children = append(children,
			tok(n.Package, len("package")))

	case *ast.ForStmt:
		children = append(children,
			tok(n.For, len("for")))

	case *ast.FuncDecl:
		

		
		
		
		
		
		
		
		children = nil 
		children = append(children, tok(n.Type.Func, len("func")))
		if n.Recv != nil {
			children = append(children, n.Recv)
		}
		children = append(children, n.Name)
		if n.Type.Params != nil {
			children = append(children, n.Type.Params)
		}
		if n.Type.Results != nil {
			children = append(children, n.Type.Results)
		}
		if n.Body != nil {
			children = append(children, n.Body)
		}

	case *ast.FuncLit:
		

	case *ast.FuncType:
		if n.Func != 0 {
			children = append(children,
				tok(n.Func, len("func")))
		}

	case *ast.GenDecl:
		children = append(children,
			tok(n.TokPos, len(n.Tok.String())))
		if n.Lparen != 0 {
			children = append(children,
				tok(n.Lparen, len("(")),
				tok(n.Rparen, len(")")))
		}

	case *ast.GoStmt:
		children = append(children,
			tok(n.Go, len("go")))

	case *ast.Ident:
		children = append(children,
			tok(n.NamePos, len(n.Name)))

	case *ast.IfStmt:
		children = append(children,
			tok(n.If, len("if")))

	case *ast.ImportSpec:
		

	case *ast.IncDecStmt:
		children = append(children,
			tok(n.TokPos, len(n.Tok.String())))

	case *ast.IndexExpr:
		children = append(children,
			tok(n.Lbrack, len("{")),
			tok(n.Rbrack, len("}")))

	case *ast.InterfaceType:
		children = append(children,
			tok(n.Interface, len("interface")))

	case *ast.KeyValueExpr:
		children = append(children,
			tok(n.Colon, len(":")))

	case *ast.LabeledStmt:
		children = append(children,
			tok(n.Colon, len(":")))

	case *ast.MapType:
		children = append(children,
			tok(n.Map, len("map")))

	case *ast.ParenExpr:
		children = append(children,
			tok(n.Lparen, len("(")),
			tok(n.Rparen, len(")")))

	case *ast.RangeStmt:
		children = append(children,
			tok(n.For, len("for")),
			tok(n.TokPos, len(n.Tok.String())))

	case *ast.ReturnStmt:
		children = append(children,
			tok(n.Return, len("return")))

	case *ast.SelectStmt:
		children = append(children,
			tok(n.Select, len("select")))

	case *ast.SelectorExpr:
		

	case *ast.SendStmt:
		children = append(children,
			tok(n.Arrow, len("<-")))

	case *ast.SliceExpr:
		children = append(children,
			tok(n.Lbrack, len("[")),
			tok(n.Rbrack, len("]")))

	case *ast.StarExpr:
		children = append(children, tok(n.Star, len("*")))

	case *ast.StructType:
		children = append(children, tok(n.Struct, len("struct")))

	case *ast.SwitchStmt:
		children = append(children, tok(n.Switch, len("switch")))

	case *ast.TypeAssertExpr:
		children = append(children,
			tok(n.Lparen-1, len(".")),
			tok(n.Lparen, len("(")),
			tok(n.Rparen, len(")")))

	case *ast.TypeSpec:
		

	case *ast.TypeSwitchStmt:
		children = append(children, tok(n.Switch, len("switch")))

	case *ast.UnaryExpr:
		children = append(children, tok(n.OpPos, len(n.Op.String())))

	case *ast.ValueSpec:
		

	case *ast.BadDecl, *ast.BadExpr, *ast.BadStmt:
		
	}

	
	
	
	
	sort.Sort(byPos(children))

	return children
}

type byPos []ast.Node

func (sl byPos) Len() int {
	return len(sl)
}
func (sl byPos) Less(i, j int) bool {
	return sl[i].Pos() < sl[j].Pos()
}
func (sl byPos) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}








func NodeDescription(n ast.Node) string {
	switch n := n.(type) {
	case *ast.ArrayType:
		return "array type"
	case *ast.AssignStmt:
		return "assignment"
	case *ast.BadDecl:
		return "bad declaration"
	case *ast.BadExpr:
		return "bad expression"
	case *ast.BadStmt:
		return "bad statement"
	case *ast.BasicLit:
		return "basic literal"
	case *ast.BinaryExpr:
		return fmt.Sprintf("binary %s operation", n.Op)
	case *ast.BlockStmt:
		return "block"
	case *ast.BranchStmt:
		switch n.Tok {
		case token.BREAK:
			return "break statement"
		case token.CONTINUE:
			return "continue statement"
		case token.GOTO:
			return "goto statement"
		case token.FALLTHROUGH:
			return "fall-through statement"
		}
	case *ast.CallExpr:
		if len(n.Args) == 1 && !n.Ellipsis.IsValid() {
			return "function call (or conversion)"
		}
		return "function call"
	case *ast.CaseClause:
		return "case clause"
	case *ast.ChanType:
		return "channel type"
	case *ast.CommClause:
		return "communication clause"
	case *ast.Comment:
		return "comment"
	case *ast.CommentGroup:
		return "comment group"
	case *ast.CompositeLit:
		return "composite literal"
	case *ast.DeclStmt:
		return NodeDescription(n.Decl) + " statement"
	case *ast.DeferStmt:
		return "defer statement"
	case *ast.Ellipsis:
		return "ellipsis"
	case *ast.EmptyStmt:
		return "empty statement"
	case *ast.ExprStmt:
		return "expression statement"
	case *ast.Field:
		
		
		
		
		
		
		return "field/method/parameter"
	case *ast.FieldList:
		return "field/method/parameter list"
	case *ast.File:
		return "source file"
	case *ast.ForStmt:
		return "for loop"
	case *ast.FuncDecl:
		return "function declaration"
	case *ast.FuncLit:
		return "function literal"
	case *ast.FuncType:
		return "function type"
	case *ast.GenDecl:
		switch n.Tok {
		case token.IMPORT:
			return "import declaration"
		case token.CONST:
			return "constant declaration"
		case token.TYPE:
			return "type declaration"
		case token.VAR:
			return "variable declaration"
		}
	case *ast.GoStmt:
		return "go statement"
	case *ast.Ident:
		return "identifier"
	case *ast.IfStmt:
		return "if statement"
	case *ast.ImportSpec:
		return "import specification"
	case *ast.IncDecStmt:
		if n.Tok == token.INC {
			return "increment statement"
		}
		return "decrement statement"
	case *ast.IndexExpr:
		return "index expression"
	case *ast.InterfaceType:
		return "interface type"
	case *ast.KeyValueExpr:
		return "key/value association"
	case *ast.LabeledStmt:
		return "statement label"
	case *ast.MapType:
		return "map type"
	case *ast.Package:
		return "package"
	case *ast.ParenExpr:
		return "parenthesized " + NodeDescription(n.X)
	case *ast.RangeStmt:
		return "range loop"
	case *ast.ReturnStmt:
		return "return statement"
	case *ast.SelectStmt:
		return "select statement"
	case *ast.SelectorExpr:
		return "selector"
	case *ast.SendStmt:
		return "channel send"
	case *ast.SliceExpr:
		return "slice expression"
	case *ast.StarExpr:
		return "*-operation" 
	case *ast.StructType:
		return "struct type"
	case *ast.SwitchStmt:
		return "switch statement"
	case *ast.TypeAssertExpr:
		return "type assertion"
	case *ast.TypeSpec:
		return "type specification"
	case *ast.TypeSwitchStmt:
		return "type switch"
	case *ast.UnaryExpr:
		return fmt.Sprintf("unary %s operation", n.Op)
	case *ast.ValueSpec:
		return "value specification"

	}
	panic(fmt.Sprintf("unexpected node type: %T", n))
}
