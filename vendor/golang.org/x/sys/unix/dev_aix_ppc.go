









package unix


func Major(dev uint64) uint32 {
	return uint32((dev >> 16) & 0xffff)
}


func Minor(dev uint64) uint32 {
	return uint32(dev & 0xffff)
}



func Mkdev(major, minor uint32) uint64 {
	return uint64(((major) << 16) | (minor))
}
