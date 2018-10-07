package pflag

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)


type stringToIntValue struct {
	value   *map[string]int
	changed bool
}

func newStringToIntValue(val map[string]int, p *map[string]int) *stringToIntValue {
	ssv := new(stringToIntValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}


func (s *stringToIntValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make(map[string]int, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("%s must be formatted as key=value", pair)
		}
		var err error
		out[kv[0]], err = strconv.Atoi(kv[1])
		if err != nil {
			return err
		}
	}
	if !s.changed {
		*s.value = out
	} else {
		for k, v := range out {
			(*s.value)[k] = v
		}
	}
	s.changed = true
	return nil
}

func (s *stringToIntValue) Type() string {
	return "stringToInt"
}

func (s *stringToIntValue) String() string {
	var buf bytes.Buffer
	i := 0
	for k, v := range *s.value {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.Itoa(v))
		i++
	}
	return "[" + buf.String() + "]"
}

func stringToIntConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	
	if len(val) == 0 {
		return map[string]int{}, nil
	}
	ss := strings.Split(val, ",")
	out := make(map[string]int, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("%s must be formatted as key=value", pair)
		}
		var err error
		out[kv[0]], err = strconv.Atoi(kv[1])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}


func (f *FlagSet) GetStringToInt(name string) (map[string]int, error) {
	val, err := f.getFlagType(name, "stringToInt", stringToIntConv)
	if err != nil {
		return map[string]int{}, err
	}
	return val.(map[string]int), nil
}




func (f *FlagSet) StringToIntVar(p *map[string]int, name string, value map[string]int, usage string) {
	f.VarP(newStringToIntValue(value, p), name, "", usage)
}


func (f *FlagSet) StringToIntVarP(p *map[string]int, name, shorthand string, value map[string]int, usage string) {
	f.VarP(newStringToIntValue(value, p), name, shorthand, usage)
}




func StringToIntVar(p *map[string]int, name string, value map[string]int, usage string) {
	CommandLine.VarP(newStringToIntValue(value, p), name, "", usage)
}


func StringToIntVarP(p *map[string]int, name, shorthand string, value map[string]int, usage string) {
	CommandLine.VarP(newStringToIntValue(value, p), name, shorthand, usage)
}




func (f *FlagSet) StringToInt(name string, value map[string]int, usage string) *map[string]int {
	p := map[string]int{}
	f.StringToIntVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) StringToIntP(name, shorthand string, value map[string]int, usage string) *map[string]int {
	p := map[string]int{}
	f.StringToIntVarP(&p, name, shorthand, value, usage)
	return &p
}




func StringToInt(name string, value map[string]int, usage string) *map[string]int {
	return CommandLine.StringToIntP(name, "", value, usage)
}


func StringToIntP(name, shorthand string, value map[string]int, usage string) *map[string]int {
	return CommandLine.StringToIntP(name, shorthand, value, usage)
}
