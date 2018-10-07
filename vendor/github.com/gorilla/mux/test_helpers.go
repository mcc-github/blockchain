



package mux

import "net/http"









func SetURLVars(r *http.Request, val map[string]string) *http.Request {
	return setVars(r, val)
}
