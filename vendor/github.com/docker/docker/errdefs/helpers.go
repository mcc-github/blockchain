package errdefs 

import "context"

type errNotFound struct{ error }

func (errNotFound) NotFound() {}

func (e errNotFound) Cause() error {
	return e.error
}


func NotFound(err error) error {
	if err == nil {
		return nil
	}
	return errNotFound{err}
}

type errInvalidParameter struct{ error }

func (errInvalidParameter) InvalidParameter() {}

func (e errInvalidParameter) Cause() error {
	return e.error
}


func InvalidParameter(err error) error {
	if err == nil {
		return nil
	}
	return errInvalidParameter{err}
}

type errConflict struct{ error }

func (errConflict) Conflict() {}

func (e errConflict) Cause() error {
	return e.error
}


func Conflict(err error) error {
	if err == nil {
		return nil
	}
	return errConflict{err}
}

type errUnauthorized struct{ error }

func (errUnauthorized) Unauthorized() {}

func (e errUnauthorized) Cause() error {
	return e.error
}


func Unauthorized(err error) error {
	if err == nil {
		return nil
	}
	return errUnauthorized{err}
}

type errUnavailable struct{ error }

func (errUnavailable) Unavailable() {}

func (e errUnavailable) Cause() error {
	return e.error
}


func Unavailable(err error) error {
	return errUnavailable{err}
}

type errForbidden struct{ error }

func (errForbidden) Forbidden() {}

func (e errForbidden) Cause() error {
	return e.error
}


func Forbidden(err error) error {
	if err == nil {
		return nil
	}
	return errForbidden{err}
}

type errSystem struct{ error }

func (errSystem) System() {}

func (e errSystem) Cause() error {
	return e.error
}


func System(err error) error {
	if err == nil {
		return nil
	}
	return errSystem{err}
}

type errNotModified struct{ error }

func (errNotModified) NotModified() {}

func (e errNotModified) Cause() error {
	return e.error
}


func NotModified(err error) error {
	if err == nil {
		return nil
	}
	return errNotModified{err}
}

type errAlreadyExists struct{ error }

func (errAlreadyExists) AlreadyExists() {}

func (e errAlreadyExists) Cause() error {
	return e.error
}


func AlreadyExists(err error) error {
	if err == nil {
		return nil
	}
	return errAlreadyExists{err}
}

type errNotImplemented struct{ error }

func (errNotImplemented) NotImplemented() {}

func (e errNotImplemented) Cause() error {
	return e.error
}


func NotImplemented(err error) error {
	if err == nil {
		return nil
	}
	return errNotImplemented{err}
}

type errUnknown struct{ error }

func (errUnknown) Unknown() {}

func (e errUnknown) Cause() error {
	return e.error
}


func Unknown(err error) error {
	if err == nil {
		return nil
	}
	return errUnknown{err}
}

type errCancelled struct{ error }

func (errCancelled) Cancelled() {}

func (e errCancelled) Cause() error {
	return e.error
}


func Cancelled(err error) error {
	if err == nil {
		return nil
	}
	return errCancelled{err}
}

type errDeadline struct{ error }

func (errDeadline) DeadlineExceeded() {}

func (e errDeadline) Cause() error {
	return e.error
}


func Deadline(err error) error {
	if err == nil {
		return nil
	}
	return errDeadline{err}
}

type errDataLoss struct{ error }

func (errDataLoss) DataLoss() {}

func (e errDataLoss) Cause() error {
	return e.error
}


func DataLoss(err error) error {
	if err == nil {
		return nil
	}
	return errDataLoss{err}
}


func FromContext(ctx context.Context) error {
	e := ctx.Err()
	if e == nil {
		return nil
	}

	if e == context.Canceled {
		return Cancelled(e)
	}
	if e == context.DeadlineExceeded {
		return Deadline(e)
	}
	return Unknown(e)
}
