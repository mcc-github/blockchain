

package assert

import (
	http "net/http"
	url "net/url"
	time "time"
)


func Conditionf(t TestingT, comp Comparison, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Condition(t, comp, append([]interface{}{msg}, args...)...)
}







func Containsf(t TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Contains(t, s, contains, append([]interface{}{msg}, args...)...)
}


func DirExistsf(t TestingT, path string, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return DirExists(t, path, append([]interface{}{msg}, args...)...)
}






func ElementsMatchf(t TestingT, listA interface{}, listB interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return ElementsMatch(t, listA, listB, append([]interface{}{msg}, args...)...)
}





func Emptyf(t TestingT, object interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Empty(t, object, append([]interface{}{msg}, args...)...)
}








func Equalf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Equal(t, expected, actual, append([]interface{}{msg}, args...)...)
}






func EqualErrorf(t TestingT, theError error, errString string, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return EqualError(t, theError, errString, append([]interface{}{msg}, args...)...)
}





func EqualValuesf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return EqualValues(t, expected, actual, append([]interface{}{msg}, args...)...)
}







func Errorf(t TestingT, err error, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Error(t, err, append([]interface{}{msg}, args...)...)
}




func Exactlyf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Exactly(t, expected, actual, append([]interface{}{msg}, args...)...)
}


func Failf(t TestingT, failureMessage string, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Fail(t, failureMessage, append([]interface{}{msg}, args...)...)
}


func FailNowf(t TestingT, failureMessage string, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return FailNow(t, failureMessage, append([]interface{}{msg}, args...)...)
}




func Falsef(t TestingT, value bool, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return False(t, value, append([]interface{}{msg}, args...)...)
}


func FileExistsf(t TestingT, path string, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return FileExists(t, path, append([]interface{}{msg}, args...)...)
}







func HTTPBodyContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return HTTPBodyContains(t, handler, method, url, values, str, append([]interface{}{msg}, args...)...)
}







func HTTPBodyNotContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return HTTPBodyNotContains(t, handler, method, url, values, str, append([]interface{}{msg}, args...)...)
}






func HTTPErrorf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return HTTPError(t, handler, method, url, values, append([]interface{}{msg}, args...)...)
}






func HTTPRedirectf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return HTTPRedirect(t, handler, method, url, values, append([]interface{}{msg}, args...)...)
}






func HTTPSuccessf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return HTTPSuccess(t, handler, method, url, values, append([]interface{}{msg}, args...)...)
}




func Implementsf(t TestingT, interfaceObject interface{}, object interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Implements(t, interfaceObject, object, append([]interface{}{msg}, args...)...)
}




func InDeltaf(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return InDelta(t, expected, actual, delta, append([]interface{}{msg}, args...)...)
}


func InDeltaMapValuesf(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return InDeltaMapValues(t, expected, actual, delta, append([]interface{}{msg}, args...)...)
}


func InDeltaSlicef(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return InDeltaSlice(t, expected, actual, delta, append([]interface{}{msg}, args...)...)
}


func InEpsilonf(t TestingT, expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return InEpsilon(t, expected, actual, epsilon, append([]interface{}{msg}, args...)...)
}


func InEpsilonSlicef(t TestingT, expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return InEpsilonSlice(t, expected, actual, epsilon, append([]interface{}{msg}, args...)...)
}


func IsTypef(t TestingT, expectedType interface{}, object interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return IsType(t, expectedType, object, append([]interface{}{msg}, args...)...)
}




func JSONEqf(t TestingT, expected string, actual string, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return JSONEq(t, expected, actual, append([]interface{}{msg}, args...)...)
}





func Lenf(t TestingT, object interface{}, length int, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Len(t, object, length, append([]interface{}{msg}, args...)...)
}




func Nilf(t TestingT, object interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Nil(t, object, append([]interface{}{msg}, args...)...)
}







func NoErrorf(t TestingT, err error, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NoError(t, err, append([]interface{}{msg}, args...)...)
}







func NotContainsf(t TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotContains(t, s, contains, append([]interface{}{msg}, args...)...)
}







func NotEmptyf(t TestingT, object interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotEmpty(t, object, append([]interface{}{msg}, args...)...)
}







func NotEqualf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotEqual(t, expected, actual, append([]interface{}{msg}, args...)...)
}




func NotNilf(t TestingT, object interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotNil(t, object, append([]interface{}{msg}, args...)...)
}




func NotPanicsf(t TestingT, f PanicTestFunc, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotPanics(t, f, append([]interface{}{msg}, args...)...)
}





func NotRegexpf(t TestingT, rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotRegexp(t, rx, str, append([]interface{}{msg}, args...)...)
}





func NotSubsetf(t TestingT, list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotSubset(t, list, subset, append([]interface{}{msg}, args...)...)
}


func NotZerof(t TestingT, i interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotZero(t, i, append([]interface{}{msg}, args...)...)
}




func Panicsf(t TestingT, f PanicTestFunc, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Panics(t, f, append([]interface{}{msg}, args...)...)
}





func PanicsWithValuef(t TestingT, expected interface{}, f PanicTestFunc, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return PanicsWithValue(t, expected, f, append([]interface{}{msg}, args...)...)
}





func Regexpf(t TestingT, rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Regexp(t, rx, str, append([]interface{}{msg}, args...)...)
}





func Subsetf(t TestingT, list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Subset(t, list, subset, append([]interface{}{msg}, args...)...)
}




func Truef(t TestingT, value bool, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return True(t, value, append([]interface{}{msg}, args...)...)
}




func WithinDurationf(t TestingT, expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return WithinDuration(t, expected, actual, delta, append([]interface{}{msg}, args...)...)
}


func Zerof(t TestingT, i interface{}, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return Zero(t, i, append([]interface{}{msg}, args...)...)
}
