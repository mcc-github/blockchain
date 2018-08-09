




















































































package multierr 

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"go.uber.org/atomic"
)

var (
	
	_singlelineSeparator = []byte("; ")

	_newline = []byte("\n")

	
	_multilinePrefix = []byte("the following errors occurred:")

	
	
	
	
	
	
	
	
	
	
	
	
	_multilineSeparator = []byte("\n -  ")
	_multilineIndent    = []byte("    ")
)


var _bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

type errorGroup interface {
	Errors() []error
}











func Errors(err error) []error {
	if err == nil {
		return nil
	}

	
	
	
	
	
	
	
	eg, ok := err.(*multiError)
	if !ok {
		return []error{err}
	}

	errors := eg.Errors()
	result := make([]error, len(errors))
	copy(result, errors)
	return result
}








type multiError struct {
	copyNeeded atomic.Bool
	errors     []error
}

var _ errorGroup = (*multiError)(nil)




func (merr *multiError) Errors() []error {
	if merr == nil {
		return nil
	}
	return merr.errors
}

func (merr *multiError) Error() string {
	if merr == nil {
		return ""
	}

	buff := _bufferPool.Get().(*bytes.Buffer)
	buff.Reset()

	merr.writeSingleline(buff)

	result := buff.String()
	_bufferPool.Put(buff)
	return result
}

func (merr *multiError) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('+') {
		merr.writeMultiline(f)
	} else {
		merr.writeSingleline(f)
	}
}

func (merr *multiError) writeSingleline(w io.Writer) {
	first := true
	for _, item := range merr.errors {
		if first {
			first = false
		} else {
			w.Write(_singlelineSeparator)
		}
		io.WriteString(w, item.Error())
	}
}

func (merr *multiError) writeMultiline(w io.Writer) {
	w.Write(_multilinePrefix)
	for _, item := range merr.errors {
		w.Write(_multilineSeparator)
		writePrefixLine(w, _multilineIndent, fmt.Sprintf("%+v", item))
	}
}



func writePrefixLine(w io.Writer, prefix []byte, s string) {
	first := true
	for len(s) > 0 {
		if first {
			first = false
		} else {
			w.Write(prefix)
		}

		idx := strings.IndexByte(s, '\n')
		if idx < 0 {
			idx = len(s) - 1
		}

		io.WriteString(w, s[:idx+1])
		s = s[idx+1:]
	}
}

type inspectResult struct {
	
	Count int

	
	Capacity int

	
	
	FirstErrorIdx int

	
	ContainsMultiError bool
}



func inspect(errors []error) (res inspectResult) {
	first := true
	for i, err := range errors {
		if err == nil {
			continue
		}

		res.Count++
		if first {
			first = false
			res.FirstErrorIdx = i
		}

		if merr, ok := err.(*multiError); ok {
			res.Capacity += len(merr.errors)
			res.ContainsMultiError = true
		} else {
			res.Capacity++
		}
	}
	return
}


func fromSlice(errors []error) error {
	res := inspect(errors)
	switch res.Count {
	case 0:
		return nil
	case 1:
		
		return errors[res.FirstErrorIdx]
	case len(errors):
		if !res.ContainsMultiError {
			
			return &multiError{errors: errors}
		}
	}

	nonNilErrs := make([]error, 0, res.Capacity)
	for _, err := range errors[res.FirstErrorIdx:] {
		if err == nil {
			continue
		}

		if nested, ok := err.(*multiError); ok {
			nonNilErrs = append(nonNilErrs, nested.errors...)
		} else {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	return &multiError{errors: nonNilErrs}
}
































func Combine(errors ...error) error {
	return fromSlice(errors)
}
















func Append(left error, right error) error {
	switch {
	case left == nil:
		return right
	case right == nil:
		return left
	}

	if _, ok := right.(*multiError); !ok {
		if l, ok := left.(*multiError); ok && !l.copyNeeded.Swap(true) {
			
			
			errs := append(l.errors, right)
			return &multiError{errors: errs}
		} else if !ok {
			
			return &multiError{errors: []error{left, right}}
		}
	}

	
	
	errors := [2]error{left, right}
	return fromSlice(errors[0:])
}
