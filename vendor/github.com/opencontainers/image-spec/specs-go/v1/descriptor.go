













package v1

import digest "github.com/opencontainers/go-digest"




type Descriptor struct {
	
	MediaType string `json:"mediaType,omitempty"`

	
	Digest digest.Digest `json:"digest"`

	
	Size int64 `json:"size"`

	
	URLs []string `json:"urls,omitempty"`

	
	Annotations map[string]string `json:"annotations,omitempty"`

	
	
	
	Platform *Platform `json:"platform,omitempty"`
}


type Platform struct {
	
	
	Architecture string `json:"architecture"`

	
	OS string `json:"os"`

	
	
	OSVersion string `json:"os.version,omitempty"`

	
	
	OSFeatures []string `json:"os.features,omitempty"`

	
	
	Variant string `json:"variant,omitempty"`
}
