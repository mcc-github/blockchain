package types






type Plugin struct {

	
	
	Config PluginConfig `json:"Config"`

	
	
	Enabled bool `json:"Enabled"`

	
	ID string `json:"Id,omitempty"`

	
	
	Name string `json:"Name"`

	
	PluginReference string `json:"PluginReference,omitempty"`

	
	
	Settings PluginSettings `json:"Settings"`
}



type PluginConfig struct {

	
	
	Args PluginConfigArgs `json:"Args"`

	
	
	Description string `json:"Description"`

	
	DockerVersion string `json:"DockerVersion,omitempty"`

	
	
	Documentation string `json:"Documentation"`

	
	
	Entrypoint []string `json:"Entrypoint"`

	
	
	Env []PluginEnv `json:"Env"`

	
	
	Interface PluginConfigInterface `json:"Interface"`

	
	
	IpcHost bool `json:"IpcHost"`

	
	
	Linux PluginConfigLinux `json:"Linux"`

	
	
	Mounts []PluginMount `json:"Mounts"`

	
	
	Network PluginConfigNetwork `json:"Network"`

	
	
	PidHost bool `json:"PidHost"`

	
	
	PropagatedMount string `json:"PropagatedMount"`

	
	User PluginConfigUser `json:"User,omitempty"`

	
	
	WorkDir string `json:"WorkDir"`

	
	Rootfs *PluginConfigRootfs `json:"rootfs,omitempty"`
}



type PluginConfigArgs struct {

	
	
	Description string `json:"Description"`

	
	
	Name string `json:"Name"`

	
	
	Settable []string `json:"Settable"`

	
	
	Value []string `json:"Value"`
}



type PluginConfigInterface struct {

	
	
	Socket string `json:"Socket"`

	
	
	Types []PluginInterfaceType `json:"Types"`
}



type PluginConfigLinux struct {

	
	
	AllowAllDevices bool `json:"AllowAllDevices"`

	
	
	Capabilities []string `json:"Capabilities"`

	
	
	Devices []PluginDevice `json:"Devices"`
}



type PluginConfigNetwork struct {

	
	
	Type string `json:"Type"`
}



type PluginConfigRootfs struct {

	
	DiffIds []string `json:"diff_ids"`

	
	Type string `json:"type,omitempty"`
}



type PluginConfigUser struct {

	
	GID uint32 `json:"GID,omitempty"`

	
	UID uint32 `json:"UID,omitempty"`
}



type PluginSettings struct {

	
	
	Args []string `json:"Args"`

	
	
	Devices []PluginDevice `json:"Devices"`

	
	
	Env []string `json:"Env"`

	
	
	Mounts []PluginMount `json:"Mounts"`
}
