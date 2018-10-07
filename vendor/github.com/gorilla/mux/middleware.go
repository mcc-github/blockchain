package mux

import (
	"net/http"
	"strings"
)




type MiddlewareFunc func(http.Handler) http.Handler


type middleware interface {
	Middleware(handler http.Handler) http.Handler
}


func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mw(handler)
}


func (r *Router) Use(mwf ...MiddlewareFunc) {
	for _, fn := range mwf {
		r.middlewares = append(r.middlewares, fn)
	}
}


func (r *Router) useInterface(mw middleware) {
	r.middlewares = append(r.middlewares, mw)
}





func CORSMethodMiddleware(r *Router) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var allMethods []string

			err := r.Walk(func(route *Route, _ *Router, _ []*Route) error {
				for _, m := range route.matchers {
					if _, ok := m.(*routeRegexp); ok {
						if m.Match(req, &RouteMatch{}) {
							methods, err := route.GetMethods()
							if err != nil {
								return err
							}

							allMethods = append(allMethods, methods...)
						}
						break
					}
				}
				return nil
			})

			if err == nil {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(append(allMethods, "OPTIONS"), ","))

				if req.Method == "OPTIONS" {
					return
				}
			}

			next.ServeHTTP(w, req)
		})
	}
}
