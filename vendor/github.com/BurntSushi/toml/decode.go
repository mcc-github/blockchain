package toml

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"strings"
	"time"
)

func e(format string, args ...interface{}) error {
	return fmt.Errorf("toml: "+format, args...)
}



type Unmarshaler interface {
	UnmarshalTOML(interface{}) error
}


func Unmarshal(p []byte, v interface{}) error {
	_, err := Decode(string(p), v)
	return err
}













type Primitive struct {
	undecoded interface{}
	context   Key
}




func PrimitiveDecode(primValue Primitive, v interface{}) error {
	md := MetaData{decoded: make(map[string]bool)}
	return md.unify(primValue.undecoded, rvalue(v))
}












func (md *MetaData) PrimitiveDecode(primValue Primitive, v interface{}) error {
	md.context = primValue.context
	defer func() { md.context = nil }()
	return md.unify(primValue.undecoded, rvalue(v))
}





































func Decode(data string, v interface{}) (MetaData, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return MetaData{}, e("Decode of non-pointer %s", reflect.TypeOf(v))
	}
	if rv.IsNil() {
		return MetaData{}, e("Decode of nil %s", reflect.TypeOf(v))
	}
	p, err := parse(data)
	if err != nil {
		return MetaData{}, err
	}
	md := MetaData{
		p.mapping, p.types, p.ordered,
		make(map[string]bool, len(p.ordered)), nil,
	}
	return md, md.unify(p.mapping, indirect(rv))
}



func DecodeFile(fpath string, v interface{}) (MetaData, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return MetaData{}, err
	}
	return Decode(string(bs), v)
}



func DecodeReader(r io.Reader, v interface{}) (MetaData, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return MetaData{}, err
	}
	return Decode(string(bs), v)
}






func (md *MetaData) unify(data interface{}, rv reflect.Value) error {

	
	if rv.Type() == reflect.TypeOf((*Primitive)(nil)).Elem() {
		
		
		context := make(Key, len(md.context))
		copy(context, md.context)
		rv.Set(reflect.ValueOf(Primitive{
			undecoded: data,
			context:   context,
		}))
		return nil
	}

	
	if rv.CanAddr() {
		if v, ok := rv.Addr().Interface().(Unmarshaler); ok {
			return v.UnmarshalTOML(data)
		}
	}

	
	
	
	
	if rv.Type().AssignableTo(rvalue(time.Time{}).Type()) {
		return md.unifyDatetime(data, rv)
	}

	
	if v, ok := rv.Interface().(TextUnmarshaler); ok {
		return md.unifyText(data, v)
	}
	
	
	
	
	
	
	

	k := rv.Kind()

	
	if k >= reflect.Int && k <= reflect.Uint64 {
		return md.unifyInt(data, rv)
	}
	switch k {
	case reflect.Ptr:
		elem := reflect.New(rv.Type().Elem())
		err := md.unify(data, reflect.Indirect(elem))
		if err != nil {
			return err
		}
		rv.Set(elem)
		return nil
	case reflect.Struct:
		return md.unifyStruct(data, rv)
	case reflect.Map:
		return md.unifyMap(data, rv)
	case reflect.Array:
		return md.unifyArray(data, rv)
	case reflect.Slice:
		return md.unifySlice(data, rv)
	case reflect.String:
		return md.unifyString(data, rv)
	case reflect.Bool:
		return md.unifyBool(data, rv)
	case reflect.Interface:
		
		if rv.NumMethod() > 0 {
			return e("unsupported type %s", rv.Type())
		}
		return md.unifyAnything(data, rv)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return md.unifyFloat64(data, rv)
	}
	return e("unsupported type %s", rv.Kind())
}

func (md *MetaData) unifyStruct(mapping interface{}, rv reflect.Value) error {
	tmap, ok := mapping.(map[string]interface{})
	if !ok {
		if mapping == nil {
			return nil
		}
		return e("type mismatch for %s: expected table but found %T",
			rv.Type().String(), mapping)
	}

	for key, datum := range tmap {
		var f *field
		fields := cachedTypeFields(rv.Type())
		for i := range fields {
			ff := &fields[i]
			if ff.name == key {
				f = ff
				break
			}
			if f == nil && strings.EqualFold(ff.name, key) {
				f = ff
			}
		}
		if f != nil {
			subv := rv
			for _, i := range f.index {
				subv = indirect(subv.Field(i))
			}
			if isUnifiable(subv) {
				md.decoded[md.context.add(key).String()] = true
				md.context = append(md.context, key)
				if err := md.unify(datum, subv); err != nil {
					return err
				}
				md.context = md.context[0 : len(md.context)-1]
			} else if f.name != "" {
				
				return e("cannot write unexported field %s.%s",
					rv.Type().String(), f.name)
			}
		}
	}
	return nil
}

func (md *MetaData) unifyMap(mapping interface{}, rv reflect.Value) error {
	tmap, ok := mapping.(map[string]interface{})
	if !ok {
		if tmap == nil {
			return nil
		}
		return badtype("map", mapping)
	}
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rv.Type()))
	}
	for k, v := range tmap {
		md.decoded[md.context.add(k).String()] = true
		md.context = append(md.context, k)

		rvkey := indirect(reflect.New(rv.Type().Key()))
		rvval := reflect.Indirect(reflect.New(rv.Type().Elem()))
		if err := md.unify(v, rvval); err != nil {
			return err
		}
		md.context = md.context[0 : len(md.context)-1]

		rvkey.SetString(k)
		rv.SetMapIndex(rvkey, rvval)
	}
	return nil
}

