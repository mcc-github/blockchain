
























































package triegen 







import (
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"io"
	"log"
	"unicode/utf8"
)



type builder struct {
	Name string

	
	ValueType string

	
	ValueSize int

	
	
	IndexType string

	
	IndexSize int

	
	
	
	SourceType string

	Trie []*Trie

	IndexBlocks []*node
	ValueBlocks [][]uint64
	Compactions []compaction
	Checksum    uint64

	ASCIIBlock   string
	StarterBlock string

	indexBlockIdx map[uint64]int
	valueBlockIdx map[uint64]nodeIndex
	asciiBlockIdx map[uint64]int

	
	Stats struct {
		NValueEntries int
		NValueBytes   int
		NIndexEntries int
		NIndexBytes   int
		NHandleBytes  int
	}

	err error
}




type nodeIndex struct {
	compaction int
	index      int
}


type compaction struct {
	c         Compacter
	blocks    []*node
	maxHandle uint32
	totalSize int

	
	Cutoff  uint32
	Offset  uint32
	Handler string
}

func (b *builder) setError(err error) {
	if b.err == nil {
		b.err = err
	}
}


type Option func(b *builder) error


func Compact(c Compacter) Option {
	return func(b *builder) error {
		b.Compactions = append(b.Compactions, compaction{
			c:       c,
			Handler: c.Handler() + "(n, b)"})
		return nil
	}
}






func Gen(w io.Writer, name string, tries []*Trie, opts ...Option) (sz int, err error) {
	
	
	
	b := &builder{
		Name:        name,
		Trie:        tries,
		IndexBlocks: []*node{{}, {}, {}},
		Compactions: []compaction{{
			Handler: name + "Values[n<<6+uint32(b)]",
		}},
		
		
		indexBlockIdx: map[uint64]int{0: 0},
		valueBlockIdx: map[uint64]nodeIndex{0: {}},
		asciiBlockIdx: map[uint64]int{},
	}
	b.Compactions[0].c = (*simpleCompacter)(b)

	for _, f := range opts {
		if err := f(b); err != nil {
			return 0, err
		}
	}
	b.build()
	if b.err != nil {
		return 0, b.err
	}
	if err = b.print(w); err != nil {
		return 0, err
	}
	return b.Size(), nil
}



type Trie struct {
	root *node

	hiddenTrie
}



type hiddenTrie struct {
	Name         string
	Checksum     uint64
	ASCIIIndex   int
	StarterIndex int
}


func NewTrie(name string) *Trie {
	return &Trie{
		&node{
			children: make([]*node, blockSize),
			values:   make([]uint64, utf8.RuneSelf),
		},
		hiddenTrie{Name: name},
	}
}




func (t *Trie) Gen(w io.Writer, opts ...Option) (sz int, err error) {
	return Gen(w, t.Name, []*Trie{t}, opts...)
}


type node struct {
	
	
	children []*node

	
	
	
	
	values []uint64

	index nodeIndex
}



func (t *Trie) Insert(r rune, value uint64) {
	if value == 0 {
		return
	}
	s := string(r)
	if []rune(s)[0] != r && value != 0 {
		
		
		
		
		panic(fmt.Sprintf("triegen: non-zero value for invalid rune %U", r))
	}
	if len(s) == 1 {
		
		t.root.values[s[0]] = value
		return
	}

	n := t.root
	for ; len(s) > 1; s = s[1:] {
		if n.children == nil {
			n.children = make([]*node, blockSize)
		}
		p := s[0] % blockSize
		c := n.children[p]
		if c == nil {
			c = &node{}
			n.children[p] = c
		}
		if len(s) > 2 && c.values != nil {
			log.Fatalf("triegen: insert(%U): found internal node with values", r)
		}
		n = c
	}
	if n.values == nil {
		n.values = make([]uint64, blockSize)
	}
	if n.children != nil {
		log.Fatalf("triegen: insert(%U): found leaf node that also has child nodes", r)
	}
	n.values[s[0]-0x80] = value
}



func (b *builder) Size() int {
	
	sz := len(b.IndexBlocks) * blockSize * b.IndexSize

	
	
	
	sz += len(b.ValueBlocks) * blockSize * b.ValueSize
	for _, c := range b.Compactions[1:] {
		sz += c.totalSize
	}

	
	
	
	
	

	
	if len(b.Trie) > 1 {
		sz += 2 * b.IndexSize * len(b.Trie)
	}
	return sz
}

