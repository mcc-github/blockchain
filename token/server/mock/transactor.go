
package mock

import (
	sync "sync"

	token "github.com/mcc-github/blockchain/protos/token"
	server "github.com/mcc-github/blockchain/token/server"
)

type Transactor struct {
	ListTokensStub        func() (*token.UnspentTokens, error)
	listTokensMutex       sync.RWMutex
	listTokensArgsForCall []struct {
	}
	listTokensReturns struct {
		result1 *token.UnspentTokens
		result2 error
	}
	listTokensReturnsOnCall map[int]struct {
		result1 *token.UnspentTokens
		result2 error
	}
	RequestApproveStub        func(*token.ApproveRequest) (*token.TokenTransaction, error)
	requestApproveMutex       sync.RWMutex
	requestApproveArgsForCall []struct {
		arg1 *token.ApproveRequest
	}
	requestApproveReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestApproveReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	RequestRedeemStub        func(*token.RedeemRequest) (*token.TokenTransaction, error)
	requestRedeemMutex       sync.RWMutex
	requestRedeemArgsForCall []struct {
		arg1 *token.RedeemRequest
	}
	requestRedeemReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestRedeemReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	RequestTransferStub        func(*token.TransferRequest) (*token.TokenTransaction, error)
	requestTransferMutex       sync.RWMutex
	requestTransferArgsForCall []struct {
		arg1 *token.TransferRequest
	}
	requestTransferReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestTransferReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	RequestTransferFromStub        func(*token.TransferRequest) (*token.TokenTransaction, error)
	requestTransferFromMutex       sync.RWMutex
	requestTransferFromArgsForCall []struct {
		arg1 *token.TransferRequest
	}
	requestTransferFromReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestTransferFromReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Transactor) ListTokens() (*token.UnspentTokens, error) {
	fake.listTokensMutex.Lock()
	ret, specificReturn := fake.listTokensReturnsOnCall[len(fake.listTokensArgsForCall)]
	fake.listTokensArgsForCall = append(fake.listTokensArgsForCall, struct {
	}{})
	fake.recordInvocation("ListTokens", []interface{}{})
	fake.listTokensMutex.Unlock()
	if fake.ListTokensStub != nil {
		return fake.ListTokensStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listTokensReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Transactor) ListTokensCallCount() int {
	fake.listTokensMutex.RLock()
	defer fake.listTokensMutex.RUnlock()
	return len(fake.listTokensArgsForCall)
}

func (fake *Transactor) ListTokensCalls(stub func() (*token.UnspentTokens, error)) {
	fake.listTokensMutex.Lock()
	defer fake.listTokensMutex.Unlock()
	fake.ListTokensStub = stub
}

func (fake *Transactor) ListTokensReturns(result1 *token.UnspentTokens, result2 error) {
	fake.listTokensMutex.Lock()
	defer fake.listTokensMutex.Unlock()
	fake.ListTokensStub = nil
	fake.listTokensReturns = struct {
		result1 *token.UnspentTokens
		result2 error
	}{result1, result2}
}

func (fake *Transactor) ListTokensReturnsOnCall(i int, result1 *token.UnspentTokens, result2 error) {
	fake.listTokensMutex.Lock()
	defer fake.listTokensMutex.Unlock()
	fake.ListTokensStub = nil
	if fake.listTokensReturnsOnCall == nil {
		fake.listTokensReturnsOnCall = make(map[int]struct {
			result1 *token.UnspentTokens
			result2 error
		})
	}
	fake.listTokensReturnsOnCall[i] = struct {
		result1 *token.UnspentTokens
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestApprove(arg1 *token.ApproveRequest) (*token.TokenTransaction, error) {
	fake.requestApproveMutex.Lock()
	ret, specificReturn := fake.requestApproveReturnsOnCall[len(fake.requestApproveArgsForCall)]
	fake.requestApproveArgsForCall = append(fake.requestApproveArgsForCall, struct {
		arg1 *token.ApproveRequest
	}{arg1})
	fake.recordInvocation("RequestApprove", []interface{}{arg1})
	fake.requestApproveMutex.Unlock()
	if fake.RequestApproveStub != nil {
		return fake.RequestApproveStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.requestApproveReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Transactor) RequestApproveCallCount() int {
	fake.requestApproveMutex.RLock()
	defer fake.requestApproveMutex.RUnlock()
	return len(fake.requestApproveArgsForCall)
}

func (fake *Transactor) RequestApproveCalls(stub func(*token.ApproveRequest) (*token.TokenTransaction, error)) {
	fake.requestApproveMutex.Lock()
	defer fake.requestApproveMutex.Unlock()
	fake.RequestApproveStub = stub
}

func (fake *Transactor) RequestApproveArgsForCall(i int) *token.ApproveRequest {
	fake.requestApproveMutex.RLock()
	defer fake.requestApproveMutex.RUnlock()
	argsForCall := fake.requestApproveArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Transactor) RequestApproveReturns(result1 *token.TokenTransaction, result2 error) {
	fake.requestApproveMutex.Lock()
	defer fake.requestApproveMutex.Unlock()
	fake.RequestApproveStub = nil
	fake.requestApproveReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestApproveReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.requestApproveMutex.Lock()
	defer fake.requestApproveMutex.Unlock()
	fake.RequestApproveStub = nil
	if fake.requestApproveReturnsOnCall == nil {
		fake.requestApproveReturnsOnCall = make(map[int]struct {
			result1 *token.TokenTransaction
			result2 error
		})
	}
	fake.requestApproveReturnsOnCall[i] = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestRedeem(arg1 *token.RedeemRequest) (*token.TokenTransaction, error) {
	fake.requestRedeemMutex.Lock()
	ret, specificReturn := fake.requestRedeemReturnsOnCall[len(fake.requestRedeemArgsForCall)]
	fake.requestRedeemArgsForCall = append(fake.requestRedeemArgsForCall, struct {
		arg1 *token.RedeemRequest
	}{arg1})
	fake.recordInvocation("RequestRedeem", []interface{}{arg1})
	fake.requestRedeemMutex.Unlock()
	if fake.RequestRedeemStub != nil {
		return fake.RequestRedeemStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.requestRedeemReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Transactor) RequestRedeemCallCount() int {
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	return len(fake.requestRedeemArgsForCall)
}

func (fake *Transactor) RequestRedeemCalls(stub func(*token.RedeemRequest) (*token.TokenTransaction, error)) {
	fake.requestRedeemMutex.Lock()
	defer fake.requestRedeemMutex.Unlock()
	fake.RequestRedeemStub = stub
}

func (fake *Transactor) RequestRedeemArgsForCall(i int) *token.RedeemRequest {
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	argsForCall := fake.requestRedeemArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Transactor) RequestRedeemReturns(result1 *token.TokenTransaction, result2 error) {
	fake.requestRedeemMutex.Lock()
	defer fake.requestRedeemMutex.Unlock()
	fake.RequestRedeemStub = nil
	fake.requestRedeemReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestRedeemReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.requestRedeemMutex.Lock()
	defer fake.requestRedeemMutex.Unlock()
	fake.RequestRedeemStub = nil
	if fake.requestRedeemReturnsOnCall == nil {
		fake.requestRedeemReturnsOnCall = make(map[int]struct {
			result1 *token.TokenTransaction
			result2 error
		})
	}
	fake.requestRedeemReturnsOnCall[i] = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestTransfer(arg1 *token.TransferRequest) (*token.TokenTransaction, error) {
	fake.requestTransferMutex.Lock()
	ret, specificReturn := fake.requestTransferReturnsOnCall[len(fake.requestTransferArgsForCall)]
	fake.requestTransferArgsForCall = append(fake.requestTransferArgsForCall, struct {
		arg1 *token.TransferRequest
	}{arg1})
	fake.recordInvocation("RequestTransfer", []interface{}{arg1})
	fake.requestTransferMutex.Unlock()
	if fake.RequestTransferStub != nil {
		return fake.RequestTransferStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.requestTransferReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Transactor) RequestTransferCallCount() int {
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	return len(fake.requestTransferArgsForCall)
}

func (fake *Transactor) RequestTransferCalls(stub func(*token.TransferRequest) (*token.TokenTransaction, error)) {
	fake.requestTransferMutex.Lock()
	defer fake.requestTransferMutex.Unlock()
	fake.RequestTransferStub = stub
}

func (fake *Transactor) RequestTransferArgsForCall(i int) *token.TransferRequest {
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	argsForCall := fake.requestTransferArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Transactor) RequestTransferReturns(result1 *token.TokenTransaction, result2 error) {
	fake.requestTransferMutex.Lock()
	defer fake.requestTransferMutex.Unlock()
	fake.RequestTransferStub = nil
	fake.requestTransferReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestTransferReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.requestTransferMutex.Lock()
	defer fake.requestTransferMutex.Unlock()
	fake.RequestTransferStub = nil
	if fake.requestTransferReturnsOnCall == nil {
		fake.requestTransferReturnsOnCall = make(map[int]struct {
			result1 *token.TokenTransaction
			result2 error
		})
	}
	fake.requestTransferReturnsOnCall[i] = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestTransferFrom(arg1 *token.TransferRequest) (*token.TokenTransaction, error) {
	fake.requestTransferFromMutex.Lock()
	ret, specificReturn := fake.requestTransferFromReturnsOnCall[len(fake.requestTransferFromArgsForCall)]
	fake.requestTransferFromArgsForCall = append(fake.requestTransferFromArgsForCall, struct {
		arg1 *token.TransferRequest
	}{arg1})
	fake.recordInvocation("RequestTransferFrom", []interface{}{arg1})
	fake.requestTransferFromMutex.Unlock()
	if fake.RequestTransferFromStub != nil {
		return fake.RequestTransferFromStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.requestTransferFromReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Transactor) RequestTransferFromCallCount() int {
	fake.requestTransferFromMutex.RLock()
	defer fake.requestTransferFromMutex.RUnlock()
	return len(fake.requestTransferFromArgsForCall)
}

func (fake *Transactor) RequestTransferFromCalls(stub func(*token.TransferRequest) (*token.TokenTransaction, error)) {
	fake.requestTransferFromMutex.Lock()
	defer fake.requestTransferFromMutex.Unlock()
	fake.RequestTransferFromStub = stub
}

func (fake *Transactor) RequestTransferFromArgsForCall(i int) *token.TransferRequest {
	fake.requestTransferFromMutex.RLock()
	defer fake.requestTransferFromMutex.RUnlock()
	argsForCall := fake.requestTransferFromArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Transactor) RequestTransferFromReturns(result1 *token.TokenTransaction, result2 error) {
	fake.requestTransferFromMutex.Lock()
	defer fake.requestTransferFromMutex.Unlock()
	fake.RequestTransferFromStub = nil
	fake.requestTransferFromReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestTransferFromReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
	fake.requestTransferFromMutex.Lock()
	defer fake.requestTransferFromMutex.Unlock()
	fake.RequestTransferFromStub = nil
	if fake.requestTransferFromReturnsOnCall == nil {
		fake.requestTransferFromReturnsOnCall = make(map[int]struct {
			result1 *token.TokenTransaction
			result2 error
		})
	}
	fake.requestTransferFromReturnsOnCall[i] = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listTokensMutex.RLock()
	defer fake.listTokensMutex.RUnlock()
	fake.requestApproveMutex.RLock()
	defer fake.requestApproveMutex.RUnlock()
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	fake.requestTransferFromMutex.RLock()
	defer fake.requestTransferFromMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Transactor) recordInvocation(key string, args []interface{}) {
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

var _ server.Transactor = new(Transactor)
