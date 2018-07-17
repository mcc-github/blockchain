package pflag

import "strconv"


type uint32Value uint32

func newUint32Value(val uint32, p *uint32) *uint32Value {
	*p = val
	return (*uint32Value)(p)
}

func (i *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	*i = uint32Value(v)
	return err
}

func (i *uint32Value) Type() string {
	return "uint32"
}

func (i *uint32Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uint32Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}


func (f *FlagSet) GetUint32(name string) (uint32, error) {
	val, err := f.getFlagType(name, "uint32", uint32Conv)
	if err != nil {
		return 0, err
	}
	return val.(uint32), nil
}



func (f *FlagSet) Uint32Var(p *uint32, name string, value uint32, usage string) {
	f.VarP(newUint32Value(value, p), name, "", usage)
}


func (f *FlagSet) Uint32VarP(p *uint32, name, shorthand string, value uint32, usage string) {
	f.VarP(newUint32Value(value, p), name, shorthand, usage)
}



func Uint32Var(p *uint32, name string, value uint32, usage string) {
	CommandLine.VarP(newUint32Value(value, p), name, "", usage)
}


func Uint32VarP(p *uint32, name, shorthand string, value uint32, usage string) {
	CommandLine.VarP(newUint32Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Uint32(name string, value uint32, usage string) *uint32 {
	p := new(uint32)
	f.Uint32VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Uint32P(name, shorthand string, value uint32, usage string) *uint32 {
	p := new(uint32)
	f.Uint32VarP(p, name, shorthand, value, usage)
	return p
}



func Uint32(name string, value uint32, usage string) *uint32 {
	return CommandLine.Uint32P(name, "", value, usage)
}


func Uint32P(name, shorthand string, value uint32, usage string) *uint32 {
	return CommandLine.Uint32P(name, shorthand, value, usage)
}
