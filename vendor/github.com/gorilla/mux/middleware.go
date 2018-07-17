package mux

import "net/http"




type MiddlewareFunc func(http.Handler) http.Handler


type middleware interface {
	Middleware(handler http.Handler) http.Handler
}


func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mw(handler)
}


func (r *Router) Use(mwf MiddlewareFunc) {
	r.middlewares = append(r.middlewares, mwf)
}


func (r *Router) useInterface(mw middleware) {
	r.middlewares = append(r.middlewares, mw)
}
