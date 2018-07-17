



package snappy

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	
	ErrCorrupt = errors.New("snappy: corrupt input")
	
	ErrTooLarge = errors.New("snappy: decoded block is too large")
	
	ErrUnsupported = errors.New("snappy: unsupported input")

	errUnsupportedLiteralLength = errors.New("snappy: unsupported literal length")
)


func DecodedLen(src []byte) (int, error) {
	v, _, err := decodedLen(src)
	return v, err
}



func decodedLen(src []byte) (blockLen, headerLen int, err error) {
	v, n := binary.Uvarint(src)
	if n <= 0 || v > 0xffffffff {
		return 0, 0, ErrCorrupt
	}

	const wordSize = 32 << (^uint(0) >> 32 & 1)
	if wordSize == 32 && v > 0x7fffffff {
		return 0, 0, ErrTooLarge
	}
	return int(v), n, nil
}

const (
	decodeErrCodeCorrupt                  = 1
	decodeErrCodeUnsupportedLiteralLength = 2
)






func Decode(dst, src []byte) ([]byte, error) {
	dLen, s, err := decodedLen(src)
	if err != nil {
		return nil, err
	}
	if dLen <= len(dst) {
		dst = dst[:dLen]
	} else {
		dst = make([]byte, dLen)
	}
	switch decode(dst, src[s:]) {
	case 0:
		return dst, nil
	case decodeErrCodeUnsupportedLiteralLength:
		return nil, errUnsupportedLiteralLength
	}
	return nil, ErrCorrupt
}




func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:       r,
		decoded: make([]byte, maxBlockSize),
		buf:     make([]byte, maxEncodedLenOfMaxBlockSize+checksumSize),
	}
}


type Reader struct {
	r       io.Reader
	err     error
	decoded []byte
	buf     []byte
	
	i, j       int
	readHeader bool
}




func (r *Reader) Reset(reader io.Reader) {
	r.r = reader
	r.err = nil
	r.i = 0
	r.j = 0
	r.readHeader = false
}

func (r *Reader) readFull(p []byte, allowEOF bool) (ok bool) {
	if _, r.err = io.ReadFull(r.r, p); r.err != nil {
		if r.err == io.ErrUnexpectedEOF || (r.err == io.EOF && !allowEOF) {
			r.err = ErrCorrupt
		}
		return false
	}
	return true
}


func (r *Reader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	for {
		if r.i < r.j {
			n := copy(p, r.decoded[r.i:r.j])
			r.i += n
			return n, nil
		}
		if !r.readFull(r.buf[:4], true) {
			return 0, r.err
		}
		chunkType := r.buf[0]
		if !r.readHeader {
			if chunkType != chunkTypeStreamIdentifier {
				r.err = ErrCorrupt
				return 0, r.err
			}
			r.readHeader = true
		}
		chunkLen := int(r.buf[1]) | int(r.buf[2])<<8 | int(r.buf[3])<<16
		if chunkLen > len(r.buf) {
			r.err = ErrUnsupported
			return 0, r.err
		}

		
		
		switch chunkType {
		case chunkTypeCompressedData:
			
			if chunkLen < checksumSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			buf := r.buf[:chunkLen]
			if !r.readFull(buf, false) {
				return 0, r.err
			}
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			buf = buf[checksumSize:]

			n, err := DecodedLen(buf)
			if err != nil {
				r.err = err
				return 0, r.err
			}
			if n > len(r.decoded) {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if _, err := Decode(r.decoded, buf); err != nil {
				r.err = err
				return 0, r.err
			}
			if crc(r.decoded[:n]) != checksum {
				r.err = ErrCorrupt
				return 0, r.err
			}
			r.i, r.j = 0, n
			continue

		case chunkTypeUncompressedData:
			
			if chunkLen < checksumSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			buf := r.buf[:checksumSize]
			if !r.readFull(buf, false) {
				return 0, r.err
			}
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			
			n := chunkLen - checksumSize
			if n > len(r.decoded) {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if !r.readFull(r.decoded[:n], false) {
				return 0, r.err
			}
			if crc(r.decoded[:n]) != checksum {
				r.err = ErrCorrupt
				return 0, r.err
			}
			r.i, r.j = 0, n
			continue

		case chunkTypeStreamIdentifier:
			
			if chunkLen != len(magicBody) {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if !r.readFull(r.buf[:len(magicBody)], false) {
				return 0, r.err
			}
			for i := 0; i < len(magicBody); i++ {
				if r.buf[i] != magicBody[i] {
					r.err = ErrCorrupt
					return 0, r.err
				}
			}
			continue
		}

		if chunkType <= 0x7f {
			
			r.err = ErrUnsupported
			return 0, r.err
		}
		
		
		if !r.readFull(r.buf[:chunkLen], false) {
			return 0, r.err
		}
	}
}
