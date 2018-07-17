package idtools 

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)




type IDMap struct {
	ContainerID int `json:"container_id"`
	HostID      int `json:"host_id"`
	Size        int `json:"size"`
}

type subIDRange struct {
	Start  int
	Length int
}

type ranges []subIDRange

func (e ranges) Len() int           { return len(e) }
func (e ranges) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ranges) Less(i, j int) bool { return e[i].Start < e[j].Start }

const (
	subuidFileName string = "/etc/subuid"
	subgidFileName string = "/etc/subgid"
)




func MkdirAllAndChown(path string, mode os.FileMode, owner IDPair) error {
	return mkdirAs(path, mode, owner.UID, owner.GID, true, true)
}





func MkdirAndChown(path string, mode os.FileMode, owner IDPair) error {
	return mkdirAs(path, mode, owner.UID, owner.GID, false, true)
}




func MkdirAllAndChownNew(path string, mode os.FileMode, owner IDPair) error {
	return mkdirAs(path, mode, owner.UID, owner.GID, true, false)
}



func GetRootUIDGID(uidMap, gidMap []IDMap) (int, int, error) {
	uid, err := toHost(0, uidMap)
	if err != nil {
		return -1, -1, err
	}
	gid, err := toHost(0, gidMap)
	if err != nil {
		return -1, -1, err
	}
	return uid, gid, nil
}




func toContainer(hostID int, idMap []IDMap) (int, error) {
	if idMap == nil {
		return hostID, nil
	}
	for _, m := range idMap {
		if (hostID >= m.HostID) && (hostID <= (m.HostID + m.Size - 1)) {
			contID := m.ContainerID + (hostID - m.HostID)
			return contID, nil
		}
	}
	return -1, fmt.Errorf("Host ID %d cannot be mapped to a container ID", hostID)
}




func toHost(contID int, idMap []IDMap) (int, error) {
	if idMap == nil {
		return contID, nil
	}
	for _, m := range idMap {
		if (contID >= m.ContainerID) && (contID <= (m.ContainerID + m.Size - 1)) {
			hostID := m.HostID + (contID - m.ContainerID)
			return hostID, nil
		}
	}
	return -1, fmt.Errorf("Container ID %d cannot be mapped to a host ID", contID)
}


type IDPair struct {
	UID int
	GID int
}


type IDMappings struct {
	uids []IDMap
	gids []IDMap
}




func NewIDMappings(username, groupname string) (*IDMappings, error) {
	subuidRanges, err := parseSubuid(username)
	if err != nil {
		return nil, err
	}
	subgidRanges, err := parseSubgid(groupname)
	if err != nil {
		return nil, err
	}
	if len(subuidRanges) == 0 {
		return nil, fmt.Errorf("No subuid ranges found for user %q", username)
	}
	if len(subgidRanges) == 0 {
		return nil, fmt.Errorf("No subgid ranges found for group %q", groupname)
	}

	return &IDMappings{
		uids: createIDMap(subuidRanges),
		gids: createIDMap(subgidRanges),
	}, nil
}



func NewIDMappingsFromMaps(uids []IDMap, gids []IDMap) *IDMappings {
	return &IDMappings{uids: uids, gids: gids}
}




func (i *IDMappings) RootPair() IDPair {
	uid, gid, _ := GetRootUIDGID(i.uids, i.gids)
	return IDPair{UID: uid, GID: gid}
}



func (i *IDMappings) ToHost(pair IDPair) (IDPair, error) {
	var err error
	target := i.RootPair()

	if pair.UID != target.UID {
		target.UID, err = toHost(pair.UID, i.uids)
		if err != nil {
			return target, err
		}
	}

	if pair.GID != target.GID {
		target.GID, err = toHost(pair.GID, i.gids)
	}
	return target, err
}


func (i *IDMappings) ToContainer(pair IDPair) (int, int, error) {
	uid, err := toContainer(pair.UID, i.uids)
	if err != nil {
		return -1, -1, err
	}
	gid, err := toContainer(pair.GID, i.gids)
	return uid, gid, err
}


func (i *IDMappings) Empty() bool {
	return len(i.uids) == 0 && len(i.gids) == 0
}



func (i *IDMappings) UIDs() []IDMap {
	return i.uids
}



func (i *IDMappings) GIDs() []IDMap {
	return i.gids
}

func createIDMap(subidRanges ranges) []IDMap {
	idMap := []IDMap{}

	
	sort.Sort(subidRanges)
	containerID := 0
	for _, idrange := range subidRanges {
		idMap = append(idMap, IDMap{
			ContainerID: containerID,
			HostID:      idrange.Start,
			Size:        idrange.Length,
		})
		containerID = containerID + idrange.Length
	}
	return idMap
}

func parseSubuid(username string) (ranges, error) {
	return parseSubidFile(subuidFileName, username)
}

func parseSubgid(username string) (ranges, error) {
	return parseSubidFile(subgidFileName, username)
}




func parseSubidFile(path, username string) (ranges, error) {
	var rangeList ranges

	subidFile, err := os.Open(path)
	if err != nil {
		return rangeList, err
	}
	defer subidFile.Close()

	s := bufio.NewScanner(subidFile)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return rangeList, err
		}

		text := strings.TrimSpace(s.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		parts := strings.Split(text, ":")
		if len(parts) != 3 {
			return rangeList, fmt.Errorf("Cannot parse subuid/gid information: Format not correct for %s file", path)
		}
		if parts[0] == username || username == "ALL" {
			startid, err := strconv.Atoi(parts[1])
			if err != nil {
				return rangeList, fmt.Errorf("String to int conversion failed during subuid/gid parsing of %s: %v", path, err)
			}
			length, err := strconv.Atoi(parts[2])
			if err != nil {
				return rangeList, fmt.Errorf("String to int conversion failed during subuid/gid parsing of %s: %v", path, err)
			}
			rangeList = append(rangeList, subIDRange{startid, length})
		}
	}
	return rangeList, nil
}
