





package unix

import (
	"strings"
	"unsafe"
)



func xattrnamespace(fullattr string) (ns int, attr string, err error) {
	s := strings.IndexByte(fullattr, '.')
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
