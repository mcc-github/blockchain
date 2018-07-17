







package iterator

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	ErrIterReleased = errors.New("leveldb/iterator: iterator released")
)


type IteratorSeeker interface {
	
	
	
	
	First() bool

	
	
	
	
	Last() bool

	
	
	
	
	
	Seek(key []byte) bool

	
	
	Next() bool

	
	
	Prev() bool
}


type CommonIterator interface {
	IteratorSeeker

	
	
	
	util.Releaser

	
	
	util.ReleaseSetter

	
	Valid() bool

	
	
	Error() error
}












type Iterator interface {
	CommonIterator

	
	
	
	Key() []byte

	
	
	
	Value() []byte
}





type ErrorCallbackSetter interface {
	
	
	SetErrorCallback(f func(err error))
}

type emptyIterator struct {
	util.BasicReleaser
	err error
}

func (i *emptyIterator) rErr() {
	if i.err == nil && i.Released() {
		i.err = ErrIterReleased
	}
}

func (*emptyIterator) Valid() bool            { return false }
func (i *emptyIterator) First() bool          { i.rErr(); return false }
func (i *emptyIterator) Last() bool           { i.rErr(); return false }
func (i *emptyIterator) Seek(key []byte) bool { i.rErr(); return false }
func (i *emptyIterator) Next() bool           { i.rErr(); return false }
func (i *emptyIterator) Prev() bool           { i.rErr(); return false }
func (*emptyIterator) Key() []byte            { return nil }
func (*emptyIterator) Value() []byte          { return nil }
func (i *emptyIterator) Error() error         { return i.err }



func NewEmptyIterator(err error) Iterator {
	return &emptyIterator{err: err}
}
