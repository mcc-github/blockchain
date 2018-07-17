



package build

import (
	"fmt"
	"io"
	"reflect"

	"golang.org/x/text/internal/colltab"
)


type table struct {
	colltab.Table
	trie trie
	root *trieHandle
}




func (t *table) fprint(w io.Writer, name string) (n, size int, err error) {
	update := func(nn, sz int, e error) {
		n += nn
		if err == nil {
			err = e
		}
		size += sz
	}
	
	update(printColElems(w, t.ExpandElem, name+"ExpandElem"))
	update(printColElems(w, t.ContractElem, name+"ContractElem"))
	update(t.trie.printArrays(w, name))
	update(printArray(t.ContractTries, w, name))

	nn, e := fmt.Fprintf(w, "
	update(nn, 0, e)
	return
}

func (t *table) fprintIndex(w io.Writer, h *trieHandle, id string) (n int, err error) {
	p := func(f string, a ...interface{}) {
		nn, e := fmt.Fprintf(w, f, a...)
		n += nn
		if err == nil {
			err = e
		}
	}
	p("\t{ 
	p("\t\tlookupOffset: 0x%x,\n", h.lookupStart)
	p("\t\tvaluesOffset: 0x%x,\n", h.valueStart)
	p("\t},\n")
	return
}

func printColElems(w io.Writer, a []uint32, name string) (n, sz int, err error) {
	p := func(f string, a ...interface{}) {
		nn, e := fmt.Fprintf(w, f, a...)
		n += nn
		if err == nil {
			err = e
		}
	}
	sz = len(a) * int(reflect.TypeOf(uint32(0)).Size())
	p("
	p("var %s = [%d]uint32 {", name, len(a))
	for i, c := range a {
		switch {
		case i%64 == 0:
			p("\n\t
		case (i%64)%6 == 0:
			p("\n\t")
		}
		p("0x%.8X, ", c)
	}
	p("\n}\n\n")
	return
}
