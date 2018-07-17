












package procfs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)
















const userHZ = 100



type ProcStat struct {
	
	PID int
	
	Comm string
	
	State string
	
	PPID int
	
	PGRP int
	
	Session int
	
	TTY int
	
	
	TPGID int
	
	Flags uint
	
	
	MinFlt uint
	
	
	CMinFlt uint
	
	
	MajFlt uint
	
	
	CMajFlt uint
	
	
	UTime uint
	
	
	STime uint
	
	
	CUTime uint
	
	
	CSTime uint
	
	
	Priority int
	
	
	Nice int
	
	NumThreads int
	
	
	Starttime uint64
	
	VSize int
	
	RSS int

	fs FS
}


func (p Proc) NewStat() (ProcStat, error) {
	f, err := os.Open(p.path("stat"))
	if err != nil {
		return ProcStat{}, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return ProcStat{}, err
	}

	var (
		ignore int

		s = ProcStat{PID: p.PID, fs: p.fs}
		l = bytes.Index(data, []byte("("))
		r = bytes.LastIndex(data, []byte(")"))
	)

	if l < 0 || r < 0 {
		return ProcStat{}, fmt.Errorf(
			"unexpected format, couldn't extract comm: %s",
			data,
		)
	}

	s.Comm = string(data[l+1 : r])
	_, err = fmt.Fscan(
		bytes.NewBuffer(data[r+2:]),
		&s.State,
		&s.PPID,
		&s.PGRP,
		&s.Session,
		&s.TTY,
		&s.TPGID,
		&s.Flags,
		&s.MinFlt,
		&s.CMinFlt,
		&s.MajFlt,
		&s.CMajFlt,
		&s.UTime,
		&s.STime,
		&s.CUTime,
		&s.CSTime,
		&s.Priority,
		&s.Nice,
		&s.NumThreads,
		&ignore,
		&s.Starttime,
		&s.VSize,
		&s.RSS,
	)
	if err != nil {
		return ProcStat{}, err
	}

	return s, nil
}


func (s ProcStat) VirtualMemory() int {
	return s.VSize
}


func (s ProcStat) ResidentMemory() int {
	return s.RSS * os.Getpagesize()
}


func (s ProcStat) StartTime() (float64, error) {
	stat, err := s.fs.NewStat()
	if err != nil {
		return 0, err
	}
	return float64(stat.BootTime) + (float64(s.Starttime) / userHZ), nil
}


func (s ProcStat) CPUTime() float64 {
	return float64(s.UTime+s.STime) / userHZ
}
