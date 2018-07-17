










package build

import (
	"fmt"
	"hash/fnv"
	"io"
	"reflect"
)

const (
	blockSize   = 64
	blockOffset = 2 
)

type trieHandle struct {
	lookupStart uint16 
	valueStart  uint16 
}

type trie struct {
	index  []uint16
	values []uint32
}


type trieNode struct {
	index    []*trieNode
	value    []uint32
	b        byte
	refValue uint16
	refIndex uint16
}

func newNode() *trieNode {
	return &trieNode{
		index: make([]*trieNode, 64),
		value: make([]uint32, 128), 
	}
}

func (n *trieNode) isInternal() bool {
	return n.value != nil
}

func (n *trieNode) insert(r rune, value uint32) {
	const maskx = 0x3F 
	str := string(r)
	if len(str) == 1 {
		n.value[str[0]] = value
		return
	}
	for i := 0; i < len(str)-1; i++ {
		b := str[i] & maskx
		if n.index == nil {
			n.index = make([]*trieNode, blockSize)
		}
		nn := n.index[b]
		if nn == nil {
			nn = &trieNode{}
			nn.b = b
			n.index[b] = nn
		}
		n = nn
	}
	if n.value == nil {
		n.value = make([]uint32, blockSize)
	}
	b := str[len(str)-1] & maskx
	n.value[b] = value
}

type trieBuilder struct {
	t *trie

	roots []*trieHandle

	lookupBlocks []*trieNode
	valueBlocks  []*trieNode

	lookupBlockIdx map[uint32]*trieNode
	valueBlockIdx  map[uint32]*trieNode
}

func newTrieBuilder() *trieBuilder {
	index := &trieBuilder{}
	index.lookupBlocks = make([]*trieNode, 0)
	index.valueBlocks = make([]*trieNode, 0)
	index.lookupBlockIdx = make(map[uint32]*trieNode)
	index.valueBlockIdx = make(map[uint32]*trieNode)
	
	
	index.lookupBlocks = append(index.lookupBlocks, nil, nil, nil)
	index.t = &trie{}
	return index
}

func (b *trieBuilder) computeOffsets(n *trieNode) *trieNode {
	hasher := fnv.New32()
	if n.index != nil {
		for i, nn := range n.index {
			var vi, vv uint16
			if nn != nil {
				nn = b.computeOffsets(nn)
				n.index[i] = nn
				vi = nn.refIndex
				vv = nn.refValue
			}
			hasher.Write([]byte{byte(vi >> 8), byte(vi)})
			hasher.Write([]byte{byte(vv >> 8), byte(vv)})
		}
		h := hasher.Sum32()
		nn, ok := b.lookupBlockIdx[h]
		if !ok {
			n.refIndex = uint16(len(b.lookupBlocks)) - blockOffset
			b.lookupBlocks = append(b.lookupBlocks, n)
			b.lookupBlockIdx[h] = n
		} else {
			n = nn
		}
	} else {
		for _, v := range n.value {
			hasher.Write([]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
		}
		h := hasher.Sum32()
		nn, ok := b.valueBlockIdx[h]
		if !ok {
			n.refValue = uint16(len(b.valueBlocks)) - blockOffset
			n.refIndex = n.refValue
			b.valueBlocks = append(b.valueBlocks, n)
			b.valueBlockIdx[h] = n
		} else {
			n = nn
		}
	}
	return n
}

func (b *trieBuilder) addStartValueBlock(n *trieNode) uint16 {
	hasher := fnv.New32()
	for _, v := range n.value[:2*blockSize] {
		hasher.Write([]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
	}
	h := hasher.Sum32()
	nn, ok := b.valueBlockIdx[h]
	if !ok {
		n.refValue = uint16(len(b.valueBlocks))
		n.refIndex = n.refValue
		b.valueBlocks = append(b.valueBlocks, n)
		
		b.valueBlocks = append(b.valueBlocks, nil)
		b.valueBlockIdx[h] = n
	} else {
		n = nn
	}
	return n.refValue
}

func genValueBlock(t *trie, n *trieNode) {
	if n != nil {
		for _, v := range n.value {
			t.values = append(t.values, v)
		}
	}
}

func genLookupBlock(t *trie, n *trieNode) {
	for _, nn := range n.index {
		v := uint16(0)
		if nn != nil {
			if n.index != nil {
				v = nn.refIndex
			} else {
				v = nn.refValue
			}
		}
		t.index = append(t.index, v)
	}
}

func (b *trieBuilder) addTrie(n *trieNode) *trieHandle {
	h := &trieHandle{}
	b.roots = append(b.roots, h)
	h.valueStart = b.addStartValueBlock(n)
	if len(b.roots) == 1 {
		
		
		
		
		b.valueBlocks = append(b.valueBlocks, nil)
	}
	n = b.computeOffsets(n)
	
	h.lookupStart = n.refIndex - 1
	return h
}


func (b *trieBuilder) generate() (t *trie, err error) {
	t = b.t
	if len(b.valueBlocks) >= 1<<16 {
		return nil, fmt.Errorf("maximum number of value blocks exceeded (%d > %d)", len(b.valueBlocks), 1<<16)
	}
	if len(b.lookupBlocks) >= 1<<16 {
		return nil, fmt.Errorf("maximum number of lookup blocks exceeded (%d > %d)", len(b.lookupBlocks), 1<<16)
	}
	genValueBlock(t, b.valueBlocks[0])
	genValueBlock(t, &trieNode{value: make([]uint32, 64)})
	for i := 2; i < len(b.valueBlocks); i++ {
		genValueBlock(t, b.valueBlocks[i])
	}
	n := &trieNode{index: make([]*trieNode, 64)}
	genLookupBlock(t, n)
	genLookupBlock(t, n)
	genLookupBlock(t, n)
	for i := 3; i < len(b.lookupBlocks); i++ {
		genLookupBlock(t, b.lookupBlocks[i])
	}
	return b.t, nil
}

func (t *trie) printArrays(w io.Writer, name string) (n, size int, err error) {
	p := func(f string, a ...interface{}) {
		nn, e := fmt.Fprintf(w, f, a...)
		n += nn
		if err == nil {
			err = e
		}
	}
	nv := len(t.values)
	p("
	p("
	p("var %sValues = [%d]uint32 {", name, nv)
	var printnewline bool
	for i, v := range t.values {
		if i%blockSize == 0 {
			p("\n\t
		}
		if i%4 == 0 {
			printnewline = true
		}
		if v != 0 {
			if printnewline {
				p("\n\t")
				printnewline = false
			}
			p("%#04x:%#08x, ", i, v)
		}
	}
	p("\n}\n\n")
	ni := len(t.index)
	p("
	p("
	p("var %sLookup = [%d]uint16 {", name, ni)
	printnewline = false
	for i, v := range t.index {
		if i%blockSize == 0 {
			p("\n\t
		}
		if i%8 == 0 {
			printnewline = true
		}
		if v != 0 {
			if printnewline {
				p("\n\t")
				printnewline = false
			}
			p("%#03x:%#02x, ", i, v)
		}
	}
	p("\n}\n\n")
	return n, nv*4 + ni*2, err
}

func (t *trie) printStruct(w io.Writer, handle *trieHandle, name string) (n, sz int, err error) {
	const msg = "trie{ %sLookup[%d:], %sValues[%d:], %sLookup[:], %sValues[:]}"
	n, err = fmt.Fprintf(w, msg, name, handle.lookupStart*blockSize, name, handle.valueStart*blockSize, name, name)
	sz += int(reflect.TypeOf(trie{}).Size())
	return
}
