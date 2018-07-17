







package unix

import "syscall"

func Getpagesize() int {
	return syscall.Getpagesize()
}
