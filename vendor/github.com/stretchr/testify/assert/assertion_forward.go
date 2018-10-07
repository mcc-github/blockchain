

package assert

import (
	http "net/http"
	url "net/url"
	time "time"
)


func (a *Assertions) Condition(comp Comparison, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Condition(a.t, comp, msgAndArgs...)
}


func (a *Assertions) Conditionf(comp Comparison, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Conditionf(a.t, comp, msg, args...)
}







func (a *Assertions) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Contains(a.t, s, contains, msgAndArgs...)
}







func (a *Assertions) Containsf(s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Containsf(a.t, s, contains, msg, args...)
}


func (a *Assertions) DirExists(path string, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return DirExists(a.t, path, msgAndArgs...)
}


func (a *Assertions) DirExistsf(path string, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return DirExistsf(a.t, path, msg, args...)
}






func (a *Assertions) ElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return ElementsMatch(a.t, listA, listB, msgAndArgs...)
}






func (a *Assertions) ElementsMatchf(listA interface{}, listB interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return ElementsMatchf(a.t, listA, listB, msg, args...)
}





func (a *Assertions) Empty(object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Empty(a.t, object, msgAndArgs...)
}





func (a *Assertions) Emptyf(object interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Emptyf(a.t, object, msg, args...)
}








func (a *Assertions) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Equal(a.t, expected, actual, msgAndArgs...)
}






func (a *Assertions) EqualError(theError error, errString string, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return EqualError(a.t, theError, errString, msgAndArgs...)
}






func (a *Assertions) EqualErrorf(theError error, errString string, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return EqualErrorf(a.t, theError, errString, msg, args...)
}





func (a *Assertions) EqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return EqualValues(a.t, expected, actual, msgAndArgs...)
}





func (a *Assertions) EqualValuesf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return EqualValuesf(a.t, expected, actual, msg, args...)
}








func (a *Assertions) Equalf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Equalf(a.t, expected, actual, msg, args...)
}







func (a *Assertions) Error(err error, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Error(a.t, err, msgAndArgs...)
}







func (a *Assertions) Errorf(err error, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Errorf(a.t, err, msg, args...)
}




func (a *Assertions) Exactly(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Exactly(a.t, expected, actual, msgAndArgs...)
}




func (a *Assertions) Exactlyf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Exactlyf(a.t, expected, actual, msg, args...)
}


func (a *Assertions) Fail(failureMessage string, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Fail(a.t, failureMessage, msgAndArgs...)
}


func (a *Assertions) FailNow(failureMessage string, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return FailNow(a.t, failureMessage, msgAndArgs...)
}


func (a *Assertions) FailNowf(failureMessage string, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return FailNowf(a.t, failureMessage, msg, args...)
}


func (a *Assertions) Failf(failureMessage string, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Failf(a.t, failureMessage, msg, args...)
}




func (a *Assertions) False(value bool, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return False(a.t, value, msgAndArgs...)
}




func (a *Assertions) Falsef(value bool, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Falsef(a.t, value, msg, args...)
}


func (a *Assertions) FileExists(path string, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return FileExists(a.t, path, msgAndArgs...)
}


func (a *Assertions) FileExistsf(path string, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return FileExistsf(a.t, path, msg, args...)
}







func (a *Assertions) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPBodyContains(a.t, handler, method, url, values, str, msgAndArgs...)
}







func (a *Assertions) HTTPBodyContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPBodyContainsf(a.t, handler, method, url, values, str, msg, args...)
}







func (a *Assertions) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPBodyNotContains(a.t, handler, method, url, values, str, msgAndArgs...)
}







func (a *Assertions) HTTPBodyNotContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPBodyNotContainsf(a.t, handler, method, url, values, str, msg, args...)
}






func (a *Assertions) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPError(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPErrorf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPErrorf(a.t, handler, method, url, values, msg, args...)
}






func (a *Assertions) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPRedirect(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPRedirectf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPRedirectf(a.t, handler, method, url, values, msg, args...)
}






func (a *Assertions) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPSuccess(a.t, handler, method, url, values, msgAndArgs...)
}






func (a *Assertions) HTTPSuccessf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return HTTPSuccessf(a.t, handler, method, url, values, msg, args...)
}




func (a *Assertions) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Implements(a.t, interfaceObject, object, msgAndArgs...)
}




func (a *Assertions) Implementsf(interfaceObject interface{}, object interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Implementsf(a.t, interfaceObject, object, msg, args...)
}




