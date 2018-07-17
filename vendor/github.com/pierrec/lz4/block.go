package lz4

import (
	"encoding/binary"
	"errors"
)



type block struct {
	compressed bool
	zdata      []byte 
	data       []byte 
	offset     int    
	checksum   uint32 
	err        error  
}

var (
	
	ErrInvalidSource = errors.New("lz4: invalid source")
	
	
	ErrShortBuffer = errors.New("lz4: short buffer")
)


func CompressBlockBound(n int) int {
	return n + n/255 + 16
}







func UncompressBlock(src, dst []byte, di int) (int, error) {
	si, sn, di0 := 0, len(src), di
	if sn == 0 {
		return 0, nil
	}

	for {
		
		lLen := int(src[si] >> 4)
		mLen := int(src[si] & 0xF)
		if si++; si == sn {
			return di, ErrInvalidSource
		}

		
		if lLen > 0 {
			if lLen == 0xF {
				for src[si] == 0xFF {
					lLen += 0xFF
					if si++; si == sn {
						return di - di0, ErrInvalidSource
					}
				}
				lLen += int(src[si])
				if si++; si == sn {
					return di - di0, ErrInvalidSource
				}
			}
			if len(dst)-di < lLen || si+lLen > sn {
				return di - di0, ErrShortBuffer
			}
			di += copy(dst[di:], src[si:si+lLen])

			if si += lLen; si >= sn {
				return di - di0, nil
			}
		}

		if si += 2; si >= sn {
			return di, ErrInvalidSource
		}
		offset := int(src[si-2]) | int(src[si-1])<<8
		if di-offset < 0 || offset == 0 {
			return di - di0, ErrInvalidSource
		}

		
		if mLen == 0xF {
			for src[si] == 0xFF {
				mLen += 0xFF
				if si++; si == sn {
					return di - di0, ErrInvalidSource
				}
			}
			mLen += int(src[si])
			if si++; si == sn {
				return di - di0, ErrInvalidSource
			}
		}
		
		mLen += 4
		if len(dst)-di <= mLen {
			return di - di0, ErrShortBuffer
		}

		
		if mLen >= offset {
			bytesToCopy := offset * (mLen / offset)
			
			
			expanded := dst[di-offset : di+bytesToCopy]
			n := offset
			for n <= bytesToCopy+offset {
				copy(expanded[n:], expanded[:n])
				n *= 2
			}
			di += bytesToCopy
			mLen -= bytesToCopy
		}

		di += copy(dst[di:], dst[di-offset:di-offset+mLen])
	}
}







func CompressBlock(src, dst []byte, soffset int) (int, error) {
	sn, dn := len(src)-mfLimit, len(dst)
	if sn <= 0 || dn == 0 || soffset >= sn {
		return 0, nil
	}
	var si, di int

	
	
	var hashTable [1 << hashLog]int
	var hashShift = uint((minMatch * 8) - hashLog)

	
	
	for si < soffset {
		h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift
		si++
		hashTable[h] = si
	}

	anchor := si
	fma := 1 << skipStrength
	for si < sn-minMatch {
		
		h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift
		
		ref := hashTable[h] - 1
		
		hashTable[h] = si + 1
		
		
		
		
		
		
		
		
		

		
		if ref < 0 || fma&(1<<skipStrength-1) < 4 ||
			(si-ref)>>winSizeLog > 0 ||
			src[ref] != src[si] ||
			src[ref+1] != src[si+1] ||
			src[ref+2] != src[si+2] ||
			src[ref+3] != src[si+3] {
			
			si += fma >> skipStrength
			fma++
			continue
		}
		
		fma = 1 << skipStrength
		lLen := si - anchor
		offset := si - ref

		
		si += minMatch
		mLen := si 
		for si <= sn && src[si] == src[si-offset] {
			si++
		}
		mLen = si - mLen
		if mLen < 0xF {
			dst[di] = byte(mLen)
		} else {
			dst[di] = 0xF
		}

		
		if lLen < 0xF {
			dst[di] |= byte(lLen << 4)
		} else {
			dst[di] |= 0xF0
			if di++; di == dn {
				return di, ErrShortBuffer
			}
			l := lLen - 0xF
			for ; l >= 0xFF; l -= 0xFF {
				dst[di] = 0xFF
				if di++; di == dn {
					return di, ErrShortBuffer
				}
			}
			dst[di] = byte(l)
		}
		if di++; di == dn {
			return di, ErrShortBuffer
		}

		
		if di+lLen >= dn {
			return di, ErrShortBuffer
		}
		di += copy(dst[di:], src[anchor:anchor+lLen])
		anchor = si

		
		if di += 2; di >= dn {
			return di, ErrShortBuffer
		}
		dst[di-2], dst[di-1] = byte(offset), byte(offset>>8)

		
		if mLen >= 0xF {
			for mLen -= 0xF; mLen >= 0xFF; mLen -= 0xFF {
				dst[di] = 0xFF
				if di++; di == dn {
					return di, ErrShortBuffer
				}
			}
			dst[di] = byte(mLen)
			if di++; di == dn {
				return di, ErrShortBuffer
			}
		}
	}

	if anchor == 0 {
		
		return 0, nil
	}

	
	lLen := len(src) - anchor
	if lLen < 0xF {
		dst[di] = byte(lLen << 4)
	} else {
		dst[di] = 0xF0
		if di++; di == dn {
			return di, ErrShortBuffer
		}
		lLen -= 0xF
		for ; lLen >= 0xFF; lLen -= 0xFF {
			dst[di] = 0xFF
			if di++; di == dn {
				return di, ErrShortBuffer
			}
		}
		dst[di] = byte(lLen)
	}
	if di++; di == dn {
		return di, ErrShortBuffer
	}

	
	src = src[anchor:]
	switch n := di + len(src); {
	case n > dn:
		return di, ErrShortBuffer
	case n >= sn:
		
		return 0, nil
	}
	di += copy(dst[di:], src)
	return di, nil
}







