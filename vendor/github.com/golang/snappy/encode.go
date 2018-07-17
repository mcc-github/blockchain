



package snappy

import (
	"encoding/binary"
	"errors"
	"io"
)






func Encode(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if len(dst) < n {
		dst = make([]byte, n)
	}

	
	d := binary.PutUvarint(dst, uint64(len(src)))

	for len(src) > 0 {
		p := src
		src = nil
		if len(p) > maxBlockSize {
			p, src = p[:maxBlockSize], p[maxBlockSize:]
		}
		if len(p) < minNonLiteralBlockSize {
			d += emitLiteral(dst[d:], p)
		} else {
			d += encodeBlock(dst[d:], p)
		}
	}
	return dst[:d]
}









const inputMargin = 16 - 1


















const minNonLiteralBlockSize = 1 + 1 + inputMargin





func MaxEncodedLen(srcLen int) int {
	n := uint64(srcLen)
	if n > 0xffffffff {
		return -1
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	n = 32 + n + n/6
	if n > 0xffffffff {
		return -1
	}
	return int(n)
}

var errClosed = errors.New("snappy: Writer is closed")










func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:    w,
		obuf: make([]byte, obufLen),
	}
}








func NewBufferedWriter(w io.Writer) *Writer {
	return &Writer{
		w:    w,
		ibuf: make([]byte, 0, maxBlockSize),
		obuf: make([]byte, obufLen),
	}
}


type Writer struct {
	w   io.Writer
	err error

	
	
	
	
	
	ibuf []byte

	
	obuf []byte

	
	wroteStreamHeader bool
}



func (w *Writer) Reset(writer io.Writer) {
	w.w = writer
	w.err = nil
	if w.ibuf != nil {
		w.ibuf = w.ibuf[:0]
	}
	w.wroteStreamHeader = false
}


func (w *Writer) Write(p []byte) (nRet int, errRet error) {
	if w.ibuf == nil {
		
		
		
		
		return w.write(p)
	}

	
	

	for len(p) > (cap(w.ibuf)-len(w.ibuf)) && w.err == nil {
		var n int
		if len(w.ibuf) == 0 {
			
			
			n, _ = w.write(p)
		} else {
			n = copy(w.ibuf[len(w.ibuf):cap(w.ibuf)], p)
			w.ibuf = w.ibuf[:len(w.ibuf)+n]
			w.Flush()
		}
		nRet += n
		p = p[n:]
	}
	if w.err != nil {
		return nRet, w.err
	}
	n := copy(w.ibuf[len(w.ibuf):cap(w.ibuf)], p)
	w.ibuf = w.ibuf[:len(w.ibuf)+n]
	nRet += n
	return nRet, nil
}

func (w *Writer) write(p []byte) (nRet int, errRet error) {
	if w.err != nil {
		return 0, w.err
	}
	for len(p) > 0 {
		obufStart := len(magicChunk)
		if !w.wroteStreamHeader {
			w.wroteStreamHeader = true
			copy(w.obuf, magicChunk)
			obufStart = 0
		}

		var uncompressed []byte
		if len(p) > maxBlockSize {
			uncompressed, p = p[:maxBlockSize], p[maxBlockSize:]
		} else {
			uncompressed, p = p, nil
		}
		checksum := crc(uncompressed)

		
		
		compressed := Encode(w.obuf[obufHeaderLen:], uncompressed)
		chunkType := uint8(chunkTypeCompressedData)
		chunkLen := 4 + len(compressed)
		obufEnd := obufHeaderLen + len(compressed)
		if len(compressed) >= len(uncompressed)-len(uncompressed)/8 {
			chunkType = chunkTypeUncompressedData
			chunkLen = 4 + len(uncompressed)
			obufEnd = obufHeaderLen
		}

		
		w.obuf[len(magicChunk)+0] = chunkType
		w.obuf[len(magicChunk)+1] = uint8(chunkLen >> 0)
		w.obuf[len(magicChunk)+2] = uint8(chunkLen >> 8)
		w.obuf[len(magicChunk)+3] = uint8(chunkLen >> 16)
		w.obuf[len(magicChunk)+4] = uint8(checksum >> 0)
		w.obuf[len(magicChunk)+5] = uint8(checksum >> 8)
		w.obuf[len(magicChunk)+6] = uint8(checksum >> 16)
		w.obuf[len(magicChunk)+7] = uint8(checksum >> 24)

		if _, err := w.w.Write(w.obuf[obufStart:obufEnd]); err != nil {
			w.err = err
			return nRet, err
		}
		if chunkType == chunkTypeUncompressedData {
			if _, err := w.w.Write(uncompressed); err != nil {
				w.err = err
				return nRet, err
			}
		}
		nRet += len(uncompressed)
	}
	return nRet, nil
}


func (w *Writer) Flush() error {
	if w.err != nil {
		return w.err
	}
	if len(w.ibuf) == 0 {
		return nil
	}
	w.write(w.ibuf)
	w.ibuf = w.ibuf[:0]
	return w.err
}


func (w *Writer) Close() error {
	w.Flush()
	ret := w.err
	if w.err == nil {
		w.err = errClosed
	}
	return ret
}
