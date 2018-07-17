





package snappy

func load32(b []byte, i int) uint32 {
	b = b[i : i+4 : len(b)] 
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func load64(b []byte, i int) uint64 {
	b = b[i : i+8 : len(b)] 
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}






func emitLiteral(dst, lit []byte) int {
	i, n := 0, uint(len(lit)-1)
	switch {
	case n < 60:
		dst[0] = uint8(n)<<2 | tagLiteral
		i = 1
	case n < 1<<8:
		dst[0] = 60<<2 | tagLiteral
		dst[1] = uint8(n)
		i = 2
	default:
		dst[0] = 61<<2 | tagLiteral
		dst[1] = uint8(n)
		dst[2] = uint8(n >> 8)
		i = 3
	}
	return i + copy(dst[i:], lit)
}







func emitCopy(dst []byte, offset, length int) int {
	i := 0
	
	
	
	
	
	
	
	
	
	for length >= 68 {
		
		dst[i+0] = 63<<2 | tagCopy2
		dst[i+1] = uint8(offset)
		dst[i+2] = uint8(offset >> 8)
		i += 3
		length -= 64
	}
	if length > 64 {
		
		dst[i+0] = 59<<2 | tagCopy2
		dst[i+1] = uint8(offset)
		dst[i+2] = uint8(offset >> 8)
		i += 3
		length -= 60
	}
	if length >= 12 || offset >= 2048 {
		
		dst[i+0] = uint8(length-1)<<2 | tagCopy2
		dst[i+1] = uint8(offset)
		dst[i+2] = uint8(offset >> 8)
		return i + 3
	}
	
	dst[i+0] = uint8(offset>>8)<<5 | uint8(length-4)<<2 | tagCopy1
	dst[i+1] = uint8(offset)
	return i + 2
}






func extendMatch(src []byte, i, j int) int {
	for ; j < len(src) && src[i] == src[j]; i, j = i+1, j+1 {
	}
	return j
}

func hash(u, shift uint32) uint32 {
	return (u * 0x1e35a7bd) >> shift
}








func encodeBlock(dst, src []byte) (d int) {
	
	
	
	const (
		maxTableSize = 1 << 14
		
		
		tableMask = maxTableSize - 1
	)
	shift := uint32(32 - 8)
	for tableSize := 1 << 8; tableSize < maxTableSize && tableSize < len(src); tableSize *= 2 {
		shift--
	}
	
	
	
	
	var table [maxTableSize]uint16

	
	
	
	sLimit := len(src) - inputMargin

	
	nextEmit := 0

	
	
	s := 1
	nextHash := hash(load32(src, s), shift)

	for {
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		skip := 32

		nextS := s
		candidate := 0
		for {
			s = nextS
			bytesBetweenHashLookups := skip >> 5
			nextS = s + bytesBetweenHashLookups
			skip += bytesBetweenHashLookups
			if nextS > sLimit {
				goto emitRemainder
			}
			candidate = int(table[nextHash&tableMask])
			table[nextHash&tableMask] = uint16(s)
			nextHash = hash(load32(src, nextS), shift)
			if load32(src, s) == load32(src, candidate) {
				break
			}
		}

		
		
		
		d += emitLiteral(dst[d:], src[nextEmit:s])

		
		
		
		
		
		
		
		
		for {
			
			
			base := s

			
			
			
			
			s += 4
			for i := candidate + 4; s < len(src) && src[i] == src[s]; i, s = i+1, s+1 {
			}

			d += emitCopy(dst[d:], base-candidate, s-base)
			nextEmit = s
			if s >= sLimit {
				goto emitRemainder
			}

			
			
			
			
			
			
			x := load64(src, s-1)
			prevHash := hash(uint32(x>>0), shift)
			table[prevHash&tableMask] = uint16(s - 1)
			currHash := hash(uint32(x>>8), shift)
			candidate = int(table[currHash&tableMask])
			table[currHash&tableMask] = uint16(s)
			if uint32(x>>8) != load32(src, candidate) {
				nextHash = hash(uint32(x>>16), shift)
				s++
				break
			}
		}
	}

emitRemainder:
	if nextEmit < len(src) {
		d += emitLiteral(dst[d:], src[nextEmit:])
	}
	return d
}
