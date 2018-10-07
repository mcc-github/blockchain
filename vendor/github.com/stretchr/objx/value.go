package objx

import (
	"fmt"
	"strconv"
)



type Value struct {
	
	data interface{}
}


func (v *Value) Data() interface{} {
	return v.data
}


func (v *Value) String() string {
	switch {
	case v.IsStr():
		return v.Str()
	case v.IsBool():
		return strconv.FormatBool(v.Bool())
	case v.IsFloat32():
		return strconv.FormatFloat(float64(v.Float32()), 'f', -1, 32)
	case v.IsFloat64():
		return strconv.FormatFloat(v.Float64(), 'f', -1, 64)
	case v.IsInt():
		return strconv.FormatInt(int64(v.Int()), 10)
	case v.IsInt8():
		return strconv.FormatInt(int64(v.Int8()), 10)
	case v.IsInt16():
		return strconv.FormatInt(int64(v.Int16()), 10)
	case v.IsInt32():
		return strconv.FormatInt(int64(v.Int32()), 10)
	case v.IsInt64():
		return strconv.FormatInt(v.Int64(), 10)
	case v.IsUint():
		return strconv.FormatUint(uint64(v.Uint()), 10)
	case v.IsUint8():
		return strconv.FormatUint(uint64(v.Uint8()), 10)
	case v.IsUint16():
		return strconv.FormatUint(uint64(v.Uint16()), 10)
	case v.IsUint32():
		return strconv.FormatUint(uint64(v.Uint32()), 10)
	case v.IsUint64():
		return strconv.FormatUint(v.Uint64(), 10)
	}
	return fmt.Sprintf("%#v", v.Data())
}
