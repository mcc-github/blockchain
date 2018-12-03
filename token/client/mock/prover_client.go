
package mock

import (
	context "context"
	sync "sync"

	token "github.com/mcc-github/blockchain/protos/token"
	grpc "google.golang.org/grpc"
)

type ProverClient struct {
	ProcessCommandStub        func(context.Context, *token.SignedCommand, ...grpc.CallOption) (*token.SignedCommandResponse, error)
	processCommandMutex       sync.RWMutex
	processCommandArgsForCall []struct {
		arg1 context.Context
		arg2 *token.SignedCommand
		arg3 []grpc.CallOption
	}
	processCommandReturns struct {
		result1 *token.SignedCommandResponse
		result2 error
	}
	processCommandReturnsOnCall map[int]struct {
		result1 *token.SignedCommandResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ProverClient) ProcessCommand(arg1 context.Context, arg2 *token.SignedCommand, arg3 ...grpc.CallOption) (*token.SignedCommandResponse, error) {
	fake.processCommandMutex.Lock()
	ret, specificReturn := fake.processCommandReturnsOnCall[len(fake.processCommandArgsForCall)]
	fake.processCommandArgsForCall = append(fake.processCommandArgsForCall, struct {
		arg1 context.Context
		arg2 *token.SignedCommand
		arg3 []grpc.CallOption
	}{arg1, arg2, arg3})
	fake.recordInvocation("ProcessCommand", []interface{}{arg1, arg2, arg3})
	fake.processCommandMutex.Unlock()
	if fake.ProcessCommandStub != nil {
		return fake.ProcessCommandStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.processCommandReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ProverClient) ProcessCommandCallCount() int {
	fake.processCommandMutex.RLock()
	defer fake.processCommandMutex.RUnlock()
	return len(fake.processCommandArgsForCall)
}

func (fake *ProverClient) ProcessCommandCalls(stub func(context.Context, *token.SignedCommand, ...grpc.CallOption) (*token.SignedCommandResponse, error)) {
	fake.processCommandMutex.Lock()
	defer fake.processCommandMutex.Unlock()
	fake.ProcessCommandStub = stub
}

func (fake *ProverClient) ProcessCommandArgsForCall(i int) (context.Context, *token.SignedCommand, []grpc.CallOption) {
	fake.processCommandMutex.RLock()
	defer fake.processCommandMutex.RUnlock()
	argsForCall := fake.processCommandArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ProverClient) ProcessCommandReturns(result1 *token.SignedCommandResponse, result2 error) {
	fake.processCommandMutex.Lock()
	defer fake.processCommandMutex.Unlock()
	fake.ProcessCommandStub = nil
	fake.processCommandReturns = struct {
		result1 *token.SignedCommandResponse
		result2 error
	}{result1, result2}
}

func (fake *ProverClient) ProcessCommandReturnsOnCall(i int, result1 *token.SignedCommandResponse, result2 error) {
	fake.processCommandMutex.Lock()
	defer fake.processCommandMutex.Unlock()
	fake.ProcessCommandStub = nil
	if fake.processCommandReturnsOnCall == nil {
		fake.processCommandReturnsOnCall = make(map[int]struct {
			result1 *token.SignedCommandResponse
			result2 error
		})
	}
	fake.processCommandReturnsOnCall[i] = struct {
		result1 *token.SignedCommandResponse
		result2 error
	}{result1, result2}
}

func (fake *ProverClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.processCommandMutex.RLock()
	defer fake.processCommandMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ProverClient) recordInvocation(key string, args []interface{}) {
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