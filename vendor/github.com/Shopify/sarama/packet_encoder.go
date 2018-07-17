package sarama

import "github.com/rcrowley/go-metrics"




type packetEncoder interface {
	
	putInt8(in int8)
	putInt16(in int16)
	putInt32(in int32)
	putInt64(in int64)
	putVarint(in int64)
	putArrayLength(in int) error
	putBool(in bool)

	
	putBytes(in []byte) error
	putVarintBytes(in []byte) error
	putRawBytes(in []byte) error
	putString(in string) error
	putNullableString(in *string) error
	putStringArray(in []string) error
	putInt32Array(in []int32) error
	putInt64Array(in []int64) error

	
	offset() int

	
	push(in pushEncoder)
	pop() error

	
	metricRegistry() metrics.Registry
}





type pushEncoder interface {
	
	saveOffset(in int)

	
	reserveLength() int

	
	
	
	run(curOffset int, buf []byte) error
}




type dynamicPushEncoder interface {
	pushEncoder

	
	
	adjustLength(currOffset int) int
}
