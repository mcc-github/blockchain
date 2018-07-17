package swarm 

import "os"


type Secret struct {
	ID string
	Meta
	Spec SecretSpec
}


type SecretSpec struct {
	Annotations
	Data   []byte  `json:",omitempty"`
	Driver *Driver `json:",omitempty"` 

	
	
	Templating *Driver `json:",omitempty"`
}


type SecretReferenceFileTarget struct {
	Name string
	UID  string
	GID  string
	Mode os.FileMode
}


type SecretReference struct {
	File       *SecretReferenceFileTarget
	SecretID   string
	SecretName string
}
