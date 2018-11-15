/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler


type Chain struct {
	mw []Middleware
}



func NewChain(middlewares ...Middleware) Chain {
	return Chain{
		mw: append([]Middleware{}, middlewares...),
	}
}


func (c Chain) Handler(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := len(c.mw) - 1; i >= 0; i-- {
		h = c.mw[i](h)
	}
	return h
}
