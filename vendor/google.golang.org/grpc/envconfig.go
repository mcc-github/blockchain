

package grpc

import (
	"os"
	"strings"
)

const (
	envConfigPrefix        = "GRPC_GO_"
	envConfigStickinessStr = envConfigPrefix + "STICKINESS"
)

var (
	envConfigStickinessOn bool
)

func init() {
	envConfigStickinessOn = strings.EqualFold(os.Getenv(envConfigStickinessStr), "on")
}
