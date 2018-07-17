

package remote

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/onsi/ginkgo/internal/spec_iterator"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/ginkgo/types"
)


type Server struct {
	listener        net.Listener
	reporters       []reporters.Reporter
	alives          []func() bool
	lock            *sync.Mutex
	beforeSuiteData types.RemoteBeforeSuiteData
	parallelTotal   int
	counter         int
}


func NewServer(parallelTotal int) (*Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	return &Server{
		listener:        listener,
		lock:            &sync.Mutex{},
		alives:          make([]func() bool, parallelTotal),
		beforeSuiteData: types.RemoteBeforeSuiteData{Data: nil, State: types.RemoteBeforeSuiteStatePending},
		parallelTotal:   parallelTotal,
	}, nil
}


func (server *Server) Start() {
	httpServer := &http.Server{}
	mux := http.NewServeMux()
	httpServer.Handler = mux

	
	mux.HandleFunc("/SpecSuiteWillBegin", server.specSuiteWillBegin)
	mux.HandleFunc("/BeforeSuiteDidRun", server.beforeSuiteDidRun)
	mux.HandleFunc("/AfterSuiteDidRun", server.afterSuiteDidRun)
	mux.HandleFunc("/SpecWillRun", server.specWillRun)
	mux.HandleFunc("/SpecDidComplete", server.specDidComplete)
	mux.HandleFunc("/SpecSuiteDidEnd", server.specSuiteDidEnd)

	
	mux.HandleFunc("/BeforeSuiteState", server.handleBeforeSuiteState)
	mux.HandleFunc("/RemoteAfterSuiteData", server.handleRemoteAfterSuiteData)
	mux.HandleFunc("/counter", server.handleCounter)
	mux.HandleFunc("/has-counter", server.handleHasCounter) 

	go httpServer.Serve(server.listener)
}


func (server *Server) Close() {
	server.listener.Close()
}


func (server *Server) Address() string {
	return "http://" + server.listener.Addr().String()
}






func (server *Server) readAll(request *http.Request) []byte {
	defer request.Body.Close()
	body, _ := ioutil.ReadAll(request.Body)
	return body
}

func (server *Server) RegisterReporters(reporters ...reporters.Reporter) {
	server.reporters = reporters
}

func (server *Server) specSuiteWillBegin(writer http.ResponseWriter, request *http.Request) {
	body := server.readAll(request)

	var data struct {
		Config  config.GinkgoConfigType `json:"config"`
		Summary *types.SuiteSummary     `json:"suite-summary"`
	}

	json.Unmarshal(body, &data)

	for _, reporter := range server.reporters {
		reporter.SpecSuiteWillBegin(data.Config, data.Summary)
	}
}

func (server *Server) beforeSuiteDidRun(writer http.ResponseWriter, request *http.Request) {
	body := server.readAll(request)
	var setupSummary *types.SetupSummary
	json.Unmarshal(body, &setupSummary)

	for _, reporter := range server.reporters {
		reporter.BeforeSuiteDidRun(setupSummary)
	}
}

func (server *Server) afterSuiteDidRun(writer http.ResponseWriter, request *http.Request) {
	body := server.readAll(request)
	var setupSummary *types.SetupSummary
	json.Unmarshal(body, &setupSummary)

	for _, reporter := range server.reporters {
		reporter.AfterSuiteDidRun(setupSummary)
	}
}

func (server *Server) specWillRun(writer http.ResponseWriter, request *http.Request) {
	body := server.readAll(request)
	var specSummary *types.SpecSummary
	json.Unmarshal(body, &specSummary)

	for _, reporter := range server.reporters {
		reporter.SpecWillRun(specSummary)
	}
}

func (server *Server) specDidComplete(writer http.ResponseWriter, request *http.Request) {
	body := server.readAll(request)
	var specSummary *types.SpecSummary
	json.Unmarshal(body, &specSummary)

	for _, reporter := range server.reporters {
		reporter.SpecDidComplete(specSummary)
	}
}

func (server *Server) specSuiteDidEnd(writer http.ResponseWriter, request *http.Request) {
	body := server.readAll(request)
	var suiteSummary *types.SuiteSummary
	json.Unmarshal(body, &suiteSummary)

	for _, reporter := range server.reporters {
		reporter.SpecSuiteDidEnd(suiteSummary)
	}
}





func (server *Server) RegisterAlive(node int, alive func() bool) {
	server.lock.Lock()
	defer server.lock.Unlock()
	server.alives[node-1] = alive
}

func (server *Server) nodeIsAlive(node int) bool {
	server.lock.Lock()
	defer server.lock.Unlock()
	alive := server.alives[node-1]
	if alive == nil {
		return true
	}
	return alive()
}

func (server *Server) handleBeforeSuiteState(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		dec := json.NewDecoder(request.Body)
		dec.Decode(&(server.beforeSuiteData))
	} else {
		beforeSuiteData := server.beforeSuiteData
		if beforeSuiteData.State == types.RemoteBeforeSuiteStatePending && !server.nodeIsAlive(1) {
			beforeSuiteData.State = types.RemoteBeforeSuiteStateDisappeared
		}
		enc := json.NewEncoder(writer)
		enc.Encode(beforeSuiteData)
	}
}

func (server *Server) handleRemoteAfterSuiteData(writer http.ResponseWriter, request *http.Request) {
	afterSuiteData := types.RemoteAfterSuiteData{
		CanRun: true,
	}
	for i := 2; i <= server.parallelTotal; i++ {
		afterSuiteData.CanRun = afterSuiteData.CanRun && !server.nodeIsAlive(i)
	}

	enc := json.NewEncoder(writer)
	enc.Encode(afterSuiteData)
}

func (server *Server) handleCounter(writer http.ResponseWriter, request *http.Request) {
	c := spec_iterator.Counter{}
	server.lock.Lock()
	c.Index = server.counter
	server.counter = server.counter + 1
	server.lock.Unlock()

	json.NewEncoder(writer).Encode(c)
}

func (server *Server) handleHasCounter(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(""))
}
