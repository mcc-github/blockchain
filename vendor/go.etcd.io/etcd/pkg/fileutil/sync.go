















package fileutil

import "os"


func Fsync(f *os.File) error {
	return f.Sync()
}


func Fdatasync(f *os.File) error {
	return f.Sync()
}
