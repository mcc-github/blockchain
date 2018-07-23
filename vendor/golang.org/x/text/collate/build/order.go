



package build

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/unicode/norm"
)

type logicalAnchor int

const (
	firstAnchor logicalAnchor = -1
	noAnchor                  = 0
	lastAnchor                = 1
)





type entry struct {
	str    string 
	runes  []rune
	elems  []rawCE 
	extend string  
	before bool    
	lock   bool    

	
	prev, next *entry
	level      colltab.Level 
	skipRemove bool          

	decompose bool 
	exclude   bool 
	implicit  bool 
	modified  bool 
	logical   logicalAnchor

	expansionIndex    int 
	contractionHandle ctHandle
	contractionIndex  int 
}

func (e *entry) String() string {
	return fmt.Sprintf("%X (%q) -> %X (ch:%x; ci:%d, ei:%d)",
		e.runes, e.str, e.elems, e.contractionHandle, e.contractionIndex, e.expansionIndex)
}

func (e *entry) skip() bool {
	return e.contraction()
}

func (e *entry) expansion() bool {
	return !e.decompose && len(e.elems) > 1
}

func (e *entry) contraction() bool {
	return len(e.runes) > 1
}

func (e *entry) contractionStarter() bool {
	return e.contractionHandle.n != 0
}






func (e *entry) nextIndexed() (*entry, colltab.Level) {
	level := e.level
	for e = e.next; e != nil && (e.exclude || len(e.elems) == 0); e = e.next {
		if e.level < level {
			level = e.level
		}
	}
	return e, level
}





func (e *entry) remove() {
	if e.logical != noAnchor {
		log.Fatalf("may not remove anchor %q", e.str)
	}
	
	e.elems = nil
	if !e.skipRemove {
		if e.prev != nil {
			e.prev.next = e.next
		}
		if e.next != nil {
			e.next.prev = e.prev
		}
	}
	e.skipRemove = false
}


func (e *entry) insertAfter(n *entry) {
	if e == n {
		panic("e == anchor")
	}
	if e == nil {
		panic("unexpected nil anchor")
	}
	n.remove()
	n.decompose = false 

	n.next = e.next
	n.prev = e
	if e.next != nil {
		e.next.prev = n
	}
	e.next = n
}


func (e *entry) insertBefore(n *entry) {
	if e == n {
		panic("e == anchor")
	}
	if e == nil {
		panic("unexpected nil anchor")
	}
	n.remove()
	n.decompose = false 

	n.prev = e.prev
	n.next = e
	if e.prev != nil {
		e.prev.next = n
	}
	e.prev = n
}

func (e *entry) encodeBase() (ce uint32, err error) {
	switch {
	case e.expansion():
		ce, err = makeExpandIndex(e.expansionIndex)
	default:
		if e.decompose {
			log.Fatal("decompose should be handled elsewhere")
		}
		ce, err = makeCE(e.elems[0])
	}
	return
}

func (e *entry) encode() (ce uint32, err error) {
	if e.skip() {
		log.Fatal("cannot build colElem for entry that should be skipped")
	}
	switch {
	case e.decompose:
		t1 := e.elems[0].w[2]
		t2 := 0
		if len(e.elems) > 1 {
			t2 = e.elems[1].w[2]
		}
		ce, err = makeDecompose(t1, t2)
	case e.contractionStarter():
		ce, err = makeContractIndex(e.contractionHandle, e.contractionIndex)
	default:
		if len(e.runes) > 1 {
			log.Fatal("colElem: contractions are handled in contraction trie")
		}
		ce, err = e.encodeBase()
	}
	return
}


func entryLess(a, b *entry) bool {
	if res, _ := compareWeights(a.elems, b.elems); res != 0 {
		return res == -1
	}
	if a.logical != noAnchor {
		return a.logical == firstAnchor
	}
	if b.logical != noAnchor {
		return b.logical == lastAnchor
	}
	return a.str < b.str
}

type sortedEntries []*entry

func (s sortedEntries) Len() int {
	return len(s)
}

func (s sortedEntries) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortedEntries) Less(i, j int) bool {
	return entryLess(s[i], s[j])
}

type ordering struct {
	id       string
	entryMap map[string]*entry
	ordered  []*entry
	handle   *trieHandle
}




