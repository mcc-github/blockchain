





package comparer

import "bytes"

type bytesComparer struct{}

func (bytesComparer) Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

func (bytesComparer) Name() string {
	return "leveldb.BytewiseComparator"
}

func (bytesComparer) Separator(dst, a, b []byte) []byte {
	i, n := 0, len(a)
	if n > len(b) {
		n = len(b)
	}
	for ; i < n && a[i] == b[i]; i++ {
	}
	if i >= n {
		
	} else if c := a[i]; c < 0xff && c+1 < b[i] {
		dst = append(dst, a[:i+1]...)
		dst[i]++
		return dst
	}
	return nil
}

func (bytesComparer) Successor(dst, b []byte) []byte {
	for i, c := range b {
		if c != 0xff {
			dst = append(dst, b[:i+1]...)
			dst[i]++
			return dst
		}
	}
	return nil
}



var DefaultComparer = bytesComparer{}
