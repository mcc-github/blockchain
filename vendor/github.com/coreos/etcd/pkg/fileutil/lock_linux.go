















package fileutil

import (
	"io"
	"os"
	"syscall"
)







const (
	F_OFD_GETLK  = 37
	F_OFD_SETLK  = 37
	F_OFD_SETLKW = 38
)

var (
	wrlck = syscall.Flock_t{
		Type:   syscall.F_WRLCK,
		Whence: int16(io.SeekStart),
		Start:  0,
		Len:    0,
	}

	linuxTryLockFile = flockTryLockFile
	linuxLockFile    = flockLockFile
)

func init() {
	
	getlk := syscall.Flock_t{Type: syscall.F_RDLCK}
	if err := syscall.FcntlFlock(0, F_OFD_GETLK, &getlk); err == nil {
		linuxTryLockFile = ofdTryLockFile
		linuxLockFile = ofdLockFile
	}
}

func TryLockFile(path string, flag int, perm os.FileMode) (*LockedFile, error) {
	return linuxTryLockFile(path, flag, perm)
}

func ofdTryLockFile(path string, flag int, perm os.FileMode) (*LockedFile, error) {
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	flock := wrlck
	if err = syscall.FcntlFlock(f.Fd(), F_OFD_SETLK, &flock); err != nil {
		f.Close()
		if err == syscall.EWOULDBLOCK {
			err = ErrLocked
		}
		return nil, err
	}
	return &LockedFile{f}, nil
}

func LockFile(path string, flag int, perm os.FileMode) (*LockedFile, error) {
	return linuxLockFile(path, flag, perm)
}

func ofdLockFile(path string, flag int, perm os.FileMode) (*LockedFile, error) {
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	flock := wrlck
	err = syscall.FcntlFlock(f.Fd(), F_OFD_SETLKW, &flock)

	if err != nil {
		f.Close()
		return nil, err
	}
	return &LockedFile{f}, err
}