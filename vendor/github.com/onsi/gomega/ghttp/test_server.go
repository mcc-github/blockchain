
package ghttp

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"sync"

	. "github.com/onsi/gomega"
)

func new() *Server {
	return &Server{
		AllowUnhandledRequests:     false,
		UnhandledRequestStatusCode: http.StatusInternalServerError,
		writeLock:                  &sync.Mutex{},
	}
}

type routedHandler struct {
	method     string
	pathRegexp *regexp.Regexp
	path       string
	handler    http.HandlerFunc
}


func NewServer() *Server {
	s := new()
	s.HTTPTestServer = httptest.NewServer(s)
	return s
}


func NewUnstartedServer() *Server {
	s := new()
	s.HTTPTestServer = httptest.NewUnstartedServer(s)
	return s
}


func NewTLSServer() *Server {
	s := new()
	s.HTTPTestServer = httptest.NewTLSServer(s)
	return s
}

type Server struct {
	
	HTTPTestServer *httptest.Server

	
	AllowUnhandledRequests bool

	
	
	
	UnhandledRequestStatusCode int

	
	
	
	Writer io.Writer

	receivedRequests []*http.Request
	requestHandlers  []http.HandlerFunc
	routedHandlers   []routedHandler

	writeLock *sync.Mutex
	calls     int
}


func (s *Server) Start() {
	s.HTTPTestServer.Start()
}


func (s *Server) URL() string {
	return s.HTTPTestServer.URL
}


func (s *Server) Addr() string {
	return s.HTTPTestServer.Listener.Addr().String()
}


func (s *Server) Close() {
	s.writeLock.Lock()
	server := s.HTTPTestServer
	s.HTTPTestServer = nil
	s.writeLock.Unlock()

	if server != nil {
		server.Close()
	}
}









func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.writeLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		
		
		
		

		
		eAsString, ok := e.(string)
		if ok && strings.Contains(eAsString, "defer GinkgoRecover()") {
			return
		}

		
		
		
		defer func() {
			recover()
		}()
		Ω(e).Should(BeNil(), "Handler Panicked")
	}()

	if s.Writer != nil {
		s.Writer.Write([]byte(fmt.Sprintf("GHTTP Received Request: %s - %s\n", req.Method, req.URL)))
	}

	s.receivedRequests = append(s.receivedRequests, req)
	if routedHandler, ok := s.handlerForRoute(req.Method, req.URL.Path); ok {
		s.writeLock.Unlock()
		routedHandler(w, req)
	} else if s.calls < len(s.requestHandlers) {
		h := s.requestHandlers[s.calls]
		s.calls++
		s.writeLock.Unlock()
		h(w, req)
	} else {
		s.writeLock.Unlock()
		if s.AllowUnhandledRequests {
			ioutil.ReadAll(req.Body)
			req.Body.Close()
			w.WriteHeader(s.UnhandledRequestStatusCode)
		} else {
			Ω(req).Should(BeNil(), "Received Unhandled Request")
		}
	}
}


func (s *Server) ReceivedRequests() []*http.Request {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	return s.receivedRequests
}





func (s *Server) RouteToHandler(method string, path interface{}, handler http.HandlerFunc) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	rh := routedHandler{
		method:  method,
		handler: handler,
	}

	switch p := path.(type) {
	case *regexp.Regexp:
		rh.pathRegexp = p
	case string:
		rh.path = p
	default:
		panic("path must be a string or a regular expression")
	}

	for i, existingRH := range s.routedHandlers {
		if existingRH.method == method &&
			reflect.DeepEqual(existingRH.pathRegexp, rh.pathRegexp) &&
			existingRH.path == rh.path {
			s.routedHandlers[i] = rh
			return
		}
	}
	s.routedHandlers = append(s.routedHandlers, rh)
}

func (s *Server) handlerForRoute(method string, path string) (http.HandlerFunc, bool) {
	for _, rh := range s.routedHandlers {
		if rh.method == method {
			if rh.pathRegexp != nil {
				if rh.pathRegexp.Match([]byte(path)) {
					return rh.handler, true
				}
			} else if rh.path == path {
				return rh.handler, true
			}
		}
	}

	return nil, false
}


func (s *Server) AppendHandlers(handlers ...http.HandlerFunc) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	s.requestHandlers = append(s.requestHandlers, handlers...)
}




func (s *Server) SetHandler(index int, handler http.HandlerFunc) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	s.requestHandlers[index] = handler
}


func (s *Server) GetHandler(index int) http.HandlerFunc {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	return s.requestHandlers[index]
}

func (s *Server) Reset() {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	s.HTTPTestServer.CloseClientConnections()
	s.calls = 0
	s.receivedRequests = nil
	s.requestHandlers = nil
	s.routedHandlers = nil
}







func (s *Server) WrapHandler(index int, handler http.HandlerFunc) {
	existingHandler := s.GetHandler(index)
	s.SetHandler(index, CombineHandlers(existingHandler, handler))
}

func (s *Server) CloseClientConnections() {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	s.HTTPTestServer.CloseClientConnections()
}
