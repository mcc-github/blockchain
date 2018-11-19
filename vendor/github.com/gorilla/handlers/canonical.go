package handlers

import (
	"net/http"
	"net/url"
	"strings"
)

type canonical struct {
	h      http.Handler
	domain string
	code   int
}
















func CanonicalHost(domain string, code int) func(h http.Handler) http.Handler {
	fn := func(h http.Handler) http.Handler {
		return canonical{h, domain, code}
	}

	return fn
}

func (c canonical) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dest, err := url.Parse(c.domain)
	if err != nil {
		
		c.h.ServeHTTP(w, r)
		return
	}

	if dest.Scheme == "" || dest.Host == "" {
		
		
		c.h.ServeHTTP(w, r)
		return
	}

	if !strings.EqualFold(cleanHost(r.Host), dest.Host) {
		
		dest := dest.Scheme + "://" + dest.Host + r.URL.Path
		if r.URL.RawQuery != "" {
			dest += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, dest, c.code)
		return
	}

	c.h.ServeHTTP(w, r)
}




func cleanHost(in string) string {
	if i := strings.IndexAny(in, " /"); i != -1 {
		return in[:i]
	}
	return in
}
