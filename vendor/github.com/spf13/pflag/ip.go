package pflag

import (
	"fmt"
	"net"
	"strings"
)


type ipValue net.IP

func newIPValue(val net.IP, p *net.IP) *ipValue {
	*p = val
	return (*ipValue)(p)
}

func (i *ipValue) String() string { return net.IP(*i).String() }
func (i *ipValue) Set(s string) error {
	ip := net.ParseIP(strings.TrimSpace(s))
	if ip == nil {
		return fmt.Errorf("failed to parse IP: %q", s)
	}
	*i = ipValue(ip)
	return nil
}

func (i *ipValue) Type() string {
	return "ip"
}

func ipConv(sval string) (interface{}, error) {
	ip := net.ParseIP(sval)
	if ip != nil {
		return ip, nil
	}
	return nil, fmt.Errorf("invalid string being converted to IP address: %s", sval)
}


func (f *FlagSet) GetIP(name string) (net.IP, error) {
	val, err := f.getFlagType(name, "ip", ipConv)
	if err != nil {
		return nil, err
	}
	return val.(net.IP), nil
}



func (f *FlagSet) IPVar(p *net.IP, name string, value net.IP, usage string) {
	f.VarP(newIPValue(value, p), name, "", usage)
}


func (f *FlagSet) IPVarP(p *net.IP, name, shorthand string, value net.IP, usage string) {
	f.VarP(newIPValue(value, p), name, shorthand, usage)
}



func IPVar(p *net.IP, name string, value net.IP, usage string) {
	CommandLine.VarP(newIPValue(value, p), name, "", usage)
}


func IPVarP(p *net.IP, name, shorthand string, value net.IP, usage string) {
	CommandLine.VarP(newIPValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) IP(name string, value net.IP, usage string) *net.IP {
	p := new(net.IP)
	f.IPVarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) IPP(name, shorthand string, value net.IP, usage string) *net.IP {
	p := new(net.IP)
	f.IPVarP(p, name, shorthand, value, usage)
	return p
}



func IP(name string, value net.IP, usage string) *net.IP {
	return CommandLine.IPP(name, "", value, usage)
}


func IPP(name, shorthand string, value net.IP, usage string) *net.IP {
	return CommandLine.IPP(name, shorthand, value, usage)
}
