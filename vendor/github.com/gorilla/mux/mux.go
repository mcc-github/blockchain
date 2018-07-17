



package mux

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
)

var (
	ErrMethodMismatch = errors.New("method is not allowed")
	ErrNotFound       = errors.New("no matching route was found")
)


func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route), KeepContext: false}
}



















type Router struct {
	
	NotFoundHandler http.Handler

	
	MethodNotAllowedHandler http.Handler

	
	parent parentRoute
	
	routes []*Route
	
	namedRoutes map[string]*Route
	
	strictSlash bool
	
	skipClean bool
	
	
	
	KeepContext bool
	
	useEncodedPath bool
	
	middlewares []middleware
}












func (r *Router) Match(req *http.Request, match *RouteMatch) bool {
	for _, route := range r.routes {
		if route.Match(req, match) {
			
			if match.MatchErr == nil {
				for i := len(r.middlewares) - 1; i >= 0; i-- {
					match.Handler = r.middlewares[i].Middleware(match.Handler)
				}
			}
			return true
		}
	}

	if match.MatchErr == ErrMethodMismatch {
		if r.MethodNotAllowedHandler != nil {
			match.Handler = r.MethodNotAllowedHandler
			return true
		} else {
			return false
		}
	}

	
	if r.NotFoundHandler != nil {
		match.Handler = r.NotFoundHandler
		match.MatchErr = ErrNotFound
		return true
	}

	match.MatchErr = ErrNotFound
	return false
}





func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.skipClean {
		path := req.URL.Path
		if r.useEncodedPath {
			path = req.URL.EscapedPath()
		}
		
		if p := cleanPath(path); p != path {

			
			
			
			url := *req.URL
			url.Path = p
			p = url.String()

			w.Header().Set("Location", p)
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
	}
	var match RouteMatch
	var handler http.Handler
	if r.Match(req, &match) {
		handler = match.Handler
		req = setVars(req, match.Vars)
		req = setCurrentRoute(req, match.Route)
	}

	if handler == nil && match.MatchErr == ErrMethodMismatch {
		handler = methodNotAllowedHandler()
	}

	if handler == nil {
		handler = http.NotFoundHandler()
	}

	if !r.KeepContext {
		defer contextClear(req)
	}

	handler.ServeHTTP(w, req)
}


func (r *Router) Get(name string) *Route {
	return r.getNamedRoutes()[name]
}



func (r *Router) GetRoute(name string) *Route {
	return r.getNamedRoutes()[name]
}




















func (r *Router) StrictSlash(value bool) *Router {
	r.strictSlash = value
	return r
}









func (r *Router) SkipClean(value bool) *Router {
	r.skipClean = value
	return r
}







func (r *Router) UseEncodedPath() *Router {
	r.useEncodedPath = true
	return r
}





func (r *Router) getBuildScheme() string {
	if r.parent != nil {
		return r.parent.getBuildScheme()
	}
	return ""
}


func (r *Router) getNamedRoutes() map[string]*Route {
	if r.namedRoutes == nil {
		if r.parent != nil {
			r.namedRoutes = r.parent.getNamedRoutes()
		} else {
			r.namedRoutes = make(map[string]*Route)
		}
	}
	return r.namedRoutes
}


func (r *Router) getRegexpGroup() *routeRegexpGroup {
	if r.parent != nil {
		return r.parent.getRegexpGroup()
	}
	return nil
}

func (r *Router) buildVars(m map[string]string) map[string]string {
	if r.parent != nil {
		m = r.parent.buildVars(m)
	}
	return m
}






func (r *Router) NewRoute() *Route {
	route := &Route{parent: r, strictSlash: r.strictSlash, skipClean: r.skipClean, useEncodedPath: r.useEncodedPath}
	r.routes = append(r.routes, route)
	return route
}



func (r *Router) Handle(path string, handler http.Handler) *Route {
	return r.NewRoute().Path(path).Handler(handler)
}



func (r *Router) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *Route {
	return r.NewRoute().Path(path).HandlerFunc(f)
}



func (r *Router) Headers(pairs ...string) *Route {
	return r.NewRoute().Headers(pairs...)
}



func (r *Router) Host(tpl string) *Route {
	return r.NewRoute().Host(tpl)
}



func (r *Router) MatcherFunc(f MatcherFunc) *Route {
	return r.NewRoute().MatcherFunc(f)
}



func (r *Router) Methods(methods ...string) *Route {
	return r.NewRoute().Methods(methods...)
}



func (r *Router) Path(tpl string) *Route {
	return r.NewRoute().Path(tpl)
}



