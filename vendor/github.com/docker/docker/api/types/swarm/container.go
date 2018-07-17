package swarm 

import (
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)






type DNSConfig struct {
	
	Nameservers []string `json:",omitempty"`
	
	Search []string `json:",omitempty"`
	
	Options []string `json:",omitempty"`
}


type SELinuxContext struct {
	Disable bool

	User  string
	Role  string
	Type  string
	Level string
}


type CredentialSpec struct {
	File     string
	Registry string
}


type Privileges struct {
	CredentialSpec *CredentialSpec
	SELinuxContext *SELinuxContext
}


type ContainerSpec struct {
	Image           string                  `json:",omitempty"`
	Labels          map[string]string       `json:",omitempty"`
	Command         []string                `json:",omitempty"`
	Args            []string                `json:",omitempty"`
	Hostname        string                  `json:",omitempty"`
	Env             []string                `json:",omitempty"`
	Dir             string                  `json:",omitempty"`
	User            string                  `json:",omitempty"`
	Groups          []string                `json:",omitempty"`
	Privileges      *Privileges             `json:",omitempty"`
	StopSignal      string                  `json:",omitempty"`
	TTY             bool                    `json:",omitempty"`
	OpenStdin       bool                    `json:",omitempty"`
	ReadOnly        bool                    `json:",omitempty"`
	Mounts          []mount.Mount           `json:",omitempty"`
	StopGracePeriod *time.Duration          `json:",omitempty"`
	Healthcheck     *container.HealthConfig `json:",omitempty"`
	
	
	
	Hosts     []string            `json:",omitempty"`
	DNSConfig *DNSConfig          `json:",omitempty"`
	Secrets   []*SecretReference  `json:",omitempty"`
	Configs   []*ConfigReference  `json:",omitempty"`
	Isolation container.Isolation `json:",omitempty"`
}
