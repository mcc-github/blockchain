
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"golang.org/x/net/context"
)

type Invoker struct {
	InvokeStub        func(ctxt context.Context, cccid *ccprovider.CCContext, spec ccprovider.ChaincodeSpecGetter) (*pb.ChaincodeMessage, error)
	invokeMutex       sync.RWMutex
	invokeArgsForCall []struct {
		ctxt  context.Context
		cccid *ccprovider.CCContext
		spec  ccprovider.ChaincodeSpecGetter
	}
	invokeReturns struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}
	invokeReturnsOnCall map[int]struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Invoker) Invoke(ctxt context.Context, cccid *ccprovider.CCContext, spec ccprovider.ChaincodeSpecGetter) (*pb.ChaincodeMessage, error) {
	fake.invokeMutex.Lock()
	ret, specificReturn := fake.invokeReturnsOnCall[len(fake.invokeArgsForCall)]
	fake.invokeArgsForCall = append(fake.invokeArgsForCall, struct {
		ctxt  context.Context
		cccid *ccprovider.CCContext
		spec  ccprovider.ChaincodeSpecGetter
	}{ctxt, cccid, spec})
	fake.recordInvocation("Invoke", []interface{}{ctxt, cccid, spec})
	fake.invokeMutex.Unlock()
	if fake.InvokeStub != nil {
		return fake.InvokeStub(ctxt, cccid, spec)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.invokeReturns.result1, fake.invokeReturns.result2
}

func (fake *Invoker) InvokeCallCount() int {
	fake.invokeMutex.RLock()
	defer fake.invokeMutex.RUnlock()
	return len(fake.invokeArgsForCall)
}

func (fake *Invoker) InvokeArgsForCall(i int) (context.Context, *ccprovider.CCContext, ccprovider.ChaincodeSpecGetter) {
	fake.invokeMutex.RLock()
	defer fake.invokeMutex.RUnlock()
	return fake.invokeArgsForCall[i].ctxt, fake.invokeArgsForCall[i].cccid, fake.invokeArgsForCall[i].spec
}

func (fake *Invoker) InvokeReturns(result1 *pb.ChaincodeMessage, result2 error) {
	fake.InvokeStub = nil
	fake.invokeReturns = struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}{result1, result2}
}

func (fake *Invoker) InvokeReturnsOnCall(i int, result1 *pb.ChaincodeMessage, result2 error) {
	fake.InvokeStub = nil
	if fake.invokeReturnsOnCall == nil {
		fake.invokeReturnsOnCall = make(map[int]struct {
			result1 *pb.ChaincodeMessage
			result2 error
		})
	}
	fake.invokeReturnsOnCall[i] = struct {
		result1 *pb.ChaincodeMessage
		result2 error
	}{result1, result2}
}

func (fake *Invoker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.invokeMutex.RLock()
	defer fake.invokeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Invoker) recordInvocation(key string, args []interface{}) {
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
