







package unix

import (
	"syscall"
	"unsafe"
)

func seek(fd int, offset int64, whence int) (newoffset int64, err syscall.Errno) {
	offsetLow := uint32(offset & 0xffffffff)
	offsetHigh := uint32((offset >> 32) & 0xffffffff)
	_, _, err = Syscall6(SYS__LLSEEK, uintptr(fd), uintptr(offsetHigh), uintptr(offsetLow), uintptr(unsafe.Pointer(&newoffset)), uintptr(whence), 0)
	return newoffset, err
}