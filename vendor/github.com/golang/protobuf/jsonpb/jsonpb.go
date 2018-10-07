































package jsonpb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"

	stpb "github.com/golang/protobuf/ptypes/struct"
)

const secondInNanos = int64(time.Second / time.Nanosecond)



type Marshaler struct {
	
	EnumsAsInts bool

	
	EmitDefaults bool

	
	
	
	
	Indent string

	
	OrigName bool

	
	
	
	
	AnyResolver AnyResolver
}



type AnyResolver interface {
	Resolve(typeUrl string) (proto.Message, error)
}

func defaultResolveAny(typeUrl string) (proto.Message, error) {
	
	mname := typeUrl
	if slash := strings.LastIndex(mname, "/"); slash >= 0 {
		mname = mname[slash+1:]
	}
	mt := proto.MessageType(mname)
	if mt == nil {
		return nil, fmt.Errorf("unknown message type %q", mname)
	}
	return reflect.New(mt.Elem()).Interface().(proto.Message), nil
}








type JSONPBMarshaler interface {
	MarshalJSONPB(*Marshaler) ([]byte, error)
}








type JSONPBUnmarshaler interface {
	UnmarshalJSONPB(*Unmarshaler, []byte) error
}


func (m *Marshaler) Marshal(out io.Writer, pb proto.Message) error {
	v := reflect.ValueOf(pb)
	if pb == nil || (v.Kind() == reflect.Ptr && v.IsNil()) {
		return errors.New("Marshal called with nil")
	}
	
	if err := checkRequiredFields(pb); err != nil {
		return err
	}
	writer := &errWriter{writer: out}
	return m.marshalObject(writer, pb, "", "")
}


