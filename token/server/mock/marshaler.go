
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/server"
)

type Marshaler struct {
	MarshalCommandResponseStub        func([]byte, interface{}) (*token.SignedCommandResponse, error)
	marshalCommandResponseMutex       sync.RWMutex
	marshalCommandResponseArgsForCall []struct {
		arg1 []byte
		arg2 interface{}
	}
	marshalCommandResponseReturns struct {
		result1 *token.SignedCommandResponse
		result2 error
	}
	marshalCommandResponseReturnsOnCall map[int]struct {
		result1 *token.SignedCommandResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Marshaler) MarshalCommandResponse(arg1 []byte, arg2 interface{}) (*token.SignedCommandResponse, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.marshalCommandResponseMutex.Lock()
	ret, specificReturn := fake.marshalCommandResponseReturnsOnCall[len(fake.marshalCommandResponseArgsForCall)]
	fake.marshalCommandResponseArgsForCall = append(fake.marshalCommandResponseArgsForCall, struct {
		arg1 []byte
		arg2 interface{}
	}{arg1Copy, arg2})
	fake.recordInvocation("MarshalCommandResponse", []interface{}{arg1Copy, arg2})
	fake.marshalCommandResponseMutex.Unlock()
	if fake.MarshalCommandResponseStub != nil {
		return fake.MarshalCommandResponseStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.marshalCommandResponseReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Marshaler) MarshalCommandResponseCallCount() int {
	fake.marshalCommandResponseMutex.RLock()
	defer fake.marshalCommandResponseMutex.RUnlock()
	return len(fake.marshalCommandResponseArgsForCall)
}

func (fake *Marshaler) MarshalCommandResponseCalls(stub func([]byte, interface{}) (*token.SignedCommandResponse, error)) {
	fake.marshalCommandResponseMutex.Lock()
	defer fake.marshalCommandResponseMutex.Unlock()
	fake.MarshalCommandResponseStub = stub
}

func (fake *Marshaler) MarshalCommandResponseArgsForCall(i int) ([]byte, interface{}) {
	fake.marshalCommandResponseMutex.RLock()
	defer fake.marshalCommandResponseMutex.RUnlock()
	argsForCall := fake.marshalCommandResponseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Marshaler) MarshalCommandResponseReturns(result1 *token.SignedCommandResponse, result2 error) {
	fake.marshalCommandResponseMutex.Lock()
	defer fake.marshalCommandResponseMutex.Unlock()
	fake.MarshalCommandResponseStub = nil
	fake.marshalCommandResponseReturns = struct {
		result1 *token.SignedCommandResponse
		result2 error
	}{result1, result2}
}

func (fake *Marshaler) MarshalCommandResponseReturnsOnCall(i int, result1 *token.SignedCommandResponse, result2 error) {
	fake.marshalCommandResponseMutex.Lock()
	defer fake.marshalCommandResponseMutex.Unlock()
	fake.MarshalCommandResponseStub = nil
	if fake.marshalCommandResponseReturnsOnCall == nil {
		fake.marshalCommandResponseReturnsOnCall = make(map[int]struct {
			result1 *token.SignedCommandResponse
			result2 error
		})
	}
	fake.marshalCommandResponseReturnsOnCall[i] = struct {
		result1 *token.SignedCommandResponse
		result2 error
	}{result1, result2}
}

func (fake *Marshaler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.marshalCommandResponseMutex.RLock()
	defer fake.marshalCommandResponseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Marshaler) recordInvocation(key string, args []interface{}) {
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

var _ server.Marshaler = new(Marshaler)
