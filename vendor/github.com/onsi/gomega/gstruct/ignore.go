package gstruct

import (
	"github.com/onsi/gomega/types"
)




func Ignore() types.GomegaMatcher {
	return &IgnoreMatcher{true}
}





func Reject() types.GomegaMatcher {
	return &IgnoreMatcher{false}
}


type IgnoreMatcher struct {
	Succeed bool
}

func (m *IgnoreMatcher) Match(actual interface{}) (bool, error) {
	return m.Succeed, nil
}

func (m *IgnoreMatcher) FailureMessage(_ interface{}) (message string) {
	return "Unconditional failure"
}

func (m *IgnoreMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return "Unconditional success"
}
