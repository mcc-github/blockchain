

package require

import (
	assert "github.com/stretchr/testify/assert"
	http "net/http"
	url "net/url"
	time "time"
)


func (a *Assertions) Condition(comp assert.Comparison, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Condition(a.t, comp, msgAndArgs...)
}


func (a *Assertions) Conditionf(comp assert.Comparison, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Conditionf(a.t, comp, msg, args...)
}







func (a *Assertions) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Contains(a.t, s, contains, msgAndArgs...)
}







func (a *Assertions) Containsf(s interface{}, contains interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Containsf(a.t, s, contains, msg, args...)
}


func (a *Assertions) DirExists(path string, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	DirExists(a.t, path, msgAndArgs...)
}


func (a *Assertions) DirExistsf(path string, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	DirExistsf(a.t, path, msg, args...)
}






func (a *Assertions) ElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	ElementsMatch(a.t, listA, listB, msgAndArgs...)
}






func (a *Assertions) ElementsMatchf(listA interface{}, listB interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	ElementsMatchf(a.t, listA, listB, msg, args...)
}





func (a *Assertions) Empty(object interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Empty(a.t, object, msgAndArgs...)
}





func (a *Assertions) Emptyf(object interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Emptyf(a.t, object, msg, args...)
}








func (a *Assertions) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Equal(a.t, expected, actual, msgAndArgs...)
}






func (a *Assertions) EqualError(theError error, errString string, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	EqualError(a.t, theError, errString, msgAndArgs...)
}






func (a *Assertions) EqualErrorf(theError error, errString string, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	EqualErrorf(a.t, theError, errString, msg, args...)
}





func (a *Assertions) EqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	EqualValues(a.t, expected, actual, msgAndArgs...)
}





func (a *Assertions) EqualValuesf(expected interface{}, actual interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	EqualValuesf(a.t, expected, actual, msg, args...)
}








func (a *Assertions) Equalf(expected interface{}, actual interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Equalf(a.t, expected, actual, msg, args...)
}







func (a *Assertions) Error(err error, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Error(a.t, err, msgAndArgs...)
}







func (a *Assertions) Errorf(err error, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Errorf(a.t, err, msg, args...)
}




func (a *Assertions) Exactly(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Exactly(a.t, expected, actual, msgAndArgs...)
}




func (a *Assertions) Exactlyf(expected interface{}, actual interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Exactlyf(a.t, expected, actual, msg, args...)
}


func (a *Assertions) Fail(failureMessage string, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Fail(a.t, failureMessage, msgAndArgs...)
}


func (a *Assertions) FailNow(failureMessage string, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	FailNow(a.t, failureMessage, msgAndArgs...)
}


func (a *Assertions) FailNowf(failureMessage string, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	FailNowf(a.t, failureMessage, msg, args...)
}


func (a *Assertions) Failf(failureMessage string, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Failf(a.t, failureMessage, msg, args...)
}




func (a *Assertions) False(value bool, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	False(a.t, value, msgAndArgs...)
}




func (a *Assertions) Falsef(value bool, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Falsef(a.t, value, msg, args...)
}


func (a *Assertions) FileExists(path string, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	FileExists(a.t, path, msgAndArgs...)
}


func (a *Assertions) FileExistsf(path string, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	FileExistsf(a.t, path, msg, args...)
}







func (a *Assertions) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPBodyContains(a.t, handler, method, url, values, str, msgAndArgs...)
}







func (a *Assertions) HTTPBodyContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPBodyContainsf(a.t, handler, method, url, values, str, msg, args...)
}







func (a *Assertions) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPBodyNotContains(a.t, handler, method, url, values, str, msgAndArgs...)
}







func (a *Assertions) HTTPBodyNotContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPBodyNotContainsf(a.t, handler, method, url, values, str, msg, args...)
}






func (a *Assertions) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPError(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPErrorf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPErrorf(a.t, handler, method, url, values, msg, args...)
}






func (a *Assertions) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPRedirect(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPRedirectf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPRedirectf(a.t, handler, method, url, values, msg, args...)
}






func (a *Assertions) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPSuccess(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPSuccessf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	HTTPSuccessf(a.t, handler, method, url, values, msg, args...)
}




func (a *Assertions) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Implements(a.t, interfaceObject, object, msgAndArgs...)
}




func (a *Assertions) Implementsf(interfaceObject interface{}, object interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Implementsf(a.t, interfaceObject, object, msg, args...)
}




