



package norm

import "unicode/utf8"

const (
	maxNonStarters = 30
	
	
	maxBufferSize    = maxNonStarters + 2
	maxNFCExpansion  = 3  
	maxNFKCExpansion = 18 

	maxByteBufferSize = utf8.UTFMax * maxBufferSize 
)



type ssState int

const (
	
	ssSuccess ssState = iota
	
	ssStarter
	
	ssOverflow
)


type streamSafe uint8



func (ss *streamSafe) first(p Properties) {
	*ss = streamSafe(p.nTrailingNonStarters())
}



func (ss *streamSafe) next(p Properties) ssState {
	if *ss > maxNonStarters {
		panic("streamSafe was not reset")
	}
	n := p.nLeadingNonStarters()
	if *ss += streamSafe(n); *ss > maxNonStarters {
		*ss = 0
		return ssOverflow
	}
	
	
	
	
	
	
	
	if n == 0 {
		*ss = streamSafe(p.nTrailingNonStarters())
		return ssStarter
	}
	return ssSuccess
}





func (ss *streamSafe) backwards(p Properties) ssState {
	if *ss > maxNonStarters {
		panic("streamSafe was not reset")
	}
	c := *ss + streamSafe(p.nTrailingNonStarters())
	if c > maxNonStarters {
		return ssOverflow
	}
	*ss = c
	if p.nLeadingNonStarters() == 0 {
		return ssStarter
	}
	return ssSuccess
}

func (ss streamSafe) isMax() bool {
	return ss == maxNonStarters
}


const GraphemeJoiner = "\u034F"






type reorderBuffer struct {
	rune  [maxBufferSize]Properties 
	byte  [maxByteBufferSize]byte   
	nbyte uint8                     
	ss    streamSafe                
	nrune int                       
	f     formInfo

	src      input
	nsrc     int
	tmpBytes input

	out    []byte
	flushF func(*reorderBuffer) bool
}

func (rb *reorderBuffer) init(f Form, src []byte) {
	rb.f = *formTable[f]
	rb.src.setBytes(src)
	rb.nsrc = len(src)
	rb.ss = 0
}

func (rb *reorderBuffer) initString(f Form, src string) {
	rb.f = *formTable[f]
	rb.src.setString(src)
	rb.nsrc = len(src)
	rb.ss = 0
}

func (rb *reorderBuffer) setFlusher(out []byte, f func(*reorderBuffer) bool) {
	rb.out = out
	rb.flushF = f
}


func (rb *reorderBuffer) reset() {
	rb.nrune = 0
	rb.nbyte = 0
}

func (rb *reorderBuffer) doFlush() bool {
	if rb.f.composing {
		rb.compose()
	}
	res := rb.flushF(rb)
	rb.reset()
	return res
}


func appendFlush(rb *reorderBuffer) bool {
	for i := 0; i < rb.nrune; i++ {
		start := rb.rune[i].pos
		end := start + rb.rune[i].size
		rb.out = append(rb.out, rb.byte[start:end]...)
	}
	return true
}


func (rb *reorderBuffer) flush(out []byte) []byte {
	for i := 0; i < rb.nrune; i++ {
		start := rb.rune[i].pos
		end := start + rb.rune[i].size
		out = append(out, rb.byte[start:end]...)
	}
	rb.reset()
	return out
}



func (rb *reorderBuffer) flushCopy(buf []byte) int {
	p := 0
	for i := 0; i < rb.nrune; i++ {
		runep := rb.rune[i]
		p += copy(buf[p:], rb.byte[runep.pos:runep.pos+runep.size])
	}
	rb.reset()
	return p
}




func (rb *reorderBuffer) insertOrdered(info Properties) {
	n := rb.nrune
	b := rb.rune[:]
	cc := info.ccc
	if cc > 0 {
		
		for ; n > 0; n-- {
			if b[n-1].ccc <= cc {
				break
			}
			b[n] = b[n-1]
		}
	}
	rb.nrune += 1
	pos := uint8(rb.nbyte)
	rb.nbyte += utf8.UTFMax
	info.pos = pos
	b[n] = info
}



