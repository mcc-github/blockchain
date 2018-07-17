package swarm 

import "os"


type Config struct {
	ID string
	Meta
	Spec ConfigSpec
}


type ConfigSpec struct {
	Annotations
	Data []byte `json:",omitempty"`

	
	
	Templating *Driver `json:",omitempty"`
}


type ConfigReferenceFileTarget struct {
	Name string
	UID  string
	GID  string
	Mode os.FileMode
}


type ConfigReference struct {
	File       *ConfigReferenceFileTarget
	ConfigID   string
	ConfigName string
}
