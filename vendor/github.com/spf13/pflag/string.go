package pflag


type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string { return string(*s) }

func stringConv(sval string) (interface{}, error) {
	return sval, nil
}


func (f *FlagSet) GetString(name string) (string, error) {
	val, err := f.getFlagType(name, "string", stringConv)
	if err != nil {
		return "", err
	}
	return val.(string), nil
}



func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.VarP(newStringValue(value, p), name, "", usage)
}


func (f *FlagSet) StringVarP(p *string, name, shorthand string, value string, usage string) {
	f.VarP(newStringValue(value, p), name, shorthand, usage)
}



func StringVar(p *string, name string, value string, usage string) {
	CommandLine.VarP(newStringValue(value, p), name, "", usage)
}


func StringVarP(p *string, name, shorthand string, value string, usage string) {
	CommandLine.VarP(newStringValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) String(name string, value string, usage string) *string {
	p := new(string)
	f.StringVarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) StringP(name, shorthand string, value string, usage string) *string {
	p := new(string)
	f.StringVarP(p, name, shorthand, value, usage)
	return p
}



func String(name string, value string, usage string) *string {
	return CommandLine.StringP(name, "", value, usage)
}


func StringP(name, shorthand string, value string, usage string) *string {
	return CommandLine.StringP(name, shorthand, value, usage)
}
