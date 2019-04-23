
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/container"
)

type ExitedFunc struct {
	Stub        func(int, error)
	mutex       sync.RWMutex
	argsForCall []struct {
		arg1 int
		arg2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ExitedFunc) Spy(arg1 int, arg2 error) {
	fake.mutex.Lock()
	fake.argsForCall = append(fake.argsForCall, struct {
		arg1 int
		arg2 error
	}{arg1, arg2})
	fake.recordInvocation("ExitedFunc", []interface{}{arg1, arg2})
	fake.mutex.Unlock()
	if fake.Stub != nil {
		fake.Stub(arg1, arg2)
	}
}

func (fake *ExitedFunc) CallCount() int {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return len(fake.argsForCall)
}

func (fake *ExitedFunc) Calls(stub func(int, error)) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = stub
}

func (fake *ExitedFunc) ArgsForCall(i int) (int, error) {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return fake.argsForCall[i].arg1, fake.argsForCall[i].arg2
}

func (fake *ExitedFunc) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ExitedFunc) recordInvocation(key string, args []interface{}) {
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

var _ container.ExitedFunc = new(ExitedFunc).Spy
