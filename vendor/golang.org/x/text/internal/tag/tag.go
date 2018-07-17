




package tag 

import "sort"





type Index string


func (s Index) Elem(x int) string {
	return string(s[x*4 : x*4+4])
}




func (s Index) Index(key []byte) int {
	n := len(key)
	
	
	index := sort.Search(len(s)/4, func(i int) bool {
		return cmp(s[i*4:i*4+n], key) != -1
	})
	i := index * 4
	if cmp(s[i:i+len(key)], key) != 0 {
		return -1
	}
	return index
}



func (s Index) Next(key []byte, x int) int {
	if x++; x*4 < len(s) && cmp(s[x*4:x*4+len(key)], key) == 0 {
		return x
	}
	return -1
}


func cmp(a Index, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i, c := range b[:n] {
		switch {
		case a[i] > c:
			return 1
		case a[i] < c:
			return -1
		}
	}
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	}
	return 0
}


func Compare(a string, b []byte) int {
	return cmp(Index(a), b)
}



func FixCase(form string, b []byte) bool {
	if len(form) != len(b) {
		return false
	}
	for i, c := range b {
		if form[i] <= 'Z' {
			if c >= 'a' {
				c -= 'z' - 'Z'
			}
			if c < 'A' || 'Z' < c {
				return false
			}
		} else {
			if c <= 'Z' {
				c += 'z' - 'Z'
			}
			if c < 'a' || 'z' < c {
				return false
			}
		}
		b[i] = c
	}
	return true
}
