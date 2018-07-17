
package filters 

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/versions"
)


type Args struct {
	fields map[string]map[string]bool
}


type KeyValuePair struct {
	Key   string
	Value string
}


func Arg(key, value string) KeyValuePair {
	return KeyValuePair{Key: key, Value: value}
}


func NewArgs(initialArgs ...KeyValuePair) Args {
	args := Args{fields: map[string]map[string]bool{}}
	for _, arg := range initialArgs {
		args.Add(arg.Key, arg.Value)
	}
	return args
}




func ParseFlag(arg string, prev Args) (Args, error) {
	filters := prev
	if len(arg) == 0 {
		return filters, nil
	}

	if !strings.Contains(arg, "=") {
		return filters, ErrBadFormat
	}

	f := strings.SplitN(arg, "=", 2)

	name := strings.ToLower(strings.TrimSpace(f[0]))
	value := strings.TrimSpace(f[1])

	filters.Add(name, value)

	return filters, nil
}




var ErrBadFormat = errors.New("bad format of filter (expected name=value)")




func ToParam(a Args) (string, error) {
	return ToJSON(a)
}


func (args Args) MarshalJSON() ([]byte, error) {
	if len(args.fields) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(args.fields)
}


func ToJSON(a Args) (string, error) {
	if a.Len() == 0 {
		return "", nil
	}
	buf, err := json.Marshal(a)
	return string(buf), err
}






func ToParamWithVersion(version string, a Args) (string, error) {
	if a.Len() == 0 {
		return "", nil
	}

	if version != "" && versions.LessThan(version, "1.22") {
		buf, err := json.Marshal(convertArgsToSlice(a.fields))
		return string(buf), err
	}

	return ToJSON(a)
}




func FromParam(p string) (Args, error) {
	return FromJSON(p)
}


func FromJSON(p string) (Args, error) {
	args := NewArgs()

	if p == "" {
		return args, nil
	}

	raw := []byte(p)
	err := json.Unmarshal(raw, &args)
	if err == nil {
		return args, nil
	}

	
	deprecated := map[string][]string{}
	if legacyErr := json.Unmarshal(raw, &deprecated); legacyErr != nil {
		return args, err
	}

	args.fields = deprecatedArgs(deprecated)
	return args, nil
}


func (args Args) UnmarshalJSON(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, &args.fields)
}


func (args Args) Get(key string) []string {
	values := args.fields[key]
	if values == nil {
		return make([]string, 0)
	}
	slice := make([]string, 0, len(values))
	for key := range values {
		slice = append(slice, key)
	}
	return slice
}


func (args Args) Add(key, value string) {
	if _, ok := args.fields[key]; ok {
		args.fields[key][value] = true
	} else {
		args.fields[key] = map[string]bool{value: true}
	}
}


func (args Args) Del(key, value string) {
	if _, ok := args.fields[key]; ok {
		delete(args.fields[key], value)
		if len(args.fields[key]) == 0 {
			delete(args.fields, key)
		}
	}
}


func (args Args) Len() int {
	return len(args.fields)
}



func (args Args) MatchKVList(key string, sources map[string]string) bool {
	fieldValues := args.fields[key]

	
	if len(fieldValues) == 0 {
		return true
	}

	if len(sources) == 0 {
		return false
	}

	for value := range fieldValues {
		testKV := strings.SplitN(value, "=", 2)

		v, ok := sources[testKV[0]]
		if !ok {
			return false
		}
		if len(testKV) == 2 && testKV[1] != v {
			return false
		}
	}

	return true
}


func (args Args) Match(field, source string) bool {
	if args.ExactMatch(field, source) {
		return true
	}

	fieldValues := args.fields[field]
	for name2match := range fieldValues {
		match, err := regexp.MatchString(name2match, source)
		if err != nil {
			continue
		}
		if match {
			return true
		}
	}
	return false
}


func (args Args) ExactMatch(key, source string) bool {
	fieldValues, ok := args.fields[key]
	
	if !ok || len(fieldValues) == 0 {
		return true
	}

	
	return fieldValues[source]
}



func (args Args) UniqueExactMatch(key, source string) bool {
	fieldValues := args.fields[key]
	
	if len(fieldValues) == 0 {
		return true
	}
	if len(args.fields[key]) != 1 {
		return false
	}

	
	return fieldValues[source]
}



func (args Args) FuzzyMatch(key, source string) bool {
	if args.ExactMatch(key, source) {
		return true
	}

	fieldValues := args.fields[key]
	for prefix := range fieldValues {
		if strings.HasPrefix(source, prefix) {
			return true
		}
	}
	return false
}




func (args Args) Include(field string) bool {
	_, ok := args.fields[field]
	return ok
}


func (args Args) Contains(field string) bool {
	_, ok := args.fields[field]
	return ok
}

type invalidFilter string

func (e invalidFilter) Error() string {
	return "Invalid filter '" + string(e) + "'"
}

func (invalidFilter) InvalidParameter() {}



func (args Args) Validate(accepted map[string]bool) error {
	for name := range args.fields {
		if !accepted[name] {
			return invalidFilter(name)
		}
	}
	return nil
}




func (args Args) WalkValues(field string, op func(value string) error) error {
	if _, ok := args.fields[field]; !ok {
		return nil
	}
	for v := range args.fields[field] {
		if err := op(v); err != nil {
			return err
		}
	}
	return nil
}

func deprecatedArgs(d map[string][]string) map[string]map[string]bool {
	m := map[string]map[string]bool{}
	for k, v := range d {
		values := map[string]bool{}
		for _, vv := range v {
			values[vv] = true
		}
		m[k] = values
	}
	return m
}

func convertArgsToSlice(f map[string]map[string]bool) map[string][]string {
	m := map[string][]string{}
	for k, v := range f {
		values := []string{}
		for kk := range v {
			if v[kk] {
				values = append(values, kk)
			}
		}
		m[k] = values
	}
	return m
}
