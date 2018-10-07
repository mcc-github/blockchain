package pflag

import "strconv"


type int16Value int16

func newInt16Value(val int16, p *int16) *int16Value {
	*p = val
	return (*int16Value)(p)
}

func (i *int16Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 16)
	*i = int16Value(v)
	return err
}

func (i *int16Value) Type() string {
	return "int16"
}

func (i *int16Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int16Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseInt(sval, 0, 16)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}


func (f *FlagSet) GetInt16(name string) (int16, error) {
	val, err := f.getFlagType(name, "int16", int16Conv)
	if err != nil {
		return 0, err
	}
	return val.(int16), nil
}



func (f *FlagSet) Int16Var(p *int16, name string, value int16, usage string) {
	f.VarP(newInt16Value(value, p), name, "", usage)
}


func (f *FlagSet) Int16VarP(p *int16, name, shorthand string, value int16, usage string) {
	f.VarP(newInt16Value(value, p), name, shorthand, usage)
}



func Int16Var(p *int16, name string, value int16, usage string) {
	CommandLine.VarP(newInt16Value(value, p), name, "", usage)
}


func Int16VarP(p *int16, name, shorthand string, value int16, usage string) {
	CommandLine.VarP(newInt16Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Int16(name string, value int16, usage string) *int16 {
	p := new(int16)
	f.Int16VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Int16P(name, shorthand string, value int16, usage string) *int16 {
	p := new(int16)
	f.Int16VarP(p, name, shorthand, value, usage)
	return p
}



func Int16(name string, value int16, usage string) *int16 {
	return CommandLine.Int16P(name, "", value, usage)
}


func Int16P(name, shorthand string, value int16, usage string) *int16 {
	return CommandLine.Int16P(name, shorthand, value, usage)
}
