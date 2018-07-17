






package errors

import (
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)


var (
	ErrNotFound    = New("leveldb: not found")
	ErrReleased    = util.ErrReleased
	ErrHasReleaser = util.ErrHasReleaser
)


func New(text string) error {
	return errors.New(text)
}



type ErrCorrupted struct {
	Fd  storage.FileDesc
	Err error
}

func (e *ErrCorrupted) Error() string {
	if !e.Fd.Zero() {
		return fmt.Sprintf("%v [file=%v]", e.Err, e.Fd)
	}
	return e.Err.Error()
}


func NewErrCorrupted(fd storage.FileDesc, err error) error {
	return &ErrCorrupted{fd, err}
}



func IsCorrupted(err error) bool {
	switch err.(type) {
	case *ErrCorrupted:
		return true
	case *storage.ErrCorrupted:
		return true
	}
	return false
}



type ErrMissingFiles struct {
	Fds []storage.FileDesc
}

func (e *ErrMissingFiles) Error() string { return "file missing" }



func SetFd(err error, fd storage.FileDesc) error {
	switch x := err.(type) {
	case *ErrCorrupted:
		x.Fd = fd
		return x
	}
	return err
}
