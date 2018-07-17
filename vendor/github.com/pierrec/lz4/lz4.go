
















package lz4

import (
	"hash"
	"sync"

	"github.com/pierrec/xxHash/xxHash32"
)

const (
	
	Extension = ".lz4"
	
	Version = 1

	frameMagic     = uint32(0x184D2204)
	frameSkipMagic = uint32(0x184D2A50)

	
	minMatch   = 4  
	winSizeLog = 16 
	winSize    = 1 << winSizeLog
	winMask    = winSize - 1 

	
	
	
	
	hashLog       = 16
	hashTableSize = 1 << hashLog
	hashShift     = uint((minMatch * 8) - hashLog)

	mfLimit      = 8 + minMatch 
	skipStrength = 6            

	hasher = uint32(2654435761) 
)


var bsMapID = map[byte]int{4: 64 << 10, 5: 256 << 10, 6: 1 << 20, 7: 4 << 20}
var bsMapValue = map[int]byte{}


func init() {
	for i, v := range bsMapID {
		bsMapValue[v] = i
	}
}






type Header struct {
	BlockDependency bool   
	BlockChecksum   bool   
	NoChecksum      bool   
	BlockMaxSize    int    
	Size            uint64 
	HighCompression bool   
	done            bool   
	
	
	
}



type xxhPool struct {
	sync.Pool
}

func (p *xxhPool) Get() hash.Hash32 {
	return p.Pool.Get().(hash.Hash32)
}

func (p *xxhPool) Put(h hash.Hash32) {
	h.Reset()
	p.Pool.Put(h)
}


var hashPool = xxhPool{
	Pool: sync.Pool{
		New: func() interface{} { return xxHash32.New(0) },
	},
}
