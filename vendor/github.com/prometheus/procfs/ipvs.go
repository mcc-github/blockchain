












package procfs

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)


type IPVSStats struct {
	
	Connections uint64
	
	IncomingPackets uint64
	
	OutgoingPackets uint64
	
	IncomingBytes uint64
	
	OutgoingBytes uint64
}


type IPVSBackendStatus struct {
	
	LocalAddress net.IP
	
	RemoteAddress net.IP
	
	LocalPort uint16
	
	RemotePort uint16
	
	LocalMark string
	
	Proto string
	
	ActiveConn uint64
	
	InactConn uint64
	
	Weight uint64
}


func NewIPVSStats() (IPVSStats, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return IPVSStats{}, err
	}

	return fs.NewIPVSStats()
}


func (fs FS) NewIPVSStats() (IPVSStats, error) {
	file, err := os.Open(fs.Path("net/ip_vs_stats"))
	if err != nil {
		return IPVSStats{}, err
	}
	defer file.Close()

	return parseIPVSStats(file)
}


func parseIPVSStats(file io.Reader) (IPVSStats, error) {
	var (
		statContent []byte
		statLines   []string
		statFields  []string
		stats       IPVSStats
	)

	statContent, err := ioutil.ReadAll(file)
	if err != nil {
		return IPVSStats{}, err
	}

	statLines = strings.SplitN(string(statContent), "\n", 4)
	if len(statLines) != 4 {
		return IPVSStats{}, errors.New("ip_vs_stats corrupt: too short")
	}

	statFields = strings.Fields(statLines[2])
	if len(statFields) != 5 {
		return IPVSStats{}, errors.New("ip_vs_stats corrupt: unexpected number of fields")
	}

	stats.Connections, err = strconv.ParseUint(statFields[0], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}
	stats.IncomingPackets, err = strconv.ParseUint(statFields[1], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}
	stats.OutgoingPackets, err = strconv.ParseUint(statFields[2], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}
	stats.IncomingBytes, err = strconv.ParseUint(statFields[3], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}
	stats.OutgoingBytes, err = strconv.ParseUint(statFields[4], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}

	return stats, nil
}


func NewIPVSBackendStatus() ([]IPVSBackendStatus, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return []IPVSBackendStatus{}, err
	}

	return fs.NewIPVSBackendStatus()
}


func (fs FS) NewIPVSBackendStatus() ([]IPVSBackendStatus, error) {
	file, err := os.Open(fs.Path("net/ip_vs"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseIPVSBackendStatus(file)
}

func parseIPVSBackendStatus(file io.Reader) ([]IPVSBackendStatus, error) {
	var (
		status       []IPVSBackendStatus
		scanner      = bufio.NewScanner(file)
		proto        string
		localMark    string
		localAddress net.IP
		localPort    uint16
		err          error
	)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		switch {
		case fields[0] == "IP" || fields[0] == "Prot" || fields[1] == "RemoteAddress:Port":
			continue
		case fields[0] == "TCP" || fields[0] == "UDP":
			if len(fields) < 2 {
				continue
			}
			proto = fields[0]
			localMark = ""
			localAddress, localPort, err = parseIPPort(fields[1])
			if err != nil {
				return nil, err
			}
		case fields[0] == "FWM":
			if len(fields) < 2 {
				continue
			}
			proto = fields[0]
			localMark = fields[1]
			localAddress = nil
			localPort = 0
		case fields[0] == "->":
			if len(fields) < 6 {
				continue
			}
			remoteAddress, remotePort, err := parseIPPort(fields[1])
			if err != nil {
				return nil, err
			}
			weight, err := strconv.ParseUint(fields[3], 10, 64)
			if err != nil {
				return nil, err
			}
			activeConn, err := strconv.ParseUint(fields[4], 10, 64)
			if err != nil {
				return nil, err
			}
			inactConn, err := strconv.ParseUint(fields[5], 10, 64)
			if err != nil {
				return nil, err
			}
			status = append(status, IPVSBackendStatus{
				LocalAddress:  localAddress,
				LocalPort:     localPort,
				LocalMark:     localMark,
				RemoteAddress: remoteAddress,
				RemotePort:    remotePort,
				Proto:         proto,
				Weight:        weight,
				ActiveConn:    activeConn,
				InactConn:     inactConn,
			})
		}
	}
	return status, nil
}

func parseIPPort(s string) (net.IP, uint16, error) {
	var (
		ip  net.IP
		err error
	)

	switch len(s) {
	case 13:
		ip, err = hex.DecodeString(s[0:8])
		if err != nil {
			return nil, 0, err
		}
	case 46:
		ip = net.ParseIP(s[1:40])
		if ip == nil {
			return nil, 0, fmt.Errorf("invalid IPv6 address: %s", s[1:40])
		}
	default:
		return nil, 0, fmt.Errorf("unexpected IP:Port: %s", s)
	}

	portString := s[len(s)-4:]
	if len(portString) != 4 {
		return nil, 0, fmt.Errorf("unexpected port string format: %s", portString)
	}
	port, err := strconv.ParseUint(portString, 16, 16)
	if err != nil {
		return nil, 0, err
	}

	return ip, uint16(port), nil
}
