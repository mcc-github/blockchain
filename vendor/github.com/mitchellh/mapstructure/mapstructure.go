






package mapstructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)














type DecodeHookFunc interface{}



type DecodeHookFuncType func(reflect.Type, reflect.Type, interface{}) (interface{}, error)



type DecodeHookFuncKind func(reflect.Kind, reflect.Kind, interface{}) (interface{}, error)



type DecoderConfig struct {
	
	
	
	
	
	
	DecodeHook DecodeHookFunc

	
	
	
	ErrorUnused bool

	
	
	
	ZeroFields bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	WeaklyTypedInput bool

	
	
	Metadata *Metadata

	
	
	Result interface{}

	
	
	TagName string
}







type Decoder struct {
	config *DecoderConfig
}



type Metadata struct {
	
	Keys []string

	
	
	Unused []string
}



func Decode(input interface{}, output interface{}) error {
	config := &DecoderConfig{
		Metadata: nil,
		Result:   output,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}



func WeakDecode(input, output interface{}) error {
	config := &DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}



func DecodeMetadata(input interface{}, output interface{}, metadata *Metadata) error {
	config := &DecoderConfig{
		Metadata: metadata,
		Result:   output,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}




func WeakDecodeMetadata(input interface{}, output interface{}, metadata *Metadata) error {
	config := &DecoderConfig{
		Metadata:         metadata,
		Result:           output,
		WeaklyTypedInput: true,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}




func NewDecoder(config *DecoderConfig) (*Decoder, error) {
	val := reflect.ValueOf(config.Result)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("result must be a pointer")
	}

	val = val.Elem()
	if !val.CanAddr() {
		return nil, errors.New("result must be addressable (a pointer)")
	}

	if config.Metadata != nil {
		if config.Metadata.Keys == nil {
			config.Metadata.Keys = make([]string, 0)
		}

		if config.Metadata.Unused == nil {
			config.Metadata.Unused = make([]string, 0)
		}
	}

	if config.TagName == "" {
		config.TagName = "mapstructure"
	}

	result := &Decoder{
		config: config,
	}

	return result, nil
}



func (d *Decoder) Decode(input interface{}) error {
	return d.decode("", input, reflect.ValueOf(d.config.Result).Elem())
}


func (d *Decoder) decode(name string, input interface{}, outVal reflect.Value) error {
	var inputVal reflect.Value
	if input != nil {
		inputVal = reflect.ValueOf(input)

		
		
		if inputVal.Kind() == reflect.Ptr && inputVal.IsNil() {
			input = nil
		}
	}

	if input == nil {
		
		
		if d.config.ZeroFields {
			outVal.Set(reflect.Zero(outVal.Type()))

			if d.config.Metadata != nil && name != "" {
				d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
			}
		}
		return nil
	}

	if !inputVal.IsValid() {
		
		
		outVal.Set(reflect.Zero(outVal.Type()))
		if d.config.Metadata != nil && name != "" {
			d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
		}
		return nil
	}

	if d.config.DecodeHook != nil {
		
		var err error
		input, err = DecodeHookExec(
			d.config.DecodeHook,
			inputVal.Type(), outVal.Type(), input)
		if err != nil {
			return fmt.Errorf("error decoding '%s': %s", name, err)
		}
	}

	var err error
	outputKind := getKind(outVal)
	switch outputKind {
	case reflect.Bool:
		err = d.decodeBool(name, input, outVal)
	case reflect.Interface:
		err = d.decodeBasic(name, input, outVal)
	case reflect.String:
		err = d.decodeString(name, input, outVal)
	case reflect.Int:
		err = d.decodeInt(name, input, outVal)
	case reflect.Uint:
		err = d.decodeUint(name, input, outVal)
	case reflect.Float32:
		err = d.decodeFloat(name, input, outVal)
	case reflect.Struct:
		err = d.decodeStruct(name, input, outVal)
	case reflect.Map:
		err = d.decodeMap(name, input, outVal)
	case reflect.Ptr:
		err = d.decodePtr(name, input, outVal)
	case reflect.Slice:
		err = d.decodeSlice(name, input, outVal)
	case reflect.Array:
		err = d.decodeArray(name, input, outVal)
	case reflect.Func:
		err = d.decodeFunc(name, input, outVal)
	default:
		
		return fmt.Errorf("%s: unsupported type: %s", name, outputKind)
	}

	
	
	if d.config.Metadata != nil && name != "" {
		d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
	}

	return err
}



func (d *Decoder) decodeBasic(name string, data interface{}, val reflect.Value) error {
	if val.IsValid() && val.Elem().IsValid() {
		return d.decode(name, data, val.Elem())
	}
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if !dataVal.IsValid() {
		dataVal = reflect.Zero(val.Type())
	}

	dataValType := dataVal.Type()
	if !dataValType.AssignableTo(val.Type()) {
		return fmt.Errorf(
			"'%s' expected type '%s', got '%s'",
			name, val.Type(), dataValType)
	}

	val.Set(dataVal)
	return nil
}

func (d *Decoder) decodeString(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	converted := true
	switch {
	case dataKind == reflect.String:
		val.SetString(dataVal.String())
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetString("1")
		} else {
			val.SetString("0")
		}
	case dataKind == reflect.Int && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatInt(dataVal.Int(), 10))
	case dataKind == reflect.Uint && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatUint(dataVal.Uint(), 10))
	case dataKind == reflect.Float32 && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatFloat(dataVal.Float(), 'f', -1, 64))
	case dataKind == reflect.Slice && d.config.WeaklyTypedInput,
		dataKind == reflect.Array && d.config.WeaklyTypedInput:
		dataType := dataVal.Type()
		elemKind := dataType.Elem().Kind()
		switch elemKind {
		case reflect.Uint8:
			var uints []uint8
			if dataKind == reflect.Array {
				uints = make([]uint8, dataVal.Len(), dataVal.Len())
				for i := range uints {
					uints[i] = dataVal.Index(i).Interface().(uint8)
				}
			} else {
				uints = dataVal.Interface().([]uint8)
			}
			val.SetString(string(uints))
		default:
			converted = false
		}
	default:
		converted = false
	}

	if !converted {
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeInt(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetInt(dataVal.Int())
	case dataKind == reflect.Uint:
		val.SetInt(int64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetInt(int64(dataVal.Float()))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetInt(1)
		} else {
			val.SetInt(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		i, err := strconv.ParseInt(dataVal.String(), 0, val.Type().Bits())
		if err == nil {
			val.SetInt(i)
		} else {
			return fmt.Errorf("cannot parse '%s' as int: %s", name, err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Int64()
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into %s: %s", name, err)
		}
		val.SetInt(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeUint(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	switch {
	case dataKind == reflect.Int:
		i := dataVal.Int()
		if i < 0 && !d.config.WeaklyTypedInput {
			return fmt.Errorf("cannot parse '%s', %d overflows uint",
				name, i)
		}
		val.SetUint(uint64(i))
	case dataKind == reflect.Uint:
		val.SetUint(dataVal.Uint())
	case dataKind == reflect.Float32:
		f := dataVal.Float()
		if f < 0 && !d.config.WeaklyTypedInput {
			return fmt.Errorf("cannot parse '%s', %f overflows uint",
				name, f)
		}
		val.SetUint(uint64(f))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetUint(1)
		} else {
			val.SetUint(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		i, err := strconv.ParseUint(dataVal.String(), 0, val.Type().Bits())
		if err == nil {
			val.SetUint(i)
		} else {
			return fmt.Errorf("cannot parse '%s' as uint: %s", name, err)
		}
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeBool(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	switch {
	case dataKind == reflect.Bool:
		val.SetBool(dataVal.Bool())
	case dataKind == reflect.Int && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Int() != 0)
	case dataKind == reflect.Uint && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Uint() != 0)
	case dataKind == reflect.Float32 && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Float() != 0)
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		b, err := strconv.ParseBool(dataVal.String())
		if err == nil {
			val.SetBool(b)
		} else if dataVal.String() == "" {
			val.SetBool(false)
		} else {
			return fmt.Errorf("cannot parse '%s' as bool: %s", name, err)
		}
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeFloat(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetFloat(float64(dataVal.Int()))
	case dataKind == reflect.Uint:
		val.SetFloat(float64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetFloat(dataVal.Float())
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetFloat(1)
		} else {
			val.SetFloat(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		f, err := strconv.ParseFloat(dataVal.String(), val.Type().Bits())
		if err == nil {
			val.SetFloat(f)
		} else {
			return fmt.Errorf("cannot parse '%s' as float: %s", name, err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Float64()
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into %s: %s", name, err)
		}
		val.SetFloat(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeMap(name string, data interface{}, val reflect.Value) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	
	valMap := val

	
	if valMap.IsNil() || d.config.ZeroFields {
		
		mapType := reflect.MapOf(valKeyType, valElemType)
		valMap = reflect.MakeMap(mapType)
	}

	
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	switch dataVal.Kind() {
	case reflect.Map:
		return d.decodeMapFromMap(name, dataVal, val, valMap)

	case reflect.Struct:
		return d.decodeMapFromStruct(name, dataVal, val, valMap)

	case reflect.Array, reflect.Slice:
		if d.config.WeaklyTypedInput {
			return d.decodeMapFromSlice(name, dataVal, val, valMap)
		}

		fallthrough

	default:
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataVal.Kind())
	}
}

func (d *Decoder) decodeMapFromSlice(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error {
	
	if dataVal.Len() == 0 {
		val.Set(valMap)
		return nil
	}

	for i := 0; i < dataVal.Len(); i++ {
		err := d.decode(
			fmt.Sprintf("%s[%d]", name, i),
			dataVal.Index(i).Interface(), val)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodeMapFromMap(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	
	errors := make([]string, 0)

	
	if dataVal.Len() == 0 {
		if dataVal.IsNil() {
			if !val.IsNil() {
				val.Set(dataVal)
			}
		} else {
			
			val.Set(valMap)
		}

		return nil
	}

	for _, k := range dataVal.MapKeys() {
		fieldName := fmt.Sprintf("%s[%s]", name, k)

		
		currentKey := reflect.Indirect(reflect.New(valKeyType))
		if err := d.decode(fieldName, k.Interface(), currentKey); err != nil {
			errors = appendErrors(errors, err)
			continue
		}

		
		v := dataVal.MapIndex(k).Interface()
		currentVal := reflect.Indirect(reflect.New(valElemType))
		if err := d.decode(fieldName, v, currentVal); err != nil {
			errors = appendErrors(errors, err)
			continue
		}

		valMap.SetMapIndex(currentKey, currentVal)
	}

	
	val.Set(valMap)

	
	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func (d *Decoder) decodeMapFromStruct(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error {
	typ := dataVal.Type()
	for i := 0; i < typ.NumField(); i++ {
		
		
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}

		
		
		v := dataVal.Field(i)
		if !v.Type().AssignableTo(valMap.Type().Elem()) {
			return fmt.Errorf("cannot assign type '%s' to map value field of type '%s'", v.Type(), valMap.Type().Elem())
		}

		tagValue := f.Tag.Get(d.config.TagName)
		tagParts := strings.Split(tagValue, ",")

		
		keyName := f.Name
		if tagParts[0] != "" {
			if tagParts[0] == "-" {
				continue
			}
			keyName = tagParts[0]
		}

		
		squash := false
		for _, tag := range tagParts[1:] {
			if tag == "squash" {
				squash = true
				break
			}
		}
		if squash && v.Kind() != reflect.Struct {
			return fmt.Errorf("cannot squash non-struct type '%s'", v.Type())
		}

		switch v.Kind() {
		
		case reflect.Struct:
			x := reflect.New(v.Type())
			x.Elem().Set(v)

			vType := valMap.Type()
			vKeyType := vType.Key()
			vElemType := vType.Elem()
			mType := reflect.MapOf(vKeyType, vElemType)
			vMap := reflect.MakeMap(mType)

			err := d.decode(keyName, x.Interface(), vMap)
			if err != nil {
				return err
			}

			if squash {
				for _, k := range vMap.MapKeys() {
					valMap.SetMapIndex(k, vMap.MapIndex(k))
				}
			} else {
				valMap.SetMapIndex(reflect.ValueOf(keyName), vMap)
			}

		default:
			valMap.SetMapIndex(reflect.ValueOf(keyName), v)
		}
	}

	if val.CanAddr() {
		val.Set(valMap)
	}

	return nil
}

func (d *Decoder) decodePtr(name string, data interface{}, val reflect.Value) error {
	
	
	isNil := data == nil
	if !isNil {
		switch v := reflect.Indirect(reflect.ValueOf(data)); v.Kind() {
		case reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice:
			isNil = v.IsNil()
		}
	}
	if isNil {
		if !val.IsNil() && val.CanSet() {
			nilValue := reflect.New(val.Type()).Elem()
			val.Set(nilValue)
		}

		return nil
	}

	
	
	valType := val.Type()
	valElemType := valType.Elem()
	if val.CanSet() {
		realVal := val
		if realVal.IsNil() || d.config.ZeroFields {
			realVal = reflect.New(valElemType)
		}

		if err := d.decode(name, data, reflect.Indirect(realVal)); err != nil {
			return err
		}

		val.Set(realVal)
	} else {
		if err := d.decode(name, data, reflect.Indirect(val)); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) decodeFunc(name string, data interface{}, val reflect.Value) error {
	
	
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if val.Type() != dataVal.Type() {
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	}
	val.Set(dataVal)
	return nil
}

func (d *Decoder) decodeSlice(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	sliceType := reflect.SliceOf(valElemType)

	valSlice := val
	if valSlice.IsNil() || d.config.ZeroFields {
		if d.config.WeaklyTypedInput {
			switch {
			
			case dataValKind == reflect.Slice, dataValKind == reflect.Array:
				break

			
			case dataValKind == reflect.Map:
				if dataVal.Len() == 0 {
					val.Set(reflect.MakeSlice(sliceType, 0, 0))
					return nil
				}
				
				return d.decodeSlice(name, []interface{}{data}, val)

			case dataValKind == reflect.String && valElemType.Kind() == reflect.Uint8:
				return d.decodeSlice(name, []byte(dataVal.String()), val)

			
			
			default:
				
				return d.decodeSlice(name, []interface{}{data}, val)
			}
		}

		
		if dataValKind != reflect.Array && dataValKind != reflect.Slice {
			return fmt.Errorf(
				"'%s': source data must be an array or slice, got %s", name, dataValKind)

		}

		
		if dataVal.Len() == 0 {
			return nil
		}

		
		valSlice = reflect.MakeSlice(sliceType, dataVal.Len(), dataVal.Len())
	}

	
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		for valSlice.Len() <= i {
			valSlice = reflect.Append(valSlice, reflect.Zero(valElemType))
		}
		currentField := valSlice.Index(i)

		fieldName := fmt.Sprintf("%s[%d]", name, i)
		if err := d.decode(fieldName, currentData, currentField); err != nil {
			errors = appendErrors(errors, err)
		}
	}

	
	val.Set(valSlice)

	
	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func (d *Decoder) decodeArray(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	arrayType := reflect.ArrayOf(valType.Len(), valElemType)

	valArray := val

	if valArray.Interface() == reflect.Zero(valArray.Type()).Interface() || d.config.ZeroFields {
		
		if dataValKind != reflect.Array && dataValKind != reflect.Slice {
			if d.config.WeaklyTypedInput {
				switch {
				
				case dataValKind == reflect.Map:
					if dataVal.Len() == 0 {
						val.Set(reflect.Zero(arrayType))
						return nil
					}

				
				
				default:
					
					return d.decodeArray(name, []interface{}{data}, val)
				}
			}

			return fmt.Errorf(
				"'%s': source data must be an array or slice, got %s", name, dataValKind)

		}
		if dataVal.Len() > arrayType.Len() {
			return fmt.Errorf(
				"'%s': expected source data to have length less or equal to %d, got %d", name, arrayType.Len(), dataVal.Len())

		}

		
		valArray = reflect.New(arrayType).Elem()
	}

	
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		currentField := valArray.Index(i)

		fieldName := fmt.Sprintf("%s[%d]", name, i)
		if err := d.decode(fieldName, currentData, currentField); err != nil {
			errors = appendErrors(errors, err)
		}
	}

	
	val.Set(valArray)

	
	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func (d *Decoder) decodeStruct(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))

	
	
	if dataVal.Type() == val.Type() {
		val.Set(dataVal)
		return nil
	}

	dataValKind := dataVal.Kind()
	switch dataValKind {
	case reflect.Map:
		return d.decodeStructFromMap(name, dataVal, val)

	case reflect.Struct:
		
		
		
		m := make(map[string]interface{})
		mval := reflect.Indirect(reflect.ValueOf(&m))
		if err := d.decodeMapFromStruct(name, dataVal, mval, mval); err != nil {
			return err
		}

		result := d.decodeStructFromMap(name, mval, val)
		return result

	default:
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataVal.Kind())
	}
}

func (d *Decoder) decodeStructFromMap(name string, dataVal, val reflect.Value) error {
	dataValType := dataVal.Type()
	if kind := dataValType.Key().Kind(); kind != reflect.String && kind != reflect.Interface {
		return fmt.Errorf(
			"'%s' needs a map with string keys, has '%s' keys",
			name, dataValType.Key().Kind())
	}

	dataValKeys := make(map[reflect.Value]struct{})
	dataValKeysUnused := make(map[interface{}]struct{})
	for _, dataValKey := range dataVal.MapKeys() {
		dataValKeys[dataValKey] = struct{}{}
		dataValKeysUnused[dataValKey.Interface()] = struct{}{}
	}

	errors := make([]string, 0)

	
	
	
	structs := make([]reflect.Value, 1, 5)
	structs[0] = val

	
	
	type field struct {
		field reflect.StructField
		val   reflect.Value
	}
	fields := []field{}
	for len(structs) > 0 {
		structVal := structs[0]
		structs = structs[1:]

		structType := structVal.Type()

		for i := 0; i < structType.NumField(); i++ {
			fieldType := structType.Field(i)
			fieldKind := fieldType.Type.Kind()

			
			squash := false
			tagParts := strings.Split(fieldType.Tag.Get(d.config.TagName), ",")
			for _, tag := range tagParts[1:] {
				if tag == "squash" {
					squash = true
					break
				}
			}

			if squash {
				if fieldKind != reflect.Struct {
					errors = appendErrors(errors,
						fmt.Errorf("%s: unsupported type for squash: %s", fieldType.Name, fieldKind))
				} else {
					structs = append(structs, structVal.FieldByName(fieldType.Name))
				}
				continue
			}

			
			fields = append(fields, field{fieldType, structVal.Field(i)})
		}
	}

	
	for _, f := range fields {
		field, fieldValue := f.field, f.val
		fieldName := field.Name

		tagValue := field.Tag.Get(d.config.TagName)
		tagValue = strings.SplitN(tagValue, ",", 2)[0]
		if tagValue != "" {
			fieldName = tagValue
		}

		rawMapKey := reflect.ValueOf(fieldName)
		rawMapVal := dataVal.MapIndex(rawMapKey)
		if !rawMapVal.IsValid() {
			
			
			for dataValKey := range dataValKeys {
				mK, ok := dataValKey.Interface().(string)
				if !ok {
					
					continue
				}

				if strings.EqualFold(mK, fieldName) {
					rawMapKey = dataValKey
					rawMapVal = dataVal.MapIndex(dataValKey)
					break
				}
			}

			if !rawMapVal.IsValid() {
				
				
				continue
			}
		}

		
		delete(dataValKeysUnused, rawMapKey.Interface())

		if !fieldValue.IsValid() {
			
			panic("field is not valid")
		}

		
		
		if !fieldValue.CanSet() {
			continue
		}

		
		
		if name != "" {
			fieldName = fmt.Sprintf("%s.%s", name, fieldName)
		}

		if err := d.decode(fieldName, rawMapVal.Interface(), fieldValue); err != nil {
			errors = appendErrors(errors, err)
		}
	}

	if d.config.ErrorUnused && len(dataValKeysUnused) > 0 {
		keys := make([]string, 0, len(dataValKeysUnused))
		for rawKey := range dataValKeysUnused {
			keys = append(keys, rawKey.(string))
		}
		sort.Strings(keys)

		err := fmt.Errorf("'%s' has invalid keys: %s", name, strings.Join(keys, ", "))
		errors = appendErrors(errors, err)
	}

	if len(errors) > 0 {
		return &Error{errors}
	}

	
	if d.config.Metadata != nil {
		for rawKey := range dataValKeysUnused {
			key := rawKey.(string)
			if name != "" {
				key = fmt.Sprintf("%s.%s", name, key)
			}

			d.config.Metadata.Unused = append(d.config.Metadata.Unused, key)
		}
	}

	return nil
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}
