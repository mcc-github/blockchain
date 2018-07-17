package types 

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/go-connections/nat"
)


type RootFS struct {
	Type      string
	Layers    []string `json:",omitempty"`
	BaseLayer string   `json:",omitempty"`
}



type ImageInspect struct {
	ID              string `json:"Id"`
	RepoTags        []string
	RepoDigests     []string
	Parent          string
	Comment         string
	Created         string
	Container       string
	ContainerConfig *container.Config
	DockerVersion   string
	Author          string
	Config          *container.Config
	Architecture    string
	Os              string
	OsVersion       string `json:",omitempty"`
	Size            int64
	VirtualSize     int64
	GraphDriver     GraphDriverData
	RootFS          RootFS
	Metadata        ImageMetadata
}


type ImageMetadata struct {
	LastTagTime time.Time `json:",omitempty"`
}



type Container struct {
	ID         string `json:"Id"`
	Names      []string
	Image      string
	ImageID    string
	Command    string
	Created    int64
	Ports      []Port
	SizeRw     int64 `json:",omitempty"`
	SizeRootFs int64 `json:",omitempty"`
	Labels     map[string]string
	State      string
	Status     string
	HostConfig struct {
		NetworkMode string `json:",omitempty"`
	}
	NetworkSettings *SummaryNetworkSettings
	Mounts          []MountPoint
}



type CopyConfig struct {
	Resource string
}




type ContainerPathStat struct {
	Name       string      `json:"name"`
	Size       int64       `json:"size"`
	Mode       os.FileMode `json:"mode"`
	Mtime      time.Time   `json:"mtime"`
	LinkTarget string      `json:"linkTarget"`
}



type ContainerStats struct {
	Body   io.ReadCloser `json:"body"`
	OSType string        `json:"ostype"`
}



type Ping struct {
	APIVersion   string
	OSType       string
	Experimental bool
}


type ComponentVersion struct {
	Name    string
	Version string
	Details map[string]string `json:",omitempty"`
}



type Version struct {
	Platform   struct{ Name string } `json:",omitempty"`
	Components []ComponentVersion    `json:",omitempty"`

	

	Version       string
	APIVersion    string `json:"ApiVersion"`
	MinAPIVersion string `json:"MinAPIVersion,omitempty"`
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
	KernelVersion string `json:",omitempty"`
	Experimental  bool   `json:",omitempty"`
	BuildTime     string `json:",omitempty"`
}



type Commit struct {
	ID       string 
	Expected string 
}



type Info struct {
	ID                 string
	Containers         int
	ContainersRunning  int
	ContainersPaused   int
	ContainersStopped  int
	Images             int
	Driver             string
	DriverStatus       [][2]string
	SystemStatus       [][2]string
	Plugins            PluginsInfo
	MemoryLimit        bool
	SwapLimit          bool
	KernelMemory       bool
	CPUCfsPeriod       bool `json:"CpuCfsPeriod"`
	CPUCfsQuota        bool `json:"CpuCfsQuota"`
	CPUShares          bool
	CPUSet             bool
	IPv4Forwarding     bool
	BridgeNfIptables   bool
	BridgeNfIP6tables  bool `json:"BridgeNfIp6tables"`
	Debug              bool
	NFd                int
	OomKillDisable     bool
	NGoroutines        int
	SystemTime         string
	LoggingDriver      string
	CgroupDriver       string
	NEventsListener    int
	KernelVersion      string
	OperatingSystem    string
	OSType             string
	Architecture       string
	IndexServerAddress string
	RegistryConfig     *registry.ServiceConfig
	NCPU               int
	MemTotal           int64
	GenericResources   []swarm.GenericResource
	DockerRootDir      string
	HTTPProxy          string `json:"HttpProxy"`
	HTTPSProxy         string `json:"HttpsProxy"`
	NoProxy            string
	Name               string
	Labels             []string
	ExperimentalBuild  bool
	ServerVersion      string
	ClusterStore       string
	ClusterAdvertise   string
	Runtimes           map[string]Runtime
	DefaultRuntime     string
	Swarm              swarm.Info
	
	
	
	LiveRestoreEnabled bool
	Isolation          container.Isolation
	InitBinary         string
	ContainerdCommit   Commit
	RuncCommit         Commit
	InitCommit         Commit
	SecurityOptions    []string
}


type KeyValue struct {
	Key, Value string
}


type SecurityOpt struct {
	Name    string
	Options []KeyValue
}



func DecodeSecurityOptions(opts []string) ([]SecurityOpt, error) {
	so := []SecurityOpt{}
	for _, opt := range opts {
		
		if !strings.Contains(opt, "=") {
			so = append(so, SecurityOpt{Name: opt})
			continue
		}
		secopt := SecurityOpt{}
		split := strings.Split(opt, ",")
		for _, s := range split {
			kv := strings.SplitN(s, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid security option %q", s)
			}
			if kv[0] == "" || kv[1] == "" {
				return nil, errors.New("invalid empty security option")
			}
			if kv[0] == "name" {
				secopt.Name = kv[1]
				continue
			}
			secopt.Options = append(secopt.Options, KeyValue{Key: kv[0], Value: kv[1]})
		}
		so = append(so, secopt)
	}
	return so, nil
}



type PluginsInfo struct {
	
	Volume []string
	
	Network []string
	
	Authorization []string
	
	Log []string
}



type ExecStartCheck struct {
	
	Detach bool
	
	Tty bool
}


type HealthcheckResult struct {
	Start    time.Time 
	End      time.Time 
	ExitCode int       
	Output   string    
}


