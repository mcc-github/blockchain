





package iterator

import (
	"github.com/syndtr/goleveldb/leveldb/util"
)


type BasicArray interface {
	
	Len() int

	
	
	Search(key []byte) int
}


type Array interface {
	BasicArray

	
	Index(i int) (key, value []byte)
}


type ArrayIndexer interface {
	BasicArray

	
	Get(i int) Iterator
}

type basicArrayIterator struct {
	util.BasicReleaser
	array BasicArray
	pos   int
	err   error
}

func (i *basicArrayIterator) Valid() bool {
	return i.pos >= 0 && i.pos < i.array.Len() && !i.Released()
}

func (i *basicArrayIterator) First() bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	if i.array.Len() == 0 {
		i.pos = -1
		return false
	}
	i.pos = 0
	return true
}

func (i *basicArrayIterator) Last() bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	n := i.array.Len()
	if n == 0 {
		i.pos = 0
		return false
	}
	i.pos = n - 1
	return true
}

func (i *basicArrayIterator) Seek(key []byte) bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	n := i.array.Len()
	if n == 0 {
		i.pos = 0
		return false
	}
	i.pos = i.array.Search(key)
	if i.pos >= n {
		return false
	}
	return true
}

func (i *basicArrayIterator) Next() bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	i.pos++
	if n := i.array.Len(); i.pos >= n {
		i.pos = n
		return false
	}
	return true
}

func (i *basicArrayIterator) Prev() bool {
	if i.Released() {
		i.err = ErrIterReleased
		return false
	}

	i.pos--
	if i.pos < 0 {
		i.pos = -1
		return false
	}
	return true
}

func (i *basicArrayIterator) Error() error { return i.err }

type arrayIterator struct {
	basicArrayIterator
	array      Array
	pos        int
	key, value []byte
}

func (i *arrayIterator) updateKV() {
	if i.pos == i.basicArrayIterator.pos {
		return
	}
	i.pos = i.basicArrayIterator.pos
	if i.Valid() {
		i.key, i.value = i.array.Index(i.pos)
	} else {
		i.key = nil
		i.value = nil
	}
}

func (i *arrayIterator) Key() []byte {
	i.updateKV()
	return i.key
}

func (i *arrayIterator) Value() []byte {
	i.updateKV()
	return i.value
}

type arrayIteratorIndexer struct {
	basicArrayIterator
	array ArrayIndexer
}

func (i *arrayIteratorIndexer) Get() Iterator {
	if i.Valid() {
		return i.array.Get(i.basicArrayIterator.pos)
	}
	return nil
}


func NewArrayIterator(array Array) Iterator {
	return &arrayIterator{
		basicArrayIterator: basicArrayIterator{array: array, pos: -1},
		array:              array,
		pos:                -1,
	}
}


func NewArrayIndexer(array ArrayIndexer) IteratorIndexer {
	return &arrayIteratorIndexer{
		basicArrayIterator: basicArrayIterator{array: array, pos: -1},
		array:              array,
	}
}
