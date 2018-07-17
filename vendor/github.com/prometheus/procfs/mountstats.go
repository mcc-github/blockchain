












package procfs








import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)


const (
	deviceEntryLen = 8

	fieldBytesLen  = 8
	fieldEventsLen = 27

	statVersion10 = "1.0"
	statVersion11 = "1.1"

	fieldTransport10Len = 10
	fieldTransport11Len = 13
)


type Mount struct {
	
	Device string
	
	Mount string
	
	Type string
	
	
	Stats MountStats
}



type MountStats interface {
	mountStats()
}


type MountStatsNFS struct {
	
	StatVersion string
	
	Age time.Duration
	
	Bytes NFSBytesStats
	
	Events NFSEventsStats
	
	Operations []NFSOperationStats
	
	Transport NFSTransportStats
}


func (m MountStatsNFS) mountStats() {}



type NFSBytesStats struct {
	
	Read uint64
	
	Write uint64
	
	DirectRead uint64
	
	DirectWrite uint64
	
	ReadTotal uint64
	
	WriteTotal uint64
	
	ReadPages uint64
	
	WritePages uint64
}


type NFSEventsStats struct {
	
	InodeRevalidate uint64
	
	DnodeRevalidate uint64
	
	DataInvalidate uint64
	
	AttributeInvalidate uint64
	
	VFSOpen uint64
	
	VFSLookup uint64
	
	VFSAccess uint64
	
	VFSUpdatePage uint64
	
	VFSReadPage uint64
	
	VFSReadPages uint64
	
	VFSWritePage uint64
	
	VFSWritePages uint64
	
	VFSGetdents uint64
	
	VFSSetattr uint64
	
	VFSFlush uint64
	
	VFSFsync uint64
	
	VFSLock uint64
	
	VFSFileRelease uint64
	
	CongestionWait uint64
	
	Truncation uint64
	
	WriteExtension uint64
	
	SillyRename uint64
	
	ShortRead uint64
	
	ShortWrite uint64
	
	
	JukeboxDelay uint64
	
	PNFSRead uint64
	
	PNFSWrite uint64
}


type NFSOperationStats struct {
	
	Operation string
	
	Requests uint64
	
	Transmissions uint64
	
	MajorTimeouts uint64
	
	BytesSent uint64
	
	BytesReceived uint64
	
	CumulativeQueueTime time.Duration
	
	CumulativeTotalResponseTime time.Duration
	
	CumulativeTotalRequestTime time.Duration
}



type NFSTransportStats struct {
	
	Port uint64
	
	
	Bind uint64
	
	Connect uint64
	
	
	ConnectIdleTime uint64
	
	IdleTime time.Duration
	
	Sends uint64
	
	Receives uint64
	
	
	BadTransactionIDs uint64
	
	
	CumulativeActiveRequests uint64
	
	
	CumulativeBacklog uint64

	

	
	MaximumRPCSlotsUsed uint64
	
	
	CumulativeSendingQueue uint64
	
	
	CumulativePendingQueue uint64
}




func parseMountStats(r io.Reader) ([]*Mount, error) {
	const (
		device            = "device"
		statVersionPrefix = "statvers="

		nfs3Type = "nfs"
		nfs4Type = "nfs4"
	)

	var mounts []*Mount

	s := bufio.NewScanner(r)
	for s.Scan() {
		
		ss := strings.Fields(string(s.Bytes()))
		if len(ss) == 0 || ss[0] != device {
			continue
		}

		m, err := parseMount(ss)
		if err != nil {
			return nil, err
		}

		
		if len(ss) > deviceEntryLen {
			
			if m.Type != nfs3Type && m.Type != nfs4Type {
				return nil, fmt.Errorf("cannot parse MountStats for fstype %q", m.Type)
			}

			statVersion := strings.TrimPrefix(ss[8], statVersionPrefix)

			stats, err := parseMountStatsNFS(s, statVersion)
			if err != nil {
				return nil, err
			}

			m.Stats = stats
		}

		mounts = append(mounts, m)
	}

	return mounts, s.Err()
}



func parseMount(ss []string) (*Mount, error) {
	if len(ss) < deviceEntryLen {
		return nil, fmt.Errorf("invalid device entry: %v", ss)
	}

	
	
	format := []struct {
		i int
		s string
	}{
		{i: 0, s: "device"},
		{i: 2, s: "mounted"},
		{i: 3, s: "on"},
		{i: 5, s: "with"},
		{i: 6, s: "fstype"},
	}

	for _, f := range format {
		if ss[f.i] != f.s {
			return nil, fmt.Errorf("invalid device entry: %v", ss)
		}
	}

	return &Mount{
		Device: ss[1],
		Mount:  ss[4],
		Type:   ss[7],
	}, nil
}



