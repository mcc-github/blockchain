











package unix

import "unsafe"


type SockaddrDatalink struct {
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [12]int8
	Rcf    uint16
	Route  [16]uint16
	raw    RawSockaddrDatalink
}


func nametomib(name string) (mib []_C_int, err error) {
	const siz = unsafe.Sizeof(mib[0])

	
	
	
	
	
	
	
	var buf [CTL_MAXNAME + 2]_C_int
	n := uintptr(CTL_MAXNAME) * siz

	p := (*byte)(unsafe.Pointer(&buf[0]))
	bytes, err := ByteSliceFromString(name)
	if err != nil {
		return nil, err
	}

	
	
	if err = sysctl([]_C_int{0, 3}, p, &n, &bytes[0], uintptr(len(name))); err != nil {
		return nil, err
	}
	return buf[0 : n/siz], nil
}



func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	p[0], p[1], err = pipe()
	return
}


func Pread(fd int, p []byte, offset int64) (n int, err error) {
	return extpread(fd, p, 0, offset)
}


func Pwrite(fd int, p []byte, offset int64) (n int, err error) {
	return extpwrite(fd, p, 0, offset)
}

func Accept4(fd, flags int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept4(fd, &rsa, &len, flags)
	if err != nil {
		return
	}
	if len > SizeofSockaddrAny {
		panic("RawSockaddrAny too small")
	}
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
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

func sysctlUname(mib []_C_int, old *byte, oldlen *uintptr) error {
	err := sysctl(mib, old, oldlen, nil, 0)
	if err != nil {
		
		
		
		if err == ENOMEM {
			err = nil
		}
	}
	return err
}

func Uname(uname *Utsname) error {
	mib := []_C_int{CTL_KERN, KERN_OSTYPE}
	n := unsafe.Sizeof(uname.Sysname)
	if err := sysctlUname(mib, &uname.Sysname[0], &n); err != nil {
		return err
	}
	uname.Sysname[unsafe.Sizeof(uname.Sysname)-1] = 0

	mib = []_C_int{CTL_KERN, KERN_HOSTNAME}
	n = unsafe.Sizeof(uname.Nodename)
	if err := sysctlUname(mib, &uname.Nodename[0], &n); err != nil {
		return err
	}
	uname.Nodename[unsafe.Sizeof(uname.Nodename)-1] = 0

	mib = []_C_int{CTL_KERN, KERN_OSRELEASE}
	n = unsafe.Sizeof(uname.Release)
	if err := sysctlUname(mib, &uname.Release[0], &n); err != nil {
		return err
	}
	uname.Release[unsafe.Sizeof(uname.Release)-1] = 0

	mib = []_C_int{CTL_KERN, KERN_VERSION}
	n = unsafe.Sizeof(uname.Version)
	if err := sysctlUname(mib, &uname.Version[0], &n); err != nil {
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
	if err := sysctlUname(mib, &uname.Machine[0], &n); err != nil {
		return err
	}
	uname.Machine[unsafe.Sizeof(uname.Machine)-1] = 0

	return nil
}



























































































































































































































































































