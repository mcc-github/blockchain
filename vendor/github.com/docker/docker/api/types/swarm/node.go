package swarm 


type Node struct {
	ID string
	Meta
	
	
	Spec NodeSpec `json:",omitempty"`
	
	
	Description NodeDescription `json:",omitempty"`
	
	Status NodeStatus `json:",omitempty"`
	
	
	ManagerStatus *ManagerStatus `json:",omitempty"`
}


type NodeSpec struct {
	Annotations
	Role         NodeRole         `json:",omitempty"`
	Availability NodeAvailability `json:",omitempty"`
}


type NodeRole string

const (
	
	NodeRoleWorker NodeRole = "worker"
	
	NodeRoleManager NodeRole = "manager"
)


type NodeAvailability string

const (
	
	NodeAvailabilityActive NodeAvailability = "active"
	
	NodeAvailabilityPause NodeAvailability = "pause"
	
	NodeAvailabilityDrain NodeAvailability = "drain"
)


type NodeDescription struct {
	Hostname  string            `json:",omitempty"`
	Platform  Platform          `json:",omitempty"`
	Resources Resources         `json:",omitempty"`
	Engine    EngineDescription `json:",omitempty"`
	TLSInfo   TLSInfo           `json:",omitempty"`
}


type Platform struct {
	Architecture string `json:",omitempty"`
	OS           string `json:",omitempty"`
}


type EngineDescription struct {
	EngineVersion string              `json:",omitempty"`
	Labels        map[string]string   `json:",omitempty"`
	Plugins       []PluginDescription `json:",omitempty"`
}


type PluginDescription struct {
	Type string `json:",omitempty"`
	Name string `json:",omitempty"`
}


type NodeStatus struct {
	State   NodeState `json:",omitempty"`
	Message string    `json:",omitempty"`
	Addr    string    `json:",omitempty"`
}


type Reachability string

const (
	
	ReachabilityUnknown Reachability = "unknown"
	
	ReachabilityUnreachable Reachability = "unreachable"
	
	ReachabilityReachable Reachability = "reachable"
)


type ManagerStatus struct {
	Leader       bool         `json:",omitempty"`
	Reachability Reachability `json:",omitempty"`
	Addr         string       `json:",omitempty"`
}


type NodeState string

const (
	
	NodeStateUnknown NodeState = "unknown"
	
	NodeStateDown NodeState = "down"
	
	NodeStateReady NodeState = "ready"
	
	NodeStateDisconnected NodeState = "disconnected"
)
