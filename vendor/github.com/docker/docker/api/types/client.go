package types 

import (
	"bufio"
	"io"
	"net"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	units "github.com/docker/go-units"
)


type CheckpointCreateOptions struct {
	CheckpointID  string
	CheckpointDir string
	Exit          bool
}


type CheckpointListOptions struct {
	CheckpointDir string
}


type CheckpointDeleteOptions struct {
	CheckpointID  string
	CheckpointDir string
}


type ContainerAttachOptions struct {
	Stream     bool
	Stdin      bool
	Stdout     bool
	Stderr     bool
	DetachKeys string
	Logs       bool
}


type ContainerCommitOptions struct {
	Reference string
	Comment   string
	Author    string
	Changes   []string
	Pause     bool
	Config    *container.Config
}


type ContainerExecInspect struct {
	ExecID      string
	ContainerID string
	Running     bool
	ExitCode    int
	Pid         int
}


type ContainerListOptions struct {
	Quiet   bool
	Size    bool
	All     bool
	Latest  bool
	Since   string
	Before  string
	Limit   int
	Filters filters.Args
}


type ContainerLogsOptions struct {
	ShowStdout bool
	ShowStderr bool
	Since      string
	Until      string
	Timestamps bool
	Follow     bool
	Tail       string
	Details    bool
}


type ContainerRemoveOptions struct {
	RemoveVolumes bool
	RemoveLinks   bool
	Force         bool
}


type ContainerStartOptions struct {
	CheckpointID  string
	CheckpointDir string
}



type CopyToContainerOptions struct {
	AllowOverwriteDirWithFile bool
	CopyUIDGID                bool
}


type EventsOptions struct {
	Since   string
	Until   string
	Filters filters.Args
}


type NetworkListOptions struct {
	Filters filters.Args
}


type HijackedResponse struct {
	Conn   net.Conn
	Reader *bufio.Reader
}


func (h *HijackedResponse) Close() {
	h.Conn.Close()
}



type CloseWriter interface {
	CloseWrite() error
}


func (h *HijackedResponse) CloseWrite() error {
	if conn, ok := h.Conn.(CloseWriter); ok {
		return conn.CloseWrite()
	}
	return nil
}



type ImageBuildOptions struct {
	Tags           []string
	SuppressOutput bool
	RemoteContext  string
	NoCache        bool
	Remove         bool
	ForceRemove    bool
	PullParent     bool
	Isolation      container.Isolation
	CPUSetCPUs     string
	CPUSetMems     string
	CPUShares      int64
	CPUQuota       int64
	CPUPeriod      int64
	Memory         int64
	MemorySwap     int64
	CgroupParent   string
	NetworkMode    string
	ShmSize        int64
	Dockerfile     string
	Ulimits        []*units.Ulimit
	
	
	
	
	BuildArgs   map[string]*string
	AuthConfigs map[string]AuthConfig
	Context     io.Reader
	Labels      map[string]string
	
	
	
	Squash bool
	
	
	CacheFrom   []string
	SecurityOpt []string
	ExtraHosts  []string 
	Target      string
	SessionID   string
	Platform    string
	
	Version BuilderVersion
	
	
	
	BuildID string
}


type BuilderVersion string

const (
	
	BuilderV1 BuilderVersion = "1"
	
	BuilderBuildKit = "2"
)




type ImageBuildResponse struct {
	Body   io.ReadCloser
	OSType string
}


type ImageCreateOptions struct {
	RegistryAuth string 
	Platform     string 
}


type ImageImportSource struct {
	Source     io.Reader 
	SourceName string    
}


type ImageImportOptions struct {
	Tag      string   
	Message  string   
	Changes  []string 
	Platform string   
}


type ImageListOptions struct {
	All     bool
	Filters filters.Args
}


type ImageLoadResponse struct {
	
	Body io.ReadCloser
	JSON bool
}


type ImagePullOptions struct {
	All           bool
	RegistryAuth  string 
	PrivilegeFunc RequestPrivilegeFunc
	Platform      string
}







type RequestPrivilegeFunc func() (string, error)


type ImagePushOptions ImagePullOptions


type ImageRemoveOptions struct {
	Force         bool
	PruneChildren bool
}


type ImageSearchOptions struct {
	RegistryAuth  string
	PrivilegeFunc RequestPrivilegeFunc
	Filters       filters.Args
	Limit         int
}




type ResizeOptions struct {
	Height uint
	Width  uint
}


type NodeListOptions struct {
	Filters filters.Args
}


type NodeRemoveOptions struct {
	Force bool
}


type ServiceCreateOptions struct {
	
	
	
	
	EncodedRegistryAuth string

	
	
	
	
	QueryRegistry bool
}



type ServiceCreateResponse struct {
	
	ID string
	
	Warnings []string `json:",omitempty"`
}


const (
	RegistryAuthFromSpec         = "spec"
	RegistryAuthFromPreviousSpec = "previous-spec"
)


type ServiceUpdateOptions struct {
	
	
	
	
	EncodedRegistryAuth string

	
	
	

	
	
	
	RegistryAuthFrom string

	
	
	
	
	Rollback string

	
	
	
	
	QueryRegistry bool
}


type ServiceListOptions struct {
	Filters filters.Args
}



type ServiceInspectOptions struct {
	InsertDefaults bool
}


type TaskListOptions struct {
	Filters filters.Args
}


type PluginRemoveOptions struct {
	Force bool
}


type PluginEnableOptions struct {
	Timeout int
}


type PluginDisableOptions struct {
	Force bool
}


type PluginInstallOptions struct {
	Disabled              bool
	AcceptAllPermissions  bool
	RegistryAuth          string 
	RemoteRef             string 
	PrivilegeFunc         RequestPrivilegeFunc
	AcceptPermissionsFunc func(PluginPrivileges) (bool, error)
	Args                  []string
}



type SwarmUnlockKeyResponse struct {
	
	UnlockKey string
}


type PluginCreateOptions struct {
	RepoName string
}
