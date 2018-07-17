

package assert

import (
	http "net/http"
	url "net/url"
	time "time"
)


func (a *Assertions) Condition(comp Comparison, msgAndArgs ...interface{}) bool {
	return Condition(a.t, comp, msgAndArgs...)
}


func (a *Assertions) Conditionf(comp Comparison, msg string, args ...interface{}) bool {
	return Conditionf(a.t, comp, msg, args...)
}







func (a *Assertions) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) bool {
	return Contains(a.t, s, contains, msgAndArgs...)
}







func (a *Assertions) Containsf(s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	return Containsf(a.t, s, contains, msg, args...)
}


func (a *Assertions) DirExists(path string, msgAndArgs ...interface{}) bool {
	return DirExists(a.t, path, msgAndArgs...)
}


func (a *Assertions) DirExistsf(path string, msg string, args ...interface{}) bool {
	return DirExistsf(a.t, path, msg, args...)
}






func (a *Assertions) ElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) bool {
	return ElementsMatch(a.t, listA, listB, msgAndArgs...)
}






func (a *Assertions) ElementsMatchf(listA interface{}, listB interface{}, msg string, args ...interface{}) bool {
	return ElementsMatchf(a.t, listA, listB, msg, args...)
}





func (a *Assertions) Empty(object interface{}, msgAndArgs ...interface{}) bool {
	return Empty(a.t, object, msgAndArgs...)
}





func (a *Assertions) Emptyf(object interface{}, msg string, args ...interface{}) bool {
	return Emptyf(a.t, object, msg, args...)
}








func (a *Assertions) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	return Equal(a.t, expected, actual, msgAndArgs...)
}






func (a *Assertions) EqualError(theError error, errString string, msgAndArgs ...interface{}) bool {
	return EqualError(a.t, theError, errString, msgAndArgs...)
}






func (a *Assertions) EqualErrorf(theError error, errString string, msg string, args ...interface{}) bool {
	return EqualErrorf(a.t, theError, errString, msg, args...)
}





func (a *Assertions) EqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	return EqualValues(a.t, expected, actual, msgAndArgs...)
}





func (a *Assertions) EqualValuesf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	return EqualValuesf(a.t, expected, actual, msg, args...)
}








func (a *Assertions) Equalf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	return Equalf(a.t, expected, actual, msg, args...)
}







func (a *Assertions) Error(err error, msgAndArgs ...interface{}) bool {
	return Error(a.t, err, msgAndArgs...)
}







func (a *Assertions) Errorf(err error, msg string, args ...interface{}) bool {
	return Errorf(a.t, err, msg, args...)
}




func (a *Assertions) Exactly(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	return Exactly(a.t, expected, actual, msgAndArgs...)
}




func (a *Assertions) Exactlyf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	return Exactlyf(a.t, expected, actual, msg, args...)
}


func (a *Assertions) Fail(failureMessage string, msgAndArgs ...interface{}) bool {
	return Fail(a.t, failureMessage, msgAndArgs...)
}


func (a *Assertions) FailNow(failureMessage string, msgAndArgs ...interface{}) bool {
	return FailNow(a.t, failureMessage, msgAndArgs...)
}


func (a *Assertions) FailNowf(failureMessage string, msg string, args ...interface{}) bool {
	return FailNowf(a.t, failureMessage, msg, args...)
}


func (a *Assertions) Failf(failureMessage string, msg string, args ...interface{}) bool {
	return Failf(a.t, failureMessage, msg, args...)
}




func (a *Assertions) False(value bool, msgAndArgs ...interface{}) bool {
	return False(a.t, value, msgAndArgs...)
}




func (a *Assertions) Falsef(value bool, msg string, args ...interface{}) bool {
	return Falsef(a.t, value, msg, args...)
}


func (a *Assertions) FileExists(path string, msgAndArgs ...interface{}) bool {
	return FileExists(a.t, path, msgAndArgs...)
}


func (a *Assertions) FileExistsf(path string, msg string, args ...interface{}) bool {
	return FileExistsf(a.t, path, msg, args...)
}







