package prometheus



const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)


func hashNew() uint64 {
	return offset64
}


func hashAdd(h uint64, s string) uint64 {
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}


func hashAddByte(h uint64, b byte) uint64 {
	h ^= uint64(b)
	h *= prime64
	return h
}
