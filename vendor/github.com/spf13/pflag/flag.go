




package pflag

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)


var ErrHelp = errors.New("pflag: help requested")


type ErrorHandling int

const (
	
	ContinueOnError ErrorHandling = iota
	
	ExitOnError
	
	PanicOnError
)


type ParseErrorsWhitelist struct {
	
	UnknownFlags bool
}



type NormalizedName string


type FlagSet struct {
	
	
	
	Usage func()

	
	
	SortFlags bool

	
	ParseErrorsWhitelist ParseErrorsWhitelist

	name              string
	parsed            bool
	actual            map[NormalizedName]*Flag
	orderedActual     []*Flag
	sortedActual      []*Flag
	formal            map[NormalizedName]*Flag
	orderedFormal     []*Flag
	sortedFormal      []*Flag
	shorthands        map[byte]*Flag
	args              []string 
	argsLenAtDash     int      
	errorHandling     ErrorHandling
	output            io.Writer 
	interspersed      bool      
	normalizeNameFunc func(f *FlagSet, name string) NormalizedName

	addedGoFlagSets []*goflag.FlagSet
}


type Flag struct {
	Name                string              
	Shorthand           string              
	Usage               string              
	Value               Value               
	DefValue            string              
	Changed             bool                
	NoOptDefVal         string              
	Deprecated          string              
	Hidden              bool                
	ShorthandDeprecated string              
	Annotations         map[string][]string 
}



type Value interface {
	String() string
	Set(string) error
	Type() string
}


func sortFlags(flags map[NormalizedName]*Flag) []*Flag {
	list := make(sort.StringSlice, len(flags))
	i := 0
	for k := range flags {
		list[i] = string(k)
		i++
	}
	list.Sort()
	result := make([]*Flag, len(list))
	for i, name := range list {
		result[i] = flags[NormalizedName(name)]
	}
	return result
}






func (f *FlagSet) SetNormalizeFunc(n func(f *FlagSet, name string) NormalizedName) {
	f.normalizeNameFunc = n
	f.sortedFormal = f.sortedFormal[:0]
	for fname, flag := range f.formal {
		nname := f.normalizeFlagName(flag.Name)
		if fname == nname {
			continue
		}
		flag.Name = string(nname)
		delete(f.formal, fname)
		f.formal[nname] = flag
		if _, set := f.actual[fname]; set {
			delete(f.actual, fname)
			f.actual[nname] = flag
		}
	}
}



func (f *FlagSet) GetNormalizeFunc() func(f *FlagSet, name string) NormalizedName {
	if f.normalizeNameFunc != nil {
		return f.normalizeNameFunc
	}
	return func(f *FlagSet, name string) NormalizedName { return NormalizedName(name) }
}

func (f *FlagSet) normalizeFlagName(name string) NormalizedName {
	n := f.GetNormalizeFunc()
	return n(f, name)
}

func (f *FlagSet) out() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}



func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}




func (f *FlagSet) VisitAll(fn func(*Flag)) {
	if len(f.formal) == 0 {
		return
	}

	var flags []*Flag
	if f.SortFlags {
		if len(f.formal) != len(f.sortedFormal) {
			f.sortedFormal = sortFlags(f.formal)
		}
		flags = f.sortedFormal
	} else {
		flags = f.orderedFormal
	}

	for _, flag := range flags {
		fn(flag)
	}
}


func (f *FlagSet) HasFlags() bool {
	return len(f.formal) > 0
}



func (f *FlagSet) HasAvailableFlags() bool {
	for _, flag := range f.formal {
		if !flag.Hidden {
			return true
		}
	}
	return false
}




func VisitAll(fn func(*Flag)) {
	CommandLine.VisitAll(fn)
}




func (f *FlagSet) Visit(fn func(*Flag)) {
	if len(f.actual) == 0 {
		return
	}

	var flags []*Flag
	if f.SortFlags {
		if len(f.actual) != len(f.sortedActual) {
			f.sortedActual = sortFlags(f.actual)
		}
		flags = f.sortedActual
	} else {
		flags = f.orderedActual
	}

	for _, flag := range flags {
		fn(flag)
	}
}




func Visit(fn func(*Flag)) {
	CommandLine.Visit(fn)
}


func (f *FlagSet) Lookup(name string) *Flag {
	return f.lookup(f.normalizeFlagName(name))
}




