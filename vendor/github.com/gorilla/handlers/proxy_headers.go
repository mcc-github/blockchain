package handlers

import (
	"net/http"
	"regexp"
	"strings"
)

var (
	
	xForwardedFor    = http.CanonicalHeaderKey("X-Forwarded-For")
	xForwardedHost   = http.CanonicalHeaderKey("X-Forwarded-Host")
	xForwardedProto  = http.CanonicalHeaderKey("X-Forwarded-Proto")
	xForwardedScheme = http.CanonicalHeaderKey("X-Forwarded-Scheme")
	xRealIP          = http.CanonicalHeaderKey("X-Real-IP")
)

var (
	
	
	
	forwarded = http.CanonicalHeaderKey("Forwarded")
	
	
	forRegex = regexp.MustCompile(`(?i)(?:for=)([^(;|,| )]+)`)
	
	
	protoRegex = regexp.MustCompile(`(?i)(?:proto=)(https|http)`)
)













func ProxyHeaders(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		
		if fwd := getIP(r); fwd != "" {
			r.RemoteAddr = fwd
		}

		
		if scheme := getScheme(r); scheme != "" {
			r.URL.Scheme = scheme
		}
		
		if r.Header.Get(xForwardedHost) != "" {
			r.Host = r.Header.Get(xForwardedHost)
		}
		
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}



func getIP(r *http.Request) string {
	var addr string

	if fwd := r.Header.Get(xForwardedFor); fwd != "" {
		
		
		
		s := strings.Index(fwd, ", ")
		if s == -1 {
			s = len(fwd)
		}
		addr = fwd[:s]
	} else if fwd := r.Header.Get(xRealIP); fwd != "" {
		
		
		addr = fwd
	} else if fwd := r.Header.Get(forwarded); fwd != "" {
		
		
		
		
		
		if match := forRegex.FindStringSubmatch(fwd); len(match) > 1 {
			
			
			addr = strings.Trim(match[1], `"`)
		}
	}

	return addr
}



func getScheme(r *http.Request) string {
	var scheme string

	
	if proto := r.Header.Get(xForwardedProto); proto != "" {
		scheme = strings.ToLower(proto)
	} else if proto = r.Header.Get(xForwardedScheme); proto != "" {
		scheme = strings.ToLower(proto)
	} else if proto = r.Header.Get(forwarded); proto != "" {
		
		
		
		
		if match := protoRegex.FindStringSubmatch(proto); len(match) > 1 {
			scheme = strings.ToLower(match[1])
		}
	}

	return scheme
}
