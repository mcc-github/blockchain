package pflag

import "strconv"


type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Type() string {
	return "uint64"
}

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uint64Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 64)
	if err != nil {
		return 0, err
	}
	return uint64(v), nil
}


func (f *FlagSet) GetUint64(name string) (uint64, error) {
	val, err := f.getFlagType(name, "uint64", uint64Conv)
	if err != nil {
		return 0, err
	}
	return val.(uint64), nil
}



func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	f.VarP(newUint64Value(value, p), name, "", usage)
}


func (f *FlagSet) Uint64VarP(p *uint64, name, shorthand string, value uint64, usage string) {
	f.VarP(newUint64Value(value, p), name, shorthand, usage)
}



func Uint64Var(p *uint64, name string, value uint64, usage string) {
	CommandLine.VarP(newUint64Value(value, p), name, "", usage)
}


func Uint64VarP(p *uint64, name, shorthand string, value uint64, usage string) {
	CommandLine.VarP(newUint64Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	p := new(uint64)
	f.Uint64VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Uint64P(name, shorthand string, value uint64, usage string) *uint64 {
	p := new(uint64)
	f.Uint64VarP(p, name, shorthand, value, usage)
	return p
}



func Uint64(name string, value uint64, usage string) *uint64 {
	return CommandLine.Uint64P(name, "", value, usage)
}


func Uint64P(name, shorthand string, value uint64, usage string) *uint64 {
	return CommandLine.Uint64P(name, shorthand, value, usage)
}
