

package require

import (
	assert "github.com/stretchr/testify/assert"
	http "net/http"
	url "net/url"
	time "time"
)


func Condition(t TestingT, comp assert.Comparison, msgAndArgs ...interface{}) {
	if assert.Condition(t, comp, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func Conditionf(t TestingT, comp assert.Comparison, msg string, args ...interface{}) {
	if assert.Conditionf(t, comp, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func Contains(t TestingT, s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	if assert.Contains(t, s, contains, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func Containsf(t TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) {
	if assert.Containsf(t, s, contains, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func DirExists(t TestingT, path string, msgAndArgs ...interface{}) {
	if assert.DirExists(t, path, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func DirExistsf(t TestingT, path string, msg string, args ...interface{}) {
	if assert.DirExistsf(t, path, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func ElementsMatch(t TestingT, listA interface{}, listB interface{}, msgAndArgs ...interface{}) {
	if assert.ElementsMatch(t, listA, listB, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func ElementsMatchf(t TestingT, listA interface{}, listB interface{}, msg string, args ...interface{}) {
	if assert.ElementsMatchf(t, listA, listB, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func Empty(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if assert.Empty(t, object, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func Emptyf(t TestingT, object interface{}, msg string, args ...interface{}) {
	if assert.Emptyf(t, object, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}








func Equal(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if assert.Equal(t, expected, actual, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func EqualError(t TestingT, theError error, errString string, msgAndArgs ...interface{}) {
	if assert.EqualError(t, theError, errString, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func EqualErrorf(t TestingT, theError error, errString string, msg string, args ...interface{}) {
	if assert.EqualErrorf(t, theError, errString, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func EqualValues(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if assert.EqualValues(t, expected, actual, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func EqualValuesf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) {
	if assert.EqualValuesf(t, expected, actual, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}








func Equalf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) {
	if assert.Equalf(t, expected, actual, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func Error(t TestingT, err error, msgAndArgs ...interface{}) {
	if assert.Error(t, err, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func Errorf(t TestingT, err error, msg string, args ...interface{}) {
	if assert.Errorf(t, err, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Exactly(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if assert.Exactly(t, expected, actual, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Exactlyf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) {
	if assert.Exactlyf(t, expected, actual, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func Fail(t TestingT, failureMessage string, msgAndArgs ...interface{}) {
	if assert.Fail(t, failureMessage, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func FailNow(t TestingT, failureMessage string, msgAndArgs ...interface{}) {
	if assert.FailNow(t, failureMessage, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func FailNowf(t TestingT, failureMessage string, msg string, args ...interface{}) {
	if assert.FailNowf(t, failureMessage, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func Failf(t TestingT, failureMessage string, msg string, args ...interface{}) {
	if assert.Failf(t, failureMessage, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func False(t TestingT, value bool, msgAndArgs ...interface{}) {
	if assert.False(t, value, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Falsef(t TestingT, value bool, msg string, args ...interface{}) {
	if assert.Falsef(t, value, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func FileExists(t TestingT, path string, msgAndArgs ...interface{}) {
	if assert.FileExists(t, path, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func FileExistsf(t TestingT, path string, msg string, args ...interface{}) {
	if assert.FileExistsf(t, path, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func HTTPBodyContains(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	if assert.HTTPBodyContains(t, handler, method, url, values, str, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func HTTPBodyContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) {
	if assert.HTTPBodyContainsf(t, handler, method, url, values, str, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func HTTPBodyNotContains(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	if assert.HTTPBodyNotContains(t, handler, method, url, values, str, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func HTTPBodyNotContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) {
	if assert.HTTPBodyNotContainsf(t, handler, method, url, values, str, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func HTTPError(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if assert.HTTPError(t, handler, method, url, values, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func HTTPErrorf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	if assert.HTTPErrorf(t, handler, method, url, values, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func HTTPRedirect(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if assert.HTTPRedirect(t, handler, method, url, values, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func HTTPRedirectf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	if assert.HTTPRedirectf(t, handler, method, url, values, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func HTTPSuccess(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	if assert.HTTPSuccess(t, handler, method, url, values, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}






func HTTPSuccessf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) {
	if assert.HTTPSuccessf(t, handler, method, url, values, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Implements(t TestingT, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	if assert.Implements(t, interfaceObject, object, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Implementsf(t TestingT, interfaceObject interface{}, object interface{}, msg string, args ...interface{}) {
	if assert.Implementsf(t, interfaceObject, object, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func InDelta(t TestingT, expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if assert.InDelta(t, expected, actual, delta, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func InDeltaMapValues(t TestingT, expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if assert.InDeltaMapValues(t, expected, actual, delta, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func InDeltaMapValuesf(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	if assert.InDeltaMapValuesf(t, expected, actual, delta, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func InDeltaSlice(t TestingT, expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	if assert.InDeltaSlice(t, expected, actual, delta, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func InDeltaSlicef(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	if assert.InDeltaSlicef(t, expected, actual, delta, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func InDeltaf(t TestingT, expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) {
	if assert.InDeltaf(t, expected, actual, delta, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func InEpsilon(t TestingT, expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	if assert.InEpsilon(t, expected, actual, epsilon, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func InEpsilonSlice(t TestingT, expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	if assert.InEpsilonSlice(t, expected, actual, epsilon, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func InEpsilonSlicef(t TestingT, expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) {
	if assert.InEpsilonSlicef(t, expected, actual, epsilon, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func InEpsilonf(t TestingT, expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) {
	if assert.InEpsilonf(t, expected, actual, epsilon, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func IsType(t TestingT, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	if assert.IsType(t, expectedType, object, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func IsTypef(t TestingT, expectedType interface{}, object interface{}, msg string, args ...interface{}) {
	if assert.IsTypef(t, expectedType, object, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func JSONEq(t TestingT, expected string, actual string, msgAndArgs ...interface{}) {
	if assert.JSONEq(t, expected, actual, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func JSONEqf(t TestingT, expected string, actual string, msg string, args ...interface{}) {
	if assert.JSONEqf(t, expected, actual, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func Len(t TestingT, object interface{}, length int, msgAndArgs ...interface{}) {
	if assert.Len(t, object, length, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func Lenf(t TestingT, object interface{}, length int, msg string, args ...interface{}) {
	if assert.Lenf(t, object, length, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Nil(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if assert.Nil(t, object, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Nilf(t TestingT, object interface{}, msg string, args ...interface{}) {
	if assert.Nilf(t, object, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func NoError(t TestingT, err error, msgAndArgs ...interface{}) {
	if assert.NoError(t, err, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func NoErrorf(t TestingT, err error, msg string, args ...interface{}) {
	if assert.NoErrorf(t, err, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func NotContains(t TestingT, s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	if assert.NotContains(t, s, contains, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func NotContainsf(t TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) {
	if assert.NotContainsf(t, s, contains, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func NotEmpty(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if assert.NotEmpty(t, object, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func NotEmptyf(t TestingT, object interface{}, msg string, args ...interface{}) {
	if assert.NotEmptyf(t, object, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func NotEqual(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if assert.NotEqual(t, expected, actual, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}







func NotEqualf(t TestingT, expected interface{}, actual interface{}, msg string, args ...interface{}) {
	if assert.NotEqualf(t, expected, actual, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func NotNil(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if assert.NotNil(t, object, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func NotNilf(t TestingT, object interface{}, msg string, args ...interface{}) {
	if assert.NotNilf(t, object, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func NotPanics(t TestingT, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if assert.NotPanics(t, f, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func NotPanicsf(t TestingT, f assert.PanicTestFunc, msg string, args ...interface{}) {
	if assert.NotPanicsf(t, f, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func NotRegexp(t TestingT, rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	if assert.NotRegexp(t, rx, str, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func NotRegexpf(t TestingT, rx interface{}, str interface{}, msg string, args ...interface{}) {
	if assert.NotRegexpf(t, rx, str, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func NotSubset(t TestingT, list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	if assert.NotSubset(t, list, subset, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func NotSubsetf(t TestingT, list interface{}, subset interface{}, msg string, args ...interface{}) {
	if assert.NotSubsetf(t, list, subset, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func NotZero(t TestingT, i interface{}, msgAndArgs ...interface{}) {
	if assert.NotZero(t, i, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func NotZerof(t TestingT, i interface{}, msg string, args ...interface{}) {
	if assert.NotZerof(t, i, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Panics(t TestingT, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if assert.Panics(t, f, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func PanicsWithValue(t TestingT, expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if assert.PanicsWithValue(t, expected, f, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func PanicsWithValuef(t TestingT, expected interface{}, f assert.PanicTestFunc, msg string, args ...interface{}) {
	if assert.PanicsWithValuef(t, expected, f, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Panicsf(t TestingT, f assert.PanicTestFunc, msg string, args ...interface{}) {
	if assert.Panicsf(t, f, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func Regexp(t TestingT, rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	if assert.Regexp(t, rx, str, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func Regexpf(t TestingT, rx interface{}, str interface{}, msg string, args ...interface{}) {
	if assert.Regexpf(t, rx, str, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func Subset(t TestingT, list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	if assert.Subset(t, list, subset, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}





func Subsetf(t TestingT, list interface{}, subset interface{}, msg string, args ...interface{}) {
	if assert.Subsetf(t, list, subset, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func True(t TestingT, value bool, msgAndArgs ...interface{}) {
	if assert.True(t, value, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func Truef(t TestingT, value bool, msg string, args ...interface{}) {
	if assert.Truef(t, value, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func WithinDuration(t TestingT, expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	if assert.WithinDuration(t, expected, actual, delta, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}




func WithinDurationf(t TestingT, expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface{}) {
	if assert.WithinDurationf(t, expected, actual, delta, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func Zero(t TestingT, i interface{}, msgAndArgs ...interface{}) {
	if assert.Zero(t, i, msgAndArgs...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}


func Zerof(t TestingT, i interface{}, msg string, args ...interface{}) {
	if assert.Zerof(t, i, msg, args...) {
		return
	}
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	t.FailNow()
}
