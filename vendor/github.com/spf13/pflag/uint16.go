package pflag

import "strconv"


type uint16Value uint16

func newUint16Value(val uint16, p *uint16) *uint16Value {
	*p = val
	return (*uint16Value)(p)
}

func (i *uint16Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 16)
	*i = uint16Value(v)
	return err
}

func (i *uint16Value) Type() string {
	return "uint16"
}

func (i *uint16Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uint16Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}


func (f *FlagSet) GetUint16(name string) (uint16, error) {
	val, err := f.getFlagType(name, "uint16", uint16Conv)
	if err != nil {
		return 0, err
	}
	return val.(uint16), nil
}



func (f *FlagSet) Uint16Var(p *uint16, name string, value uint16, usage string) {
	f.VarP(newUint16Value(value, p), name, "", usage)
}


func (f *FlagSet) Uint16VarP(p *uint16, name, shorthand string, value uint16, usage string) {
	f.VarP(newUint16Value(value, p), name, shorthand, usage)
}



func Uint16Var(p *uint16, name string, value uint16, usage string) {
	CommandLine.VarP(newUint16Value(value, p), name, "", usage)
}


func Uint16VarP(p *uint16, name, shorthand string, value uint16, usage string) {
	CommandLine.VarP(newUint16Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Uint16(name string, value uint16, usage string) *uint16 {
	p := new(uint16)
	f.Uint16VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Uint16P(name, shorthand string, value uint16, usage string) *uint16 {
	p := new(uint16)
	f.Uint16VarP(p, name, shorthand, value, usage)
	return p
}



func Uint16(name string, value uint16, usage string) *uint16 {
	return CommandLine.Uint16P(name, "", value, usage)
}


func Uint16P(name, shorthand string, value uint16, usage string) *uint16 {
	return CommandLine.Uint16P(name, shorthand, value, usage)
}
