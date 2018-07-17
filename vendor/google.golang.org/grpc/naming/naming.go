



package naming


type Operation uint8

const (
	
	Add Operation = iota
	
	Delete
)



type Update struct {
	
	Op Operation
	
	Addr string
	
	
	Metadata interface{}
}


type Resolver interface {
	
	Resolve(target string) (Watcher, error)
}


type Watcher interface {
	
	
	
	Next() ([]*Update, error)
	
	Close()
}
