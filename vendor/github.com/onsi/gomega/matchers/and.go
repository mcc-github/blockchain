package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/internal/oraclematcher"
	"github.com/onsi/gomega/types"
)

type AndMatcher struct {
	Matchers []types.GomegaMatcher

	
	firstFailedMatcher types.GomegaMatcher
}

func (m *AndMatcher) Match(actual interface{}) (success bool, err error) {
	m.firstFailedMatcher = nil
	for _, matcher := range m.Matchers {
		success, err := matcher.Match(actual)
		if !success || err != nil {
			m.firstFailedMatcher = matcher
			return false, err
		}
	}
	return true, nil
}

func (m *AndMatcher) FailureMessage(actual interface{}) (message string) {
	return m.firstFailedMatcher.FailureMessage(actual)
}

func (m *AndMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	
	return format.Message(actual, fmt.Sprintf("To not satisfy all of these matchers: %s", m.Matchers))
}

func (m *AndMatcher) MatchMayChangeInTheFuture(actual interface{}) bool {
	

	if m.firstFailedMatcher == nil {
		
		for _, matcher := range m.Matchers {
			if oraclematcher.MatchMayChangeInTheFuture(matcher, actual) {
				return true
			}
		}
		return false 
	}
	
	return oraclematcher.MatchMayChangeInTheFuture(m.firstFailedMatcher, actual)
}
