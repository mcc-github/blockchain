package pflag

import (
	"fmt"
	"strings"
	"time"
)


type durationSliceValue struct {
	value   *[]time.Duration
	changed bool
}

func newDurationSliceValue(val []time.Duration, p *[]time.Duration) *durationSliceValue {
	dsv := new(durationSliceValue)
	dsv.value = p
	*dsv.value = val
	return dsv
}

func (s *durationSliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]time.Duration, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = time.ParseDuration(d)
		if err != nil {
			return err
		}

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *durationSliceValue) Type() string {
	return "durationSlice"
}

func (s *durationSliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%s", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func durationSliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	
	if len(val) == 0 {
		return []time.Duration{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]time.Duration, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = time.ParseDuration(d)
		if err != nil {
			return nil, err
		}

	}
	return out, nil
}


func (f *FlagSet) GetDurationSlice(name string) ([]time.Duration, error) {
	val, err := f.getFlagType(name, "durationSlice", durationSliceConv)
	if err != nil {
		return []time.Duration{}, err
	}
	return val.([]time.Duration), nil
}



func (f *FlagSet) DurationSliceVar(p *[]time.Duration, name string, value []time.Duration, usage string) {
	f.VarP(newDurationSliceValue(value, p), name, "", usage)
}


func (f *FlagSet) DurationSliceVarP(p *[]time.Duration, name, shorthand string, value []time.Duration, usage string) {
	f.VarP(newDurationSliceValue(value, p), name, shorthand, usage)
}



func DurationSliceVar(p *[]time.Duration, name string, value []time.Duration, usage string) {
	CommandLine.VarP(newDurationSliceValue(value, p), name, "", usage)
}


func DurationSliceVarP(p *[]time.Duration, name, shorthand string, value []time.Duration, usage string) {
	CommandLine.VarP(newDurationSliceValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) DurationSlice(name string, value []time.Duration, usage string) *[]time.Duration {
	p := []time.Duration{}
	f.DurationSliceVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) DurationSliceP(name, shorthand string, value []time.Duration, usage string) *[]time.Duration {
	p := []time.Duration{}
	f.DurationSliceVarP(&p, name, shorthand, value, usage)
	return &p
}



func DurationSlice(name string, value []time.Duration, usage string) *[]time.Duration {
	return CommandLine.DurationSliceP(name, "", value, usage)
}


func DurationSliceP(name, shorthand string, value []time.Duration, usage string) *[]time.Duration {
	return CommandLine.DurationSliceP(name, shorthand, value, usage)
}
