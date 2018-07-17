



package template

import (
	"fmt"
	"reflect"

	"github.com/alecthomas/template/parse"
)


type common struct {
	tmpl map[string]*Template
	
	
	
	parseFuncs FuncMap
	execFuncs  map[string]reflect.Value
}




type Template struct {
	name string
	*parse.Tree
	*common
	leftDelim  string
	rightDelim string
}


func New(name string) *Template {
	return &Template{
		name: name,
	}
}


func (t *Template) Name() string {
	return t.name
}




func (t *Template) New(name string) *Template {
	t.init()
	return &Template{
		name:       name,
		common:     t.common,
		leftDelim:  t.leftDelim,
		rightDelim: t.rightDelim,
	}
}

func (t *Template) init() {
	if t.common == nil {
		t.common = new(common)
		t.tmpl = make(map[string]*Template)
		t.parseFuncs = make(FuncMap)
		t.execFuncs = make(map[string]reflect.Value)
	}
}







func (t *Template) Clone() (*Template, error) {
	nt := t.copy(nil)
	nt.init()
	nt.tmpl[t.name] = nt
	for k, v := range t.tmpl {
		if k == t.name { 
			continue
		}
		
		tmpl := v.copy(nt.common)
		nt.tmpl[k] = tmpl
	}
	for k, v := range t.parseFuncs {
		nt.parseFuncs[k] = v
	}
	for k, v := range t.execFuncs {
		nt.execFuncs[k] = v
	}
	return nt, nil
}


func (t *Template) copy(c *common) *Template {
	nt := New(t.name)
	nt.Tree = t.Tree
	nt.common = c
	nt.leftDelim = t.leftDelim
	nt.rightDelim = t.rightDelim
	return nt
}



func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error) {
	if t.common != nil && t.tmpl[name] != nil {
		return nil, fmt.Errorf("template: redefinition of template %q", name)
	}
	nt := t.New(name)
	nt.Tree = tree
	t.tmpl[name] = nt
	return nt, nil
}



func (t *Template) Templates() []*Template {
	if t.common == nil {
		return nil
	}
	
	m := make([]*Template, 0, len(t.tmpl))
	for _, v := range t.tmpl {
		m = append(m, v)
	}
	return m
}






func (t *Template) Delims(left, right string) *Template {
	t.leftDelim = left
	t.rightDelim = right
	return t
}





func (t *Template) Funcs(funcMap FuncMap) *Template {
	t.init()
	addValueFuncs(t.execFuncs, funcMap)
	addFuncs(t.parseFuncs, funcMap)
	return t
}



func (t *Template) Lookup(name string) *Template {
	if t.common == nil {
		return nil
	}
	return t.tmpl[name]
}








func (t *Template) Parse(text string) (*Template, error) {
	t.init()
	trees, err := parse.Parse(t.name, text, t.leftDelim, t.rightDelim, t.parseFuncs, builtins)
	if err != nil {
		return nil, err
	}
	
	for name, tree := range trees {
		
		
		tmpl := t
		if name != t.name {
			tmpl = t.New(name)
		}
		
		if replace, err := t.associate(tmpl, tree); err != nil {
			return nil, err
		} else if replace {
			tmpl.Tree = tree
		}
		tmpl.leftDelim = t.leftDelim
		tmpl.rightDelim = t.rightDelim
	}
	return t, nil
}





func (t *Template) associate(new *Template, tree *parse.Tree) (bool, error) {
	if new.common != t.common {
		panic("internal error: associate not common")
	}
	name := new.name
	if old := t.tmpl[name]; old != nil {
		oldIsEmpty := parse.IsEmptyTree(old.Root)
		newIsEmpty := parse.IsEmptyTree(tree.Root)
		if newIsEmpty {
			
			return false, nil
		}
		if !oldIsEmpty {
			return false, fmt.Errorf("template: redefinition of template %q", name)
		}
	}
	t.tmpl[name] = new
	return true, nil
}
