package pflag

import (
	"io"
	"strconv"
	"strings"
)


type boolSliceValue struct {
	value   *[]bool
	changed bool
}

func newBoolSliceValue(val []bool, p *[]bool) *boolSliceValue {
	bsv := new(boolSliceValue)
	bsv.value = p
	*bsv.value = val
	return bsv
}



func (s *boolSliceValue) Set(val string) error {

	
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	
	boolStrSlice, err := readAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	
	out := make([]bool, 0, len(boolStrSlice))
	for _, boolStr := range boolStrSlice {
		b, err := strconv.ParseBool(strings.TrimSpace(boolStr))
		if err != nil {
			return err
		}
		out = append(out, b)
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

	return nil
}


func (s *boolSliceValue) Type() string {
	return "boolSlice"
}


func (s *boolSliceValue) String() string {

	boolStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		boolStrSlice[i] = strconv.FormatBool(b)
	}

	out, _ := writeAsCSV(boolStrSlice)

	return "[" + out + "]"
}

func boolSliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	
	if len(val) == 0 {
		return []bool{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]bool, len(ss))
	for i, t := range ss {
		var err error
		out[i], err = strconv.ParseBool(t)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}


func (f *FlagSet) GetBoolSlice(name string) ([]bool, error) {
	val, err := f.getFlagType(name, "boolSlice", boolSliceConv)
	if err != nil {
		return []bool{}, err
	}
	return val.([]bool), nil
}



func (f *FlagSet) BoolSliceVar(p *[]bool, name string, value []bool, usage string) {
	f.VarP(newBoolSliceValue(value, p), name, "", usage)
}


func (f *FlagSet) BoolSliceVarP(p *[]bool, name, shorthand string, value []bool, usage string) {
	f.VarP(newBoolSliceValue(value, p), name, shorthand, usage)
}



func BoolSliceVar(p *[]bool, name string, value []bool, usage string) {
	CommandLine.VarP(newBoolSliceValue(value, p), name, "", usage)
}


func BoolSliceVarP(p *[]bool, name, shorthand string, value []bool, usage string) {
	CommandLine.VarP(newBoolSliceValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) BoolSlice(name string, value []bool, usage string) *[]bool {
	p := []bool{}
	f.BoolSliceVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) BoolSliceP(name, shorthand string, value []bool, usage string) *[]bool {
	p := []bool{}
	f.BoolSliceVarP(&p, name, shorthand, value, usage)
	return &p
}



func BoolSlice(name string, value []bool, usage string) *[]bool {
	return CommandLine.BoolSliceP(name, "", value, usage)
}


func BoolSliceP(name, shorthand string, value []bool, usage string) *[]bool {
	return CommandLine.BoolSliceP(name, shorthand, value, usage)
}
