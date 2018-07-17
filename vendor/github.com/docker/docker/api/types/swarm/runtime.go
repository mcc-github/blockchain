package swarm 


type RuntimeType string


type RuntimeURL string

const (
	
	RuntimeContainer RuntimeType = "container"
	
	RuntimePlugin RuntimeType = "plugin"

	
	RuntimeURLContainer RuntimeURL = "types.docker.com/RuntimeContainer"
	
	RuntimeURLPlugin RuntimeURL = "types.docker.com/RuntimePlugin"
)
