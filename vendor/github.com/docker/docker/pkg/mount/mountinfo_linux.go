package mount 

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func parseInfoFile(r io.Reader, filter FilterFunc) ([]*Info, error) {
	s := bufio.NewScanner(r)
	out := []*Info{}
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}
		

		text := s.Text()
		fields := strings.Split(text, " ")
		numFields := len(fields)
		if numFields < 10 {
			
			return nil, fmt.Errorf("Parsing '%s' failed: not enough fields (%d)", text, numFields)
		}

		p := &Info{}
		
		p.ID, _ = strconv.Atoi(fields[0])
		p.Parent, _ = strconv.Atoi(fields[1])
		mm := strings.Split(fields[2], ":")
		if len(mm) != 2 {
			return nil, fmt.Errorf("Parsing '%s' failed: unexpected minor:major pair %s", text, mm)
		}
		p.Major, _ = strconv.Atoi(mm[0])
		p.Minor, _ = strconv.Atoi(mm[1])

		p.Root = fields[3]
		p.Mountpoint = fields[4]
		p.Opts = fields[5]

		var skip, stop bool
		if filter != nil {
			
			skip, stop = filter(p)
			if skip {
				continue
			}
		}

		
		i := 6
		for ; i < numFields && fields[i] != "-"; i++ {
			switch i {
			case 6:
				p.Optional = fields[6]
			default:
				
				break
			}
		}
		if i == numFields {
			return nil, fmt.Errorf("Parsing '%s' failed: missing separator ('-')", text)
		}

		
		if i+4 > numFields {
			return nil, fmt.Errorf("Parsing '%s' failed: not enough fields after a separator", text)
		}
		
		
		
		
		
		p.Fstype = fields[i+1]
		p.Source = fields[i+2]
		p.VfsOpts = fields[i+3]

		out = append(out, p)
		if stop {
			break
		}
	}
	return out, nil
}



func parseMountTable(filter FilterFunc) ([]*Info, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseInfoFile(f, filter)
}




func PidMountInfo(pid int) ([]*Info, error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/mountinfo", pid))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseInfoFile(f, nil)
}
