
package mock

import (
	sync "sync"

	token "github.com/mcc-github/blockchain/protos/token"
	client "github.com/mcc-github/blockchain/token/client"
	grpc "google.golang.org/grpc"
)

type ProverPeerClient struct {
	CreateProverClientStub        func() (*grpc.ClientConn, token.ProverClient, error)
	createProverClientMutex       sync.RWMutex
	createProverClientArgsForCall []struct {
	}
	createProverClientReturns struct {
		result1 *grpc.ClientConn
		result2 token.ProverClient
		result3 error
	}
	createProverClientReturnsOnCall map[int]struct {
		result1 *grpc.ClientConn
		result2 token.ProverClient
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ProverPeerClient) CreateProverClient() (*grpc.ClientConn, token.ProverClient, error) {
	fake.createProverClientMutex.Lock()
	ret, specificReturn := fake.createProverClientReturnsOnCall[len(fake.createProverClientArgsForCall)]
	fake.createProverClientArgsForCall = append(fake.createProverClientArgsForCall, struct {
	}{})
	fake.recordInvocation("CreateProverClient", []interface{}{})
	fake.createProverClientMutex.Unlock()
	if fake.CreateProverClientStub != nil {
		return fake.CreateProverClientStub()
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.createProverClientReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *ProverPeerClient) CreateProverClientCallCount() int {
	fake.createProverClientMutex.RLock()
	defer fake.createProverClientMutex.RUnlock()
	return len(fake.createProverClientArgsForCall)
}

func (fake *ProverPeerClient) CreateProverClientCalls(stub func() (*grpc.ClientConn, token.ProverClient, error)) {
	fake.createProverClientMutex.Lock()
	defer fake.createProverClientMutex.Unlock()
	fake.CreateProverClientStub = stub
}

func (fake *ProverPeerClient) CreateProverClientReturns(result1 *grpc.ClientConn, result2 token.ProverClient, result3 error) {
	fake.createProverClientMutex.Lock()
	defer fake.createProverClientMutex.Unlock()
	fake.CreateProverClientStub = nil
	fake.createProverClientReturns = struct {
		result1 *grpc.ClientConn
		result2 token.ProverClient
		result3 error
	}{result1, result2, result3}
}

func (fake *ProverPeerClient) CreateProverClientReturnsOnCall(i int, result1 *grpc.ClientConn, result2 token.ProverClient, result3 error) {
	fake.createProverClientMutex.Lock()
	defer fake.createProverClientMutex.Unlock()
	fake.CreateProverClientStub = nil
	if fake.createProverClientReturnsOnCall == nil {
		fake.createProverClientReturnsOnCall = make(map[int]struct {
			result1 *grpc.ClientConn
			result2 token.ProverClient
			result3 error
		})
	}
	fake.createProverClientReturnsOnCall[i] = struct {
		result1 *grpc.ClientConn
		result2 token.ProverClient
		result3 error
	}{result1, result2, result3}
}

func (fake *ProverPeerClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createProverClientMutex.RLock()
	defer fake.createProverClientMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ProverPeerClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ client.ProverPeerClient = new(ProverPeerClient)
