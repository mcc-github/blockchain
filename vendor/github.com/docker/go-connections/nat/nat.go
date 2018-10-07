
package nat

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	
	portSpecTemplate = "ip:hostPort:containerPort"
)


type PortBinding struct {
	
	HostIP string `json:"HostIp"`
	
	HostPort string
}


type PortMap map[Port][]PortBinding


type PortSet map[Port]struct{}


type Port string


func NewPort(proto, port string) (Port, error) {
	
	

	portStartInt, portEndInt, err := ParsePortRangeToInt(port)
	if err != nil {
		return "", err
	}

	if portStartInt == portEndInt {
		return Port(fmt.Sprintf("%d/%s", portStartInt, proto)), nil
	}
	return Port(fmt.Sprintf("%d-%d/%s", portStartInt, portEndInt, proto)), nil
}


func ParsePort(rawPort string) (int, error) {
	if len(rawPort) == 0 {
		return 0, nil
	}
	port, err := strconv.ParseUint(rawPort, 10, 16)
	if err != nil {
		return 0, err
	}
	return int(port), nil
}


func ParsePortRangeToInt(rawPort string) (int, int, error) {
	if len(rawPort) == 0 {
		return 0, 0, nil
	}
	start, end, err := ParsePortRange(rawPort)
	if err != nil {
		return 0, 0, err
	}
	return int(start), int(end), nil
}


func (p Port) Proto() string {
	proto, _ := SplitProtoPort(string(p))
	return proto
}


func (p Port) Port() string {
	_, port := SplitProtoPort(string(p))
	return port
}


func (p Port) Int() int {
	portStr := p.Port()
	
	
	port, _ := ParsePort(portStr)
	return port
}


func (p Port) Range() (int, int, error) {
	return ParsePortRangeToInt(p.Port())
}


func SplitProtoPort(rawPort string) (string, string) {
	parts := strings.Split(rawPort, "/")
	l := len(parts)
	if len(rawPort) == 0 || l == 0 || len(parts[0]) == 0 {
		return "", ""
	}
	if l == 1 {
		return "tcp", rawPort
	}
	if len(parts[1]) == 0 {
		return "tcp", parts[0]
	}
	return parts[1], parts[0]
}

func validateProto(proto string) bool {
	for _, availableProto := range []string{"tcp", "udp", "sctp"} {
		if availableProto == proto {
			return true
		}
	}
	return false
}



func ParsePortSpecs(ports []string) (map[Port]struct{}, map[Port][]PortBinding, error) {
	var (
		exposedPorts = make(map[Port]struct{}, len(ports))
		bindings     = make(map[Port][]PortBinding)
	)
	for _, rawPort := range ports {
		portMappings, err := ParsePortSpec(rawPort)
		if err != nil {
			return nil, nil, err
		}

		for _, portMapping := range portMappings {
			port := portMapping.Port
			if _, exists := exposedPorts[port]; !exists {
				exposedPorts[port] = struct{}{}
			}
			bslice, exists := bindings[port]
			if !exists {
				bslice = []PortBinding{}
			}
			bindings[port] = append(bslice, portMapping.Binding)
		}
	}
	return exposedPorts, bindings, nil
}


type PortMapping struct {
	Port    Port
	Binding PortBinding
}

func splitParts(rawport string) (string, string, string) {
	parts := strings.Split(rawport, ":")
	n := len(parts)
	containerport := parts[n-1]

	switch n {
	case 1:
		return "", "", containerport
	case 2:
		return "", parts[0], containerport
	case 3:
		return parts[0], parts[1], containerport
	default:
		return strings.Join(parts[:n-2], ":"), parts[n-2], containerport
	}
}


func ParsePortSpec(rawPort string) ([]PortMapping, error) {
	var proto string
	rawIP, hostPort, containerPort := splitParts(rawPort)
	proto, containerPort = SplitProtoPort(containerPort)

	
	ip, _, err := net.SplitHostPort(rawIP + ":")
	if err != nil {
		return nil, fmt.Errorf("Invalid ip address %v: %s", rawIP, err)
	}
	if ip != "" && net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("Invalid ip address: %s", ip)
	}
	if containerPort == "" {
		return nil, fmt.Errorf("No port specified: %s<empty>", rawPort)
	}

	startPort, endPort, err := ParsePortRange(containerPort)
	if err != nil {
		return nil, fmt.Errorf("Invalid containerPort: %s", containerPort)
	}

	var startHostPort, endHostPort uint64 = 0, 0
	if len(hostPort) > 0 {
		startHostPort, endHostPort, err = ParsePortRange(hostPort)
		if err != nil {
			return nil, fmt.Errorf("Invalid hostPort: %s", hostPort)
		}
	}

	if hostPort != "" && (endPort-startPort) != (endHostPort-startHostPort) {
		
		
		
		if endPort != startPort {
			return nil, fmt.Errorf("Invalid ranges specified for container and host Ports: %s and %s", containerPort, hostPort)
		}
	}

	if !validateProto(strings.ToLower(proto)) {
		return nil, fmt.Errorf("Invalid proto: %s", proto)
	}

	ports := []PortMapping{}
	for i := uint64(0); i <= (endPort - startPort); i++ {
		containerPort = strconv.FormatUint(startPort+i, 10)
		if len(hostPort) > 0 {
			hostPort = strconv.FormatUint(startHostPort+i, 10)
		}
		
		
		if startPort == endPort && startHostPort != endHostPort {
			hostPort = fmt.Sprintf("%s-%s", hostPort, strconv.FormatUint(endHostPort, 10))
		}
		port, err := NewPort(strings.ToLower(proto), containerPort)
		if err != nil {
			return nil, err
		}

		binding := PortBinding{
			HostIP:   ip,
			HostPort: hostPort,
		}
		ports = append(ports, PortMapping{Port: port, Binding: binding})
	}
	return ports, nil
}
