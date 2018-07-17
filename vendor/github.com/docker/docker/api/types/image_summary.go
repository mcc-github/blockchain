package types






type ImageSummary struct {

	
	
	Containers int64 `json:"Containers"`

	
	
	Created int64 `json:"Created"`

	
	
	ID string `json:"Id"`

	
	
	Labels map[string]string `json:"Labels"`

	
	
	ParentID string `json:"ParentId"`

	
	
	RepoDigests []string `json:"RepoDigests"`

	
	
	RepoTags []string `json:"RepoTags"`

	
	
	SharedSize int64 `json:"SharedSize"`

	
	
	Size int64 `json:"Size"`

	
	
	VirtualSize int64 `json:"VirtualSize"`
}
