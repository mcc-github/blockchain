
package fakes

import (
	"sync"

	"github.com/mcc-github/blockchain/common/grpclogging/testpb"
	"golang.org/x/net/context"
)

type EchoServiceServer struct {
	EchoStub        func(context.Context, *testpb.Message) (*testpb.Message, error)
	echoMutex       sync.RWMutex
	echoArgsForCall []struct {
		arg1 context.Context
		arg2 *testpb.Message
	}
	echoReturns struct {
		result1 *testpb.Message
		result2 error
	}
	echoReturnsOnCall map[int]struct {
		result1 *testpb.Message
		result2 error
	}
	EchoStreamStub        func(testpb.EchoService_EchoStreamServer) error
	echoStreamMutex       sync.RWMutex
	echoStreamArgsForCall []struct {
		arg1 testpb.EchoService_EchoStreamServer
	}
	echoStreamReturns struct {
		result1 error
	}
	echoStreamReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *EchoServiceServer) Echo(arg1 context.Context, arg2 *testpb.Message) (*testpb.Message, error) {
	fake.echoMutex.Lock()
	ret, specificReturn := fake.echoReturnsOnCall[len(fake.echoArgsForCall)]
	fake.echoArgsForCall = append(fake.echoArgsForCall, struct {
		arg1 context.Context
		arg2 *testpb.Message
	}{arg1, arg2})
	fake.recordInvocation("Echo", []interface{}{arg1, arg2})
	fake.echoMutex.Unlock()
	if fake.EchoStub != nil {
		return fake.EchoStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.echoReturns.result1, fake.echoReturns.result2
}

func (fake *EchoServiceServer) EchoCallCount() int {
	fake.echoMutex.RLock()
	defer fake.echoMutex.RUnlock()
	return len(fake.echoArgsForCall)
}

func (fake *EchoServiceServer) EchoArgsForCall(i int) (context.Context, *testpb.Message) {
	fake.echoMutex.RLock()
	defer fake.echoMutex.RUnlock()
	return fake.echoArgsForCall[i].arg1, fake.echoArgsForCall[i].arg2
}

func (fake *EchoServiceServer) EchoReturns(result1 *testpb.Message, result2 error) {
	fake.EchoStub = nil
	fake.echoReturns = struct {
		result1 *testpb.Message
		result2 error
	}{result1, result2}
}

func (fake *EchoServiceServer) EchoReturnsOnCall(i int, result1 *testpb.Message, result2 error) {
	fake.EchoStub = nil
	if fake.echoReturnsOnCall == nil {
		fake.echoReturnsOnCall = make(map[int]struct {
			result1 *testpb.Message
			result2 error
		})
	}
	fake.echoReturnsOnCall[i] = struct {
		result1 *testpb.Message
		result2 error
	}{result1, result2}
}

func (fake *EchoServiceServer) EchoStream(arg1 testpb.EchoService_EchoStreamServer) error {
	fake.echoStreamMutex.Lock()
	ret, specificReturn := fake.echoStreamReturnsOnCall[len(fake.echoStreamArgsForCall)]
	fake.echoStreamArgsForCall = append(fake.echoStreamArgsForCall, struct {
		arg1 testpb.EchoService_EchoStreamServer
	}{arg1})
	fake.recordInvocation("EchoStream", []interface{}{arg1})
	fake.echoStreamMutex.Unlock()
	if fake.EchoStreamStub != nil {
		return fake.EchoStreamStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.echoStreamReturns.result1
}

func (fake *EchoServiceServer) EchoStreamCallCount() int {
	fake.echoStreamMutex.RLock()
	defer fake.echoStreamMutex.RUnlock()
	return len(fake.echoStreamArgsForCall)
}

func (fake *EchoServiceServer) EchoStreamArgsForCall(i int) testpb.EchoService_EchoStreamServer {
	fake.echoStreamMutex.RLock()
	defer fake.echoStreamMutex.RUnlock()
	return fake.echoStreamArgsForCall[i].arg1
}

func (fake *EchoServiceServer) EchoStreamReturns(result1 error) {
	fake.EchoStreamStub = nil
	fake.echoStreamReturns = struct {
		result1 error
	}{result1}
}

func (fake *EchoServiceServer) EchoStreamReturnsOnCall(i int, result1 error) {
	fake.EchoStreamStub = nil
	if fake.echoStreamReturnsOnCall == nil {
		fake.echoStreamReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.echoStreamReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *EchoServiceServer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.echoMutex.RLock()
	defer fake.echoMutex.RUnlock()
	fake.echoStreamMutex.RLock()
	defer fake.echoStreamMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *EchoServiceServer) recordInvocation(key string, args []interface{}) {
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