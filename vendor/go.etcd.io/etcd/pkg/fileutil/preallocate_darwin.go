















package fileutil

import (
	"os"
	"syscall"
	"unsafe"
)

func preallocExtend(f *os.File, sizeInBytes int64) error {
	if err := preallocFixed(f, sizeInBytes); err != nil {
		return err
	}
	return preallocExtendTrunc(f, sizeInBytes)
}

func preallocFixed(f *os.File, sizeInBytes int64) error {
	
	
	fstore := &syscall.Fstore_t{
		Flags:   syscall.F_ALLOCATEALL,
		Posmode: syscall.F_PEOFPOSMODE,
		Length:  sizeInBytes}
	p := unsafe.Pointer(fstore)
	_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), uintptr(syscall.F_PREALLOCATE), uintptr(p))
	if errno == 0 || errno == syscall.ENOTSUP {
		return nil
	}

	
	if errno == syscall.EINVAL {
		
		
		var stat syscall.Stat_t
		syscall.Fstat(int(f.Fd()), &stat)

		
		
		
		var statfs syscall.Statfs_t
		syscall.Fstatfs(int(f.Fd()), &statfs)
		blockSize := int64(statfs.Bsize)

		if stat.Blocks*blockSize >= sizeInBytes {
			
			return nil
		}
	}
	return errno
}
