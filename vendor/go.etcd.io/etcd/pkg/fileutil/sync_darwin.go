















package fileutil

import (
	"os"
	"syscall"
)





func Fsync(f *os.File) error {
	_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), uintptr(syscall.F_FULLFSYNC), uintptr(0))
	if errno == 0 {
		return nil
	}
	return errno
}



func Fdatasync(f *os.File) error {
	return Fsync(f)
}