func CompressBlockHC(src, dst []byte, soffset int) (int, error) {
	sn, dn := len(src)-mfLimit, len(dst)
	if sn <= 0 || dn == 0 || soffset >= sn {
		return 0, nil
	}
	var si, di int

	
	
	
	var hashTable [1 << hashLog]int
	var chainTable [winSize]int
	var hashShift = uint((minMatch * 8) - hashLog)

	
	
	for si < soffset {
		h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift
		chainTable[si&winMask] = hashTable[h]
		si++
		hashTable[h] = si
	}

	anchor := si
	for si < sn-minMatch {
		
		h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift

		
		mLen := 0
		offset := 0
		for next := hashTable[h] - 1; next > 0 && next > si-winSize; next = chainTable[next&winMask] - 1 {
			
			if src[next+mLen] == src[si+mLen] {
				for ml := 0; ; ml++ {
					if src[next+ml] != src[si+ml] || si+ml > sn {
						
						if mLen < ml && ml >= minMatch {
							mLen = ml
							offset = si - next
						}
						break
					}
				}
			}
		}
		chainTable[si&winMask] = hashTable[h]
		hashTable[h] = si + 1

		
		if mLen == 0 {
			si++
			continue
		}

		
		
		
		for si, ml := si+1, si+mLen; si < ml; {
			h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift
			chainTable[si&winMask] = hashTable[h]
			si++
			hashTable[h] = si
		}

		lLen := si - anchor
		si += mLen
		mLen -= minMatch 

		if mLen < 0xF {
			dst[di] = byte(mLen)
		} else {
			dst[di] = 0xF
		}

		
		if lLen < 0xF {
			dst[di] |= byte(lLen << 4)
		} else {
			dst[di] |= 0xF0
			if di++; di == dn {
				return di, ErrShortBuffer
			}
			l := lLen - 0xF
			for ; l >= 0xFF; l -= 0xFF {
				dst[di] = 0xFF
				if di++; di == dn {
					return di, ErrShortBuffer
				}
			}
			dst[di] = byte(l)
		}
		if di++; di == dn {
			return di, ErrShortBuffer
		}

		
		if di+lLen >= dn {
			return di, ErrShortBuffer
		}
		di += copy(dst[di:], src[anchor:anchor+lLen])
		anchor = si

		
		if di += 2; di >= dn {
			return di, ErrShortBuffer
		}
		dst[di-2], dst[di-1] = byte(offset), byte(offset>>8)

		
		if mLen >= 0xF {
			for mLen -= 0xF; mLen >= 0xFF; mLen -= 0xFF {
				dst[di] = 0xFF
				if di++; di == dn {
					return di, ErrShortBuffer
				}
			}
			dst[di] = byte(mLen)
			if di++; di == dn {
				return di, ErrShortBuffer
			}
		}
	}

	if anchor == 0 {
		
		return 0, nil
	}

	
	lLen := len(src) - anchor
	if lLen < 0xF {
		dst[di] = byte(lLen << 4)
	} else {
		dst[di] = 0xF0
		if di++; di == dn {
			return di, ErrShortBuffer
		}
		lLen -= 0xF
		for ; lLen >= 0xFF; lLen -= 0xFF {
			dst[di] = 0xFF
			if di++; di == dn {
				return di, ErrShortBuffer
			}
		}
		dst[di] = byte(lLen)
	}
	if di++; di == dn {
		return di, ErrShortBuffer
	}

	
	src = src[anchor:]
	switch n := di + len(src); {
	case n > dn:
		return di, ErrShortBuffer
	case n >= sn:
		
		return 0, nil
	}
	di += copy(dst[di:], src)
	return di, nil
}
