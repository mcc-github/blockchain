













package v1

import "github.com/opencontainers/image-spec/specs-go"



type Index struct {
	specs.Versioned

	
	Manifests []Descriptor `json:"manifests"`

	
	Annotations map[string]string `json:"annotations,omitempty"`
}
