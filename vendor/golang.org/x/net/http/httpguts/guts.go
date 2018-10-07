








package httpguts

import (
	"net/textproto"
	"strings"
)




func ValidTrailerHeader(name string) bool {
	name = textproto.CanonicalMIMEHeaderKey(name)
	if strings.HasPrefix(name, "If-") || badTrailer[name] {
		return false
	}
	return true
}

var badTrailer = map[string]bool{
	"Authorization":       true,
	"Cache-Control":       true,
	"Connection":          true,
	"Content-Encoding":    true,
	"Content-Length":      true,
	"Content-Range":       true,
	"Content-Type":        true,
	"Expect":              true,
	"Host":                true,
	"Keep-Alive":          true,
	"Max-Forwards":        true,
	"Pragma":              true,
	"Proxy-Authenticate":  true,
	"Proxy-Authorization": true,
	"Proxy-Connection":    true,
	"Range":               true,
	"Realm":               true,
	"Te":                  true,
	"Trailer":             true,
	"Transfer-Encoding":   true,
	"Www-Authenticate":    true,
}
