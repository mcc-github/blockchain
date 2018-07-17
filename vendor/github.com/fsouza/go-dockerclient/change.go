



package docker

import "fmt"



type ChangeType int

const (
	
	ChangeModify ChangeType = iota

	
	ChangeAdd

	
	ChangeDelete
)




type Change struct {
	Path string
	Kind ChangeType
}

func (change *Change) String() string {
	var kind string
	switch change.Kind {
	case ChangeModify:
		kind = "C"
	case ChangeAdd:
		kind = "A"
	case ChangeDelete:
		kind = "D"
	}
	return fmt.Sprintf("%s %s", kind, change.Path)
}
