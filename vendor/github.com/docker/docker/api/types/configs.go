package types 

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)






type ContainerCreateConfig struct {
	Name             string
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	AdjustCPUShares  bool
}




type ContainerRmConfig struct {
	ForceRemove, RemoveVolume, RemoveLink bool
}



type ExecConfig struct {
	User         string   
	Privileged   bool     
	Tty          bool     
	AttachStdin  bool     
	AttachStderr bool     
	AttachStdout bool     
	Detach       bool     
	DetachKeys   string   
	Env          []string 
	WorkingDir   string   
	Cmd          []string 
}


type PluginRmConfig struct {
	ForceRemove bool
}


type PluginEnableConfig struct {
	Timeout int
}


type PluginDisableConfig struct {
	ForceDisable bool
}
