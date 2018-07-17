













package v1

import "github.com/opencontainers/image-spec/specs-go"


type Manifest struct {
	specs.Versioned

	
	
	Config Descriptor `json:"config"`

	
	Layers []Descriptor `json:"layers"`

	
	Annotations map[string]string `json:"annotations,omitempty"`
}
