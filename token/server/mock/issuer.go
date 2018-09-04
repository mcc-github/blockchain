
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/server"
)

type Issuer struct {
	RequestImportStub        func(tokensToIssue []*token.TokenToIssue) (*token.TokenTransaction, error)
	requestImportMutex       sync.RWMutex
	requestImportArgsForCall []struct {
		tokensToIssue []*token.TokenToIssue
	}
	requestImportReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestImportReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Issuer) RequestImport(tokensToIssue []*token.TokenToIssue) (*token.TokenTransaction, error) {
	var tokensToIssueCopy []*token.TokenToIssue
	if tokensToIssue != nil {
		tokensToIssueCopy = make([]*token.TokenToIssue, len(tokensToIssue))
		copy(tokensToIssueCopy, tokensToIssue)
	}
	fake.requestImportMutex.Lock()
	ret, specificReturn := fake.requestImportReturnsOnCall[len(fake.requestImportArgsForCall)]
	fake.requestImportArgsForCall = append(fake.requestImportArgsForCall, struct {
		tokensToIssue []*token.TokenToIssue
	}{tokensToIssueCopy})
	fake.recordInvocation("RequestImport", []interface{}{tokensToIssueCopy})
	fake.requestImportMutex.Unlock()
	if fake.RequestImportStub != nil {
		return fake.RequestImportStub(tokensToIssue)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.requestImportReturns.result1, fake.requestImportReturns.result2
}

func (fake *Issuer) RequestImportCallCount() int {
	fake.requestImportMutex.RLock()
	defer fake.requestImportMutex.RUnlock()
	return len(fake.requestImportArgsForCall)
}

func (fake *Issuer) RequestImportArgsForCall(i int) []*token.TokenToIssue {
	fake.requestImportMutex.RLock()
	defer fake.requestImportMutex.RUnlock()
	return fake.requestImportArgsForCall[i].tokensToIssue
}

func (fake *Issuer) RequestImportReturns(result1 *token.TokenTransaction, result2 error) {
	fake.RequestImportStub = nil
	fake.requestImportReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Issuer) RequestImportReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.RequestImportStub = nil
	if fake.requestImportReturnsOnCall == nil {
		fake.requestImportReturnsOnCall = make(map[int]struct {
			result1 *token.TokenTransaction
			result2 error
		})
	}
	fake.requestImportReturnsOnCall[i] = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Issuer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.requestImportMutex.RLock()
	defer fake.requestImportMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Issuer) recordInvocation(key string, args []interface{}) {
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

var _ server.Issuer = new(Issuer)
