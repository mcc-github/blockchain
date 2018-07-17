












































































package journal

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)


const (
	fullChunkType   = 1
	firstChunkType  = 2
	middleChunkType = 3
	lastChunkType   = 4
)

const (
	blockSize  = 32 * 1024
	headerSize = 7
)

type flusher interface {
	Flush() error
}


type ErrCorrupted struct {
	Size   int
	Reason string
}

func (e *ErrCorrupted) Error() string {
	return fmt.Sprintf("leveldb/journal: block/chunk corrupted: %s (%d bytes)", e.Reason, e.Size)
}



type Dropper interface {
	Drop(err error)
}


type Reader struct {
	
	r io.Reader
	
	dropper Dropper
	
	strict bool
	
	checksum bool
	
	seq int
	
	
	i, j int
	
	
	n int
	
	last bool
	
	err error
	
	buf [blockSize]byte
}




func NewReader(r io.Reader, dropper Dropper, strict, checksum bool) *Reader {
	return &Reader{
		r:        r,
		dropper:  dropper,
		strict:   strict,
		checksum: checksum,
		last:     true,
	}
}

var errSkip = errors.New("leveldb/journal: skipped")

func (r *Reader) corrupt(n int, reason string, skip bool) error {
	if r.dropper != nil {
		r.dropper.Drop(&ErrCorrupted{n, reason})
	}
	if r.strict && !skip {
		r.err = errors.NewErrCorrupted(storage.FileDesc{}, &ErrCorrupted{n, reason})
		return r.err
	}
	return errSkip
}



func (r *Reader) nextChunk(first bool) error {
	for {
		if r.j+headerSize <= r.n {
			checksum := binary.LittleEndian.Uint32(r.buf[r.j+0 : r.j+4])
			length := binary.LittleEndian.Uint16(r.buf[r.j+4 : r.j+6])
			chunkType := r.buf[r.j+6]
			unprocBlock := r.n - r.j
			if checksum == 0 && length == 0 && chunkType == 0 {
				
				r.i = r.n
				r.j = r.n
				return r.corrupt(unprocBlock, "zero header", false)
			}
			if chunkType < fullChunkType || chunkType > lastChunkType {
				
				r.i = r.n
				r.j = r.n
				return r.corrupt(unprocBlock, fmt.Sprintf("invalid chunk type %#x", chunkType), false)
			}
			r.i = r.j + headerSize
			r.j = r.j + headerSize + int(length)
			if r.j > r.n {
				
				r.i = r.n
				r.j = r.n
				return r.corrupt(unprocBlock, "chunk length overflows block", false)
			} else if r.checksum && checksum != util.NewCRC(r.buf[r.i-1:r.j]).Value() {
				
				r.i = r.n
				r.j = r.n
				return r.corrupt(unprocBlock, "checksum mismatch", false)
			}
			if first && chunkType != fullChunkType && chunkType != firstChunkType {
				chunkLength := (r.j - r.i) + headerSize
				r.i = r.j
				
				return r.corrupt(chunkLength, "orphan chunk", true)
			}
			r.last = chunkType == fullChunkType || chunkType == lastChunkType
			return nil
		}

		
		if r.n < blockSize && r.n > 0 {
			if !first {
				return r.corrupt(0, "missing chunk part", false)
			}
			r.err = io.EOF
			return r.err
		}

		
		n, err := io.ReadFull(r.r, r.buf[:])
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
		if n == 0 {
			if !first {
				return r.corrupt(0, "missing chunk part", false)
			}
			r.err = io.EOF
			return r.err
		}
		r.i, r.j, r.n = 0, 0, n
	}
}





func (r *Reader) Next() (io.Reader, error) {
	r.seq++
	if r.err != nil {
		return nil, r.err
	}
	r.i = r.j
	for {
		if err := r.nextChunk(true); err == nil {
			break
		} else if err != errSkip {
			return nil, err
		}
	}
	return &singleReader{r, r.seq, nil}, nil
}



func (r *Reader) Reset(reader io.Reader, dropper Dropper, strict, checksum bool) error {
	r.seq++
	err := r.err
	r.r = reader
	r.dropper = dropper
	r.strict = strict
	r.checksum = checksum
	r.i = 0
	r.j = 0
	r.n = 0
	r.last = true
	r.err = nil
	return err
}

type singleReader struct {
	r   *Reader
	seq int
	err error
}

