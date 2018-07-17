package pflag

import "strconv"


type float32Value float32

func newFloat32Value(val float32, p *float32) *float32Value {
	*p = val
	return (*float32Value)(p)
}

func (f *float32Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	*f = float32Value(v)
	return err
}

func (f *float32Value) Type() string {
	return "float32"
}

func (f *float32Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 32) }

func float32Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseFloat(sval, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}


func (f *FlagSet) GetFloat32(name string) (float32, error) {
	val, err := f.getFlagType(name, "float32", float32Conv)
	if err != nil {
		return 0, err
	}
	return val.(float32), nil
}



func (f *FlagSet) Float32Var(p *float32, name string, value float32, usage string) {
	f.VarP(newFloat32Value(value, p), name, "", usage)
}


func (f *FlagSet) Float32VarP(p *float32, name, shorthand string, value float32, usage string) {
	f.VarP(newFloat32Value(value, p), name, shorthand, usage)
}



func Float32Var(p *float32, name string, value float32, usage string) {
	CommandLine.VarP(newFloat32Value(value, p), name, "", usage)
}


func Float32VarP(p *float32, name, shorthand string, value float32, usage string) {
	CommandLine.VarP(newFloat32Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) Float32(name string, value float32, usage string) *float32 {
	p := new(float32)
	f.Float32VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) Float32P(name, shorthand string, value float32, usage string) *float32 {
	p := new(float32)
	f.Float32VarP(p, name, shorthand, value, usage)
	return p
}



func Float32(name string, value float32, usage string) *float32 {
	return CommandLine.Float32P(name, "", value, usage)
}


func Float32P(name, shorthand string, value float32, usage string) *float32 {
	return CommandLine.Float32P(name, shorthand, value, usage)
}
