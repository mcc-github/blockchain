



















package zap

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

var _errArrayElemPool = sync.Pool{New: func() interface{} {
	return &errArrayElem{}
}}


func Error(err error) Field {
	return NamedError("error", err)
}








func NamedError(key string, err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{Key: key, Type: zapcore.ErrorType, Interface: err}
}

type errArray []error

func (errs errArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range errs {
		if errs[i] == nil {
			continue
		}
		
		
		
		
		elem := _errArrayElemPool.Get().(*errArrayElem)
		elem.error = errs[i]
		arr.AppendObject(elem)
		elem.error = nil
		_errArrayElemPool.Put(elem)
	}
	return nil
}

type errArrayElem struct {
	error
}

func (e *errArrayElem) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	
	Error(e.error).AddTo(enc)
	return nil
}
