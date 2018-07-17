






package unix


func Major(dev uint64) uint32 {
	return uint32((dev & 0x000fff00) >> 8)
}


func Minor(dev uint64) uint32 {
	minor := uint32((dev & 0x000000ff) >> 0)
	minor |= uint32((dev & 0xfff00000) >> 12)
	return minor
}



func Mkdev(major, minor uint32) uint64 {
	dev := (uint64(major) << 8) & 0x000fff00
	dev |= (uint64(minor) << 12) & 0xfff00000
	dev |= (uint64(minor) << 0) & 0x000000ff
	return dev
}
