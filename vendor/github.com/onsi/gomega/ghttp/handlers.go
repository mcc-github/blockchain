package ghttp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)



func CombineHandlers(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, handler := range handlers {
			handler(w, req)
		}
	}
}






func VerifyRequest(method string, path interface{}, rawQuery ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Expect(req.Method).Should(Equal(method), "Method mismatch")
		switch p := path.(type) {
		case types.GomegaMatcher:
			Expect(req.URL.Path).Should(p, "Path mismatch")
		default:
			Expect(req.URL.Path).Should(Equal(path), "Path mismatch")
		}
		if len(rawQuery) > 0 {
			values, err := url.ParseQuery(rawQuery[0])
			Expect(err).ShouldNot(HaveOccurred(), "Expected RawQuery is malformed")

			Expect(req.URL.Query()).Should(Equal(values), "RawQuery mismatch")
		}
	}
}



func VerifyContentType(contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Expect(req.Header.Get("Content-Type")).Should(Equal(contentType))
	}
}



func VerifyMimeType(mimeType string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Expect(strings.Split(req.Header.Get("Content-Type"), ";")[0]).Should(Equal(mimeType))
	}
}



func VerifyBasicAuth(username string, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		Expect(auth).ShouldNot(Equal(""), "Authorization header must be specified")

		decoded, err := base64.StdEncoding.DecodeString(auth[6:])
		Expect(err).ShouldNot(HaveOccurred())

		Expect(string(decoded)).Should(Equal(fmt.Sprintf("%s:%s", username, password)), "Authorization mismatch")
	}
}






func VerifyHeader(header http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for key, values := range header {
			key = http.CanonicalHeaderKey(key)
			Expect(req.Header[key]).Should(Equal(values), "Header mismatch for key: %s", key)
		}
	}
}




func VerifyHeaderKV(key string, values ...string) http.HandlerFunc {
	return VerifyHeader(http.Header{key: values})
}



func VerifyBody(expectedBody []byte) http.HandlerFunc {
	return CombineHandlers(
		func(w http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			req.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(body).Should(Equal(expectedBody), "Body Mismatch")
		},
	)
}





func VerifyJSON(expectedJSON string) http.HandlerFunc {
	return CombineHandlers(
		VerifyMimeType("application/json"),
		func(w http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			req.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(body).Should(MatchJSON(expectedJSON), "JSON Mismatch")
		},
	)
}




func VerifyJSONRepresenting(object interface{}) http.HandlerFunc {
	data, err := json.Marshal(object)
	Expect(err).ShouldNot(HaveOccurred())
	return CombineHandlers(
		VerifyContentType("application/json"),
		VerifyJSON(string(data)),
	)
}





func VerifyForm(values url.Values) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		Expect(err).ShouldNot(HaveOccurred())
		for key, vals := range values {
			Expect(r.Form[key]).Should(Equal(vals), "Form mismatch for key: %s", key)
		}
	}
}




func VerifyFormKV(key string, values ...string) http.HandlerFunc {
	return VerifyForm(url.Values{key: values})
}





func VerifyProtoRepresenting(expected proto.Message) http.HandlerFunc {
	return CombineHandlers(
		VerifyContentType("application/x-protobuf"),
		func(w http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			Expect(err).ShouldNot(HaveOccurred())
			req.Body.Close()

			expectedType := reflect.TypeOf(expected)
			actualValuePtr := reflect.New(expectedType.Elem())

			actual, ok := actualValuePtr.Interface().(proto.Message)
			Expect(ok).Should(BeTrue(), "Message value is not a proto.Message")

			err = proto.Unmarshal(body, actual)
			Expect(err).ShouldNot(HaveOccurred(), "Failed to unmarshal protobuf")

			Expect(actual).Should(Equal(expected), "ProtoBuf Mismatch")
		},
	)
}

func copyHeader(src http.Header, dst http.Header) {
	for key, value := range src {
		dst[key] = value
	}
}


func RespondWith(statusCode int, body interface{}, optionalHeader ...http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if len(optionalHeader) == 1 {
			copyHeader(optionalHeader[0], w.Header())
		}
		w.WriteHeader(statusCode)
		switch x := body.(type) {
		case string:
			w.Write([]byte(x))
		case []byte:
			w.Write(x)
		default:
			Expect(body).Should(BeNil(), "Invalid type for body.  Should be string or []byte.")
		}
	}
}


func RespondWithPtr(statusCode *int, body interface{}, optionalHeader ...http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if len(optionalHeader) == 1 {
			copyHeader(optionalHeader[0], w.Header())
		}
		w.WriteHeader(*statusCode)
		if body != nil {
			switch x := (body).(type) {
			case *string:
				w.Write([]byte(*x))
			case *[]byte:
				w.Write(*x)
			default:
				Expect(body).Should(BeNil(), "Invalid type for body.  Should be string or []byte.")
			}
		}
	}
}


func RespondWithJSONEncoded(statusCode int, object interface{}, optionalHeader ...http.Header) http.HandlerFunc {
	data, err := json.Marshal(object)
	Expect(err).ShouldNot(HaveOccurred())

	var headers http.Header
	if len(optionalHeader) == 1 {
		headers = optionalHeader[0]
	} else {
		headers = make(http.Header)
	}
	if _, found := headers["Content-Type"]; !found {
		headers["Content-Type"] = []string{"application/json"}
	}
	return RespondWith(statusCode, string(data), headers)
}


func RespondWithJSONEncodedPtr(statusCode *int, object interface{}, optionalHeader ...http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := json.Marshal(object)
		Expect(err).ShouldNot(HaveOccurred())
		var headers http.Header
		if len(optionalHeader) == 1 {
			headers = optionalHeader[0]
		} else {
			headers = make(http.Header)
		}
		if _, found := headers["Content-Type"]; !found {
			headers["Content-Type"] = []string{"application/json"}
		}
		copyHeader(headers, w.Header())
		w.WriteHeader(*statusCode)
		w.Write(data)
	}
}





func RespondWithProto(statusCode int, message proto.Message, optionalHeader ...http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := proto.Marshal(message)
		Expect(err).ShouldNot(HaveOccurred())

		var headers http.Header
		if len(optionalHeader) == 1 {
			headers = optionalHeader[0]
		} else {
			headers = make(http.Header)
		}
		if _, found := headers["Content-Type"]; !found {
			headers["Content-Type"] = []string{"application/x-protobuf"}
		}
		copyHeader(headers, w.Header())

		w.WriteHeader(statusCode)
		w.Write(data)
	}
}
