












package procfs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


type XfrmStat struct {
	
	XfrmInError int
	
	XfrmInBufferError int
	
	XfrmInHdrError int
	
	
	XfrmInNoStates int
	
	
	XfrmInStateProtoError int
	
	XfrmInStateModeError int
	
	
	XfrmInStateSeqError int
	
	XfrmInStateExpired int
	
	
	XfrmInStateMismatch int
	
	XfrmInStateInvalid int
	
	
	XfrmInTmplMismatch int
	
	
	XfrmInNoPols int
	
	XfrmInPolBlock int
	
	XfrmInPolError int
	
	XfrmOutError int
	
	XfrmOutBundleGenError int
	
	XfrmOutBundleCheckError int
	
	XfrmOutNoStates int
	
	XfrmOutStateProtoError int
	
	XfrmOutStateModeError int
	
	
	XfrmOutStateSeqError int
	
	XfrmOutStateExpired int
	
	XfrmOutPolBlock int
	
	XfrmOutPolDead int
	
	XfrmOutPolError     int
	XfrmFwdHdrError     int
	XfrmOutStateInvalid int
	XfrmAcquireError    int
}


func NewXfrmStat() (XfrmStat, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return XfrmStat{}, err
	}

	return fs.NewXfrmStat()
}


func (fs FS) NewXfrmStat() (XfrmStat, error) {
	file, err := os.Open(fs.Path("net/xfrm_stat"))
	if err != nil {
		return XfrmStat{}, err
	}
	defer file.Close()

	var (
		x = XfrmStat{}
		s = bufio.NewScanner(file)
	)

	for s.Scan() {
		fields := strings.Fields(s.Text())

		if len(fields) != 2 {
			return XfrmStat{}, fmt.Errorf(
				"couldn't parse %s line %s", file.Name(), s.Text())
		}

		name := fields[0]
		value, err := strconv.Atoi(fields[1])
		if err != nil {
			return XfrmStat{}, err
		}

		switch name {
		case "XfrmInError":
			x.XfrmInError = value
		case "XfrmInBufferError":
			x.XfrmInBufferError = value
		case "XfrmInHdrError":
			x.XfrmInHdrError = value
		case "XfrmInNoStates":
			x.XfrmInNoStates = value
		case "XfrmInStateProtoError":
			x.XfrmInStateProtoError = value
		case "XfrmInStateModeError":
			x.XfrmInStateModeError = value
		case "XfrmInStateSeqError":
			x.XfrmInStateSeqError = value
		case "XfrmInStateExpired":
			x.XfrmInStateExpired = value
		case "XfrmInStateInvalid":
			x.XfrmInStateInvalid = value
		case "XfrmInTmplMismatch":
			x.XfrmInTmplMismatch = value
		case "XfrmInNoPols":
			x.XfrmInNoPols = value
		case "XfrmInPolBlock":
			x.XfrmInPolBlock = value
		case "XfrmInPolError":
			x.XfrmInPolError = value
		case "XfrmOutError":
			x.XfrmOutError = value
		case "XfrmInStateMismatch":
			x.XfrmInStateMismatch = value
		case "XfrmOutBundleGenError":
			x.XfrmOutBundleGenError = value
		case "XfrmOutBundleCheckError":
			x.XfrmOutBundleCheckError = value
		case "XfrmOutNoStates":
			x.XfrmOutNoStates = value
		case "XfrmOutStateProtoError":
			x.XfrmOutStateProtoError = value
		case "XfrmOutStateModeError":
			x.XfrmOutStateModeError = value
		case "XfrmOutStateSeqError":
			x.XfrmOutStateSeqError = value
		case "XfrmOutStateExpired":
			x.XfrmOutStateExpired = value
		case "XfrmOutPolBlock":
			x.XfrmOutPolBlock = value
		case "XfrmOutPolDead":
			x.XfrmOutPolDead = value
		case "XfrmOutPolError":
			x.XfrmOutPolError = value
		case "XfrmFwdHdrError":
			x.XfrmFwdHdrError = value
		case "XfrmOutStateInvalid":
			x.XfrmOutStateInvalid = value
		case "XfrmAcquireError":
			x.XfrmAcquireError = value
		}

	}

	return x, s.Err()
}
