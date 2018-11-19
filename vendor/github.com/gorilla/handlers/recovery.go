package handlers

import (
	"log"
	"net/http"
	"runtime/debug"
)


type RecoveryHandlerLogger interface {
	Println(...interface{})
}

type recoveryHandler struct {
	handler    http.Handler
	logger     RecoveryHandlerLogger
	printStack bool
}




type RecoveryOption func(http.Handler)

func parseRecoveryOptions(h http.Handler, opts ...RecoveryOption) http.Handler {
	for _, option := range opts {
		option(h)
	}

	return h
}













func RecoveryHandler(opts ...RecoveryOption) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		r := &recoveryHandler{handler: h}
		return parseRecoveryOptions(r, opts...)
	}
}



func RecoveryLogger(logger RecoveryHandlerLogger) RecoveryOption {
	return func(h http.Handler) {
		r := h.(*recoveryHandler)
		r.logger = logger
	}
}



func PrintRecoveryStack(print bool) RecoveryOption {
	return func(h http.Handler) {
		r := h.(*recoveryHandler)
		r.printStack = print
	}
}

func (h recoveryHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log(err)
		}
	}()

	h.handler.ServeHTTP(w, req)
}

func (h recoveryHandler) log(v ...interface{}) {
	if h.logger != nil {
		h.logger.Println(v...)
	} else {
		log.Println(v...)
	}

	if h.printStack {
		debug.PrintStack()
	}
}
