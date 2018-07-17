



package build 

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)


















type Builder struct {
	index  *trieBuilder
	root   ordering
	locale []*Tailoring
	t      *table
	err    error
	built  bool

	minNonVar int 
	varTop    int 

	
	expIndex map[string]int      
	ctHandle map[string]ctHandle 
	ctElem   map[string]int      
}






type Tailoring struct {
	id      string
	builder *Builder
	index   *ordering

	anchor *entry
	before bool
}


func NewBuilder() *Builder {
	return &Builder{
		index:    newTrieBuilder(),
		root:     makeRootOrdering(),
		expIndex: make(map[string]int),
		ctHandle: make(map[string]ctHandle),
		ctElem:   make(map[string]int),
	}
}



func (b *Builder) Tailoring(loc language.Tag) *Tailoring {
	t := &Tailoring{
		id:      loc.String(),
		builder: b,
		index:   b.root.clone(),
	}
	t.index.id = t.id
	b.locale = append(b.locale, t)
	return t
}










func (b *Builder) Add(runes []rune, colelems [][]int, variables []int) error {
	str := string(runes)
	elems := make([]rawCE, len(colelems))
	for i, ce := range colelems {
		if len(ce) == 0 {
			break
		}
		elems[i] = makeRawCE(ce, 0)
		if len(ce) == 1 {
			elems[i].w[1] = defaultSecondary
		}
		if len(ce) <= 2 {
			elems[i].w[2] = defaultTertiary
		}
		if len(ce) <= 3 {
			elems[i].w[3] = ce[0]
		}
	}
	for i, ce := range elems {
		p := ce.w[0]
		isvar := false
		for _, j := range variables {
			if i == j {
				isvar = true
			}
		}
		if isvar {
			if p >= b.minNonVar && b.minNonVar > 0 {
				return fmt.Errorf("primary value %X of variable is larger than the smallest non-variable %X", p, b.minNonVar)
			}
			if p > b.varTop {
				b.varTop = p
			}
		} else if p > 1 { 
			if p <= b.varTop {
				return fmt.Errorf("primary value %X of non-variable is smaller than the highest variable %X", p, b.varTop)
			}
			if b.minNonVar == 0 || p < b.minNonVar {
				b.minNonVar = p
			}
		}
	}
	elems, err := convertLargeWeights(elems)
	if err != nil {
		return err
	}
	cccs := []uint8{}
	nfd := norm.NFD.String(str)
	for i := range nfd {
		cccs = append(cccs, norm.NFD.PropertiesString(nfd[i:]).CCC())
	}
	if len(cccs) < len(elems) {
		if len(cccs) > 2 {
			return fmt.Errorf("number of decomposed characters should be greater or equal to the number of collation elements for len(colelems) > 3 (%d < %d)", len(cccs), len(elems))
		}
		p := len(elems) - 1
		for ; p > 0 && elems[p].w[0] == 0; p-- {
			elems[p].ccc = cccs[len(cccs)-1]
		}
		for ; p >= 0; p-- {
			elems[p].ccc = cccs[0]
		}
	} else {
		for i := range elems {
			elems[i].ccc = cccs[i]
		}
	}
	
	if len(elems) > 1 && len(cccs) > 1 && cccs[0] != 0 && cccs[0] != cccs[len(cccs)-1] {
		return fmt.Errorf("incompatible CCC values for expansion %X (%d)", runes, cccs)
	}
	b.root.newEntry(str, elems)
	return nil
}

func (t *Tailoring) setAnchor(anchor string) error {
	anchor = norm.NFC.String(anchor)
	a := t.index.find(anchor)
	if a == nil {
		a = t.index.newEntry(anchor, nil)
		a.implicit = true
		a.modified = true
		for _, r := range []rune(anchor) {
			e := t.index.find(string(r))
			e.lock = true
		}
	}
	t.anchor = a
	return nil
}







func (t *Tailoring) SetAnchor(anchor string) error {
	if err := t.setAnchor(anchor); err != nil {
		return err
	}
	t.before = false
	return nil
}



