


package envconfig

import (
	"os"
	"strings"
)

const (
	prefix   = "GRPC_GO_"
	retryStr = prefix + "RETRY"
)

var (
	
	Retry = strings.EqualFold(os.Getenv(retryStr), "on")
)
