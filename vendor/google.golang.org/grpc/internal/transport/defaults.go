

package transport

import (
	"math"
	"time"
)

const (
	
	defaultWindowSize = 65535
	
	initialWindowSize             = defaultWindowSize 
	infinity                      = time.Duration(math.MaxInt64)
	defaultClientKeepaliveTime    = infinity
	defaultClientKeepaliveTimeout = 20 * time.Second
	defaultMaxStreamsClient       = 100
	defaultMaxConnectionIdle      = infinity
	defaultMaxConnectionAge       = infinity
	defaultMaxConnectionAgeGrace  = infinity
	defaultServerKeepaliveTime    = 2 * time.Hour
	defaultServerKeepaliveTimeout = 20 * time.Second
	defaultKeepalivePolicyMinTime = 5 * time.Minute
	
	maxWindowSize = math.MaxInt32
	
	
	
	defaultWriteQuota              = 64 * 1024
	defaultClientMaxHeaderListSize = uint32(16 << 20)
	defaultServerMaxHeaderListSize = uint32(16 << 20)
)
