package sarama




type packetDecoder interface {
	
	getInt8() (int8, error)
	getInt16() (int16, error)
	getInt32() (int32, error)
	getInt64() (int64, error)
	getVarint() (int64, error)
	getArrayLength() (int, error)
	getBool() (bool, error)

	
	getBytes() ([]byte, error)
	getVarintBytes() ([]byte, error)
	getRawBytes(length int) ([]byte, error)
	getString() (string, error)
	getNullableString() (*string, error)
	getInt32Array() ([]int32, error)
	getInt64Array() ([]int64, error)
	getStringArray() ([]string, error)

	
	remaining() int
	getSubset(length int) (packetDecoder, error)
	peek(offset, length int) (packetDecoder, error) 

	
	push(in pushDecoder) error
	pop() error
}





type pushDecoder interface {
	
	saveOffset(in int)

	
	reserveLength() int

	
	
	
	check(curOffset int, buf []byte) error
}





type dynamicPushDecoder interface {
	pushDecoder
	decoder
}
