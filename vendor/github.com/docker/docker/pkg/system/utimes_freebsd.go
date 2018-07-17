package system 

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)



func LUtimesNano(path string, ts []syscall.Timespec) error {
	var _path *byte
	_path, err := unix.BytePtrFromString(path)
	if err != nil {
		return err
	}

	if _, _, err := unix.Syscall(unix.SYS_LUTIMES, uintptr(unsafe.Pointer(_path)), uintptr(unsafe.Pointer(&ts[0])), 0); err != 0 && err != unix.ENOSYS {
		return err
	}

	return nil
}
