package convert

import (
	"go/ast"
)


func rewriteOtherFuncsToUseGinkgoT(declarations []ast.Decl) {
	for _, decl := range declarations {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		for _, param := range decl.Type.Params.List {
			starExpr, ok := param.Type.(*ast.StarExpr)
			if !ok {
				continue
			}

			selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			xIdent, ok := selectorExpr.X.(*ast.Ident)
			if !ok || xIdent.Name != "testing" {
				continue
			}

			if selectorExpr.Sel.Name != "T" {
				continue
			}

			param.Type = newGinkgoTInterface()
		}
	}
}


func walkNodesInRootNodeReplacingTestingT(rootNode *ast.File) {
	ast.Inspect(rootNode, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		switch node := node.(type) {
		case *ast.StructType:
			replaceTestingTsInStructType(node)
		case *ast.FuncLit:
			replaceTypeDeclTestingTsInFuncLiteral(node)
		}

		return true
	})
}


func replaceNamedTestingTsInKeyValueExpression(kve *ast.KeyValueExpr, testingT string) {
	ident, ok := kve.Value.(*ast.Ident)
	if !ok {
		return
	}

	if ident.Name == testingT {
		kve.Value = newGinkgoTFromIdent(ident)
	}
}


func replaceTypeDeclTestingTsInFuncLiteral(functionLiteral *ast.FuncLit) {
	for _, arg := range functionLiteral.Type.Params.List {
		starExpr, ok := arg.Type.(*ast.StarExpr)
		if !ok {
			continue
		}

		selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		target, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			continue
		}

		if target.Name == "testing" && selectorExpr.Sel.Name == "T" {
			arg.Type = newGinkgoTInterface()
		}
	}
}


func replaceTestingTsInStructType(structType *ast.StructType) {
	for _, field := range structType.Fields.List {
		starExpr, ok := field.Type.(*ast.StarExpr)
		if !ok {
			continue
		}

		selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		xIdent, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			continue
		}

		if xIdent.Name == "testing" && selectorExpr.Sel.Name == "T" {
			field.Type = newGinkgoTInterface()
		}
	}
}
