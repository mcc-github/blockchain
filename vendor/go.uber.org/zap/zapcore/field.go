



















package zapcore

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"time"
)



type FieldType uint8

const (
	
	UnknownType FieldType = iota
	
	ArrayMarshalerType
	
	ObjectMarshalerType
	
	BinaryType
	
	BoolType
	
	ByteStringType
	
	Complex128Type
	
	Complex64Type
	
	DurationType
	
	Float64Type
	
	Float32Type
	
	Int64Type
	
	Int32Type
	
	Int16Type
	
	Int8Type
	
	StringType
	
	TimeType
	
	Uint64Type
	
	Uint32Type
	
	Uint16Type
	
	Uint8Type
	
	UintptrType
	
	
	ReflectType
	
	
	NamespaceType
	
	StringerType
	
	ErrorType
	
	SkipType
)




type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Interface interface{}
}



func (f Field) AddTo(enc ObjectEncoder) {
	var err error

	switch f.Type {
	case ArrayMarshalerType:
		err = enc.AddArray(f.Key, f.Interface.(ArrayMarshaler))
	case ObjectMarshalerType:
		err = enc.AddObject(f.Key, f.Interface.(ObjectMarshaler))
	case BinaryType:
		enc.AddBinary(f.Key, f.Interface.([]byte))
	case BoolType:
		enc.AddBool(f.Key, f.Integer == 1)
	case ByteStringType:
		enc.AddByteString(f.Key, f.Interface.([]byte))
	case Complex128Type:
		enc.AddComplex128(f.Key, f.Interface.(complex128))
	case Complex64Type:
		enc.AddComplex64(f.Key, f.Interface.(complex64))
	case DurationType:
		enc.AddDuration(f.Key, time.Duration(f.Integer))
	case Float64Type:
		enc.AddFloat64(f.Key, math.Float64frombits(uint64(f.Integer)))
	case Float32Type:
		enc.AddFloat32(f.Key, math.Float32frombits(uint32(f.Integer)))
	case Int64Type:
		enc.AddInt64(f.Key, f.Integer)
	case Int32Type:
		enc.AddInt32(f.Key, int32(f.Integer))
	case Int16Type:
		enc.AddInt16(f.Key, int16(f.Integer))
	case Int8Type:
		enc.AddInt8(f.Key, int8(f.Integer))
	case StringType:
		enc.AddString(f.Key, f.String)
	case TimeType:
		if f.Interface != nil {
			enc.AddTime(f.Key, time.Unix(0, f.Integer).In(f.Interface.(*time.Location)))
		} else {
			
			enc.AddTime(f.Key, time.Unix(0, f.Integer))
		}
	case Uint64Type:
		enc.AddUint64(f.Key, uint64(f.Integer))
	case Uint32Type:
		enc.AddUint32(f.Key, uint32(f.Integer))
	case Uint16Type:
		enc.AddUint16(f.Key, uint16(f.Integer))
	case Uint8Type:
		enc.AddUint8(f.Key, uint8(f.Integer))
	case UintptrType:
		enc.AddUintptr(f.Key, uintptr(f.Integer))
	case ReflectType:
		err = enc.AddReflected(f.Key, f.Interface)
	case NamespaceType:
		enc.OpenNamespace(f.Key)
	case StringerType:
		enc.AddString(f.Key, f.Interface.(fmt.Stringer).String())
	case ErrorType:
		encodeError(f.Key, f.Interface.(error), enc)
	case SkipType:
		break
	default:
		panic(fmt.Sprintf("unknown field type: %v", f))
	}

	if err != nil {
		enc.AddString(fmt.Sprintf("%sError", f.Key), err.Error())
	}
}



func (f Field) Equals(other Field) bool {
	if f.Type != other.Type {
		return false
	}
	if f.Key != other.Key {
		return false
	}

	switch f.Type {
	case BinaryType, ByteStringType:
		return bytes.Equal(f.Interface.([]byte), other.Interface.([]byte))
	case ArrayMarshalerType, ObjectMarshalerType, ErrorType, ReflectType:
		return reflect.DeepEqual(f.Interface, other.Interface)
	default:
		return f == other
	}
}

func addFields(enc ObjectEncoder, fields []Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
