





package filter

import (
	"github.com/syndtr/goleveldb/leveldb/util"
)

func bloomHash(key []byte) uint32 {
	return util.Hash(key, 0xbc9f1d34)
}

type bloomFilter int




func (bloomFilter) Name() string {
	return "leveldb.BuiltinBloomFilter"
}

func (f bloomFilter) Contains(filter, key []byte) bool {
	nBytes := len(filter) - 1
	if nBytes < 1 {
		return false
	}
	nBits := uint32(nBytes * 8)

	
	
	k := filter[nBytes]
	if k > 30 {
		
		
		return true
	}

	kh := bloomHash(key)
	delta := (kh >> 17) | (kh << 15) 
	for j := uint8(0); j < k; j++ {
		bitpos := kh % nBits
		if (uint32(filter[bitpos/8]) & (1 << (bitpos % 8))) == 0 {
			return false
		}
		kh += delta
	}
	return true
}

func (f bloomFilter) NewGenerator() FilterGenerator {
	
	k := uint8(f * 69 / 100) 
	if k < 1 {
		k = 1
	} else if k > 30 {
		k = 30
	}
	return &bloomFilterGenerator{
		n: int(f),
		k: k,
	}
}

type bloomFilterGenerator struct {
	n int
	k uint8

	keyHashes []uint32
}

func (g *bloomFilterGenerator) Add(key []byte) {
	
	
	g.keyHashes = append(g.keyHashes, bloomHash(key))
}

func (g *bloomFilterGenerator) Generate(b Buffer) {
	
	nBits := uint32(len(g.keyHashes) * g.n)
	
	
	if nBits < 64 {
		nBits = 64
	}
	nBytes := (nBits + 7) / 8
	nBits = nBytes * 8

	dest := b.Alloc(int(nBytes) + 1)
	dest[nBytes] = g.k
	for _, kh := range g.keyHashes {
		delta := (kh >> 17) | (kh << 15) 
		for j := uint8(0); j < g.k; j++ {
			bitpos := kh % nBits
			dest[bitpos/8] |= (1 << (bitpos % 8))
			kh += delta
		}
	}

	g.keyHashes = g.keyHashes[:0]
}









func NewBloomFilter(bitsPerKey int) Filter {
	return bloomFilter(bitsPerKey)
}
