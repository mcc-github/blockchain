









package lz4

const (
	
	Extension = ".lz4"
	
	Version = 1

	frameMagic     uint32 = 0x184D2204
	frameSkipMagic uint32 = 0x184D2A50

	
	minMatch            = 4  
	winSizeLog          = 16 
	winSize             = 1 << winSizeLog
	winMask             = winSize - 1 
	compressedBlockFlag = 1 << 31
	compressedBlockMask = compressedBlockFlag - 1

	
	
	
	
	hashLog       = 16
	hashTableSize = 1 << hashLog
	hashShift     = uint((minMatch * 8) - hashLog)

	mfLimit      = 8 + minMatch 
	skipStrength = 6            
)


var (
	bsMapID    = map[byte]int{4: 64 << 10, 5: 256 << 10, 6: 1 << 20, 7: 4 << 20}
	bsMapValue = make(map[int]byte, len(bsMapID))
)


func init() {
	for i, v := range bsMapID {
		bsMapValue[v] = i
	}
}







type Header struct {
	BlockChecksum    bool   
	NoChecksum       bool   
	BlockMaxSize     int    
	Size             uint64 
	CompressionLevel int    
	done             bool   
}
