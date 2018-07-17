



package collate

import (
	"bytes"
	"sort"
)

const (
	maxSortBuffer  = 40960
	maxSortEntries = 4096
)

type swapper interface {
	Swap(i, j int)
}

type sorter struct {
	buf  *Buffer
	keys [][]byte
	src  swapper
}

func (s *sorter) init(n int) {
	if s.buf == nil {
		s.buf = &Buffer{}
		s.buf.init()
	}
	if cap(s.keys) < n {
		s.keys = make([][]byte, n)
	}
	s.keys = s.keys[0:n]
}

func (s *sorter) sort(src swapper) {
	s.src = src
	sort.Sort(s)
}

func (s sorter) Len() int {
	return len(s.keys)
}

func (s sorter) Less(i, j int) bool {
	return bytes.Compare(s.keys[i], s.keys[j]) == -1
}

func (s sorter) Swap(i, j int) {
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
	s.src.Swap(i, j)
}


type Lister interface {
	Len() int
	Swap(i, j int)
	
	Bytes(i int) []byte
}


func (c *Collator) Sort(x Lister) {
	n := x.Len()
	c.sorter.init(n)
	for i := 0; i < n; i++ {
		c.sorter.keys[i] = c.Key(c.sorter.buf, x.Bytes(i))
	}
	c.sorter.sort(x)
}


func (c *Collator) SortStrings(x []string) {
	c.sorter.init(len(x))
	for i, s := range x {
		c.sorter.keys[i] = c.KeyFromString(c.sorter.buf, s)
	}
	c.sorter.sort(sort.StringSlice(x))
}
