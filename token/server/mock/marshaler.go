
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/server"
)

type Marshaler struct {
	MarshalCommandResponseStub        func(command []byte, responsePayload interface{}) (*token.SignedCommandResponse, error)
	marshalCommandResponseMutex       sync.RWMutex
	marshalCommandResponseArgsForCall []struct {
		command         []byte
		responsePayload interface{}
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

func (fake *Marshaler) MarshalCommandResponse(command []byte, responsePayload interface{}) (*token.SignedCommandResponse, error) {
	var commandCopy []byte
	if command != nil {
		commandCopy = make([]byte, len(command))
		copy(commandCopy, command)
	}
	fake.marshalCommandResponseMutex.Lock()
	ret, specificReturn := fake.marshalCommandResponseReturnsOnCall[len(fake.marshalCommandResponseArgsForCall)]
	fake.marshalCommandResponseArgsForCall = append(fake.marshalCommandResponseArgsForCall, struct {
		command         []byte
		responsePayload interface{}
	}{commandCopy, responsePayload})
	fake.recordInvocation("MarshalCommandResponse", []interface{}{commandCopy, responsePayload})
	fake.marshalCommandResponseMutex.Unlock()
	if fake.MarshalCommandResponseStub != nil {
		return fake.MarshalCommandResponseStub(command, responsePayload)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.marshalCommandResponseReturns.result1, fake.marshalCommandResponseReturns.result2
}

func (fake *Marshaler) MarshalCommandResponseCallCount() int {
	fake.marshalCommandResponseMutex.RLock()
	defer fake.marshalCommandResponseMutex.RUnlock()
	return len(fake.marshalCommandResponseArgsForCall)
}

func (fake *Marshaler) MarshalCommandResponseArgsForCall(i int) ([]byte, interface{}) {
	fake.marshalCommandResponseMutex.RLock()
	defer fake.marshalCommandResponseMutex.RUnlock()
	return fake.marshalCommandResponseArgsForCall[i].command, fake.marshalCommandResponseArgsForCall[i].responsePayload
}

func (fake *Marshaler) MarshalCommandResponseReturns(result1 *token.SignedCommandResponse, result2 error) {
	fake.MarshalCommandResponseStub = nil
	fake.marshalCommandResponseReturns = struct {
		result1 *token.SignedCommandResponse
		result2 error
	}{result1, result2}
}

func (fake *Marshaler) MarshalCommandResponseReturnsOnCall(i int, result1 *token.SignedCommandResponse, result2 error) {
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