func (f *FlagSet) ShorthandLookup(name string) *Flag {
	if name == "" {
		return nil
	}
	if len(name) > 1 {
		msg := fmt.Sprintf("can not look up shorthand which is more than one ASCII character: %q", name)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	}
	c := name[0]
	return f.shorthands[c]
}


func (f *FlagSet) lookup(name NormalizedName) *Flag {
	return f.formal[name]
}


func (f *FlagSet) getFlagType(name string, ftype string, convFunc func(sval string) (interface{}, error)) (interface{}, error) {
	flag := f.Lookup(name)
	if flag == nil {
		err := fmt.Errorf("flag accessed but not defined: %s", name)
		return nil, err
	}

	if flag.Value.Type() != ftype {
		err := fmt.Errorf("trying to get %s value of flag of type %s", ftype, flag.Value.Type())
		return nil, err
	}

	sval := flag.Value.String()
	result, err := convFunc(sval)
	if err != nil {
		return nil, err
	}
	return result, nil
}




func (f *FlagSet) ArgsLenAtDash() int {
	return f.argsLenAtDash
}




func (f *FlagSet) MarkDeprecated(name string, usageMessage string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return fmt.Errorf("flag %q does not exist", name)
	}
	if usageMessage == "" {
		return fmt.Errorf("deprecated message for flag %q must be set", name)
	}
	flag.Deprecated = usageMessage
	flag.Hidden = true
	return nil
}




func (f *FlagSet) MarkShorthandDeprecated(name string, usageMessage string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return fmt.Errorf("flag %q does not exist", name)
	}
	if usageMessage == "" {
		return fmt.Errorf("deprecated message for flag %q must be set", name)
	}
	flag.ShorthandDeprecated = usageMessage
	return nil
}



func (f *FlagSet) MarkHidden(name string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return fmt.Errorf("flag %q does not exist", name)
	}
	flag.Hidden = true
	return nil
}



func Lookup(name string) *Flag {
	return CommandLine.Lookup(name)
}



func ShorthandLookup(name string) *Flag {
	return CommandLine.ShorthandLookup(name)
}


func (f *FlagSet) Set(name, value string) error {
	normalName := f.normalizeFlagName(name)
	flag, ok := f.formal[normalName]
	if !ok {
		return fmt.Errorf("no such flag -%v", name)
	}

	err := flag.Value.Set(value)
	if err != nil {
		var flagName string
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
			flagName = fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
		} else {
			flagName = fmt.Sprintf("--%s", flag.Name)
		}
		return fmt.Errorf("invalid argument %q for %q flag: %v", value, flagName, err)
	}

	if !flag.Changed {
		if f.actual == nil {
			f.actual = make(map[NormalizedName]*Flag)
		}
		f.actual[normalName] = flag
		f.orderedActual = append(f.orderedActual, flag)

		flag.Changed = true
	}

	if flag.Deprecated != "" {
		fmt.Fprintf(f.out(), "Flag --%s has been deprecated, %s\n", flag.Name, flag.Deprecated)
	}
	return nil
}




func (f *FlagSet) SetAnnotation(name, key string, values []string) error {
	normalName := f.normalizeFlagName(name)
	flag, ok := f.formal[normalName]
	if !ok {
		return fmt.Errorf("no such flag -%v", name)
	}
	if flag.Annotations == nil {
		flag.Annotations = map[string][]string{}
	}
	flag.Annotations[key] = values
	return nil
}



func (f *FlagSet) Changed(name string) bool {
	flag := f.Lookup(name)
	
	if flag == nil {
		return false
	}
	return flag.Changed
}


func Set(name, value string) error {
	return CommandLine.Set(name, value)
}



func (f *FlagSet) PrintDefaults() {
	usages := f.FlagUsages()
	fmt.Fprint(f.out(), usages)
}



func (f *Flag) defaultIsZeroValue() bool {
	switch f.Value.(type) {
	case boolFlag:
		return f.DefValue == "false"
	case *durationValue:
		
		return f.DefValue == "0" || f.DefValue == "0s"
	case *intValue, *int8Value, *int32Value, *int64Value, *uintValue, *uint8Value, *uint16Value, *uint32Value, *uint64Value, *countValue, *float32Value, *float64Value:
		return f.DefValue == "0"
	case *stringValue:
		return f.DefValue == ""
	case *ipValue, *ipMaskValue, *ipNetValue:
		return f.DefValue == "<nil>"
	case *intSliceValue, *stringSliceValue, *stringArrayValue:
		return f.DefValue == "[]"
	default:
		switch f.Value.String() {
		case "false":
			return true
		case "<nil>":
			return true
		case "":
			return true
		case "0":
			return true
		}
		return false
	}
}






