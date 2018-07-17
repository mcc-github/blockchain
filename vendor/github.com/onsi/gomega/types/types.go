package types

type GomegaFailHandler func(message string, callerSkip ...int)


type GomegaTestingT interface {
	Fatalf(format string, args ...interface{})
}




type GomegaMatcher interface {
	Match(actual interface{}) (success bool, err error)
	FailureMessage(actual interface{}) (message string)
	NegatedFailureMessage(actual interface{}) (message string)
}