const (
	NoHealthcheck = "none"      
	Starting      = "starting"  
	Healthy       = "healthy"   
	Unhealthy     = "unhealthy" 
)


type Health struct {
	Status        string               
	FailingStreak int                  
	Log           []*HealthcheckResult 
}



type ContainerState struct {
	Status     string 
	Running    bool
	Paused     bool
	Restarting bool
	OOMKilled  bool
	Dead       bool
	Pid        int
	ExitCode   int
	Error      string
	StartedAt  string
	FinishedAt string
	Health     *Health `json:",omitempty"`
}



type ContainerNode struct {
	ID        string
	IPAddress string `json:"IP"`
	Addr      string
	Name      string
	Cpus      int
	Memory    int64
	Labels    map[string]string
}



type ContainerJSONBase struct {
	ID              string `json:"Id"`
	Created         string
	Path            string
	Args            []string
	State           *ContainerState
	Image           string
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	LogPath         string
	Node            *ContainerNode `json:",omitempty"`
	Name            string
	RestartCount    int
	Driver          string
	Platform        string
	MountLabel      string
	ProcessLabel    string
	AppArmorProfile string
	ExecIDs         []string
	HostConfig      *container.HostConfig
	GraphDriver     GraphDriverData
	SizeRw          *int64 `json:",omitempty"`
	SizeRootFs      *int64 `json:",omitempty"`
}


type ContainerJSON struct {
	*ContainerJSONBase
	Mounts          []MountPoint
	Config          *container.Config
	NetworkSettings *NetworkSettings
}


type NetworkSettings struct {
	NetworkSettingsBase
	DefaultNetworkSettings
	Networks map[string]*network.EndpointSettings
}



type SummaryNetworkSettings struct {
	Networks map[string]*network.EndpointSettings
}


type NetworkSettingsBase struct {
	Bridge                 string      
	SandboxID              string      
	HairpinMode            bool        
	LinkLocalIPv6Address   string      
	LinkLocalIPv6PrefixLen int         
	Ports                  nat.PortMap 
	SandboxKey             string      
	SecondaryIPAddresses   []network.Address
	SecondaryIPv6Addresses []network.Address
}




type DefaultNetworkSettings struct {
	EndpointID          string 
	Gateway             string 
	GlobalIPv6Address   string 
	GlobalIPv6PrefixLen int    
	IPAddress           string 
	IPPrefixLen         int    
	IPv6Gateway         string 
	MacAddress          string 
}



type MountPoint struct {
	Type        mount.Type `json:",omitempty"`
	Name        string     `json:",omitempty"`
	Source      string
	Destination string
	Driver      string `json:",omitempty"`
	Mode        string
	RW          bool
	Propagation mount.Propagation
}


type NetworkResource struct {
	Name       string                         
	ID         string                         `json:"Id"` 
	Created    time.Time                      
	Scope      string                         
	Driver     string                         
	EnableIPv6 bool                           
	IPAM       network.IPAM                   
	Internal   bool                           
	Attachable bool                           
	Ingress    bool                           
	ConfigFrom network.ConfigReference        
	ConfigOnly bool                           
	Containers map[string]EndpointResource    
	Options    map[string]string              
	Labels     map[string]string              
	Peers      []network.PeerInfo             `json:",omitempty"` 
	Services   map[string]network.ServiceInfo `json:",omitempty"`
}


type EndpointResource struct {
	Name        string
	EndpointID  string
	MacAddress  string
	IPv4Address string
	IPv6Address string
}


type NetworkCreate struct {
	
	
	
	
	
	
	
	CheckDuplicate bool
	Driver         string
	Scope          string
	EnableIPv6     bool
	IPAM           *network.IPAM
	Internal       bool
	Attachable     bool
	Ingress        bool
	ConfigOnly     bool
	ConfigFrom     *network.ConfigReference
	Options        map[string]string
	Labels         map[string]string
}


type NetworkCreateRequest struct {
	NetworkCreate
	Name string
}


type NetworkCreateResponse struct {
	ID      string `json:"Id"`
	Warning string
}


type NetworkConnect struct {
	Container      string
	EndpointConfig *network.EndpointSettings `json:",omitempty"`
}


type NetworkDisconnect struct {
	Container string
	Force     bool
}


type NetworkInspectOptions struct {
	Scope   string
	Verbose bool
}


type Checkpoint struct {
	Name string 
}


type Runtime struct {
	Path string   `json:"path"`
	Args []string `json:"runtimeArgs,omitempty"`
}



type DiskUsage struct {
	LayersSize  int64
	Images      []*ImageSummary
	Containers  []*Container
	Volumes     []*Volume
	BuilderSize int64
}



type ContainersPruneReport struct {
	ContainersDeleted []string
	SpaceReclaimed    uint64
}



type VolumesPruneReport struct {
	VolumesDeleted []string
	SpaceReclaimed uint64
}



type ImagesPruneReport struct {
	ImagesDeleted  []ImageDeleteResponseItem
	SpaceReclaimed uint64
}



type BuildCachePruneReport struct {
	SpaceReclaimed uint64
}



type NetworksPruneReport struct {
	NetworksDeleted []string
}



type SecretCreateResponse struct {
	
	ID string
}


type SecretListOptions struct {
	Filters filters.Args
}



type ConfigCreateResponse struct {
	
	ID string
}


type ConfigListOptions struct {
	Filters filters.Args
}




type PushResult struct {
	Tag    string
	Digest string
	Size   int
}


type BuildResult struct {
	ID string
}
