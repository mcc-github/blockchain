

package assert

import (
	http "net/http"
	url "net/url"
	time "time"
)


func Conditionf(t TestingT, comp Comparison, msg string, args ...interface{}) bool {
	return Condition(t, comp, append([]interface{}{msg}, args...)...)
}







func Containsf(t TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	return Contains(t, s, contains, append([]interface{}{msg}, args...)...)
}


func DirExistsf(t TestingT, path string, msg string, args ...interface{}) bool {
	return DirExists(t, path, append([]interface{}{msg}, args...)...)
}






func ElementsMatchf(t TestingT, listA interface{}, listB interface{}, msg string, args ...interface{}) bool {
	return ElementsMatch(t, listA, listB, append([]interface{}{msg}, args...)...)
}





func Emptyf(t TestingT, object interface{}, msg string, args ...interface{}) bool {
	return Empty(t, object, append([]interface{}{msg}, args...)...)
}








func Equalf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	return Equal(t, expected, actual, append([]interface{}{msg}, args...)...)
}






func EqualErrorf(t TestingT, theError error, errString string, msg string, args ...interface{}) bool {
	return EqualError(t, theError, errString, append([]interface{}{msg}, args...)...)
}





func EqualValuesf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	return EqualValues(t, expected, actual, append([]interface{}{msg}, args...)...)
}







func Errorf(t TestingT, err error, msg string, args ...interface{}) bool {
	return Error(t, err, append([]interface{}{msg}, args...)...)
}




func Exactlyf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	return Exactly(t, expected, actual, append([]interface{}{msg}, args...)...)
}


func Failf(t TestingT, failureMessage string, msg string, args ...interface{}) bool {
	return Fail(t, failureMessage, append([]interface{}{msg}, args...)...)
}


func FailNowf(t TestingT, failureMessage string, msg string, args ...interface{}) bool {
	return FailNow(t, failureMessage, append([]interface{}{msg}, args...)...)
}




func Falsef(t TestingT, value bool, msg string, args ...interface{}) bool {
	return False(t, value, append([]interface{}{msg}, args...)...)
}


func FileExistsf(t TestingT, path string, msg string, args ...interface{}) bool {
	return FileExists(t, path, append([]interface{}{msg}, args...)...)
}







func HTTPBodyContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	return HTTPBodyContains(t, handler, method, url, values, str, append([]interface{}{msg}, args...)...)
}







func HTTPBodyNotContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	return HTTPBodyNotContains(t, handler, method, url, values, str, append([]interface{}{msg}, args...)...)
}






func HTTPErrorf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	return HTTPError(t, handler, method, url, values, append([]interface{}{msg}, args...)...)
}






func HTTPRedirectf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	return HTTPRedirect(t, handler, method, url, values, append([]interface{}{msg}, args...)...)
}






func HTTPSuccessf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	return HTTPSuccess(t, handler, method, url, values, append([]interface{}{msg}, args...)...)
}




func Implementsf(t TestingT, interfaceObject interface{}, object interface{}, msg string, args ...interface{}) bool {
	return Implements(t, interfaceObject, object, append([]interface{}{msg}, args...)...)
}




func InDeltaf(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	return InDelta(t, expected, actual, delta, append([]interface{}{msg}, args...)...)
}


func InDeltaMapValuesf(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	return InDeltaMapValues(t, expected, actual, delta, append([]interface{}{msg}, args...)...)
}


func InDeltaSlicef(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	return InDeltaSlice(t, expected, actual, delta, append([]interface{}{msg}, args...)...)
}


func InEpsilonf(t TestingT, expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	return InEpsilon(t, expected, actual, epsilon, append([]interface{}{msg}, args...)...)
}


func InEpsilonSlicef(t TestingT, expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	return InEpsilonSlice(t, expected, actual, epsilon, append([]interface{}{msg}, args...)...)
}


func IsTypef(t TestingT, expectedType interface{}, object interface{}, msg string, args ...interface{}) bool {
	return IsType(t, expectedType, object, append([]interface{}{msg}, args...)...)
}




func JSONEqf(t TestingT, expected string, actual string, msg string, args ...interface{}) bool {
	return JSONEq(t, expected, actual, append([]interface{}{msg}, args...)...)
}





func Lenf(t TestingT, object interface{}, length int, msg string, args ...interface{}) bool {
	return Len(t, object, length, append([]interface{}{msg}, args...)...)
}




func Nilf(t TestingT, object interface{}, msg string, args ...interface{}) bool {
	return Nil(t, object, append([]interface{}{msg}, args...)...)
}







func NoErrorf(t TestingT, err error, msg string, args ...interface{}) bool {
	return NoError(t, err, append([]interface{}{msg}, args...)...)
}







func NotContainsf(t TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	return NotContains(t, s, contains, append([]interface{}{msg}, args...)...)
}







func NotEmptyf(t TestingT, object interface{}, msg string, args ...interface{}) bool {
	return NotEmpty(t, object, append([]interface{}{msg}, args...)...)
}







func NotEqualf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	return NotEqual(t, expected, actual, append([]interface{}{msg}, args...)...)
}




func NotNilf(t TestingT, object interface{}, msg string, args ...interface{}) bool {
	return NotNil(t, object, append([]interface{}{msg}, args...)...)
}




func NotPanicsf(t TestingT, f PanicTestFunc, msg string, args ...interface{}) bool {
	return NotPanics(t, f, append([]interface{}{msg}, args...)...)
}





func NotRegexpf(t TestingT, rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	return NotRegexp(t, rx, str, append([]interface{}{msg}, args...)...)
}





func NotSubsetf(t TestingT, list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	return NotSubset(t, list, subset, append([]interface{}{msg}, args...)...)
}


func NotZerof(t TestingT, i interface{}, msg string, args ...interface{}) bool {
	return NotZero(t, i, append([]interface{}{msg}, args...)...)
}




func Panicsf(t TestingT, f PanicTestFunc, msg string, args ...interface{}) bool {
	return Panics(t, f, append([]interface{}{msg}, args...)...)
}





func PanicsWithValuef(t TestingT, expected interface{}, f PanicTestFunc, msg string, args ...interface{}) bool {
	return PanicsWithValue(t, expected, f, append([]interface{}{msg}, args...)...)
}





func Regexpf(t TestingT, rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	return Regexp(t, rx, str, append([]interface{}{msg}, args...)...)
}





func Subsetf(t TestingT, list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	return Subset(t, list, subset, append([]interface{}{msg}, args...)...)
}




func Truef(t TestingT, value bool, msg string, args ...interface{}) bool {
	return True(t, value, append([]interface{}{msg}, args...)...)
}




func WithinDurationf(t TestingT, expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface{}) bool {
	return WithinDuration(t, expected, actual, delta, append([]interface{}{msg}, args...)...)
}


func Zerof(t TestingT, i interface{}, msg string, args ...interface{}) bool {
	return Zero(t, i, append([]interface{}{msg}, args...)...)
}
