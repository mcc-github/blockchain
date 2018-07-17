



package naming

import "net"

var (
	lookupHost = net.DefaultResolver.LookupHost
	lookupSRV  = net.DefaultResolver.LookupSRV
)
