package swarm 

import "time"


type Version struct {
	Index uint64 `json:",omitempty"`
}


type Meta struct {
	Version   Version   `json:",omitempty"`
	CreatedAt time.Time `json:",omitempty"`
	UpdatedAt time.Time `json:",omitempty"`
}


type Annotations struct {
	Name   string            `json:",omitempty"`
	Labels map[string]string `json:"Labels"`
}


type Driver struct {
	Name    string            `json:",omitempty"`
	Options map[string]string `json:",omitempty"`
}



type TLSInfo struct {
	
	TrustRoot string `json:",omitempty"`

	
	CertIssuerSubject []byte `json:",omitempty"`

	
	CertIssuerPublicKey []byte `json:",omitempty"`
}
