











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
	Data   [46]int8
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

func GetsockoptIPMreqn(fd, level, opt int) (*IPMreqn, error) {
	var value IPMreqn
	vallen := _Socklen(SizeofIPMreqn)
	errno := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, errno
}

func SetsockoptIPMreqn(fd, level, opt int, mreq *IPMreqn) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), unsafe.Sizeof(*mreq))
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
	sa, err = anyToSockaddr(&rsa)
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



func xattrnamespace(fullattr string) (ns int, attr string, err error) {
	s := -1
	for idx, val := range fullattr {
		if val == '.' {
			s = idx
			break
		}
	}

	if s == -1 {
		return -1, "", ENOATTR
	}

	namespace := fullattr[0:s]
	attr = fullattr[s+1:]

	switch namespace {
	case "user":
		return EXTATTR_NAMESPACE_USER, attr, nil
	case "system":
		return EXTATTR_NAMESPACE_SYSTEM, attr, nil
	default:
		return -1, "", ENOATTR
	}
}

func initxattrdest(dest []byte, idx int) (d unsafe.Pointer) {
	if len(dest) > idx {
		return unsafe.Pointer(&dest[idx])
	} else {
		return unsafe.Pointer(_zero)
	}
}



func Getxattr(file string, attr string, dest []byte) (sz int, err error) {
	d := initxattrdest(dest, 0)
	destsize := len(dest)

	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return -1, err
	}

	return ExtattrGetFile(file, nsid, a, uintptr(d), destsize)
}

func Fgetxattr(fd int, attr string, dest []byte) (sz int, err error) {
	d := initxattrdest(dest, 0)
	destsize := len(dest)

	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return -1, err
	}

	return ExtattrGetFd(fd, nsid, a, uintptr(d), destsize)
}

func Lgetxattr(link string, attr string, dest []byte) (sz int, err error) {
	d := initxattrdest(dest, 0)
	destsize := len(dest)

	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return -1, err
	}

	return ExtattrGetLink(link, nsid, a, uintptr(d), destsize)
}



func Fsetxattr(fd int, attr string, data []byte, flags int) (err error) {
	d := unsafe.Pointer(&data[0])
	datasiz := len(data)

	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return
	}

	_, err = ExtattrSetFd(fd, nsid, a, uintptr(d), datasiz)
	return
}

func Setxattr(file string, attr string, data []byte, flags int) (err error) {
	d := unsafe.Pointer(&data[0])
	datasiz := len(data)

	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return
	}

	_, err = ExtattrSetFile(file, nsid, a, uintptr(d), datasiz)
	return
}

func Lsetxattr(link string, attr string, data []byte, flags int) (err error) {
	d := unsafe.Pointer(&data[0])
	datasiz := len(data)

	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return
	}

	_, err = ExtattrSetLink(link, nsid, a, uintptr(d), datasiz)
	return
}

func Removexattr(file string, attr string) (err error) {
	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return
	}

	err = ExtattrDeleteFile(file, nsid, a)
	return
}

func Fremovexattr(fd int, attr string) (err error) {
	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return
	}

	err = ExtattrDeleteFd(fd, nsid, a)
	return
}

func Lremovexattr(link string, attr string) (err error) {
	nsid, a, err := xattrnamespace(attr)
	if err != nil {
		return
	}

	err = ExtattrDeleteLink(link, nsid, a)
	return
}

func Listxattr(file string, dest []byte) (sz int, err error) {
	d := initxattrdest(dest, 0)
	destsiz := len(dest)

	
	s := 0
	for _, nsid := range [...]int{EXTATTR_NAMESPACE_USER, EXTATTR_NAMESPACE_SYSTEM} {
		stmp, e := ExtattrListFile(file, nsid, uintptr(d), destsiz)

		
		if e != nil && e == EPERM && nsid != EXTATTR_NAMESPACE_USER {
			continue
		} else if e != nil {
			return s, e
		}

		s += stmp
		destsiz -= s
		if destsiz < 0 {
			destsiz = 0
		}
		d = initxattrdest(dest, s)
	}

	return s, nil
}

func Flistxattr(fd int, dest []byte) (sz int, err error) {
	d := initxattrdest(dest, 0)
	destsiz := len(dest)

	s := 0
	for _, nsid := range [...]int{EXTATTR_NAMESPACE_USER, EXTATTR_NAMESPACE_SYSTEM} {
		stmp, e := ExtattrListFd(fd, nsid, uintptr(d), destsiz)
		if e != nil && e == EPERM && nsid != EXTATTR_NAMESPACE_USER {
			continue
		} else if e != nil {
			return s, e
		}

		s += stmp
		destsiz -= s
		if destsiz < 0 {
			destsiz = 0
		}
		d = initxattrdest(dest, s)
	}

	return s, nil
}

func Llistxattr(link string, dest []byte) (sz int, err error) {
	d := initxattrdest(dest, 0)
	destsiz := len(dest)

	s := 0
	for _, nsid := range [...]int{EXTATTR_NAMESPACE_USER, EXTATTR_NAMESPACE_SYSTEM} {
		stmp, e := ExtattrListLink(link, nsid, uintptr(d), destsiz)
		if e != nil && e == EPERM && nsid != EXTATTR_NAMESPACE_USER {
			continue
		} else if e != nil {
			return s, e
		}

		s += stmp
		destsiz -= s
		if destsiz < 0 {
			destsiz = 0
		}
		d = initxattrdest(dest, s)
	}

	return s, nil
}








func IoctlSetInt(fd int, req uint, value int) error {
	return ioctl(fd, req, uintptr(value))
}

func IoctlSetWinsize(fd int, req uint, value *Winsize) error {
	return ioctl(fd, req, uintptr(unsafe.Pointer(value)))
}

func IoctlSetTermios(fd int, req uint, value *Termios) error {
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

























































































































































































































































































