func (t *Tailoring) SetAnchorBefore(anchor string) error {
	if err := t.setAnchor(anchor); err != nil {
		return err
	}
	t.before = true
	return nil
}

































func (t *Tailoring) Insert(level colltab.Level, str, extend string) error {
	if t.anchor == nil {
		return fmt.Errorf("%s:Insert: no anchor point set for tailoring of %s", t.id, str)
	}
	str = norm.NFC.String(str)
	e := t.index.find(str)
	if e == nil {
		e = t.index.newEntry(str, nil)
	} else if e.logical != noAnchor {
		return fmt.Errorf("%s:Insert: cannot reinsert logical reset position %q", t.id, e.str)
	}
	if e.lock {
		return fmt.Errorf("%s:Insert: cannot reinsert element %q", t.id, e.str)
	}
	a := t.anchor
	
	
	
	e.before = t.before
	if t.before {
		t.before = false
		if a.prev == nil {
			a.insertBefore(e)
		} else {
			for a = a.prev; a.level > level; a = a.prev {
			}
			a.insertAfter(e)
		}
		e.level = level
	} else {
		for ; a.level > level; a = a.next {
		}
		e.level = a.level
		if a != e {
			a.insertAfter(e)
			a.level = level
		} else {
			
			
			
			a.prev.level = level
		}
	}
	e.extend = norm.NFD.String(extend)
	e.exclude = false
	e.modified = true
	e.elems = nil
	t.anchor = e
	return nil
}

func (o *ordering) getWeight(e *entry) []rawCE {
	if len(e.elems) == 0 && e.logical == noAnchor {
		if e.implicit {
			for _, r := range e.runes {
				e.elems = append(e.elems, o.getWeight(o.find(string(r)))...)
			}
		} else if e.before {
			count := [colltab.Identity + 1]int{}
			a := e
			for ; a.elems == nil && !a.implicit; a = a.next {
				count[a.level]++
			}
			e.elems = []rawCE{makeRawCE(a.elems[0].w, a.elems[0].ccc)}
			for i := colltab.Primary; i < colltab.Quaternary; i++ {
				if count[i] != 0 {
					e.elems[0].w[i] -= count[i]
					break
				}
			}
			if e.prev != nil {
				o.verifyWeights(e.prev, e, e.prev.level)
			}
		} else {
			prev := e.prev
			e.elems = nextWeight(prev.level, o.getWeight(prev))
			o.verifyWeights(e, e.next, e.level)
		}
	}
	return e.elems
}

func (o *ordering) addExtension(e *entry) {
	if ex := o.find(e.extend); ex != nil {
		e.elems = append(e.elems, ex.elems...)
	} else {
		for _, r := range []rune(e.extend) {
			e.elems = append(e.elems, o.find(string(r)).elems...)
		}
	}
	e.extend = ""
}

func (o *ordering) verifyWeights(a, b *entry, level colltab.Level) error {
	if level == colltab.Identity || b == nil || b.elems == nil || a.elems == nil {
		return nil
	}
	for i := colltab.Primary; i < level; i++ {
		if a.elems[0].w[i] < b.elems[0].w[i] {
			return nil
		}
	}
	if a.elems[0].w[level] >= b.elems[0].w[level] {
		err := fmt.Errorf("%s:overflow: collation elements of %q (%X) overflows those of %q (%X) at level %d (%X >= %X)", o.id, a.str, a.runes, b.str, b.runes, level, a.elems, b.elems)
		log.Println(err)
		
	}
	return nil
}

func (b *Builder) error(e error) {
	if e != nil {
		b.err = e
	}
}

func (b *Builder) errorID(locale string, e error) {
	if e != nil {
		b.err = fmt.Errorf("%s:%v", locale, e)
	}
}