func (md *MetaData) unifyArray(data interface{}, rv reflect.Value) error {
	datav := reflect.ValueOf(data)
	if datav.Kind() != reflect.Slice {
		if !datav.IsValid() {
			return nil
		}
		return badtype("slice", data)
	}
	sliceLen := datav.Len()
	if sliceLen != rv.Len() {
		return e("expected array length %d; got TOML array of length %d",
			rv.Len(), sliceLen)
	}
	return md.unifySliceArray(datav, rv)
}

func (md *MetaData) unifySlice(data interface{}, rv reflect.Value) error {
	datav := reflect.ValueOf(data)
	if datav.Kind() != reflect.Slice {
		if !datav.IsValid() {
			return nil
		}
		return badtype("slice", data)
	}
	n := datav.Len()
	if rv.IsNil() || rv.Cap() < n {
		rv.Set(reflect.MakeSlice(rv.Type(), n, n))
	}
	rv.SetLen(n)
	return md.unifySliceArray(datav, rv)
}

func (md *MetaData) unifySliceArray(data, rv reflect.Value) error {
	sliceLen := data.Len()
	for i := 0; i < sliceLen; i++ {
		v := data.Index(i).Interface()
		sliceval := indirect(rv.Index(i))
		if err := md.unify(v, sliceval); err != nil {
			return err
		}
	}
	return nil
}

func (md *MetaData) unifyDatetime(data interface{}, rv reflect.Value) error {
	if _, ok := data.(time.Time); ok {
		rv.Set(reflect.ValueOf(data))
		return nil
	}
	return badtype("time.Time", data)
}

func (md *MetaData) unifyString(data interface{}, rv reflect.Value) error {
	if s, ok := data.(string); ok {
		rv.SetString(s)
		return nil
	}
	return badtype("string", data)
}

func (md *MetaData) unifyFloat64(data interface{}, rv reflect.Value) error {
	if num, ok := data.(float64); ok {
		switch rv.Kind() {
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			rv.SetFloat(num)
		default:
			panic("bug")
		}
		return nil
	}
	return badtype("float", data)
}

func (md *MetaData) unifyInt(data interface{}, rv reflect.Value) error {
	if num, ok := data.(int64); ok {
		if rv.Kind() >= reflect.Int && rv.Kind() <= reflect.Int64 {
			switch rv.Kind() {
			case reflect.Int, reflect.Int64:
				
			case reflect.Int8:
				if num < math.MinInt8 || num > math.MaxInt8 {
					return e("value %d is out of range for int8", num)
				}
			case reflect.Int16:
				if num < math.MinInt16 || num > math.MaxInt16 {
					return e("value %d is out of range for int16", num)
				}
			case reflect.Int32:
				if num < math.MinInt32 || num > math.MaxInt32 {
					return e("value %d is out of range for int32", num)
				}
			}
			rv.SetInt(num)
		} else if rv.Kind() >= reflect.Uint && rv.Kind() <= reflect.Uint64 {
			unum := uint64(num)
			switch rv.Kind() {
			case reflect.Uint, reflect.Uint64:
				
			case reflect.Uint8:
				if num < 0 || unum > math.MaxUint8 {
					return e("value %d is out of range for uint8", num)
				}
			case reflect.Uint16:
				if num < 0 || unum > math.MaxUint16 {
					return e("value %d is out of range for uint16", num)
				}
			case reflect.Uint32:
				if num < 0 || unum > math.MaxUint32 {
					return e("value %d is out of range for uint32", num)
				}
			}
			rv.SetUint(unum)
		} else {
			panic("unreachable")
		}
		return nil
	}
	return badtype("integer", data)
}

func (md *MetaData) unifyBool(data interface{}, rv reflect.Value) error {
	if b, ok := data.(bool); ok {
		rv.SetBool(b)
		return nil
	}
	return badtype("boolean", data)
}

func (md *MetaData) unifyAnything(data interface{}, rv reflect.Value) error {
	rv.Set(reflect.ValueOf(data))
	return nil
}

func (md *MetaData) unifyText(data interface{}, v TextUnmarshaler) error {
	var s string
	switch sdata := data.(type) {
	case TextMarshaler:
		text, err := sdata.MarshalText()
		if err != nil {
			return err
		}
		s = string(text)
	case fmt.Stringer:
		s = sdata.String()
	case string:
		s = sdata
	case bool:
		s = fmt.Sprintf("%v", sdata)
	case int64:
		s = fmt.Sprintf("%d", sdata)
	case float64:
		s = fmt.Sprintf("%f", sdata)
	default:
		return badtype("primitive (string-like)", data)
	}
	if err := v.UnmarshalText([]byte(s)); err != nil {
		return err
	}
	return nil
}


func rvalue(v interface{}) reflect.Value {
	return indirect(reflect.ValueOf(v))
}







func indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		if v.CanSet() {
			pv := v.Addr()
			if _, ok := pv.Interface().(TextUnmarshaler); ok {
				return pv
			}
		}
		return v
	}
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return indirect(reflect.Indirect(v))
}

func isUnifiable(rv reflect.Value) bool {
	if rv.CanSet() {
		return true
	}
	if _, ok := rv.Interface().(TextUnmarshaler); ok {
		return true
	}
	return false
}

func badtype(expected string, data interface{}) error {
	return e("cannot load TOML value of type %T into a Go %s", data, expected)
}
