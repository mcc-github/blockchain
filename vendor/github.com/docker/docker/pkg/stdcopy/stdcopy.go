package stdcopy 

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
)



type StdType byte

const (
	
	Stdin StdType = iota
	
	Stdout
	
	Stderr
	
	
	Systemerr

	stdWriterPrefixLen = 8
	stdWriterFdIndex   = 0
	stdWriterSizeIndex = 4

	startingBufLen = 32*1024 + stdWriterPrefixLen + 1
)

var bufPool = &sync.Pool{New: func() interface{} { return bytes.NewBuffer(nil) }}


type stdWriter struct {
	io.Writer
	prefix byte
}





func (w *stdWriter) Write(p []byte) (n int, err error) {
	if w == nil || w.Writer == nil {
		return 0, errors.New("Writer not instantiated")
	}
	if p == nil {
		return 0, nil
	}

	header := [stdWriterPrefixLen]byte{stdWriterFdIndex: w.prefix}
	binary.BigEndian.PutUint32(header[stdWriterSizeIndex:], uint32(len(p)))
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Write(header[:])
	buf.Write(p)

	n, err = w.Writer.Write(buf.Bytes())
	n -= stdWriterPrefixLen
	if n < 0 {
		n = 0
	}

	buf.Reset()
	bufPool.Put(buf)
	return
}







func NewStdWriter(w io.Writer, t StdType) io.Writer {
	return &stdWriter{
		Writer: w,
		prefix: byte(t),
	}
}











func StdCopy(dstout, dsterr io.Writer, src io.Reader) (written int64, err error) {
	var (
		buf       = make([]byte, startingBufLen)
		bufLen    = len(buf)
		nr, nw    int
		er, ew    error
		out       io.Writer
		frameSize int
	)

	for {
		
		for nr < stdWriterPrefixLen {
			var nr2 int
			nr2, er = src.Read(buf[nr:])
			nr += nr2
			if er == io.EOF {
				if nr < stdWriterPrefixLen {
					return written, nil
				}
				break
			}
			if er != nil {
				return 0, er
			}
		}

		stream := StdType(buf[stdWriterFdIndex])
		
		switch stream {
		case Stdin:
			fallthrough
		case Stdout:
			
			out = dstout
		case Stderr:
			
			out = dsterr
		case Systemerr:
			
			
			
			out = nil
		default:
			return 0, fmt.Errorf("Unrecognized input header: %d", buf[stdWriterFdIndex])
		}

		
		frameSize = int(binary.BigEndian.Uint32(buf[stdWriterSizeIndex : stdWriterSizeIndex+4]))

		
		
		if frameSize+stdWriterPrefixLen > bufLen {
			buf = append(buf, make([]byte, frameSize+stdWriterPrefixLen-bufLen+1)...)
			bufLen = len(buf)
		}

		
		for nr < frameSize+stdWriterPrefixLen {
			var nr2 int
			nr2, er = src.Read(buf[nr:])
			nr += nr2
			if er == io.EOF {
				if nr < frameSize+stdWriterPrefixLen {
					return written, nil
				}
				break
			}
			if er != nil {
				return 0, er
			}
		}

		
		
		if stream == Systemerr {
			return written, fmt.Errorf("error from daemon in stream: %s", string(buf[stdWriterPrefixLen:frameSize+stdWriterPrefixLen]))
		}

		
		nw, ew = out.Write(buf[stdWriterPrefixLen : frameSize+stdWriterPrefixLen])
		if ew != nil {
			return 0, ew
		}

		
		if nw != frameSize {
			return 0, io.ErrShortWrite
		}
		written += int64(nw)

		
		copy(buf, buf[frameSize+stdWriterPrefixLen:])
		
		nr -= frameSize + stdWriterPrefixLen
	}
}
