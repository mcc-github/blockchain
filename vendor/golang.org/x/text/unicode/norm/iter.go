



package norm

import (
	"fmt"
	"unicode/utf8"
)



const MaxSegmentSize = maxByteBufferSize



type Iter struct {
	rb     reorderBuffer
	buf    [maxByteBufferSize]byte
	info   Properties 
	next   iterFunc   
	asciiF iterFunc

	p        int    
	multiSeg []byte 
}

type iterFunc func(*Iter) []byte


func (i *Iter) Init(f Form, src []byte) {
	i.p = 0
	if len(src) == 0 {
		i.setDone()
		i.rb.nsrc = 0
		return
	}
	i.multiSeg = nil
	i.rb.init(f, src)
	i.next = i.rb.f.nextMain
	i.asciiF = nextASCIIBytes
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.rb.ss.first(i.info)
}


func (i *Iter) InitString(f Form, src string) {
	i.p = 0
	if len(src) == 0 {
		i.setDone()
		i.rb.nsrc = 0
		return
	}
	i.multiSeg = nil
	i.rb.initString(f, src)
	i.next = i.rb.f.nextMain
	i.asciiF = nextASCIIString
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.rb.ss.first(i.info)
}




func (i *Iter) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case 0:
		abs = offset
	case 1:
		abs = int64(i.p) + offset
	case 2:
		abs = int64(i.rb.nsrc) + offset
	default:
		return 0, fmt.Errorf("norm: invalid whence")
	}
	if abs < 0 {
		return 0, fmt.Errorf("norm: negative position")
	}
	if int(abs) >= i.rb.nsrc {
		i.setDone()
		return int64(i.p), nil
	}
	i.p = int(abs)
	i.multiSeg = nil
	i.next = i.rb.f.nextMain
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.rb.ss.first(i.info)
	return abs, nil
}





func (i *Iter) returnSlice(a, b int) []byte {
	if i.rb.src.bytes == nil {
		return i.buf[:copy(i.buf[:], i.rb.src.str[a:b])]
	}
	return i.rb.src.bytes[a:b]
}


func (i *Iter) Pos() int {
	return i.p
}

func (i *Iter) setDone() {
	i.next = nextDone
	i.p = i.rb.nsrc
}


func (i *Iter) Done() bool {
	return i.p >= i.rb.nsrc
}






func (i *Iter) Next() []byte {
	return i.next(i)
}

func nextASCIIBytes(i *Iter) []byte {
	p := i.p + 1
	if p >= i.rb.nsrc {
		i.setDone()
		return i.rb.src.bytes[i.p:p]
	}
	if i.rb.src.bytes[p] < utf8.RuneSelf {
		p0 := i.p
		i.p = p
		return i.rb.src.bytes[p0:p]
	}
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.next = i.rb.f.nextMain
	return i.next(i)
}

func nextASCIIString(i *Iter) []byte {
	p := i.p + 1
	if p >= i.rb.nsrc {
		i.buf[0] = i.rb.src.str[i.p]
		i.setDone()
		return i.buf[:1]
	}
	if i.rb.src.str[p] < utf8.RuneSelf {
		i.buf[0] = i.rb.src.str[i.p]
		i.p = p
		return i.buf[:1]
	}
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.next = i.rb.f.nextMain
	return i.next(i)
}

func nextHangul(i *Iter) []byte {
	p := i.p
	next := p + hangulUTF8Size
	if next >= i.rb.nsrc {
		i.setDone()
	} else if i.rb.src.hangul(next) == 0 {
		i.rb.ss.next(i.info)
		i.info = i.rb.f.info(i.rb.src, i.p)
		i.next = i.rb.f.nextMain
		return i.next(i)
	}
	i.p = next
	return i.buf[:decomposeHangul(i.buf[:], i.rb.src.hangul(p))]
}

func nextDone(i *Iter) []byte {
	return nil
}



func nextMulti(i *Iter) []byte {
	j := 0
	d := i.multiSeg
	
	for j = 1; j < len(d) && !utf8.RuneStart(d[j]); j++ {
	}
	for j < len(d) {
		info := i.rb.f.info(input{bytes: d}, j)
		if info.BoundaryBefore() {
			i.multiSeg = d[j:]
			return d[:j]
		}
		j += int(info.size)
	}
	
	i.next = i.rb.f.nextMain
	return i.next(i)
}



func nextMultiNorm(i *Iter) []byte {
	j := 0
	d := i.multiSeg
	for j < len(d) {
		info := i.rb.f.info(input{bytes: d}, j)
		if info.BoundaryBefore() {
			i.rb.compose()
			seg := i.buf[:i.rb.flushCopy(i.buf[:])]
			i.rb.insertUnsafe(input{bytes: d}, j, info)
			i.multiSeg = d[j+int(info.size):]
			return seg
		}
		i.rb.insertUnsafe(input{bytes: d}, j, info)
		j += int(info.size)
	}
	i.multiSeg = nil
	i.next = nextComposed
	return doNormComposed(i)
}