func (b *builder) build() {
	
	var vmax uint64
	for _, t := range b.Trie {
		vmax = maxValue(t.root, vmax)
	}
	b.ValueType, b.ValueSize = getIntType(vmax)

	
	
	
	
	
	for _, t := range b.Trie {
		b.Checksum += b.buildTrie(t)
	}

	
	offset := uint32(0)
	for i := range b.Compactions {
		c := &b.Compactions[i]
		c.Offset = offset
		offset += c.maxHandle + 1
		c.Cutoff = offset
	}

	
	
	
	imax := uint64(b.Compactions[len(b.Compactions)-1].Cutoff)
	for _, ib := range b.IndexBlocks {
		if x := uint64(ib.index.index); x > imax {
			imax = x
		}
	}
	b.IndexType, b.IndexSize = getIntType(imax)
}

func maxValue(n *node, max uint64) uint64 {
	if n == nil {
		return max
	}
	for _, c := range n.children {
		max = maxValue(c, max)
	}
	for _, v := range n.values {
		if max < v {
			max = v
		}
	}
	return max
}

func getIntType(v uint64) (string, int) {
	switch {
	case v < 1<<8:
		return "uint8", 1
	case v < 1<<16:
		return "uint16", 2
	case v < 1<<32:
		return "uint32", 4
	}
	return "uint64", 8
}

const (
	blockSize = 64

	
	blockOffset = 2

	
	rootBlockOffset = 3
)

var crcTable = crc64.MakeTable(crc64.ISO)

func (b *builder) buildTrie(t *Trie) uint64 {
	n := t.root

	
	
	hasher := crc64.New(crcTable)
	binary.Write(hasher, binary.BigEndian, n.values)
	hash := hasher.Sum64()

	v, ok := b.asciiBlockIdx[hash]
	if !ok {
		v = len(b.ValueBlocks)
		b.asciiBlockIdx[hash] = v

		b.ValueBlocks = append(b.ValueBlocks, n.values[:blockSize], n.values[blockSize:])
		if v == 0 {
			
			
			
			
			
			b.ValueBlocks = append(b.ValueBlocks, make([]uint64, blockSize))
		}
	}
	t.ASCIIIndex = v

	
	t.Checksum = b.computeOffsets(n, true)
	
	
	t.StarterIndex = n.index.index - (rootBlockOffset - blockOffset)
	return t.Checksum
}

func (b *builder) computeOffsets(n *node, root bool) uint64 {
	
	
	first := len(b.IndexBlocks) == rootBlockOffset
	if first {
		b.IndexBlocks = append(b.IndexBlocks, n)
	}

	
	
	hash := uint64(0)
	if n.children != nil || n.values != nil {
		hasher := crc64.New(crcTable)
		for _, c := range n.children {
			var v uint64
			if c != nil {
				v = b.computeOffsets(c, false)
			}
			binary.Write(hasher, binary.BigEndian, v)
		}
		binary.Write(hasher, binary.BigEndian, n.values)
		hash = hasher.Sum64()
	}

	if first {
		b.indexBlockIdx[hash] = rootBlockOffset - blockOffset
	}

	
	if n.children != nil {
		v, ok := b.indexBlockIdx[hash]
		if !ok {
			v = len(b.IndexBlocks) - blockOffset
			b.IndexBlocks = append(b.IndexBlocks, n)
			b.indexBlockIdx[hash] = v
		}
		n.index = nodeIndex{0, v}
	} else {
		h, ok := b.valueBlockIdx[hash]
		if !ok {
			bestI, bestSize := 0, blockSize*b.ValueSize
			for i, c := range b.Compactions[1:] {
				if sz, ok := c.c.Size(n.values); ok && bestSize > sz {
					bestI, bestSize = i+1, sz
				}
			}
			c := &b.Compactions[bestI]
			c.totalSize += bestSize
			v := c.c.Store(n.values)
			if c.maxHandle < v {
				c.maxHandle = v
			}
			h = nodeIndex{bestI, int(v)}
			b.valueBlockIdx[hash] = h
		}
		n.index = h
	}
	return hash
}
