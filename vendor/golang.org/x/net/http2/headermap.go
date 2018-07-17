



package http2

import (
	"net/http"
	"strings"
)

var (
	commonLowerHeader = map[string]string{} 
	commonCanonHeader = map[string]string{} 
)

func init() {
	for _, v := range []string{
		"accept",
		"accept-charset",
		"accept-encoding",
		"accept-language",
		"accept-ranges",
		"age",
		"access-control-allow-origin",
		"allow",
		"authorization",
		"cache-control",
		"content-disposition",
		"content-encoding",
		"content-language",
		"content-length",
		"content-location",
		"content-range",
		"content-type",
		"cookie",
		"date",
		"etag",
		"expect",
		"expires",
		"from",
		"host",
		"if-match",
		"if-modified-since",
		"if-none-match",
		"if-unmodified-since",
		"last-modified",
		"link",
		"location",
		"max-forwards",
		"proxy-authenticate",
		"proxy-authorization",
		"range",
		"referer",
		"refresh",
		"retry-after",
		"server",
		"set-cookie",
		"strict-transport-security",
		"trailer",
		"transfer-encoding",
		"user-agent",
		"vary",
		"via",
		"www-authenticate",
	} {
		chk := http.CanonicalHeaderKey(v)
		commonLowerHeader[chk] = v
		commonCanonHeader[v] = chk
	}
}

func lowerHeader(v string) string {
	if s, ok := commonLowerHeader[v]; ok {
		return s
	}
	return strings.ToLower(v)
}
