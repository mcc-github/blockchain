



package norm

























const (
	qcInfoMask      = 0x3F 
	headerLenMask   = 0x3F 
	headerFlagsMask = 0xC0 
)


type Properties struct {
	pos   uint8  
	size  uint8  
	ccc   uint8  
	tccc  uint8  
	nLead uint8  
	flags qcInfo 
	index uint16
}


type lookupFunc func(b input, i int) Properties


type formInfo struct {
	form                     Form
	composing, compatibility bool 
	info                     lookupFunc
	nextMain                 iterFunc
}

var formTable = []*formInfo{{
	form:          NFC,
	composing:     true,
	compatibility: false,
	info:          lookupInfoNFC,
	nextMain:      nextComposed,
}, {
	form:          NFD,
	composing:     false,
	compatibility: false,
	info:          lookupInfoNFC,
	nextMain:      nextDecomposed,
}, {
	form:          NFKC,
	composing:     true,
	compatibility: true,
	info:          lookupInfoNFKC,
	nextMain:      nextComposed,
}, {
	form:          NFKD,
	composing:     false,
	compatibility: true,
	info:          lookupInfoNFKC,
	nextMain:      nextDecomposed,
}}









func (p Properties) BoundaryBefore() bool {
	if p.ccc == 0 && !p.combinesBackward() {
		return true
	}
	
	
	
	return false
}



func (p Properties) BoundaryAfter() bool {
	
	return p.isInert()
}









type qcInfo uint8

func (p Properties) isYesC() bool { return p.flags&0x10 == 0 }
func (p Properties) isYesD() bool { return p.flags&0x4 == 0 }

func (p Properties) combinesForward() bool  { return p.flags&0x20 != 0 }
func (p Properties) combinesBackward() bool { return p.flags&0x8 != 0 } 
func (p Properties) hasDecomposition() bool { return p.flags&0x4 != 0 } 

func (p Properties) isInert() bool {
	return p.flags&qcInfoMask == 0 && p.ccc == 0
}

func (p Properties) multiSegment() bool {
	return p.index >= firstMulti && p.index < endMulti
}

func (p Properties) nLeadingNonStarters() uint8 {
	return p.nLead
}

func (p Properties) nTrailingNonStarters() uint8 {
	return uint8(p.flags & 0x03)
}



func (p Properties) Decomposition() []byte {
	
	if p.index == 0 {
		return nil
	}
	i := p.index
	n := decomps[i] & headerLenMask
	i++
	return decomps[i : i+uint16(n)]
}


func (p Properties) Size() int {
	return int(p.size)
}


func (p Properties) CCC() uint8 {
	if p.index >= firstCCCZeroExcept {
		return 0
	}
	return ccc[p.ccc]
}



func (p Properties) LeadCCC() uint8 {
	return ccc[p.ccc]
}



func (p Properties) TrailCCC() uint8 {
	return ccc[p.tccc]
}









func combine(a, b rune) rune {
	key := uint32(uint16(a))<<16 + uint32(uint16(b))
	return recompMap[key]
}

func lookupInfoNFC(b input, i int) Properties {
	v, sz := b.charinfoNFC(i)
	return compInfo(v, sz)
}

func lookupInfoNFKC(b input, i int) Properties {
	v, sz := b.charinfoNFKC(i)
	return compInfo(v, sz)
}


func (f Form) Properties(s []byte) Properties {
	if f == NFC || f == NFD {
		return compInfo(nfcData.lookup(s))
	}
	return compInfo(nfkcData.lookup(s))
}


func (f Form) PropertiesString(s string) Properties {
	if f == NFC || f == NFD {
		return compInfo(nfcData.lookupString(s))
	}
	return compInfo(nfkcData.lookupString(s))
}




func compInfo(v uint16, sz int) Properties {
	if v == 0 {
		return Properties{size: uint8(sz)}
	} else if v >= 0x8000 {
		p := Properties{
			size:  uint8(sz),
			ccc:   uint8(v),
			tccc:  uint8(v),
			flags: qcInfo(v >> 8),
		}
		if p.ccc > 0 || p.combinesBackward() {
			p.nLead = uint8(p.flags & 0x3)
		}
		return p
	}
	
	h := decomps[v]
	f := (qcInfo(h&headerFlagsMask) >> 2) | 0x4
	p := Properties{size: uint8(sz), flags: f, index: v}
	if v >= firstCCC {
		v += uint16(h&headerLenMask) + 1
		c := decomps[v]
		p.tccc = c >> 2
		p.flags |= qcInfo(c & 0x3)
		if v >= firstLeadingCCC {
			p.nLead = c & 0x3
			if v >= firstStarterWithNLead {
				
				p.flags &= 0x03
				p.index = 0
				return p
			}
			p.ccc = decomps[v+1]
		}
	}
	return p
}
