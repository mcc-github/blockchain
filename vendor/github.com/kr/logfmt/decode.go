



















package logfmt

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)




type Handler interface {
	HandleLogfmt(key, val []byte) error
}




type HandlerFunc func(key, val []byte) error

func (f HandlerFunc) HandleLogfmt(key, val []byte) error {
	return f(key, val)
}







func Unmarshal(data []byte, v interface{}) (err error) {
	h, ok := v.(Handler)
	if !ok {
		h, err = NewStructHandler(v)
		if err != nil {
			return err
		}
	}
	return gotoScanner(data, h)
}


















type StructHandler struct {
	rv reflect.Value
}

func NewStructHandler(v interface{}) (Handler, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	return &StructHandler{rv: rv}, nil
}

func (h *StructHandler) HandleLogfmt(key, val []byte) error {
	el := h.rv.Elem()
	skey := string(key)
	for i := 0; i < el.NumField(); i++ {
		fv := el.Field(i)
		ft := el.Type().Field(i)
		switch {
		case ft.Name == skey:
		case ft.Tag.Get("logfmt") == skey:
		case strings.EqualFold(ft.Name, skey):
		default:
			continue
		}
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				t := fv.Type().Elem()
				v := reflect.New(t)
				fv.Set(v)
				fv = v
			}
			fv = fv.Elem()
		}
		switch fv.Interface().(type) {
		case time.Duration:
			d, err := time.ParseDuration(string(val))
			if err != nil {
				return &UnmarshalTypeError{string(val), fv.Type()}
			}
			fv.Set(reflect.ValueOf(d))
		case string:
			fv.SetString(string(val))
		case []byte:
			b := make([]byte, len(val))
			copy(b, val)
			fv.SetBytes(b)
		case bool:
			fv.SetBool(true)
		default:
			switch {
			case reflect.Int <= fv.Kind() && fv.Kind() <= reflect.Int64:
				v, err := strconv.ParseInt(string(val), 10, 64)
				if err != nil {
					return err
				}
				fv.SetInt(v)
			case reflect.Uint32 <= fv.Kind() && fv.Kind() <= reflect.Uint64:
				v, err := strconv.ParseUint(string(val), 10, 64)
				if err != nil {
					return err
				}
				fv.SetUint(v)
			case reflect.Float32 <= fv.Kind() && fv.Kind() <= reflect.Float64:
				v, err := strconv.ParseFloat(string(val), 10)
				if err != nil {
					return err
				}
				fv.SetFloat(v)
			default:
				return &UnmarshalTypeError{string(val), fv.Type()}
			}
		}

	}
	return nil
}



type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "logfmt: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "logfmt: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "logfmt: Unmarshal(nil " + e.Type.String() + ")"
}



type UnmarshalTypeError struct {
	Value string       
	Type  reflect.Type 
}

func (e *UnmarshalTypeError) Error() string {
	return "logfmt: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}
