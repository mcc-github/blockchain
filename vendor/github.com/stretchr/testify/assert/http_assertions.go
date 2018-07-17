package assert

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)



func httpCode(handler http.HandlerFunc, method, url string, values url.Values) (int, error) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url+"?"+values.Encode(), nil)
	if err != nil {
		return -1, err
	}
	handler(w, req)
	return w.Code, nil
}






func HTTPSuccess(t TestingT, handler http.HandlerFunc, method, url string, values url.Values, msgAndArgs ...interface{}) bool {
	code, err := httpCode(handler, method, url, values)
	if err != nil {
		Fail(t, fmt.Sprintf("Failed to build test request, got error: %s", err))
		return false
	}

	isSuccessCode := code >= http.StatusOK && code <= http.StatusPartialContent
	if !isSuccessCode {
		Fail(t, fmt.Sprintf("Expected HTTP success status code for %q but received %d", url+"?"+values.Encode(), code))
	}

	return isSuccessCode
}






func HTTPRedirect(t TestingT, handler http.HandlerFunc, method, url string, values url.Values, msgAndArgs ...interface{}) bool {
	code, err := httpCode(handler, method, url, values)
	if err != nil {
		Fail(t, fmt.Sprintf("Failed to build test request, got error: %s", err))
		return false
	}

	isRedirectCode := code >= http.StatusMultipleChoices && code <= http.StatusTemporaryRedirect
	if !isRedirectCode {
		Fail(t, fmt.Sprintf("Expected HTTP redirect status code for %q but received %d", url+"?"+values.Encode(), code))
	}

	return isRedirectCode
}






func HTTPError(t TestingT, handler http.HandlerFunc, method, url string, values url.Values, msgAndArgs ...interface{}) bool {
	code, err := httpCode(handler, method, url, values)
	if err != nil {
		Fail(t, fmt.Sprintf("Failed to build test request, got error: %s", err))
		return false
	}

	isErrorCode := code >= http.StatusBadRequest
	if !isErrorCode {
		Fail(t, fmt.Sprintf("Expected HTTP error status code for %q but received %d", url+"?"+values.Encode(), code))
	}

	return isErrorCode
}



func HTTPBody(handler http.HandlerFunc, method, url string, values url.Values) string {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url+"?"+values.Encode(), nil)
	if err != nil {
		return ""
	}
	handler(w, req)
	return w.Body.String()
}







func HTTPBodyContains(t TestingT, handler http.HandlerFunc, method, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	body := HTTPBody(handler, method, url, values)

	contains := strings.Contains(body, fmt.Sprint(str))
	if !contains {
		Fail(t, fmt.Sprintf("Expected response body for \"%s\" to contain \"%s\" but found \"%s\"", url+"?"+values.Encode(), str, body))
	}

	return contains
}







func HTTPBodyNotContains(t TestingT, handler http.HandlerFunc, method, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	body := HTTPBody(handler, method, url, values)

	contains := strings.Contains(body, fmt.Sprint(str))
	if contains {
		Fail(t, fmt.Sprintf("Expected response body for \"%s\" to NOT contain \"%s\" but found \"%s\"", url+"?"+values.Encode(), str, body))
	}

	return !contains
}
