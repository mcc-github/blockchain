package pflag

import "strconv"



type boolFlag interface {
	Value
	IsBoolFlag() bool
}


type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) Type() string {
	return "bool"
}

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

func boolConv(sval string) (interface{}, error) {
	return strconv.ParseBool(sval)
}


func (f *FlagSet) GetBool(name string) (bool, error) {
	val, err := f.getFlagType(name, "bool", boolConv)
	if err != nil {
		return false, err
	}
	return val.(bool), nil
}



func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	f.BoolVarP(p, name, "", value, usage)
}


func (f *FlagSet) BoolVarP(p *bool, name, shorthand string, value bool, usage string) {
	flag := f.VarPF(newBoolValue(value, p), name, shorthand, usage)
	flag.NoOptDefVal = "true"
}



func BoolVar(p *bool, name string, value bool, usage string) {
	BoolVarP(p, name, "", value, usage)
}


func BoolVarP(p *bool, name, shorthand string, value bool, usage string) {
	flag := CommandLine.VarPF(newBoolValue(value, p), name, shorthand, usage)
	flag.NoOptDefVal = "true"
}



func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	return f.BoolP(name, "", value, usage)
}


func (f *FlagSet) BoolP(name, shorthand string, value bool, usage string) *bool {
	p := new(bool)
	f.BoolVarP(p, name, shorthand, value, usage)
	return p
}



func Bool(name string, value bool, usage string) *bool {
	return BoolP(name, "", value, usage)
}


func BoolP(name, shorthand string, value bool, usage string) *bool {
	b := CommandLine.BoolP(name, shorthand, value, usage)
	return b
}
