












package unix


func Major(dev uint64) uint32 {
	return uint32((dev >> 8) & 0xff)
}


func Minor(dev uint64) uint32 {
	return uint32(dev & 0xffff00ff)
}



func Mkdev(major, minor uint32) uint64 {
	return (uint64(major) << 8) | uint64(minor)
}
