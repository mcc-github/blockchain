package swarm 

import (
	"github.com/docker/docker/api/types/network"
)


type Endpoint struct {
	Spec       EndpointSpec        `json:",omitempty"`
	Ports      []PortConfig        `json:",omitempty"`
	VirtualIPs []EndpointVirtualIP `json:",omitempty"`
}


type EndpointSpec struct {
	Mode  ResolutionMode `json:",omitempty"`
	Ports []PortConfig   `json:",omitempty"`
}


type ResolutionMode string

const (
	
	ResolutionModeVIP ResolutionMode = "vip"
	
	ResolutionModeDNSRR ResolutionMode = "dnsrr"
)


type PortConfig struct {
	Name     string             `json:",omitempty"`
	Protocol PortConfigProtocol `json:",omitempty"`
	
	TargetPort uint32 `json:",omitempty"`
	
	PublishedPort uint32 `json:",omitempty"`
	
	PublishMode PortConfigPublishMode `json:",omitempty"`
}



type PortConfigPublishMode string

const (
	
	
	PortConfigPublishModeIngress PortConfigPublishMode = "ingress"
	
	
	PortConfigPublishModeHost PortConfigPublishMode = "host"
)


type PortConfigProtocol string

const (
	

	
	PortConfigProtocolTCP PortConfigProtocol = "tcp"
	
	PortConfigProtocolUDP PortConfigProtocol = "udp"
	
	PortConfigProtocolSCTP PortConfigProtocol = "sctp"
)


type EndpointVirtualIP struct {
	NetworkID string `json:",omitempty"`
	Addr      string `json:",omitempty"`
}


type Network struct {
	ID string
	Meta
	Spec        NetworkSpec  `json:",omitempty"`
	DriverState Driver       `json:",omitempty"`
	IPAMOptions *IPAMOptions `json:",omitempty"`
}


type NetworkSpec struct {
	Annotations
	DriverConfiguration *Driver                  `json:",omitempty"`
	IPv6Enabled         bool                     `json:",omitempty"`
	Internal            bool                     `json:",omitempty"`
	Attachable          bool                     `json:",omitempty"`
	Ingress             bool                     `json:",omitempty"`
	IPAMOptions         *IPAMOptions             `json:",omitempty"`
	ConfigFrom          *network.ConfigReference `json:",omitempty"`
	Scope               string                   `json:",omitempty"`
}


type NetworkAttachmentConfig struct {
	Target     string            `json:",omitempty"`
	Aliases    []string          `json:",omitempty"`
	DriverOpts map[string]string `json:",omitempty"`
}


type NetworkAttachment struct {
	Network   Network  `json:",omitempty"`
	Addresses []string `json:",omitempty"`
}


type IPAMOptions struct {
	Driver  Driver       `json:",omitempty"`
	Configs []IPAMConfig `json:",omitempty"`
}


type IPAMConfig struct {
	Subnet  string `json:",omitempty"`
	Range   string `json:",omitempty"`
	Gateway string `json:",omitempty"`
}
