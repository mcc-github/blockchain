package mount 

import (
	"os"
)


type Type string


const (
	
	TypeBind Type = "bind"
	
	TypeVolume Type = "volume"
	
	TypeTmpfs Type = "tmpfs"
	
	TypeNamedPipe Type = "npipe"
)


type Mount struct {
	Type Type `json:",omitempty"`
	
	
	
	Source      string      `json:",omitempty"`
	Target      string      `json:",omitempty"`
	ReadOnly    bool        `json:",omitempty"`
	Consistency Consistency `json:",omitempty"`

	BindOptions   *BindOptions   `json:",omitempty"`
	VolumeOptions *VolumeOptions `json:",omitempty"`
	TmpfsOptions  *TmpfsOptions  `json:",omitempty"`
}


type Propagation string

const (
	
	PropagationRPrivate Propagation = "rprivate"
	
	PropagationPrivate Propagation = "private"
	
	PropagationRShared Propagation = "rshared"
	
	PropagationShared Propagation = "shared"
	
	PropagationRSlave Propagation = "rslave"
	
	PropagationSlave Propagation = "slave"
)


var Propagations = []Propagation{
	PropagationRPrivate,
	PropagationPrivate,
	PropagationRShared,
	PropagationShared,
	PropagationRSlave,
	PropagationSlave,
}


type Consistency string

const (
	
	ConsistencyFull Consistency = "consistent"
	
	ConsistencyCached Consistency = "cached"
	
	ConsistencyDelegated Consistency = "delegated"
	
	ConsistencyDefault Consistency = "default"
)


type BindOptions struct {
	Propagation Propagation `json:",omitempty"`
}


type VolumeOptions struct {
	NoCopy       bool              `json:",omitempty"`
	Labels       map[string]string `json:",omitempty"`
	DriverConfig *Driver           `json:",omitempty"`
}


type Driver struct {
	Name    string            `json:",omitempty"`
	Options map[string]string `json:",omitempty"`
}


type TmpfsOptions struct {
	
	
	
	
	
	
	
	
	SizeBytes int64 `json:",omitempty"`
	
	Mode os.FileMode `json:",omitempty"`

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
}
