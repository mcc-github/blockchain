package lz4

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"runtime"
)


type Writer struct {
	Header
	dst      io.Writer
	checksum hash.Hash32 
	data     []byte      
	window   []byte      

	zbCompressBuf []byte 
	writeSizeBuf  []byte 
}





func NewWriter(dst io.Writer) *Writer {
	return &Writer{
		dst:      dst,
		checksum: hashPool.Get(),
		Header: Header{
			BlockMaxSize: 4 << 20,
		},
		writeSizeBuf: make([]byte, 4),
	}
}


func (z *Writer) writeHeader() error {
	
	if z.Header.BlockMaxSize == 0 {
		z.Header.BlockMaxSize = 4 << 20
	}
	
	bSize, ok := bsMapValue[z.Header.BlockMaxSize]
	if !ok {
		return fmt.Errorf("lz4: invalid block max size: %d", z.Header.BlockMaxSize)
	}

	
	
	var buf [19]byte

	
	binary.LittleEndian.PutUint32(buf[0:], frameMagic)
	flg := byte(Version << 6)
	if !z.Header.BlockDependency {
		flg |= 1 << 5
	}
	if z.Header.BlockChecksum {
		flg |= 1 << 4
	}
	if z.Header.Size > 0 {
		flg |= 1 << 3
	}
	if !z.Header.NoChecksum {
		flg |= 1 << 2
	}
	
	
	
	buf[4] = flg
	buf[5] = bSize << 4

	
	n := 6
	
	if z.Header.Size > 0 {
		binary.LittleEndian.PutUint64(buf[n:], z.Header.Size)
		n += 8
	}
	
	
	
	

	
	z.checksum.Write(buf[4:n])
	buf[n] = byte(z.checksum.Sum32() >> 8 & 0xFF)
	z.checksum.Reset()

	
	if _, err := z.dst.Write(buf[0 : n+1]); err != nil {
		return err
	}
	z.Header.done = true

	
	z.zbCompressBuf = make([]byte, winSize+z.BlockMaxSize)

	return nil
}









func (z *Writer) Write(buf []byte) (n int, err error) {
	if !z.Header.done {
		if err = z.writeHeader(); err != nil {
			return
		}
	}

	if len(buf) == 0 {
		return
	}

	if !z.NoChecksum {
		z.checksum.Write(buf)
	}

	
	
	bl := 0
	if z.BlockDependency && len(z.window) == 0 {
		bl = len(z.data)
		z.data = append(z.data, buf...)
		if len(z.data) < winSize {
			return len(buf), nil
		}
		buf = z.data
		z.data = nil
	}

	
	
	var (
		zb       block
		wbuf     = buf
		zn       = len(wbuf) / z.BlockMaxSize
		zi       = 0
		leftover = len(buf) % z.BlockMaxSize
	)

loop:
	for zi < zn {
		if z.BlockDependency {
			if zi == 0 {
				
				zb.data = append(z.window, wbuf[:z.BlockMaxSize]...)
				zb.offset = len(z.window)
				wbuf = wbuf[z.BlockMaxSize-winSize:]
			} else {
				
				zb.data = wbuf[:z.BlockMaxSize+winSize]
				zb.offset = winSize
				wbuf = wbuf[z.BlockMaxSize:]
			}
		} else {
			zb.data = wbuf[:z.BlockMaxSize]
			wbuf = wbuf[z.BlockMaxSize:]
		}

		goto write
	}

	
	if leftover > 0 {
		zb = block{data: wbuf}
		if z.BlockDependency {
			if zn == 0 {
				zb.data = append(z.window, zb.data...)
				zb.offset = len(z.window)
			} else {
				zb.offset = winSize
			}
		}

		leftover = 0
		goto write
	}

	if z.BlockDependency {
		if len(z.window) == 0 {
			z.window = make([]byte, winSize)
		}
		
		if len(buf) >= winSize {
			copy(z.window, buf[len(buf)-winSize:])
		} else {
			copy(z.window, z.window[len(buf):])
			copy(z.window[len(buf)+1:], buf)
		}
	}

	return

write:
	zb = z.compressBlock(zb)
	_, err = z.writeBlock(zb)

	written := len(zb.data)
	if bl > 0 {
		if written >= bl {
			written -= bl
			bl = 0
		} else {
			bl -= written
			written = 0
		}
	}

	n += written
	
	if z.BlockDependency {
		if zi == 0 {
			n -= len(z.window)
		} else {
			n -= winSize
		}
	}
	if err != nil {
		return
	}
	zi++
	goto loop
}


func (z *Writer) compressBlock(zb block) block {
	
	var (
		n    int
		err  error
		zbuf = z.zbCompressBuf
	)
	if z.HighCompression {
		n, err = CompressBlockHC(zb.data, zbuf, zb.offset)
	} else {
		n, err = CompressBlock(zb.data, zbuf, zb.offset)
	}

	
	if err == nil && n > 0 && len(zb.zdata) < len(zb.data) {
		zb.compressed = true
		zb.zdata = zbuf[:n]
	} else {
		zb.compressed = false
		zb.zdata = zb.data[zb.offset:]
	}

	if z.BlockChecksum {
		xxh := hashPool.Get()
		xxh.Write(zb.zdata)
		zb.checksum = xxh.Sum32()
		hashPool.Put(xxh)
	}

	return zb
}


func (z *Writer) writeBlock(zb block) (int, error) {
	bLen := uint32(len(zb.zdata))
	if !zb.compressed {
		bLen |= 1 << 31
	}

	n := 0

	binary.LittleEndian.PutUint32(z.writeSizeBuf, bLen)
	n, err := z.dst.Write(z.writeSizeBuf)
	if err != nil {
		return n, err
	}

	m, err := z.dst.Write(zb.zdata)
	n += m
	if err != nil {
		return n, err
	}

	if z.BlockChecksum {
		binary.LittleEndian.PutUint32(z.writeSizeBuf, zb.checksum)
		m, err := z.dst.Write(z.writeSizeBuf)
		n += m

		if err != nil {
			return n, err
		}
	}

	return n, nil
}







func (z *Writer) Flush() error {
	if len(z.data) == 0 {
		return nil
	}

	zb := z.compressBlock(block{data: z.data})
	if _, err := z.writeBlock(zb); err != nil {
		return err
	}
	return nil
}


func (z *Writer) Close() error {
	if !z.Header.done {
		if err := z.writeHeader(); err != nil {
			return err
		}
	}

	
	if z.BlockDependency && len(z.data) > 0 {
		zb := block{data: z.data}
		if _, err := z.writeBlock(z.compressBlock(zb)); err != nil {
			return err
		}
	}

	if err := binary.Write(z.dst, binary.LittleEndian, uint32(0)); err != nil {
		return err
	}
	if !z.NoChecksum {
		if err := binary.Write(z.dst, binary.LittleEndian, z.checksum.Sum32()); err != nil {
			return err
		}
	}
	return nil
}




func (z *Writer) Reset(w io.Writer) {
	z.Header = Header{}
	z.dst = w
	z.checksum.Reset()
	z.data = nil
	z.window = nil
}




func (z *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	cpus := runtime.GOMAXPROCS(0)
	buf := make([]byte, cpus*z.BlockMaxSize)
	for {
		m, er := io.ReadFull(r, buf)
		n += int64(m)
		if er == nil || er == io.ErrUnexpectedEOF || er == io.EOF {
			if _, err = z.Write(buf[:m]); err != nil {
				return
			}
			if er == nil {
				continue
			}
			return
		}
		return n, er
	}
}
