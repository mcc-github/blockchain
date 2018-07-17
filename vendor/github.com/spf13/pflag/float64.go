package pflag

import "strconv"


type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func (f *float64Value) Type() string {
	return "float64"
}

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

func float64Conv(sval string) (interface{}, error) {
	return strconv.ParseFloat(sval, 64)
}


func (f *FlagSet) GetFloat64(name string) (float64, error) {
	val, err := f.getFlagType(name, "float64", float64Conv)
	if err != nil {
		return 0, err
	}
	return val.(float64), nil
}



func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	f.VarP(newFloat64Value(value, p), name, "", usage)
}


func (f *FlagSet) Float64VarP(p *float64, name, shorthand string, value float64, usage string) {
	f.VarP(newFloat64Value(value, p), name, shorthand, usage)
}



func Float64Var(p *float64, name string, value float64, usage string) {
	CommandLine.VarP(newFloat64Value(value, p), name, "", usage)
}


func Float64VarP(p *float64, name, shorthand string, value float64, usage string) {
	CommandLine.VarP(newFloat64Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Float64(name string, value float64, usage string) *float64 {
	p := new(float64)
	f.Float64VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Float64P(name, shorthand string, value float64, usage string) *float64 {
	p := new(float64)
	f.Float64VarP(p, name, shorthand, value, usage)
	return p
}



func Float64(name string, value float64, usage string) *float64 {
	return CommandLine.Float64P(name, "", value, usage)
}


func Float64P(name, shorthand string, value float64, usage string) *float64 {
	return CommandLine.Float64P(name, shorthand, value, usage)
}
