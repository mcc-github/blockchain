

package winio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"syscall"
	"unicode/utf16"
)




const (
	BackupData = uint32(iota + 1)
	BackupEaData
	BackupSecurity
	BackupAlternateData
	BackupLink
	BackupPropertyData
	BackupObjectId
	BackupReparseData
	BackupSparseBlock
	BackupTxfsData
)

const (
	StreamSparseAttributes = uint32(8)
)

const (
	WRITE_DAC              = 0x40000
	WRITE_OWNER            = 0x80000
	ACCESS_SYSTEM_SECURITY = 0x1000000
)


type BackupHeader struct {
	Id         uint32 
	Attributes uint32 
	Size       int64  
	Name       string 
	Offset     int64  
}

type win32StreamId struct {
	StreamId   uint32
	Attributes uint32
	Size       uint64
	NameSize   uint32
}



type BackupStreamReader struct {
	r         io.Reader
	bytesLeft int64
}


func NewBackupStreamReader(r io.Reader) *BackupStreamReader {
	return &BackupStreamReader{r, 0}
}



func (r *BackupStreamReader) Next() (*BackupHeader, error) {
	if r.bytesLeft > 0 {
		if s, ok := r.r.(io.Seeker); ok {
			
			
			if _, err := s.Seek(0, io.SeekCurrent); err == nil {
				if _, err = s.Seek(r.bytesLeft, io.SeekCurrent); err != nil {
					return nil, err
				}
				r.bytesLeft = 0
			}
		}
		if _, err := io.Copy(ioutil.Discard, r); err != nil {
			return nil, err
		}
	}
	var wsi win32StreamId
	if err := binary.Read(r.r, binary.LittleEndian, &wsi); err != nil {
		return nil, err
	}
	hdr := &BackupHeader{
		Id:         wsi.StreamId,
		Attributes: wsi.Attributes,
		Size:       int64(wsi.Size),
	}
	if wsi.NameSize != 0 {
		name := make([]uint16, int(wsi.NameSize/2))
		if err := binary.Read(r.r, binary.LittleEndian, name); err != nil {
			return nil, err
		}
		hdr.Name = syscall.UTF16ToString(name)
	}
	if wsi.StreamId == BackupSparseBlock {
		if err := binary.Read(r.r, binary.LittleEndian, &hdr.Offset); err != nil {
			return nil, err
		}
		hdr.Size -= 8
	}
	r.bytesLeft = hdr.Size
	return hdr, nil
}


func (r *BackupStreamReader) Read(b []byte) (int, error) {
	if r.bytesLeft == 0 {
		return 0, io.EOF
	}
	if int64(len(b)) > r.bytesLeft {
		b = b[:r.bytesLeft]
	}
	n, err := r.r.Read(b)
	r.bytesLeft -= int64(n)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	} else if r.bytesLeft == 0 && err == nil {
		err = io.EOF
	}
	return n, err
}


type BackupStreamWriter struct {
	w         io.Writer
	bytesLeft int64
}


func NewBackupStreamWriter(w io.Writer) *BackupStreamWriter {
	return &BackupStreamWriter{w, 0}
}


func (w *BackupStreamWriter) WriteHeader(hdr *BackupHeader) error {
	if w.bytesLeft != 0 {
		return fmt.Errorf("missing %d bytes", w.bytesLeft)
	}
	name := utf16.Encode([]rune(hdr.Name))
	wsi := win32StreamId{
		StreamId:   hdr.Id,
		Attributes: hdr.Attributes,
		Size:       uint64(hdr.Size),
		NameSize:   uint32(len(name) * 2),
	}
	if hdr.Id == BackupSparseBlock {
		
		wsi.Size += 8
	}
	if err := binary.Write(w.w, binary.LittleEndian, &wsi); err != nil {
		return err
	}
	if len(name) != 0 {
		if err := binary.Write(w.w, binary.LittleEndian, name); err != nil {
			return err
		}
	}
	if hdr.Id == BackupSparseBlock {
		if err := binary.Write(w.w, binary.LittleEndian, hdr.Offset); err != nil {
			return err
		}
	}
	w.bytesLeft = hdr.Size
	return nil
}


func (w *BackupStreamWriter) Write(b []byte) (int, error) {
	if w.bytesLeft < int64(len(b)) {
		return 0, fmt.Errorf("too many bytes by %d", int64(len(b))-w.bytesLeft)
	}
	n, err := w.w.Write(b)
	w.bytesLeft -= int64(n)
	return n, err
}


type BackupFileReader struct {
	f               *os.File
	includeSecurity bool
	ctx             uintptr
}



func NewBackupFileReader(f *os.File, includeSecurity bool) *BackupFileReader {
	r := &BackupFileReader{f, includeSecurity, 0}
	return r
}


func (r *BackupFileReader) Read(b []byte) (int, error) {
	var bytesRead uint32
	err := backupRead(syscall.Handle(r.f.Fd()), b, &bytesRead, false, r.includeSecurity, &r.ctx)
	if err != nil {
		return 0, &os.PathError{"BackupRead", r.f.Name(), err}
	}
	runtime.KeepAlive(r.f)
	if bytesRead == 0 {
		return 0, io.EOF
	}
	return int(bytesRead), nil
}



func (r *BackupFileReader) Close() error {
	if r.ctx != 0 {
		backupRead(syscall.Handle(r.f.Fd()), nil, nil, true, false, &r.ctx)
		runtime.KeepAlive(r.f)
		r.ctx = 0
	}
	return nil
}


type BackupFileWriter struct {
	f               *os.File
	includeSecurity bool
	ctx             uintptr
}



func NewBackupFileWriter(f *os.File, includeSecurity bool) *BackupFileWriter {
	w := &BackupFileWriter{f, includeSecurity, 0}
	return w
}


func (w *BackupFileWriter) Write(b []byte) (int, error) {
	var bytesWritten uint32
	err := backupWrite(syscall.Handle(w.f.Fd()), b, &bytesWritten, false, w.includeSecurity, &w.ctx)
	if err != nil {
		return 0, &os.PathError{"BackupWrite", w.f.Name(), err}
	}
	runtime.KeepAlive(w.f)
	if int(bytesWritten) != len(b) {
		return int(bytesWritten), errors.New("not all bytes could be written")
	}
	return len(b), nil
}



func (w *BackupFileWriter) Close() error {
	if w.ctx != 0 {
		backupWrite(syscall.Handle(w.f.Fd()), nil, nil, true, false, &w.ctx)
		runtime.KeepAlive(w.f)
		w.ctx = 0
	}
	return nil
}





func OpenForBackup(path string, access uint32, share uint32, createmode uint32) (*os.File, error) {
	winPath, err := syscall.UTF16FromString(path)
	if err != nil {
		return nil, err
	}
	h, err := syscall.CreateFile(&winPath[0], access, share, nil, createmode, syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OPEN_REPARSE_POINT, 0)
	if err != nil {
		err = &os.PathError{Op: "open", Path: path, Err: err}
		return nil, err
	}
	return os.NewFile(uintptr(h), path), nil
}
