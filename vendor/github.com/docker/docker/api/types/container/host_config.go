package container 

import (
	"strings"

	"github.com/docker/docker/api/types/blkiodev"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
)



type Isolation string



func (i Isolation) IsDefault() bool {
	return strings.ToLower(string(i)) == "default" || string(i) == ""
}


func (i Isolation) IsHyperV() bool {
	return strings.ToLower(string(i)) == "hyperv"
}


func (i Isolation) IsProcess() bool {
	return strings.ToLower(string(i)) == "process"
}

const (
	
	IsolationEmpty = Isolation("")
	
	IsolationDefault = Isolation("default")
	
	IsolationProcess = Isolation("process")
	
	IsolationHyperV = Isolation("hyperv")
)


type IpcMode string


func (n IpcMode) IsPrivate() bool {
	return n == "private"
}


func (n IpcMode) IsHost() bool {
	return n == "host"
}


func (n IpcMode) IsShareable() bool {
	return n == "shareable"
}


func (n IpcMode) IsContainer() bool {
	parts := strings.SplitN(string(n), ":", 2)
	return len(parts) > 1 && parts[0] == "container"
}


func (n IpcMode) IsNone() bool {
	return n == "none"
}


func (n IpcMode) IsEmpty() bool {
	return n == ""
}


func (n IpcMode) Valid() bool {
	return n.IsEmpty() || n.IsNone() || n.IsPrivate() || n.IsHost() || n.IsShareable() || n.IsContainer()
}


func (n IpcMode) Container() string {
	parts := strings.SplitN(string(n), ":", 2)
	if len(parts) > 1 && parts[0] == "container" {
		return parts[1]
	}
	return ""
}


type NetworkMode string


func (n NetworkMode) IsNone() bool {
	return n == "none"
}


func (n NetworkMode) IsDefault() bool {
	return n == "default"
}


func (n NetworkMode) IsPrivate() bool {
	return !(n.IsHost() || n.IsContainer())
}


func (n NetworkMode) IsContainer() bool {
	parts := strings.SplitN(string(n), ":", 2)
	return len(parts) > 1 && parts[0] == "container"
}


