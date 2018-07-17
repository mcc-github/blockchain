package pflag

import "strconv"


type int8Value int8

func newInt8Value(val int8, p *int8) *int8Value {
	*p = val
	return (*int8Value)(p)
}

func (i *int8Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 8)
	*i = int8Value(v)
	return err
}

func (i *int8Value) Type() string {
	return "int8"
}

func (i *int8Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int8Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseInt(sval, 0, 8)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}


func (f *FlagSet) GetInt8(name string) (int8, error) {
	val, err := f.getFlagType(name, "int8", int8Conv)
	if err != nil {
		return 0, err
	}
	return val.(int8), nil
}



func (f *FlagSet) Int8Var(p *int8, name string, value int8, usage string) {
	f.VarP(newInt8Value(value, p), name, "", usage)
}


func (f *FlagSet) Int8VarP(p *int8, name, shorthand string, value int8, usage string) {
	f.VarP(newInt8Value(value, p), name, shorthand, usage)
}



func Int8Var(p *int8, name string, value int8, usage string) {
	CommandLine.VarP(newInt8Value(value, p), name, "", usage)
}


func Int8VarP(p *int8, name, shorthand string, value int8, usage string) {
	CommandLine.VarP(newInt8Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Int8(name string, value int8, usage string) *int8 {
	p := new(int8)
	f.Int8VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Int8P(name, shorthand string, value int8, usage string) *int8 {
	p := new(int8)
	f.Int8VarP(p, name, shorthand, value, usage)
	return p
}



func Int8(name string, value int8, usage string) *int8 {
	return CommandLine.Int8P(name, "", value, usage)
}


func Int8P(name, shorthand string, value int8, usage string) *int8 {
	return CommandLine.Int8P(name, shorthand, value, usage)
}
