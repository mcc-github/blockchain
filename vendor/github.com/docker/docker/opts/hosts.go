package opts 

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

var (
	
	
	
	DefaultHTTPPort = 2375 
	
	DefaultTLSHTTPPort = 2376 
	
	
	DefaultUnixSocket = "/var/run/docker.sock"
	
	DefaultTCPHost = fmt.Sprintf("tcp://%s:%d", DefaultHTTPHost, DefaultHTTPPort)
	
	DefaultTLSHost = fmt.Sprintf("tcp://%s:%d", DefaultHTTPHost, DefaultTLSHTTPPort)
	
	DefaultNamedPipe = `
)


func ValidateHost(val string) (string, error) {
	host := strings.TrimSpace(val)
	
	if host != "" {
		_, err := parseDaemonHost(host)
		if err != nil {
			return val, err
		}
	}
	
	
	return val, nil
}


func ParseHost(defaultToTLS bool, val string) (string, error) {
	host := strings.TrimSpace(val)
	if host == "" {
		if defaultToTLS {
			host = DefaultTLSHost
		} else {
			host = DefaultHost
		}
	} else {
		var err error
		host, err = parseDaemonHost(host)
		if err != nil {
			return val, err
		}
	}
	return host, nil
}



func parseDaemonHost(addr string) (string, error) {
	addrParts := strings.SplitN(addr, "://", 2)
	if len(addrParts) == 1 && addrParts[0] != "" {
		addrParts = []string{"tcp", addrParts[0]}
	}

	switch addrParts[0] {
	case "tcp":
		return ParseTCPAddr(addrParts[1], DefaultTCPHost)
	case "unix":
		return parseSimpleProtoAddr("unix", addrParts[1], DefaultUnixSocket)
	case "npipe":
		return parseSimpleProtoAddr("npipe", addrParts[1], DefaultNamedPipe)
	case "fd":
		return addr, nil
	default:
		return "", fmt.Errorf("Invalid bind address format: %s", addr)
	}
}





func parseSimpleProtoAddr(proto, addr, defaultAddr string) (string, error) {
	addr = strings.TrimPrefix(addr, proto+"://")
	if strings.Contains(addr, "://") {
		return "", fmt.Errorf("Invalid proto, expected %s: %s", proto, addr)
	}
	if addr == "" {
		addr = defaultAddr
	}
	return fmt.Sprintf("%s://%s", proto, addr), nil
}






func ParseTCPAddr(tryAddr string, defaultAddr string) (string, error) {
	if tryAddr == "" || tryAddr == "tcp://" {
		return defaultAddr, nil
	}
	addr := strings.TrimPrefix(tryAddr, "tcp://")
	if strings.Contains(addr, "://") || addr == "" {
		return "", fmt.Errorf("Invalid proto, expected tcp: %s", tryAddr)
	}

	defaultAddr = strings.TrimPrefix(defaultAddr, "tcp://")
	defaultHost, defaultPort, err := net.SplitHostPort(defaultAddr)
	if err != nil {
		return "", err
	}
	
	
	
	if strings.HasSuffix(addr, "]:") {
		addr += defaultPort
	}

	u, err := url.Parse("tcp://" + addr)
	if err != nil {
		return "", err
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		
		host, port, err = net.SplitHostPort(net.JoinHostPort(u.Host, defaultPort))
	}
	if err != nil {
		return "", fmt.Errorf("Invalid bind address format: %s", tryAddr)
	}

	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	p, err := strconv.Atoi(port)
	if err != nil && p == 0 {
		return "", fmt.Errorf("Invalid bind address format: %s", tryAddr)
	}

	return fmt.Sprintf("tcp://%s%s", net.JoinHostPort(host, port), u.Path), nil
}



func ValidateExtraHost(val string) (string, error) {
	
	arr := strings.SplitN(val, ":", 2)
	if len(arr) != 2 || len(arr[0]) == 0 {
		return "", fmt.Errorf("bad format for add-host: %q", val)
	}
	if _, err := ValidateIPAddress(arr[1]); err != nil {
		return "", fmt.Errorf("invalid IP address in add-host: %q", arr[1])
	}
	return val, nil
}
