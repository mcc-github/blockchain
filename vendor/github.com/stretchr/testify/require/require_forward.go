

package require

import (
	assert "github.com/stretchr/testify/assert"
	http "net/http"
	url "net/url"
	time "time"
)


func (a *Assertions) Condition(comp assert.Comparison, msgAndArgs ...interface{}) {
	Condition(a.t, comp, msgAndArgs...)
}


func (a *Assertions) Conditionf(comp assert.Comparison, msg string, args ...interface{}) {
	Conditionf(a.t, comp, msg, args...)
}







func (a *Assertions) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	Contains(a.t, s, contains, msgAndArgs...)
}







func (a *Assertions) Containsf(s interface{}, contains interface{}, msg string, args ...interface{}) {
	Containsf(a.t, s, contains, msg, args...)
}


func (a *Assertions) DirExists(path string, msgAndArgs ...interface{}) {
	DirExists(a.t, path, msgAndArgs...)
}


func (a *Assertions) DirExistsf(path string, msg string, args ...interface{}) {
	DirExistsf(a.t, path, msg, args...)
}






func (a *Assertions) ElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) {
	ElementsMatch(a.t, listA, listB, msgAndArgs...)
}






func (a *Assertions) ElementsMatchf(listA interface{}, listB interface{}, msg string, args ...interface{}) {
	ElementsMatchf(a.t, listA, listB, msg, args...)
}





func (a *Assertions) Empty(object interface{}, msgAndArgs ...interface{}) {
	Empty(a.t, object, msgAndArgs...)
}





func (a *Assertions) Emptyf(object interface{}, msg string, args ...interface{}) {
	Emptyf(a.t, object, msg, args...)
}








func (a *Assertions) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	Equal(a.t, expected, actual, msgAndArgs...)
}






func (a *Assertions) EqualError(theError error, errString string, msgAndArgs ...interface{}) {
	EqualError(a.t, theError, errString, msgAndArgs...)
}






func (a *Assertions) EqualErrorf(theError error, errString string, msg string, args ...interface{}) {
	EqualErrorf(a.t, theError, errString, msg, args...)
}





func (a *Assertions) EqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	EqualValues(a.t, expected, actual, msgAndArgs...)
}





func (a *Assertions) EqualValuesf(expected interface{}, actual interface{}, msg string, args ...interface{}) {
	EqualValuesf(a.t, expected, actual, msg, args...)
}








func (a *Assertions) Equalf(expected interface{}, actual interface{}, msg string, args ...interface{}) {
	Equalf(a.t, expected, actual, msg, args...)
}







func (a *Assertions) Error(err error, msgAndArgs ...interface{}) {
	Error(a.t, err, msgAndArgs...)
}







func (a *Assertions) Errorf(err error, msg string, args ...interface{}) {
	Errorf(a.t, err, msg, args...)
}




func (a *Assertions) Exactly(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	Exactly(a.t, expected, actual, msgAndArgs...)
}




func (a *Assertions) Exactlyf(expected interface{}, actual interface{}, msg string, args ...interface{}) {
	Exactlyf(a.t, expected, actual, msg, args...)
}


func (a *Assertions) Fail(failureMessage string, msgAndArgs ...interface{}) {
	Fail(a.t, failureMessage, msgAndArgs...)
}


func (a *Assertions) FailNow(failureMessage string, msgAndArgs ...interface{}) {
	FailNow(a.t, failureMessage, msgAndArgs...)
}


func (a *Assertions) FailNowf(failureMessage string, msg string, args ...interface{}) {
	FailNowf(a.t, failureMessage, msg, args...)
}


func (a *Assertions) Failf(failureMessage string, msg string, args ...interface{}) {
	Failf(a.t, failureMessage, msg, args...)
}




func (a *Assertions) False(value bool, msgAndArgs ...interface{}) {
	False(a.t, value, msgAndArgs...)
}




func (a *Assertions) Falsef(value bool, msg string, args ...interface{}) {
	Falsef(a.t, value, msg, args...)
}


