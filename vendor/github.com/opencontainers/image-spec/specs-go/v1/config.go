













package v1

import (
	"time"

	digest "github.com/opencontainers/go-digest"
)


type ImageConfig struct {
	
	User string `json:"User,omitempty"`

	
	ExposedPorts map[string]struct{} `json:"ExposedPorts,omitempty"`

	
	Env []string `json:"Env,omitempty"`

	
	Entrypoint []string `json:"Entrypoint,omitempty"`

	
	Cmd []string `json:"Cmd,omitempty"`

	
	Volumes map[string]struct{} `json:"Volumes,omitempty"`

	
	WorkingDir string `json:"WorkingDir,omitempty"`

	
	Labels map[string]string `json:"Labels,omitempty"`

	
	StopSignal string `json:"StopSignal,omitempty"`
}


type RootFS struct {
	
	Type string `json:"type"`

	
	DiffIDs []digest.Digest `json:"diff_ids"`
}


type History struct {
	
	Created *time.Time `json:"created,omitempty"`

	
	CreatedBy string `json:"created_by,omitempty"`

	
	Author string `json:"author,omitempty"`

	
	Comment string `json:"comment,omitempty"`

	
	EmptyLayer bool `json:"empty_layer,omitempty"`
}



type Image struct {
	
	Created *time.Time `json:"created,omitempty"`

	
	Author string `json:"author,omitempty"`

	
	Architecture string `json:"architecture"`

	
	OS string `json:"os"`

	
	Config ImageConfig `json:"config,omitempty"`

	
	RootFS RootFS `json:"rootfs"`

	
	History []History `json:"history,omitempty"`
}
