












package procfs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)


type CPUStat struct {
	User      float64
	Nice      float64
	System    float64
	Idle      float64
	Iowait    float64
	IRQ       float64
	SoftIRQ   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}




type SoftIRQStat struct {
	Hi          uint64
	Timer       uint64
	NetTx       uint64
	NetRx       uint64
	Block       uint64
	BlockIoPoll uint64
	Tasklet     uint64
	Sched       uint64
	Hrtimer     uint64
	Rcu         uint64
}


type Stat struct {
	
	BootTime uint64
	
	CPUTotal CPUStat
	
	CPU []CPUStat
	
	IRQTotal uint64
	
	IRQ []uint64
	
	ContextSwitches uint64
	
	ProcessCreated uint64
	
	ProcessesRunning uint64
	
	ProcessesBlocked uint64
	
	SoftIRQTotal uint64
	
	SoftIRQ SoftIRQStat
}


func NewStat() (Stat, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Stat{}, err
	}

	return fs.NewStat()
}


func parseCPUStat(line string) (CPUStat, int64, error) {
	cpuStat := CPUStat{}
	var cpu string

	count, err := fmt.Sscanf(line, "%s %f %f %f %f %f %f %f %f %f %f",
		&cpu,
		&cpuStat.User, &cpuStat.Nice, &cpuStat.System, &cpuStat.Idle,
		&cpuStat.Iowait, &cpuStat.IRQ, &cpuStat.SoftIRQ, &cpuStat.Steal,
		&cpuStat.Guest, &cpuStat.GuestNice)

	if err != nil && err != io.EOF {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %s (cpu): %s", line, err)
	}
	if count == 0 {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %s (cpu): 0 elements parsed", line)
	}

	cpuStat.User /= userHZ
	cpuStat.Nice /= userHZ
	cpuStat.System /= userHZ
	cpuStat.Idle /= userHZ
	cpuStat.Iowait /= userHZ
	cpuStat.IRQ /= userHZ
	cpuStat.SoftIRQ /= userHZ
	cpuStat.Steal /= userHZ
	cpuStat.Guest /= userHZ
	cpuStat.GuestNice /= userHZ

	if cpu == "cpu" {
		return cpuStat, -1, nil
	}

	cpuID, err := strconv.ParseInt(cpu[3:], 10, 64)
	if err != nil {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %s (cpu/cpuid): %s", line, err)
	}

	return cpuStat, cpuID, nil
}


func parseSoftIRQStat(line string) (SoftIRQStat, uint64, error) {
	softIRQStat := SoftIRQStat{}
	var total uint64
	var prefix string

	_, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d %d",
		&prefix, &total,
		&softIRQStat.Hi, &softIRQStat.Timer, &softIRQStat.NetTx, &softIRQStat.NetRx,
		&softIRQStat.Block, &softIRQStat.BlockIoPoll,
		&softIRQStat.Tasklet, &softIRQStat.Sched,
		&softIRQStat.Hrtimer, &softIRQStat.Rcu)

	if err != nil {
		return SoftIRQStat{}, 0, fmt.Errorf("couldn't parse %s (softirq): %s", line, err)
	}

	return softIRQStat, total, nil
}


func (fs FS) NewStat() (Stat, error) {
	

	f, err := os.Open(fs.Path("stat"))
	if err != nil {
		return Stat{}, err
	}
	defer f.Close()

	stat := Stat{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(scanner.Text())
		
		if len(parts) < 2 {
			continue
		}
		switch {
		case parts[0] == "btime":
			if stat.BootTime, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (btime): %s", parts[1], err)
			}
		case parts[0] == "intr":
			if stat.IRQTotal, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (intr): %s", parts[1], err)
			}
			numberedIRQs := parts[2:]
			stat.IRQ = make([]uint64, len(numberedIRQs))
			for i, count := range numberedIRQs {
				if stat.IRQ[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Stat{}, fmt.Errorf("couldn't parse %s (intr%d): %s", count, i, err)
				}
			}
		case parts[0] == "ctxt":
			if stat.ContextSwitches, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (ctxt): %s", parts[1], err)
			}
		case parts[0] == "processes":
			if stat.ProcessCreated, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (processes): %s", parts[1], err)
			}
		case parts[0] == "procs_running":
			if stat.ProcessesRunning, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (procs_running): %s", parts[1], err)
			}
		case parts[0] == "procs_blocked":
			if stat.ProcessesBlocked, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (procs_blocked): %s", parts[1], err)
			}
		case parts[0] == "softirq":
			softIRQStats, total, err := parseSoftIRQStat(line)
			if err != nil {
				return Stat{}, err
			}
			stat.SoftIRQTotal = total
			stat.SoftIRQ = softIRQStats
		case strings.HasPrefix(parts[0], "cpu"):
			cpuStat, cpuID, err := parseCPUStat(line)
			if err != nil {
				return Stat{}, err
			}
			if cpuID == -1 {
				stat.CPUTotal = cpuStat
			} else {
				for int64(len(stat.CPU)) <= cpuID {
					stat.CPU = append(stat.CPU, CPUStat{})
				}
				stat.CPU[cpuID] = cpuStat
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Stat{}, fmt.Errorf("couldn't parse %s: %s", f.Name(), err)
	}

	return stat, nil
}
