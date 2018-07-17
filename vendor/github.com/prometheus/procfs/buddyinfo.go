












package procfs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)




type BuddyInfo struct {
	Node  string
	Zone  string
	Sizes []float64
}


func NewBuddyInfo() ([]BuddyInfo, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewBuddyInfo()
}


func (fs FS) NewBuddyInfo() ([]BuddyInfo, error) {
	file, err := os.Open(fs.Path("buddyinfo"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseBuddyInfo(file)
}

func parseBuddyInfo(r io.Reader) ([]BuddyInfo, error) {
	var (
		buddyInfo   = []BuddyInfo{}
		scanner     = bufio.NewScanner(r)
		bucketCount = -1
	)

	for scanner.Scan() {
		var err error
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid number of fields when parsing buddyinfo")
		}

		node := strings.TrimRight(parts[1], ",")
		zone := strings.TrimRight(parts[3], ",")
		arraySize := len(parts[4:])

		if bucketCount == -1 {
			bucketCount = arraySize
		} else {
			if bucketCount != arraySize {
				return nil, fmt.Errorf("mismatch in number of buddyinfo buckets, previous count %d, new count %d", bucketCount, arraySize)
			}
		}

		sizes := make([]float64, arraySize)
		for i := 0; i < arraySize; i++ {
			sizes[i], err = strconv.ParseFloat(parts[i+4], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value in buddyinfo: %s", err)
			}
		}

		buddyInfo = append(buddyInfo, BuddyInfo{node, zone, sizes})
	}

	return buddyInfo, scanner.Err()
}
