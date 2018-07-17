












package nfs

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)


func ParseClientRPCStats(r io.Reader) (*ClientRPCStats, error) {
	stats := &ClientRPCStats{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(scanner.Text())
		
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid NFS metric line %q", line)
		}

		values, err := util.ParseUint64s(parts[1:])
		if err != nil {
			return nil, fmt.Errorf("error parsing NFS metric line: %s", err)
		}

		switch metricLine := parts[0]; metricLine {
		case "net":
			stats.Network, err = parseNetwork(values)
		case "rpc":
			stats.ClientRPC, err = parseClientRPC(values)
		case "proc2":
			stats.V2Stats, err = parseV2Stats(values)
		case "proc3":
			stats.V3Stats, err = parseV3Stats(values)
		case "proc4":
			stats.ClientV4Stats, err = parseClientV4Stats(values)
		default:
			return nil, fmt.Errorf("unknown NFS metric line %q", metricLine)
		}
		if err != nil {
			return nil, fmt.Errorf("errors parsing NFS metric line: %s", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning NFS file: %s", err)
	}

	return stats, nil
}
