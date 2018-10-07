package system 

import (
	"runtime"
	"strings"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)


func IsOSSupported(os string) bool {
	if strings.EqualFold(runtime.GOOS, os) {
		return true
	}
	if LCOWSupported() && strings.EqualFold(os, "linux") {
		return true
	}
	return false
}




func ValidatePlatform(platform specs.Platform) error {
	if runtime.GOOS == "windows" {
		if !(platform.OS == runtime.GOOS || (LCOWSupported() && platform.OS == "linux")) {
			return errors.Errorf("unsupported os %s", platform.OS)
		}
	}
	return nil
}