func (x *singleReader) Read(p []byte) (int, error) {
	r := x.r
	if r.seq != x.seq {
		return 0, errors.New("leveldb/journal: stale reader")
	}
	if x.err != nil {
		return 0, x.err
	}
	if r.err != nil {
		return 0, r.err
	}
	for r.i == r.j {
		if r.last {
			return 0, io.EOF
		}
		x.err = r.nextChunk(false)
		if x.err != nil {
			if x.err == errSkip {
				x.err = io.ErrUnexpectedEOF
			}
			return 0, x.err
		}
	}
	n := copy(p, r.buf[r.i:r.j])
	r.i += n
	return n, nil
}

func (x *singleReader) ReadByte() (byte, error) {
	r := x.r
	if r.seq != x.seq {
		return 0, errors.New("leveldb/journal: stale reader")
	}
	if x.err != nil {
		return 0, x.err
	}
	if r.err != nil {
		return 0, r.err
	}
	for r.i == r.j {
		if r.last {
			return 0, io.EOF
		}
		x.err = r.nextChunk(false)
		if x.err != nil {
			if x.err == errSkip {
				x.err = io.ErrUnexpectedEOF
			}
			return 0, x.err
		}
	}
	c := r.buf[r.i]
	r.i++
	return c, nil
}


type Writer struct {
	
	w io.Writer
	
	seq int
	
	f flusher
	
	
	i, j int
	
	
	written int
	
	first bool
	
	pending bool
	
	err error
	
	buf [blockSize]byte
}


func NewWriter(w io.Writer) *Writer {
	f, _ := w.(flusher)
	return &Writer{
		w: w,
		f: f,
	}
}


func (w *Writer) fillHeader(last bool) {
	if w.i+headerSize > w.j || w.j > blockSize {
		panic("leveldb/journal: bad writer state")
	}
	if last {
		if w.first {
			w.buf[w.i+6] = fullChunkType
		} else {
			w.buf[w.i+6] = lastChunkType
		}
	} else {
		if w.first {
			w.buf[w.i+6] = firstChunkType
		} else {
			w.buf[w.i+6] = middleChunkType
		}
	}
	binary.LittleEndian.PutUint32(w.buf[w.i+0:w.i+4], util.NewCRC(w.buf[w.i+6:w.j]).Value())
	binary.LittleEndian.PutUint16(w.buf[w.i+4:w.i+6], uint16(w.j-w.i-headerSize))
}



func (w *Writer) writeBlock() {
	_, w.err = w.w.Write(w.buf[w.written:])
	w.i = 0
	w.j = headerSize
	w.written = 0
}



func (w *Writer) writePending() {
	if w.err != nil {
		return
	}
	if w.pending {
		w.fillHeader(true)
		w.pending = false
	}
	_, w.err = w.w.Write(w.buf[w.written:w.j])
	w.written = w.j
}


func (w *Writer) Close() error {
	w.seq++
	w.writePending()
	if w.err != nil {
		return w.err
	}
	w.err = errors.New("leveldb/journal: closed Writer")
	return nil
}



func (w *Writer) Flush() error {
	w.seq++
	w.writePending()
	if w.err != nil {
		return w.err
	}
	if w.f != nil {
		w.err = w.f.Flush()
		return w.err
	}
	return nil
}



func (w *Writer) Reset(writer io.Writer) (err error) {
	w.seq++
	if w.err == nil {
		w.writePending()
		err = w.err
	}
	w.w = writer
	w.f, _ = writer.(flusher)
	w.i = 0
	w.j = 0
	w.written = 0
	w.first = false
	w.pending = false
	w.err = nil
	return
}



func (w *Writer) Next() (io.Writer, error) {
	w.seq++
	if w.err != nil {
		return nil, w.err
	}
	if w.pending {
		w.fillHeader(true)
	}
	w.i = w.j
	w.j = w.j + headerSize
	
	if w.j > blockSize {
		
		for k := w.i; k < blockSize; k++ {
			w.buf[k] = 0
		}
		w.writeBlock()
		if w.err != nil {
			return nil, w.err
		}
	}
	w.first = true
	w.pending = true
	return singleWriter{w, w.seq}, nil
}

type singleWriter struct {
	w   *Writer
	seq int
}

func (x singleWriter) Write(p []byte) (int, error) {
	w := x.w
	if w.seq != x.seq {
		return 0, errors.New("leveldb/journal: stale writer")
	}
	if w.err != nil {
		return 0, w.err
	}
	n0 := len(p)
	for len(p) > 0 {
		
		if w.j == blockSize {
			w.fillHeader(false)
			w.writeBlock()
			if w.err != nil {
				return 0, w.err
			}
			w.first = false
		}
		
		n := copy(w.buf[w.j:], p)
		w.j += n
		p = p[n:]
	}
	return n0, nil
}
