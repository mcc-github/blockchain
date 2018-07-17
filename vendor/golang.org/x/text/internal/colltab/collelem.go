



package colltab

import (
	"fmt"
	"unicode"
)







type Level int

const (
	Primary Level = iota
	Secondary
	Tertiary
	Quaternary
	Identity

	NumLevels
)

const (
	defaultSecondary = 0x20
	defaultTertiary  = 0x2
	maxTertiary      = 0x1F
	MaxQuaternary    = 0x1FFFFF 
)





type Elem uint32

const (
	maxCE       Elem = 0xAFFFFFFF
	PrivateUse       = minContract
	minContract      = 0xC0000000
	maxContract      = 0xDFFFFFFF
	minExpand        = 0xE0000000
	maxExpand        = 0xEFFFFFFF
	minDecomp        = 0xF0000000
)

type ceType int

const (
	ceNormal           ceType = iota 
	ceContractionIndex               
	ceExpansionIndex                 
	ceDecompose                      
)

func (ce Elem) ctype() ceType {
	if ce <= maxCE {
		return ceNormal
	}
	if ce <= maxContract {
		return ceContractionIndex
	} else {
		if ce <= maxExpand {
			return ceExpansionIndex
		}
		return ceDecompose
	}
	panic("should not reach here")
	return ceType(-1)
}






















const (
	ceTypeMask              = 0xC0000000
	ceTypeMaskExt           = 0xE0000000
	ceIgnoreMask            = 0xF00FFFFF
	ceType1                 = 0x40000000
	ceType2                 = 0x00000000
	ceType3or4              = 0x80000000
	ceType4                 = 0xA0000000
	ceTypeQ                 = 0xC0000000
	Ignore                  = ceType4
	firstNonPrimary         = 0x80000000
	lastSpecialPrimary      = 0xA0000000
	secondaryMask           = 0x80000000
	hasTertiaryMask         = 0x40000000
	primaryValueMask        = 0x3FFFFE00
	maxPrimaryBits          = 21
	compactPrimaryBits      = 16
	maxSecondaryBits        = 12
	maxTertiaryBits         = 8
	maxCCCBits              = 8
	maxSecondaryCompactBits = 8
	maxSecondaryDiffBits    = 4
	maxTertiaryCompactBits  = 5
	primaryShift            = 9
	compactSecondaryShift   = 5
	minCompactSecondary     = defaultSecondary - 4
)

func makeImplicitCE(primary int) Elem {
	return ceType1 | Elem(primary<<primaryShift) | defaultSecondary
}



func MakeElem(primary, secondary, tertiary int, ccc uint8) (Elem, error) {
	if w := primary; w >= 1<<maxPrimaryBits || w < 0 {
		return 0, fmt.Errorf("makeCE: primary weight out of bounds: %x >= %x", w, 1<<maxPrimaryBits)
	}
	if w := secondary; w >= 1<<maxSecondaryBits || w < 0 {
		return 0, fmt.Errorf("makeCE: secondary weight out of bounds: %x >= %x", w, 1<<maxSecondaryBits)
	}
	if w := tertiary; w >= 1<<maxTertiaryBits || w < 0 {
		return 0, fmt.Errorf("makeCE: tertiary weight out of bounds: %x >= %x", w, 1<<maxTertiaryBits)
	}
	ce := Elem(0)
	if primary != 0 {
		if ccc != 0 {
			if primary >= 1<<compactPrimaryBits {
				return 0, fmt.Errorf("makeCE: primary weight with non-zero CCC out of bounds: %x >= %x", primary, 1<<compactPrimaryBits)
			}
			if secondary != defaultSecondary {
				return 0, fmt.Errorf("makeCE: cannot combine non-default secondary value (%x) with non-zero CCC (%x)", secondary, ccc)
			}
			ce = Elem(tertiary << (compactPrimaryBits + maxCCCBits))
			ce |= Elem(ccc) << compactPrimaryBits
			ce |= Elem(primary)
			ce |= ceType3or4
		} else if tertiary == defaultTertiary {
			if secondary >= 1<<maxSecondaryCompactBits {
				return 0, fmt.Errorf("makeCE: secondary weight with non-zero primary out of bounds: %x >= %x", secondary, 1<<maxSecondaryCompactBits)
			}
			ce = Elem(primary<<(maxSecondaryCompactBits+1) + secondary)
			ce |= ceType1
		} else {
			d := secondary - defaultSecondary + maxSecondaryDiffBits
			if d >= 1<<maxSecondaryDiffBits || d < 0 {
				return 0, fmt.Errorf("makeCE: secondary weight diff out of bounds: %x < 0 || %x > %x", d, d, 1<<maxSecondaryDiffBits)
			}
			if tertiary >= 1<<maxTertiaryCompactBits {
				return 0, fmt.Errorf("makeCE: tertiary weight with non-zero primary out of bounds: %x > %x", tertiary, 1<<maxTertiaryCompactBits)
			}
			ce = Elem(primary<<maxSecondaryDiffBits + d)
			ce = ce<<maxTertiaryCompactBits + Elem(tertiary)
		}
	} else {
		ce = Elem(secondary<<maxTertiaryBits + tertiary)
		ce += Elem(ccc) << (maxSecondaryBits + maxTertiaryBits)
		ce |= ceType4
	}
	return ce, nil
}


