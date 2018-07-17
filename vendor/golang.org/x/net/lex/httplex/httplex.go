








package httplex

import (
	"net"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

var isTokenTable = [127]bool{
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'0':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	'A':  true,
	'B':  true,
	'C':  true,
	'D':  true,
	'E':  true,
	'F':  true,
	'G':  true,
	'H':  true,
	'I':  true,
	'J':  true,
	'K':  true,
	'L':  true,
	'M':  true,
	'N':  true,
	'O':  true,
	'P':  true,
	'Q':  true,
	'R':  true,
	'S':  true,
	'T':  true,
	'U':  true,
	'W':  true,
	'V':  true,
	'X':  true,
	'Y':  true,
	'Z':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'a':  true,
	'b':  true,
	'c':  true,
	'd':  true,
	'e':  true,
	'f':  true,
	'g':  true,
	'h':  true,
	'i':  true,
	'j':  true,
	'k':  true,
	'l':  true,
	'm':  true,
	'n':  true,
	'o':  true,
	'p':  true,
	'q':  true,
	'r':  true,
	's':  true,
	't':  true,
	'u':  true,
	'v':  true,
	'w':  true,
	'x':  true,
	'y':  true,
	'z':  true,
	'|':  true,
	'~':  true,
}

func IsTokenRune(r rune) bool {
	i := int(r)
	return i < len(isTokenTable) && isTokenTable[i]
}

func isNotToken(r rune) bool {
	return !IsTokenRune(r)
}



func HeaderValuesContainsToken(values []string, token string) bool {
	for _, v := range values {
		if headerValueContainsToken(v, token) {
			return true
		}
	}
	return false
}



func isOWS(b byte) bool { return b == ' ' || b == '\t' }



func trimOWS(x string) string {
	
	
	
	
	for len(x) > 0 && isOWS(x[0]) {
		x = x[1:]
	}
	for len(x) > 0 && isOWS(x[len(x)-1]) {
		x = x[:len(x)-1]
	}
	return x
}





func headerValueContainsToken(v string, token string) bool {
	v = trimOWS(v)
	if comma := strings.IndexByte(v, ','); comma != -1 {
		return tokenEqual(trimOWS(v[:comma]), token) || headerValueContainsToken(v[comma+1:], token)
	}
	return tokenEqual(v, token)
}


func lowerASCII(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}


func tokenEqual(t1, t2 string) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i, b := range t1 {
		if b >= utf8.RuneSelf {
			
			return false
		}
		if lowerASCII(byte(b)) != lowerASCII(t2[i]) {
			return false
		}
	}
	return true
}




func isLWS(b byte) bool { return b == ' ' || b == '\t' }





func isCTL(b byte) bool {
	const del = 0x7f 
	return b < ' ' || b == del
}











func ValidHeaderFieldName(v string) bool {
	if len(v) == 0 {
		return false
	}
	for _, r := range v {
		if !IsTokenRune(r) {
			return false
		}
	}
	return true
}


func ValidHostHeader(h string) bool {
	
	
	
	
	
	
	
	
	
	
	
	for i := 0; i < len(h); i++ {
		if !validHostByte[h[i]] {
			return false
		}
	}
	return true
}


var validHostByte = [256]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true,
	'8': true, '9': true,

	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true,
	'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true,
	'q': true, 'r': true, 's': true, 't': true, 'u': true, 'v': true, 'w': true, 'x': true,
	'y': true, 'z': true,

	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true, 'H': true,
	'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true, 'O': true, 'P': true,
	'Q': true, 'R': true, 'S': true, 'T': true, 'U': true, 'V': true, 'W': true, 'X': true,
	'Y': true, 'Z': true,

	'!':  true, 
	'$':  true, 
	'%':  true, 
	'&':  true, 
	'(':  true, 
	')':  true, 
	'*':  true, 
	'+':  true, 
	',':  true, 
	'-':  true, 
	'.':  true, 
	':':  true, 
	';':  true, 
	'=':  true, 
	'[':  true,
	'\'': true, 
	']':  true,
	'_':  true, 
	'~':  true, 
}






































func ValidHeaderFieldValue(v string) bool {
	for i := 0; i < len(v); i++ {
		b := v[i]
		if isCTL(b) && !isLWS(b) {
			return false
		}
	}
	return true
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}



func PunycodeHostPort(v string) (string, error) {
	if isASCII(v) {
		return v, nil
	}

	host, port, err := net.SplitHostPort(v)
	if err != nil {
		
		
		
		host = v
		port = ""
	}
	host, err = idna.ToASCII(host)
	if err != nil {
		
		
		return "", err
	}
	if port == "" {
		return host, nil
	}
	return net.JoinHostPort(host, port), nil
}
