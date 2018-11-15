
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

func (fake *Transactor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	fake.listTokensMutex.RLock()
	defer fake.listTokensMutex.RUnlock()
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