func MakeQuaternary(v int) Elem {
	return ceTypeQ | Elem(v<<primaryShift)
}




func (ce Elem) Mask(l Level) uint32 {
	return 0
}



func (ce Elem) CCC() uint8 {
	if ce&ceType3or4 != 0 {
		if ce&ceType4 == ceType3or4 {
			return uint8(ce >> 16)
		}
		return uint8(ce >> 20)
	}
	return 0
}


func (ce Elem) Primary() int {
	if ce >= firstNonPrimary {
		if ce > lastSpecialPrimary {
			return 0
		}
		return int(uint16(ce))
	}
	return int(ce&primaryValueMask) >> primaryShift
}


func (ce Elem) Secondary() int {
	switch ce & ceTypeMask {
	case ceType1:
		return int(uint8(ce))
	case ceType2:
		return minCompactSecondary + int((ce>>compactSecondaryShift)&0xF)
	case ceType3or4:
		if ce < ceType4 {
			return defaultSecondary
		}
		return int(ce>>8) & 0xFFF
	case ceTypeQ:
		return 0
	}
	panic("should not reach here")
}


func (ce Elem) Tertiary() uint8 {
	if ce&hasTertiaryMask == 0 {
		if ce&ceType3or4 == 0 {
			return uint8(ce & 0x1F)
		}
		if ce&ceType4 == ceType4 {
			return uint8(ce)
		}
		return uint8(ce>>24) & 0x1F 
	} else if ce&ceTypeMask == ceType1 {
		return defaultTertiary
	}
	
	return 0
}

func (ce Elem) updateTertiary(t uint8) Elem {
	if ce&ceTypeMask == ceType1 {
		
		nce := ce & primaryValueMask
		nce |= Elem(uint8(ce)-minCompactSecondary) << compactSecondaryShift
		ce = nce
	} else if ce&ceTypeMaskExt == ceType3or4 {
		ce &= ^Elem(maxTertiary << 24)
		return ce | (Elem(t) << 24)
	} else {
		
		ce &= ^Elem(maxTertiary)
	}
	return ce | Elem(t)
}




func (ce Elem) Quaternary() int {
	if ce&ceTypeMask == ceTypeQ {
		return int(ce&primaryValueMask) >> primaryShift
	} else if ce&ceIgnoreMask == Ignore {
		return 0
	}
	return MaxQuaternary
}


func (ce Elem) Weight(l Level) int {
	switch l {
	case Primary:
		return ce.Primary()
	case Secondary:
		return ce.Secondary()
	case Tertiary:
		return int(ce.Tertiary())
	case Quaternary:
		return ce.Quaternary()
	}
	return 0 
}







const (
	maxNBits              = 4
	maxTrieIndexBits      = 12
	maxContractOffsetBits = 13
)

func splitContractIndex(ce Elem) (index, n, offset int) {
	n = int(ce & (1<<maxNBits - 1))
	ce >>= maxNBits
	index = int(ce & (1<<maxTrieIndexBits - 1))
	ce >>= maxTrieIndexBits
	offset = int(ce & (1<<maxContractOffsetBits - 1))
	return
}



const maxExpandIndexBits = 16

func splitExpandIndex(ce Elem) (index int) {
	return int(uint16(ce))
}









func splitDecompose(ce Elem) (t1, t2 uint8) {
	return uint8(ce), uint8(ce >> 8)
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
