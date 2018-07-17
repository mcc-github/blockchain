





package http2

import (
	"io"
	"net/http"
)

func configureServer18(h1 *http.Server, h2 *Server) error {
	
	return nil
}

func shouldLogPanic(panicValue interface{}) bool {
	return panicValue != nil
}

func reqGetBody(req *http.Request) func() (io.ReadCloser, error) {
	return nil
}

func reqBodyIsNoBody(io.ReadCloser) bool { return false }

func go18httpNoBody() io.ReadCloser { return nil } 
