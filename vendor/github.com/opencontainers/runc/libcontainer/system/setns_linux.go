package system

import (
	"fmt"
	"runtime"
	"syscall"
)





var setNsMap = map[string]uintptr{
	"linux/386":     346,
	"linux/arm64":   268,
	"linux/amd64":   308,
	"linux/arm":     375,
	"linux/ppc":     350,
	"linux/ppc64":   350,
	"linux/ppc64le": 350,
	"linux/s390x":   339,
}

var sysSetns = setNsMap[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]

func SysSetns() uint32 {
	return uint32(sysSetns)
}

func Setns(fd uintptr, flags uintptr) error {
	ns, exists := setNsMap[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !exists {
		return fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	_, _, err := syscall.RawSyscall(ns, fd, flags, 0)
	if err != 0 {
		return err
	}
	return nil
}
