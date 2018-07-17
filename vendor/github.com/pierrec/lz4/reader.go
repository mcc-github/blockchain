package lz4

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"runtime"
	"sync"
	"sync/atomic"
)



var ErrInvalid = errors.New("invalid lz4 data")



var errEndOfBlock = errors.New("end of block")




type Reader struct {
	Pos int64 
	Header
	src      io.Reader
	checksum hash.Hash32    
	wg       sync.WaitGroup 
	data     []byte         
	window   []byte         
}



func NewReader(src io.Reader) *Reader {
	return &Reader{
		src:      src,
		checksum: hashPool.Get(),
	}
}




func (z *Reader) readHeader(first bool) error {
	defer z.checksum.Reset()

	for {
		var magic uint32
		if err := binary.Read(z.src, binary.LittleEndian, &magic); err != nil {
			if !first && err == io.ErrUnexpectedEOF {
				return io.EOF
			}
			return err
		}
		z.Pos += 4
		if magic>>8 == frameSkipMagic>>8 {
			var skipSize uint32
			if err := binary.Read(z.src, binary.LittleEndian, &skipSize); err != nil {
				return err
			}
			z.Pos += 4
			m, err := io.CopyN(ioutil.Discard, z.src, int64(skipSize))
			z.Pos += m
			if err != nil {
				return err
			}
			continue
		}
		if magic != frameMagic {
			return ErrInvalid
		}
		break
	}

	
	var buf [8]byte
	if _, err := io.ReadFull(z.src, buf[:2]); err != nil {
		return err
	}
	z.Pos += 2

	b := buf[0]
	if b>>6 != Version {
		return fmt.Errorf("lz4.Read: invalid version: got %d expected %d", b>>6, Version)
	}
	z.BlockDependency = b>>5&1 == 0
	z.BlockChecksum = b>>4&1 > 0
	frameSize := b>>3&1 > 0
	z.NoChecksum = b>>2&1 == 0
	

	bmsID := buf[1] >> 4 & 0x7
	bSize, ok := bsMapID[bmsID]
	if !ok {
		return fmt.Errorf("lz4.Read: invalid block max size: %d", bmsID)
	}
	z.BlockMaxSize = bSize

	z.checksum.Write(buf[0:2])

	if frameSize {
		if err := binary.Read(z.src, binary.LittleEndian, &z.Size); err != nil {
			return err
		}
		z.Pos += 8
		binary.LittleEndian.PutUint64(buf[:], z.Size)
		z.checksum.Write(buf[0:8])
	}

	
	
	
	
	
	
	
	

	
	if _, err := io.ReadFull(z.src, buf[:1]); err != nil {
		return err
	}
	z.Pos++
	if h := byte(z.checksum.Sum32() >> 8 & 0xFF); h != buf[0] {
		return fmt.Errorf("lz4.Read: invalid header checksum: got %v expected %v", buf[0], h)
	}

	z.Header.done = true

	return nil
}











func (z *Reader) Read(buf []byte) (n int, err error) {
	if !z.Header.done {
		if err = z.readHeader(true); err != nil {
			return
		}
	}

	if len(buf) == 0 {
		return
	}

	
	if len(z.data) > 0 {
		n = copy(buf, z.data)
		z.data = z.data[n:]
		if len(z.data) == 0 {
			z.data = nil
		}
		return
	}

	
	
	
	
	wbuf := buf
	zn := (len(wbuf) + z.BlockMaxSize - 1) / z.BlockMaxSize
	zblocks := make([]block, zn)
	for zi, abort := 0, uint32(0); zi < zn && atomic.LoadUint32(&abort) == 0; zi++ {
		zb := &zblocks[zi]
		
		if len(wbuf) < z.BlockMaxSize+len(z.window) {
			wbuf = make([]byte, z.BlockMaxSize+len(z.window))
		}
		copy(wbuf, z.window)
		if zb.err = z.readBlock(wbuf, zb); zb.err != nil {
			break
		}
		wbuf = wbuf[z.BlockMaxSize:]
		if !z.BlockDependency {
			z.wg.Add(1)
			go z.decompressBlock(zb, &abort)
			continue
		}
		
		z.decompressBlock(zb, nil)
		
		if len(z.window) == 0 {
			z.window = make([]byte, winSize)
		}
		if len(zb.data) >= winSize {
			copy(z.window, zb.data[len(zb.data)-winSize:])
		} else {
			copy(z.window, z.window[len(zb.data):])
			copy(z.window[len(zb.data)+1:], zb.data)
		}
	}
	z.wg.Wait()

	
	for _, zb := range zblocks {
		if zb.err != nil {
			if zb.err == errEndOfBlock {
				return n, z.close()
			}
			return n, zb.err
		}
		bLen := len(zb.data)
		if !z.NoChecksum {
			z.checksum.Write(zb.data)
		}
		m := copy(buf[n:], zb.data)
		
		if m < bLen {
			z.data = zb.data[m:]
		}
		n += m
	}

	return
}