func nextDecomposed(i *Iter) (next []byte) {
	outp := 0
	inCopyStart, outCopyStart := i.p, 0
	for {
		if sz := int(i.info.size); sz <= 1 {
			i.rb.ss = 0
			p := i.p
			i.p++ 
			if i.p >= i.rb.nsrc {
				i.setDone()
				return i.returnSlice(p, i.p)
			} else if i.rb.src._byte(i.p) < utf8.RuneSelf {
				i.next = i.asciiF
				return i.returnSlice(p, i.p)
			}
			outp++
		} else if d := i.info.Decomposition(); d != nil {
			
			
			
			
			p := outp + len(d)
			if outp > 0 {
				i.rb.src.copySlice(i.buf[outCopyStart:], inCopyStart, i.p)
				
				
				if p > len(i.buf) {
					return i.buf[:outp]
				}
			} else if i.info.multiSegment() {
				
				
				if i.multiSeg == nil {
					i.multiSeg = d
					i.next = nextMulti
					return nextMulti(i)
				}
				
				d = i.multiSeg
				i.multiSeg = nil
				p = len(d)
			}
			prevCC := i.info.tccc
			if i.p += sz; i.p >= i.rb.nsrc {
				i.setDone()
				i.info = Properties{} 
			} else {
				i.info = i.rb.f.info(i.rb.src, i.p)
			}
			switch i.rb.ss.next(i.info) {
			case ssOverflow:
				i.next = nextCGJDecompose
				fallthrough
			case ssStarter:
				if outp > 0 {
					copy(i.buf[outp:], d)
					return i.buf[:p]
				}
				return d
			}
			copy(i.buf[outp:], d)
			outp = p
			inCopyStart, outCopyStart = i.p, outp
			if i.info.ccc < prevCC {
				goto doNorm
			}
			continue
		} else if r := i.rb.src.hangul(i.p); r != 0 {
			outp = decomposeHangul(i.buf[:], r)
			i.p += hangulUTF8Size
			inCopyStart, outCopyStart = i.p, outp
			if i.p >= i.rb.nsrc {
				i.setDone()
				break
			} else if i.rb.src.hangul(i.p) != 0 {
				i.next = nextHangul
				return i.buf[:outp]
			}
		} else {
			p := outp + sz
			if p > len(i.buf) {
				break
			}
			outp = p
			i.p += sz
		}
		if i.p >= i.rb.nsrc {
			i.setDone()
			break
		}
		prevCC := i.info.tccc
		i.info = i.rb.f.info(i.rb.src, i.p)
		if v := i.rb.ss.next(i.info); v == ssStarter {
			break
		} else if v == ssOverflow {
			i.next = nextCGJDecompose
			break
		}
		if i.info.ccc < prevCC {
			goto doNorm
		}
	}
	if outCopyStart == 0 {
		return i.returnSlice(inCopyStart, i.p)
	} else if inCopyStart < i.p {
		i.rb.src.copySlice(i.buf[outCopyStart:], inCopyStart, i.p)
	}
	return i.buf[:outp]
doNorm:
	
	
	i.rb.src.copySlice(i.buf[outCopyStart:], inCopyStart, i.p)
	i.rb.insertDecomposed(i.buf[0:outp])
	return doNormDecomposed(i)
}

func doNormDecomposed(i *Iter) []byte {
	for {
		i.rb.insertUnsafe(i.rb.src, i.p, i.info)
		if i.p += int(i.info.size); i.p >= i.rb.nsrc {
			i.setDone()
			break
		}
		i.info = i.rb.f.info(i.rb.src, i.p)
		if i.info.ccc == 0 {
			break
		}
		if s := i.rb.ss.next(i.info); s == ssOverflow {
			i.next = nextCGJDecompose
			break
		}
	}
	
	return i.buf[:i.rb.flushCopy(i.buf[:])]
}

func nextCGJDecompose(i *Iter) []byte {
	i.rb.ss = 0
	i.rb.insertCGJ()
	i.next = nextDecomposed
	i.rb.ss.first(i.info)
	buf := doNormDecomposed(i)
	return buf
}


func nextComposed(i *Iter) []byte {
	outp, startp := 0, i.p
	var prevCC uint8
	for {
		if !i.info.isYesC() {
			goto doNorm
		}
		prevCC = i.info.tccc
		sz := int(i.info.size)
		if sz == 0 {
			sz = 1 
		}
		p := outp + sz
		if p > len(i.buf) {
			break
		}
		outp = p
		i.p += sz
		if i.p >= i.rb.nsrc {
			i.setDone()
			break
		} else if i.rb.src._byte(i.p) < utf8.RuneSelf {
			i.rb.ss = 0
			i.next = i.asciiF
			break
		}
		i.info = i.rb.f.info(i.rb.src, i.p)
		if v := i.rb.ss.next(i.info); v == ssStarter {
			break
		} else if v == ssOverflow {
			i.next = nextCGJCompose
			break
		}
		if i.info.ccc < prevCC {
			goto doNorm
		}
	}
	return i.returnSlice(startp, i.p)
doNorm:
	
	i.p = startp
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.rb.ss.first(i.info)
	if i.info.multiSegment() {
		d := i.info.Decomposition()
		info := i.rb.f.info(input{bytes: d}, 0)
		i.rb.insertUnsafe(input{bytes: d}, 0, info)
		i.multiSeg = d[int(info.size):]
		i.next = nextMultiNorm
		return nextMultiNorm(i)
	}
	i.rb.ss.first(i.info)
	i.rb.insertUnsafe(i.rb.src, i.p, i.info)
	return doNormComposed(i)
}

func doNormComposed(i *Iter) []byte {
	
	for {
		if i.p += int(i.info.size); i.p >= i.rb.nsrc {
			i.setDone()
			break
		}
		i.info = i.rb.f.info(i.rb.src, i.p)
		if s := i.rb.ss.next(i.info); s == ssStarter {
			break
		} else if s == ssOverflow {
			i.next = nextCGJCompose
			break
		}
		i.rb.insertUnsafe(i.rb.src, i.p, i.info)
	}
	i.rb.compose()
	seg := i.buf[:i.rb.flushCopy(i.buf[:])]
	return seg
}

func nextCGJCompose(i *Iter) []byte {
	i.rb.ss = 0 
	i.rb.insertCGJ()
	i.next = nextComposed
	
	
	
	i.rb.ss.first(i.info)
	i.rb.insertUnsafe(i.rb.src, i.p, i.info)
	return doNormComposed(i)
}
