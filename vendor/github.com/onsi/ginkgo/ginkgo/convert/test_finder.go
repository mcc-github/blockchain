package convert

import (
	"go/ast"
	"regexp"
)


func findTestFuncs(rootNode *ast.File) (testsToRewrite []*ast.FuncDecl) {
	testNameRegexp := regexp.MustCompile("^Test[0-9A-Z].+")

	ast.Inspect(rootNode, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		switch node := node.(type) {
		case *ast.FuncDecl:
			matches := testNameRegexp.MatchString(node.Name.Name)

			if matches && receivesTestingT(node) {
				testsToRewrite = append(testsToRewrite, node)
			}
		}

		return true
	})

	return
}


func receivesTestingT(node *ast.FuncDecl) bool {
	if len(node.Type.Params.List) != 1 {
		return false
	}

	base, ok := node.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}

	intermediate := base.X.(*ast.SelectorExpr)
	isTestingPackage := intermediate.X.(*ast.Ident).Name == "testing"
	isTestingT := intermediate.Sel.Name == "T"

	return isTestingPackage && isTestingT
}