func parseMountStatsNFS(s *bufio.Scanner, statVersion string) (*MountStatsNFS, error) {
	
	const (
		fieldAge        = "age:"
		fieldBytes      = "bytes:"
		fieldEvents     = "events:"
		fieldPerOpStats = "per-op"
		fieldTransport  = "xprt:"
	)

	stats := &MountStatsNFS{
		StatVersion: statVersion,
	}

	for s.Scan() {
		ss := strings.Fields(string(s.Bytes()))
		if len(ss) == 0 {
			break
		}
		if len(ss) < 2 {
			return nil, fmt.Errorf("not enough information for NFS stats: %v", ss)
		}

		switch ss[0] {
		case fieldAge:
			
			d, err := time.ParseDuration(ss[1] + "s")
			if err != nil {
				return nil, err
			}

			stats.Age = d
		case fieldBytes:
			bstats, err := parseNFSBytesStats(ss[1:])
			if err != nil {
				return nil, err
			}

			stats.Bytes = *bstats
		case fieldEvents:
			estats, err := parseNFSEventsStats(ss[1:])
			if err != nil {
				return nil, err
			}

			stats.Events = *estats
		case fieldTransport:
			if len(ss) < 3 {
				return nil, fmt.Errorf("not enough information for NFS transport stats: %v", ss)
			}

			tstats, err := parseNFSTransportStats(ss[2:], statVersion)
			if err != nil {
				return nil, err
			}

			stats.Transport = *tstats
		}

		
		
		
		
		if ss[0] == fieldPerOpStats {
			break
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	
	perOpStats, err := parseNFSOperationStats(s)
	if err != nil {
		return nil, err
	}

	stats.Operations = perOpStats

	return stats, nil
}



func parseNFSBytesStats(ss []string) (*NFSBytesStats, error) {
	if len(ss) != fieldBytesLen {
		return nil, fmt.Errorf("invalid NFS bytes stats: %v", ss)
	}

	ns := make([]uint64, 0, fieldBytesLen)
	for _, s := range ss {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		ns = append(ns, n)
	}

	return &NFSBytesStats{
		Read:        ns[0],
		Write:       ns[1],
		DirectRead:  ns[2],
		DirectWrite: ns[3],
		ReadTotal:   ns[4],
		WriteTotal:  ns[5],
		ReadPages:   ns[6],
		WritePages:  ns[7],
	}, nil
}



func parseNFSEventsStats(ss []string) (*NFSEventsStats, error) {
	if len(ss) != fieldEventsLen {
		return nil, fmt.Errorf("invalid NFS events stats: %v", ss)
	}

	ns := make([]uint64, 0, fieldEventsLen)
	for _, s := range ss {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		ns = append(ns, n)
	}

	return &NFSEventsStats{
		InodeRevalidate:     ns[0],
		DnodeRevalidate:     ns[1],
		DataInvalidate:      ns[2],
		AttributeInvalidate: ns[3],
		VFSOpen:             ns[4],
		VFSLookup:           ns[5],
		VFSAccess:           ns[6],
		VFSUpdatePage:       ns[7],
		VFSReadPage:         ns[8],
		VFSReadPages:        ns[9],
		VFSWritePage:        ns[10],
		VFSWritePages:       ns[11],
		VFSGetdents:         ns[12],
		VFSSetattr:          ns[13],
		VFSFlush:            ns[14],
		VFSFsync:            ns[15],
		VFSLock:             ns[16],
		VFSFileRelease:      ns[17],
		CongestionWait:      ns[18],
		Truncation:          ns[19],
		WriteExtension:      ns[20],
		SillyRename:         ns[21],
		ShortRead:           ns[22],
		ShortWrite:          ns[23],
		JukeboxDelay:        ns[24],
		PNFSRead:            ns[25],
		PNFSWrite:           ns[26],
	}, nil
}




func parseNFSOperationStats(s *bufio.Scanner) ([]NFSOperationStats, error) {
	const (
		
		numFields = 9
	)

	var ops []NFSOperationStats

	for s.Scan() {
		ss := strings.Fields(string(s.Bytes()))
		if len(ss) == 0 {
			
			
			break
		}

		if len(ss) != numFields {
			return nil, fmt.Errorf("invalid NFS per-operations stats: %v", ss)
		}

		
		ns := make([]uint64, 0, numFields-1)
		for _, st := range ss[1:] {
			n, err := strconv.ParseUint(st, 10, 64)
			if err != nil {
				return nil, err
			}

			ns = append(ns, n)
		}

		ops = append(ops, NFSOperationStats{
			Operation:                   strings.TrimSuffix(ss[0], ":"),
			Requests:                    ns[0],
			Transmissions:               ns[1],
			MajorTimeouts:               ns[2],
			BytesSent:                   ns[3],
			BytesReceived:               ns[4],
			CumulativeQueueTime:         time.Duration(ns[5]) * time.Millisecond,
			CumulativeTotalResponseTime: time.Duration(ns[6]) * time.Millisecond,
			CumulativeTotalRequestTime:  time.Duration(ns[7]) * time.Millisecond,
		})
	}

	return ops, s.Err()
}



func parseNFSTransportStats(ss []string, statVersion string) (*NFSTransportStats, error) {
	switch statVersion {
	case statVersion10:
		if len(ss) != fieldTransport10Len {
			return nil, fmt.Errorf("invalid NFS transport stats 1.0 statement: %v", ss)
		}
	case statVersion11:
		if len(ss) != fieldTransport11Len {
			return nil, fmt.Errorf("invalid NFS transport stats 1.1 statement: %v", ss)
		}
	default:
		return nil, fmt.Errorf("unrecognized NFS transport stats version: %q", statVersion)
	}

	
	
	
	
	
	
	ns := make([]uint64, fieldTransport11Len)
	for i, s := range ss {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		ns[i] = n
	}

	return &NFSTransportStats{
		Port:                     ns[0],
		Bind:                     ns[1],
		Connect:                  ns[2],
		ConnectIdleTime:          ns[3],
		IdleTime:                 time.Duration(ns[4]) * time.Second,
		Sends:                    ns[5],
		Receives:                 ns[6],
		BadTransactionIDs:        ns[7],
		CumulativeActiveRequests: ns[8],
		CumulativeBacklog:        ns[9],
		MaximumRPCSlotsUsed:      ns[10],
		CumulativeSendingQueue:   ns[11],
		CumulativePendingQueue:   ns[12],
	}, nil
}
