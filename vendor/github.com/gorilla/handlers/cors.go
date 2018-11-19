package handlers

import (
	"net/http"
	"strconv"
	"strings"
)


type CORSOption func(*cors) error

type cors struct {
	h                      http.Handler
	allowedHeaders         []string
	allowedMethods         []string
	allowedOrigins         []string
	allowedOriginValidator OriginValidator
	exposedHeaders         []string
	maxAge                 int
	ignoreOptions          bool
	allowCredentials       bool
}


type OriginValidator func(string) bool

var (
	defaultCorsMethods = []string{"GET", "HEAD", "POST"}
	defaultCorsHeaders = []string{"Accept", "Accept-Language", "Content-Language", "Origin"}
	
)

const (
	corsOptionMethod           string = "OPTIONS"
	corsAllowOriginHeader      string = "Access-Control-Allow-Origin"
	corsExposeHeadersHeader    string = "Access-Control-Expose-Headers"
	corsMaxAgeHeader           string = "Access-Control-Max-Age"
	corsAllowMethodsHeader     string = "Access-Control-Allow-Methods"
	corsAllowHeadersHeader     string = "Access-Control-Allow-Headers"
	corsAllowCredentialsHeader string = "Access-Control-Allow-Credentials"
	corsRequestMethodHeader    string = "Access-Control-Request-Method"
	corsRequestHeadersHeader   string = "Access-Control-Request-Headers"
	corsOriginHeader           string = "Origin"
	corsVaryHeader             string = "Vary"
	corsOriginMatchAll         string = "*"
)

func (ch *cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get(corsOriginHeader)
	if !ch.isOriginAllowed(origin) {
		if r.Method != corsOptionMethod || ch.ignoreOptions {
			ch.h.ServeHTTP(w, r)
		}

		return
	}

	if r.Method == corsOptionMethod {
		if ch.ignoreOptions {
			ch.h.ServeHTTP(w, r)
			return
		}

		if _, ok := r.Header[corsRequestMethodHeader]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		method := r.Header.Get(corsRequestMethodHeader)
		if !ch.isMatch(method, ch.allowedMethods) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		requestHeaders := strings.Split(r.Header.Get(corsRequestHeadersHeader), ",")
		allowedHeaders := []string{}
		for _, v := range requestHeaders {
			canonicalHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
			if canonicalHeader == "" || ch.isMatch(canonicalHeader, defaultCorsHeaders) {
				continue
			}

			if !ch.isMatch(canonicalHeader, ch.allowedHeaders) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			allowedHeaders = append(allowedHeaders, canonicalHeader)
		}

		if len(allowedHeaders) > 0 {
			w.Header().Set(corsAllowHeadersHeader, strings.Join(allowedHeaders, ","))
		}

		if ch.maxAge > 0 {
			w.Header().Set(corsMaxAgeHeader, strconv.Itoa(ch.maxAge))
		}

		if !ch.isMatch(method, defaultCorsMethods) {
			w.Header().Set(corsAllowMethodsHeader, method)
		}
	} else {
		if len(ch.exposedHeaders) > 0 {
			w.Header().Set(corsExposeHeadersHeader, strings.Join(ch.exposedHeaders, ","))
		}
	}

	if ch.allowCredentials {
		w.Header().Set(corsAllowCredentialsHeader, "true")
	}

	if len(ch.allowedOrigins) > 1 {
		w.Header().Set(corsVaryHeader, corsOriginHeader)
	}

	returnOrigin := origin
	if ch.allowedOriginValidator == nil && len(ch.allowedOrigins) == 0 {
		returnOrigin = "*"
	} else {
		for _, o := range ch.allowedOrigins {
			
			
			
			if o == corsOriginMatchAll {
				returnOrigin = "*"
				break
			}
		}
	}
	w.Header().Set(corsAllowOriginHeader, returnOrigin)

	if r.Method == corsOptionMethod {
		return
	}
	ch.h.ServeHTTP(w, r)
}




















func CORS(opts ...CORSOption) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		ch := parseCORSOptions(opts...)
		ch.h = h
		return ch
	}
}

func parseCORSOptions(opts ...CORSOption) *cors {
	ch := &cors{
		allowedMethods: defaultCorsMethods,
		allowedHeaders: defaultCorsHeaders,
		allowedOrigins: []string{},
	}

	for _, option := range opts {
		option(ch)
	}

	return ch
}











func AllowedHeaders(headers []string) CORSOption {
	return func(ch *cors) error {
		for _, v := range headers {
			normalizedHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
			if normalizedHeader == "" {
				continue
			}

			if !ch.isMatch(normalizedHeader, ch.allowedHeaders) {
				ch.allowedHeaders = append(ch.allowedHeaders, normalizedHeader)
			}
		}

		return nil
	}
}





func AllowedMethods(methods []string) CORSOption {
	return func(ch *cors) error {
		ch.allowedMethods = []string{}
		for _, v := range methods {
			normalizedMethod := strings.ToUpper(strings.TrimSpace(v))
			if normalizedMethod == "" {
				continue
			}

			if !ch.isMatch(normalizedMethod, ch.allowedMethods) {
				ch.allowedMethods = append(ch.allowedMethods, normalizedMethod)
			}
		}

		return nil
	}
}




func AllowedOrigins(origins []string) CORSOption {
	return func(ch *cors) error {
		for _, v := range origins {
			if v == corsOriginMatchAll {
				ch.allowedOrigins = []string{corsOriginMatchAll}
				return nil
			}
		}

		ch.allowedOrigins = origins
		return nil
	}
}



func AllowedOriginValidator(fn OriginValidator) CORSOption {
	return func(ch *cors) error {
		ch.allowedOriginValidator = fn
		return nil
	}
}



func ExposedHeaders(headers []string) CORSOption {
	return func(ch *cors) error {
		ch.exposedHeaders = []string{}
		for _, v := range headers {
			normalizedHeader := http.CanonicalHeaderKey(strings.TrimSpace(v))
			if normalizedHeader == "" {
				continue
			}

			if !ch.isMatch(normalizedHeader, ch.exposedHeaders) {
				ch.exposedHeaders = append(ch.exposedHeaders, normalizedHeader)
			}
		}

		return nil
	}
}




func MaxAge(age int) CORSOption {
	return func(ch *cors) error {
		
		if age > 600 {
			age = 600
		}

		ch.maxAge = age
		return nil
	}
}




func IgnoreOptions() CORSOption {
	return func(ch *cors) error {
		ch.ignoreOptions = true
		return nil
	}
}



func AllowCredentials() CORSOption {
	return func(ch *cors) error {
		ch.allowCredentials = true
		return nil
	}
}

func (ch *cors) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	if ch.allowedOriginValidator != nil {
		return ch.allowedOriginValidator(origin)
	}

	if len(ch.allowedOrigins) == 0 {
		return true
	}

	for _, allowedOrigin := range ch.allowedOrigins {
		if allowedOrigin == origin || allowedOrigin == corsOriginMatchAll {
			return true
		}
	}

	return false
}

func (ch *cors) isMatch(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}
