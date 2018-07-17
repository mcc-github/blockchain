












package procfs

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	statuslineRE = regexp.MustCompile(`(\d+) blocks .*\[(\d+)/(\d+)\] \[[U_]+\]`)
	buildlineRE  = regexp.MustCompile(`\((\d+)/\d+\)`)
)


type MDStat struct {
	
	Name string
	
	ActivityState string
	
	DisksActive int64
	
	DisksTotal int64
	
	BlocksTotal int64
	
	BlocksSynced int64
}


func (fs FS) ParseMDStat() (mdstates []MDStat, err error) {
	mdStatusFilePath := fs.Path("mdstat")
	content, err := ioutil.ReadFile(mdStatusFilePath)
	if err != nil {
		return []MDStat{}, fmt.Errorf("error parsing %s: %s", mdStatusFilePath, err)
	}

	mdStates := []MDStat{}
	lines := strings.Split(string(content), "\n")
	for i, l := range lines {
		if l == "" {
			continue
		}
		if l[0] == ' ' {
			continue
		}
		if strings.HasPrefix(l, "Personalities") || strings.HasPrefix(l, "unused") {
			continue
		}

		mainLine := strings.Split(l, " ")
		if len(mainLine) < 3 {
			return mdStates, fmt.Errorf("error parsing mdline: %s", l)
		}
		mdName := mainLine[0]
		activityState := mainLine[2]

		if len(lines) <= i+3 {
			return mdStates, fmt.Errorf(
				"error parsing %s: too few lines for md device %s",
				mdStatusFilePath,
				mdName,
			)
		}

		active, total, size, err := evalStatusline(lines[i+1])
		if err != nil {
			return mdStates, fmt.Errorf("error parsing %s: %s", mdStatusFilePath, err)
		}

		
		j := i + 2
		if strings.Contains(lines[i+2], "bitmap") { 
			j = i + 3
		}

		
		
		syncedBlocks := size
		if strings.Contains(lines[j], "recovery") || strings.Contains(lines[j], "resync") {
			syncedBlocks, err = evalBuildline(lines[j])
			if err != nil {
				return mdStates, fmt.Errorf("error parsing %s: %s", mdStatusFilePath, err)
			}
		}

		mdStates = append(mdStates, MDStat{
			Name:          mdName,
			ActivityState: activityState,
			DisksActive:   active,
			DisksTotal:    total,
			BlocksTotal:   size,
			BlocksSynced:  syncedBlocks,
		})
	}

	return mdStates, nil
}

func evalStatusline(statusline string) (active, total, size int64, err error) {
	matches := statuslineRE.FindStringSubmatch(statusline)
	if len(matches) != 4 {
		return 0, 0, 0, fmt.Errorf("unexpected statusline: %s", statusline)
	}

	size, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unexpected statusline %s: %s", statusline, err)
	}

	total, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unexpected statusline %s: %s", statusline, err)
	}

	active, err = strconv.ParseInt(matches[3], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unexpected statusline %s: %s", statusline, err)
	}

	return active, total, size, nil
}

func evalBuildline(buildline string) (syncedBlocks int64, err error) {
	matches := buildlineRE.FindStringSubmatch(buildline)
	if len(matches) != 2 {
		return 0, fmt.Errorf("unexpected buildline: %s", buildline)
	}

	syncedBlocks, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s in buildline: %s", err, buildline)
	}

	return syncedBlocks, nil
}
