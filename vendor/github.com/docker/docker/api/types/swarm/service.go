package swarm 

import "time"


type Service struct {
	ID string
	Meta
	Spec         ServiceSpec   `json:",omitempty"`
	PreviousSpec *ServiceSpec  `json:",omitempty"`
	Endpoint     Endpoint      `json:",omitempty"`
	UpdateStatus *UpdateStatus `json:",omitempty"`
}


type ServiceSpec struct {
	Annotations

	
	
	TaskTemplate   TaskSpec      `json:",omitempty"`
	Mode           ServiceMode   `json:",omitempty"`
	UpdateConfig   *UpdateConfig `json:",omitempty"`
	RollbackConfig *UpdateConfig `json:",omitempty"`

	
	
	
	Networks     []NetworkAttachmentConfig `json:",omitempty"`
	EndpointSpec *EndpointSpec             `json:",omitempty"`
}


type ServiceMode struct {
	Replicated *ReplicatedService `json:",omitempty"`
	Global     *GlobalService     `json:",omitempty"`
}


type UpdateState string

const (
	
	UpdateStateUpdating UpdateState = "updating"
	
	UpdateStatePaused UpdateState = "paused"
	
	UpdateStateCompleted UpdateState = "completed"
	
	UpdateStateRollbackStarted UpdateState = "rollback_started"
	
	UpdateStateRollbackPaused UpdateState = "rollback_paused"
	
	UpdateStateRollbackCompleted UpdateState = "rollback_completed"
)


type UpdateStatus struct {
	State       UpdateState `json:",omitempty"`
	StartedAt   *time.Time  `json:",omitempty"`
	CompletedAt *time.Time  `json:",omitempty"`
	Message     string      `json:",omitempty"`
}


type ReplicatedService struct {
	Replicas *uint64 `json:",omitempty"`
}


type GlobalService struct{}

const (
	
	UpdateFailureActionPause = "pause"
	
	UpdateFailureActionContinue = "continue"
	
	UpdateFailureActionRollback = "rollback"

	
	UpdateOrderStopFirst = "stop-first"
	
	UpdateOrderStartFirst = "start-first"
)


type UpdateConfig struct {
	
	
	Parallelism uint64

	
	Delay time.Duration `json:",omitempty"`

	
	FailureAction string `json:",omitempty"`

	
	
	
	
	
	
	Monitor time.Duration `json:",omitempty"`

	
	
	
	
	
	
	
	
	
	
	
	MaxFailureRatio float32

	
	
	
	Order string
}
