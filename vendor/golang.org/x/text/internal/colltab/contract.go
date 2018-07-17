



package colltab

import "unicode/utf8"



type ContractTrieSet []struct{ L, H, N, I uint8 }









type ctScanner struct {
	states ContractTrieSet
	s      []byte
	n      int
	index  int
	pindex int
	done   bool
}

type ctScannerString struct {
	states ContractTrieSet
	s      string
	n      int
	index  int
	pindex int
	done   bool
}

func (t ContractTrieSet) scanner(index, n int, b []byte) ctScanner {
	return ctScanner{s: b, states: t[index:], n: n}
}

func (t ContractTrieSet) scannerString(index, n int, str string) ctScannerString {
	return ctScannerString{s: str, states: t[index:], n: n}
}



func (s *ctScanner) result() (i, p int) {
	return s.index, s.pindex
}

func (s *ctScannerString) result() (i, p int) {
	return s.index, s.pindex
}

const (
	final   = 0
	noIndex = 0xFF
)



func (s *ctScanner) scan(p int) int {
	pr := p 
	str := s.s
	states, n := s.states, s.n
	for i := 0; i < n && p < len(str); {
		e := states[i]
		c := str[p]
		
		
		
		
		if c >= e.L {
			if e.L == c {
				p++
				if e.I != noIndex {
					s.index = int(e.I)
					s.pindex = p
				}
				if e.N != final {
					i, states, n = 0, states[int(e.H)+n:], int(e.N)
					if p >= len(str) || utf8.RuneStart(str[p]) {
						s.states, s.n, pr = states, n, p
					}
				} else {
					s.done = true
					return p
				}
				continue
			} else if e.N == final && c <= e.H {
				p++
				s.done = true
				s.index = int(c-e.L) + int(e.I)
				s.pindex = p
				return p
			}
		}
		i++
	}
	return pr
}


func (s *ctScannerString) scan(p int) int {
	pr := p 
	str := s.s
	states, n := s.states, s.n
	for i := 0; i < n && p < len(str); {
		e := states[i]
		c := str[p]
		
		
		
		
		if c >= e.L {
			if e.L == c {
				p++
				if e.I != noIndex {
					s.index = int(e.I)
					s.pindex = p
				}
				if e.N != final {
					i, states, n = 0, states[int(e.H)+n:], int(e.N)
					if p >= len(str) || utf8.RuneStart(str[p]) {
						s.states, s.n, pr = states, n, p
					}
				} else {
					s.done = true
					return p
				}
				continue
			} else if e.N == final && c <= e.H {
				p++
				s.done = true
				s.index = int(c-e.L) + int(e.I)
				s.pindex = p
				return p
			}
		}
		i++
	}
	return pr
}
