
package fake

import (
	"sync"

	chaincode_test "github.com/mcc-github/blockchain/core/chaincode"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

type MessageHandler struct {
	HandleStub        func(*pb.ChaincodeMessage, *chaincode_test.TransactionContext) (*pb.ChaincodeMessage, error)
	handleMutex       sync.RWMutex
	handleArgsForCall []struct {
		arg1 *pb.ChaincodeMessage
		arg2 *chaincode_test.TransactionContext
	}
	handleReturns struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}
	handleReturnsOnCall map[int]struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MessageHandler) Handle(arg1 *pb.ChaincodeMessage, arg2 *chaincode_test.TransactionContext) (*pb.ChaincodeMessage, error) {
	fake.handleMutex.Lock()
	ret, specificReturn := fake.handleReturnsOnCall[len(fake.handleArgsForCall)]
	fake.handleArgsForCall = append(fake.handleArgsForCall, struct {
		arg1 *pb.ChaincodeMessage
		arg2 *chaincode_test.TransactionContext
	}{arg1, arg2})
	fake.recordInvocation("Handle", []interface{}{arg1, arg2})
	fake.handleMutex.Unlock()
	if fake.HandleStub != nil {
		return fake.HandleStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.handleReturns.result1, fake.handleReturns.result2
}

func (fake *MessageHandler) HandleCallCount() int {
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	return len(fake.handleArgsForCall)
}

func (fake *MessageHandler) HandleArgsForCall(i int) (*pb.ChaincodeMessage, *chaincode_test.TransactionContext) {
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	return fake.handleArgsForCall[i].arg1, fake.handleArgsForCall[i].arg2
}

func (fake *MessageHandler) HandleReturns(result1 *pb.ChaincodeMessage, result2 error) {
	fake.HandleStub = nil
	fake.handleReturns = struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}{result1, result2}
}

func (fake *MessageHandler) HandleReturnsOnCall(i int, result1 *pb.ChaincodeMessage, result2 error) {
	fake.HandleStub = nil
	if fake.handleReturnsOnCall == nil {
		fake.handleReturnsOnCall = make(map[int]struct {
			result1 *pb.ChaincodeMessage
			result2 error
		})
	}
	fake.handleReturnsOnCall[i] = struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}{result1, result2}
}

func (fake *MessageHandler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MessageHandler) recordInvocation(key string, args []interface{}) {
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
