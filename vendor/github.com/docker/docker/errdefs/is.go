package errdefs 

type causer interface {
	Cause() error
}

func getImplementer(err error) error {
	switch e := err.(type) {
	case
		ErrNotFound,
		ErrInvalidParameter,
		ErrConflict,
		ErrUnauthorized,
		ErrUnavailable,
		ErrForbidden,
		ErrSystem,
		ErrNotModified,
		ErrAlreadyExists,
		ErrNotImplemented,
		ErrCancelled,
		ErrDeadline,
		ErrDataLoss,
		ErrUnknown:
		return err
	case causer:
		return getImplementer(e.Cause())
	default:
		return err
	}
}


func IsNotFound(err error) bool {
	_, ok := getImplementer(err).(ErrNotFound)
	return ok
}


func IsInvalidParameter(err error) bool {
	_, ok := getImplementer(err).(ErrInvalidParameter)
	return ok
}


func IsConflict(err error) bool {
	_, ok := getImplementer(err).(ErrConflict)
	return ok
}


func IsUnauthorized(err error) bool {
	_, ok := getImplementer(err).(ErrUnauthorized)
	return ok
}


func IsUnavailable(err error) bool {
	_, ok := getImplementer(err).(ErrUnavailable)
	return ok
}


func IsForbidden(err error) bool {
	_, ok := getImplementer(err).(ErrForbidden)
	return ok
}


func IsSystem(err error) bool {
	_, ok := getImplementer(err).(ErrSystem)
	return ok
}


func IsNotModified(err error) bool {
	_, ok := getImplementer(err).(ErrNotModified)
	return ok
}


func IsAlreadyExists(err error) bool {
	_, ok := getImplementer(err).(ErrAlreadyExists)
	return ok
}


func IsNotImplemented(err error) bool {
	_, ok := getImplementer(err).(ErrNotImplemented)
	return ok
}


func IsUnknown(err error) bool {
	_, ok := getImplementer(err).(ErrUnknown)
	return ok
}


func IsCancelled(err error) bool {
	_, ok := getImplementer(err).(ErrCancelled)
	return ok
}


func IsDeadline(err error) bool {
	_, ok := getImplementer(err).(ErrDeadline)
	return ok
}


func IsDataLoss(err error) bool {
	_, ok := getImplementer(err).(ErrDataLoss)
	return ok
}
