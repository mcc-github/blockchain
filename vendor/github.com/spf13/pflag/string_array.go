package pflag


type stringArrayValue struct {
	value   *[]string
	changed bool
}

func newStringArrayValue(val []string, p *[]string) *stringArrayValue {
	ssv := new(stringArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringArrayValue) Set(val string) error {
	if !s.changed {
		*s.value = []string{val}
		s.changed = true
	} else {
		*s.value = append(*s.value, val)
	}
	return nil
}

func (s *stringArrayValue) Type() string {
	return "stringArray"
}

func (s *stringArrayValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}

func stringArrayConv(sval string) (interface{}, error) {
	sval = sval[1 : len(sval)-1]
	
	if len(sval) == 0 {
		return []string{}, nil
	}
	return readAsCSV(sval)
}


func (f *FlagSet) GetStringArray(name string) ([]string, error) {
	val, err := f.getFlagType(name, "stringArray", stringArrayConv)
	if err != nil {
		return []string{}, err
	}
	return val.([]string), nil
}




func (f *FlagSet) StringArrayVar(p *[]string, name string, value []string, usage string) {
	f.VarP(newStringArrayValue(value, p), name, "", usage)
}


func (f *FlagSet) StringArrayVarP(p *[]string, name, shorthand string, value []string, usage string) {
	f.VarP(newStringArrayValue(value, p), name, shorthand, usage)
}




func StringArrayVar(p *[]string, name string, value []string, usage string) {
	CommandLine.VarP(newStringArrayValue(value, p), name, "", usage)
}


func StringArrayVarP(p *[]string, name, shorthand string, value []string, usage string) {
	CommandLine.VarP(newStringArrayValue(value, p), name, shorthand, usage)
}




func (f *FlagSet) StringArray(name string, value []string, usage string) *[]string {
	p := []string{}
	f.StringArrayVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) StringArrayP(name, shorthand string, value []string, usage string) *[]string {
	p := []string{}
	f.StringArrayVarP(&p, name, shorthand, value, usage)
	return &p
}




func StringArray(name string, value []string, usage string) *[]string {
	return CommandLine.StringArrayP(name, "", value, usage)
}


func StringArrayP(name, shorthand string, value []string, usage string) *[]string {
	return CommandLine.StringArrayP(name, shorthand, value, usage)
}