func (n NetworkMode) ConnectedContainer() string {
	parts := strings.SplitN(string(n), ":", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}


func (n NetworkMode) UserDefined() string {
	if n.IsUserDefined() {
		return string(n)
	}
	return ""
}


type UsernsMode string


func (n UsernsMode) IsHost() bool {
	return n == "host"
}


func (n UsernsMode) IsPrivate() bool {
	return !(n.IsHost())
}


func (n UsernsMode) Valid() bool {
	parts := strings.Split(string(n), ":")
	switch mode := parts[0]; mode {
	case "", "host":
	default:
		return false
	}
	return true
}


type CgroupSpec string


func (c CgroupSpec) IsContainer() bool {
	parts := strings.SplitN(string(c), ":", 2)
	return len(parts) > 1 && parts[0] == "container"
}


func (c CgroupSpec) Valid() bool {
	return c.IsContainer() || c == ""
}


func (c CgroupSpec) Container() string {
	parts := strings.SplitN(string(c), ":", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}


type UTSMode string


func (n UTSMode) IsPrivate() bool {
	return !(n.IsHost())
}


func (n UTSMode) IsHost() bool {
	return n == "host"
}


func (n UTSMode) Valid() bool {
	parts := strings.Split(string(n), ":")
	switch mode := parts[0]; mode {
	case "", "host":
	default:
		return false
	}
	return true
}


type PidMode string


func (n PidMode) IsPrivate() bool {
	return !(n.IsHost() || n.IsContainer())
}


func (n PidMode) IsHost() bool {
	return n == "host"
}


func (n PidMode) IsContainer() bool {
	parts := strings.SplitN(string(n), ":", 2)
	return len(parts) > 1 && parts[0] == "container"
}


func (n PidMode) Valid() bool {
	parts := strings.Split(string(n), ":")
	switch mode := parts[0]; mode {
	case "", "host":
	case "container":
		if len(parts) != 2 || parts[1] == "" {
			return false
		}
	default:
		return false
	}
	return true
}


func (n PidMode) Container() string {
	parts := strings.SplitN(string(n), ":", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}


type DeviceMapping struct {
	PathOnHost        string
	PathInContainer   string
	CgroupPermissions string
}


type RestartPolicy struct {
	Name              string
	MaximumRetryCount int
}



func (rp *RestartPolicy) IsNone() bool {
	return rp.Name == "no" || rp.Name == ""
}



func (rp *RestartPolicy) IsAlways() bool {
	return rp.Name == "always"
}



func (rp *RestartPolicy) IsOnFailure() bool {
	return rp.Name == "on-failure"
}




func (rp *RestartPolicy) IsUnlessStopped() bool {
	return rp.Name == "unless-stopped"
}


func (rp *RestartPolicy) IsSame(tp *RestartPolicy) bool {
	return rp.Name == tp.Name && rp.MaximumRetryCount == tp.MaximumRetryCount
}



type LogMode string


const (
	LogModeUnset            = ""
	LogModeBlocking LogMode = "blocking"
	LogModeNonBlock LogMode = "non-blocking"
)


type LogConfig struct {
	Type   string
	Config map[string]string
}


type Resources struct {
	
	CPUShares int64 `json:"CpuShares"` 
	Memory    int64 
	NanoCPUs  int64 `json:"NanoCpus"` 

	
	CgroupParent         string 
	BlkioWeight          uint16 
	BlkioWeightDevice    []*blkiodev.WeightDevice
	BlkioDeviceReadBps   []*blkiodev.ThrottleDevice
	BlkioDeviceWriteBps  []*blkiodev.ThrottleDevice
	BlkioDeviceReadIOps  []*blkiodev.ThrottleDevice
	BlkioDeviceWriteIOps []*blkiodev.ThrottleDevice
	CPUPeriod            int64           `json:"CpuPeriod"`          
	CPUQuota             int64           `json:"CpuQuota"`           
	CPURealtimePeriod    int64           `json:"CpuRealtimePeriod"`  
	CPURealtimeRuntime   int64           `json:"CpuRealtimeRuntime"` 
	CpusetCpus           string          
	CpusetMems           string          
	Devices              []DeviceMapping 
	DeviceCgroupRules    []string        
	DiskQuota            int64           
	KernelMemory         int64           
	MemoryReservation    int64           
	MemorySwap           int64           
	MemorySwappiness     *int64          
	OomKillDisable       *bool           
	PidsLimit            int64           
	Ulimits              []*units.Ulimit 

	
	CPUCount           int64  `json:"CpuCount"`   
	CPUPercent         int64  `json:"CpuPercent"` 
	IOMaximumIOps      uint64 
	IOMaximumBandwidth uint64 
}



type UpdateConfig struct {
	
	Resources
	RestartPolicy RestartPolicy
}




type HostConfig struct {
	
	Binds           []string      
	ContainerIDFile string        
	LogConfig       LogConfig     
	NetworkMode     NetworkMode   
	PortBindings    nat.PortMap   
	RestartPolicy   RestartPolicy 
	AutoRemove      bool          
	VolumeDriver    string        
	VolumesFrom     []string      

	
	CapAdd          strslice.StrSlice 
	CapDrop         strslice.StrSlice 
	DNS             []string          `json:"Dns"`        
	DNSOptions      []string          `json:"DnsOptions"` 
	DNSSearch       []string          `json:"DnsSearch"`  
	ExtraHosts      []string          
	GroupAdd        []string          
	IpcMode         IpcMode           
	Cgroup          CgroupSpec        
	Links           []string          
	OomScoreAdj     int               
	PidMode         PidMode           
	Privileged      bool              
	PublishAllPorts bool              
	ReadonlyRootfs  bool              
	SecurityOpt     []string          
	StorageOpt      map[string]string `json:",omitempty"` 
	Tmpfs           map[string]string `json:",omitempty"` 
	UTSMode         UTSMode           
	UsernsMode      UsernsMode        
	ShmSize         int64             
	Sysctls         map[string]string `json:",omitempty"` 
	Runtime         string            `json:",omitempty"` 

	
	ConsoleSize [2]uint   
	Isolation   Isolation 

	
	Resources

	
	Mounts []mount.Mount `json:",omitempty"`

	
	Init *bool `json:",omitempty"`
}
