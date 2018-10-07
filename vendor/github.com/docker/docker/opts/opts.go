package opts 

import (
	"fmt"
	"net"
	"path"
	"regexp"
	"strings"

	"github.com/docker/go-units"
)

var (
	alphaRegexp  = regexp.MustCompile(`[a-zA-Z]`)
	domainRegexp = regexp.MustCompile(`^(:?(:?[a-zA-Z0-9]|(:?[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]))(:?\.(:?[a-zA-Z0-9]|(:?[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])))*)\.?\s*$`)
)


type ListOpts struct {
	values    *[]string
	validator ValidatorFctType
}


func NewListOpts(validator ValidatorFctType) ListOpts {
	var values []string
	return *NewListOptsRef(&values, validator)
}


func NewListOptsRef(values *[]string, validator ValidatorFctType) *ListOpts {
	return &ListOpts{
		values:    values,
		validator: validator,
	}
}

func (opts *ListOpts) String() string {
	if len(*opts.values) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", *opts.values)
}



func (opts *ListOpts) Set(value string) error {
	if opts.validator != nil {
		v, err := opts.validator(value)
		if err != nil {
			return err
		}
		value = v
	}
	*opts.values = append(*opts.values, value)
	return nil
}


func (opts *ListOpts) Delete(key string) {
	for i, k := range *opts.values {
		if k == key {
			*opts.values = append((*opts.values)[:i], (*opts.values)[i+1:]...)
			return
		}
	}
}



func (opts *ListOpts) GetMap() map[string]struct{} {
	ret := make(map[string]struct{})
	for _, k := range *opts.values {
		ret[k] = struct{}{}
	}
	return ret
}


func (opts *ListOpts) GetAll() []string {
	return *opts.values
}



func (opts *ListOpts) GetAllOrEmpty() []string {
	v := *opts.values
	if v == nil {
		return make([]string, 0)
	}
	return v
}


func (opts *ListOpts) Get(key string) bool {
	for _, k := range *opts.values {
		if k == key {
			return true
		}
	}
	return false
}


func (opts *ListOpts) Len() int {
	return len(*opts.values)
}


func (opts *ListOpts) Type() string {
	return "list"
}


func (opts *ListOpts) WithValidator(validator ValidatorFctType) *ListOpts {
	opts.validator = validator
	return opts
}



type NamedOption interface {
	Name() string
}




type NamedListOpts struct {
	name string
	ListOpts
}

var _ NamedOption = &NamedListOpts{}


func NewNamedListOptsRef(name string, values *[]string, validator ValidatorFctType) *NamedListOpts {
	return &NamedListOpts{
		name:     name,
		ListOpts: *NewListOptsRef(values, validator),
	}
}


func (o *NamedListOpts) Name() string {
	return o.name
}


type MapOpts struct {
	values    map[string]string
	validator ValidatorFctType
}



func (opts *MapOpts) Set(value string) error {
	if opts.validator != nil {
		v, err := opts.validator(value)
		if err != nil {
			return err
		}
		value = v
	}
	vals := strings.SplitN(value, "=", 2)
	if len(vals) == 1 {
		(opts.values)[vals[0]] = ""
	} else {
		(opts.values)[vals[0]] = vals[1]
	}
	return nil
}


func (opts *MapOpts) GetAll() map[string]string {
	return opts.values
}

func (opts *MapOpts) String() string {
	return fmt.Sprintf("%v", opts.values)
}


func (opts *MapOpts) Type() string {
	return "map"
}


func NewMapOpts(values map[string]string, validator ValidatorFctType) *MapOpts {
	if values == nil {
		values = make(map[string]string)
	}
	return &MapOpts{
		values:    values,
		validator: validator,
	}
}




type NamedMapOpts struct {
	name string
	MapOpts
}

var _ NamedOption = &NamedMapOpts{}


func NewNamedMapOpts(name string, values map[string]string, validator ValidatorFctType) *NamedMapOpts {
	return &NamedMapOpts{
		name:    name,
		MapOpts: *NewMapOpts(values, validator),
	}
}


func (o *NamedMapOpts) Name() string {
	return o.name
}


type ValidatorFctType func(val string) (string, error)


type ValidatorFctListType func(val string) ([]string, error)


func ValidateIPAddress(val string) (string, error) {
	var ip = net.ParseIP(strings.TrimSpace(val))
	if ip != nil {
		return ip.String(), nil
	}
	return "", fmt.Errorf("%s is not an ip address", val)
}



func ValidateDNSSearch(val string) (string, error) {
	if val = strings.Trim(val, " "); val == "." {
		return val, nil
	}
	return validateDomain(val)
}

func validateDomain(val string) (string, error) {
	if alphaRegexp.FindString(val) == "" {
		return "", fmt.Errorf("%s is not a valid domain", val)
	}
	ns := domainRegexp.FindSubmatch([]byte(val))
	if len(ns) > 0 && len(ns[1]) < 255 {
		return string(ns[1]), nil
	}
	return "", fmt.Errorf("%s is not a valid domain", val)
}



func ValidateLabel(val string) (string, error) {
	if strings.Count(val, "=") < 1 {
		return "", fmt.Errorf("bad attribute format: %s", val)
	}
	return val, nil
}




func ValidateSingleGenericResource(val string) (string, error) {
	if strings.Count(val, "=") < 1 {
		return "", fmt.Errorf("invalid node-generic-resource format `%s` expected `name=value`", val)
	}
	return val, nil
}


func ParseLink(val string) (string, string, error) {
	if val == "" {
		return "", "", fmt.Errorf("empty string specified for links")
	}
	arr := strings.Split(val, ":")
	if len(arr) > 2 {
		return "", "", fmt.Errorf("bad format for links: %s", val)
	}
	if len(arr) == 1 {
		return val, val, nil
	}
	
	
	
	if strings.HasPrefix(arr[0], "/") {
		_, alias := path.Split(arr[1])
		return arr[0][1:], alias, nil
	}
	return arr[0], arr[1], nil
}


type MemBytes int64


func (m *MemBytes) String() string {
	
	
	
	if m.Value() != 0 {
		return units.BytesSize(float64(m.Value()))
	}
	return "0"
}


func (m *MemBytes) Set(value string) error {
	val, err := units.RAMInBytes(value)
	*m = MemBytes(val)
	return err
}


func (m *MemBytes) Type() string {
	return "bytes"
}


func (m *MemBytes) Value() int64 {
	return int64(*m)
}


func (m *MemBytes) UnmarshalJSON(s []byte) error {
	if len(s) <= 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return fmt.Errorf("invalid size: %q", s)
	}
	val, err := units.RAMInBytes(string(s[1 : len(s)-1]))
	*m = MemBytes(val)
	return err
}