func (a *Assertions) InDelta(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InDelta(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaMapValues(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InDeltaMapValues(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaMapValuesf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InDeltaMapValuesf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) InDeltaSlice(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InDeltaSlice(a.t, expected, actual, delta, msgAndArgs...)
}


func (a *Assertions) InDeltaSlicef(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InDeltaSlicef(a.t, expected, actual, delta, msg, args...)
}




func (a *Assertions) InDeltaf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InDeltaf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) InEpsilon(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InEpsilon(a.t, expected, actual, epsilon, msgAndArgs...)
}


func (a *Assertions) InEpsilonSlice(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InEpsilonSlice(a.t, expected, actual, epsilon, msgAndArgs...)
}


func (a *Assertions) InEpsilonSlicef(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InEpsilonSlicef(a.t, expected, actual, epsilon, msg, args...)
}


func (a *Assertions) InEpsilonf(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return InEpsilonf(a.t, expected, actual, epsilon, msg, args...)
}


func (a *Assertions) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return IsType(a.t, expectedType, object, msgAndArgs...)
}


func (a *Assertions) IsTypef(expectedType interface{}, object interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return IsTypef(a.t, expectedType, object, msg, args...)
}




func (a *Assertions) JSONEq(expected string, actual string, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return JSONEq(a.t, expected, actual, msgAndArgs...)
}




func (a *Assertions) JSONEqf(expected string, actual string, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return JSONEqf(a.t, expected, actual, msg, args...)
}





func (a *Assertions) Len(object interface{}, length int, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Len(a.t, object, length, msgAndArgs...)
}





func (a *Assertions) Lenf(object interface{}, length int, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Lenf(a.t, object, length, msg, args...)
}




func (a *Assertions) Nil(object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Nil(a.t, object, msgAndArgs...)
}




func (a *Assertions) Nilf(object interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Nilf(a.t, object, msg, args...)
}







func (a *Assertions) NoError(err error, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NoError(a.t, err, msgAndArgs...)
}







func (a *Assertions) NoErrorf(err error, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NoErrorf(a.t, err, msg, args...)
}







func (a *Assertions) NotContains(s interface{}, contains interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotContains(a.t, s, contains, msgAndArgs...)
}







func (a *Assertions) NotContainsf(s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotContainsf(a.t, s, contains, msg, args...)
}







func (a *Assertions) NotEmpty(object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotEmpty(a.t, object, msgAndArgs...)
}







func (a *Assertions) NotEmptyf(object interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotEmptyf(a.t, object, msg, args...)
}







func (a *Assertions) NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotEqual(a.t, expected, actual, msgAndArgs...)
}







func (a *Assertions) NotEqualf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotEqualf(a.t, expected, actual, msg, args...)
}




func (a *Assertions) NotNil(object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotNil(a.t, object, msgAndArgs...)
}




func (a *Assertions) NotNilf(object interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotNilf(a.t, object, msg, args...)
}




func (a *Assertions) NotPanics(f PanicTestFunc, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotPanics(a.t, f, msgAndArgs...)
}




func (a *Assertions) NotPanicsf(f PanicTestFunc, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotPanicsf(a.t, f, msg, args...)
}





func (a *Assertions) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotRegexp(a.t, rx, str, msgAndArgs...)
}





func (a *Assertions) NotRegexpf(rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotRegexpf(a.t, rx, str, msg, args...)
}





func (a *Assertions) NotSubset(list interface{}, subset interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotSubset(a.t, list, subset, msgAndArgs...)
}





func (a *Assertions) NotSubsetf(list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotSubsetf(a.t, list, subset, msg, args...)
}


func (a *Assertions) NotZero(i interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotZero(a.t, i, msgAndArgs...)
}


func (a *Assertions) NotZerof(i interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotZerof(a.t, i, msg, args...)
}




func (a *Assertions) Panics(f PanicTestFunc, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Panics(a.t, f, msgAndArgs...)
}





func (a *Assertions) PanicsWithValue(expected interface{}, f PanicTestFunc, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return PanicsWithValue(a.t, expected, f, msgAndArgs...)
}





func (a *Assertions) PanicsWithValuef(expected interface{}, f PanicTestFunc, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return PanicsWithValuef(a.t, expected, f, msg, args...)
}




func (a *Assertions) Panicsf(f PanicTestFunc, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Panicsf(a.t, f, msg, args...)
}





func (a *Assertions) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Regexp(a.t, rx, str, msgAndArgs...)
}





func (a *Assertions) Regexpf(rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Regexpf(a.t, rx, str, msg, args...)
}





func (a *Assertions) Subset(list interface{}, subset interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Subset(a.t, list, subset, msgAndArgs...)
}





func (a *Assertions) Subsetf(list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Subsetf(a.t, list, subset, msg, args...)
}




func (a *Assertions) True(value bool, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return True(a.t, value, msgAndArgs...)
}




func (a *Assertions) Truef(value bool, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Truef(a.t, value, msg, args...)
}




func (a *Assertions) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return WithinDuration(a.t, expected, actual, delta, msgAndArgs...)
}




func (a *Assertions) WithinDurationf(expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return WithinDurationf(a.t, expected, actual, delta, msg, args...)
}


func (a *Assertions) Zero(i interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Zero(a.t, i, msgAndArgs...)
}


func (a *Assertions) Zerof(i interface{}, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Zerof(a.t, i, msg, args...)
}
