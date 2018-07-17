





package rangetable

import (
	"sort"
	"unicode"
)


func New(r ...rune) *unicode.RangeTable {
	if len(r) == 0 {
		return &unicode.RangeTable{}
	}

	sort.Sort(byRune(r))

	
	k := 1
	for i := 1; i < len(r); i++ {
		if r[k-1] != r[i] {
			r[k] = r[i]
			k++
		}
	}

	var rt unicode.RangeTable
	for _, r := range r[:k] {
		if r <= 0xFFFF {
			rt.R16 = append(rt.R16, unicode.Range16{Lo: uint16(r), Hi: uint16(r), Stride: 1})
		} else {
			rt.R32 = append(rt.R32, unicode.Range32{Lo: uint32(r), Hi: uint32(r), Stride: 1})
		}
	}

	
	return Merge(&rt)
}

type byRune []rune

func (r byRune) Len() int           { return len(r) }
func (r byRune) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r byRune) Less(i, j int) bool { return r[i] < r[j] }


func Visit(rt *unicode.RangeTable, fn func(rune)) {
	for _, r16 := range rt.R16 {
		for r := rune(r16.Lo); r <= rune(r16.Hi); r += rune(r16.Stride) {
			fn(r)
		}
	}
	for _, r32 := range rt.R32 {
		for r := rune(r32.Lo); r <= rune(r32.Hi); r += rune(r32.Stride) {
			fn(r)
		}
	}
}





func Assigned(version string) *unicode.RangeTable {
	return assigned[version]
}
