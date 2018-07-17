
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"golang.org/x/net/context"
)

type ChaincodeDefinitionGetter struct {
	GetChaincodeDefinitionStub        func(ctxt context.Context, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, chainID string, chaincodeID string) (ccprovider.ChaincodeDefinition, error)
	getChaincodeDefinitionMutex       sync.RWMutex
	getChaincodeDefinitionArgsForCall []struct {
		ctxt        context.Context
		txid        string
		signedProp  *pb.SignedProposal
		prop        *pb.Proposal
		chainID     string
		chaincodeID string
	}
	getChaincodeDefinitionReturns struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}
	getChaincodeDefinitionReturnsOnCall map[int]struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeDefinitionGetter) GetChaincodeDefinition(ctxt context.Context, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, chainID string, chaincodeID string) (ccprovider.ChaincodeDefinition, error) {
	fake.getChaincodeDefinitionMutex.Lock()
	ret, specificReturn := fake.getChaincodeDefinitionReturnsOnCall[len(fake.getChaincodeDefinitionArgsForCall)]
	fake.getChaincodeDefinitionArgsForCall = append(fake.getChaincodeDefinitionArgsForCall, struct {
		ctxt        context.Context
		txid        string
		signedProp  *pb.SignedProposal
		prop        *pb.Proposal
		chainID     string
		chaincodeID string
	}{ctxt, txid, signedProp, prop, chainID, chaincodeID})
	fake.recordInvocation("GetChaincodeDefinition", []interface{}{ctxt, txid, signedProp, prop, chainID, chaincodeID})
	fake.getChaincodeDefinitionMutex.Unlock()
	if fake.GetChaincodeDefinitionStub != nil {
		return fake.GetChaincodeDefinitionStub(ctxt, txid, signedProp, prop, chainID, chaincodeID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChaincodeDefinitionReturns.result1, fake.getChaincodeDefinitionReturns.result2
}

func (fake *ChaincodeDefinitionGetter) GetChaincodeDefinitionCallCount() int {
	fake.getChaincodeDefinitionMutex.RLock()
	defer fake.getChaincodeDefinitionMutex.RUnlock()
	return len(fake.getChaincodeDefinitionArgsForCall)
}

func (fake *ChaincodeDefinitionGetter) GetChaincodeDefinitionArgsForCall(i int) (context.Context, string, *pb.SignedProposal, *pb.Proposal, string, string) {
	fake.getChaincodeDefinitionMutex.RLock()
	defer fake.getChaincodeDefinitionMutex.RUnlock()
	return fake.getChaincodeDefinitionArgsForCall[i].ctxt, fake.getChaincodeDefinitionArgsForCall[i].txid, fake.getChaincodeDefinitionArgsForCall[i].signedProp, fake.getChaincodeDefinitionArgsForCall[i].prop, fake.getChaincodeDefinitionArgsForCall[i].chainID, fake.getChaincodeDefinitionArgsForCall[i].chaincodeID
}

func (fake *ChaincodeDefinitionGetter) GetChaincodeDefinitionReturns(result1 ccprovider.ChaincodeDefinition, result2 error) {
	fake.GetChaincodeDefinitionStub = nil
	fake.getChaincodeDefinitionReturns = struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeDefinitionGetter) GetChaincodeDefinitionReturnsOnCall(i int, result1 ccprovider.ChaincodeDefinition, result2 error) {
	fake.GetChaincodeDefinitionStub = nil
	if fake.getChaincodeDefinitionReturnsOnCall == nil {
		fake.getChaincodeDefinitionReturnsOnCall = make(map[int]struct {
			result1 ccprovider.ChaincodeDefinition
			result2 error
		})
	}
	fake.getChaincodeDefinitionReturnsOnCall[i] = struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeDefinitionGetter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChaincodeDefinitionMutex.RLock()
	defer fake.getChaincodeDefinitionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeDefinitionGetter) recordInvocation(key string, args []interface{}) {
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
