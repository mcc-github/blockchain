package network 
import (
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/errdefs"
)


type Address struct {
	Addr      string
	PrefixLen int
}


type IPAM struct {
	Driver  string
	Options map[string]string 
	Config  []IPAMConfig
}


type IPAMConfig struct {
	Subnet     string            `json:",omitempty"`
	IPRange    string            `json:",omitempty"`
	Gateway    string            `json:",omitempty"`
	AuxAddress map[string]string `json:"AuxiliaryAddresses,omitempty"`
}


type EndpointIPAMConfig struct {
	IPv4Address  string   `json:",omitempty"`
	IPv6Address  string   `json:",omitempty"`
	LinkLocalIPs []string `json:",omitempty"`
}


func (cfg *EndpointIPAMConfig) Copy() *EndpointIPAMConfig {
	cfgCopy := *cfg
	cfgCopy.LinkLocalIPs = make([]string, 0, len(cfg.LinkLocalIPs))
	cfgCopy.LinkLocalIPs = append(cfgCopy.LinkLocalIPs, cfg.LinkLocalIPs...)
	return &cfgCopy
}


type PeerInfo struct {
	Name string
	IP   string
}


type EndpointSettings struct {
	
	IPAMConfig *EndpointIPAMConfig
	Links      []string
	Aliases    []string
	
	NetworkID           string
	EndpointID          string
	Gateway             string
	IPAddress           string
	IPPrefixLen         int
	IPv6Gateway         string
	GlobalIPv6Address   string
	GlobalIPv6PrefixLen int
	MacAddress          string
	DriverOpts          map[string]string
}


type Task struct {
	Name       string
	EndpointID string
	EndpointIP string
	Info       map[string]string
}


type ServiceInfo struct {
	VIP          string
	Ports        []string
	LocalLBIndex int
	Tasks        []Task
}


func (es *EndpointSettings) Copy() *EndpointSettings {
	epCopy := *es
	if es.IPAMConfig != nil {
		epCopy.IPAMConfig = es.IPAMConfig.Copy()
	}

	if es.Links != nil {
		links := make([]string, 0, len(es.Links))
		epCopy.Links = append(links, es.Links...)
	}

	if es.Aliases != nil {
		aliases := make([]string, 0, len(es.Aliases))
		epCopy.Aliases = append(aliases, es.Aliases...)
	}
	return &epCopy
}



type NetworkingConfig struct {
	EndpointsConfig map[string]*EndpointSettings 
}


type ConfigReference struct {
	Network string
}

var acceptedFilters = map[string]bool{
	"driver": true,
	"type":   true,
	"name":   true,
	"id":     true,
	"label":  true,
	"scope":  true,
}


func ValidateFilters(filter filters.Args) error {
	return errdefs.InvalidParameter(filter.Validate(acceptedFilters))
}
