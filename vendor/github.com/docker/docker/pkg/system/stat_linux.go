package system 

import "syscall"


func fromStatT(s *syscall.Stat_t) (*StatT, error) {
	return &StatT{size: s.Size,
		mode: s.Mode,
		uid:  s.Uid,
		gid:  s.Gid,
		rdev: s.Rdev,
		mtim: s.Mtim}, nil
}



func FromStatT(s *syscall.Stat_t) (*StatT, error) {
	return fromStatT(s)
}
