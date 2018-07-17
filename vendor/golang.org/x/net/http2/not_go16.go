





package http2

import (
	"net/http"
	"time"
)

func configureTransport(t1 *http.Transport) (*Transport, error) {
	return nil, errTransportVersion
}

func transportExpectContinueTimeout(t1 *http.Transport) time.Duration {
	return 0

}
