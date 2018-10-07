package snappy

import (
	"bytes"
	"encoding/binary"
	"errors"

	master "github.com/golang/snappy"
)

const (
	sizeOffset = 16
	sizeBytes  = 4
)

var (
	xerialHeader = []byte{130, 83, 78, 65, 80, 80, 89, 0}

	
	xerialVersionInfo = []byte{0, 0, 0, 1, 0, 0, 0, 1}

	
	
	ErrMalformed = errors.New("malformed xerial framing")
)

func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}


func Encode(src []byte) []byte {
	return master.Encode(nil, src)
}





func EncodeStream(dst, src []byte) []byte {
	if len(dst) == 0 {
		dst = append(dst, xerialHeader...)
		dst = append(dst, xerialVersionInfo...)
	}

	
	var (
		max = len(src)
		blockSize = 32 * 1024
		pos   = 0
		chunk []byte
	)

	for pos < max {
		newPos := min(pos + blockSize, max)
		chunk = master.Encode(chunk[:cap(chunk)], src[pos:newPos])

		
		
		origLen := len(dst)
		dst = append(dst, dst[0:4]...)
		binary.BigEndian.PutUint32(dst[origLen:], uint32(len(chunk)))

		
		dst = append(dst, chunk...)
		pos = newPos
	}
	return dst
}



func Decode(src []byte) ([]byte, error) {
	return DecodeInto(nil, src)
}






func DecodeInto(dst, src []byte) ([]byte, error) {
	var max = len(src)
	if max < len(xerialHeader) {
		return nil, ErrMalformed
	}

	if !bytes.Equal(src[:8], xerialHeader) {
		return master.Decode(dst[:cap(dst)], src)
	}

	if max < sizeOffset+sizeBytes {
		return nil, ErrMalformed
	}

	if dst == nil {
		dst = make([]byte, 0, len(src))
	}

	dst = dst[:0]
	var (
		pos   = sizeOffset
		chunk []byte
		err       error
	)

	for pos+sizeBytes <= max {
		size := int(binary.BigEndian.Uint32(src[pos : pos+sizeBytes]))
		pos += sizeBytes

		nextPos := pos + size
		
		
		
		if nextPos < pos || nextPos > max {
			return nil, ErrMalformed
		}

		chunk, err = master.Decode(chunk[:cap(chunk)], src[pos:nextPos])

		if err != nil {
			return nil, err
		}
		pos = nextPos
		dst = append(dst, chunk...)
	}
	return dst, nil
}
