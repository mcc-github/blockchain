package pflag

import "strconv"


type int32Value int32

func newInt32Value(val int32, p *int32) *int32Value {
	*p = val
	return (*int32Value)(p)
}

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
}

func (i *int32Value) Type() string {
	return "int32"
}

func (i *int32Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int32Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseInt(sval, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}


func (f *FlagSet) GetInt32(name string) (int32, error) {
	val, err := f.getFlagType(name, "int32", int32Conv)
	if err != nil {
		return 0, err
	}
	return val.(int32), nil
}



func (f *FlagSet) Int32Var(p *int32, name string, value int32, usage string) {
	f.VarP(newInt32Value(value, p), name, "", usage)
}


func (f *FlagSet) Int32VarP(p *int32, name, shorthand string, value int32, usage string) {
	f.VarP(newInt32Value(value, p), name, shorthand, usage)
}



func Int32Var(p *int32, name string, value int32, usage string) {
	CommandLine.VarP(newInt32Value(value, p), name, "", usage)
}


func Int32VarP(p *int32, name, shorthand string, value int32, usage string) {
	CommandLine.VarP(newInt32Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Int32(name string, value int32, usage string) *int32 {
	p := new(int32)
	f.Int32VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Int32P(name, shorthand string, value int32, usage string) *int32 {
	p := new(int32)
	f.Int32VarP(p, name, shorthand, value, usage)
	return p
}



func Int32(name string, value int32, usage string) *int32 {
	return CommandLine.Int32P(name, "", value, usage)
}


func Int32P(name, shorthand string, value int32, usage string) *int32 {
	return CommandLine.Int32P(name, shorthand, value, usage)
}
