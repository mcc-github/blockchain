package opts 

import (
	"fmt"

	"github.com/docker/go-units"
)


type UlimitOpt struct {
	values *map[string]*units.Ulimit
}


func NewUlimitOpt(ref *map[string]*units.Ulimit) *UlimitOpt {
	if ref == nil {
		ref = &map[string]*units.Ulimit{}
	}
	return &UlimitOpt{ref}
}


func (o *UlimitOpt) Set(val string) error {
	l, err := units.ParseUlimit(val)
	if err != nil {
		return err
	}

	(*o.values)[l.Name] = l

	return nil
}


func (o *UlimitOpt) String() string {
	var out []string
	for _, v := range *o.values {
		out = append(out, v.String())
	}

	return fmt.Sprintf("%v", out)
}


func (o *UlimitOpt) GetList() []*units.Ulimit {
	var ulimits []*units.Ulimit
	for _, v := range *o.values {
		ulimits = append(ulimits, v)
	}

	return ulimits
}


func (o *UlimitOpt) Type() string {
	return "ulimit"
}


type NamedUlimitOpt struct {
	name string
	UlimitOpt
}

var _ NamedOption = &NamedUlimitOpt{}


func NewNamedUlimitOpt(name string, ref *map[string]*units.Ulimit) *NamedUlimitOpt {
	if ref == nil {
		ref = &map[string]*units.Ulimit{}
	}
	return &NamedUlimitOpt{
		name:      name,
		UlimitOpt: *NewUlimitOpt(ref),
	}
}


func (o *NamedUlimitOpt) Name() string {
	return o.name
}