func (a *Assertions) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	return HTTPBodyContains(a.t, handler, method, url, values, str, msgAndArgs...)
}







func (a *Assertions) HTTPBodyContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	return HTTPBodyContainsf(a.t, handler, method, url, values, str, msg, args...)
}







func (a *Assertions) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	return HTTPBodyNotContains(a.t, handler, method, url, values, str, msgAndArgs...)
}







func (a *Assertions) HTTPBodyNotContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	return HTTPBodyNotContainsf(a.t, handler, method, url, values, str, msg, args...)
}






func (a *Assertions) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	return HTTPError(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPErrorf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	return HTTPErrorf(a.t, handler, method, url, values, msg, args...)
}






func (a *Assertions) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	return HTTPRedirect(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPRedirectf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	return HTTPRedirectf(a.t, handler, method, url, values, msg, args...)
}






func (a *Assertions) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	return HTTPSuccess(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPSuccessf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	return HTTPSuccessf(a.t, handler, method, url, values, msg, args...)
}




func (a *Assertions) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	return Implements(a.t, interfaceObject, object, msgAndArgs...)
}




func (a *Assertions) Implementsf(interfaceObject interface{}, object interface{}, msg string, args ...interface{}) bool {
	return Implementsf(a.t, interfaceObject, object, msg, args...)
}




func (a *Assertions) InDelta(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	return InDelta(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaMapValues(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	return InDeltaMapValues(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaMapValuesf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	return InDeltaMapValuesf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) InDeltaSlice(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	return InDeltaSlice(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaSlicef(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	return InDeltaSlicef(a.t, expected, actual, delta, msg, args...)
}




func (a *Assertions) InDeltaf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	return InDeltaf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) InEpsilon(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	return InEpsilon(a.t, expected, actual, epsilon, msgAndArgs...)
}


func (a *Assertions) InEpsilonSlice(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	return InEpsilonSlice(a.t, expected, actual, epsilon, msgAndArgs...)
}


func (a *Assertions) InEpsilonSlicef(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	return InEpsilonSlicef(a.t, expected, actual, epsilon, msg, args...)
}


func (a *Assertions) InEpsilonf(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	return InEpsilonf(a.t, expected, actual, epsilon, msg, args...)
}


func (a *Assertions) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	return IsType(a.t, expectedType, object, msgAndArgs...)
}


func (a *Assertions) IsTypef(expectedType interface{}, object interface{}, msg string, args ...interface{}) bool {
	return IsTypef(a.t, expectedType, object, msg, args...)
}




func (a *Assertions) JSONEq(expected string, actual string, msgAndArgs ...interface{}) bool {
	return JSONEq(a.t, expected, actual, msgAndArgs...)
}




func (a *Assertions) JSONEqf(expected string, actual string, msg string, args ...interface{}) bool {
	return JSONEqf(a.t, expected, actual, msg, args...)
}





func (a *Assertions) Len(object interface{}, length int, msgAndArgs ...interface{}) bool {
	return Len(a.t, object, length, msgAndArgs...)
}





func (a *Assertions) Lenf(object interface{}, length int, msg string, args ...interface{}) bool {
	return Lenf(a.t, object, length, msg, args...)
}




func (a *Assertions) Nil(object interface{}, msgAndArgs ...interface{}) bool {
	return Nil(a.t, object, msgAndArgs...)
}




func (a *Assertions) Nilf(object interface{}, msg string, args ...interface{}) bool {
	return Nilf(a.t, object, msg, args...)
}







func (a *Assertions) NoError(err error, msgAndArgs ...interface{}) bool {
	return NoError(a.t, err, msgAndArgs...)
}







func (a *Assertions) NoErrorf(err error, msg string, args ...interface{}) bool {
	return NoErrorf(a.t, err, msg, args...)
}







func (a *Assertions) NotContains(s interface{}, contains interface{}, msgAndArgs ...interface{}) bool {
	return NotContains(a.t, s, contains, msgAndArgs...)
}







func (a *Assertions) NotContainsf(s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	return NotContainsf(a.t, s, contains, msg, args...)
}







func (a *Assertions) NotEmpty(object interface{}, msgAndArgs ...interface{}) bool {
	return NotEmpty(a.t, object, msgAndArgs...)
}







func (a *Assertions) NotEmptyf(object interface{}, msg string, args ...interface{}) bool {
	return NotEmptyf(a.t, object, msg, args...)
}







func (a *Assertions) NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	return NotEqual(a.t, expected, actual, msgAndArgs...)
}







func (a *Assertions) NotEqualf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	return NotEqualf(a.t, expected, actual, msg, args...)
}




func (a *Assertions) NotNil(object interface{}, msgAndArgs ...interface{}) bool {
	return NotNil(a.t, object, msgAndArgs...)
}




func (a *Assertions) NotNilf(object interface{}, msg string, args ...interface{}) bool {
	return NotNilf(a.t, object, msg, args...)
}




func (a *Assertions) NotPanics(f PanicTestFunc, msgAndArgs ...interface{}) bool {
	return NotPanics(a.t, f, msgAndArgs...)
}




func (a *Assertions) NotPanicsf(f PanicTestFunc, msg string, args ...interface{}) bool {
	return NotPanicsf(a.t, f, msg, args...)
}





func (a *Assertions) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	return NotRegexp(a.t, rx, str, msgAndArgs...)
}





