

package idna




































type info uint16

const (
	catSmallMask = 0x3
	catBigMask   = 0xF8
	indexShift   = 3
	xorBit       = 0x4    
	inlineXOR    = 0xE000 

	joinShift = 8
	joinMask  = 0x07

	
	attributesMask = 0x1800
	viramaModifier = 0x1800
	modifier       = 0x1000
	rtl            = 0x0800

	mayNeedNorm = 0x2000
)


type category uint16

const (
	unknown              category = 0 
	mapped               category = 1
	disallowedSTD3Mapped category = 2
	deviation            category = 3
)

const (
	valid               category = 0x08
	validNV8            category = 0x18
	validXV8            category = 0x28
	disallowed          category = 0x40
	disallowedSTD3Valid category = 0x80
	ignored             category = 0xC0
)


const (
	joiningL = (iota + 1)
	joiningD
	joiningT
	joiningR

	
	joinZWJ
	joinZWNJ
	joinVirama
	numJoinTypes
)

func (c info) isMapped() bool {
	return c&0x3 != 0
}

func (c info) category() category {
	small := c & catSmallMask
	if small != 0 {
		return category(small)
	}
	return category(c & catBigMask)
}

func (c info) joinType() info {
	if c.isMapped() {
		return 0
	}
	return (c >> joinShift) & joinMask
}

func (c info) isModifier() bool {
	return c&(modifier|catSmallMask) == modifier
}

func (c info) isViramaModifier() bool {
	return c&(attributesMask|catSmallMask) == viramaModifier
}
