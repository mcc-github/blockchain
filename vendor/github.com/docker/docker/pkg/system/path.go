package system 

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/containerd/continuity/pathdriver"
)

const defaultUnixPathEnv = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"




func DefaultPathEnv(os string) string {
	if runtime.GOOS == "windows" {
		if os != runtime.GOOS {
			return defaultUnixPathEnv
		}
		
		
		return ""
	}
	return defaultUnixPathEnv

}
















func CheckSystemDriveAndRemoveDriveLetter(path string, driver pathdriver.PathDriver) (string, error) {
	if runtime.GOOS != "windows" || LCOWSupported() {
		return path, nil
	}

	if len(path) == 2 && string(path[1]) == ":" {
		return "", fmt.Errorf("No relative path specified in %q", path)
	}
	if !driver.IsAbs(path) || len(path) < 2 {
		return filepath.FromSlash(path), nil
	}
	if string(path[1]) == ":" && !strings.EqualFold(string(path[0]), "c") {
		return "", fmt.Errorf("The specified path is not on the system drive (C:)")
	}
	return filepath.FromSlash(path[2:]), nil
}
