






package fastwalk

import (
	"bytes"
	"syscall"
	"unsafe"
)

func direntNamlen(dirent *syscall.Dirent) uint64 {
	const fixedHdr = uint16(unsafe.Offsetof(syscall.Dirent{}.Name))
	nameBuf := (*[unsafe.Sizeof(dirent.Name)]byte)(unsafe.Pointer(&dirent.Name[0]))
	nameLen := bytes.IndexByte(nameBuf[:dirent.Reclen-fixedHdr], 0)
	if nameLen < 0 {
		panic("failed to find terminating 0 byte in dirent")
	}
	return uint64(nameLen)
}
