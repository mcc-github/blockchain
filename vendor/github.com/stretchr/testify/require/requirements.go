package require


type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

type tHelper interface {
	Helper()
}



type ComparisonAssertionFunc func(TestingT, interface{}, interface{}, ...interface{})



type ValueAssertionFunc func(TestingT, interface{}, ...interface{})



type BoolAssertionFunc func(TestingT, bool, ...interface{})



type ErrorAssertionFunc func(TestingT, error, ...interface{})


