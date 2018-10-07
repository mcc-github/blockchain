package lz4

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pierrec/lz4/internal/xxh32"
)


type Writer struct {
	Header

	buf       [19]byte      
	dst       io.Writer     
	checksum  xxh32.XXHZero 
	zdata     []byte        
	data      []byte        
	idx       int           
	hashtable [winSize]int  
}





func NewWriter(dst io.Writer) *Writer {
	return &Writer{dst: dst}
}


func (z *Writer) writeHeader() error {
	
	if z.Header.BlockMaxSize == 0 {
		z.Header.BlockMaxSize = bsMapID[7]
	}
	
	bSize := z.Header.BlockMaxSize
	bSizeID, ok := bsMapValue[bSize]
	if !ok {
		return fmt.Errorf("lz4: invalid block max size: %d", bSize)
	}
	
	
	if n := 2 * bSize; cap(z.zdata) < n {
		z.zdata = make([]byte, n, n)
	}
	z.zdata = z.zdata[:bSize]
	z.data = z.zdata[:cap(z.zdata)][bSize:]
	z.idx = 0

	
	buf := z.buf[:]

	
	binary.LittleEndian.PutUint32(buf[0:], frameMagic)
	flg := byte(Version << 6)
	flg |= 1 << 5 
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
	buf[5] = bSizeID << 4

	
	n := 6
	
	if z.Header.Size > 0 {
		binary.LittleEndian.PutUint64(buf[n:], z.Header.Size)
		n += 8
	}

	
	buf[n] = byte(xxh32.ChecksumZero(buf[4:n]) >> 8 & 0xFF)
	z.checksum.Reset()

	
	if _, err := z.dst.Write(buf[0 : n+1]); err != nil {
		return err
	}
	z.Header.done = true
	if debugFlag {
		debug("wrote header %v", z.Header)
	}

	return nil
}



func (z *Writer) Write(buf []byte) (int, error) {
	if !z.Header.done {
		if err := z.writeHeader(); err != nil {
			return 0, err
		}
	}
	if debugFlag {
		debug("input buffer len=%d index=%d", len(buf), z.idx)
	}

	zn := len(z.data)
	var n int
	for len(buf) > 0 {
		if z.idx == 0 && len(buf) >= zn {
			
			if err := z.compressBlock(buf[:zn]); err != nil {
				return n, err
			}
			n += zn
			buf = buf[zn:]
			continue
		}
		
		m := copy(z.data[z.idx:], buf)
		n += m
		z.idx += m
		buf = buf[m:]
		if debugFlag {
			debug("%d bytes copied to buf, current index %d", n, z.idx)
		}

		if z.idx < len(z.data) {
			
			if debugFlag {
				debug("need more data for compression")
			}
			return n, nil
		}

		
		if err := z.compressBlock(z.data); err != nil {
			return n, err
		}
		z.idx = 0
	}

	return n, nil
}


func (z *Writer) compressBlock(data []byte) error {
	if !z.NoChecksum {
		z.checksum.Write(data)
	}

	
	var zn int
	var err error

	if level := z.Header.CompressionLevel; level != 0 {
		zn, err = CompressBlockHC(data, z.zdata, level)
	} else {
		zn, err = CompressBlock(data, z.zdata, z.hashtable[:])
	}

	var zdata []byte
	var bLen uint32
	if debugFlag {
		debug("block compression %d => %d", len(data), zn)
	}
	if err == nil && zn > 0 && zn < len(data) {
		
		bLen = uint32(zn)
		zdata = z.zdata[:zn]
	} else {
		
		bLen = uint32(len(data)) | compressedBlockFlag
		zdata = data
	}
	if debugFlag {
		debug("block compression to be written len=%d data len=%d", bLen, len(zdata))
	}

	
	if err := z.writeUint32(bLen); err != nil {
		return err
	}
	if _, err := z.dst.Write(zdata); err != nil {
		return err
	}

	if z.BlockChecksum {
		checksum := xxh32.ChecksumZero(zdata)
		if debugFlag {
			debug("block checksum %x", checksum)
		}
		if err := z.writeUint32(checksum); err != nil {
			return err
		}
	}
	if debugFlag {
		debug("current frame checksum %x", z.checksum.Sum32())
	}

	return nil
}




func (z *Writer) Flush() error {
	if debugFlag {
		debug("flush with index %d", z.idx)
	}
	if z.idx == 0 {
		return nil
	}

	return z.compressBlock(z.data[:z.idx])
}


func (z *Writer) Close() error {
	if !z.Header.done {
		if err := z.writeHeader(); err != nil {
			return err
		}
	}

	if err := z.Flush(); err != nil {
		return err
	}

	if debugFlag {
		debug("writing last empty block")
	}
	if err := z.writeUint32(0); err != nil {
		return err
	}
	if !z.NoChecksum {
		checksum := z.checksum.Sum32()
		if debugFlag {
			debug("stream checksum %x", checksum)
		}
		if err := z.writeUint32(checksum); err != nil {
			return err
		}
	}
	return nil
}




func (z *Writer) Reset(w io.Writer) {
	z.Header = Header{}
	z.dst = w
	z.checksum.Reset()
	z.zdata = z.zdata[:0]
	z.data = z.data[:0]
	z.idx = 0
}


func (z *Writer) writeUint32(x uint32) error {
	buf := z.buf[:4]
	binary.LittleEndian.PutUint32(buf, x)
	_, err := z.dst.Write(buf)
	return err
}
