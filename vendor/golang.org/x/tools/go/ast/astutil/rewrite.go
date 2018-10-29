



package astutil

import (
	"fmt"
	"go/ast"
	"reflect"
	"sort"
)







type ApplyFunc func(*Cursor) bool






















func Apply(root ast.Node, pre, post ApplyFunc) (result ast.Node) {
	parent := &struct{ ast.Node }{root}
	defer func() {
		if r := recover(); r != nil && r != abort {
			panic(r)
		}
		result = parent.Node
	}()
	a := &application{pre: pre, post: post}
	a.apply(parent, "Node", nil, root)
	return
}

var abort = new(int) 














type Cursor struct {
	parent ast.Node
	name   string
	iter   *iterator 
	node   ast.Node
}


func (c *Cursor) Node() ast.Node { return c.node }


func (c *Cursor) Parent() ast.Node { return c.parent }




func (c *Cursor) Name() string { return c.name }





func (c *Cursor) Index() int {
	if c.iter != nil {
		return c.iter.index
	}
	return -1
}


func (c *Cursor) field() reflect.Value {
	return reflect.Indirect(reflect.ValueOf(c.parent)).FieldByName(c.name)
}



func (c *Cursor) Replace(n ast.Node) {
	if _, ok := c.node.(*ast.File); ok {
		file, ok := n.(*ast.File)
		if !ok {
			panic("attempt to replace *ast.File with non-*ast.File")
		}
		c.parent.(*ast.Package).Files[c.name] = file
		return
	}

	v := c.field()
	if i := c.Index(); i >= 0 {
		v = v.Index(i)
	}
	v.Set(reflect.ValueOf(n))
}





func (c *Cursor) Delete() {
	if _, ok := c.node.(*ast.File); ok {
		delete(c.parent.(*ast.Package).Files, c.name)
		return
	}

	i := c.Index()
	if i < 0 {
		panic("Delete node not contained in slice")
	}
	v := c.field()
	l := v.Len()
	reflect.Copy(v.Slice(i, l), v.Slice(i+1, l))
	v.Index(l - 1).Set(reflect.Zero(v.Type().Elem()))
	v.SetLen(l - 1)
	c.iter.step--
}




func (c *Cursor) InsertAfter(n ast.Node) {
	i := c.Index()
	if i < 0 {
		panic("InsertAfter node not contained in slice")
	}
	v := c.field()
	v.Set(reflect.Append(v, reflect.Zero(v.Type().Elem())))
	l := v.Len()
	reflect.Copy(v.Slice(i+2, l), v.Slice(i+1, l))
	v.Index(i + 1).Set(reflect.ValueOf(n))
	c.iter.step++
}




func (c *Cursor) InsertBefore(n ast.Node) {
	i := c.Index()
	if i < 0 {
		panic("InsertBefore node not contained in slice")
	}
	v := c.field()
	v.Set(reflect.Append(v, reflect.Zero(v.Type().Elem())))
	l := v.Len()
	reflect.Copy(v.Slice(i+1, l), v.Slice(i, l))
	v.Index(i).Set(reflect.ValueOf(n))
	c.iter.index++
}


type application struct {
	pre, post ApplyFunc
	cursor    Cursor
	iter      iterator
}

