



















package zapcore

import (
	"time"

	"go.uber.org/zap/buffer"
)




const DefaultLineEnding = "\n"


type LevelEncoder func(Level, PrimitiveArrayEncoder)



func LowercaseLevelEncoder(l Level, enc PrimitiveArrayEncoder) {
	enc.AppendString(l.String())
}



func LowercaseColorLevelEncoder(l Level, enc PrimitiveArrayEncoder) {
	s, ok := _levelToLowercaseColorString[l]
	if !ok {
		s = _unknownLevelColor.Add(l.String())
	}
	enc.AppendString(s)
}



func CapitalLevelEncoder(l Level, enc PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}



func CapitalColorLevelEncoder(l Level, enc PrimitiveArrayEncoder) {
	s, ok := _levelToCapitalColorString[l]
	if !ok {
		s = _unknownLevelColor.Add(l.CapitalString())
	}
	enc.AppendString(s)
}





func (e *LevelEncoder) UnmarshalText(text []byte) error {
	switch string(text) {
	case "capital":
		*e = CapitalLevelEncoder
	case "capitalColor":
		*e = CapitalColorLevelEncoder
	case "color":
		*e = LowercaseColorLevelEncoder
	default:
		*e = LowercaseLevelEncoder
	}
	return nil
}


type TimeEncoder func(time.Time, PrimitiveArrayEncoder)



func EpochTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) {
	nanos := t.UnixNano()
	sec := float64(nanos) / float64(time.Second)
	enc.AppendFloat64(sec)
}



func EpochMillisTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) {
	nanos := t.UnixNano()
	millis := float64(nanos) / float64(time.Millisecond)
	enc.AppendFloat64(millis)
}



func EpochNanosTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano())
}



func ISO8601TimeEncoder(t time.Time, enc PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
}




func (e *TimeEncoder) UnmarshalText(text []byte) error {
	switch string(text) {
	case "iso8601", "ISO8601":
		*e = ISO8601TimeEncoder
	case "millis":
		*e = EpochMillisTimeEncoder
	case "nanos":
		*e = EpochNanosTimeEncoder
	default:
		*e = EpochTimeEncoder
	}
	return nil
}


type DurationEncoder func(time.Duration, PrimitiveArrayEncoder)


func SecondsDurationEncoder(d time.Duration, enc PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Second))
}



func NanosDurationEncoder(d time.Duration, enc PrimitiveArrayEncoder) {
	enc.AppendInt64(int64(d))
}



func StringDurationEncoder(d time.Duration, enc PrimitiveArrayEncoder) {
	enc.AppendString(d.String())
}




func (e *DurationEncoder) UnmarshalText(text []byte) error {
	switch string(text) {
	case "string":
		*e = StringDurationEncoder
	case "nanos":
		*e = NanosDurationEncoder
	default:
		*e = SecondsDurationEncoder
	}
	return nil
}


type CallerEncoder func(EntryCaller, PrimitiveArrayEncoder)



func FullCallerEncoder(caller EntryCaller, enc PrimitiveArrayEncoder) {
	
	enc.AppendString(caller.String())
}



func ShortCallerEncoder(caller EntryCaller, enc PrimitiveArrayEncoder) {
	
	enc.AppendString(caller.TrimmedPath())
}



func (e *CallerEncoder) UnmarshalText(text []byte) error {
	switch string(text) {
	case "full":
		*e = FullCallerEncoder
	default:
		*e = ShortCallerEncoder
	}
	return nil
}



type NameEncoder func(string, PrimitiveArrayEncoder)


func FullNameEncoder(loggerName string, enc PrimitiveArrayEncoder) {
	enc.AppendString(loggerName)
}



func (e *NameEncoder) UnmarshalText(text []byte) error {
	switch string(text) {
	case "full":
		*e = FullNameEncoder
	default:
		*e = FullNameEncoder
	}
	return nil
}



type EncoderConfig struct {
	
	
	MessageKey    string `json:"messageKey" yaml:"messageKey"`
	LevelKey      string `json:"levelKey" yaml:"levelKey"`
	TimeKey       string `json:"timeKey" yaml:"timeKey"`
	NameKey       string `json:"nameKey" yaml:"nameKey"`
	CallerKey     string `json:"callerKey" yaml:"callerKey"`
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
	LineEnding    string `json:"lineEnding" yaml:"lineEnding"`
	
	
	
	EncodeLevel    LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime     TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller   CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
	
	
	EncodeName NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`
}




type ObjectEncoder interface {
	
	AddArray(key string, marshaler ArrayMarshaler) error
	AddObject(key string, marshaler ObjectMarshaler) error

	
	AddBinary(key string, value []byte)     
	AddByteString(key string, value []byte) 
	AddBool(key string, value bool)
	AddComplex128(key string, value complex128)
	AddComplex64(key string, value complex64)
	AddDuration(key string, value time.Duration)
	AddFloat64(key string, value float64)
	AddFloat32(key string, value float32)
	AddInt(key string, value int)
	AddInt64(key string, value int64)
	AddInt32(key string, value int32)
	AddInt16(key string, value int16)
	AddInt8(key string, value int8)
	AddString(key, value string)
	AddTime(key string, value time.Time)
	AddUint(key string, value uint)
	AddUint64(key string, value uint64)
	AddUint32(key string, value uint32)
	AddUint16(key string, value uint16)
	AddUint8(key string, value uint8)
	AddUintptr(key string, value uintptr)

	
	
	AddReflected(key string, value interface{}) error
	
	
	
	OpenNamespace(key string)
}





type ArrayEncoder interface {
	
	PrimitiveArrayEncoder

	
	AppendDuration(time.Duration)
	AppendTime(time.Time)

	
	AppendArray(ArrayMarshaler) error
	AppendObject(ObjectMarshaler) error

	
	
	AppendReflected(value interface{}) error
}




type PrimitiveArrayEncoder interface {
	
	AppendBool(bool)
	AppendByteString([]byte) 
	AppendComplex128(complex128)
	AppendComplex64(complex64)
	AppendFloat64(float64)
	AppendFloat32(float32)
	AppendInt(int)
	AppendInt64(int64)
	AppendInt32(int32)
	AppendInt16(int16)
	AppendInt8(int8)
	AppendString(string)
	AppendUint(uint)
	AppendUint64(uint64)
	AppendUint32(uint32)
	AppendUint16(uint16)
	AppendUint8(uint8)
	AppendUintptr(uintptr)
}









type Encoder interface {
	ObjectEncoder

	
	
	Clone() Encoder

	
	
	EncodeEntry(Entry, []Field) (*buffer.Buffer, error)
}
