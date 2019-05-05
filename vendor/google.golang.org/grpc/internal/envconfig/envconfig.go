


package envconfig

import (
	"os"
	"strings"
)

const (
	prefix              = "GRPC_GO_"
	retryStr            = prefix + "RETRY"
	requireHandshakeStr = prefix + "REQUIRE_HANDSHAKE"
)


type RequireHandshakeSetting int

const (
	
	
	RequireHandshakeOn RequireHandshakeSetting = iota
	
	
	RequireHandshakeOff
)

var (
	
	Retry = strings.EqualFold(os.Getenv(retryStr), "on")
	
	
	
	
	RequireHandshake = RequireHandshakeOn
)

func init() {
	switch strings.ToLower(os.Getenv(requireHandshakeStr)) {
	case "on":
		fallthrough
	default:
		RequireHandshake = RequireHandshakeOn
	case "off":
		RequireHandshake = RequireHandshakeOff
	}
}
