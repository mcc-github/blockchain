package gstruct

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)





func PointTo(matcher types.GomegaMatcher) types.GomegaMatcher {
	return &PointerMatcher{
		Matcher: matcher,
	}
}

type PointerMatcher struct {
	Matcher types.GomegaMatcher

	
	failure string
}

func (m *PointerMatcher) Match(actual interface{}) (bool, error) {
	val := reflect.ValueOf(actual)

	
	if val.Kind() != reflect.Ptr {
		return false, fmt.Errorf("PointerMatcher expects a pointer but we have '%s'", val.Kind())
	}

	if !val.IsValid() || val.IsNil() {
		m.failure = format.Message(actual, "not to be <nil>")
		return false, nil
	}

	
	elem := val.Elem().Interface()
	match, err := m.Matcher.Match(elem)
	if !match {
		m.failure = m.Matcher.FailureMessage(elem)
	}
	return match, err
}

func (m *PointerMatcher) FailureMessage(_ interface{}) (message string) {
	return m.failure
}

func (m *PointerMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(actual)
}
