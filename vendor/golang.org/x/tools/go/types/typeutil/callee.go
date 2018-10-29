



package typeutil

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
)



func Callee(info *types.Info, call *ast.CallExpr) types.Object {
	var obj types.Object
	switch fun := astutil.Unparen(call.Fun).(type) {
	case *ast.Ident:
		obj = info.Uses[fun] 
	case *ast.SelectorExpr:
		if sel, ok := info.Selections[fun]; ok {
			obj = sel.Obj() 
		} else {
			obj = info.Uses[fun.Sel] 
		}
	}
	if _, ok := obj.(*types.TypeName); ok {
		return nil 
	}
	return obj
}



func StaticCallee(info *types.Info, call *ast.CallExpr) *types.Func {
	if f, ok := Callee(info, call).(*types.Func); ok && !interfaceMethod(f) {
		return f
	}
	return nil
}

func interfaceMethod(f *types.Func) bool {
	recv := f.Type().(*types.Signature).Recv()
	return recv != nil && types.IsInterface(recv.Type())
}
