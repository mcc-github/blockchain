package pflag

import (
	"time"
)


type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) Type() string {
	return "duration"
}

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

func durationConv(sval string) (interface{}, error) {
	return time.ParseDuration(sval)
}


func (f *FlagSet) GetDuration(name string) (time.Duration, error) {
	val, err := f.getFlagType(name, "duration", durationConv)
	if err != nil {
		return 0, err
	}
	return val.(time.Duration), nil
}



func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	f.VarP(newDurationValue(value, p), name, "", usage)
}


func (f *FlagSet) DurationVarP(p *time.Duration, name, shorthand string, value time.Duration, usage string) {
	f.VarP(newDurationValue(value, p), name, shorthand, usage)
}



func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	CommandLine.VarP(newDurationValue(value, p), name, "", usage)
}


func DurationVarP(p *time.Duration, name, shorthand string, value time.Duration, usage string) {
	CommandLine.VarP(newDurationValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	f.DurationVarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) DurationP(name, shorthand string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	f.DurationVarP(p, name, shorthand, value, usage)
	return p
}



func Duration(name string, value time.Duration, usage string) *time.Duration {
	return CommandLine.DurationP(name, "", value, usage)
}


func DurationP(name, shorthand string, value time.Duration, usage string) *time.Duration {
	return CommandLine.DurationP(name, shorthand, value, usage)
}