func UnquoteUsage(flag *Flag) (name string, usage string) {
	
	usage = flag.Usage
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break 
		}
	}

	name = flag.Value.Type()
	switch name {
	case "bool":
		name = ""
	case "float64":
		name = "float"
	case "int64":
		name = "int"
	case "uint64":
		name = "uint"
	case "stringSlice":
		name = "strings"
	case "intSlice":
		name = "ints"
	case "uintSlice":
		name = "uints"
	case "boolSlice":
		name = "bools"
	}

	return
}





func wrapN(i, slop int, s string) (string, string) {
	if i+slop > len(s) {
		return s, ""
	}

	w := strings.LastIndexAny(s[:i], " \t\n")
	if w <= 0 {
		return s, ""
	}
	nlPos := strings.LastIndex(s[:i], "\n")
	if nlPos > 0 && nlPos < w {
		return s[:nlPos], s[nlPos+1:]
	}
	return s[:w], s[w+1:]
}




func wrap(i, w int, s string) string {
	if w == 0 {
		return strings.Replace(s, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	
	
	wrap := w - i

	var r, l string

	
	
	if wrap < 24 {
		i = 16
		wrap = w - i
		r += "\n" + strings.Repeat(" ", i)
	}
	
	if wrap < 24 {
		return strings.Replace(s, "\n", r, -1)
	}

	
	
	
	slop := 5
	wrap = wrap - slop

	
	
	l, s = wrapN(wrap, slop, s)
	r = r + strings.Replace(l, "\n", "\n"+strings.Repeat(" ", i), -1)

	
	for s != "" {
		var t string

		t, s = wrapN(wrap, slop, s)
		r = r + "\n" + strings.Repeat(" ", i) + strings.Replace(t, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	return r

}




func (f *FlagSet) FlagUsagesWrapped(cols int) string {
	buf := new(bytes.Buffer)

	lines := make([]string, 0, len(f.formal))

	maxlen := 0
	f.VisitAll(func(flag *Flag) {
		if flag.Hidden {
			return
		}

		line := ""
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
			line = fmt.Sprintf("  -%s, --%s", flag.Shorthand, flag.Name)
		} else {
			line = fmt.Sprintf("      --%s", flag.Name)
		}

		varname, usage := UnquoteUsage(flag)
		if varname != "" {
			line += " " + varname
		}
		if flag.NoOptDefVal != "" {
			switch flag.Value.Type() {
			case "string":
				line += fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
			case "bool":
				if flag.NoOptDefVal != "true" {
					line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
				}
			case "count":
				if flag.NoOptDefVal != "+1" {
					line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
				}
			default:
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		}

		
		
		line += "\x00"
		if len(line) > maxlen {
			maxlen = len(line)
		}

		line += usage
		if !flag.defaultIsZeroValue() {
			if flag.Value.Type() == "string" {
				line += fmt.Sprintf(" (default %q)", flag.DefValue)
			} else {
				line += fmt.Sprintf(" (default %s)", flag.DefValue)
			}
		}
		if len(flag.Deprecated) != 0 {
			line += fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
		}

		lines = append(lines, line)
	})

	for _, line := range lines {
		sidx := strings.Index(line, "\x00")
		spacing := strings.Repeat(" ", maxlen-sidx)
		
		fmt.Fprintln(buf, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
	}

	return buf.String()
}



func (f *FlagSet) FlagUsages() string {
	return f.FlagUsagesWrapped(0)
}


func PrintDefaults() {
	CommandLine.PrintDefaults()
}


func defaultUsage(f *FlagSet) {
	fmt.Fprintf(f.out(), "Usage of %s:\n", f.name)
	f.PrintDefaults()
}









var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}


func (f *FlagSet) NFlag() int { return len(f.actual) }


func NFlag() int { return len(CommandLine.actual) }



func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}
	return f.args[i]
}



func Arg(i int) string {
	return CommandLine.Arg(i)
}


func (f *FlagSet) NArg() int { return len(f.args) }


func NArg() int { return len(CommandLine.args) }


func (f *FlagSet) Args() []string { return f.args }


func Args() []string { return CommandLine.args }







func (f *FlagSet) Var(value Value, name string, usage string) {
	f.VarP(value, name, "", usage)
}


func (f *FlagSet) VarPF(value Value, name, shorthand, usage string) *Flag {
	
	flag := &Flag{
		Name:      name,
		Shorthand: shorthand,
		Usage:     usage,
		Value:     value,
		DefValue:  value.String(),
	}
	f.AddFlag(flag)
	return flag
}


func (f *FlagSet) VarP(value Value, name, shorthand, usage string) {
	f.VarPF(value, name, shorthand, usage)
}


func (f *FlagSet) AddFlag(flag *Flag) {
	normalizedFlagName := f.normalizeFlagName(flag.Name)

	_, alreadyThere := f.formal[normalizedFlagName]
	if alreadyThere {
		msg := fmt.Sprintf("%s flag redefined: %s", f.name, flag.Name)
		fmt.Fprintln(f.out(), msg)
		panic(msg) 
	}
	if f.formal == nil {
		f.formal = make(map[NormalizedName]*Flag)
	}

	flag.Name = string(normalizedFlagName)
	f.formal[normalizedFlagName] = flag
	f.orderedFormal = append(f.orderedFormal, flag)

	if flag.Shorthand == "" {
		return
	}
	if len(flag.Shorthand) > 1 {
		msg := fmt.Sprintf("%q shorthand is more than one ASCII character", flag.Shorthand)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	}
	if f.shorthands == nil {
		f.shorthands = make(map[byte]*Flag)
	}
	c := flag.Shorthand[0]
	used, alreadyThere := f.shorthands[c]
	if alreadyThere {
		msg := fmt.Sprintf("unable to redefine %q shorthand in %q flagset: it's already used for %q flag", c, f.name, used.Name)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	}
	f.shorthands[c] = flag
}



func (f *FlagSet) AddFlagSet(newSet *FlagSet) {
	if newSet == nil {
		return
	}
	newSet.VisitAll(func(flag *Flag) {
		if f.Lookup(flag.Name) == nil {
			f.AddFlag(flag)
		}
	})
}







func Var(value Value, name string, usage string) {
	CommandLine.VarP(value, name, "", usage)
}


func VarP(value Value, name, shorthand, usage string) {
	CommandLine.VarP(value, name, shorthand, usage)
}



func (f *FlagSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	if f.errorHandling != ContinueOnError {
		fmt.Fprintln(f.out(), err)
		f.usage()
	}
	return err
}



func (f *FlagSet) usage() {
	if f == CommandLine {
		Usage()
	} else if f.Usage == nil {
		defaultUsage(f)
	} else {
		f.Usage()
	}
}




func stripUnknownFlagValue(args []string) []string {
	if len(args) == 0 {
		
		return args
	}

	first := args[0]
	if len(first) > 0 && first[0] == '-' {
		
		return args
	}

	
	if len(args) > 1 {
		return args[1:]
	}
	return nil
}

func (f *FlagSet) parseLongArg(s string, args []string, fn parseFunc) (a []string, err error) {
	a = args
	name := s[2:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		err = f.failf("bad flag syntax: %s", s)
		return
	}

	split := strings.SplitN(name, "=", 2)
	name = split[0]
	flag, exists := f.formal[f.normalizeFlagName(name)]

	if !exists {
		switch {
		case name == "help":
			f.usage()
			return a, ErrHelp
		case f.ParseErrorsWhitelist.UnknownFlags:
			
			
			if len(split) >= 2 {
				return a, nil
			}

			return stripUnknownFlagValue(a), nil
		default:
			err = f.failf("unknown flag: --%s", name)
			return
		}
	}

	var value string
	if len(split) == 2 {
		
		value = split[1]
	} else if flag.NoOptDefVal != "" {
		
		value = flag.NoOptDefVal
	} else if len(a) > 0 {
		
		value = a[0]
		a = a[1:]
	} else {
		
		err = f.failf("flag needs an argument: %s", s)
		return
	}

	err = fn(flag, value)
	if err != nil {
		f.failf(err.Error())
	}
	return
}

func (f *FlagSet) parseSingleShortArg(shorthands string, args []string, fn parseFunc) (outShorts string, outArgs []string, err error) {
	outArgs = args

	if strings.HasPrefix(shorthands, "test.") {
		return
	}

	outShorts = shorthands[1:]
	c := shorthands[0]

	flag, exists := f.shorthands[c]
	if !exists {
		switch {
		case c == 'h':
			f.usage()
			err = ErrHelp
			return
		case f.ParseErrorsWhitelist.UnknownFlags:
			
			
			if len(shorthands) > 2 && shorthands[1] == '=' {
				outShorts = ""
				return
			}

			outArgs = stripUnknownFlagValue(outArgs)
			return
		default:
			err = f.failf("unknown shorthand flag: %q in -%s", c, shorthands)
			return
		}
	}

	var value string
	if len(shorthands) > 2 && shorthands[1] == '=' {
		
		value = shorthands[2:]
		outShorts = ""
	} else if flag.NoOptDefVal != "" {
		
		value = flag.NoOptDefVal
	} else if len(shorthands) > 1 {
		
		value = shorthands[1:]
		outShorts = ""
	} else if len(args) > 0 {
		
		value = args[0]
		outArgs = args[1:]
	} else {
		
		err = f.failf("flag needs an argument: %q in -%s", c, shorthands)
		return
	}

	if flag.ShorthandDeprecated != "" {
		fmt.Fprintf(f.out(), "Flag shorthand -%s has been deprecated, %s\n", flag.Shorthand, flag.ShorthandDeprecated)
	}

	err = fn(flag, value)
	if err != nil {
		f.failf(err.Error())
	}
	return
}

func (f *FlagSet) parseShortArg(s string, args []string, fn parseFunc) (a []string, err error) {
	a = args
	shorthands := s[1:]

	
	for len(shorthands) > 0 {
		shorthands, a, err = f.parseSingleShortArg(shorthands, args, fn)
		if err != nil {
			return
		}
	}

	return
}

func (f *FlagSet) parseArgs(args []string, fn parseFunc) (err error) {
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			if !f.interspersed {
				f.args = append(f.args, s)
				f.args = append(f.args, args...)
				return nil
			}
			f.args = append(f.args, s)
			continue
		}

		if s[1] == '-' {
			if len(s) == 2 { 
				f.argsLenAtDash = len(f.args)
				f.args = append(f.args, args...)
				break
			}
			args, err = f.parseLongArg(s, args, fn)
		} else {
			args, err = f.parseShortArg(s, args, fn)
		}
		if err != nil {
			return
		}
	}
	return
}





