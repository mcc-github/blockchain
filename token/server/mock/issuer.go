
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/server"
)

type Issuer struct {
	RequestIssueStub        func([]*token.Token) (*token.TokenTransaction, error)
	requestIssueMutex       sync.RWMutex
	requestIssueArgsForCall []struct {
		arg1 []*token.Token
	}
	requestIssueReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestIssueReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	RequestTokenOperationStub        func(*token.TokenOperation) (*token.TokenTransaction, error)
	requestTokenOperationMutex       sync.RWMutex
	requestTokenOperationArgsForCall []struct {
		arg1 *token.TokenOperation
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

func (fake *Issuer) RequestIssue(arg1 []*token.Token) (*token.TokenTransaction, error) {
	var arg1Copy []*token.Token
	if arg1 != nil {
		arg1Copy = make([]*token.Token, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.requestIssueMutex.Lock()
	ret, specificReturn := fake.requestIssueReturnsOnCall[len(fake.requestIssueArgsForCall)]
	fake.requestIssueArgsForCall = append(fake.requestIssueArgsForCall, struct {
		arg1 []*token.Token
	}{arg1Copy})
	fake.recordInvocation("RequestIssue", []interface{}{arg1Copy})
	fake.requestIssueMutex.Unlock()
	if fake.RequestIssueStub != nil {
		return fake.RequestIssueStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.requestIssueReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Issuer) RequestIssueCallCount() int {
	fake.requestIssueMutex.RLock()
	defer fake.requestIssueMutex.RUnlock()
	return len(fake.requestIssueArgsForCall)
}

func (fake *Issuer) RequestIssueCalls(stub func([]*token.Token) (*token.TokenTransaction, error)) {
	fake.requestIssueMutex.Lock()
	defer fake.requestIssueMutex.Unlock()
	fake.RequestIssueStub = stub
}

func (fake *Issuer) RequestIssueArgsForCall(i int) []*token.Token {
	fake.requestIssueMutex.RLock()
	defer fake.requestIssueMutex.RUnlock()
	argsForCall := fake.requestIssueArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Issuer) RequestIssueReturns(result1 *token.TokenTransaction, result2 error) {
	fake.requestIssueMutex.Lock()
	defer fake.requestIssueMutex.Unlock()
	fake.RequestIssueStub = nil
	fake.requestIssueReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Issuer) RequestIssueReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.requestIssueMutex.Lock()
	defer fake.requestIssueMutex.Unlock()
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

func (fake *Issuer) RequestTokenOperation(arg1 *token.TokenOperation) (*token.TokenTransaction, error) {
	fake.requestTokenOperationMutex.Lock()
	ret, specificReturn := fake.requestTokenOperationReturnsOnCall[len(fake.requestTokenOperationArgsForCall)]
	fake.requestTokenOperationArgsForCall = append(fake.requestTokenOperationArgsForCall, struct {
		arg1 *token.TokenOperation
	}{arg1})
	fake.recordInvocation("RequestTokenOperation", []interface{}{arg1})
	fake.requestTokenOperationMutex.Unlock()
	if fake.RequestTokenOperationStub != nil {
		return fake.RequestTokenOperationStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.requestTokenOperationReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Issuer) RequestTokenOperationCallCount() int {
	fake.requestTokenOperationMutex.RLock()
	defer fake.requestTokenOperationMutex.RUnlock()
	return len(fake.requestTokenOperationArgsForCall)
}

func (fake *Issuer) RequestTokenOperationCalls(stub func(*token.TokenOperation) (*token.TokenTransaction, error)) {
	fake.requestTokenOperationMutex.Lock()
	defer fake.requestTokenOperationMutex.Unlock()
	fake.RequestTokenOperationStub = stub
}

func (fake *Issuer) RequestTokenOperationArgsForCall(i int) *token.TokenOperation {
	fake.requestTokenOperationMutex.RLock()
	defer fake.requestTokenOperationMutex.RUnlock()
	argsForCall := fake.requestTokenOperationArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Issuer) RequestTokenOperationReturns(result1 *token.TokenTransaction, result2 error) {
	fake.requestTokenOperationMutex.Lock()
	defer fake.requestTokenOperationMutex.Unlock()
	fake.RequestTokenOperationStub = nil
	fake.requestTokenOperationReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Issuer) RequestTokenOperationReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.requestTokenOperationMutex.Lock()
	defer fake.requestTokenOperationMutex.Unlock()
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
