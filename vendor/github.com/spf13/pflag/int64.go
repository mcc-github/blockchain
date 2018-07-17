package pflag

import "strconv"


type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Type() string {
	return "int64"
}

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int64Conv(sval string) (interface{}, error) {
	return strconv.ParseInt(sval, 0, 64)
}


func (f *FlagSet) GetInt64(name string) (int64, error) {
	val, err := f.getFlagType(name, "int64", int64Conv)
	if err != nil {
		return 0, err
	}
	return val.(int64), nil
}



func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	f.VarP(newInt64Value(value, p), name, "", usage)
}


func (f *FlagSet) Int64VarP(p *int64, name, shorthand string, value int64, usage string) {
	f.VarP(newInt64Value(value, p), name, shorthand, usage)
}



func Int64Var(p *int64, name string, value int64, usage string) {
	CommandLine.VarP(newInt64Value(value, p), name, "", usage)
}


func Int64VarP(p *int64, name, shorthand string, value int64, usage string) {
	CommandLine.VarP(newInt64Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Int64(name string, value int64, usage string) *int64 {
	p := new(int64)
	f.Int64VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Int64P(name, shorthand string, value int64, usage string) *int64 {
	p := new(int64)
	f.Int64VarP(p, name, shorthand, value, usage)
	return p
}



func Int64(name string, value int64, usage string) *int64 {
	return CommandLine.Int64P(name, "", value, usage)
}


func Int64P(name, shorthand string, value int64, usage string) *int64 {
	return CommandLine.Int64P(name, shorthand, value, usage)
}