func (m *Marshaler) MarshalToString(pb proto.Message) (string, error) {
	var buf bytes.Buffer
	if err := m.Marshal(&buf, pb); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type int32Slice []int32

var nonFinite = map[string]float64{
	`"NaN"`:       math.NaN(),
	`"Infinity"`:  math.Inf(1),
	`"-Infinity"`: math.Inf(-1),
}


func (s int32Slice) Len() int           { return len(s) }
func (s int32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s int32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type wkt interface {
	XXX_WellKnownType() string
}


func (m *Marshaler) marshalObject(out *errWriter, v proto.Message, indent, typeURL string) error {
	if jsm, ok := v.(JSONPBMarshaler); ok {
		b, err := jsm.MarshalJSONPB(m)
		if err != nil {
			return err
		}
		if typeURL != "" {
			
			var js map[string]*json.RawMessage
			if err = json.Unmarshal(b, &js); err != nil {
				return fmt.Errorf("type %T produced invalid JSON: %v", v, err)
			}
			turl, err := json.Marshal(typeURL)
			if err != nil {
				return fmt.Errorf("failed to marshal type URL %q to JSON: %v", typeURL, err)
			}
			js["@type"] = (*json.RawMessage)(&turl)
			if b, err = json.Marshal(js); err != nil {
				return err
			}
		}

		out.write(string(b))
		return out.err
	}

	s := reflect.ValueOf(v).Elem()

	
	if wkt, ok := v.(wkt); ok {
		switch wkt.XXX_WellKnownType() {
		case "DoubleValue", "FloatValue", "Int64Value", "UInt64Value",
			"Int32Value", "UInt32Value", "BoolValue", "StringValue", "BytesValue":
			
			
			sprop := proto.GetProperties(s.Type())
			return m.marshalValue(out, sprop.Prop[0], s.Field(0), indent)
		case "Any":
			
			return m.marshalAny(out, v, indent)
		case "Duration":
			
			
			s, ns := s.Field(0).Int(), s.Field(1).Int()
			if ns <= -secondInNanos || ns >= secondInNanos {
				return fmt.Errorf("ns out of range (%v, %v)", -secondInNanos, secondInNanos)
			}
			if (s > 0 && ns < 0) || (s < 0 && ns > 0) {
				return errors.New("signs of seconds and nanos do not match")
			}
			if s < 0 {
				ns = -ns
			}
			x := fmt.Sprintf("%d.%09d", s, ns)
			x = strings.TrimSuffix(x, "000")
			x = strings.TrimSuffix(x, "000")
			x = strings.TrimSuffix(x, ".000")
			out.write(`"`)
			out.write(x)
			out.write(`s"`)
			return out.err
		case "Struct", "ListValue":
			
			
			return m.marshalValue(out, &proto.Properties{}, s.Field(0), indent)
		case "Timestamp":
			
			
			s, ns := s.Field(0).Int(), s.Field(1).Int()
			if ns < 0 || ns >= secondInNanos {
				return fmt.Errorf("ns out of range [0, %v)", secondInNanos)
			}
			t := time.Unix(s, ns).UTC()
			
			x := t.Format("2006-01-02T15:04:05.000000000")
			x = strings.TrimSuffix(x, "000")
			x = strings.TrimSuffix(x, "000")
			x = strings.TrimSuffix(x, ".000")
			out.write(`"`)
			out.write(x)
			out.write(`Z"`)
			return out.err
		case "Value":
			
			kind := s.Field(0)
			if kind.IsNil() {
				
				return errors.New("nil Value")
			}
			
			x := kind.Elem().Elem().Field(0)
			
			return m.marshalValue(out, &proto.Properties{}, x, indent)
		}
	}

	out.write("{")
	if m.Indent != "" {
		out.write("\n")
	}

	firstField := true

	if typeURL != "" {
		if err := m.marshalTypeURL(out, indent, typeURL); err != nil {
			return err
		}
		firstField = false
	}

	for i := 0; i < s.NumField(); i++ {
		value := s.Field(i)
		valueField := s.Type().Field(i)
		if strings.HasPrefix(valueField.Name, "XXX_") {
			continue
		}

		
		switch value.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface:
			if value.IsNil() {
				continue
			}
		}

		if !m.EmitDefaults {
			switch value.Kind() {
			case reflect.Bool:
				if !value.Bool() {
					continue
				}
			case reflect.Int32, reflect.Int64:
				if value.Int() == 0 {
					continue
				}
			case reflect.Uint32, reflect.Uint64:
				if value.Uint() == 0 {
					continue
				}
			case reflect.Float32, reflect.Float64:
				if value.Float() == 0 {
					continue
				}
			case reflect.String:
				if value.Len() == 0 {
					continue
				}
			case reflect.Map, reflect.Ptr, reflect.Slice:
				if value.IsNil() {
					continue
				}
			}
		}

		
		if valueField.Tag.Get("protobuf_oneof") != "" {
			
			sv := value.Elem().Elem() 
			value = sv.Field(0)
			valueField = sv.Type().Field(0)
		}
		prop := jsonProperties(valueField, m.OrigName)
		if !firstField {
			m.writeSep(out)
		}
		if err := m.marshalField(out, prop, value, indent); err != nil {
			return err
		}
		firstField = false
	}

	
	if ep, ok := v.(proto.Message); ok {
		extensions := proto.RegisteredExtensions(v)
		
		ids := make([]int32, 0, len(extensions))
		for id, desc := range extensions {
			if !proto.HasExtension(ep, desc) {
				continue
			}
			ids = append(ids, id)
		}
		sort.Sort(int32Slice(ids))
		for _, id := range ids {
			desc := extensions[id]
			if desc == nil {
				
				continue
			}
			ext, extErr := proto.GetExtension(ep, desc)
			if extErr != nil {
				return extErr
			}
			value := reflect.ValueOf(ext)
			var prop proto.Properties
			prop.Parse(desc.Tag)
			prop.JSONName = fmt.Sprintf("[%s]", desc.Name)
			if !firstField {
				m.writeSep(out)
			}
			if err := m.marshalField(out, &prop, value, indent); err != nil {
				return err
			}
			firstField = false
		}

	}

	if m.Indent != "" {
		out.write("\n")
		out.write(indent)
	}
	out.write("}")
	return out.err
}

func (m *Marshaler) writeSep(out *errWriter) {
	if m.Indent != "" {
		out.write(",\n")
	} else {
		out.write(",")
	}
}

func (m *Marshaler) marshalAny(out *errWriter, any proto.Message, indent string) error {
	
	
	
	
	v := reflect.ValueOf(any).Elem()
	turl := v.Field(0).String()
	val := v.Field(1).Bytes()

	var msg proto.Message
	var err error
	if m.AnyResolver != nil {
		msg, err = m.AnyResolver.Resolve(turl)
	} else {
		msg, err = defaultResolveAny(turl)
	}
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(val, msg); err != nil {
		return err
	}

	if _, ok := msg.(wkt); ok {
		out.write("{")
		if m.Indent != "" {
			out.write("\n")
		}
		if err := m.marshalTypeURL(out, indent, turl); err != nil {
			return err
		}
		m.writeSep(out)
		if m.Indent != "" {
			out.write(indent)
			out.write(m.Indent)
			out.write(`"value": `)
		} else {
			out.write(`"value":`)
		}
		if err := m.marshalObject(out, msg, indent+m.Indent, ""); err != nil {
			return err
		}
		if m.Indent != "" {
			out.write("\n")
			out.write(indent)
		}
		out.write("}")
		return out.err
	}

	return m.marshalObject(out, msg, indent, turl)
}

func (m *Marshaler) marshalTypeURL(out *errWriter, indent, typeURL string) error {
	if m.Indent != "" {
		out.write(indent)
		out.write(m.Indent)
	}
	out.write(`"@type":`)
	if m.Indent != "" {
		out.write(" ")
	}
	b, err := json.Marshal(typeURL)
	if err != nil {
		return err
	}
	out.write(string(b))
	return out.err
}


func (m *Marshaler) marshalField(out *errWriter, prop *proto.Properties, v reflect.Value, indent string) error {
	if m.Indent != "" {
		out.write(indent)
		out.write(m.Indent)
	}
	out.write(`"`)
	out.write(prop.JSONName)
	out.write(`":`)
	if m.Indent != "" {
		out.write(" ")
	}
	if err := m.marshalValue(out, prop, v, indent); err != nil {
		return err
	}
	return nil
}


func (m *Marshaler) marshalValue(out *errWriter, prop *proto.Properties, v reflect.Value, indent string) error {
	var err error
	v = reflect.Indirect(v)

	
	if v.Kind() == reflect.Invalid {
		out.write("null")
		return out.err
	}

	
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint8 {
		out.write("[")
		comma := ""
		for i := 0; i < v.Len(); i++ {
			sliceVal := v.Index(i)
			out.write(comma)
			if m.Indent != "" {
				out.write("\n")
				out.write(indent)
				out.write(m.Indent)
				out.write(m.Indent)
			}
			if err := m.marshalValue(out, prop, sliceVal, indent+m.Indent); err != nil {
				return err
			}
			comma = ","
		}
		if m.Indent != "" {
			out.write("\n")
			out.write(indent)
			out.write(m.Indent)
		}
		out.write("]")
		return out.err
	}

	
	
	if wkt, ok := v.Interface().(wkt); ok {
		switch wkt.XXX_WellKnownType() {
		case "NullValue":
			out.write("null")
			return out.err
		}
	}

	
	if !m.EnumsAsInts && prop.Enum != "" {
		
		
		
		enumStr := v.Interface().(fmt.Stringer).String()
		var valStr string
		if v.Kind() == reflect.Ptr {
			valStr = strconv.Itoa(int(v.Elem().Int()))
		} else {
			valStr = strconv.Itoa(int(v.Int()))
		}
		isKnownEnum := enumStr != valStr
		if isKnownEnum {
			out.write(`"`)
		}
		out.write(enumStr)
		if isKnownEnum {
			out.write(`"`)
		}
		return out.err
	}

	
	if v.Kind() == reflect.Struct {
		return m.marshalObject(out, v.Addr().Interface().(proto.Message), indent+m.Indent, "")
	}

	
	
	if v.Kind() == reflect.Map {
		out.write(`{`)
		keys := v.MapKeys()
		sort.Sort(mapKeys(keys))
		for i, k := range keys {
			if i > 0 {
				out.write(`,`)
			}
			if m.Indent != "" {
				out.write("\n")
				out.write(indent)
				out.write(m.Indent)
				out.write(m.Indent)
			}

			
			b, err := json.Marshal(k.Interface())
			if err != nil {
				return err
			}
			s := string(b)

			
			if !strings.HasPrefix(s, `"`) {
				b, err := json.Marshal(s)
				if err != nil {
					return err
				}
				s = string(b)
			}

			out.write(s)
			out.write(`:`)
			if m.Indent != "" {
				out.write(` `)
			}

			vprop := prop
			if prop != nil && prop.MapValProp != nil {
				vprop = prop.MapValProp
			}
			if err := m.marshalValue(out, vprop, v.MapIndex(k), indent+m.Indent); err != nil {
				return err
			}
		}
		if m.Indent != "" {
			out.write("\n")
			out.write(indent)
			out.write(m.Indent)
		}
		out.write(`}`)
		return out.err
	}

	
	if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
		f := v.Float()
		var sval string
		switch {
		case math.IsInf(f, 1):
			sval = `"Infinity"`
		case math.IsInf(f, -1):
			sval = `"-Infinity"`
		case math.IsNaN(f):
			sval = `"NaN"`
		}
		if sval != "" {
			out.write(sval)
			return out.err
		}
	}

	
	b, err := json.Marshal(v.Interface())
	if err != nil {
		return err
	}
	needToQuote := string(b[0]) != `"` && (v.Kind() == reflect.Int64 || v.Kind() == reflect.Uint64)
	if needToQuote {
		out.write(`"`)
	}
	out.write(string(b))
	if needToQuote {
		out.write(`"`)
	}
	return out.err
}



