

package winio

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)




const (
	fileBasicInfo = 0
	fileIDInfo    = 0x12
)


type FileBasicInfo struct {
	CreationTime, LastAccessTime, LastWriteTime, ChangeTime syscall.Filetime
	FileAttributes                                          uint32
	pad                                                     uint32 
}


func GetFileBasicInfo(f *os.File) (*FileBasicInfo, error) {
	bi := &FileBasicInfo{}
	if err := getFileInformationByHandleEx(syscall.Handle(f.Fd()), fileBasicInfo, (*byte)(unsafe.Pointer(bi)), uint32(unsafe.Sizeof(*bi))); err != nil {
		return nil, &os.PathError{Op: "GetFileInformationByHandleEx", Path: f.Name(), Err: err}
	}
	runtime.KeepAlive(f)
	return bi, nil
}


func SetFileBasicInfo(f *os.File, bi *FileBasicInfo) error {
	if err := setFileInformationByHandle(syscall.Handle(f.Fd()), fileBasicInfo, (*byte)(unsafe.Pointer(bi)), uint32(unsafe.Sizeof(*bi))); err != nil {
		return &os.PathError{Op: "SetFileInformationByHandle", Path: f.Name(), Err: err}
	}
	runtime.KeepAlive(f)
	return nil
}



type FileIDInfo struct {
	VolumeSerialNumber uint64
	FileID             [16]byte
}


func GetFileID(f *os.File) (*FileIDInfo, error) {
	fileID := &FileIDInfo{}
	if err := getFileInformationByHandleEx(syscall.Handle(f.Fd()), fileIDInfo, (*byte)(unsafe.Pointer(fileID)), uint32(unsafe.Sizeof(*fileID))); err != nil {
		return nil, &os.PathError{Op: "GetFileInformationByHandleEx", Path: f.Name(), Err: err}
	}
	runtime.KeepAlive(f)
	return fileID, nil
}
