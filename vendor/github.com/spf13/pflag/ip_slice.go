package pflag

import (
	"fmt"
	"io"
	"net"
	"strings"
)


type ipSliceValue struct {
	value   *[]net.IP
	changed bool
}

func newIPSliceValue(val []net.IP, p *[]net.IP) *ipSliceValue {
	ipsv := new(ipSliceValue)
	ipsv.value = p
	*ipsv.value = val
	return ipsv
}



func (s *ipSliceValue) Set(val string) error {

	
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	
	ipStrSlice, err := readAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	
	out := make([]net.IP, 0, len(ipStrSlice))
	for _, ipStr := range ipStrSlice {
		ip := net.ParseIP(strings.TrimSpace(ipStr))
		if ip == nil {
			return fmt.Errorf("invalid string being converted to IP address: %s", ipStr)
		}
		out = append(out, ip)
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

	return nil
}


func (s *ipSliceValue) Type() string {
	return "ipSlice"
}


func (s *ipSliceValue) String() string {

	ipStrSlice := make([]string, len(*s.value))
	for i, ip := range *s.value {
		ipStrSlice[i] = ip.String()
	}

	out, _ := writeAsCSV(ipStrSlice)

	return "[" + out + "]"
}

func ipSliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	
	if len(val) == 0 {
		return []net.IP{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]net.IP, len(ss))
	for i, sval := range ss {
		ip := net.ParseIP(strings.TrimSpace(sval))
		if ip == nil {
			return nil, fmt.Errorf("invalid string being converted to IP address: %s", sval)
		}
		out[i] = ip
	}
	return out, nil
}


func (f *FlagSet) GetIPSlice(name string) ([]net.IP, error) {
	val, err := f.getFlagType(name, "ipSlice", ipSliceConv)
	if err != nil {
		return []net.IP{}, err
	}
	return val.([]net.IP), nil
}



func (f *FlagSet) IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string) {
	f.VarP(newIPSliceValue(value, p), name, "", usage)
}


func (f *FlagSet) IPSliceVarP(p *[]net.IP, name, shorthand string, value []net.IP, usage string) {
	f.VarP(newIPSliceValue(value, p), name, shorthand, usage)
}



func IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string) {
	CommandLine.VarP(newIPSliceValue(value, p), name, "", usage)
}


func IPSliceVarP(p *[]net.IP, name, shorthand string, value []net.IP, usage string) {
	CommandLine.VarP(newIPSliceValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) IPSlice(name string, value []net.IP, usage string) *[]net.IP {
	p := []net.IP{}
	f.IPSliceVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) IPSliceP(name, shorthand string, value []net.IP, usage string) *[]net.IP {
	p := []net.IP{}
	f.IPSliceVarP(&p, name, shorthand, value, usage)
	return &p
}



func IPSlice(name string, value []net.IP, usage string) *[]net.IP {
	return CommandLine.IPSliceP(name, "", value, usage)
}


func IPSliceP(name, shorthand string, value []net.IP, usage string) *[]net.IP {
	return CommandLine.IPSliceP(name, shorthand, value, usage)
}