type Unmarshaler struct {
	
	
	AllowUnknownFields bool

	
	
	
	
	AnyResolver AnyResolver
}




func (u *Unmarshaler) UnmarshalNext(dec *json.Decoder, pb proto.Message) error {
	inputValue := json.RawMessage{}
	if err := dec.Decode(&inputValue); err != nil {
		return err
	}
	if err := u.unmarshalValue(reflect.ValueOf(pb).Elem(), inputValue, nil); err != nil {
		return err
	}
	return checkRequiredFields(pb)
}




func (u *Unmarshaler) Unmarshal(r io.Reader, pb proto.Message) error {
	dec := json.NewDecoder(r)
	return u.UnmarshalNext(dec, pb)
}




func UnmarshalNext(dec *json.Decoder, pb proto.Message) error {
	return new(Unmarshaler).UnmarshalNext(dec, pb)
}




func Unmarshal(r io.Reader, pb proto.Message) error {
	return new(Unmarshaler).Unmarshal(r, pb)
}




func UnmarshalString(str string, pb proto.Message) error {
	return new(Unmarshaler).Unmarshal(strings.NewReader(str), pb)
}



func (u *Unmarshaler) unmarshalValue(target reflect.Value, inputValue json.RawMessage, prop *proto.Properties) error {
	targetType := target.Type()

	
	if targetType.Kind() == reflect.Ptr {
		
		
		_, isJSONPBUnmarshaler := target.Interface().(JSONPBUnmarshaler)
		if string(inputValue) == "null" && targetType != reflect.TypeOf(&stpb.Value{}) && !isJSONPBUnmarshaler {
			return nil
		}
		target.Set(reflect.New(targetType.Elem()))

		return u.unmarshalValue(target.Elem(), inputValue, prop)
	}

	if jsu, ok := target.Addr().Interface().(JSONPBUnmarshaler); ok {
		return jsu.UnmarshalJSONPB(u, []byte(inputValue))
	}

	
	if w, ok := target.Addr().Interface().(wkt); ok {
		switch w.XXX_WellKnownType() {
		case "DoubleValue", "FloatValue", "Int64Value", "UInt64Value",
			"Int32Value", "UInt32Value", "BoolValue", "StringValue", "BytesValue":
			return u.unmarshalValue(target.Field(0), inputValue, prop)
		case "Any":
			
			
			
			var jsonFields map[string]*json.RawMessage
			if err := json.Unmarshal(inputValue, &jsonFields); err != nil {
				return err
			}

			val, ok := jsonFields["@type"]
			if !ok || val == nil {
				return errors.New("Any JSON doesn't have '@type'")
			}

			var turl string
			if err := json.Unmarshal([]byte(*val), &turl); err != nil {
				return fmt.Errorf("can't unmarshal Any's '@type': %q", *val)
			}
			target.Field(0).SetString(turl)

			var m proto.Message
			var err error
			if u.AnyResolver != nil {
				m, err = u.AnyResolver.Resolve(turl)
			} else {
				m, err = defaultResolveAny(turl)
			}
			if err != nil {
				return err
			}

			if _, ok := m.(wkt); ok {
				val, ok := jsonFields["value"]
				if !ok {
					return errors.New("Any JSON doesn't have 'value'")
				}

				if err := u.unmarshalValue(reflect.ValueOf(m).Elem(), *val, nil); err != nil {
					return fmt.Errorf("can't unmarshal Any nested proto %T: %v", m, err)
				}
			} else {
				delete(jsonFields, "@type")
				nestedProto, err := json.Marshal(jsonFields)
				if err != nil {
					return fmt.Errorf("can't generate JSON for Any's nested proto to be unmarshaled: %v", err)
				}

				if err = u.unmarshalValue(reflect.ValueOf(m).Elem(), nestedProto, nil); err != nil {
					return fmt.Errorf("can't unmarshal Any nested proto %T: %v", m, err)
				}
			}

			b, err := proto.Marshal(m)
			if err != nil {
				return fmt.Errorf("can't marshal proto %T into Any.Value: %v", m, err)
			}
			target.Field(1).SetBytes(b)

			return nil
		case "Duration":
			unq, err := unquote(string(inputValue))
			if err != nil {
				return err
			}

			d, err := time.ParseDuration(unq)
			if err != nil {
				return fmt.Errorf("bad Duration: %v", err)
			}

			ns := d.Nanoseconds()
			s := ns / 1e9
			ns %= 1e9
			target.Field(0).SetInt(s)
			target.Field(1).SetInt(ns)
			return nil
		case "Timestamp":
			unq, err := unquote(string(inputValue))
			if err != nil {
				return err
			}

			t, err := time.Parse(time.RFC3339Nano, unq)
			if err != nil {
				return fmt.Errorf("bad Timestamp: %v", err)
			}

			target.Field(0).SetInt(t.Unix())
			target.Field(1).SetInt(int64(t.Nanosecond()))
			return nil
		case "Struct":
			var m map[string]json.RawMessage
			if err := json.Unmarshal(inputValue, &m); err != nil {
				return fmt.Errorf("bad StructValue: %v", err)
			}

			target.Field(0).Set(reflect.ValueOf(map[string]*stpb.Value{}))
			for k, jv := range m {
				pv := &stpb.Value{}
				if err := u.unmarshalValue(reflect.ValueOf(pv).Elem(), jv, prop); err != nil {
					return fmt.Errorf("bad value in StructValue for key %q: %v", k, err)
				}
				target.Field(0).SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(pv))
			}
			return nil
		case "ListValue":
			var s []json.RawMessage
			if err := json.Unmarshal(inputValue, &s); err != nil {
				return fmt.Errorf("bad ListValue: %v", err)
			}

			target.Field(0).Set(reflect.ValueOf(make([]*stpb.Value, len(s))))
			for i, sv := range s {
				if err := u.unmarshalValue(target.Field(0).Index(i), sv, prop); err != nil {
					return err
				}
			}
			return nil
		case "Value":
			ivStr := string(inputValue)
			if ivStr == "null" {
				target.Field(0).Set(reflect.ValueOf(&stpb.Value_NullValue{}))
			} else if v, err := strconv.ParseFloat(ivStr, 0); err == nil {
				target.Field(0).Set(reflect.ValueOf(&stpb.Value_NumberValue{v}))
			} else if v, err := unquote(ivStr); err == nil {
				target.Field(0).Set(reflect.ValueOf(&stpb.Value_StringValue{v}))
			} else if v, err := strconv.ParseBool(ivStr); err == nil {
				target.Field(0).Set(reflect.ValueOf(&stpb.Value_BoolValue{v}))
			} else if err := json.Unmarshal(inputValue, &[]json.RawMessage{}); err == nil {
				lv := &stpb.ListValue{}
				target.Field(0).Set(reflect.ValueOf(&stpb.Value_ListValue{lv}))
				return u.unmarshalValue(reflect.ValueOf(lv).Elem(), inputValue, prop)
			} else if err := json.Unmarshal(inputValue, &map[string]json.RawMessage{}); err == nil {
				sv := &stpb.Struct{}
				target.Field(0).Set(reflect.ValueOf(&stpb.Value_StructValue{sv}))
				return u.unmarshalValue(reflect.ValueOf(sv).Elem(), inputValue, prop)
			} else {
				return fmt.Errorf("unrecognized type for Value %q", ivStr)
			}
			return nil
		}
	}

	
	
	
	
	if inputValue[0] == '"' && prop != nil && prop.Enum != "" {
		vmap := proto.EnumValueMap(prop.Enum)
		
		
		s := inputValue[1 : len(inputValue)-1]
		n, ok := vmap[string(s)]
		if !ok {
			return fmt.Errorf("unknown value %q for enum %s", s, prop.Enum)
		}
		if target.Kind() == reflect.Ptr { 
			target.Set(reflect.New(targetType.Elem()))
			target = target.Elem()
		}
		if targetType.Kind() != reflect.Int32 {
			return fmt.Errorf("invalid target %q for enum %s", targetType.Kind(), prop.Enum)
		}
		target.SetInt(int64(n))
		return nil
	}

	
	if targetType.Kind() == reflect.Struct {
		var jsonFields map[string]json.RawMessage
		if err := json.Unmarshal(inputValue, &jsonFields); err != nil {
			return err
		}

		consumeField := func(prop *proto.Properties) (json.RawMessage, bool) {
			
			fieldNames := acceptedJSONFieldNames(prop)

			vOrig, okOrig := jsonFields[fieldNames.orig]
			vCamel, okCamel := jsonFields[fieldNames.camel]
			if !okOrig && !okCamel {
				return nil, false
			}
			
			var raw json.RawMessage
			if okOrig {
				raw = vOrig
				delete(jsonFields, fieldNames.orig)
			}
			if okCamel {
				raw = vCamel
				delete(jsonFields, fieldNames.camel)
			}
			return raw, true
		}

		sprops := proto.GetProperties(targetType)
		for i := 0; i < target.NumField(); i++ {
			ft := target.Type().Field(i)
			if strings.HasPrefix(ft.Name, "XXX_") {
				continue
			}

			valueForField, ok := consumeField(sprops.Prop[i])
			if !ok {
				continue
			}

			if err := u.unmarshalValue(target.Field(i), valueForField, sprops.Prop[i]); err != nil {
				return err
			}
		}
		
		if len(jsonFields) > 0 {
			for _, oop := range sprops.OneofTypes {
				raw, ok := consumeField(oop.Prop)
				if !ok {
					continue
				}
				nv := reflect.New(oop.Type.Elem())
				target.Field(oop.Field).Set(nv)
				if err := u.unmarshalValue(nv.Elem().Field(0), raw, oop.Prop); err != nil {
					return err
				}
			}
		}
		
		if len(jsonFields) > 0 {
			if ep, ok := target.Addr().Interface().(proto.Message); ok {
				for _, ext := range proto.RegisteredExtensions(ep) {
					name := fmt.Sprintf("[%s]", ext.Name)
					raw, ok := jsonFields[name]
					if !ok {
						continue
					}
					delete(jsonFields, name)
					nv := reflect.New(reflect.TypeOf(ext.ExtensionType).Elem())
					if err := u.unmarshalValue(nv.Elem(), raw, nil); err != nil {
						return err
					}
					if err := proto.SetExtension(ep, ext, nv.Interface()); err != nil {
						return err
					}
				}
			}
		}
		if !u.AllowUnknownFields && len(jsonFields) > 0 {
			
			var f string
			for fname := range jsonFields {
				f = fname
				break
			}
			return fmt.Errorf("unknown field %q in %v", f, targetType)
		}
		return nil
	}

	
	if targetType.Kind() == reflect.Slice && targetType.Elem().Kind() != reflect.Uint8 {
		var slc []json.RawMessage
		if err := json.Unmarshal(inputValue, &slc); err != nil {
			return err
		}
		if slc != nil {
			l := len(slc)
			target.Set(reflect.MakeSlice(targetType, l, l))
			for i := 0; i < l; i++ {
				if err := u.unmarshalValue(target.Index(i), slc[i], prop); err != nil {
					return err
				}
			}
		}
		return nil
	}

	
	if targetType.Kind() == reflect.Map {
		var mp map[string]json.RawMessage
		if err := json.Unmarshal(inputValue, &mp); err != nil {
			return err
		}
		if mp != nil {
			target.Set(reflect.MakeMap(targetType))
			for ks, raw := range mp {
				
				
				var k reflect.Value
				if targetType.Key().Kind() == reflect.String {
					k = reflect.ValueOf(ks)
				} else {
					k = reflect.New(targetType.Key()).Elem()
					var kprop *proto.Properties
					if prop != nil && prop.MapKeyProp != nil {
						kprop = prop.MapKeyProp
					}
					if err := u.unmarshalValue(k, json.RawMessage(ks), kprop); err != nil {
						return err
					}
				}

				
				v := reflect.New(targetType.Elem()).Elem()
				var vprop *proto.Properties
				if prop != nil && prop.MapValProp != nil {
					vprop = prop.MapValProp
				}
				if err := u.unmarshalValue(v, raw, vprop); err != nil {
					return err
				}
				target.SetMapIndex(k, v)
			}
		}
		return nil
	}

	
	isFloat := targetType.Kind() == reflect.Float32 || targetType.Kind() == reflect.Float64
	if isFloat {
		if num, ok := nonFinite[string(inputValue)]; ok {
			target.SetFloat(num)
			return nil
		}
	}

	
	
	isNum := targetType.Kind() == reflect.Int64 || targetType.Kind() == reflect.Uint64 ||
		targetType.Kind() == reflect.Int32 || targetType.Kind() == reflect.Uint32 ||
		targetType.Kind() == reflect.Float32 || targetType.Kind() == reflect.Float64
	if isNum && strings.HasPrefix(string(inputValue), `"`) {
		inputValue = inputValue[1 : len(inputValue)-1]
	}

	
	return json.Unmarshal(inputValue, target.Addr().Interface())
}

