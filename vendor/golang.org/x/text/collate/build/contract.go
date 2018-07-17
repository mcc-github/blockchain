



package build

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/text/internal/colltab"
)






























const (
	final   = 0
	noIndex = 0xFF
)























type ctEntry struct {
	L uint8 
	H uint8 
	N uint8 
	I uint8 
}



type contractTrieSet []struct{ l, h, n, i uint8 }



type ctHandle struct {
	index, n int
}



func appendTrie(ct *colltab.ContractTrieSet, suffixes []string) (ctHandle, error) {
	es := make([]stridx, len(suffixes))
	for i, s := range suffixes {
		es[i].str = s
	}
	sort.Sort(offsetSort(es))
	for i := range es {
		es[i].index = i + 1
	}
	sort.Sort(genidxSort(es))
	i := len(*ct)
	n, err := genStates(ct, es)
	if err != nil {
		*ct = (*ct)[:i]
		return ctHandle{}, err
	}
	return ctHandle{i, n}, nil
}



func genStates(ct *colltab.ContractTrieSet, sis []stridx) (int, error) {
	if len(sis) == 0 {
		return 0, fmt.Errorf("genStates: list of suffices must be non-empty")
	}
	start := len(*ct)
	
	for _, si := range sis {
		s := si.str
		if len(s) == 0 {
			continue
		}
		added := false
		c := s[0]
		if len(s) > 1 {
			for j := len(*ct) - 1; j >= start; j-- {
				if (*ct)[j].L == c {
					added = true
					break
				}
			}
			if !added {
				*ct = append(*ct, ctEntry{L: c, I: noIndex})
			}
		} else {
			for j := len(*ct) - 1; j >= start; j-- {
				
				if (*ct)[j].L == c {
					(*ct)[j].I = uint8(si.index)
					added = true
				}
				
				if (*ct)[j].H+1 == c {
					(*ct)[j].H = c
					added = true
				}
			}
			if !added {
				*ct = append(*ct, ctEntry{L: c, H: c, N: final, I: uint8(si.index)})
			}
		}
	}
	n := len(*ct) - start
	
	sp := 0
	for i, end := start, len(*ct); i < end; i++ {
		fe := (*ct)[i]
		if fe.H == 0 { 
			ln := len(*ct) - start - n
			if ln > 0xFF {
				return 0, fmt.Errorf("genStates: relative block offset too large: %d > 255", ln)
			}
			fe.H = uint8(ln)
			
			for ; sis[sp].str[0] != fe.L; sp++ {
			}
			se := sp + 1
			for ; se < len(sis) && len(sis[se].str) > 1 && sis[se].str[0] == fe.L; se++ {
			}
			sl := sis[sp:se]
			sp = se
			for i, si := range sl {
				sl[i].str = si.str[1:]
			}
			nn, err := genStates(ct, sl)
			if err != nil {
				return 0, err
			}
			fe.N = uint8(nn)
			(*ct)[i] = fe
		}
	}
	sort.Sort(entrySort((*ct)[start : start+n]))
	return n, nil
}




type entrySort colltab.ContractTrieSet

func (fe entrySort) Len() int      { return len(fe) }
func (fe entrySort) Swap(i, j int) { fe[i], fe[j] = fe[j], fe[i] }
func (fe entrySort) Less(i, j int) bool {
	return fe[i].L > fe[j].L
}


type stridx struct {
	str   string
	index int
}





type offsetSort []stridx

func (si offsetSort) Len() int      { return len(si) }
func (si offsetSort) Swap(i, j int) { si[i], si[j] = si[j], si[i] }
func (si offsetSort) Less(i, j int) bool {
	if len(si[i].str) != len(si[j].str) {
		return len(si[i].str) > len(si[j].str)
	}
	return si[i].str < si[j].str
}



type genidxSort []stridx

func (si genidxSort) Len() int      { return len(si) }
func (si genidxSort) Swap(i, j int) { si[i], si[j] = si[j], si[i] }
func (si genidxSort) Less(i, j int) bool {
	if strings.HasPrefix(si[j].str, si[i].str) {
		return false
	}
	if strings.HasPrefix(si[i].str, si[j].str) {
		return true
	}
	return si[i].str < si[j].str
}



func lookup(ct *colltab.ContractTrieSet, h ctHandle, str []byte) (index, ns int) {
	states := (*ct)[h.index:]
	p := 0
	n := h.n
	for i := 0; i < n && p < len(str); {
		e := states[i]
		c := str[p]
		if c >= e.L {
			if e.L == c {
				p++
				if e.I != noIndex {
					index, ns = int(e.I), p
				}
				if e.N != final {
					
					i, states, n = 0, states[int(e.H)+n:], int(e.N)
				} else {
					return
				}
				continue
			} else if e.N == final && c <= e.H {
				p++
				return int(c-e.L) + int(e.I), p
			}
		}
		i++
	}
	return
}



func print(t *colltab.ContractTrieSet, w io.Writer, name string) (n, size int, err error) {
	update3 := func(nn, sz int, e error) {
		n += nn
		if err == nil {
			err = e
		}
		size += sz
	}
	update2 := func(nn int, e error) { update3(nn, 0, e) }

	update3(printArray(*t, w, name))
	update2(fmt.Fprintf(w, "var %sContractTrieSet = ", name))
	update3(printStruct(*t, w, name))
	update2(fmt.Fprintln(w))
	return
}

func printArray(ct colltab.ContractTrieSet, w io.Writer, name string) (n, size int, err error) {
	p := func(f string, a ...interface{}) {
		nn, e := fmt.Fprintf(w, f, a...)
		n += nn
		if err == nil {
			err = e
		}
	}
	size = len(ct) * 4
	p("
	p("var %sCTEntries = [%d]struct{L,H,N,I uint8}{\n", name, len(ct))
	for _, fe := range ct {
		p("\t{0x%X, 0x%X, %d, %d},\n", fe.L, fe.H, fe.N, fe.I)
	}
	p("}\n")
	return
}

func printStruct(ct colltab.ContractTrieSet, w io.Writer, name string) (n, size int, err error) {
	n, err = fmt.Fprintf(w, "colltab.ContractTrieSet( %sCTEntries[:] )", name)
	size = int(reflect.TypeOf(ct).Size())
	return
}
