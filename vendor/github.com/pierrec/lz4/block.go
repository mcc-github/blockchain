package lz4

import (
	"encoding/binary"
	"errors"
)

var (
	
	
	ErrInvalidSourceShortBuffer = errors.New("lz4: invalid source or destination buffer too short")
	
	ErrInvalid = errors.New("lz4: bad magic number")
)


func blockHash(x uint32) uint32 {
	const hasher uint32 = 2654435761 
	return x * hasher >> hashShift
}


func CompressBlockBound(n int) int {
	return n + n/255 + 16
}







func UncompressBlock(src, dst []byte) (si int, err error) {
	defer func() {
		
		
		if recover() != nil {
			err = ErrInvalidSourceShortBuffer
		}
	}()
	sn := len(src)
	if sn == 0 {
		return 0, nil
	}
	var di int

	for {
		
		b := int(src[si])
		si++

		
		if lLen := b >> 4; lLen > 0 {
			if lLen == 0xF {
				for src[si] == 0xFF {
					lLen += 0xFF
					si++
				}
				lLen += int(src[si])
				si++
			}
			i := si
			si += lLen
			di += copy(dst[di:], src[i:si])

			if si >= sn {
				return di, nil
			}
		}

		si++
		_ = src[si] 
		offset := int(src[si-1]) | int(src[si])<<8
		si++

		
		mLen := b & 0xF
		if mLen == 0xF {
			for src[si] == 0xFF {
				mLen += 0xFF
				si++
			}
			mLen += int(src[si])
			si++
		}
		mLen += minMatch

		
		i := di - offset
		if offset > 0 && mLen >= offset {
			
			bytesToCopy := offset * (mLen / offset)
			expanded := dst[i:]
			for n := offset; n <= bytesToCopy+offset; n *= 2 {
				copy(expanded[n:], expanded[:n])
			}
			di += bytesToCopy
			mLen -= bytesToCopy
		}
		di += copy(dst[di:], dst[i:i+mLen])
	}
}








func CompressBlock(src, dst []byte, hashTable []int) (di int, err error) {
	defer func() {
		if recover() != nil {
			err = ErrInvalidSourceShortBuffer
		}
	}()

	sn, dn := len(src)-mfLimit, len(dst)
	if sn <= 0 || dn == 0 {
		return 0, nil
	}
	var si int

	
	

	anchor := si 
	

	for si < sn {
		
		match := binary.LittleEndian.Uint32(src[si:])
		h := blockHash(match)

		ref := hashTable[h]
		hashTable[h] = si
		if ref >= sn { 
			si++
			continue
		}
		offset := si - ref
		if offset <= 0 || offset >= winSize || 
			match != binary.LittleEndian.Uint32(src[ref:]) { 
			
			
			si++
			continue
		}

		
		
		lLen := si - anchor 

		
		si += minMatch
		mLen := si 
		
		for si < sn && binary.LittleEndian.Uint64(src[si:]) == binary.LittleEndian.Uint64(src[si-offset:]) {
			si += 8
		}
		
		for si < sn && src[si] == src[si-offset] {
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
			di++
			l := lLen - 0xF
			for ; l >= 0xFF; l -= 0xFF {
				dst[di] = 0xFF
				di++
			}
			dst[di] = byte(l)
		}
		di++

		
		copy(dst[di:], src[anchor:anchor+lLen])
		di += lLen + 2
		anchor = si

		
		_ = dst[di] 
		dst[di-2], dst[di-1] = byte(offset), byte(offset>>8)

		
		if mLen >= 0xF {
			for mLen -= 0xF; mLen >= 0xFF; mLen -= 0xFF {
				dst[di] = 0xFF
				di++
			}
			dst[di] = byte(mLen)
			di++
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
		di++
		for lLen -= 0xF; lLen >= 0xFF; lLen -= 0xFF {
			dst[di] = 0xFF
			di++
		}
		dst[di] = byte(lLen)
	}
	di++

	
	if di >= anchor {
		
		return 0, nil
	}
	di += copy(dst[di:], src[anchor:])
	return di, nil
}









func CompressBlockHC(src, dst []byte, depth int) (di int, err error) {
	defer func() {
		if recover() != nil {
			err = ErrInvalidSourceShortBuffer
		}
	}()

	sn, dn := len(src)-mfLimit, len(dst)
	if sn <= 0 || dn == 0 {
		return 0, nil
	}
	var si int

	
	
	var hashTable, chainTable [winSize]int

	if depth <= 0 {
		depth = winSize
	}

	anchor := si
	for si < sn {
		
		match := binary.LittleEndian.Uint32(src[si:])
		h := blockHash(match)

		
		mLen := 0
		offset := 0
		for next, try := hashTable[h], depth; try > 0 && next > 0 && si-next < winSize; next = chainTable[next&winMask] {
			
			
			if src[next+mLen] != src[si+mLen] {
				continue
			}
			ml := 0
			
			for ml < sn-si && binary.LittleEndian.Uint64(src[next+ml:]) == binary.LittleEndian.Uint64(src[si+ml:]) {
				ml += 8
			}
			for ml < sn-si && src[next+ml] == src[si+ml] {
				ml++
			}
			if ml+1 < minMatch || ml <= mLen {
				
				continue
			}
			
			mLen = ml
			offset = si - next
			
			try--
		}
		chainTable[si&winMask] = hashTable[h]
		hashTable[h] = si

		
		if mLen == 0 {
			si++
			continue
		}

		
		
		
		winStart := si + 1
		if ws := si + mLen - winSize; ws > winStart {
			winStart = ws
		}
		for si, ml := winStart, si+mLen; si < ml; {
			match >>= 8
			match |= uint32(src[si+3]) << 24
			h := blockHash(match)
			chainTable[si&winMask] = hashTable[h]
			hashTable[h] = si
			si++
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
			di++
			l := lLen - 0xF
			for ; l >= 0xFF; l -= 0xFF {
				dst[di] = 0xFF
				di++
			}
			dst[di] = byte(l)
		}
		di++

		
		copy(dst[di:], src[anchor:anchor+lLen])
		di += lLen
		anchor = si

		
		di += 2
		dst[di-2], dst[di-1] = byte(offset), byte(offset>>8)

		
		if mLen >= 0xF {
			for mLen -= 0xF; mLen >= 0xFF; mLen -= 0xFF {
				dst[di] = 0xFF
				di++
			}
			dst[di] = byte(mLen)
			di++
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
		di++
		lLen -= 0xF
		for ; lLen >= 0xFF; lLen -= 0xFF {
			dst[di] = 0xFF
			di++
		}
		dst[di] = byte(lLen)
	}
	di++

	
	if di >= anchor {
		
		return 0, nil
	}
	di += copy(dst[di:], src[anchor:])
	return di, nil
}
