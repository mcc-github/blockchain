package container 

import (
	"time"

	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)





const MinimumDuration = 1 * time.Millisecond


type HealthConfig struct {
	
	
	
	
	
	
	
	Test []string `json:",omitempty"`

	
	Interval    time.Duration `json:",omitempty"` 
	Timeout     time.Duration `json:",omitempty"` 
	StartPeriod time.Duration `json:",omitempty"` 

	
	
	Retries int `json:",omitempty"`
}







type Config struct {
	Hostname        string              
	Domainname      string              
	User            string              
	AttachStdin     bool                
	AttachStdout    bool                
	AttachStderr    bool                
	ExposedPorts    nat.PortSet         `json:",omitempty"` 
	Tty             bool                
	OpenStdin       bool                
	StdinOnce       bool                
	Env             []string            
	Cmd             strslice.StrSlice   
	Healthcheck     *HealthConfig       `json:",omitempty"` 
	ArgsEscaped     bool                `json:",omitempty"` 
	Image           string              
	Volumes         map[string]struct{} 
	WorkingDir      string              
	Entrypoint      strslice.StrSlice   
	NetworkDisabled bool                `json:",omitempty"` 
	MacAddress      string              `json:",omitempty"` 
	OnBuild         []string            
	Labels          map[string]string   
	StopSignal      string              `json:",omitempty"` 
	StopTimeout     *int                `json:",omitempty"` 
	Shell           strslice.StrSlice   `json:",omitempty"` 
}
