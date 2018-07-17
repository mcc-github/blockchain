package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/internal/oraclematcher"
	"github.com/onsi/gomega/types"
)

type OrMatcher struct {
	Matchers []types.GomegaMatcher

	
	firstSuccessfulMatcher types.GomegaMatcher
}

func (m *OrMatcher) Match(actual interface{}) (success bool, err error) {
	m.firstSuccessfulMatcher = nil
	for _, matcher := range m.Matchers {
		success, err := matcher.Match(actual)
		if err != nil {
			return false, err
		}
		if success {
			m.firstSuccessfulMatcher = matcher
			return true, nil
		}
	}
	return false, nil
}

func (m *OrMatcher) FailureMessage(actual interface{}) (message string) {
	
	return format.Message(actual, fmt.Sprintf("To satisfy at least one of these matchers: %s", m.Matchers))
}

func (m *OrMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.firstSuccessfulMatcher.NegatedFailureMessage(actual)
}

func (m *OrMatcher) MatchMayChangeInTheFuture(actual interface{}) bool {
	

	if m.firstSuccessfulMatcher != nil {
		
		return oraclematcher.MatchMayChangeInTheFuture(m.firstSuccessfulMatcher, actual)
	} else {
		
		for _, matcher := range m.Matchers {
			if oraclematcher.MatchMayChangeInTheFuture(matcher, actual) {
				return true
			}
		}
		return false 
	}
}
