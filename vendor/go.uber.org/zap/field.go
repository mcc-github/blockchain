



















package zap

import (
	"fmt"
	"math"
	"time"

	"go.uber.org/zap/zapcore"
)



type Field = zapcore.Field



func Skip() Field {
	return Field{Type: zapcore.SkipType}
}






func Binary(key string, val []byte) Field {
	return Field{Key: key, Type: zapcore.BinaryType, Interface: val}
}


func Bool(key string, val bool) Field {
	var ival int64
	if val {
		ival = 1
	}
	return Field{Key: key, Type: zapcore.BoolType, Integer: ival}
}




func ByteString(key string, val []byte) Field {
	return Field{Key: key, Type: zapcore.ByteStringType, Interface: val}
}




func Complex128(key string, val complex128) Field {
	return Field{Key: key, Type: zapcore.Complex128Type, Interface: val}
}




func Complex64(key string, val complex64) Field {
	return Field{Key: key, Type: zapcore.Complex64Type, Interface: val}
}




func Float64(key string, val float64) Field {
	return Field{Key: key, Type: zapcore.Float64Type, Integer: int64(math.Float64bits(val))}
}




func Float32(key string, val float32) Field {
	return Field{Key: key, Type: zapcore.Float32Type, Integer: int64(math.Float32bits(val))}
}


func Int(key string, val int) Field {
	return Int64(key, int64(val))
}


func Int64(key string, val int64) Field {
	return Field{Key: key, Type: zapcore.Int64Type, Integer: val}
}


func Int32(key string, val int32) Field {
	return Field{Key: key, Type: zapcore.Int32Type, Integer: int64(val)}
}


func Int16(key string, val int16) Field {
	return Field{Key: key, Type: zapcore.Int16Type, Integer: int64(val)}
}


func Int8(key string, val int8) Field {
	return Field{Key: key, Type: zapcore.Int8Type, Integer: int64(val)}
}


func String(key string, val string) Field {
	return Field{Key: key, Type: zapcore.StringType, String: val}
}


func Uint(key string, val uint) Field {
	return Uint64(key, uint64(val))
}


func Uint64(key string, val uint64) Field {
	return Field{Key: key, Type: zapcore.Uint64Type, Integer: int64(val)}
}


func Uint32(key string, val uint32) Field {
	return Field{Key: key, Type: zapcore.Uint32Type, Integer: int64(val)}
}


func Uint16(key string, val uint16) Field {
	return Field{Key: key, Type: zapcore.Uint16Type, Integer: int64(val)}
}


func Uint8(key string, val uint8) Field {
	return Field{Key: key, Type: zapcore.Uint8Type, Integer: int64(val)}
}


func Uintptr(key string, val uintptr) Field {
	return Field{Key: key, Type: zapcore.UintptrType, Integer: int64(val)}
}








func Reflect(key string, val interface{}) Field {
	return Field{Key: key, Type: zapcore.ReflectType, Interface: val}
}






func Namespace(key string) Field {
	return Field{Key: key, Type: zapcore.NamespaceType}
}



func Stringer(key string, val fmt.Stringer) Field {
	return Field{Key: key, Type: zapcore.StringerType, Interface: val}
}



func Time(key string, val time.Time) Field {
	return Field{Key: key, Type: zapcore.TimeType, Integer: val.UnixNano(), Interface: val.Location()}
}





func Stack(key string) Field {
	
	
	
	
	return String(key, takeStacktrace())
}



func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: zapcore.DurationType, Integer: int64(val)}
}





func Object(key string, val zapcore.ObjectMarshaler) Field {
	return Field{Key: key, Type: zapcore.ObjectMarshalerType, Interface: val}
}








func Any(key string, value interface{}) Field {
	switch val := value.(type) {
	case zapcore.ObjectMarshaler:
		return Object(key, val)
	case zapcore.ArrayMarshaler:
		return Array(key, val)
	case bool:
		return Bool(key, val)
	case []bool:
		return Bools(key, val)
	case complex128:
		return Complex128(key, val)
	case []complex128:
		return Complex128s(key, val)
	case complex64:
		return Complex64(key, val)
	case []complex64:
		return Complex64s(key, val)
	case float64:
		return Float64(key, val)
	case []float64:
		return Float64s(key, val)
	case float32:
		return Float32(key, val)
	case []float32:
		return Float32s(key, val)
	case int:
		return Int(key, val)
	case []int:
		return Ints(key, val)
	case int64:
		return Int64(key, val)
	case []int64:
		return Int64s(key, val)
	case int32:
		return Int32(key, val)
	case []int32:
		return Int32s(key, val)
	case int16:
		return Int16(key, val)
	case []int16:
		return Int16s(key, val)
	case int8:
		return Int8(key, val)
	case []int8:
		return Int8s(key, val)
	case string:
		return String(key, val)
	case []string:
		return Strings(key, val)
	case uint:
		return Uint(key, val)
	case []uint:
		return Uints(key, val)
	case uint64:
		return Uint64(key, val)
	case []uint64:
		return Uint64s(key, val)
	case uint32:
		return Uint32(key, val)
	case []uint32:
		return Uint32s(key, val)
	case uint16:
		return Uint16(key, val)
	case []uint16:
		return Uint16s(key, val)
	case uint8:
		return Uint8(key, val)
	case []byte:
		return Binary(key, val)
	case uintptr:
		return Uintptr(key, val)
	case []uintptr:
		return Uintptrs(key, val)
	case time.Time:
		return Time(key, val)
	case []time.Time:
		return Times(key, val)
	case time.Duration:
		return Duration(key, val)
	case []time.Duration:
		return Durations(key, val)
	case error:
		return NamedError(key, val)
	case []error:
		return Errors(key, val)
	case fmt.Stringer:
		return Stringer(key, val)
	default:
		return Reflect(key, val)
	}
}
