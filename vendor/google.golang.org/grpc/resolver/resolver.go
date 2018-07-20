



package resolver

var (
	
	m = make(map[string]Builder)
	
	defaultScheme = "passthrough"
)









func Register(b Builder) {
	m[b.Scheme()] = b
}




func Get(scheme string) Builder {
	if b, ok := m[scheme]; ok {
		return b
	}
	return nil
}



func SetDefaultScheme(scheme string) {
	defaultScheme = scheme
}


func GetDefaultScheme() string {
	return defaultScheme
}


type AddressType uint8

const (
	
	Backend AddressType = iota
	
	GRPCLB
)



type Address struct {
	
	Addr string
	
	Type AddressType
	
	
	
	
	ServerName string
	
	
	Metadata interface{}
}



type BuildOption struct {
	
	DisableServiceConfig bool
}








type ClientConn interface {
	
	
	
	NewAddress(addresses []Address)
	
	
	NewServiceConfig(serviceConfig string)
}



type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}


type Builder interface {
	
	
	
	
	Build(target Target, cc ClientConn, opts BuildOption) (Resolver, error)
	
	
	Scheme() string
}


type ResolveNowOption struct{}



type Resolver interface {
	
	
	
	
	ResolveNow(ResolveNowOption)
	
	Close()
}




func UnregisterForTesting(scheme string) {
	delete(m, scheme)
}
