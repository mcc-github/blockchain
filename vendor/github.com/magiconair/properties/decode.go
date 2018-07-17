



package properties

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)
















































































func (p *Properties) Decode(x interface{}) error {
	t, v := reflect.TypeOf(x), reflect.ValueOf(x)
	if t.Kind() != reflect.Ptr || v.Elem().Type().Kind() != reflect.Struct {
		return fmt.Errorf("not a pointer to struct: %s", t)
	}
	if err := dec(p, "", nil, nil, v); err != nil {
		return err
	}
	return nil
}

func dec(p *Properties, key string, def *string, opts map[string]string, v reflect.Value) error {
	t := v.Type()

	
	value := func() (string, error) {
		if val, ok := p.Get(key); ok {
			return val, nil
		}
		if def != nil {
			return *def, nil
		}
		return "", fmt.Errorf("missing required key %s", key)
	}

	
	conv := func(s string, t reflect.Type) (val reflect.Value, err error) {
		var v interface{}

		switch {
		case isDuration(t):
			v, err = time.ParseDuration(s)

		case isTime(t):
			layout := opts["layout"]
			if layout == "" {
				layout = time.RFC3339
			}
			v, err = time.Parse(layout, s)

		case isBool(t):
			v, err = boolVal(s), nil

		case isString(t):
			v, err = s, nil

		case isFloat(t):
			v, err = strconv.ParseFloat(s, 64)

		case isInt(t):
			v, err = strconv.ParseInt(s, 10, 64)

		case isUint(t):
			v, err = strconv.ParseUint(s, 10, 64)

		default:
			return reflect.Zero(t), fmt.Errorf("unsupported type %s", t)
		}
		if err != nil {
			return reflect.Zero(t), err
		}
		return reflect.ValueOf(v).Convert(t), nil
	}

	
	
	keydef := func(f reflect.StructField) (string, *string, map[string]string) {
		_key, _opts := parseTag(f.Tag.Get("properties"))

		var _def *string
		if d, ok := _opts["default"]; ok {
			_def = &d
		}
		if _key != "" {
			return _key, _def, _opts
		}
		return f.Name, _def, _opts
	}

	switch {
	case isDuration(t) || isTime(t) || isBool(t) || isString(t) || isFloat(t) || isInt(t) || isUint(t):
		s, err := value()
		if err != nil {
			return err
		}
		val, err := conv(s, t)
		if err != nil {
			return err
		}
		v.Set(val)

	case isPtr(t):
		return dec(p, key, def, opts, v.Elem())

	case isStruct(t):
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			fk, def, opts := keydef(t.Field(i))
			if !fv.CanSet() {
				return fmt.Errorf("cannot set %s", t.Field(i).Name)
			}
			if fk == "-" {
				continue
			}
			if key != "" {
				fk = key + "." + fk
			}
			if err := dec(p, fk, def, opts, fv); err != nil {
				return err
			}
		}
		return nil

	case isArray(t):
		val, err := value()
		if err != nil {
			return err
		}
		vals := split(val, ";")
		a := reflect.MakeSlice(t, 0, len(vals))
		for _, s := range vals {
			val, err := conv(s, t.Elem())
			if err != nil {
				return err
			}
			a = reflect.Append(a, val)
		}
		v.Set(a)

	case isMap(t):
		valT := t.Elem()
		m := reflect.MakeMap(t)
		for postfix := range p.FilterStripPrefix(key + ".").m {
			pp := strings.SplitN(postfix, ".", 2)
			mk, mv := pp[0], reflect.New(valT)
			if err := dec(p, key+"."+mk, nil, nil, mv); err != nil {
				return err
			}
			m.SetMapIndex(reflect.ValueOf(mk), mv.Elem())
		}
		v.Set(m)

	default:
		return fmt.Errorf("unsupported type %s", t)
	}
	return nil
}



func split(s string, sep string) []string {
	var a []string
	for _, v := range strings.Split(s, sep) {
		if v = strings.TrimSpace(v); v != "" {
			a = append(a, v)
		}
	}
	return a
}


func parseTag(tag string) (key string, opts map[string]string) {
	opts = map[string]string{}
	for i, s := range strings.Split(tag, ",") {
		if i == 0 {
			key = s
			continue
		}

		pp := strings.SplitN(s, "=", 2)
		if len(pp) == 1 {
			opts[pp[0]] = ""
		} else {
			opts[pp[0]] = pp[1]
		}
	}
	return key, opts
}

func isArray(t reflect.Type) bool    { return t.Kind() == reflect.Array || t.Kind() == reflect.Slice }
func isBool(t reflect.Type) bool     { return t.Kind() == reflect.Bool }
func isDuration(t reflect.Type) bool { return t == reflect.TypeOf(time.Second) }
func isMap(t reflect.Type) bool      { return t.Kind() == reflect.Map }
func isPtr(t reflect.Type) bool      { return t.Kind() == reflect.Ptr }
func isString(t reflect.Type) bool   { return t.Kind() == reflect.String }
func isStruct(t reflect.Type) bool   { return t.Kind() == reflect.Struct }
func isTime(t reflect.Type) bool     { return t == reflect.TypeOf(time.Time{}) }
func isFloat(t reflect.Type) bool {
	return t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64
}
func isInt(t reflect.Type) bool {
	return t.Kind() == reflect.Int || t.Kind() == reflect.Int8 || t.Kind() == reflect.Int16 || t.Kind() == reflect.Int32 || t.Kind() == reflect.Int64
}
func isUint(t reflect.Type) bool {
	return t.Kind() == reflect.Uint || t.Kind() == reflect.Uint8 || t.Kind() == reflect.Uint16 || t.Kind() == reflect.Uint32 || t.Kind() == reflect.Uint64
}
