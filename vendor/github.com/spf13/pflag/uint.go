package pflag

import "strconv"


type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

func (i *uintValue) Type() string {
	return "uint"
}

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

func uintConv(sval string) (interface{}, error) {
	v, err := strconv.ParseUint(sval, 0, 0)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}


func (f *FlagSet) GetUint(name string) (uint, error) {
	val, err := f.getFlagType(name, "uint", uintConv)
	if err != nil {
		return 0, err
	}
	return val.(uint), nil
}



func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	f.VarP(newUintValue(value, p), name, "", usage)
}


func (f *FlagSet) UintVarP(p *uint, name, shorthand string, value uint, usage string) {
	f.VarP(newUintValue(value, p), name, shorthand, usage)
}



func UintVar(p *uint, name string, value uint, usage string) {
	CommandLine.VarP(newUintValue(value, p), name, "", usage)
}


func UintVarP(p *uint, name, shorthand string, value uint, usage string) {
	CommandLine.VarP(newUintValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) Uint(name string, value uint, usage string) *uint {
	p := new(uint)
	f.UintVarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) UintP(name, shorthand string, value uint, usage string) *uint {
	p := new(uint)
	f.UintVarP(p, name, shorthand, value, usage)
	return p
}



func Uint(name string, value uint, usage string) *uint {
	return CommandLine.UintP(name, "", value, usage)
}


func UintP(name, shorthand string, value uint, usage string) *uint {
	return CommandLine.UintP(name, shorthand, value, usage)
}
