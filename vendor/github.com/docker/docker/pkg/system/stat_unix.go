

package system 

import (
	"syscall"
)



type StatT struct {
	mode uint32
	uid  uint32
	gid  uint32
	rdev uint64
	size int64
	mtim syscall.Timespec
}


func (s StatT) Mode() uint32 {
	return s.mode
}


func (s StatT) UID() uint32 {
	return s.uid
}


func (s StatT) GID() uint32 {
	return s.gid
}


func (s StatT) Rdev() uint64 {
	return s.rdev
}


func (s StatT) Size() int64 {
	return s.size
}


func (s StatT) Mtim() syscall.Timespec {
	return s.mtim
}


func (s StatT) IsDir() bool {
	return s.mode&syscall.S_IFDIR != 0
}





func Stat(path string) (*StatT, error) {
	s := &syscall.Stat_t{}
	if err := syscall.Stat(path, s); err != nil {
		return nil, err
	}
	return fromStatT(s)
}
