



















package zapcore

import (
	"fmt"
	"sync"
)


















func encodeError(key string, err error, enc ObjectEncoder) error {
	basic := err.Error()
	enc.AddString(key, basic)

	switch e := err.(type) {
	case errorGroup:
		return enc.AddArray(key+"Causes", errArray(e.Errors()))
	case fmt.Formatter:
		verbose := fmt.Sprintf("%+v", e)
		if verbose != basic {
			
			
			enc.AddString(key+"Verbose", verbose)
		}
	}
	return nil
}

type errorGroup interface {
	
	
	Errors() []error
}

type causer interface {
	
	Cause() error
}






type errArray []error

func (errs errArray) MarshalLogArray(arr ArrayEncoder) error {
	for i := range errs {
		if errs[i] == nil {
			continue
		}

		el := newErrArrayElem(errs[i])
		arr.AppendObject(el)
		el.Free()
	}
	return nil
}

var _errArrayElemPool = sync.Pool{New: func() interface{} {
	return &errArrayElem{}
}}




type errArrayElem struct{ err error }

func newErrArrayElem(err error) *errArrayElem {
	e := _errArrayElemPool.Get().(*errArrayElem)
	e.err = err
	return e
}

func (e *errArrayElem) MarshalLogArray(arr ArrayEncoder) error {
	return arr.AppendObject(e)
}

func (e *errArrayElem) MarshalLogObject(enc ObjectEncoder) error {
	return encodeError("error", e.err, enc)
}

func (e *errArrayElem) Free() {
	e.err = nil
	_errArrayElemPool.Put(e)
}
