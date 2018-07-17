





package template

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)







func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}





func ParseFiles(filenames ...string) (*Template, error) {
	return parseFiles(nil, filenames...)
}




func (t *Template) ParseFiles(filenames ...string) (*Template, error) {
	return parseFiles(t, filenames...)
}



func parseFiles(t *Template, filenames ...string) (*Template, error) {
	if len(filenames) == 0 {
		
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filepath.Base(filename)
		
		
		
		
		
		
		var tmpl *Template
		if t == nil {
			t = New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}






func ParseGlob(pattern string) (*Template, error) {
	return parseGlob(nil, pattern)
}






func (t *Template) ParseGlob(pattern string) (*Template, error) {
	return parseGlob(t, pattern)
}


func parseGlob(t *Template, pattern string) (*Template, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
	}
	return parseFiles(t, filenames...)
}
