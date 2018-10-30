















package fileutil

import (
	"os"
	"syscall"
)

func preallocExtend(f *os.File, sizeInBytes int64) error {
	
	err := syscall.Fallocate(int(f.Fd()), 0, 0, sizeInBytes)
	if err != nil {
		errno, ok := err.(syscall.Errno)
		
		
		if ok && (errno == syscall.ENOTSUP || errno == syscall.EINTR) {
			return preallocExtendTrunc(f, sizeInBytes)
		}
	}
	return err
}

func preallocFixed(f *os.File, sizeInBytes int64) error {
	
	err := syscall.Fallocate(int(f.Fd()), 1, 0, sizeInBytes)
	if err != nil {
		errno, ok := err.(syscall.Errno)
		
		if ok && errno == syscall.ENOTSUP {
			return nil
		}
	}
	return err
}
