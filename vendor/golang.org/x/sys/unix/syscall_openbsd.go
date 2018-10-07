











package unix

import (
	"sort"
	"syscall"
	"unsafe"
)


type SockaddrDatalink struct {
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [24]int8
	raw    RawSockaddrDatalink
}

func Syscall9(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r1, r2 uintptr, err syscall.Errno)

func nametomib(name string) (mib []_C_int, err error) {
	i := sort.Search(len(sysctlMib), func(i int) bool {
		return sysctlMib[i].ctlname >= name
	})
	if i < len(sysctlMib) && sysctlMib[i].ctlname == name {
		return sysctlMib[i].ctloid, nil
	}
	return nil, EINVAL
}

func SysctlUvmexp(name string) (*Uvmexp, error) {
	mib, err := sysctlmib(name)
	if err != nil {
		return nil, err
	}

	n := uintptr(SizeofUvmexp)
	var u Uvmexp
	if err := sysctl(mib, (*byte)(unsafe.Pointer(&u)), &n, nil, 0); err != nil {
		return nil, err
	}
	if n != SizeofUvmexp {
		return nil, EIO
	}
	return &u, nil
}


func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	err = pipe(&pp)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
}


func Getdirentries(fd int, buf []byte, basep *uintptr) (n int, err error) {
	return getdents(fd, buf)
}

const ImplementsGetwd = true



func Getwd() (string, error) {
	var buf [PathMax]byte
	_, err := Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	n := clen(buf[:])
	if n < 1 {
		return "", EINVAL
	}
	return string(buf[:n]), nil
}


func sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	return -1, ENOSYS
}

func Getfsstat(buf []Statfs_t, flags int) (n int, err error) {
	var _p0 unsafe.Pointer
	var bufsize uintptr
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
		bufsize = unsafe.Sizeof(Statfs_t{}) * uintptr(len(buf))
	}
	r0, _, e1 := Syscall(SYS_GETFSSTAT, uintptr(_p0), bufsize, uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

func setattrlistTimes(path string, times []Timespec, flags int) error {
	
	return ENOSYS
}








func IoctlSetInt(fd int, req uint, value int) error {
	return ioctl(fd, req, uintptr(value))
}

func ioctlSetWinsize(fd int, req uint, value *Winsize) error {
	return ioctl(fd, req, uintptr(unsafe.Pointer(value)))
}

func ioctlSetTermios(fd int, req uint, value *Termios) error {
	return ioctl(fd, req, uintptr(unsafe.Pointer(value)))
}



func IoctlGetInt(fd int, req uint) (int, error) {
	var value int
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return value, err
}

func IoctlGetWinsize(fd int, req uint) (*Winsize, error) {
	var value Winsize
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
}

func IoctlGetTermios(fd int, req uint) (*Termios, error) {
	var value Termios
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
}

func Uname(uname *Utsname) error {
	mib := []_C_int{CTL_KERN, KERN_OSTYPE}
	n := unsafe.Sizeof(uname.Sysname)
	if err := sysctl(mib, &uname.Sysname[0], &n, nil, 0); err != nil {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_HOSTNAME}
	n = unsafe.Sizeof(uname.Nodename)
	if err := sysctl(mib, &uname.Nodename[0], &n, nil, 0); err != nil {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_OSRELEASE}
	n = unsafe.Sizeof(uname.Release)
	if err := sysctl(mib, &uname.Release[0], &n, nil, 0); err != nil {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_VERSION}
	n = unsafe.Sizeof(uname.Version)
	if err := sysctl(mib, &uname.Version[0], &n, nil, 0); err != nil {
		return err
	}

	
	
	for i, b := range uname.Version {
		if b == '\n' || b == '\t' {
			if i == len(uname.Version)-1 {
				uname.Version[i] = 0
			} else {
				uname.Version[i] = ' '
			}
		}
	}

	mib = []_C_int{CTL_HW, HW_MACHINE}
	n = unsafe.Sizeof(uname.Machine)
	if err := sysctl(mib, &uname.Machine[0], &n, nil, 0); err != nil {
		return err
	}

	return nil
}
















































































































































































