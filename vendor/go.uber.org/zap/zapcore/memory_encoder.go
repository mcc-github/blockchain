



















package zapcore

import "time"




type MapObjectEncoder struct {
	
	Fields map[string]interface{}
	
	cur map[string]interface{}
}


func NewMapObjectEncoder() *MapObjectEncoder {
	m := make(map[string]interface{})
	return &MapObjectEncoder{
		Fields: m,
		cur:    m,
	}
}


func (m *MapObjectEncoder) AddArray(key string, v ArrayMarshaler) error {
	arr := &sliceArrayEncoder{}
	err := v.MarshalLogArray(arr)
	m.cur[key] = arr.elems
	return err
}


func (m *MapObjectEncoder) AddObject(k string, v ObjectMarshaler) error {
	newMap := NewMapObjectEncoder()
	m.cur[k] = newMap.Fields
	return v.MarshalLogObject(newMap)
}


func (m *MapObjectEncoder) AddBinary(k string, v []byte) { m.cur[k] = v }


func (m *MapObjectEncoder) AddByteString(k string, v []byte) { m.cur[k] = string(v) }


func (m *MapObjectEncoder) AddBool(k string, v bool) { m.cur[k] = v }


func (m MapObjectEncoder) AddDuration(k string, v time.Duration) { m.cur[k] = v }


func (m *MapObjectEncoder) AddComplex128(k string, v complex128) { m.cur[k] = v }


func (m *MapObjectEncoder) AddComplex64(k string, v complex64) { m.cur[k] = v }


func (m *MapObjectEncoder) AddFloat64(k string, v float64) { m.cur[k] = v }


func (m *MapObjectEncoder) AddFloat32(k string, v float32) { m.cur[k] = v }


func (m *MapObjectEncoder) AddInt(k string, v int) { m.cur[k] = v }


func (m *MapObjectEncoder) AddInt64(k string, v int64) { m.cur[k] = v }


func (m *MapObjectEncoder) AddInt32(k string, v int32) { m.cur[k] = v }


func (m *MapObjectEncoder) AddInt16(k string, v int16) { m.cur[k] = v }


func (m *MapObjectEncoder) AddInt8(k string, v int8) { m.cur[k] = v }


func (m *MapObjectEncoder) AddString(k string, v string) { m.cur[k] = v }


func (m MapObjectEncoder) AddTime(k string, v time.Time) { m.cur[k] = v }


func (m *MapObjectEncoder) AddUint(k string, v uint) { m.cur[k] = v }


func (m *MapObjectEncoder) AddUint64(k string, v uint64) { m.cur[k] = v }


func (m *MapObjectEncoder) AddUint32(k string, v uint32) { m.cur[k] = v }


func (m *MapObjectEncoder) AddUint16(k string, v uint16) { m.cur[k] = v }


func (m *MapObjectEncoder) AddUint8(k string, v uint8) { m.cur[k] = v }


func (m *MapObjectEncoder) AddUintptr(k string, v uintptr) { m.cur[k] = v }


func (m *MapObjectEncoder) AddReflected(k string, v interface{}) error {
	m.cur[k] = v
	return nil
}


func (m *MapObjectEncoder) OpenNamespace(k string) {
	ns := make(map[string]interface{})
	m.cur[k] = ns
	m.cur = ns
}



type sliceArrayEncoder struct {
	elems []interface{}
}

func (s *sliceArrayEncoder) AppendArray(v ArrayMarshaler) error {
	enc := &sliceArrayEncoder{}
	err := v.MarshalLogArray(enc)
	s.elems = append(s.elems, enc.elems)
	return err
}

func (s *sliceArrayEncoder) AppendObject(v ObjectMarshaler) error {
	m := NewMapObjectEncoder()
	err := v.MarshalLogObject(m)
	s.elems = append(s.elems, m.Fields)
	return err
}

func (s *sliceArrayEncoder) AppendReflected(v interface{}) error {
	s.elems = append(s.elems, v)
	return nil
}

func (s *sliceArrayEncoder) AppendBool(v bool)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendByteString(v []byte)      { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendComplex128(v complex128)  { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendComplex64(v complex64)    { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendDuration(v time.Duration) { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendFloat64(v float64)        { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendFloat32(v float32)        { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt(v int)                { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt64(v int64)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt32(v int32)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt16(v int16)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt8(v int8)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendString(v string)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendTime(v time.Time)         { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint(v uint)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint64(v uint64)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint32(v uint32)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint16(v uint16)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint8(v uint8)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUintptr(v uintptr)        { s.elems = append(s.elems, v) }
