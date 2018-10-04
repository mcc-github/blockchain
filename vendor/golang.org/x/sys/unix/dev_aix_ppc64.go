









package unix


func Major(dev uint64) uint32 {
	return uint32((dev & 0x3fffffff00000000) >> 32)
}


func Minor(dev uint64) uint32 {
	return uint32((dev & 0x00000000ffffffff) >> 0)
}



func Mkdev(major, minor uint32) uint64 {
	var DEVNO64 uint64
	DEVNO64 = 0x8000000000000000
	return ((uint64(major) << 32) | (uint64(minor) & 0x00000000FFFFFFFF) | DEVNO64)
}