func (o *ordering) insert(e *entry) {
	if e.logical == noAnchor {
		o.entryMap[e.str] = e
	} else {
		
		o.entryMap[fmt.Sprintf("[%s]", e.str)] = e
		
		o.entryMap[fmt.Sprintf("<%s/>", strings.Replace(e.str, " ", "_", -1))] = e
	}
	o.ordered = append(o.ordered, e)
}



func (o *ordering) newEntry(s string, ces []rawCE) *entry {
	e := &entry{
		runes: []rune(s),
		elems: ces,
		str:   s,
	}
	o.insert(e)
	return e
}




func (o *ordering) find(str string) *entry {
	e := o.entryMap[str]
	if e == nil {
		r := []rune(str)
		if len(r) == 1 {
			const (
				firstHangul = 0xAC00
				lastHangul  = 0xD7A3
			)
			if r[0] >= firstHangul && r[0] <= lastHangul {
				ce := []rawCE{}
				nfd := norm.NFD.String(str)
				for _, r := range nfd {
					ce = append(ce, o.find(string(r)).elems...)
				}
				e = o.newEntry(nfd, ce)
			} else {
				e = o.newEntry(string(r[0]), []rawCE{
					{w: []int{
						implicitPrimary(r[0]),
						defaultSecondary,
						defaultTertiary,
						int(r[0]),
					},
					},
				})
				e.modified = true
			}
			e.exclude = true 
		}
	}
	return e
}






func makeRootOrdering() ordering {
	const max = unicode.MaxRune
	o := ordering{
		entryMap: make(map[string]*entry),
	}
	insert := func(typ logicalAnchor, s string, ce []int) {
		e := &entry{
			elems:   []rawCE{{w: ce}},
			str:     s,
			exclude: true,
			logical: typ,
		}
		o.insert(e)
	}
	insert(firstAnchor, "first tertiary ignorable", []int{0, 0, 0, 0})
	insert(lastAnchor, "last tertiary ignorable", []int{0, 0, 0, max})
	insert(lastAnchor, "last primary ignorable", []int{0, defaultSecondary, defaultTertiary, max})
	insert(lastAnchor, "last non ignorable", []int{maxPrimary, defaultSecondary, defaultTertiary, max})
	insert(lastAnchor, "__END__", []int{1 << maxPrimaryBits, defaultSecondary, defaultTertiary, max})
	return o
}





func (o *ordering) patchForInsert() {
	for i := 0; i < len(o.ordered)-1; {
		e := o.ordered[i]
		lev := e.level
		n := e.next
		for ; n != nil && len(n.elems) > 1; n = n.next {
			if n.level < lev {
				lev = n.level
			}
			n.skipRemove = true
		}
		for ; o.ordered[i] != n; i++ {
			o.ordered[i].level = lev
			o.ordered[i].next = n
			o.ordered[i+1].prev = e
		}
	}
}


func (o *ordering) clone() *ordering {
	o.sort()
	oo := ordering{
		entryMap: make(map[string]*entry),
	}
	for _, e := range o.ordered {
		ne := &entry{
			runes:     e.runes,
			elems:     e.elems,
			str:       e.str,
			decompose: e.decompose,
			exclude:   e.exclude,
			logical:   e.logical,
		}
		oo.insert(ne)
	}
	oo.sort() 
	oo.patchForInsert()
	return &oo
}



func (o *ordering) front() *entry {
	e := o.ordered[0]
	if e.prev != nil {
		log.Panicf("unexpected first entry: %v", e)
	}
	
	e, _ = e.nextIndexed()
	return e
}



func (o *ordering) sort() {
	sort.Sort(sortedEntries(o.ordered))
	l := o.ordered
	for i := 1; i < len(l); i++ {
		k := i - 1
		l[k].next = l[i]
		_, l[k].level = compareWeights(l[k].elems, l[i].elems)
		l[i].prev = l[k]
	}
}



func (o *ordering) genColElems(str string) []rawCE {
	elems := []rawCE{}
	for _, r := range []rune(str) {
		for _, ce := range o.find(string(r)).elems {
			if ce.w[0] != 0 || ce.w[1] != 0 || ce.w[2] != 0 {
				elems = append(elems, ce)
			}
		}
	}
	return elems
}