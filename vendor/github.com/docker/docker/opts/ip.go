package opts 

import (
	"fmt"
	"net"
)


type IPOpt struct {
	*net.IP
}




func NewIPOpt(ref *net.IP, defaultVal string) *IPOpt {
	o := &IPOpt{
		IP: ref,
	}
	o.Set(defaultVal)
	return o
}



func (o *IPOpt) Set(val string) error {
	ip := net.ParseIP(val)
	if ip == nil {
		return fmt.Errorf("%s is not an ip address", val)
	}
	*o.IP = ip
	return nil
}



func (o *IPOpt) String() string {
	if *o.IP == nil {
		return ""
	}
	return o.IP.String()
}


func (o *IPOpt) Type() string {
	return "ip"
}
