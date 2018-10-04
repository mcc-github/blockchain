






package unix

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	_SYS_PLEDGE = 108
)










func Pledge(promises, execpromises string) error {
	maj, min, err := majmin()
	if err != nil {
		return err
	}

	
	if (maj == 5 && min != 9) || maj < 5 {
		return fmt.Errorf("pledge syscall is not available on OpenBSD %d.%d", maj, min)
	}

	
	
	if (maj < 6 || (maj == 6 && min <= 2)) && execpromises != "" {
		return fmt.Errorf("cannot use execpromises on OpenBSD %d.%d", maj, min)
	}

	pptr, err := syscall.BytePtrFromString(promises)
	if err != nil {
		return err
	}

	
	
	var expr unsafe.Pointer

	
	if maj > 6 || (maj == 6 && min > 2) {
		exptr, err := syscall.BytePtrFromString(execpromises)
		if err != nil {
			return err
		}
		expr = unsafe.Pointer(exptr)
	}

	_, _, e := syscall.Syscall(_SYS_PLEDGE, uintptr(unsafe.Pointer(pptr)), uintptr(expr), 0)
	if e != 0 {
		return e
	}

	return nil
}


func majmin() (major int, minor int, err error) {
	var v Utsname
	err = Uname(&v)
	if err != nil {
		return
	}

	major, err = strconv.Atoi(string(v.Release[0]))
	if err != nil {
		err = errors.New("cannot parse major version number returned by uname")
		return
	}

	minor, err = strconv.Atoi(string(v.Release[2]))
	if err != nil {
		err = errors.New("cannot parse minor version number returned by uname")
		return
	}

	return
}
