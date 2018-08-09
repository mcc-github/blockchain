



















package zapcore




type ObjectMarshaler interface {
	MarshalLogObject(ObjectEncoder) error
}



type ObjectMarshalerFunc func(ObjectEncoder) error


func (f ObjectMarshalerFunc) MarshalLogObject(enc ObjectEncoder) error {
	return f(enc)
}




type ArrayMarshaler interface {
	MarshalLogArray(ArrayEncoder) error
}



type ArrayMarshalerFunc func(ArrayEncoder) error


func (f ArrayMarshalerFunc) MarshalLogArray(enc ArrayEncoder) error {
	return f(enc)
}
