package pflag

import "strconv"


type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

func intConv(sval string) (interface{}, error) {
	return strconv.Atoi(sval)
}


func (f *FlagSet) GetInt(name string) (int, error) {
	val, err := f.getFlagType(name, "int", intConv)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}



func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.VarP(newIntValue(value, p), name, "", usage)
}


func (f *FlagSet) IntVarP(p *int, name, shorthand string, value int, usage string) {
	f.VarP(newIntValue(value, p), name, shorthand, usage)
}



func IntVar(p *int, name string, value int, usage string) {
	CommandLine.VarP(newIntValue(value, p), name, "", usage)
}


func IntVarP(p *int, name, shorthand string, value int, usage string) {
	CommandLine.VarP(newIntValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) Int(name string, value int, usage string) *int {
	p := new(int)
	f.IntVarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) IntP(name, shorthand string, value int, usage string) *int {
	p := new(int)
	f.IntVarP(p, name, shorthand, value, usage)
	return p
}



func Int(name string, value int, usage string) *int {
	return CommandLine.IntP(name, "", value, usage)
}


func IntP(name, shorthand string, value int, usage string) *int {
	return CommandLine.IntP(name, shorthand, value, usage)
}
