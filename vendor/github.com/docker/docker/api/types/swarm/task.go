package swarm 

import (
	"time"

	"github.com/docker/docker/api/types/swarm/runtime"
)


type TaskState string

const (
	
	TaskStateNew TaskState = "new"
	
	TaskStateAllocated TaskState = "allocated"
	
	TaskStatePending TaskState = "pending"
	
	TaskStateAssigned TaskState = "assigned"
	
	TaskStateAccepted TaskState = "accepted"
	
	TaskStatePreparing TaskState = "preparing"
	
	TaskStateReady TaskState = "ready"
	
	TaskStateStarting TaskState = "starting"
	
	TaskStateRunning TaskState = "running"
	
	TaskStateComplete TaskState = "complete"
	
	TaskStateShutdown TaskState = "shutdown"
	
	TaskStateFailed TaskState = "failed"
	
	TaskStateRejected TaskState = "rejected"
	
	TaskStateRemove TaskState = "remove"
	
	TaskStateOrphaned TaskState = "orphaned"
)


type Task struct {
	ID string
	Meta
	Annotations

	Spec                TaskSpec            `json:",omitempty"`
	ServiceID           string              `json:",omitempty"`
	Slot                int                 `json:",omitempty"`
	NodeID              string              `json:",omitempty"`
	Status              TaskStatus          `json:",omitempty"`
	DesiredState        TaskState           `json:",omitempty"`
	NetworksAttachments []NetworkAttachment `json:",omitempty"`
	GenericResources    []GenericResource   `json:",omitempty"`
}


type TaskSpec struct {
	
	
	ContainerSpec *ContainerSpec      `json:",omitempty"`
	PluginSpec    *runtime.PluginSpec `json:",omitempty"`

	Resources     *ResourceRequirements     `json:",omitempty"`
	RestartPolicy *RestartPolicy            `json:",omitempty"`
	Placement     *Placement                `json:",omitempty"`
	Networks      []NetworkAttachmentConfig `json:",omitempty"`

	
	
	
	LogDriver *Driver `json:",omitempty"`

	
	
	ForceUpdate uint64

	Runtime RuntimeType `json:",omitempty"`
}


type Resources struct {
	NanoCPUs         int64             `json:",omitempty"`
	MemoryBytes      int64             `json:",omitempty"`
	GenericResources []GenericResource `json:",omitempty"`
}



type GenericResource struct {
	NamedResourceSpec    *NamedGenericResource    `json:",omitempty"`
	DiscreteResourceSpec *DiscreteGenericResource `json:",omitempty"`
}





type NamedGenericResource struct {
	Kind  string `json:",omitempty"`
	Value string `json:",omitempty"`
}





type DiscreteGenericResource struct {
	Kind  string `json:",omitempty"`
	Value int64  `json:",omitempty"`
}


type ResourceRequirements struct {
	Limits       *Resources `json:",omitempty"`
	Reservations *Resources `json:",omitempty"`
}


type Placement struct {
	Constraints []string              `json:",omitempty"`
	Preferences []PlacementPreference `json:",omitempty"`

	
	
	
	Platforms []Platform `json:",omitempty"`
}



type PlacementPreference struct {
	Spread *SpreadOver
}



type SpreadOver struct {
	
	SpreadDescriptor string
}


type RestartPolicy struct {
	Condition   RestartPolicyCondition `json:",omitempty"`
	Delay       *time.Duration         `json:",omitempty"`
	MaxAttempts *uint64                `json:",omitempty"`
	Window      *time.Duration         `json:",omitempty"`
}


type RestartPolicyCondition string

const (
	
	RestartPolicyConditionNone RestartPolicyCondition = "none"
	
	RestartPolicyConditionOnFailure RestartPolicyCondition = "on-failure"
	
	RestartPolicyConditionAny RestartPolicyCondition = "any"
)


type TaskStatus struct {
	Timestamp       time.Time        `json:",omitempty"`
	State           TaskState        `json:",omitempty"`
	Message         string           `json:",omitempty"`
	Err             string           `json:",omitempty"`
	ContainerStatus *ContainerStatus `json:",omitempty"`
	PortStatus      PortStatus       `json:",omitempty"`
}


type ContainerStatus struct {
	ContainerID string
	PID         int
	ExitCode    int
}



type PortStatus struct {
	Ports []PortConfig `json:",omitempty"`
}
