package require


type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

