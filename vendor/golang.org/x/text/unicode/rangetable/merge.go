



package rangetable

import (
	"unicode"
)


const atEnd = unicode.MaxRune + 1









func Merge(ranges ...*unicode.RangeTable) *unicode.RangeTable {
	rt := &unicode.RangeTable{}
	if len(ranges) == 0 {
		return rt
	}

	iter := tablesIter(make([]tableIndex, len(ranges)))

	for i, t := range ranges {
		iter[i] = tableIndex{t, 0, atEnd}
		if len(t.R16) > 0 {
			iter[i].next = rune(t.R16[0].Lo)
		}
	}

	if r0 := iter.next16(); r0.Stride != 0 {
		for {
			r1 := iter.next16()
			if r1.Stride == 0 {
				rt.R16 = append(rt.R16, r0)
				break
			}
			stride := r1.Lo - r0.Hi
			if (r1.Lo == r1.Hi || stride == r1.Stride) && (r0.Lo == r0.Hi || stride == r0.Stride) {
				
				r0.Hi, r0.Stride = r1.Hi, stride
				continue
			} else if stride == r0.Stride {
				
				
				r0.Hi = r1.Lo
				r0.Stride = stride
				r1.Lo = r1.Lo + r1.Stride
				if r1.Lo > r1.Hi {
					continue
				}
			}
			rt.R16 = append(rt.R16, r0)
			r0 = r1
		}
	}

	for i, t := range ranges {
		iter[i] = tableIndex{t, 0, atEnd}
		if len(t.R32) > 0 {
			iter[i].next = rune(t.R32[0].Lo)
		}
	}

	if r0 := iter.next32(); r0.Stride != 0 {
		for {
			r1 := iter.next32()
			if r1.Stride == 0 {
				rt.R32 = append(rt.R32, r0)
				break
			}
			stride := r1.Lo - r0.Hi
			if (r1.Lo == r1.Hi || stride == r1.Stride) && (r0.Lo == r0.Hi || stride == r0.Stride) {
				
				r0.Hi, r0.Stride = r1.Hi, stride
				continue
			} else if stride == r0.Stride {
				
				
				r0.Hi = r1.Lo
				r1.Lo = r1.Lo + r1.Stride
				if r1.Lo > r1.Hi {
					continue
				}
			}
			rt.R32 = append(rt.R32, r0)
			r0 = r1
		}
	}

	for i := 0; i < len(rt.R16) && rt.R16[i].Hi <= unicode.MaxLatin1; i++ {
		rt.LatinOffset = i + 1
	}

	return rt
}

type tableIndex struct {
	t    *unicode.RangeTable
	p    uint32
	next rune
}

type tablesIter []tableIndex



func sortIter(t []tableIndex) {
	for i := range t {
		for j := i; j > 0 && t[j-1].next > t[j].next; j-- {
			t[j], t[j-1] = t[j-1], t[j]
		}
	}
}





func (ti tablesIter) next16() unicode.Range16 {
	sortIter(ti)

	t0 := ti[0]
	if t0.next == atEnd {
		return unicode.Range16{}
	}
	r0 := t0.t.R16[t0.p]
	r0.Lo = uint16(t0.next)

	
	for i := range ti {
		tn := ti[i]
		
		
		
		if rune(r0.Hi) <= tn.next {
			break
		}
		rn := tn.t.R16[tn.p]
		rn.Lo = uint16(tn.next)

		
		
		m := (rn.Lo - r0.Lo) % r0.Stride
		if m == 0 && (rn.Stride == r0.Stride || rn.Lo == rn.Hi) {
			
			
			if r0.Hi > rn.Hi {
				r0.Hi = rn.Hi
			}
		} else {
			
			
			if x := rn.Lo - m; r0.Lo <= x {
				r0.Hi = x
			}
			break
		}
	}

	
	for i := range ti {
		tn := &ti[i]
		if rune(r0.Hi) < tn.next {
			break
		}
		rn := tn.t.R16[tn.p]
		stride := rune(rn.Stride)
		tn.next += stride * (1 + ((rune(r0.Hi) - tn.next) / stride))
		if rune(rn.Hi) < tn.next {
			if tn.p++; int(tn.p) == len(tn.t.R16) {
				tn.next = atEnd
			} else {
				tn.next = rune(tn.t.R16[tn.p].Lo)
			}
		}
	}

	if r0.Lo == r0.Hi {
		r0.Stride = 1
	}

	return r0
}





func (ti tablesIter) next32() unicode.Range32 {
	sortIter(ti)

	t0 := ti[0]
	if t0.next == atEnd {
		return unicode.Range32{}
	}
	r0 := t0.t.R32[t0.p]
	r0.Lo = uint32(t0.next)

	
	for i := range ti {
		tn := ti[i]
		
		
		
		if rune(r0.Hi) <= tn.next {
			break
		}
		rn := tn.t.R32[tn.p]
		rn.Lo = uint32(tn.next)

		
		
		m := (rn.Lo - r0.Lo) % r0.Stride
		if m == 0 && (rn.Stride == r0.Stride || rn.Lo == rn.Hi) {
			
			
			if r0.Hi > rn.Hi {
				r0.Hi = rn.Hi
			}
		} else {
			
			
			if x := rn.Lo - m; r0.Lo <= x {
				r0.Hi = x
			}
			break
		}
	}

	
	for i := range ti {
		tn := &ti[i]
		if rune(r0.Hi) < tn.next {
			break
		}
		rn := tn.t.R32[tn.p]
		stride := rune(rn.Stride)
		tn.next += stride * (1 + ((rune(r0.Hi) - tn.next) / stride))
		if rune(rn.Hi) < tn.next {
			if tn.p++; int(tn.p) == len(tn.t.R32) {
				tn.next = atEnd
			} else {
				tn.next = rune(tn.t.R32[tn.p].Lo)
			}
		}
	}

	if r0.Lo == r0.Hi {
		r0.Stride = 1
	}

	return r0
}
