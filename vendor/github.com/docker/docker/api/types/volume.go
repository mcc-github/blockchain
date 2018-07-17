package types






type Volume struct {

	
	CreatedAt string `json:"CreatedAt,omitempty"`

	
	
	Driver string `json:"Driver"`

	
	
	Labels map[string]string `json:"Labels"`

	
	
	Mountpoint string `json:"Mountpoint"`

	
	
	Name string `json:"Name"`

	
	
	Options map[string]string `json:"Options"`

	
	
	Scope string `json:"Scope"`

	
	
	
	
	
	
	
	Status map[string]interface{} `json:"Status,omitempty"`

	
	UsageData *VolumeUsageData `json:"UsageData,omitempty"`
}





type VolumeUsageData struct {

	
	
	
	
	RefCount int64 `json:"RefCount"`

	
	
	
	
	
	
	Size int64 `json:"Size"`
}
