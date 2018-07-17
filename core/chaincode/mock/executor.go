
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"golang.org/x/net/context"
)

type Executor struct {
	ExecuteStub        func(ctxt context.Context, cccid *ccprovider.CCContext, cis ccprovider.ChaincodeSpecGetter) (*pb.Response, *pb.ChaincodeEvent, error)
	executeMutex       sync.RWMutex
	executeArgsForCall []struct {
		ctxt  context.Context
		cccid *ccprovider.CCContext
		cis   ccprovider.ChaincodeSpecGetter
	}
	executeReturns struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}
	executeReturnsOnCall map[int]struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Executor) Execute(ctxt context.Context, cccid *ccprovider.CCContext, cis ccprovider.ChaincodeSpecGetter) (*pb.Response, *pb.ChaincodeEvent, error) {
	fake.executeMutex.Lock()
	ret, specificReturn := fake.executeReturnsOnCall[len(fake.executeArgsForCall)]
	fake.executeArgsForCall = append(fake.executeArgsForCall, struct {
		ctxt  context.Context
		cccid *ccprovider.CCContext
		cis   ccprovider.ChaincodeSpecGetter
	}{ctxt, cccid, cis})
	fake.recordInvocation("Execute", []interface{}{ctxt, cccid, cis})
	fake.executeMutex.Unlock()
	if fake.ExecuteStub != nil {
		return fake.ExecuteStub(ctxt, cccid, cis)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.executeReturns.result1, fake.executeReturns.result2, fake.executeReturns.result3
}

func (fake *Executor) ExecuteCallCount() int {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return len(fake.executeArgsForCall)
}

func (fake *Executor) ExecuteArgsForCall(i int) (context.Context, *ccprovider.CCContext, ccprovider.ChaincodeSpecGetter) {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return fake.executeArgsForCall[i].ctxt, fake.executeArgsForCall[i].cccid, fake.executeArgsForCall[i].cis
}

func (fake *Executor) ExecuteReturns(result1 *pb.Response, result2 *pb.ChaincodeEvent, result3 error) {
	fake.ExecuteStub = nil
	fake.executeReturns = struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}{result1, result2, result3}
}

func (fake *Executor) ExecuteReturnsOnCall(i int, result1 *pb.Response, result2 *pb.ChaincodeEvent, result3 error) {
	fake.ExecuteStub = nil
	if fake.executeReturnsOnCall == nil {
		fake.executeReturnsOnCall = make(map[int]struct {
			result1 *pb.Response
			result2 *pb.ChaincodeEvent
			result3 error
		})
	}
	fake.executeReturnsOnCall[i] = struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}{result1, result2, result3}
}

func (fake *Executor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Executor) recordInvocation(key string, args []interface{}) {
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
