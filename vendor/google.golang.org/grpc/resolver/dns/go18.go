



package dns

import "net"

var (
	lookupHost = net.DefaultResolver.LookupHost
	lookupSRV  = net.DefaultResolver.LookupSRV
	lookupTXT  = net.DefaultResolver.LookupTXT
)