func (a *application) apply(parent ast.Node, name string, iter *iterator, n ast.Node) {
	
	if v := reflect.ValueOf(n); v.Kind() == reflect.Ptr && v.IsNil() {
		n = nil
	}

	
	saved := a.cursor
	a.cursor.parent = parent
	a.cursor.name = name
	a.cursor.iter = iter
	a.cursor.node = n

	if a.pre != nil && !a.pre(&a.cursor) {
		a.cursor = saved
		return
	}

	
	
	switch n := n.(type) {
	case nil:
		

	
	case *ast.Comment:
		

	case *ast.CommentGroup:
		if n != nil {
			a.applyList(n, "List")
		}

	case *ast.Field:
		a.apply(n, "Doc", nil, n.Doc)
		a.applyList(n, "Names")
		a.apply(n, "Type", nil, n.Type)
		a.apply(n, "Tag", nil, n.Tag)
		a.apply(n, "Comment", nil, n.Comment)

	case *ast.FieldList:
		a.applyList(n, "List")

	
	case *ast.BadExpr, *ast.Ident, *ast.BasicLit:
		

	case *ast.Ellipsis:
		a.apply(n, "Elt", nil, n.Elt)

	case *ast.FuncLit:
		a.apply(n, "Type", nil, n.Type)
		a.apply(n, "Body", nil, n.Body)

	case *ast.CompositeLit:
		a.apply(n, "Type", nil, n.Type)
		a.applyList(n, "Elts")

	case *ast.ParenExpr:
		a.apply(n, "X", nil, n.X)

	case *ast.SelectorExpr:
		a.apply(n, "X", nil, n.X)
		a.apply(n, "Sel", nil, n.Sel)

	case *ast.IndexExpr:
		a.apply(n, "X", nil, n.X)
		a.apply(n, "Index", nil, n.Index)

	case *ast.SliceExpr:
		a.apply(n, "X", nil, n.X)
		a.apply(n, "Low", nil, n.Low)
		a.apply(n, "High", nil, n.High)
		a.apply(n, "Max", nil, n.Max)

	case *ast.TypeAssertExpr:
		a.apply(n, "X", nil, n.X)
		a.apply(n, "Type", nil, n.Type)

	case *ast.CallExpr:
		a.apply(n, "Fun", nil, n.Fun)
		a.applyList(n, "Args")

	case *ast.StarExpr:
		a.apply(n, "X", nil, n.X)

	case *ast.UnaryExpr:
		a.apply(n, "X", nil, n.X)

	case *ast.BinaryExpr:
		a.apply(n, "X", nil, n.X)
		a.apply(n, "Y", nil, n.Y)

	case *ast.KeyValueExpr:
		a.apply(n, "Key", nil, n.Key)
		a.apply(n, "Value", nil, n.Value)

	
	case *ast.ArrayType:
		a.apply(n, "Len", nil, n.Len)
		a.apply(n, "Elt", nil, n.Elt)

	case *ast.StructType:
		a.apply(n, "Fields", nil, n.Fields)

	case *ast.FuncType:
		a.apply(n, "Params", nil, n.Params)
		a.apply(n, "Results", nil, n.Results)

	case *ast.InterfaceType:
		a.apply(n, "Methods", nil, n.Methods)

	case *ast.MapType:
		a.apply(n, "Key", nil, n.Key)
		a.apply(n, "Value", nil, n.Value)

	case *ast.ChanType:
		a.apply(n, "Value", nil, n.Value)

	
	case *ast.BadStmt:
		

	case *ast.DeclStmt:
		a.apply(n, "Decl", nil, n.Decl)

	case *ast.EmptyStmt:
		

	case *ast.LabeledStmt:
		a.apply(n, "Label", nil, n.Label)
		a.apply(n, "Stmt", nil, n.Stmt)

	case *ast.ExprStmt:
		a.apply(n, "X", nil, n.X)

	case *ast.SendStmt:
		a.apply(n, "Chan", nil, n.Chan)
		a.apply(n, "Value", nil, n.Value)

	case *ast.IncDecStmt:
		a.apply(n, "X", nil, n.X)

	case *ast.AssignStmt:
		a.applyList(n, "Lhs")
		a.applyList(n, "Rhs")

	case *ast.GoStmt:
		a.apply(n, "Call", nil, n.Call)

	case *ast.DeferStmt:
		a.apply(n, "Call", nil, n.Call)

	case *ast.ReturnStmt:
		a.applyList(n, "Results")

	case *ast.BranchStmt:
		a.apply(n, "Label", nil, n.Label)

	case *ast.BlockStmt:
		a.applyList(n, "List")

	case *ast.IfStmt:
		a.apply(n, "Init", nil, n.Init)
		a.apply(n, "Cond", nil, n.Cond)
		a.apply(n, "Body", nil, n.Body)
		a.apply(n, "Else", nil, n.Else)

	case *ast.CaseClause:
		a.applyList(n, "List")
		a.applyList(n, "Body")

	case *ast.SwitchStmt:
		a.apply(n, "Init", nil, n.Init)
		a.apply(n, "Tag", nil, n.Tag)
		a.apply(n, "Body", nil, n.Body)

	case *ast.TypeSwitchStmt:
		a.apply(n, "Init", nil, n.Init)
		a.apply(n, "Assign", nil, n.Assign)
		a.apply(n, "Body", nil, n.Body)

	case *ast.CommClause:
		a.apply(n, "Comm", nil, n.Comm)
		a.applyList(n, "Body")

	case *ast.SelectStmt:
		a.apply(n, "Body", nil, n.Body)

	case *ast.ForStmt:
		a.apply(n, "Init", nil, n.Init)
		a.apply(n, "Cond", nil, n.Cond)
		a.apply(n, "Post", nil, n.Post)
		a.apply(n, "Body", nil, n.Body)

	case *ast.RangeStmt:
		a.apply(n, "Key", nil, n.Key)
		a.apply(n, "Value", nil, n.Value)
		a.apply(n, "X", nil, n.X)
		a.apply(n, "Body", nil, n.Body)

	
	case *ast.ImportSpec:
		a.apply(n, "Doc", nil, n.Doc)
		a.apply(n, "Name", nil, n.Name)
		a.apply(n, "Path", nil, n.Path)
		a.apply(n, "Comment", nil, n.Comment)

	case *ast.ValueSpec:
		a.apply(n, "Doc", nil, n.Doc)
		a.applyList(n, "Names")
		a.apply(n, "Type", nil, n.Type)
		a.applyList(n, "Values")
		a.apply(n, "Comment", nil, n.Comment)

	case *ast.TypeSpec:
		a.apply(n, "Doc", nil, n.Doc)
		a.apply(n, "Name", nil, n.Name)
		a.apply(n, "Type", nil, n.Type)
		a.apply(n, "Comment", nil, n.Comment)

	case *ast.BadDecl:
		

	case *ast.GenDecl:
		a.apply(n, "Doc", nil, n.Doc)
		a.applyList(n, "Specs")

	case *ast.FuncDecl:
		a.apply(n, "Doc", nil, n.Doc)
		a.apply(n, "Recv", nil, n.Recv)
		a.apply(n, "Name", nil, n.Name)
		a.apply(n, "Type", nil, n.Type)
		a.apply(n, "Body", nil, n.Body)

	
	case *ast.File:
		a.apply(n, "Doc", nil, n.Doc)
		a.apply(n, "Name", nil, n.Name)
		a.applyList(n, "Decls")
		
		

	case *ast.Package:
		
		var names []string
		for name := range n.Files {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			a.apply(n, name, nil, n.Files[name])
		}

	default:
		panic(fmt.Sprintf("Apply: unexpected node type %T", n))
	}

	if a.post != nil && !a.post(&a.cursor) {
		panic(abort)
	}

	a.cursor = saved
}


type iterator struct {
	index, step int
}

func (a *application) applyList(parent ast.Node, name string) {
	
	saved := a.iter
	a.iter.index = 0
	for {
		
		v := reflect.Indirect(reflect.ValueOf(parent)).FieldByName(name)
		if a.iter.index >= v.Len() {
			break
		}

		
		var x ast.Node
		if e := v.Index(a.iter.index); e.IsValid() {
			x = e.Interface().(ast.Node)
		}

		a.iter.step = 1
		a.apply(parent, name, &a.iter, x)
		a.iter.index += a.iter.step
	}
	a.iter = saved
}
