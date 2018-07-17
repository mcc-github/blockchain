package system 

import (
	"os"
	"time"
)



type StatT struct {
	mode os.FileMode
	size int64
	mtim time.Time
}


func (s StatT) Size() int64 {
	return s.size
}


func (s StatT) Mode() os.FileMode {
	return os.FileMode(s.mode)
}


func (s StatT) Mtim() time.Time {
	return time.Time(s.mtim)
}





func Stat(path string) (*StatT, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return fromStatT(&fi)
}


func fromStatT(fi *os.FileInfo) (*StatT, error) {
	return &StatT{
		size: (*fi).Size(),
		mode: (*fi).Mode(),
		mtim: (*fi).ModTime()}, nil
}
