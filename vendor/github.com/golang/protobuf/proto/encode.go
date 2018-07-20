






























package proto



import (
	"errors"
	"fmt"
	"reflect"
)










type RequiredNotSetError struct {
	field string
}

func (e *RequiredNotSetError) Error() string {
	return fmt.Sprintf("proto: required field %q not set", e.field)
}

var (
	
	
	errRepeatedHasNil = errors.New("proto: repeated field has nil element")

	
	
	errOneofHasNil = errors.New("proto: oneof field has nil value")

	
	ErrNil = errors.New("proto: Marshal called with nil")

	
	
	ErrTooLarge = errors.New("proto: message encodes to over 2 GB")
)





const maxVarintBytes = 10 







func EncodeVarint(x uint64) []byte {
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}





func (p *Buffer) EncodeVarint(x uint64) error {
	for x >= 1<<7 {
		p.buf = append(p.buf, uint8(x&0x7f|0x80))
		x >>= 7
	}
	p.buf = append(p.buf, uint8(x))
	return nil
}


func SizeVarint(x uint64) int {
	switch {
	case x < 1<<7:
		return 1
	case x < 1<<14:
		return 2
	case x < 1<<21:
		return 3
	case x < 1<<28:
		return 4
	case x < 1<<35:
		return 5
	case x < 1<<42:
		return 6
	case x < 1<<49:
		return 7
	case x < 1<<56:
		return 8
	case x < 1<<63:
		return 9
	}
	return 10
}




func (p *Buffer) EncodeFixed64(x uint64) error {
	p.buf = append(p.buf,
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24),
		uint8(x>>32),
		uint8(x>>40),
		uint8(x>>48),
		uint8(x>>56))
	return nil
}




func (p *Buffer) EncodeFixed32(x uint64) error {
	p.buf = append(p.buf,
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24))
	return nil
}




func (p *Buffer) EncodeZigzag64(x uint64) error {
	
	return p.EncodeVarint(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}




func (p *Buffer) EncodeZigzag32(x uint64) error {
	
	return p.EncodeVarint(uint64((uint32(x) << 1) ^ uint32((int32(x) >> 31))))
}




func (p *Buffer) EncodeRawBytes(b []byte) error {
	p.EncodeVarint(uint64(len(b)))
	p.buf = append(p.buf, b...)
	return nil
}



func (p *Buffer) EncodeStringBytes(s string) error {
	p.EncodeVarint(uint64(len(s)))
	p.buf = append(p.buf, s...)
	return nil
}


type Marshaler interface {
	Marshal() ([]byte, error)
}



func (p *Buffer) EncodeMessage(pb Message) error {
	siz := Size(pb)
	p.EncodeVarint(uint64(siz))
	return p.Marshal(pb)
}


func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}
