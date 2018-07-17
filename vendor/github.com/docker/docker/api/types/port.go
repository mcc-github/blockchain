package types






type Port struct {

	
	IP string `json:"IP,omitempty"`

	
	
	PrivatePort uint16 `json:"PrivatePort"`

	
	PublicPort uint16 `json:"PublicPort,omitempty"`

	
	
	Type string `json:"Type"`
}
