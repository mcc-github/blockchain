



package pflag

import (
	goflag "flag"
	"reflect"
	"strings"
)





type flagValueWrapper struct {
	inner    goflag.Value
	flagType string
}



type goBoolFlag interface {
	goflag.Value
	IsBoolFlag() bool
}

func wrapFlagValue(v goflag.Value) Value {
	
	if pv, ok := v.(Value); ok {
		return pv
	}

	pv := &flagValueWrapper{
		inner: v,
	}

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Interface || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	pv.flagType = strings.TrimSuffix(t.Name(), "Value")
	return pv
}

func (v *flagValueWrapper) String() string {
	return v.inner.String()
}

func (v *flagValueWrapper) Set(s string) error {
	return v.inner.Set(s)
}

func (v *flagValueWrapper) Type() string {
	return v.flagType
}





func PFlagFromGoFlag(goflag *goflag.Flag) *Flag {
	
	flag := &Flag{
		Name:  goflag.Name,
		Usage: goflag.Usage,
		Value: wrapFlagValue(goflag.Value),
		
		
		DefValue: goflag.Value.String(),
	}
	
	if len(flag.Name) == 1 {
		flag.Shorthand = flag.Name
	}
	if fv, ok := goflag.Value.(goBoolFlag); ok && fv.IsBoolFlag() {
		flag.NoOptDefVal = "true"
	}
	return flag
}


func (f *FlagSet) AddGoFlag(goflag *goflag.Flag) {
	if f.Lookup(goflag.Name) != nil {
		return
	}
	newflag := PFlagFromGoFlag(goflag)
	f.AddFlag(newflag)
}


func (f *FlagSet) AddGoFlagSet(newSet *goflag.FlagSet) {
	if newSet == nil {
		return
	}
	newSet.VisitAll(func(goflag *goflag.Flag) {
		f.AddGoFlag(goflag)
	})
}