func (f *FlagSet) Parse(arguments []string) error {
	if f.addedGoFlagSets != nil {
		for _, goFlagSet := range f.addedGoFlagSets {
			goFlagSet.Parse(nil)
		}
	}
	f.parsed = true

	if len(arguments) < 0 {
		return nil
	}

	f.args = make([]string, 0, len(arguments))

	set := func(flag *Flag, value string) error {
		return f.Set(flag.Name, value)
	}

	err := f.parseArgs(arguments, set)
	if err != nil {
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			fmt.Println(err)
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

type parseFunc func(flag *Flag, value string) error






func (f *FlagSet) ParseAll(arguments []string, fn func(flag *Flag, value string) error) error {
	f.parsed = true
	f.args = make([]string, 0, len(arguments))

	err := f.parseArgs(arguments, fn)
	if err != nil {
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}


func (f *FlagSet) Parsed() bool {
	return f.parsed
}



func Parse() {
	
	CommandLine.Parse(os.Args[1:])
}




func ParseAll(fn func(flag *Flag, value string) error) {
	
	CommandLine.ParseAll(os.Args[1:], fn)
}


func SetInterspersed(interspersed bool) {
	CommandLine.SetInterspersed(interspersed)
}


func Parsed() bool {
	return CommandLine.Parsed()
}


var CommandLine = NewFlagSet(os.Args[0], ExitOnError)



func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
		argsLenAtDash: -1,
		interspersed:  true,
		SortFlags:     true,
	}
	return f
}


func (f *FlagSet) SetInterspersed(interspersed bool) {
	f.interspersed = interspersed
}




func (f *FlagSet) Init(name string, errorHandling ErrorHandling) {
	f.name = name
	f.errorHandling = errorHandling
	f.argsLenAtDash = -1
}
