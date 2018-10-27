package log

import (
	"time"

	"github.com/go-stack/stack"
)




type Valuer func() interface{}



func bindValues(keyvals []interface{}) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v()
		}
	}
}



func containsValuer(keyvals []interface{}) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}






func Timestamp(t func() time.Time) Valuer {
	return func() interface{} { return t() }
}








func TimestampFormat(t func() time.Time, layout string) Valuer {
	return func() interface{} {
		return timeFormat{
			time:   t(),
			layout: layout,
		}
	}
}



type timeFormat struct {
	time   time.Time
	layout string
}

func (tf timeFormat) String() string {
	return tf.time.Format(tf.layout)
}


func (tf timeFormat) MarshalText() (text []byte, err error) {
	
	
	
	b := make([]byte, 0, len(tf.layout)+10)
	b = tf.time.AppendFormat(b, tf.layout)
	return b, nil
}



func Caller(depth int) Valuer {
	return func() interface{} { return stack.Caller(depth) }
}

var (
	
	
	DefaultTimestamp = TimestampFormat(time.Now, time.RFC3339Nano)

	
	
	DefaultTimestampUTC = TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		time.RFC3339Nano,
	)

	
	
	DefaultCaller = Caller(3)
)