func (o *ordering) patchNorm() {
	
	for _, e := range o.ordered {
		nfd := norm.NFD.String(e.str)
		if nfd != e.str {
			if e0 := o.find(nfd); e0 != nil && !e0.modified {
				e0.elems = e.elems
			} else if e.modified && !equalCEArrays(o.genColElems(nfd), e.elems) {
				e := o.newEntry(nfd, e.elems)
				e.modified = true
			}
		}
	}
	
	for _, e := range o.ordered {
		nfd := norm.NFD.String(e.str)
		if e.modified || nfd == e.str {
			continue
		}
		if e0 := o.find(nfd); e0 != nil {
			e.elems = e0.elems
		} else {
			e.elems = o.genColElems(nfd)
			if norm.NFD.LastBoundary([]byte(nfd)) == 0 {
				r := []rune(nfd)
				head := string(r[0])
				tail := ""
				for i := 1; i < len(r); i++ {
					s := norm.NFC.String(head + string(r[i]))
					if e0 := o.find(s); e0 != nil && e0.modified {
						head = s
					} else {
						tail += string(r[i])
					}
				}
				e.elems = append(o.genColElems(head), o.genColElems(tail)...)
			}
		}
	}
	
	for _, e := range o.ordered {
		if len(e.runes) > 1 && equalCEArrays(o.genColElems(e.str), e.elems) {
			e.exclude = true
		}
	}
}

func (b *Builder) buildOrdering(o *ordering) {
	for _, e := range o.ordered {
		o.getWeight(e)
	}
	for _, e := range o.ordered {
		o.addExtension(e)
	}
	o.patchNorm()
	o.sort()
	simplify(o)
	b.processExpansions(o)   
	b.processContractions(o) 

	t := newNode()
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if !e.skip() {
			ce, err := e.encode()
			b.errorID(o.id, err)
			t.insert(e.runes[0], ce)
		}
	}
	o.handle = b.index.addTrie(t)
}

func (b *Builder) build() (*table, error) {
	if b.built {
		return b.t, b.err
	}
	b.built = true
	b.t = &table{
		Table: colltab.Table{
			MaxContractLen: utf8.UTFMax,
			VariableTop:    uint32(b.varTop),
		},
	}

	b.buildOrdering(&b.root)
	b.t.root = b.root.handle
	for _, t := range b.locale {
		b.buildOrdering(t.index)
		if b.err != nil {
			break
		}
	}
	i, err := b.index.generate()
	b.t.trie = *i
	b.t.Index = colltab.Trie{
		Index:   i.index,
		Values:  i.values,
		Index0:  i.index[blockSize*b.t.root.lookupStart:],
		Values0: i.values[blockSize*b.t.root.valueStart:],
	}
	b.error(err)
	return b.t, b.err
}


func (b *Builder) Build() (colltab.Weighter, error) {
	table, err := b.build()
	if err != nil {
		return nil, err
	}
	return table, nil
}


func (t *Tailoring) Build() (colltab.Weighter, error) {
	
	return nil, nil
}



func (b *Builder) Print(w io.Writer) (n int, err error) {
	p := func(nn int, e error) {
		n += nn
		if err == nil {
			err = e
		}
	}
	t, err := b.build()
	if err != nil {
		return 0, err
	}
	p(fmt.Fprintf(w, `var availableLocales = "und`))
	for _, loc := range b.locale {
		if loc.id != "und" {
			p(fmt.Fprintf(w, ",%s", loc.id))
		}
	}
	p(fmt.Fprint(w, "\"\n\n"))
	p(fmt.Fprintf(w, "const varTop = 0x%x\n\n", b.varTop))
	p(fmt.Fprintln(w, "var locales = [...]tableIndex{"))
	for _, loc := range b.locale {
		if loc.id == "und" {
			p(t.fprintIndex(w, loc.index.handle, loc.id))
		}
	}
	for _, loc := range b.locale {
		if loc.id != "und" {
			p(t.fprintIndex(w, loc.index.handle, loc.id))
		}
	}
	p(fmt.Fprint(w, "}\n\n"))
	n, _, err = t.fprint(w, "main")
	return
}



