package units

import (
	"fmt"
	"strconv"
	"strings"
)


type Ulimit struct {
	Name string
	Hard int64
	Soft int64
}


type Rlimit struct {
	Type int    `json:"type,omitempty"`
	Hard uint64 `json:"hard,omitempty"`
	Soft uint64 `json:"soft,omitempty"`
}

const (
	
	
	
	
	rlimitAs         = 9
	rlimitCore       = 4
	rlimitCPU        = 0
	rlimitData       = 2
	rlimitFsize      = 1
	rlimitLocks      = 10
	rlimitMemlock    = 8
	rlimitMsgqueue   = 12
	rlimitNice       = 13
	rlimitNofile     = 7
	rlimitNproc      = 6
	rlimitRss        = 5
	rlimitRtprio     = 14
	rlimitRttime     = 15
	rlimitSigpending = 11
	rlimitStack      = 3
)

var ulimitNameMapping = map[string]int{
	
	"core":       rlimitCore,
	"cpu":        rlimitCPU,
	"data":       rlimitData,
	"fsize":      rlimitFsize,
	"locks":      rlimitLocks,
	"memlock":    rlimitMemlock,
	"msgqueue":   rlimitMsgqueue,
	"nice":       rlimitNice,
	"nofile":     rlimitNofile,
	"nproc":      rlimitNproc,
	"rss":        rlimitRss,
	"rtprio":     rlimitRtprio,
	"rttime":     rlimitRttime,
	"sigpending": rlimitSigpending,
	"stack":      rlimitStack,
}


func ParseUlimit(val string) (*Ulimit, error) {
	parts := strings.SplitN(val, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ulimit argument: %s", val)
	}

	if _, exists := ulimitNameMapping[parts[0]]; !exists {
		return nil, fmt.Errorf("invalid ulimit type: %s", parts[0])
	}

	var (
		soft int64
		hard = &soft 
		temp int64
		err  error
	)
	switch limitVals := strings.Split(parts[1], ":"); len(limitVals) {
	case 2:
		temp, err = strconv.ParseInt(limitVals[1], 10, 64)
		if err != nil {
			return nil, err
		}
		hard = &temp
		fallthrough
	case 1:
		soft, err = strconv.ParseInt(limitVals[0], 10, 64)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("too many limit value arguments - %s, can only have up to two, `soft[:hard]`", parts[1])
	}

	if soft > *hard {
		return nil, fmt.Errorf("ulimit soft limit must be less than or equal to hard limit: %d > %d", soft, *hard)
	}

	return &Ulimit{Name: parts[0], Soft: soft, Hard: *hard}, nil
}


func (u *Ulimit) GetRlimit() (*Rlimit, error) {
	t, exists := ulimitNameMapping[u.Name]
	if !exists {
		return nil, fmt.Errorf("invalid ulimit name %s", u.Name)
	}

	return &Rlimit{Type: t, Soft: uint64(u.Soft), Hard: uint64(u.Hard)}, nil
}

func (u *Ulimit) String() string {
	return fmt.Sprintf("%s=%d:%d", u.Name, u.Soft, u.Hard)
}
