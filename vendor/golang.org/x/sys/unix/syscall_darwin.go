











package unix

import (
	errorspkg "errors"
	"syscall"
	"unsafe"
)

const ImplementsGetwd = true

func Getwd() (string, error) {
	buf := make([]byte, 2048)
	attrs, err := getAttrList(".", attrList{CommonAttr: attrCmnFullpath}, buf, 0)
	if err == nil && len(attrs) == 1 && len(attrs[0]) >= 2 {
		wd := string(attrs[0])
		
		
		if wd[0] == '/' && wd[len(wd)-1] == 0 {
			return wd[:len(wd)-1], nil
		}
	}
	
	
	return "", ENOTSUP
}


type SockaddrDatalink struct {
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [12]int8
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


func PtraceAttach(pid int) (err error) { return ptrace(PT_ATTACH, pid, 0, 0) }
func PtraceDetach(pid int) (err error) { return ptrace(PT_DETACH, pid, 0, 0) }

const (
	attrBitMapCount = 5
	attrCmnFullpath = 0x08000000
)

type attrList struct {
	bitmapCount uint16
	_           uint16
	CommonAttr  uint32
	VolAttr     uint32
	DirAttr     uint32
	FileAttr    uint32
	Forkattr    uint32
}

func getAttrList(path string, attrList attrList, attrBuf []byte, options uint) (attrs [][]byte, err error) {
	if len(attrBuf) < 4 {
		return nil, errorspkg.New("attrBuf too small")
	}
	attrList.bitmapCount = attrBitMapCount

	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return nil, err
	}

	_, _, e1 := Syscall6(
		SYS_GETATTRLIST,
		uintptr(unsafe.Pointer(_p0)),
		uintptr(unsafe.Pointer(&attrList)),
		uintptr(unsafe.Pointer(&attrBuf[0])),
		uintptr(len(attrBuf)),
		uintptr(options),
		0,
	)
	if e1 != 0 {
		return nil, e1
	}
	size := *(*uint32)(unsafe.Pointer(&attrBuf[0]))

	
	
	
	dat := attrBuf
	if int(size) < len(attrBuf) {
		dat = dat[:size]
	}
	dat = dat[4:] 

	for i := uint32(0); int(i) < len(dat); {
		header := dat[i:]
		if len(header) < 8 {
			return attrs, errorspkg.New("truncated attribute header")
		}
		datOff := *(*int32)(unsafe.Pointer(&header[0]))
		attrLen := *(*uint32)(unsafe.Pointer(&header[4]))
		if datOff < 0 || uint32(datOff)+attrLen > uint32(len(dat)) {
			return attrs, errorspkg.New("truncated results; attrBuf too small")
		}
		end := uint32(datOff) + attrLen
		attrs = append(attrs, dat[datOff:end])
		i = end
		if r := i % 4; r != 0 {
			i += (4 - r)
		}
	}
	return
}



func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	p[0], p[1], err = pipe()
	return
}

func Getfsstat(buf []Statfs_t, flags int) (n int, err error) {
	var _p0 unsafe.Pointer
	var bufsize uintptr
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
		bufsize = unsafe.Sizeof(Statfs_t{}) * uintptr(len(buf))
	}
	r0, _, e1 := Syscall(SYS_GETFSSTAT64, uintptr(_p0), bufsize, uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

func setattrlistTimes(path string, times []Timespec, flags int) error {
	_p0, err := BytePtrFromString(path)
	if err != nil {
		return err
	}

	var attrList attrList
	attrList.bitmapCount = ATTR_BIT_MAP_COUNT
	attrList.CommonAttr = ATTR_CMN_MODTIME | ATTR_CMN_ACCTIME

	
	attributes := [2]Timespec{times[1], times[0]}
	options := 0
	if flags&AT_SYMLINK_NOFOLLOW != 0 {
		options |= FSOPT_NOFOLLOW
	}
	_, _, e1 := Syscall6(
		SYS_SETATTRLIST,
		uintptr(unsafe.Pointer(_p0)),
		uintptr(unsafe.Pointer(&attrList)),
		uintptr(unsafe.Pointer(&attributes)),
		uintptr(unsafe.Sizeof(attributes)),
		uintptr(options),
		0,
	)
	if e1 != 0 {
		return e1
	}
	return nil
}

func utimensat(dirfd int, path string, times *[2]Timespec, flags int) error {
	
	return ENOSYS
}





func Kill(pid int, signum syscall.Signal) (err error) { return kill(pid, int(signum), 1) }








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




































































































































































































































































