func (z *Reader) readBlock(buf []byte, b *block) error {
	var bLen uint32
	if err := binary.Read(z.src, binary.LittleEndian, &bLen); err != nil {
		return err
	}
	atomic.AddInt64(&z.Pos, 4)

	switch {
	case bLen == 0:
		return errEndOfBlock
	case bLen&(1<<31) == 0:
		b.compressed = true
		b.data = buf
		b.zdata = make([]byte, bLen)
	default:
		bLen = bLen & (1<<31 - 1)
		if int(bLen) > len(buf) {
			return fmt.Errorf("lz4.Read: invalid block size: %d", bLen)
		}
		b.data = buf[:bLen]
		b.zdata = buf[:bLen]
	}
	if _, err := io.ReadFull(z.src, b.zdata); err != nil {
		return err
	}

	if z.BlockChecksum {
		if err := binary.Read(z.src, binary.LittleEndian, &b.checksum); err != nil {
			return err
		}
		xxh := hashPool.Get()
		defer hashPool.Put(xxh)
		xxh.Write(b.zdata)
		if h := xxh.Sum32(); h != b.checksum {
			return fmt.Errorf("lz4.Read: invalid block checksum: got %x expected %x", h, b.checksum)
		}
	}

	return nil
}



func (z *Reader) decompressBlock(b *block, abort *uint32) {
	if abort != nil {
		defer z.wg.Done()
	}
	if b.compressed {
		n := len(z.window)
		m, err := UncompressBlock(b.zdata, b.data, n)
		if err != nil {
			if abort != nil {
				atomic.StoreUint32(abort, 1)
			}
			b.err = err
			return
		}
		b.data = b.data[n : n+m]
	}
	atomic.AddInt64(&z.Pos, int64(len(b.data)))
}


func (z *Reader) close() error {
	if !z.NoChecksum {
		var checksum uint32
		if err := binary.Read(z.src, binary.LittleEndian, &checksum); err != nil {
			return err
		}
		if checksum != z.checksum.Sum32() {
			return fmt.Errorf("lz4.Read: invalid frame checksum: got %x expected %x", z.checksum.Sum32(), checksum)
		}
	}

	
	pos := z.Pos
	z.Reset(z.src)
	z.Pos = pos

	
	return z.readHeader(false)
}




func (z *Reader) Reset(r io.Reader) {
	z.Header = Header{}
	z.Pos = 0
	z.src = r
	z.checksum.Reset()
	z.data = nil
	z.window = nil
}



func (z *Reader) WriteTo(w io.Writer) (n int64, err error) {
	cpus := runtime.GOMAXPROCS(0)
	var buf []byte

	
	
	
	
	
	for {
		nsize := 0
		
		
		if z.BlockDependency {
			
			
			nsize = len(z.window) + z.BlockMaxSize
		} else {
			
			nsize = cpus * z.BlockMaxSize
		}
		if nsize != len(buf) {
			buf = make([]byte, nsize)
		}

		m, er := z.Read(buf)
		if er != nil && er != io.EOF {
			return n, er
		}
		m, err = w.Write(buf[:m])
		n += int64(m)
		if err != nil || er == io.EOF {
			return
		}
	}
}
