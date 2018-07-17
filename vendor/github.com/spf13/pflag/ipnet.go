package pflag

import (
	"fmt"
	"net"
	"strings"
)


type ipNetValue net.IPNet

func (ipnet ipNetValue) String() string {
	n := net.IPNet(ipnet)
	return n.String()
}

func (ipnet *ipNetValue) Set(value string) error {
	_, n, err := net.ParseCIDR(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	*ipnet = ipNetValue(*n)
	return nil
}

func (*ipNetValue) Type() string {
	return "ipNet"
}

func newIPNetValue(val net.IPNet, p *net.IPNet) *ipNetValue {
	*p = val
	return (*ipNetValue)(p)
}

func ipNetConv(sval string) (interface{}, error) {
	_, n, err := net.ParseCIDR(strings.TrimSpace(sval))
	if err == nil {
		return *n, nil
	}
	return nil, fmt.Errorf("invalid string being converted to IPNet: %s", sval)
}


func (f *FlagSet) GetIPNet(name string) (net.IPNet, error) {
	val, err := f.getFlagType(name, "ipNet", ipNetConv)
	if err != nil {
		return net.IPNet{}, err
	}
	return val.(net.IPNet), nil
}



func (f *FlagSet) IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string) {
	f.VarP(newIPNetValue(value, p), name, "", usage)
}


func (f *FlagSet) IPNetVarP(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) {
	f.VarP(newIPNetValue(value, p), name, shorthand, usage)
}



func IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string) {
	CommandLine.VarP(newIPNetValue(value, p), name, "", usage)
}


func IPNetVarP(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) {
	CommandLine.VarP(newIPNetValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) IPNet(name string, value net.IPNet, usage string) *net.IPNet {
	p := new(net.IPNet)
	f.IPNetVarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) IPNetP(name, shorthand string, value net.IPNet, usage string) *net.IPNet {
	p := new(net.IPNet)
	f.IPNetVarP(p, name, shorthand, value, usage)
	return p
}



func IPNet(name string, value net.IPNet, usage string) *net.IPNet {
	return CommandLine.IPNetP(name, "", value, usage)
}


func IPNetP(name, shorthand string, value net.IPNet, usage string) *net.IPNet {
	return CommandLine.IPNetP(name, shorthand, value, usage)
}