func (a *Assertions) FileExists(path string, msgAndArgs ...interface{}) {
	FileExists(a.t, path, msgAndArgs...)
}


func (a *Assertions) FileExistsf(path string, msg string, args ...interface{}) {
	FileExistsf(a.t, path, msg, args...)
}







func (a *Assertions) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	HTTPBodyContains(a.t, handler, method, url, values, str, msgAndArgs...)
}







func (a *Assertions) HTTPBodyContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) {
	HTTPBodyContainsf(a.t, handler, method, url, values, str, msg, args...)
}







func (a *Assertions) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	HTTPBodyNotContains(a.t, handler, method, url, values, str, msgAndArgs...)
}







func (a *Assertions) HTTPBodyNotContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) {
	HTTPBodyNotContainsf(a.t, handler, method, url, values, str, msg, args...)
}






func (a *Assertions) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	HTTPError(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPErrorf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	HTTPErrorf(a.t, handler, method, url, values, msg, args...)
}






func (a *Assertions) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	HTTPRedirect(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPRedirectf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	HTTPRedirectf(a.t, handler, method, url, values, msg, args...)
}






func (a *Assertions) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	HTTPSuccess(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPSuccessf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	HTTPSuccessf(a.t, handler, method, url, values, msg, args...)
}




func (a *Assertions) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	Implements(a.t, interfaceObject, object, msgAndArgs...)
}




func (a *Assertions) Implementsf(interfaceObject interface{}, object interface{}, msg string, args ...interface{}) {
	Implementsf(a.t, interfaceObject, object, msg, args...)
}




func (a *Assertions) InDelta(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	InDelta(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaMapValues(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	InDeltaMapValues(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaMapValuesf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	InDeltaMapValuesf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) InDeltaSlice(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	InDeltaSlice(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaSlicef(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	InDeltaSlicef(a.t, expected, actual, delta, msg, args...)
}




func (a *Assertions) InDeltaf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	InDeltaf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) InEpsilon(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	InEpsilon(a.t, expected, actual, epsilon, msgAndArgs...)
}


func (a *Assertions) InEpsilonSlice(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	InEpsilonSlice(a.t, expected, actual, epsilon, msgAndArgs...)
}


func (a *Assertions) InEpsilonSlicef(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) {
	InEpsilonSlicef(a.t, expected, actual, epsilon, msg, args...)
}


func (a *Assertions) InEpsilonf(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) {
	InEpsilonf(a.t, expected, actual, epsilon, msg, args...)
}


func (a *Assertions) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	IsType(a.t, expectedType, object, msgAndArgs...)
}


func (a *Assertions) IsTypef(expectedType interface{}, object interface{}, msg string, args ...interface{}) {
	IsTypef(a.t, expectedType, object, msg, args...)
}




func (a *Assertions) JSONEq(expected string, actual string, msgAndArgs ...interface{}) {
	JSONEq(a.t, expected, actual, msgAndArgs...)
}




func (a *Assertions) JSONEqf(expected string, actual string, msg string, args ...interface{}) {
	JSONEqf(a.t, expected, actual, msg, args...)
}





func (a *Assertions) Len(object interface{}, length int, msgAndArgs ...interface{}) {
	Len(a.t, object, length, msgAndArgs...)
}





func (a *Assertions) Lenf(object interface{}, length int, msg string, args ...interface{}) {
	Lenf(a.t, object, length, msg, args...)
}




func (a *Assertions) Nil(object interface{}, msgAndArgs ...interface{}) {
	Nil(a.t, object, msgAndArgs...)
}




func (a *Assertions) Nilf(object interface{}, msg string, args ...interface{}) {
	Nilf(a.t, object, msg, args...)
}







func (a *Assertions) NoError(err error, msgAndArgs ...interface{}) {
	NoError(a.t, err, msgAndArgs...)
}







func (a *Assertions) NoErrorf(err error, msg string, args ...interface{}) {
	NoErrorf(a.t, err, msg, args...)
}







func (a *Assertions) NotContains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	NotContains(a.t, s, contains, msgAndArgs...)
}







func (a *Assertions) NotContainsf(s interface{}, contains interface{}, msg string, args ...interface{}) {
	NotContainsf(a.t, s, contains, msg, args...)
}







func (a *Assertions) NotEmpty(object interface{}, msgAndArgs ...interface{}) {
	NotEmpty(a.t, object, msgAndArgs...)
}







func (a *Assertions) NotEmptyf(object interface{}, msg string, args ...interface{}) {
	NotEmptyf(a.t, object, msg, args...)
}







func (a *Assertions) NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	NotEqual(a.t, expected, actual, msgAndArgs...)
}







func (a *Assertions) NotEqualf(expected interface{}, actual interface{}, msg string, args ...interface{}) {
	NotEqualf(a.t, expected, actual, msg, args...)
}




func (a *Assertions) NotNil(object interface{}, msgAndArgs ...interface{}) {
	NotNil(a.t, object, msgAndArgs...)
}




func (a *Assertions) NotNilf(object interface{}, msg string, args ...interface{}) {
	NotNilf(a.t, object, msg, args...)
}




func (a *Assertions) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	NotPanics(a.t, f, msgAndArgs...)
}




