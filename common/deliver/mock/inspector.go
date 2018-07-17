
package mock

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/deliver"
	"golang.org/x/net/context"
)

type Inspector struct {
	InspectStub        func(context.Context, proto.Message) error
	inspectMutex       sync.RWMutex
	inspectArgsForCall []struct {
		arg1 context.Context
		arg2 proto.Message
	}
	inspectReturns struct {
		result1 error
	}
	inspectReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Inspector) Inspect(arg1 context.Context, arg2 proto.Message) error {
	fake.inspectMutex.Lock()
	ret, specificReturn := fake.inspectReturnsOnCall[len(fake.inspectArgsForCall)]
	fake.inspectArgsForCall = append(fake.inspectArgsForCall, struct {
		arg1 context.Context
		arg2 proto.Message
	}{arg1, arg2})
	fake.recordInvocation("Inspect", []interface{}{arg1, arg2})
	fake.inspectMutex.Unlock()
	if fake.InspectStub != nil {
		return fake.InspectStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.inspectReturns.result1
}

func (fake *Inspector) InspectCallCount() int {
	fake.inspectMutex.RLock()
	defer fake.inspectMutex.RUnlock()
	return len(fake.inspectArgsForCall)
}

func (fake *Inspector) InspectArgsForCall(i int) (context.Context, proto.Message) {
	fake.inspectMutex.RLock()
	defer fake.inspectMutex.RUnlock()
	return fake.inspectArgsForCall[i].arg1, fake.inspectArgsForCall[i].arg2
}

func (fake *Inspector) InspectReturns(result1 error) {
	fake.InspectStub = nil
	fake.inspectReturns = struct {
		result1 error
	}{result1}
}

func (fake *Inspector) InspectReturnsOnCall(i int, result1 error) {
	fake.InspectStub = nil
	if fake.inspectReturnsOnCall == nil {
		fake.inspectReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.inspectReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Inspector) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.inspectMutex.RLock()
	defer fake.inspectMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Inspector) recordInvocation(key string, args []interface{}) {
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

var _ deliver.Inspector = new(Inspector)