type insertErr int

const (
	iSuccess insertErr = -iota
	iShortDst
	iShortSrc
)





func (rb *reorderBuffer) insertFlush(src input, i int, info Properties) insertErr {
	if rune := src.hangul(i); rune != 0 {
		rb.decomposeHangul(rune)
		return iSuccess
	}
	if info.hasDecomposition() {
		return rb.insertDecomposed(info.Decomposition())
	}
	rb.insertSingle(src, i, info)
	return iSuccess
}





func (rb *reorderBuffer) insertUnsafe(src input, i int, info Properties) {
	if rune := src.hangul(i); rune != 0 {
		rb.decomposeHangul(rune)
	}
	if info.hasDecomposition() {
		
		rb.insertDecomposed(info.Decomposition())
	} else {
		rb.insertSingle(src, i, info)
	}
}




func (rb *reorderBuffer) insertDecomposed(dcomp []byte) insertErr {
	rb.tmpBytes.setBytes(dcomp)
	
	
	
	for i := 0; i < len(dcomp); {
		info := rb.f.info(rb.tmpBytes, i)
		if info.BoundaryBefore() && rb.nrune > 0 && !rb.doFlush() {
			return iShortDst
		}
		i += copy(rb.byte[rb.nbyte:], dcomp[i:i+int(info.size)])
		rb.insertOrdered(info)
	}
	return iSuccess
}



func (rb *reorderBuffer) insertSingle(src input, i int, info Properties) {
	src.copySlice(rb.byte[rb.nbyte:], i, i+int(info.size))
	rb.insertOrdered(info)
}


func (rb *reorderBuffer) insertCGJ() {
	rb.insertSingle(input{str: GraphemeJoiner}, 0, Properties{size: uint8(len(GraphemeJoiner))})
}


func (rb *reorderBuffer) appendRune(r rune) {
	bn := rb.nbyte
	sz := utf8.EncodeRune(rb.byte[bn:], rune(r))
	rb.nbyte += utf8.UTFMax
	rb.rune[rb.nrune] = Properties{pos: bn, size: uint8(sz)}
	rb.nrune++
}


func (rb *reorderBuffer) assignRune(pos int, r rune) {
	bn := rb.rune[pos].pos
	sz := utf8.EncodeRune(rb.byte[bn:], rune(r))
	rb.rune[pos] = Properties{pos: bn, size: uint8(sz)}
}


func (rb *reorderBuffer) runeAt(n int) rune {
	inf := rb.rune[n]
	r, _ := utf8.DecodeRune(rb.byte[inf.pos : inf.pos+inf.size])
	return r
}



func (rb *reorderBuffer) bytesAt(n int) []byte {
	inf := rb.rune[n]
	return rb.byte[inf.pos : int(inf.pos)+int(inf.size)]
}


const (
	hangulBase  = 0xAC00 
	hangulBase0 = 0xEA
	hangulBase1 = 0xB0
	hangulBase2 = 0x80

	hangulEnd  = hangulBase + jamoLVTCount 
	hangulEnd0 = 0xED
	hangulEnd1 = 0x9E
	hangulEnd2 = 0xA4

	jamoLBase  = 0x1100 
	jamoLBase0 = 0xE1
	jamoLBase1 = 0x84
	jamoLEnd   = 0x1113
	jamoVBase  = 0x1161
	jamoVEnd   = 0x1176
	jamoTBase  = 0x11A7
	jamoTEnd   = 0x11C3

	jamoTCount   = 28
	jamoVCount   = 21
	jamoVTCount  = 21 * 28
	jamoLVTCount = 19 * 21 * 28
)

const hangulUTF8Size = 3

func isHangul(b []byte) bool {
	if len(b) < hangulUTF8Size {
		return false
	}
	b0 := b[0]
	if b0 < hangulBase0 {
		return false
	}
	b1 := b[1]
	switch {
	case b0 == hangulBase0:
		return b1 >= hangulBase1
	case b0 < hangulEnd0:
		return true
	case b0 > hangulEnd0:
		return false
	case b1 < hangulEnd1:
		return true
	}
	return b1 == hangulEnd1 && b[2] < hangulEnd2
}

