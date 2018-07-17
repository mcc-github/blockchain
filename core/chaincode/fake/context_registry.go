
package fake

import (
	"sync"

	chaincode_test "github.com/mcc-github/blockchain/core/chaincode"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"golang.org/x/net/context"
)

type ContextRegistry struct {
	CreateStub        func(ctx context.Context, chainID, txID string, signedProp *pb.SignedProposal, proposal *pb.Proposal) (*chaincode_test.TransactionContext, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		ctx        context.Context
		chainID    string
		txID       string
		signedProp *pb.SignedProposal
		proposal   *pb.Proposal
	}
	createReturns struct {
		result1 *chaincode_test.TransactionContext
		result2 error
	}
	createReturnsOnCall map[int]struct {
		result1 *chaincode_test.TransactionContext
		result2 error
	}
	GetStub        func(chainID, txID string) *chaincode_test.TransactionContext
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		chainID string
		txID    string
	}
	getReturns struct {
		result1 *chaincode_test.TransactionContext
	}
	getReturnsOnCall map[int]struct {
		result1 *chaincode_test.TransactionContext
	}
	DeleteStub        func(chainID, txID string)
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		chainID string
		txID    string
	}
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ContextRegistry) Create(ctx context.Context, chainID string, txID string, signedProp *pb.SignedProposal, proposal *pb.Proposal) (*chaincode_test.TransactionContext, error) {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		ctx        context.Context
		chainID    string
		txID       string
		signedProp *pb.SignedProposal
		proposal   *pb.Proposal
	}{ctx, chainID, txID, signedProp, proposal})
	fake.recordInvocation("Create", []interface{}{ctx, chainID, txID, signedProp, proposal})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(ctx, chainID, txID, signedProp, proposal)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createReturns.result1, fake.createReturns.result2
}

func (fake *ContextRegistry) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *ContextRegistry) CreateArgsForCall(i int) (context.Context, string, string, *pb.SignedProposal, *pb.Proposal) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].ctx, fake.createArgsForCall[i].chainID, fake.createArgsForCall[i].txID, fake.createArgsForCall[i].signedProp, fake.createArgsForCall[i].proposal
}

func (fake *ContextRegistry) CreateReturns(result1 *chaincode_test.TransactionContext, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 *chaincode_test.TransactionContext
		result2 error
	}{result1, result2}
}

func (fake *ContextRegistry) CreateReturnsOnCall(i int, result1 *chaincode_test.TransactionContext, result2 error) {
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 *chaincode_test.TransactionContext
			result2 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 *chaincode_test.TransactionContext
		result2 error
	}{result1, result2}
}

func (fake *ContextRegistry) Get(chainID string, txID string) *chaincode_test.TransactionContext {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		chainID string
		txID    string
	}{chainID, txID})
	fake.recordInvocation("Get", []interface{}{chainID, txID})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(chainID, txID)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getReturns.result1
}

func (fake *ContextRegistry) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *ContextRegistry) GetArgsForCall(i int) (string, string) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].chainID, fake.getArgsForCall[i].txID
}

func (fake *ContextRegistry) GetReturns(result1 *chaincode_test.TransactionContext) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *chaincode_test.TransactionContext
	}{result1}
}

func (fake *ContextRegistry) GetReturnsOnCall(i int, result1 *chaincode_test.TransactionContext) {
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 *chaincode_test.TransactionContext
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 *chaincode_test.TransactionContext
	}{result1}
}

func (fake *ContextRegistry) Delete(chainID string, txID string) {
	fake.deleteMutex.Lock()
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		chainID string
		txID    string
	}{chainID, txID})
	fake.recordInvocation("Delete", []interface{}{chainID, txID})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		fake.DeleteStub(chainID, txID)
	}
}

func (fake *ContextRegistry) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *ContextRegistry) DeleteArgsForCall(i int) (string, string) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].chainID, fake.deleteArgsForCall[i].txID
}

func (fake *ContextRegistry) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *ContextRegistry) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *ContextRegistry) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ContextRegistry) recordInvocation(key string, args []interface{}) {
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
