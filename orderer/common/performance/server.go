/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package performance

import (
	"io"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const pkgLogID = "orderer/common/performance"

var logger = flogging.MustGetLogger(pkgLogID)


type BenchmarkServer struct {
	server ab.AtomicBroadcastServer
	start  chan struct{}
	halt   chan struct{}
}

var (
	servers []*BenchmarkServer
	index   int
	mutex   sync.Mutex
)


func InitializeServerPool(number int) {
	mutex = sync.Mutex{}
	index = 0
	servers = make([]*BenchmarkServer, number)
	for i := 0; i < number; i++ {
		servers[i] = &BenchmarkServer{
			server: nil,
			start:  make(chan struct{}),
			halt:   make(chan struct{}),
		}
	}
}




func GetBenchmarkServer() *BenchmarkServer {
	mutex.Lock()
	defer mutex.Unlock()

	if index >= len(servers) {
		panic("Not enough servers in the pool!")
	}

	defer func() { index++ }()
	return servers[index]
}



func GetBenchmarkServerPool() []*BenchmarkServer {
	return servers
}


func (server *BenchmarkServer) Start() {
	server.halt = make(chan struct{})
	close(server.start) 

	
	<-server.halt
}


func Halt(server *BenchmarkServer) { server.Halt() }


func (server *BenchmarkServer) Halt() {
	logger.Debug("Stopping benchmark server")
	server.server = nil
	server.start = make(chan struct{})
	close(server.halt)
}


func WaitForService(server *BenchmarkServer) { server.WaitForService() }


func (server *BenchmarkServer) WaitForService() { <-server.start }


func (server *BenchmarkServer) RegisterService(s ab.AtomicBroadcastServer) {
	server.server = s
}


func (server *BenchmarkServer) CreateBroadcastClient() *BroadcastClient {
	client := &BroadcastClient{
		requestChan:  make(chan *cb.Envelope),
		responseChan: make(chan *ab.BroadcastResponse),
		errChan:      make(chan error),
	}
	go func() {
		client.errChan <- server.server.Broadcast(client)
	}()
	return client
}



type BroadcastClient struct {
	grpc.ServerStream
	requestChan  chan *cb.Envelope
	responseChan chan *ab.BroadcastResponse
	errChan      chan error
}

func (BroadcastClient) Context() context.Context {
	return peer.NewContext(context.Background(), &peer.Peer{})
}


func (bc *BroadcastClient) SendRequest(request *cb.Envelope) {
	
	bc.requestChan <- request
}


func (bc *BroadcastClient) GetResponse() *ab.BroadcastResponse {
	return <-bc.responseChan
}


func (bc *BroadcastClient) Close() {
	close(bc.requestChan)
}


func (bc *BroadcastClient) Errors() <-chan error {
	return bc.errChan
}


func (bc *BroadcastClient) Send(br *ab.BroadcastResponse) error {
	bc.responseChan <- br
	return nil
}


func (bc *BroadcastClient) Recv() (*cb.Envelope, error) {
	msg, ok := <-bc.requestChan
	if !ok {
		return msg, io.EOF
	}
	return msg, nil
}


func (server *BenchmarkServer) CreateDeliverClient() *DeliverClient {
	client := &DeliverClient{
		requestChan:  make(chan *cb.Envelope),
		ResponseChan: make(chan *ab.DeliverResponse),
		ResultChan:   make(chan error),
	}
	go func() {
		client.ResultChan <- server.server.Deliver(client)
	}()
	return client
}



type DeliverClient struct {
	grpc.ServerStream
	requestChan  chan *cb.Envelope
	ResponseChan chan *ab.DeliverResponse
	ResultChan   chan error
}

func (DeliverClient) Context() context.Context {
	return peer.NewContext(context.Background(), &peer.Peer{})
}


func (bc *DeliverClient) SendRequest(request *cb.Envelope) {
	
	bc.requestChan <- request
}


func (bc *DeliverClient) GetResponse() *ab.DeliverResponse {
	return <-bc.ResponseChan
}


func (bc *DeliverClient) Close() {
	close(bc.requestChan)
}


func (bc *DeliverClient) Send(br *ab.DeliverResponse) error {
	bc.ResponseChan <- br
	return nil
}


func (bc *DeliverClient) Recv() (*cb.Envelope, error) {
	msg, ok := <-bc.requestChan
	if !ok {
		return msg, io.EOF
	}
	return msg, nil
}
