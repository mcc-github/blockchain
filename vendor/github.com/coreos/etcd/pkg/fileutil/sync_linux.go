















package fileutil

import (
	"os"
	"syscall"
)


func Fsync(f *os.File) error {
	return f.Sync()
}




func Fdatasync(f *os.File) error {
	return syscall.Fdatasync(int(f.Fd()))
}
