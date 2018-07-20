





package backoff

import (
	"time"

	"google.golang.org/grpc/internal/grpcrand"
)




type Strategy interface {
	
	
	Backoff(retries int) time.Duration
}

const (
	
	
	baseDelay = 1.0 * time.Second
	
	factor = 1.6
	
	jitter = 0.2
)



type Exponential struct {
	
	MaxDelay time.Duration
}



func (bc Exponential) Backoff(retries int) time.Duration {
	if retries == 0 {
		return baseDelay
	}
	backoff, max := float64(baseDelay), float64(bc.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= factor
		retries--
	}
	if backoff > max {
		backoff = max
	}
	
	
	backoff *= 1 + jitter*(grpcrand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}
