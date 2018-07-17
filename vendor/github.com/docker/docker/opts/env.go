package opts 

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)









func ValidateEnv(val string) (string, error) {
	arr := strings.Split(val, "=")
	if arr[0] == "" {
		return "", errors.Errorf("invalid environment variable: %s", val)
	}
	if len(arr) > 1 {
		return val, nil
	}
	if !doesEnvExist(val) {
		return val, nil
	}
	return fmt.Sprintf("%s=%s", val, os.Getenv(val)), nil
}

func doesEnvExist(name string) bool {
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if runtime.GOOS == "windows" {
			
			if strings.EqualFold(parts[0], name) {
				return true
			}
		}
		if parts[0] == name {
			return true
		}
	}
	return false
}
