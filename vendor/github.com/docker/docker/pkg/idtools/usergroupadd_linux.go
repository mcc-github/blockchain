package idtools 

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)






var (
	once        sync.Once
	userCommand string

	cmdTemplates = map[string]string{
		"adduser": "--system --shell /bin/false --no-create-home --disabled-login --disabled-password --group %s",
		"useradd": "-r -s /bin/false %s",
		"usermod": "-%s %d-%d %s",
	}

	idOutRegexp = regexp.MustCompile(`uid=([0-9]+).*gid=([0-9]+)`)
	
	defaultRangeLen   = 65536
	defaultRangeStart = 100000
	userMod           = "usermod"
)





func AddNamespaceRangesUser(name string) (int, int, error) {
	if err := addUser(name); err != nil {
		return -1, -1, fmt.Errorf("Error adding user %q: %v", name, err)
	}

	
	out, err := execCmd("id", name)
	if err != nil {
		return -1, -1, fmt.Errorf("Error trying to find uid/gid for new user %q: %v", name, err)
	}
	matches := idOutRegexp.FindStringSubmatch(strings.TrimSpace(string(out)))
	if len(matches) != 3 {
		return -1, -1, fmt.Errorf("Can't find uid, gid from `id` output: %q", string(out))
	}
	uid, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, -1, fmt.Errorf("Can't convert found uid (%s) to int: %v", matches[1], err)
	}
	gid, err := strconv.Atoi(matches[2])
	if err != nil {
		return -1, -1, fmt.Errorf("Can't convert found gid (%s) to int: %v", matches[2], err)
	}

	
	

	if err := createSubordinateRanges(name); err != nil {
		return -1, -1, fmt.Errorf("Couldn't create subordinate ID ranges: %v", err)
	}
	return uid, gid, nil
}

func addUser(userName string) error {
	once.Do(func() {
		
		if _, err := resolveBinary("adduser"); err == nil {
			userCommand = "adduser"
		} else if _, err := resolveBinary("useradd"); err == nil {
			userCommand = "useradd"
		}
	})
	if userCommand == "" {
		return fmt.Errorf("Cannot add user; no useradd/adduser binary found")
	}
	args := fmt.Sprintf(cmdTemplates[userCommand], userName)
	out, err := execCmd(userCommand, args)
	if err != nil {
		return fmt.Errorf("Failed to add user with error: %v; output: %q", err, string(out))
	}
	return nil
}

func createSubordinateRanges(name string) error {

	
	
	ranges, err := parseSubuid(name)
	if err != nil {
		return fmt.Errorf("Error while looking for subuid ranges for user %q: %v", name, err)
	}
	if len(ranges) == 0 {
		
		startID, err := findNextUIDRange()
		if err != nil {
			return fmt.Errorf("Can't find available subuid range: %v", err)
		}
		out, err := execCmd(userMod, fmt.Sprintf(cmdTemplates[userMod], "v", startID, startID+defaultRangeLen-1, name))
		if err != nil {
			return fmt.Errorf("Unable to add subuid range to user: %q; output: %s, err: %v", name, out, err)
		}
	}

	ranges, err = parseSubgid(name)
	if err != nil {
		return fmt.Errorf("Error while looking for subgid ranges for user %q: %v", name, err)
	}
	if len(ranges) == 0 {
		
		startID, err := findNextGIDRange()
		if err != nil {
			return fmt.Errorf("Can't find available subgid range: %v", err)
		}
		out, err := execCmd(userMod, fmt.Sprintf(cmdTemplates[userMod], "w", startID, startID+defaultRangeLen-1, name))
		if err != nil {
			return fmt.Errorf("Unable to add subgid range to user: %q; output: %s, err: %v", name, out, err)
		}
	}
	return nil
}

func findNextUIDRange() (int, error) {
	ranges, err := parseSubuid("ALL")
	if err != nil {
		return -1, fmt.Errorf("Couldn't parse all ranges in /etc/subuid file: %v", err)
	}
	sort.Sort(ranges)
	return findNextRangeStart(ranges)
}

func findNextGIDRange() (int, error) {
	ranges, err := parseSubgid("ALL")
	if err != nil {
		return -1, fmt.Errorf("Couldn't parse all ranges in /etc/subgid file: %v", err)
	}
	sort.Sort(ranges)
	return findNextRangeStart(ranges)
}

func findNextRangeStart(rangeList ranges) (int, error) {
	startID := defaultRangeStart
	for _, arange := range rangeList {
		if wouldOverlap(arange, startID) {
			startID = arange.Start + arange.Length
		}
	}
	return startID, nil
}

func wouldOverlap(arange subIDRange, ID int) bool {
	low := ID
	high := ID + defaultRangeLen
	if (low >= arange.Start && low <= arange.Start+arange.Length) ||
		(high <= arange.Start+arange.Length && high >= arange.Start) {
		return true
	}
	return false
}
