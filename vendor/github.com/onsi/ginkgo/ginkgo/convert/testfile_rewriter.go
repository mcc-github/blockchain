package convert

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
)


func rewriteTestsInFile(pathToFile string) {
	fileSet := token.NewFileSet()
	rootNode, err := parser.ParseFile(fileSet, pathToFile, nil, 0)
	if err != nil {
		panic(fmt.Sprintf("Error parsing test file '%s':\n%s\n", pathToFile, err.Error()))
	}

	addGinkgoImports(rootNode)
	removeTestingImport(rootNode)

	varUnderscoreBlock := createVarUnderscoreBlock()
	describeBlock := createDescribeBlock()
	varUnderscoreBlock.Values = []ast.Expr{describeBlock}

	for _, testFunc := range findTestFuncs(rootNode) {
		rewriteTestFuncAsItStatement(testFunc, rootNode, describeBlock)
	}

	underscoreDecl := &ast.GenDecl{
		Tok:    85, 
		TokPos: 14, 
		Specs:  []ast.Spec{varUnderscoreBlock},
	}

	imports := rootNode.Decls[0]
	tail := rootNode.Decls[1:]
	rootNode.Decls = append(append([]ast.Decl{imports}, underscoreDecl), tail...)
	rewriteOtherFuncsToUseGinkgoT(rootNode.Decls)
	walkNodesInRootNodeReplacingTestingT(rootNode)

	var buffer bytes.Buffer
	if err = format.Node(&buffer, fileSet, rootNode); err != nil {
		panic(fmt.Sprintf("Error formatting ast node after rewriting tests.\n%s\n", err.Error()))
	}

	fileInfo, err := os.Stat(pathToFile)
	if err != nil {
		panic(fmt.Sprintf("Error stat'ing file: %s\n", pathToFile))
	}

	ioutil.WriteFile(pathToFile, buffer.Bytes(), fileInfo.Mode())
	return
}


func rewriteTestFuncAsItStatement(testFunc *ast.FuncDecl, rootNode *ast.File, describe *ast.CallExpr) {
	var funcIndex int = -1
	for index, child := range rootNode.Decls {
		if child == testFunc {
			funcIndex = index
			break
		}
	}

	if funcIndex < 0 {
		panic(fmt.Sprintf("Assert failed: Error finding index for test node %s\n", testFunc.Name.Name))
	}

	var block *ast.BlockStmt = blockStatementFromDescribe(describe)
	block.List = append(block.List, createItStatementForTestFunc(testFunc))
	replaceTestingTsWithGinkgoT(block, namedTestingTArg(testFunc))

	
	rootNode.Decls = append(rootNode.Decls[:funcIndex], rootNode.Decls[funcIndex+1:]...)
	return
}


func replaceTestingTsWithGinkgoT(statementsBlock *ast.BlockStmt, testingT string) {
	ast.Inspect(statementsBlock, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		keyValueExpr, ok := node.(*ast.KeyValueExpr)
		if ok {
			replaceNamedTestingTsInKeyValueExpression(keyValueExpr, testingT)
			return true
		}

		funcLiteral, ok := node.(*ast.FuncLit)
		if ok {
			replaceTypeDeclTestingTsInFuncLiteral(funcLiteral)
			return true
		}

		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		replaceTestingTsInArgsLists(callExpr, testingT)

		funCall, ok := callExpr.Fun.(*ast.SelectorExpr)
		if ok {
			replaceTestingTsMethodCalls(funCall, testingT)
		}

		return true
	})
}


func replaceTestingTsMethodCalls(selectorExpr *ast.SelectorExpr, testingT string) {
	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return
	}

	if ident.Name == testingT {
		selectorExpr.X = newGinkgoTFromIdent(ident)
	}
}


func replaceTestingTsInArgsLists(callExpr *ast.CallExpr, testingT string) {
	for index, arg := range callExpr.Args {
		ident, ok := arg.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name == testingT {
			callExpr.Args[index] = newGinkgoTFromIdent(ident)
		}
	}
}
