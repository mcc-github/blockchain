package swarm 

import (
	"time"
)



type ClusterInfo struct {
	ID string
	Meta
	Spec                   Spec
	TLSInfo                TLSInfo
	RootRotationInProgress bool
	DefaultAddrPool        []string
	SubnetSize             uint32
}


type Swarm struct {
	ClusterInfo
	JoinTokens JoinTokens
}


type JoinTokens struct {
	
	Worker string
	
	Manager string
}


type Spec struct {
	Annotations

	Orchestration    OrchestrationConfig `json:",omitempty"`
	Raft             RaftConfig          `json:",omitempty"`
	Dispatcher       DispatcherConfig    `json:",omitempty"`
	CAConfig         CAConfig            `json:",omitempty"`
	TaskDefaults     TaskDefaults        `json:",omitempty"`
	EncryptionConfig EncryptionConfig    `json:",omitempty"`
}


type OrchestrationConfig struct {
	
	
	TaskHistoryRetentionLimit *int64 `json:",omitempty"`
}


type TaskDefaults struct {
	
	
	
	
	
	
	LogDriver *Driver `json:",omitempty"`
}


type EncryptionConfig struct {
	
	
	
	AutoLockManagers bool
}


type RaftConfig struct {
	
	SnapshotInterval uint64 `json:",omitempty"`

	
	
	KeepOldSnapshots *uint64 `json:",omitempty"`

	
	
	LogEntriesForSlowFollowers uint64 `json:",omitempty"`

	
	
	
	
	
	
	ElectionTick int

	
	
	
	
	
	
	HeartbeatTick int
}


type DispatcherConfig struct {
	
	
	HeartbeatPeriod time.Duration `json:",omitempty"`
}


type CAConfig struct {
	
	NodeCertExpiry time.Duration `json:",omitempty"`

	
	
	ExternalCAs []*ExternalCA `json:",omitempty"`

	
	
	
	SigningCACert string `json:",omitempty"`
	SigningCAKey  string `json:",omitempty"`

	
	
	ForceRotate uint64 `json:",omitempty"`
}


type ExternalCAProtocol string


const ExternalCAProtocolCFSSL ExternalCAProtocol = "cfssl"


type ExternalCA struct {
	
	Protocol ExternalCAProtocol

	
	URL string

	
	
	Options map[string]string `json:",omitempty"`

	
	
	CACert string
}


type InitRequest struct {
	ListenAddr       string
	AdvertiseAddr    string
	DataPathAddr     string
	ForceNewCluster  bool
	Spec             Spec
	AutoLockManagers bool
	Availability     NodeAvailability
	DefaultAddrPool  []string
	SubnetSize       uint32
}


type JoinRequest struct {
	ListenAddr    string
	AdvertiseAddr string
	DataPathAddr  string
	RemoteAddrs   []string
	JoinToken     string 
	Availability  NodeAvailability
}


type UnlockRequest struct {
	
	UnlockKey string
}


type LocalNodeState string

const (
	
	LocalNodeStateInactive LocalNodeState = "inactive"
	
	LocalNodeStatePending LocalNodeState = "pending"
	
	LocalNodeStateActive LocalNodeState = "active"
	
	LocalNodeStateError LocalNodeState = "error"
	
	LocalNodeStateLocked LocalNodeState = "locked"
)


type Info struct {
	NodeID   string
	NodeAddr string

	LocalNodeState   LocalNodeState
	ControlAvailable bool
	Error            string

	RemoteManagers []Peer
	Nodes          int `json:",omitempty"`
	Managers       int `json:",omitempty"`

	Cluster *ClusterInfo `json:",omitempty"`
}


type Peer struct {
	NodeID string
	Addr   string
}


type UpdateFlags struct {
	RotateWorkerToken      bool
	RotateManagerToken     bool
	RotateManagerUnlockKey bool
}
