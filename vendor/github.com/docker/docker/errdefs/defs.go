package errdefs 


type ErrNotFound interface {
	NotFound()
}


type ErrInvalidParameter interface {
	InvalidParameter()
}



type ErrConflict interface {
	Conflict()
}


type ErrUnauthorized interface {
	Unauthorized()
}


type ErrUnavailable interface {
	Unavailable()
}



type ErrForbidden interface {
	Forbidden()
}



type ErrSystem interface {
	System()
}


type ErrNotModified interface {
	NotModified()
}


type ErrAlreadyExists interface {
	AlreadyExists()
}


type ErrNotImplemented interface {
	NotImplemented()
}


type ErrUnknown interface {
	Unknown()
}


type ErrCancelled interface {
	Cancelled()
}


type ErrDeadline interface {
	DeadlineExceeded()
}


type ErrDataLoss interface {
	DataLoss()
}