func (a *Assertions) NotPanicsf(f assert.PanicTestFunc, msg string, args ...interface{}) {
	NotPanicsf(a.t, f, msg, args...)
}





func (a *Assertions) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	NotRegexp(a.t, rx, str, msgAndArgs...)
}





func (a *Assertions) NotRegexpf(rx interface{}, str interface{}, msg string, args ...interface{}) {
	NotRegexpf(a.t, rx, str, msg, args...)
}





func (a *Assertions) NotSubset(list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	NotSubset(a.t, list, subset, msgAndArgs...)
}





func (a *Assertions) NotSubsetf(list interface{}, subset interface{}, msg string, args ...interface{}) {
	NotSubsetf(a.t, list, subset, msg, args...)
}


func (a *Assertions) NotZero(i interface{}, msgAndArgs ...interface{}) {
	NotZero(a.t, i, msgAndArgs...)
}


func (a *Assertions) NotZerof(i interface{}, msg string, args ...interface{}) {
	NotZerof(a.t, i, msg, args...)
}




func (a *Assertions) Panics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	Panics(a.t, f, msgAndArgs...)
}





func (a *Assertions) PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	PanicsWithValue(a.t, expected, f, msgAndArgs...)
}





func (a *Assertions) PanicsWithValuef(expected interface{}, f assert.PanicTestFunc, msg string, args ...interface{}) {
	PanicsWithValuef(a.t, expected, f, msg, args...)
}




func (a *Assertions) Panicsf(f assert.PanicTestFunc, msg string, args ...interface{}) {
	Panicsf(a.t, f, msg, args...)
}





func (a *Assertions) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	Regexp(a.t, rx, str, msgAndArgs...)
}





func (a *Assertions) Regexpf(rx interface{}, str interface{}, msg string, args ...interface{}) {
	Regexpf(a.t, rx, str, msg, args...)
}





func (a *Assertions) Subset(list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	Subset(a.t, list, subset, msgAndArgs...)
}





func (a *Assertions) Subsetf(list interface{}, subset interface{}, msg string, args ...interface{}) {
	Subsetf(a.t, list, subset, msg, args...)
}




func (a *Assertions) True(value bool, msgAndArgs ...interface{}) {
	True(a.t, value, msgAndArgs...)
}




func (a *Assertions) Truef(value bool, msg string, args ...interface{}) {
	Truef(a.t, value, msg, args...)
}




func (a *Assertions) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	WithinDuration(a.t, expected, actual, delta, msgAndArgs...)
}




func (a *Assertions) WithinDurationf(expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface{}) {
	WithinDurationf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) Zero(i interface{}, msgAndArgs ...interface{}) {
	Zero(a.t, i, msgAndArgs...)
}


func (a *Assertions) Zerof(i interface{}, msg string, args ...interface{}) {
	Zerof(a.t, i, msg, args...)
}