func isHangulString(b string) bool {
	if len(b) < hangulUTF8Size {
		return false
	}
	b0 := b[0]
	if b0 < hangulBase0 {
		return false
	}
	b1 := b[1]
	switch {
	case b0 == hangulBase0:
		return b1 >= hangulBase1
	case b0 < hangulEnd0:
		return true
	case b0 > hangulEnd0:
		return false
	case b1 < hangulEnd1:
		return true
	}
	return b1 == hangulEnd1 && b[2] < hangulEnd2
}


func isJamoVT(b []byte) bool {
	
	return b[0] == jamoLBase0 && (b[1]&0xFC) == jamoLBase1
}

func isHangulWithoutJamoT(b []byte) bool {
	c, _ := utf8.DecodeRune(b)
	c -= hangulBase
	return c < jamoLVTCount && c%jamoTCount == 0
}



func decomposeHangul(buf []byte, r rune) int {
	const JamoUTF8Len = 3
	r -= hangulBase
	x := r % jamoTCount
	r /= jamoTCount
	utf8.EncodeRune(buf, jamoLBase+r/jamoVCount)
	utf8.EncodeRune(buf[JamoUTF8Len:], jamoVBase+r%jamoVCount)
	if x != 0 {
		utf8.EncodeRune(buf[2*JamoUTF8Len:], jamoTBase+x)
		return 3 * JamoUTF8Len
	}
	return 2 * JamoUTF8Len
}




func (rb *reorderBuffer) decomposeHangul(r rune) {
	r -= hangulBase
	x := r % jamoTCount
	r /= jamoTCount
	rb.appendRune(jamoLBase + r/jamoVCount)
	rb.appendRune(jamoVBase + r%jamoVCount)
	if x != 0 {
		rb.appendRune(jamoTBase + x)
	}
}



func (rb *reorderBuffer) combineHangul(s, i, k int) {
	b := rb.rune[:]
	bn := rb.nrune
	for ; i < bn; i++ {
		cccB := b[k-1].ccc
		cccC := b[i].ccc
		if cccB == 0 {
			s = k - 1
		}
		if s != k-1 && cccB >= cccC {
			
			b[k] = b[i]
			k++
		} else {
			l := rb.runeAt(s) 
			v := rb.runeAt(i) 
			switch {
			case jamoLBase <= l && l < jamoLEnd &&
				jamoVBase <= v && v < jamoVEnd:
				
				rb.assignRune(s, hangulBase+
					(l-jamoLBase)*jamoVTCount+(v-jamoVBase)*jamoTCount)
			case hangulBase <= l && l < hangulEnd &&
				jamoTBase < v && v < jamoTEnd &&
				((l-hangulBase)%jamoTCount) == 0:
				
				rb.assignRune(s, l+v-jamoTBase)
			default:
				b[k] = b[i]
				k++
			}
		}
	}
	rb.nrune = k
}




func (rb *reorderBuffer) compose() {
	
	
	
	
	
	bn := rb.nrune
	if bn == 0 {
		return
	}
	k := 1
	b := rb.rune[:]
	for s, i := 0, 1; i < bn; i++ {
		if isJamoVT(rb.bytesAt(i)) {
			
			
			rb.combineHangul(s, i, k)
			return
		}
		ii := b[i]
		
		
		
		
		if ii.combinesBackward() {
			cccB := b[k-1].ccc
			cccC := ii.ccc
			blocked := false 
			if cccB == 0 {
				s = k - 1
			} else {
				blocked = s != k-1 && cccB >= cccC
			}
			if !blocked {
				combined := combine(rb.runeAt(s), rb.runeAt(i))
				if combined != 0 {
					rb.assignRune(s, combined)
					continue
				}
			}
		}
		b[k] = b[i]
		k++
	}
	rb.nrune = k
}
