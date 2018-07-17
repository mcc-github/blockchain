





package http2

import (
	"crypto/tls"
	"io"
	"net/http"
)

func cloneTLSConfig(c *tls.Config) *tls.Config {
	c2 := c.Clone()
	c2.GetClientCertificate = c.GetClientCertificate 
	return c2
}

var _ http.Pusher = (*responseWriter)(nil)


func (w *responseWriter) Push(target string, opts *http.PushOptions) error {
	internalOpts := pushOptions{}
	if opts != nil {
		internalOpts.Method = opts.Method
		internalOpts.Header = opts.Header
	}
	return w.push(target, internalOpts)
}

func configureServer18(h1 *http.Server, h2 *Server) error {
	if h2.IdleTimeout == 0 {
		if h1.IdleTimeout != 0 {
			h2.IdleTimeout = h1.IdleTimeout
		} else {
			h2.IdleTimeout = h1.ReadTimeout
		}
	}
	return nil
}

func shouldLogPanic(panicValue interface{}) bool {
	return panicValue != nil && panicValue != http.ErrAbortHandler
}

func reqGetBody(req *http.Request) func() (io.ReadCloser, error) {
	return req.GetBody
}

func reqBodyIsNoBody(body io.ReadCloser) bool {
	return body == http.NoBody
}

func go18httpNoBody() io.ReadCloser { return http.NoBody } 