func reproducibleFromNFKD(e *entry, exp, nfkd []rawCE) bool {
	
	if len(exp) != len(nfkd) {
		return false
	}
	for i, ce := range exp {
		
		if ce.w[0] != nfkd[i].w[0] || ce.w[1] != nfkd[i].w[1] {
			return false
		}
		
		
		
		if i >= 2 && ce.w[2] != maxTertiary {
			return false
		}
		if _, err := makeCE(ce); err != nil {
			
			return false
		}
	}
	return true
}

func simplify(o *ordering) {
	
	
	keep := make(map[rune]bool)
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if len(e.runes) > 1 {
			keep[e.runes[0]] = true
		}
	}
	
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		s := e.str
		nfkd := norm.NFKD.String(s)
		nfd := norm.NFD.String(s)
		if e.decompose || len(e.runes) > 1 || len(e.elems) == 1 || keep[e.runes[0]] || nfkd == nfd {
			continue
		}
		if reproducibleFromNFKD(e, e.elems, o.genColElems(nfkd)) {
			e.decompose = true
		}
	}
}




func (b *Builder) appendExpansion(e *entry) int {
	t := b.t
	i := len(t.ExpandElem)
	ce := uint32(len(e.elems))
	t.ExpandElem = append(t.ExpandElem, ce)
	for _, w := range e.elems {
		ce, err := makeCE(w)
		if err != nil {
			b.error(err)
			return -1
		}
		t.ExpandElem = append(t.ExpandElem, ce)
	}
	return i
}



func (b *Builder) processExpansions(o *ordering) {
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if !e.expansion() {
			continue
		}
		key := fmt.Sprintf("%v", e.elems)
		i, ok := b.expIndex[key]
		if !ok {
			i = b.appendExpansion(e)
			b.expIndex[key] = i
		}
		e.expansionIndex = i
	}
}

func (b *Builder) processContractions(o *ordering) {
	
	starters := []rune{}
	cm := make(map[rune][]*entry)
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if e.contraction() {
			if len(e.str) > b.t.MaxContractLen {
				b.t.MaxContractLen = len(e.str)
			}
			r := e.runes[0]
			if _, ok := cm[r]; !ok {
				starters = append(starters, r)
			}
			cm[r] = append(cm[r], e)
		}
	}
	
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if !e.contraction() {
			r := e.runes[0]
			if _, ok := cm[r]; ok {
				cm[r] = append(cm[r], e)
			}
		}
	}
	
	t := b.t
	for _, r := range starters {
		l := cm[r]
		
		
		
		sufx := []string{}
		hasSingle := false
		for _, e := range l {
			if len(e.runes) > 1 {
				sufx = append(sufx, string(e.runes[1:]))
			} else {
				hasSingle = true
			}
		}
		if !hasSingle {
			b.error(fmt.Errorf("no single entry for starter rune %U found", r))
			continue
		}
		
		sort.Strings(sufx)
		key := strings.Join(sufx, "\n")
		handle, ok := b.ctHandle[key]
		if !ok {
			var err error
			handle, err = appendTrie(&t.ContractTries, sufx)
			if err != nil {
				b.error(err)
			}
			b.ctHandle[key] = handle
		}
		
		es := make([]*entry, len(l))
		for _, e := range l {
			var p, sn int
			if len(e.runes) > 1 {
				str := []byte(string(e.runes[1:]))
				p, sn = lookup(&t.ContractTries, handle, str)
				if sn != len(str) {
					log.Fatalf("%s: processContractions: unexpected length for '%X'; len=%d; want %d", o.id, e.runes, sn, len(str))
				}
			}
			if es[p] != nil {
				log.Fatalf("%s: multiple contractions for position %d for rune %U", o.id, p, e.runes[0])
			}
			es[p] = e
		}
		
		elems := []uint32{}
		for _, e := range es {
			ce, err := e.encodeBase()
			b.errorID(o.id, err)
			elems = append(elems, ce)
		}
		key = fmt.Sprintf("%v", elems)
		i, ok := b.ctElem[key]
		if !ok {
			i = len(t.ContractElem)
			b.ctElem[key] = i
			t.ContractElem = append(t.ContractElem, elems...)
		}
		
		es[0].contractionIndex = i
		es[0].contractionHandle = handle
	}
}
