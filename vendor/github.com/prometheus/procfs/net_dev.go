












package procfs

import (
	"bufio"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
)


type NetDevLine struct {
	Name         string `json:"name"`          
	RxBytes      uint64 `json:"rx_bytes"`      
	RxPackets    uint64 `json:"rx_packets"`    
	RxErrors     uint64 `json:"rx_errors"`     
	RxDropped    uint64 `json:"rx_dropped"`    
	RxFIFO       uint64 `json:"rx_fifo"`       
	RxFrame      uint64 `json:"rx_frame"`      
	RxCompressed uint64 `json:"rx_compressed"` 
	RxMulticast  uint64 `json:"rx_multicast"`  
	TxBytes      uint64 `json:"tx_bytes"`      
	TxPackets    uint64 `json:"tx_packets"`    
	TxErrors     uint64 `json:"tx_errors"`     
	TxDropped    uint64 `json:"tx_dropped"`    
	TxFIFO       uint64 `json:"tx_fifo"`       
	TxCollisions uint64 `json:"tx_collisions"` 
	TxCarrier    uint64 `json:"tx_carrier"`    
	TxCompressed uint64 `json:"tx_compressed"` 
}



type NetDev map[string]NetDevLine


func NewNetDev() (NetDev, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewNetDev()
}


func (fs FS) NewNetDev() (NetDev, error) {
	return newNetDev(fs.Path("net/dev"))
}


func (p Proc) NewNetDev() (NetDev, error) {
	return newNetDev(p.path("net/dev"))
}


func newNetDev(file string) (NetDev, error) {
	f, err := os.Open(file)
	if err != nil {
		return NetDev{}, err
	}
	defer f.Close()

	nd := NetDev{}
	s := bufio.NewScanner(f)
	for n := 0; s.Scan(); n++ {
		
		if n < 2 {
			continue
		}

		line, err := nd.parseLine(s.Text())
		if err != nil {
			return nd, err
		}

		nd[line.Name] = *line
	}

	return nd, s.Err()
}



func (nd NetDev) parseLine(rawLine string) (*NetDevLine, error) {
	parts := strings.SplitN(rawLine, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid net/dev line, missing colon")
	}
	fields := strings.Fields(strings.TrimSpace(parts[1]))

	var err error
	line := &NetDevLine{}

	
	line.Name = strings.TrimSpace(parts[0])
	if line.Name == "" {
		return nil, errors.New("invalid net/dev line, empty interface name")
	}

	
	line.RxBytes, err = strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return nil, err
	}
	line.RxPackets, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return nil, err
	}
	line.RxErrors, err = strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return nil, err
	}
	line.RxDropped, err = strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}
	line.RxFIFO, err = strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return nil, err
	}
	line.RxFrame, err = strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return nil, err
	}
	line.RxCompressed, err = strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		return nil, err
	}
	line.RxMulticast, err = strconv.ParseUint(fields[7], 10, 64)
	if err != nil {
		return nil, err
	}

	
	line.TxBytes, err = strconv.ParseUint(fields[8], 10, 64)
	if err != nil {
		return nil, err
	}
	line.TxPackets, err = strconv.ParseUint(fields[9], 10, 64)
	if err != nil {
		return nil, err
	}
	line.TxErrors, err = strconv.ParseUint(fields[10], 10, 64)
	if err != nil {
		return nil, err
	}
	line.TxDropped, err = strconv.ParseUint(fields[11], 10, 64)
	if err != nil {
		return nil, err
	}
	line.TxFIFO, err = strconv.ParseUint(fields[12], 10, 64)
	if err != nil {
		return nil, err
	}
	line.TxCollisions, err = strconv.ParseUint(fields[13], 10, 64)
	if err != nil {
		return nil, err
	}
	line.TxCarrier, err = strconv.ParseUint(fields[14], 10, 64)
	if err != nil {
		return nil, err
	}
	line.TxCompressed, err = strconv.ParseUint(fields[15], 10, 64)
	if err != nil {
		return nil, err
	}

	return line, nil
}



func (nd NetDev) Total() NetDevLine {
	total := NetDevLine{}

	names := make([]string, 0, len(nd))
	for _, ifc := range nd {
		names = append(names, ifc.Name)
		total.RxBytes += ifc.RxBytes
		total.RxPackets += ifc.RxPackets
		total.RxPackets += ifc.RxPackets
		total.RxErrors += ifc.RxErrors
		total.RxDropped += ifc.RxDropped
		total.RxFIFO += ifc.RxFIFO
		total.RxFrame += ifc.RxFrame
		total.RxCompressed += ifc.RxCompressed
		total.RxMulticast += ifc.RxMulticast
		total.TxBytes += ifc.TxBytes
		total.TxPackets += ifc.TxPackets
		total.TxErrors += ifc.TxErrors
		total.TxDropped += ifc.TxDropped
		total.TxFIFO += ifc.TxFIFO
		total.TxCollisions += ifc.TxCollisions
		total.TxCarrier += ifc.TxCarrier
		total.TxCompressed += ifc.TxCompressed
	}
	sort.Strings(names)
	total.Name = strings.Join(names, ", ")

	return total
}
