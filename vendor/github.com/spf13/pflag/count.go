package pflag

import "strconv"


type countValue int

func newCountValue(val int, p *int) *countValue {
	*p = val
	return (*countValue)(p)
}

func (i *countValue) Set(s string) error {
	
	if s == "+1" {
		*i = countValue(*i + 1)
		return nil
	}
	v, err := strconv.ParseInt(s, 0, 0)
	*i = countValue(v)
	return err
}

func (i *countValue) Type() string {
	return "count"
}

func (i *countValue) String() string { return strconv.Itoa(int(*i)) }

func countConv(sval string) (interface{}, error) {
	i, err := strconv.Atoi(sval)
	if err != nil {
		return nil, err
	}
	return i, nil
}


func (f *FlagSet) GetCount(name string) (int, error) {
	val, err := f.getFlagType(name, "count", countConv)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}




func (f *FlagSet) CountVar(p *int, name string, usage string) {
	f.CountVarP(p, name, "", usage)
}


func (f *FlagSet) CountVarP(p *int, name, shorthand string, usage string) {
	flag := f.VarPF(newCountValue(0, p), name, shorthand, usage)
	flag.NoOptDefVal = "+1"
}


func CountVar(p *int, name string, usage string) {
	CommandLine.CountVar(p, name, usage)
}


func CountVarP(p *int, name, shorthand string, usage string) {
	CommandLine.CountVarP(p, name, shorthand, usage)
}




func (f *FlagSet) Count(name string, usage string) *int {
	p := new(int)
	f.CountVarP(p, name, "", usage)
	return p
}


func (f *FlagSet) CountP(name, shorthand string, usage string) *int {
	p := new(int)
	f.CountVarP(p, name, shorthand, usage)
	return p
}




func Count(name string, usage string) *int {
	return CommandLine.CountP(name, "", usage)
}


func CountP(name, shorthand string, usage string) *int {
	return CommandLine.CountP(name, shorthand, usage)
}