func unquote(s string) (string, error) {
	var ret string
	err := json.Unmarshal([]byte(s), &ret)
	return ret, err
}


func jsonProperties(f reflect.StructField, origName bool) *proto.Properties {
	var prop proto.Properties
	prop.Init(f.Type, f.Name, f.Tag.Get("protobuf"), &f)
	if origName || prop.JSONName == "" {
		prop.JSONName = prop.OrigName
	}
	return &prop
}

type fieldNames struct {
	orig, camel string
}

func acceptedJSONFieldNames(prop *proto.Properties) fieldNames {
	opts := fieldNames{orig: prop.OrigName, camel: prop.OrigName}
	if prop.JSONName != "" {
		opts.camel = prop.JSONName
	}
	return opts
}


type errWriter struct {
	writer io.Writer
	err    error
}

func (w *errWriter) write(str string) {
	if w.err != nil {
		return
	}
	_, w.err = w.writer.Write([]byte(str))
}








type mapKeys []reflect.Value

func (s mapKeys) Len() int      { return len(s) }
func (s mapKeys) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s mapKeys) Less(i, j int) bool {
	if k := s[i].Kind(); k == s[j].Kind() {
		switch k {
		case reflect.String:
			return s[i].String() < s[j].String()
		case reflect.Int32, reflect.Int64:
			return s[i].Int() < s[j].Int()
		case reflect.Uint32, reflect.Uint64:
			return s[i].Uint() < s[j].Uint()
		}
	}
	return fmt.Sprint(s[i].Interface()) < fmt.Sprint(s[j].Interface())
}




