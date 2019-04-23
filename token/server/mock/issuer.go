
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/server"
)

type Issuer struct {
	RequestIssueStub        func(tokensToIssue []*token.Token) (*token.TokenTransaction, error)
	requestIssueMutex       sync.RWMutex
	requestIssueArgsForCall []struct {
		tokensToIssue []*token.Token
	}
	requestIssueReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestIssueReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	RequestTokenOperationStub        func(op *token.TokenOperation) (*token.TokenTransaction, error)
	requestTokenOperationMutex       sync.RWMutex
	requestTokenOperationArgsForCall []struct {
		op *token.TokenOperation
	}
	requestTokenOperationReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestTokenOperationReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Issuer) RequestIssue(tokensToIssue []*token.Token) (*token.TokenTransaction, error) {
	var tokensToIssueCopy []*token.Token
	if tokensToIssue != nil {
		tokensToIssueCopy = make([]*token.Token, len(tokensToIssue))
		copy(tokensToIssueCopy, tokensToIssue)
	}
	fake.requestIssueMutex.Lock()
	ret, specificReturn := fake.requestIssueReturnsOnCall[len(fake.requestIssueArgsForCall)]
	fake.requestIssueArgsForCall = append(fake.requestIssueArgsForCall, struct {
		tokensToIssue []*token.Token
	}{tokensToIssueCopy})
	fake.recordInvocation("RequestIssue", []interface{}{tokensToIssueCopy})
	fake.requestIssueMutex.Unlock()
	if fake.RequestIssueStub != nil {
		return fake.RequestIssueStub(tokensToIssue)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.requestIssueReturns.result1, fake.requestIssueReturns.result2
}

func (fake *Issuer) RequestIssueCallCount() int {
	fake.requestIssueMutex.RLock()
	defer fake.requestIssueMutex.RUnlock()
	return len(fake.requestIssueArgsForCall)
}

func (fake *Issuer) RequestIssueArgsForCall(i int) []*token.Token {
	fake.requestIssueMutex.RLock()
	defer fake.requestIssueMutex.RUnlock()
	return fake.requestIssueArgsForCall[i].tokensToIssue
}

func (fake *Issuer) RequestIssueReturns(result1 *token.TokenTransaction, result2 error) {
	fake.RequestIssueStub = nil
	fake.requestIssueReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Issuer) RequestIssueReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.RequestIssueStub = nil
	if fake.requestIssueReturnsOnCall == nil {
		fake.requestIssueReturnsOnCall = make(map[int]struct {
			result1 *token.TokenTransaction
			result2 error
		})
	}
	fake.requestIssueReturnsOnCall[i] = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Issuer) RequestTokenOperation(op *token.TokenOperation) (*token.TokenTransaction, error) {
	fake.requestTokenOperationMutex.Lock()
	ret, specificReturn := fake.requestTokenOperationReturnsOnCall[len(fake.requestTokenOperationArgsForCall)]
	fake.requestTokenOperationArgsForCall = append(fake.requestTokenOperationArgsForCall, struct {
		op *token.TokenOperation
	}{op})
	fake.recordInvocation("RequestTokenOperation", []interface{}{op})
	fake.requestTokenOperationMutex.Unlock()
	if fake.RequestTokenOperationStub != nil {
		return fake.RequestTokenOperationStub(op)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.requestTokenOperationReturns.result1, fake.requestTokenOperationReturns.result2
}

func (fake *Issuer) RequestTokenOperationCallCount() int {
	fake.requestTokenOperationMutex.RLock()
	defer fake.requestTokenOperationMutex.RUnlock()
	return len(fake.requestTokenOperationArgsForCall)
}

func (fake *Issuer) RequestTokenOperationArgsForCall(i int) *token.TokenOperation {
	fake.requestTokenOperationMutex.RLock()
	defer fake.requestTokenOperationMutex.RUnlock()
	return fake.requestTokenOperationArgsForCall[i].op
}

func (fake *Issuer) RequestTokenOperationReturns(result1 *token.TokenTransaction, result2 error) {
	fake.RequestTokenOperationStub = nil
	fake.requestTokenOperationReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Issuer) RequestTokenOperationReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.RequestTokenOperationStub = nil
	if fake.requestTokenOperationReturnsOnCall == nil {
		fake.requestTokenOperationReturnsOnCall = make(map[int]struct {
			result1 *token.TokenTransaction
			result2 error
		})
	}
	fake.requestTokenOperationReturnsOnCall[i] = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Issuer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.requestIssueMutex.RLock()
	defer fake.requestIssueMutex.RUnlock()
	fake.requestTokenOperationMutex.RLock()
	defer fake.requestTokenOperationMutex.RUnlock()
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