func (r *Router) PathPrefix(tpl string) *Route {
	return r.NewRoute().PathPrefix(tpl)
}



func (r *Router) Queries(pairs ...string) *Route {
	return r.NewRoute().Queries(pairs...)
}



func (r *Router) Schemes(schemes ...string) *Route {
	return r.NewRoute().Schemes(schemes...)
}



func (r *Router) BuildVarsFunc(f BuildVarsFunc) *Route {
	return r.NewRoute().BuildVarsFunc(f)
}




func (r *Router) Walk(walkFn WalkFunc) error {
	return r.walk(walkFn, []*Route{})
}



var SkipRouter = errors.New("skip this router")




type WalkFunc func(route *Route, router *Router, ancestors []*Route) error

func (r *Router) walk(walkFn WalkFunc, ancestors []*Route) error {
	for _, t := range r.routes {
		err := walkFn(t, r, ancestors)
		if err == SkipRouter {
			continue
		}
		if err != nil {
			return err
		}
		for _, sr := range t.matchers {
			if h, ok := sr.(*Router); ok {
				ancestors = append(ancestors, t)
				err := h.walk(walkFn, ancestors)
				if err != nil {
					return err
				}
				ancestors = ancestors[:len(ancestors)-1]
			}
		}
		if h, ok := t.handler.(*Router); ok {
			ancestors = append(ancestors, t)
			err := h.walk(walkFn, ancestors)
			if err != nil {
				return err
			}
			ancestors = ancestors[:len(ancestors)-1]
		}
	}
	return nil
}






type RouteMatch struct {
	Route   *Route
	Handler http.Handler
	Vars    map[string]string

	
	
	
	MatchErr error
}

type contextKey int

const (
	varsKey contextKey = iota
	routeKey
)


func Vars(r *http.Request) map[string]string {
	if rv := contextGet(r, varsKey); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}






func CurrentRoute(r *http.Request) *Route {
	if rv := contextGet(r, routeKey); rv != nil {
		return rv.(*Route)
	}
	return nil
}

func setVars(r *http.Request, val interface{}) *http.Request {
	return contextSet(r, varsKey, val)
}

func setCurrentRoute(r *http.Request, val interface{}) *http.Request {
	return contextSet(r, routeKey, val)
}







func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	
	
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}


func uniqueVars(s1, s2 []string) error {
	for _, v1 := range s1 {
		for _, v2 := range s2 {
			if v1 == v2 {
				return fmt.Errorf("mux: duplicated route variable %q", v2)
			}
		}
	}
	return nil
}



func checkPairs(pairs ...string) (int, error) {
	length := len(pairs)
	if length%2 != 0 {
		return length, fmt.Errorf(
			"mux: number of parameters must be multiple of 2, got %v", pairs)
	}
	return length, nil
}



func mapFromPairsToString(pairs ...string) (map[string]string, error) {
	length, err := checkPairs(pairs...)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m, nil
}



func mapFromPairsToRegex(pairs ...string) (map[string]*regexp.Regexp, error) {
	length, err := checkPairs(pairs...)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*regexp.Regexp, length/2)
	for i := 0; i < length; i += 2 {
		regex, err := regexp.Compile(pairs[i+1])
		if err != nil {
			return nil, err
		}
		m[pairs[i]] = regex
	}
	return m, nil
}


func matchInArray(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}


func matchMapWithString(toCheck map[string]string, toMatch map[string][]string, canonicalKey bool) bool {
	for k, v := range toCheck {
		
		if canonicalKey {
			k = http.CanonicalHeaderKey(k)
		}
		if values := toMatch[k]; values == nil {
			return false
		} else if v != "" {
			
			
			valueExists := false
			for _, value := range values {
				if v == value {
					valueExists = true
					break
				}
			}
			if !valueExists {
				return false
			}
		}
	}
	return true
}



func matchMapWithRegex(toCheck map[string]*regexp.Regexp, toMatch map[string][]string, canonicalKey bool) bool {
	for k, v := range toCheck {
		
		if canonicalKey {
			k = http.CanonicalHeaderKey(k)
		}
		if values := toMatch[k]; values == nil {
			return false
		} else if v != nil {
			
			
			valueExists := false
			for _, value := range values {
				if v.MatchString(value) {
					valueExists = true
					break
				}
			}
			if !valueExists {
				return false
			}
		}
	}
	return true
}


func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}



func methodNotAllowedHandler() http.Handler { return http.HandlerFunc(methodNotAllowed) }
