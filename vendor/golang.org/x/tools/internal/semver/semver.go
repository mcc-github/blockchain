





















package semver


type parsed struct {
	major      string
	minor      string
	patch      string
	short      string
	prerelease string
	build      string
	err        string
}


func IsValid(v string) bool {
	_, ok := parse(v)
	return ok
}






func Canonical(v string) string {
	p, ok := parse(v)
	if !ok {
		return ""
	}
	if p.build != "" {
		return v[:len(v)-len(p.build)]
	}
	if p.short != "" {
		return v + p.short
	}
	return v
}




func Major(v string) string {
	pv, ok := parse(v)
	if !ok {
		return ""
	}
	return v[:1+len(pv.major)]
}




func MajorMinor(v string) string {
	pv, ok := parse(v)
	if !ok {
		return ""
	}
	i := 1 + len(pv.major)
	if j := i + 1 + len(pv.minor); j <= len(v) && v[i] == '.' && v[i+1:j] == pv.minor {
		return v[:j]
	}
	return v[:i] + "." + pv.minor
}




func Prerelease(v string) string {
	pv, ok := parse(v)
	if !ok {
		return ""
	}
	return pv.prerelease
}




func Build(v string) string {
	pv, ok := parse(v)
	if !ok {
		return ""
	}
	return pv.build
}







func Compare(v, w string) int {
	pv, ok1 := parse(v)
	pw, ok2 := parse(w)
	if !ok1 && !ok2 {
		return 0
	}
	if !ok1 {
		return -1
	}
	if !ok2 {
		return +1
	}
	if c := compareInt(pv.major, pw.major); c != 0 {
		return c
	}
	if c := compareInt(pv.minor, pw.minor); c != 0 {
		return c
	}
	if c := compareInt(pv.patch, pw.patch); c != 0 {
		return c
	}
	return comparePrerelease(pv.prerelease, pw.prerelease)
}



func Max(v, w string) string {
	v = Canonical(v)
	w = Canonical(w)
	if Compare(v, w) > 0 {
		return v
	}
	return w
}

func parse(v string) (p parsed, ok bool) {
	if v == "" || v[0] != 'v' {
		p.err = "missing v prefix"
		return
	}
	p.major, v, ok = parseInt(v[1:])
	if !ok {
		p.err = "bad major version"
		return
	}
	if v == "" {
		p.minor = "0"
		p.patch = "0"
		p.short = ".0.0"
		return
	}
	if v[0] != '.' {
		p.err = "bad minor prefix"
		ok = false
		return
	}
	p.minor, v, ok = parseInt(v[1:])
	if !ok {
		p.err = "bad minor version"
		return
	}
	if v == "" {
		p.patch = "0"
		p.short = ".0"
		return
	}
	if v[0] != '.' {
		p.err = "bad patch prefix"
		ok = false
		return
	}
	p.patch, v, ok = parseInt(v[1:])
	if !ok {
		p.err = "bad patch version"
		return
	}
	if len(v) > 0 && v[0] == '-' {
		p.prerelease, v, ok = parsePrerelease(v)
		if !ok {
			p.err = "bad prerelease"
			return
		}
	}
	if len(v) > 0 && v[0] == '+' {
		p.build, v, ok = parseBuild(v)
		if !ok {
			p.err = "bad build"
			return
		}
	}
	if v != "" {
		p.err = "junk on end"
		ok = false
		return
	}
	ok = true
	return
}

func parseInt(v string) (t, rest string, ok bool) {
	if v == "" {
		return
	}
	if v[0] < '0' || '9' < v[0] {
		return
	}
	i := 1
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	if v[0] == '0' && i != 1 {
		return
	}
	return v[:i], v[i:], true
}

func parsePrerelease(v string) (t, rest string, ok bool) {
	
	
	
	
	if v == "" || v[0] != '-' {
		return
	}
	i := 1
	start := 1
	for i < len(v) && v[i] != '+' {
		if !isIdentChar(v[i]) && v[i] != '.' {
			return
		}
		if v[i] == '.' {
			if start == i || isBadNum(v[start:i]) {
				return
			}
			start = i + 1
		}
		i++
	}
	if start == i || isBadNum(v[start:i]) {
		return
	}
	return v[:i], v[i:], true
}

func parseBuild(v string) (t, rest string, ok bool) {
	if v == "" || v[0] != '+' {
		return
	}
	i := 1
	start := 1
	for i < len(v) {
		if !isIdentChar(v[i]) {
			return
		}
		if v[i] == '.' {
			if start == i {
				return
			}
			start = i + 1
		}
		i++
	}
	if start == i {
		return
	}
	return v[:i], v[i:], true
}

func isIdentChar(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '-'
}

func isBadNum(v string) bool {
	i := 0
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	return i == len(v) && i > 1 && v[0] == '0'
}

func isNum(v string) bool {
	i := 0
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	return i == len(v)
}

func compareInt(x, y string) int {
	if x == y {
		return 0
	}
	if len(x) < len(y) {
		return -1
	}
	if len(x) > len(y) {
		return +1
	}
	if x < y {
		return -1
	} else {
		return +1
	}
}

func comparePrerelease(x, y string) int {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	if x == y {
		return 0
	}
	if x == "" {
		return +1
	}
	if y == "" {
		return -1
	}
	for x != "" && y != "" {
		x = x[1:] 
		y = y[1:] 
		var dx, dy string
		dx, x = nextIdent(x)
		dy, y = nextIdent(y)
		if dx != dy {
			ix := isNum(dx)
			iy := isNum(dy)
			if ix != iy {
				if ix {
					return -1
				} else {
					return +1
				}
			}
			if ix {
				if len(dx) < len(dy) {
					return -1
				}
				if len(dx) > len(dy) {
					return +1
				}
			}
			if dx < dy {
				return -1
			} else {
				return +1
			}
		}
	}
	if x == "" {
		return -1
	} else {
		return +1
	}
}

func nextIdent(x string) (dx, rest string) {
	i := 0
	for i < len(x) && x[i] != '.' {
		i++
	}
	return x[:i], x[i:]
}
