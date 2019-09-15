
package mock

import (
	"sync"

	"github.com/golang/protobuf/proto"
)

type Protobuf struct {
	MarshalStub        func(proto.Message) ([]byte, error)
	marshalMutex       sync.RWMutex
	marshalArgsForCall []struct {
		arg1 proto.Message
	}
	marshalReturns struct {
		result1 []byte
		result2 error
	}
	marshalReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	UnmarshalStub        func([]byte, proto.Message) error
	unmarshalMutex       sync.RWMutex
	unmarshalArgsForCall []struct {
		arg1 []byte
		arg2 proto.Message
	}
	unmarshalReturns struct {
		result1 error
	}
	unmarshalReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Protobuf) Marshal(arg1 proto.Message) ([]byte, error) {
	fake.marshalMutex.Lock()
	ret, specificReturn := fake.marshalReturnsOnCall[len(fake.marshalArgsForCall)]
	fake.marshalArgsForCall = append(fake.marshalArgsForCall, struct {
		arg1 proto.Message
	}{arg1})
	fake.recordInvocation("Marshal", []interface{}{arg1})
	fake.marshalMutex.Unlock()
	if fake.MarshalStub != nil {
		return fake.MarshalStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.marshalReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Protobuf) MarshalCallCount() int {
	fake.marshalMutex.RLock()
	defer fake.marshalMutex.RUnlock()
	return len(fake.marshalArgsForCall)
}

func (fake *Protobuf) MarshalCalls(stub func(proto.Message) ([]byte, error)) {
	fake.marshalMutex.Lock()
	defer fake.marshalMutex.Unlock()
	fake.MarshalStub = stub
}

func (fake *Protobuf) MarshalArgsForCall(i int) proto.Message {
	fake.marshalMutex.RLock()
	defer fake.marshalMutex.RUnlock()
	argsForCall := fake.marshalArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Protobuf) MarshalReturns(result1 []byte, result2 error) {
	fake.marshalMutex.Lock()
	defer fake.marshalMutex.Unlock()
	fake.MarshalStub = nil
	fake.marshalReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Protobuf) MarshalReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.marshalMutex.Lock()
	defer fake.marshalMutex.Unlock()
	fake.MarshalStub = nil
	if fake.marshalReturnsOnCall == nil {
		fake.marshalReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.marshalReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Protobuf) Unmarshal(arg1 []byte, arg2 proto.Message) error {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.unmarshalMutex.Lock()
	ret, specificReturn := fake.unmarshalReturnsOnCall[len(fake.unmarshalArgsForCall)]
	fake.unmarshalArgsForCall = append(fake.unmarshalArgsForCall, struct {
		arg1 []byte
		arg2 proto.Message
	}{arg1Copy, arg2})
	fake.recordInvocation("Unmarshal", []interface{}{arg1Copy, arg2})
	fake.unmarshalMutex.Unlock()
	if fake.UnmarshalStub != nil {
		return fake.UnmarshalStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.unmarshalReturns
	return fakeReturns.result1
}

func (fake *Protobuf) UnmarshalCallCount() int {
	fake.unmarshalMutex.RLock()
	defer fake.unmarshalMutex.RUnlock()
	return len(fake.unmarshalArgsForCall)
}

func (fake *Protobuf) UnmarshalCalls(stub func([]byte, proto.Message) error) {
	fake.unmarshalMutex.Lock()
	defer fake.unmarshalMutex.Unlock()
	fake.UnmarshalStub = stub
}

func (fake *Protobuf) UnmarshalArgsForCall(i int) ([]byte, proto.Message) {
	fake.unmarshalMutex.RLock()
	defer fake.unmarshalMutex.RUnlock()
	argsForCall := fake.unmarshalArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Protobuf) UnmarshalReturns(result1 error) {
	fake.unmarshalMutex.Lock()
	defer fake.unmarshalMutex.Unlock()
	fake.UnmarshalStub = nil
	fake.unmarshalReturns = struct {
		result1 error
	}{result1}
}

func (fake *Protobuf) UnmarshalReturnsOnCall(i int, result1 error) {
	fake.unmarshalMutex.Lock()
	defer fake.unmarshalMutex.Unlock()
	fake.UnmarshalStub = nil
	if fake.unmarshalReturnsOnCall == nil {
		fake.unmarshalReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.unmarshalReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Protobuf) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.marshalMutex.RLock()
	defer fake.marshalMutex.RUnlock()
	fake.unmarshalMutex.RLock()
	defer fake.unmarshalMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Protobuf) recordInvocation(key string, args []interface{}) {
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
