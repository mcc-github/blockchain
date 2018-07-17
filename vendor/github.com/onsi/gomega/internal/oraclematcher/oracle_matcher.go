package oraclematcher

import "github.com/onsi/gomega/types"


type OracleMatcher interface {
	MatchMayChangeInTheFuture(actual interface{}) bool
}

func MatchMayChangeInTheFuture(matcher types.GomegaMatcher, value interface{}) bool {
	oracleMatcher, ok := matcher.(OracleMatcher)
	if !ok {
		return true
	}

	return oracleMatcher.MatchMayChangeInTheFuture(value)
}
