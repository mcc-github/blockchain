



package build

import (
	"fmt"
	"unicode"

	"golang.org/x/text/internal/colltab"
)

const (
	defaultSecondary = 0x20
	defaultTertiary  = 0x2
	maxTertiary      = 0x1F
)

type rawCE struct {
	w   []int
	ccc uint8
}

func makeRawCE(w []int, ccc uint8) rawCE {
	ce := rawCE{w: make([]int, 4), ccc: ccc}
	copy(ce.w, w)
	return ce
}








const (
	maxPrimaryBits   = 21
	maxSecondaryBits = 12
	maxTertiaryBits  = 8
)

func makeCE(ce rawCE) (uint32, error) {
	v, e := colltab.MakeElem(ce.w[0], ce.w[1], ce.w[2], ce.ccc)
	return uint32(v), e
}







const (
	contractID            = 0xC0000000
	maxNBits              = 4
	maxTrieIndexBits      = 12
	maxContractOffsetBits = 13
)

func makeContractIndex(h ctHandle, offset int) (uint32, error) {
	if h.n >= 1<<maxNBits {
		return 0, fmt.Errorf("size of contraction trie node too large: %d >= %d", h.n, 1<<maxNBits)
	}
	if h.index >= 1<<maxTrieIndexBits {
		return 0, fmt.Errorf("size of contraction trie offset too large: %d >= %d", h.index, 1<<maxTrieIndexBits)
	}
	if offset >= 1<<maxContractOffsetBits {
		return 0, fmt.Errorf("contraction offset out of bounds: %x >= %x", offset, 1<<maxContractOffsetBits)
	}
	ce := uint32(contractID)
	ce += uint32(offset << (maxNBits + maxTrieIndexBits))
	ce += uint32(h.index << maxNBits)
	ce += uint32(h.n)
	return ce, nil
}




const (
	expandID           = 0xE0000000
	maxExpandIndexBits = 16
)

func makeExpandIndex(index int) (uint32, error) {
	if index >= 1<<maxExpandIndexBits {
		return 0, fmt.Errorf("expansion index out of bounds: %x >= %x", index, 1<<maxExpandIndexBits)
	}
	return expandID + uint32(index), nil
}



func makeExpansionHeader(n int) (uint32, error) {
	return uint32(n), nil
}










const (
	decompID = 0xF0000000
)

func makeDecompose(t1, t2 int) (uint32, error) {
	if t1 >= 256 || t1 < 0 {
		return 0, fmt.Errorf("first tertiary weight out of bounds: %d >= 256", t1)
	}
	if t2 >= 256 || t2 < 0 {
		return 0, fmt.Errorf("second tertiary weight out of bounds: %d >= 256", t2)
	}
	return uint32(t2<<8+t1) + decompID, nil
}

const (
	
	minUnified       rune = 0x4E00
	maxUnified            = 0x9FFF
	minCompatibility      = 0xF900
	maxCompatibility      = 0xFAFF
	minRare               = 0x3400
	maxRare               = 0x4DBF
)
const (
	commonUnifiedOffset = 0x10000
	rareUnifiedOffset   = 0x20000 
	otherOffset         = 0x50000 
	illegalOffset       = otherOffset + int(unicode.MaxRune)
	maxPrimary          = illegalOffset + 1
)






func implicitPrimary(r rune) int {
	if unicode.Is(unicode.Ideographic, r) {
		if r >= minUnified && r <= maxUnified {
			
			return int(r) + commonUnifiedOffset
		}
		if r >= minCompatibility && r <= maxCompatibility {
			
			
			return int(r) + commonUnifiedOffset
		}
		return int(r) + rareUnifiedOffset
	}
	return int(r) + otherOffset
}









func convertLargeWeights(elems []rawCE) (res []rawCE, err error) {
	const (
		cjkPrimaryStart   = 0xFB40
		rarePrimaryStart  = 0xFB80
		otherPrimaryStart = 0xFBC0
		illegalPrimary    = 0xFFFE
		highBitsMask      = 0x3F
		lowBitsMask       = 0x7FFF
		lowBitsFlag       = 0x8000
		shiftBits         = 15
	)
	for i := 0; i < len(elems); i++ {
		ce := elems[i].w
		p := ce[0]
		if p < cjkPrimaryStart {
			continue
		}
		if p > 0xFFFF {
			return elems, fmt.Errorf("found primary weight %X; should be <= 0xFFFF", p)
		}
		if p >= illegalPrimary {
			ce[0] = illegalOffset + p - illegalPrimary
		} else {
			if i+1 >= len(elems) {
				return elems, fmt.Errorf("second part of double primary weight missing: %v", elems)
			}
			if elems[i+1].w[0]&lowBitsFlag == 0 {
				return elems, fmt.Errorf("malformed second part of double primary weight: %v", elems)
			}
			np := ((p & highBitsMask) << shiftBits) + elems[i+1].w[0]&lowBitsMask
			switch {
			case p < rarePrimaryStart:
				np += commonUnifiedOffset
			case p < otherPrimaryStart:
				np += rareUnifiedOffset
			default:
				p += otherOffset
			}
			ce[0] = np
			for j := i + 1; j+1 < len(elems); j++ {
				elems[j] = elems[j+1]
			}
			elems = elems[:len(elems)-1]
		}
	}
	return elems, nil
}



func nextWeight(level colltab.Level, elems []rawCE) []rawCE {
	if level == colltab.Identity {
		next := make([]rawCE, len(elems))
		copy(next, elems)
		return next
	}
	next := []rawCE{makeRawCE(elems[0].w, elems[0].ccc)}
	next[0].w[level]++
	if level < colltab.Secondary {
		next[0].w[colltab.Secondary] = defaultSecondary
	}
	if level < colltab.Tertiary {
		next[0].w[colltab.Tertiary] = defaultTertiary
	}
	
	for _, ce := range elems[1:] {
		skip := true
		for i := colltab.Primary; i < level; i++ {
			skip = skip && ce.w[i] == 0
		}
		if !skip {
			next = append(next, ce)
		}
	}
	return next
}

func nextVal(elems []rawCE, i int, level colltab.Level) (index, value int) {
	for ; i < len(elems) && elems[i].w[level] == 0; i++ {
	}
	if i < len(elems) {
		return i, elems[i].w[level]
	}
	return i, 0
}



func compareWeights(a, b []rawCE) (result int, level colltab.Level) {
	for level := colltab.Primary; level < colltab.Identity; level++ {
		var va, vb int
		for ia, ib := 0, 0; ia < len(a) || ib < len(b); ia, ib = ia+1, ib+1 {
			ia, va = nextVal(a, ia, level)
			ib, vb = nextVal(b, ib, level)
			if va != vb {
				if va < vb {
					return -1, level
				} else {
					return 1, level
				}
			}
		}
	}
	return 0, colltab.Identity
}

func equalCE(a, b rawCE) bool {
	for i := 0; i < 3; i++ {
		if b.w[i] != a.w[i] {
			return false
		}
	}
	return true
}

func equalCEArrays(a, b []rawCE) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !equalCE(a[i], b[i]) {
			return false
		}
	}
	return true
}
