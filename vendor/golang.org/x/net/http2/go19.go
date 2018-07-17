





package http2

import (
	"net/http"
)

func configureServer19(s *http.Server, conf *Server) error {
	s.RegisterOnShutdown(conf.state.startGracefulShutdown)
	return nil
}
