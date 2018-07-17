package pflag

import (
	"fmt"
	"strconv"
	"strings"
)


type intSliceValue struct {
	value   *[]int
	changed bool
}

func newIntSliceValue(val []int, p *[]int) *intSliceValue {
	isv := new(intSliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *intSliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]int, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.Atoi(d)
		if err != nil {
			return err
		}

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *intSliceValue) Type() string {
	return "intSlice"
}

func (s *intSliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func intSliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	
	if len(val) == 0 {
		return []int{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]int, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.Atoi(d)
		if err != nil {
			return nil, err
		}

	}
	return out, nil
}


func (f *FlagSet) GetIntSlice(name string) ([]int, error) {
	val, err := f.getFlagType(name, "intSlice", intSliceConv)
	if err != nil {
		return []int{}, err
	}
	return val.([]int), nil
}



func (f *FlagSet) IntSliceVar(p *[]int, name string, value []int, usage string) {
	f.VarP(newIntSliceValue(value, p), name, "", usage)
}


func (f *FlagSet) IntSliceVarP(p *[]int, name, shorthand string, value []int, usage string) {
	f.VarP(newIntSliceValue(value, p), name, shorthand, usage)
}



func IntSliceVar(p *[]int, name string, value []int, usage string) {
	CommandLine.VarP(newIntSliceValue(value, p), name, "", usage)
}


func IntSliceVarP(p *[]int, name, shorthand string, value []int, usage string) {
	CommandLine.VarP(newIntSliceValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) IntSlice(name string, value []int, usage string) *[]int {
	p := []int{}
	f.IntSliceVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) IntSliceP(name, shorthand string, value []int, usage string) *[]int {
	p := []int{}
	f.IntSliceVarP(&p, name, shorthand, value, usage)
	return &p
}



func IntSlice(name string, value []int, usage string) *[]int {
	return CommandLine.IntSliceP(name, "", value, usage)
}


func IntSliceP(name, shorthand string, value []int, usage string) *[]int {
	return CommandLine.IntSliceP(name, shorthand, value, usage)
}
