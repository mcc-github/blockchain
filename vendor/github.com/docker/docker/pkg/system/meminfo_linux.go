package system 

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/docker/go-units"
)



func ReadMemInfo() (*MemInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseMemInfo(file)
}




func parseMemInfo(reader io.Reader) (*MemInfo, error) {
	meminfo := &MemInfo{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		
		parts := strings.Fields(scanner.Text())

		
		if len(parts) < 3 || parts[2] != "kB" {
			continue
		}

		
		size, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		bytes := int64(size) * units.KiB

		switch parts[0] {
		case "MemTotal:":
			meminfo.MemTotal = bytes
		case "MemFree:":
			meminfo.MemFree = bytes
		case "SwapTotal:":
			meminfo.SwapTotal = bytes
		case "SwapFree:":
			meminfo.SwapFree = bytes
		}

	}

	
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return meminfo, nil
}
