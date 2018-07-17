package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/internal/oraclematcher"
	"github.com/onsi/gomega/types"
)

type WithTransformMatcher struct {
	
	Transform interface{} 
	Matcher   types.GomegaMatcher

	
	transformArgType reflect.Type

	
	transformedValue interface{}
}

func NewWithTransformMatcher(transform interface{}, matcher types.GomegaMatcher) *WithTransformMatcher {
	if transform == nil {
		panic("transform function cannot be nil")
	}
	txType := reflect.TypeOf(transform)
	if txType.NumIn() != 1 {
		panic("transform function must have 1 argument")
	}
	if txType.NumOut() != 1 {
		panic("transform function must have 1 return value")
	}

	return &WithTransformMatcher{
		Transform:        transform,
		Matcher:          matcher,
		transformArgType: reflect.TypeOf(transform).In(0),
	}
}

func (m *WithTransformMatcher) Match(actual interface{}) (bool, error) {
	
	actualType := reflect.TypeOf(actual)
	if !actualType.AssignableTo(m.transformArgType) {
		return false, fmt.Errorf("Transform function expects '%s' but we have '%s'", m.transformArgType, actualType)
	}

	
	fn := reflect.ValueOf(m.Transform)
	result := fn.Call([]reflect.Value{reflect.ValueOf(actual)})
	m.transformedValue = result[0].Interface() 

	return m.Matcher.Match(m.transformedValue)
}

func (m *WithTransformMatcher) FailureMessage(_ interface{}) (message string) {
	return m.Matcher.FailureMessage(m.transformedValue)
}

func (m *WithTransformMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(m.transformedValue)
}

func (m *WithTransformMatcher) MatchMayChangeInTheFuture(_ interface{}) bool {
	
	
	
	
	
	return oraclematcher.MatchMayChangeInTheFuture(m.Matcher, m.transformedValue)
}
