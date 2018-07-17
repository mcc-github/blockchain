






package util

import (
	"errors"
)

var (
	ErrReleased    = errors.New("leveldb: resource already relesed")
	ErrHasReleaser = errors.New("leveldb: releaser already defined")
)


type Releaser interface {
	
	
	Release()
}


type ReleaseSetter interface {
	
	
	
	
	
	
	
	SetReleaser(releaser Releaser)
}


type BasicReleaser struct {
	releaser Releaser
	released bool
}


func (r *BasicReleaser) Released() bool {
	return r.released
}


func (r *BasicReleaser) Release() {
	if !r.released {
		if r.releaser != nil {
			r.releaser.Release()
			r.releaser = nil
		}
		r.released = true
	}
}


func (r *BasicReleaser) SetReleaser(releaser Releaser) {
	if r.released {
		panic(ErrReleased)
	}
	if r.releaser != nil && releaser != nil {
		panic(ErrHasReleaser)
	}
	r.releaser = releaser
}

type NoopReleaser struct{}

func (NoopReleaser) Release() {}
