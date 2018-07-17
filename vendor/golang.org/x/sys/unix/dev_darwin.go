






package unix


func Major(dev uint64) uint32 {
	return uint32((dev >> 24) & 0xff)
}


func Minor(dev uint64) uint32 {
	return uint32(dev & 0xffffff)
}



func Mkdev(major, minor uint32) uint64 {
	return (uint64(major) << 24) | uint64(minor)
}