func (a *Assertions) InDelta(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InDelta(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaMapValues(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InDeltaMapValues(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaMapValuesf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InDeltaMapValuesf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) InDeltaSlice(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InDeltaSlice(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaSlicef(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InDeltaSlicef(a.t, expected, actual, delta, msg, args...)
}




func (a *Assertions) InDeltaf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InDeltaf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) InEpsilon(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InEpsilon(a.t, expected, actual, epsilon, msgAndArgs...)
}


func (a *Assertions) InEpsilonSlice(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InEpsilonSlice(a.t, expected, actual, epsilon, msgAndArgs...)
}


func (a *Assertions) InEpsilonSlicef(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InEpsilonSlicef(a.t, expected, actual, epsilon, msg, args...)
}


func (a *Assertions) InEpsilonf(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	InEpsilonf(a.t, expected, actual, epsilon, msg, args...)
}


func (a *Assertions) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	IsType(a.t, expectedType, object, msgAndArgs...)
}


func (a *Assertions) IsTypef(expectedType interface{}, object interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	IsTypef(a.t, expectedType, object, msg, args...)
}




func (a *Assertions) JSONEq(expected string, actual string, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	JSONEq(a.t, expected, actual, msgAndArgs...)
}




func (a *Assertions) JSONEqf(expected string, actual string, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	JSONEqf(a.t, expected, actual, msg, args...)
}





func (a *Assertions) Len(object interface{}, length int, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Len(a.t, object, length, msgAndArgs...)
}





func (a *Assertions) Lenf(object interface{}, length int, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Lenf(a.t, object, length, msg, args...)
}




func (a *Assertions) Nil(object interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Nil(a.t, object, msgAndArgs...)
}




func (a *Assertions) Nilf(object interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Nilf(a.t, object, msg, args...)
}







func (a *Assertions) NoError(err error, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NoError(a.t, err, msgAndArgs...)
}







func (a *Assertions) NoErrorf(err error, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NoErrorf(a.t, err, msg, args...)
}







func (a *Assertions) NotContains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotContains(a.t, s, contains, msgAndArgs...)
}







func (a *Assertions) NotContainsf(s interface{}, contains interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotContainsf(a.t, s, contains, msg, args...)
}







func (a *Assertions) NotEmpty(object interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotEmpty(a.t, object, msgAndArgs...)
}







func (a *Assertions) NotEmptyf(object interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotEmptyf(a.t, object, msg, args...)
}







func (a *Assertions) NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotEqual(a.t, expected, actual, msgAndArgs...)
}







func (a *Assertions) NotEqualf(expected interface{}, actual interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotEqualf(a.t, expected, actual, msg, args...)
}




func (a *Assertions) NotNil(object interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotNil(a.t, object, msgAndArgs...)
}




func (a *Assertions) NotNilf(object interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotNilf(a.t, object, msg, args...)
}




func (a *Assertions) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotPanics(a.t, f, msgAndArgs...)
}




func (a *Assertions) NotPanicsf(f assert.PanicTestFunc, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotPanicsf(a.t, f, msg, args...)
}





func (a *Assertions) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotRegexp(a.t, rx, str, msgAndArgs...)
}





func (a *Assertions) NotRegexpf(rx interface{}, str interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotRegexpf(a.t, rx, str, msg, args...)
}





func (a *Assertions) NotSubset(list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotSubset(a.t, list, subset, msgAndArgs...)
}





func (a *Assertions) NotSubsetf(list interface{}, subset interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotSubsetf(a.t, list, subset, msg, args...)
}


func (a *Assertions) NotZero(i interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotZero(a.t, i, msgAndArgs...)
}


func (a *Assertions) NotZerof(i interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	NotZerof(a.t, i, msg, args...)
}




func (a *Assertions) Panics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Panics(a.t, f, msgAndArgs...)
}





func (a *Assertions) PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	PanicsWithValue(a.t, expected, f, msgAndArgs...)
}





func (a *Assertions) PanicsWithValuef(expected interface{}, f assert.PanicTestFunc, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	PanicsWithValuef(a.t, expected, f, msg, args...)
}




func (a *Assertions) Panicsf(f assert.PanicTestFunc, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Panicsf(a.t, f, msg, args...)
}





func (a *Assertions) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Regexp(a.t, rx, str, msgAndArgs...)
}





func (a *Assertions) Regexpf(rx interface{}, str interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Regexpf(a.t, rx, str, msg, args...)
}





func (a *Assertions) Subset(list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Subset(a.t, list, subset, msgAndArgs...)
}





func (a *Assertions) Subsetf(list interface{}, subset interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Subsetf(a.t, list, subset, msg, args...)
}




func (a *Assertions) True(value bool, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	True(a.t, value, msgAndArgs...)
}




func (a *Assertions) Truef(value bool, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Truef(a.t, value, msg, args...)
}




func (a *Assertions) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	WithinDuration(a.t, expected, actual, delta, msgAndArgs...)
}




func (a *Assertions) WithinDurationf(expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	WithinDurationf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) Zero(i interface{}, msgAndArgs ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Zero(a.t, i, msgAndArgs...)
}


func (a *Assertions) Zerof(i interface{}, msg string, args ...interface{}) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	Zerof(a.t, i, msg, args...)
}
