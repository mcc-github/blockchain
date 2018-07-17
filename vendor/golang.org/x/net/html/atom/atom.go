
















package atom 


type Atom uint32


func (a Atom) String() string {
	start := uint32(a >> 8)
	n := uint32(a & 0xff)
	if start+n > uint32(len(atomText)) {
		return ""
	}
	return atomText[start : start+n]
}

func (a Atom) string() string {
	return atomText[a>>8 : a>>8+a&0xff]
}


func fnv(h uint32, s []byte) uint32 {
	for i := range s {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

func match(s string, t []byte) bool {
	for i, c := range t {
		if s[i] != c {
			return false
		}
	}
	return true
}



func Lookup(s []byte) Atom {
	if len(s) == 0 || len(s) > maxAtomLen {
		return 0
	}
	h := fnv(hash0, s)
	if a := table[h&uint32(len(table)-1)]; int(a&0xff) == len(s) && match(a.string(), s) {
		return a
	}
	if a := table[(h>>16)&uint32(len(table)-1)]; int(a&0xff) == len(s) && match(a.string(), s) {
		return a
	}
	return 0
}



func String(s []byte) string {
	if a := Lookup(s); a != 0 {
		return a.String()
	}
	return string(s)
}
