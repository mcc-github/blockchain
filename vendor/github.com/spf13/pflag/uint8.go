package pflag

import "strconv"


type uint8Value uint8

func newUint8Value(val uint8, p *uint8) *uint8Value {
	*p = val
	return (*uint8Value)(p)
}

func (i *uint8Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 8)
	*i = uint8Value(v)
	return err
}

func (i *uint8Value) Type() string {
	return "uint8"
}

func (i *uint8Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uint8Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(v), nil
}


func (f *FlagSet) GetUint8(name string) (uint8, error) {
	val, err := f.getFlagType(name, "uint8", uint8Conv)
	if err != nil {
		return 0, err
	}
	return val.(uint8), nil
}



func (f *FlagSet) Uint8Var(p *uint8, name string, value uint8, usage string) {
	f.VarP(newUint8Value(value, p), name, "", usage)
}


func (f *FlagSet) Uint8VarP(p *uint8, name, shorthand string, value uint8, usage string) {
	f.VarP(newUint8Value(value, p), name, shorthand, usage)
}



func Uint8Var(p *uint8, name string, value uint8, usage string) {
	CommandLine.VarP(newUint8Value(value, p), name, "", usage)
}


func Uint8VarP(p *uint8, name, shorthand string, value uint8, usage string) {
	CommandLine.VarP(newUint8Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Uint8(name string, value uint8, usage string) *uint8 {
	p := new(uint8)
	f.Uint8VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Uint8P(name, shorthand string, value uint8, usage string) *uint8 {
	p := new(uint8)
	f.Uint8VarP(p, name, shorthand, value, usage)
	return p
}



func Uint8(name string, value uint8, usage string) *uint8 {
	return CommandLine.Uint8P(name, "", value, usage)
}


func Uint8P(name, shorthand string, value uint8, usage string) *uint8 {
	return CommandLine.Uint8P(name, shorthand, value, usage)
}