func checkRequiredFields(pb proto.Message) error {
	
	
	
	
	
	
	
	
	
	
	if _, ok := pb.(wkt); ok {
		return nil
	}

	v := reflect.ValueOf(pb)
	
	if v.Kind() != reflect.Ptr {
		return nil
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		sfield := v.Type().Field(i)

		if sfield.PkgPath != "" {
			
			continue
		}

		if strings.HasPrefix(sfield.Name, "XXX_") {
			continue
		}

		
		
		if sfield.Tag.Get("protobuf_oneof") != "" {
			if field.Kind() != reflect.Interface {
				continue
			}
			v := field.Elem()
			if v.Kind() != reflect.Ptr || v.IsNil() {
				continue
			}
			v = v.Elem()
			if v.Kind() != reflect.Struct || v.NumField() < 1 {
				continue
			}
			field = v.Field(0)
			sfield = v.Type().Field(0)
		}

		protoTag := sfield.Tag.Get("protobuf")
		if protoTag == "" {
			continue
		}
		var prop proto.Properties
		prop.Init(sfield.Type, sfield.Name, protoTag, &sfield)

		switch field.Kind() {
		case reflect.Map:
			if field.IsNil() {
				continue
			}
			
			keys := field.MapKeys()
			for _, k := range keys {
				v := field.MapIndex(k)
				if err := checkRequiredFieldsInValue(v); err != nil {
					return err
				}
			}
		case reflect.Slice:
			
			if !prop.Repeated {
				if prop.Required && field.IsNil() {
					return fmt.Errorf("required field %q is not set", prop.Name)
				}
				continue
			}

			
			if field.IsNil() {
				continue
			}
			
			for i := 0; i < field.Len(); i++ {
				v := field.Index(i)
				if err := checkRequiredFieldsInValue(v); err != nil {
					return err
				}
			}
		case reflect.Ptr:
			if field.IsNil() {
				if prop.Required {
					return fmt.Errorf("required field %q is not set", prop.Name)
				}
				continue
			}
			if err := checkRequiredFieldsInValue(field); err != nil {
				return err
			}
		}
	}

	
	for _, ext := range proto.RegisteredExtensions(pb) {
		if !proto.HasExtension(pb, ext) {
			continue
		}
		ep, err := proto.GetExtension(pb, ext)
		if err != nil {
			return err
		}
		err = checkRequiredFieldsInValue(reflect.ValueOf(ep))
		if err != nil {
			return err
		}
	}

	return nil
}

func checkRequiredFieldsInValue(v reflect.Value) error {
	if pm, ok := v.Interface().(proto.Message); ok {
		return checkRequiredFields(pm)
	}
	return nil
}
