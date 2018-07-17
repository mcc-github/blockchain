



package util




import (
	"bytes"
	"io"
)



type Buffer struct {
	buf       []byte   
	off       int      
	bootstrap [64]byte 
}





func (b *Buffer) Bytes() []byte { return b.buf[b.off:] }



func (b *Buffer) String() string {
	if b == nil {
		
		return "<nil>"
	}
	return string(b.buf[b.off:])
}



func (b *Buffer) Len() int { return len(b.buf) - b.off }



func (b *Buffer) Truncate(n int) {
	switch {
	case n < 0 || n > b.Len():
		panic("leveldb/util.Buffer: truncation out of range")
	case n == 0:
		
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+n]
}



func (b *Buffer) Reset() { b.Truncate(0) }




func (b *Buffer) grow(n int) int {
	m := b.Len()
	
	if m == 0 && b.off != 0 {
		b.Truncate(0)
	}
	if len(b.buf)+n > cap(b.buf) {
		var buf []byte
		if b.buf == nil && n <= len(b.bootstrap) {
			buf = b.bootstrap[0:]
		} else if m+n <= cap(b.buf)/2 {
			
			
			
			
			copy(b.buf[:], b.buf[b.off:])
			buf = b.buf[:m]
		} else {
			
			buf = makeSlice(2*cap(b.buf) + n)
			copy(buf, b.buf[b.off:])
		}
		b.buf = buf
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+m+n]
	return b.off + m
}




func (b *Buffer) Alloc(n int) []byte {
	if n < 0 {
		panic("leveldb/util.Buffer.Alloc: negative count")
	}
	m := b.grow(n)
	return b.buf[m:]
}






func (b *Buffer) Grow(n int) {
	if n < 0 {
		panic("leveldb/util.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[0:m]
}




func (b *Buffer) Write(p []byte) (n int, err error) {
	m := b.grow(len(p))
	return copy(b.buf[m:], p), nil
}





const MinRead = 512





func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	
	if b.off >= len(b.buf) {
		b.Truncate(0)
	}
	for {
		if free := cap(b.buf) - len(b.buf); free < MinRead {
			
			newBuf := b.buf
			if b.off+free < MinRead {
				
				
				newBuf = makeSlice(2*cap(b.buf) + MinRead)
			}
			copy(newBuf, b.buf[b.off:])
			b.buf = newBuf[:len(b.buf)-b.off]
			b.off = 0
		}
		m, e := r.Read(b.buf[len(b.buf):cap(b.buf)])
		b.buf = b.buf[0 : len(b.buf)+m]
		n += int64(m)
		if e == io.EOF {
			break
		}
		if e != nil {
			return n, e
		}
	}
	return n, nil 
}



func makeSlice(n int) []byte {
	
	defer func() {
		if recover() != nil {
			panic(bytes.ErrTooLarge)
		}
	}()
	return make([]byte, n)
}





func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	if b.off < len(b.buf) {
		nBytes := b.Len()
		m, e := w.Write(b.buf[b.off:])
		if m > nBytes {
			panic("leveldb/util.Buffer.WriteTo: invalid Write count")
		}
		b.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		
		
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	
	b.Truncate(0)
	return
}





func (b *Buffer) WriteByte(c byte) error {
	m := b.grow(1)
	b.buf[m] = c
	return nil
}





func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.off >= len(b.buf) {
		
		b.Truncate(0)
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.off:])
	b.off += n
	return
}





func (b *Buffer) Next(n int) []byte {
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	return data
}



func (b *Buffer) ReadByte() (c byte, err error) {
	if b.off >= len(b.buf) {
		
		b.Truncate(0)
		return 0, io.EOF
	}
	c = b.buf[b.off]
	b.off++
	return c, nil
}







func (b *Buffer) ReadBytes(delim byte) (line []byte, err error) {
	slice, err := b.readSlice(delim)
	
	
	line = append(line, slice...)
	return
}


func (b *Buffer) readSlice(delim byte) (line []byte, err error) {
	i := bytes.IndexByte(b.buf[b.off:], delim)
	end := b.off + i + 1
	if i < 0 {
		end = len(b.buf)
		err = io.EOF
	}
	line = b.buf[b.off:end]
	b.off = end
	return line, err
}








func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }
