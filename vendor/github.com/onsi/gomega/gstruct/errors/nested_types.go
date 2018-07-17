package errors

import (
	"fmt"
	"strings"

	"github.com/onsi/gomega/types"
)



type NestingMatcher interface {
	types.GomegaMatcher

	
	Failures() []error
}


type NestedError struct {
	Path string
	Err  error
}

func (e *NestedError) Error() string {
	
	indented := strings.Replace(e.Err.Error(), "\n", "\n\t", -1)
	return fmt.Sprintf("%s:\n\t%v", e.Path, indented)
}




func Nest(path string, err error) error {
	if ag, ok := err.(AggregateError); ok {
		var errs AggregateError
		for _, e := range ag {
			errs = append(errs, Nest(path, e))
		}
		return errs
	}
	if ne, ok := err.(*NestedError); ok {
		return &NestedError{
			Path: path + ne.Path,
			Err:  ne.Err,
		}
	}
	return &NestedError{
		Path: path,
		Err:  err,
	}
}


type AggregateError []error


func (err AggregateError) Error() string {
	if len(err) == 0 {
		
		return ""
	}
	if len(err) == 1 {
		return err[0].Error()
	}
	result := fmt.Sprintf("[%s", err[0].Error())
	for i := 1; i < len(err); i++ {
		result += fmt.Sprintf(", %s", err[i].Error())
	}
	result += "]"
	return result
}
