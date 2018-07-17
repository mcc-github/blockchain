
package mock

import (
	"sync"

	pb "github.com/mcc-github/blockchain/protos/peer"
)

type ChaincodeStream struct {
	SendStub        func(*pb.ChaincodeMessage) error
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		arg1 *pb.ChaincodeMessage
	}
	sendReturns struct {
		result1 error
	}
	sendReturnsOnCall map[int]struct {
		result1 error
	}
	RecvStub        func() (*pb.ChaincodeMessage, error)
	recvMutex       sync.RWMutex
	recvArgsForCall []struct{}
	recvReturns     struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}
	recvReturnsOnCall map[int]struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeStream) Send(arg1 *pb.ChaincodeMessage) error {
	fake.sendMutex.Lock()
	ret, specificReturn := fake.sendReturnsOnCall[len(fake.sendArgsForCall)]
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		arg1 *pb.ChaincodeMessage
	}{arg1})
	fake.recordInvocation("Send", []interface{}{arg1})
	fake.sendMutex.Unlock()
	if fake.SendStub != nil {
		return fake.SendStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.sendReturns.result1
}

func (fake *ChaincodeStream) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *ChaincodeStream) SendArgsForCall(i int) *pb.ChaincodeMessage {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return fake.sendArgsForCall[i].arg1
}

func (fake *ChaincodeStream) SendReturns(result1 error) {
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStream) SendReturnsOnCall(i int, result1 error) {
	fake.SendStub = nil
	if fake.sendReturnsOnCall == nil {
		fake.sendReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStream) Recv() (*pb.ChaincodeMessage, error) {
	fake.recvMutex.Lock()
	ret, specificReturn := fake.recvReturnsOnCall[len(fake.recvArgsForCall)]
	fake.recvArgsForCall = append(fake.recvArgsForCall, struct{}{})
	fake.recordInvocation("Recv", []interface{}{})
	fake.recvMutex.Unlock()
	if fake.RecvStub != nil {
		return fake.RecvStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.recvReturns.result1, fake.recvReturns.result2
}

func (fake *ChaincodeStream) RecvCallCount() int {
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	return len(fake.recvArgsForCall)
}

func (fake *ChaincodeStream) RecvReturns(result1 *pb.ChaincodeMessage, result2 error) {
	fake.RecvStub = nil
	fake.recvReturns = struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStream) RecvReturnsOnCall(i int, result1 *pb.ChaincodeMessage, result2 error) {
	fake.RecvStub = nil
	if fake.recvReturnsOnCall == nil {
		fake.recvReturnsOnCall = make(map[int]struct {
			result1 *pb.ChaincodeMessage
			result2 error
		})
	}
	fake.recvReturnsOnCall[i] = struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStream) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeStream) recordInvocation(key string, args []interface{}) {
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
