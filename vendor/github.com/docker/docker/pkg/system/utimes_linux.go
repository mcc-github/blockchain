package system 

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)



func LUtimesNano(path string, ts []syscall.Timespec) error {
	atFdCwd := unix.AT_FDCWD

	var _path *byte
	_path, err := unix.BytePtrFromString(path)
	if err != nil {
		return err
	}
	if _, _, err := unix.Syscall6(unix.SYS_UTIMENSAT, uintptr(atFdCwd), uintptr(unsafe.Pointer(_path)), uintptr(unsafe.Pointer(&ts[0])), unix.AT_SYMLINK_NOFOLLOW, 0, 0); err != 0 && err != unix.ENOSYS {
		return err
	}

	return nil
}
