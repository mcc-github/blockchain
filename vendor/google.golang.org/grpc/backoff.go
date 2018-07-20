




package grpc

import (
	"time"
)



var DefaultBackoffConfig = BackoffConfig{
	MaxDelay: 120 * time.Second,
}


type BackoffConfig struct {
	
	MaxDelay time.Duration
}
