
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/server"
)

type Transactor struct {
	RequestTransferStub        func(request *token.TransferRequest) (*token.TokenTransaction, error)
	requestTransferMutex       sync.RWMutex
	requestTransferArgsForCall []struct {
		request *token.TransferRequest
	}
	requestTransferReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestTransferReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	RequestRedeemStub        func(request *token.RedeemRequest) (*token.TokenTransaction, error)
	requestRedeemMutex       sync.RWMutex
	requestRedeemArgsForCall []struct {
		request *token.RedeemRequest
	}
	requestRedeemReturns struct {
		result1 *token.TokenTransaction
		result2 error
	}
	requestRedeemReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 error
	}
	ListTokensStub        func() (*token.UnspentTokens, error)
	listTokensMutex       sync.RWMutex
	listTokensArgsForCall []struct{}
	listTokensReturns     struct {
		result1 *token.UnspentTokens
		result2 error
	}
	listTokensReturnsOnCall map[int]struct {
		result1 *token.UnspentTokens
		result2 error
	}
	RequestTokenOperationStub        func(tokenIDs []*token.TokenId, op *token.TokenOperation) (*token.TokenTransaction, int, error)
	requestTokenOperationMutex       sync.RWMutex
	requestTokenOperationArgsForCall []struct {
		tokenIDs []*token.TokenId
		op       *token.TokenOperation
	}
	requestTokenOperationReturns struct {
		result1 *token.TokenTransaction
		result2 int
		result3 error
	}
	requestTokenOperationReturnsOnCall map[int]struct {
		result1 *token.TokenTransaction
		result2 int
		result3 error
	}
	DoneStub         func()
	doneMutex        sync.RWMutex
	doneArgsForCall  []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Transactor) RequestTransfer(request *token.TransferRequest) (*token.TokenTransaction, error) {
	fake.requestTransferMutex.Lock()
	ret, specificReturn := fake.requestTransferReturnsOnCall[len(fake.requestTransferArgsForCall)]
	fake.requestTransferArgsForCall = append(fake.requestTransferArgsForCall, struct {
		request *token.TransferRequest
	}{request})
	fake.recordInvocation("RequestTransfer", []interface{}{request})
	fake.requestTransferMutex.Unlock()
	if fake.RequestTransferStub != nil {
		return fake.RequestTransferStub(request)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.requestTransferReturns.result1, fake.requestTransferReturns.result2
}

func (fake *Transactor) RequestTransferCallCount() int {
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	return len(fake.requestTransferArgsForCall)
}

func (fake *Transactor) RequestTransferArgsForCall(i int) *token.TransferRequest {
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	return fake.requestTransferArgsForCall[i].request
}

func (fake *Transactor) RequestTransferReturns(result1 *token.TokenTransaction, result2 error) {
	fake.RequestTransferStub = nil
	fake.requestTransferReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestTransferReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
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

func (fake *Transactor) RequestRedeem(request *token.RedeemRequest) (*token.TokenTransaction, error) {
	fake.requestRedeemMutex.Lock()
	ret, specificReturn := fake.requestRedeemReturnsOnCall[len(fake.requestRedeemArgsForCall)]
	fake.requestRedeemArgsForCall = append(fake.requestRedeemArgsForCall, struct {
		request *token.RedeemRequest
	}{request})
	fake.recordInvocation("RequestRedeem", []interface{}{request})
	fake.requestRedeemMutex.Unlock()
	if fake.RequestRedeemStub != nil {
		return fake.RequestRedeemStub(request)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.requestRedeemReturns.result1, fake.requestRedeemReturns.result2
}

func (fake *Transactor) RequestRedeemCallCount() int {
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	return len(fake.requestRedeemArgsForCall)
}

func (fake *Transactor) RequestRedeemArgsForCall(i int) *token.RedeemRequest {
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	return fake.requestRedeemArgsForCall[i].request
}

func (fake *Transactor) RequestRedeemReturns(result1 *token.TokenTransaction, result2 error) {
	fake.RequestRedeemStub = nil
	fake.requestRedeemReturns = struct {
		result1 *token.TokenTransaction
		result2 error
	}{result1, result2}
}

func (fake *Transactor) RequestRedeemReturnsOnCall(i int, result1 *token.TokenTransaction, result2 error) {
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

func (fake *Transactor) ListTokens() (*token.UnspentTokens, error) {
	fake.listTokensMutex.Lock()
	ret, specificReturn := fake.listTokensReturnsOnCall[len(fake.listTokensArgsForCall)]
	fake.listTokensArgsForCall = append(fake.listTokensArgsForCall, struct{}{})
	fake.recordInvocation("ListTokens", []interface{}{})
	fake.listTokensMutex.Unlock()
	if fake.ListTokensStub != nil {
		return fake.ListTokensStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.listTokensReturns.result1, fake.listTokensReturns.result2
}

func (fake *Transactor) ListTokensCallCount() int {
	fake.listTokensMutex.RLock()
	defer fake.listTokensMutex.RUnlock()
	return len(fake.listTokensArgsForCall)
}

func (fake *Transactor) ListTokensReturns(result1 *token.UnspentTokens, result2 error) {
	fake.ListTokensStub = nil
	fake.listTokensReturns = struct {
		result1 *token.UnspentTokens
		result2 error
	}{result1, result2}
}

func (fake *Transactor) ListTokensReturnsOnCall(i int, result1 *token.UnspentTokens, result2 error) {
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

func (fake *Transactor) RequestTokenOperation(tokenIDs []*token.TokenId, op *token.TokenOperation) (*token.TokenTransaction, int, error) {
	var tokenIDsCopy []*token.TokenId
	if tokenIDs != nil {
		tokenIDsCopy = make([]*token.TokenId, len(tokenIDs))
		copy(tokenIDsCopy, tokenIDs)
	}
	fake.requestTokenOperationMutex.Lock()
	ret, specificReturn := fake.requestTokenOperationReturnsOnCall[len(fake.requestTokenOperationArgsForCall)]
	fake.requestTokenOperationArgsForCall = append(fake.requestTokenOperationArgsForCall, struct {
		tokenIDs []*token.TokenId
		op       *token.TokenOperation
	}{tokenIDsCopy, op})
	fake.recordInvocation("RequestTokenOperation", []interface{}{tokenIDsCopy, op})
	fake.requestTokenOperationMutex.Unlock()
	if fake.RequestTokenOperationStub != nil {
		return fake.RequestTokenOperationStub(tokenIDs, op)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.requestTokenOperationReturns.result1, fake.requestTokenOperationReturns.result2, fake.requestTokenOperationReturns.result3
}

func (fake *Transactor) RequestTokenOperationCallCount() int {
	fake.requestTokenOperationMutex.RLock()
	defer fake.requestTokenOperationMutex.RUnlock()
	return len(fake.requestTokenOperationArgsForCall)
}

func (fake *Transactor) RequestTokenOperationArgsForCall(i int) ([]*token.TokenId, *token.TokenOperation) {
	fake.requestTokenOperationMutex.RLock()
	defer fake.requestTokenOperationMutex.RUnlock()
	return fake.requestTokenOperationArgsForCall[i].tokenIDs, fake.requestTokenOperationArgsForCall[i].op
}

func (fake *Transactor) RequestTokenOperationReturns(result1 *token.TokenTransaction, result2 int, result3 error) {
	fake.RequestTokenOperationStub = nil
	fake.requestTokenOperationReturns = struct {
		result1 *token.TokenTransaction
		result2 int
		result3 error
	}{result1, result2, result3}
}

func (fake *Transactor) RequestTokenOperationReturnsOnCall(i int, result1 *token.TokenTransaction, result2 int, result3 error) {
	fake.RequestTokenOperationStub = nil
	if fake.requestTokenOperationReturnsOnCall == nil {
		fake.requestTokenOperationReturnsOnCall = make(map[int]struct {
			result1 *token.TokenTransaction
			result2 int
			result3 error
		})
	}
	fake.requestTokenOperationReturnsOnCall[i] = struct {
		result1 *token.TokenTransaction
		result2 int
		result3 error
	}{result1, result2, result3}
}

func (fake *Transactor) Done() {
	fake.doneMutex.Lock()
	fake.doneArgsForCall = append(fake.doneArgsForCall, struct{}{})
	fake.recordInvocation("Done", []interface{}{})
	fake.doneMutex.Unlock()
	if fake.DoneStub != nil {
		fake.DoneStub()
	}
}

func (fake *Transactor) DoneCallCount() int {
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
	return len(fake.doneArgsForCall)
}

func (fake *Transactor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	fake.listTokensMutex.RLock()
	defer fake.listTokensMutex.RUnlock()
	fake.requestTokenOperationMutex.RLock()
	defer fake.requestTokenOperationMutex.RUnlock()
	fake.doneMutex.RLock()
	defer fake.doneMutex.RUnlock()
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