func (a *Assertions) NotRegexpf(rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	return NotRegexpf(a.t, rx, str, msg, args...)
}





func (a *Assertions) NotSubset(list interface{}, subset interface{}, msgAndArgs ...interface{}) bool {
	return NotSubset(a.t, list, subset, msgAndArgs...)
}





func (a *Assertions) NotSubsetf(list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	return NotSubsetf(a.t, list, subset, msg, args...)
}


func (a *Assertions) NotZero(i interface{}, msgAndArgs ...interface{}) bool {
	return NotZero(a.t, i, msgAndArgs...)
}


func (a *Assertions) NotZerof(i interface{}, msg string, args ...interface{}) bool {
	return NotZerof(a.t, i, msg, args...)
}




func (a *Assertions) Panics(f PanicTestFunc, msgAndArgs ...interface{}) bool {
	return Panics(a.t, f, msgAndArgs...)
}





func (a *Assertions) PanicsWithValue(expected interface{}, f PanicTestFunc, msgAndArgs ...interface{}) bool {
	return PanicsWithValue(a.t, expected, f, msgAndArgs...)
}





func (a *Assertions) PanicsWithValuef(expected interface{}, f PanicTestFunc, msg string, args ...interface{}) bool {
	return PanicsWithValuef(a.t, expected, f, msg, args...)
}




func (a *Assertions) Panicsf(f PanicTestFunc, msg string, args ...interface{}) bool {
	return Panicsf(a.t, f, msg, args...)
}





func (a *Assertions) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	return Regexp(a.t, rx, str, msgAndArgs...)
}





func (a *Assertions) Regexpf(rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	return Regexpf(a.t, rx, str, msg, args...)
}





func (a *Assertions) Subset(list interface{}, subset interface{}, msgAndArgs ...interface{}) bool {
	return Subset(a.t, list, subset, msgAndArgs...)
}





func (a *Assertions) Subsetf(list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	return Subsetf(a.t, list, subset, msg, args...)
}




func (a *Assertions) True(value bool, msgAndArgs ...interface{}) bool {
	return True(a.t, value, msgAndArgs...)
}




func (a *Assertions) Truef(value bool, msg string, args ...interface{}) bool {
	return Truef(a.t, value, msg, args...)
}




func (a *Assertions) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {
	return WithinDuration(a.t, expected, actual, delta, msgAndArgs...)
}




func (a *Assertions) WithinDurationf(expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface{}) bool {
	return WithinDurationf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) Zero(i interface{}, msgAndArgs ...interface{}) bool {
	return Zero(a.t, i, msgAndArgs...)
}


func (a *Assertions) Zerof(i interface{}, msg string, args ...interface{}) bool {
	return Zerof(a.t, i, msg, args...)
}
