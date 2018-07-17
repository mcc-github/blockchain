












package procfs

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)




type ProcLimits struct {
	
	CPUTime int64
	
	FileSize int64
	
	
	DataSize int64
	
	StackSize int64
	
	CoreFileSize int64
	
	ResidentSet int64
	
	
	Processes int64
	
	
	OpenFiles int64
	
	LockedMemory int64
	
	AddressSpace int64
	
	
	FileLocks int64
	
	
	PendingSignals int64
	
	
	MsqqueueSize int64
	
	NicePriority int64
	
	
	RealtimePriority int64
	
	
	
	RealtimeTimeout int64
}

const (
	limitsFields    = 3
	limitsUnlimited = "unlimited"
)

var (
	limitsDelimiter = regexp.MustCompile("  +")
)


func (p Proc) NewLimits() (ProcLimits, error) {
	f, err := os.Open(p.path("limits"))
	if err != nil {
		return ProcLimits{}, err
	}
	defer f.Close()

	var (
		l = ProcLimits{}
		s = bufio.NewScanner(f)
	)
	for s.Scan() {
		fields := limitsDelimiter.Split(s.Text(), limitsFields)
		if len(fields) != limitsFields {
			return ProcLimits{}, fmt.Errorf(
				"couldn't parse %s line %s", f.Name(), s.Text())
		}

		switch fields[0] {
		case "Max cpu time":
			l.CPUTime, err = parseInt(fields[1])
		case "Max file size":
			l.FileSize, err = parseInt(fields[1])
		case "Max data size":
			l.DataSize, err = parseInt(fields[1])
		case "Max stack size":
			l.StackSize, err = parseInt(fields[1])
		case "Max core file size":
			l.CoreFileSize, err = parseInt(fields[1])
		case "Max resident set":
			l.ResidentSet, err = parseInt(fields[1])
		case "Max processes":
			l.Processes, err = parseInt(fields[1])
		case "Max open files":
			l.OpenFiles, err = parseInt(fields[1])
		case "Max locked memory":
			l.LockedMemory, err = parseInt(fields[1])
		case "Max address space":
			l.AddressSpace, err = parseInt(fields[1])
		case "Max file locks":
			l.FileLocks, err = parseInt(fields[1])
		case "Max pending signals":
			l.PendingSignals, err = parseInt(fields[1])
		case "Max msgqueue size":
			l.MsqqueueSize, err = parseInt(fields[1])
		case "Max nice priority":
			l.NicePriority, err = parseInt(fields[1])
		case "Max realtime priority":
			l.RealtimePriority, err = parseInt(fields[1])
		case "Max realtime timeout":
			l.RealtimeTimeout, err = parseInt(fields[1])
		}
		if err != nil {
			return ProcLimits{}, err
		}
	}

	return l, s.Err()
}

func parseInt(s string) (int64, error) {
	if s == limitsUnlimited {
		return -1, nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("couldn't parse value %s: %s", s, err)
	}
	return i, nil
}
