package gstruct

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/onsi/gomega/format"
	errorsutil "github.com/onsi/gomega/gstruct/errors"
	"github.com/onsi/gomega/types"
)







func MatchAllElements(identifier Identifier, elements Elements) types.GomegaMatcher {
	return &ElementsMatcher{
		Identifier: identifier,
		Elements:   elements,
	}
}







func MatchElements(identifier Identifier, options Options, elements Elements) types.GomegaMatcher {
	return &ElementsMatcher{
		Identifier:      identifier,
		Elements:        elements,
		IgnoreExtras:    options&IgnoreExtras != 0,
		IgnoreMissing:   options&IgnoreMissing != 0,
		AllowDuplicates: options&AllowDuplicates != 0,
	}
}




type ElementsMatcher struct {
	
	Elements Elements
	
	Identifier Identifier

	
	IgnoreExtras bool
	
	IgnoreMissing bool
	
	AllowDuplicates bool

	
	failures []error
}


type Elements map[string]types.GomegaMatcher


type Identifier func(element interface{}) string

func (m *ElementsMatcher) Match(actual interface{}) (success bool, err error) {
	if reflect.TypeOf(actual).Kind() != reflect.Slice {
		return false, fmt.Errorf("%v is type %T, expected slice", actual, actual)
	}

	m.failures = m.matchElements(actual)
	if len(m.failures) > 0 {
		return false, nil
	}
	return true, nil
}

func (m *ElementsMatcher) matchElements(actual interface{}) (errs []error) {
	
	defer func() {
		if err := recover(); err != nil {
			errs = append(errs, fmt.Errorf("panic checking %+v: %v\n%s", actual, err, debug.Stack()))
		}
	}()

	val := reflect.ValueOf(actual)
	elements := map[string]bool{}
	for i := 0; i < val.Len(); i++ {
		element := val.Index(i).Interface()
		id := m.Identifier(element)
		if elements[id] {
			if !m.AllowDuplicates {
				errs = append(errs, fmt.Errorf("found duplicate element ID %s", id))
				continue
			}
		}
		elements[id] = true

		matcher, expected := m.Elements[id]
		if !expected {
			if !m.IgnoreExtras {
				errs = append(errs, fmt.Errorf("unexpected element %s", id))
			}
			continue
		}

		match, err := matcher.Match(element)
		if match {
			continue
		}

		if err == nil {
			if nesting, ok := matcher.(errorsutil.NestingMatcher); ok {
				err = errorsutil.AggregateError(nesting.Failures())
			} else {
				err = errors.New(matcher.FailureMessage(element))
			}
		}
		errs = append(errs, errorsutil.Nest(fmt.Sprintf("[%s]", id), err))
	}

	for id := range m.Elements {
		if !elements[id] && !m.IgnoreMissing {
			errs = append(errs, fmt.Errorf("missing expected element %s", id))
		}
	}

	return errs
}

func (m *ElementsMatcher) FailureMessage(actual interface{}) (message string) {
	failure := errorsutil.AggregateError(m.failures)
	return format.Message(actual, fmt.Sprintf("to match elements: %v", failure))
}

func (m *ElementsMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match elements")
}

func (m *ElementsMatcher) Failures() []error {
	return m.failures
}
