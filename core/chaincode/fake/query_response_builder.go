
package fake

import (
	"sync"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/chaincode"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

type QueryResponseBuilder struct {
	BuildQueryResponseStub        func(txContext *chaincode.TransactionContext, iter commonledger.ResultsIterator, iterID string, isPaginated bool, totalReturnLimit int32) (*pb.QueryResponse, error)
	buildQueryResponseMutex       sync.RWMutex
	buildQueryResponseArgsForCall []struct {
		txContext        *chaincode.TransactionContext
		iter             commonledger.ResultsIterator
		iterID           string
		isPaginated      bool
		totalReturnLimit int32
	}
	buildQueryResponseReturns struct {
		result1 *pb.QueryResponse
		result2 error
	}
	buildQueryResponseReturnsOnCall map[int]struct {
		result1 *pb.QueryResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *QueryResponseBuilder) BuildQueryResponse(txContext *chaincode.TransactionContext, iter commonledger.ResultsIterator, iterID string, isPaginated bool, totalReturnLimit int32) (*pb.QueryResponse, error) {
	fake.buildQueryResponseMutex.Lock()
	ret, specificReturn := fake.buildQueryResponseReturnsOnCall[len(fake.buildQueryResponseArgsForCall)]
	fake.buildQueryResponseArgsForCall = append(fake.buildQueryResponseArgsForCall, struct {
		txContext        *chaincode.TransactionContext
		iter             commonledger.ResultsIterator
		iterID           string
		isPaginated      bool
		totalReturnLimit int32
	}{txContext, iter, iterID, isPaginated, totalReturnLimit})
	fake.recordInvocation("BuildQueryResponse", []interface{}{txContext, iter, iterID, isPaginated, totalReturnLimit})
	fake.buildQueryResponseMutex.Unlock()
	if fake.BuildQueryResponseStub != nil {
		return fake.BuildQueryResponseStub(txContext, iter, iterID, isPaginated, totalReturnLimit)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.buildQueryResponseReturns.result1, fake.buildQueryResponseReturns.result2
}

func (fake *QueryResponseBuilder) BuildQueryResponseCallCount() int {
	fake.buildQueryResponseMutex.RLock()
	defer fake.buildQueryResponseMutex.RUnlock()
	return len(fake.buildQueryResponseArgsForCall)
}

func (fake *QueryResponseBuilder) BuildQueryResponseArgsForCall(i int) (*chaincode.TransactionContext, commonledger.ResultsIterator, string, bool, int32) {
	fake.buildQueryResponseMutex.RLock()
	defer fake.buildQueryResponseMutex.RUnlock()
	return fake.buildQueryResponseArgsForCall[i].txContext, fake.buildQueryResponseArgsForCall[i].iter, fake.buildQueryResponseArgsForCall[i].iterID, fake.buildQueryResponseArgsForCall[i].isPaginated, fake.buildQueryResponseArgsForCall[i].totalReturnLimit
}

func (fake *QueryResponseBuilder) BuildQueryResponseReturns(result1 *pb.QueryResponse, result2 error) {
	fake.BuildQueryResponseStub = nil
	fake.buildQueryResponseReturns = struct {
		result1 *pb.QueryResponse
		result2 error
	}{result1, result2}
}

func (fake *QueryResponseBuilder) BuildQueryResponseReturnsOnCall(i int, result1 *pb.QueryResponse, result2 error) {
	fake.BuildQueryResponseStub = nil
	if fake.buildQueryResponseReturnsOnCall == nil {
		fake.buildQueryResponseReturnsOnCall = make(map[int]struct {
			result1 *pb.QueryResponse
			result2 error
		})
	}
	fake.buildQueryResponseReturnsOnCall[i] = struct {
		result1 *pb.QueryResponse
		result2 error
	}{result1, result2}
}

func (fake *QueryResponseBuilder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildQueryResponseMutex.RLock()
	defer fake.buildQueryResponseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *QueryResponseBuilder) recordInvocation(key string, args []interface{}) {
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
