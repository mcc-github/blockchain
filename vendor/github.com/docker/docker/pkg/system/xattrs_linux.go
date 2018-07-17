package system 

import "golang.org/x/sys/unix"




func Lgetxattr(path string, attr string) ([]byte, error) {
	dest := make([]byte, 128)
	sz, errno := unix.Lgetxattr(path, attr, dest)
	if errno == unix.ENODATA {
		return nil, nil
	}
	if errno == unix.ERANGE {
		dest = make([]byte, sz)
		sz, errno = unix.Lgetxattr(path, attr, dest)
	}
	if errno != nil {
		return nil, errno
	}

	return dest[:sz], nil
}



func Lsetxattr(path string, attr string, data []byte, flags int) error {
	return unix.Lsetxattr(path, attr, data, flags)
}
