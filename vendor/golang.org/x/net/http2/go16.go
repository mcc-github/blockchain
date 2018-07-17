





package http2

import (
	"net/http"
	"time"
)

func transportExpectContinueTimeout(t1 *http.Transport) time.Duration {
	return t1.ExpectContinueTimeout
}
