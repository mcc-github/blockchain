
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	tk "github.com/mcc-github/blockchain/token"
	"github.com/mcc-github/blockchain/token/client"
)

type Prover struct {
	RequestIssueStub         func(tokensToIssue []*token.Token, signingIdentity tk.SigningIdentity) ([]byte, error)
	requestImportMutex       sync.RWMutex
	requestImportArgsForCall []struct {
		tokensToIssue   []*token.Token
		signingIdentity tk.SigningIdentity
	}
	requestImportReturns struct {
		result1 []byte
		result2 error
	}
	requestImportReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	RequestTransferStub        func(tokenIDs []*token.TokenId, shares []*token.RecipientShare, signingIdentity tk.SigningIdentity) ([]byte, error)
	requestTransferMutex       sync.RWMutex
	requestTransferArgsForCall []struct {
		tokenIDs        []*token.TokenId
		shares          []*token.RecipientShare
		signingIdentity tk.SigningIdentity
	}
	requestTransferReturns struct {
		result1 []byte
		result2 error
	}
	requestTransferReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	RequestRedeemStub        func(tokenIDs []*token.TokenId, quantity string, signingIdentity tk.SigningIdentity) ([]byte, error)
	requestRedeemMutex       sync.RWMutex
	requestRedeemArgsForCall []struct {
		tokenIDs        []*token.TokenId
		quantity        string
		signingIdentity tk.SigningIdentity
	}
	requestRedeemReturns struct {
		result1 []byte
		result2 error
	}
	requestRedeemReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	ListTokensStub        func(signingIdentity tk.SigningIdentity) ([]*token.UnspentToken, error)
	listTokensMutex       sync.RWMutex
	listTokensArgsForCall []struct {
		signingIdentity tk.SigningIdentity
	}
	listTokensReturns struct {
		result1 []*token.UnspentToken
		result2 error
	}
	listTokensReturnsOnCall map[int]struct {
		result1 []*token.UnspentToken
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Prover) RequestIssue(tokensToIssue []*token.Token, signingIdentity tk.SigningIdentity) ([]byte, error) {
	var tokensToIssueCopy []*token.Token
	if tokensToIssue != nil {
		tokensToIssueCopy = make([]*token.Token, len(tokensToIssue))
		copy(tokensToIssueCopy, tokensToIssue)
	}
	fake.requestImportMutex.Lock()
	ret, specificReturn := fake.requestImportReturnsOnCall[len(fake.requestImportArgsForCall)]
	fake.requestImportArgsForCall = append(fake.requestImportArgsForCall, struct {
		tokensToIssue   []*token.Token
		signingIdentity tk.SigningIdentity
	}{tokensToIssueCopy, signingIdentity})
	fake.recordInvocation("RequestIssue", []interface{}{tokensToIssueCopy, signingIdentity})
	fake.requestImportMutex.Unlock()
	if fake.RequestIssueStub != nil {
		return fake.RequestIssueStub(tokensToIssue, signingIdentity)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.requestImportReturns.result1, fake.requestImportReturns.result2
}

func (fake *Prover) RequestIssueCallCount() int {
	fake.requestImportMutex.RLock()
	defer fake.requestImportMutex.RUnlock()
	return len(fake.requestImportArgsForCall)
}

func (fake *Prover) RequestIssueArgsForCall(i int) ([]*token.Token, tk.SigningIdentity) {
	fake.requestImportMutex.RLock()
	defer fake.requestImportMutex.RUnlock()
	return fake.requestImportArgsForCall[i].tokensToIssue, fake.requestImportArgsForCall[i].signingIdentity
}

func (fake *Prover) RequestIssueReturns(result1 []byte, result2 error) {
	fake.RequestIssueStub = nil
	fake.requestImportReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Prover) RequestIssueReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.RequestIssueStub = nil
	if fake.requestImportReturnsOnCall == nil {
		fake.requestImportReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.requestImportReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Prover) RequestTransfer(tokenIDs []*token.TokenId, shares []*token.RecipientShare, signingIdentity tk.SigningIdentity) ([]byte, error) {
	var tokenIDsCopy []*token.TokenId
	if tokenIDs != nil {
		tokenIDsCopy = make([]*token.TokenId, len(tokenIDs))
		copy(tokenIDsCopy, tokenIDs)
	}
	var sharesCopy []*token.RecipientShare
	if shares != nil {
		sharesCopy = make([]*token.RecipientShare, len(shares))
		copy(sharesCopy, shares)
	}
	fake.requestTransferMutex.Lock()
	ret, specificReturn := fake.requestTransferReturnsOnCall[len(fake.requestTransferArgsForCall)]
	fake.requestTransferArgsForCall = append(fake.requestTransferArgsForCall, struct {
		tokenIDs        []*token.TokenId
		shares          []*token.RecipientShare
		signingIdentity tk.SigningIdentity
	}{tokenIDsCopy, sharesCopy, signingIdentity})
	fake.recordInvocation("RequestTransfer", []interface{}{tokenIDsCopy, sharesCopy, signingIdentity})
	fake.requestTransferMutex.Unlock()
	if fake.RequestTransferStub != nil {
		return fake.RequestTransferStub(tokenIDs, shares, signingIdentity)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.requestTransferReturns.result1, fake.requestTransferReturns.result2
}

func (fake *Prover) RequestTransferCallCount() int {
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	return len(fake.requestTransferArgsForCall)
}

func (fake *Prover) RequestTransferArgsForCall(i int) ([]*token.TokenId, []*token.RecipientShare, tk.SigningIdentity) {
	fake.requestTransferMutex.RLock()
	defer fake.requestTransferMutex.RUnlock()
	return fake.requestTransferArgsForCall[i].tokenIDs, fake.requestTransferArgsForCall[i].shares, fake.requestTransferArgsForCall[i].signingIdentity
}

func (fake *Prover) RequestTransferReturns(result1 []byte, result2 error) {
	fake.RequestTransferStub = nil
	fake.requestTransferReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Prover) RequestTransferReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.RequestTransferStub = nil
	if fake.requestTransferReturnsOnCall == nil {
		fake.requestTransferReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.requestTransferReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Prover) RequestRedeem(tokenIDs []*token.TokenId, quantity string, signingIdentity tk.SigningIdentity) ([]byte, error) {
	var tokenIDsCopy []*token.TokenId
	if tokenIDs != nil {
		tokenIDsCopy = make([]*token.TokenId, len(tokenIDs))
		copy(tokenIDsCopy, tokenIDs)
	}
	fake.requestRedeemMutex.Lock()
	ret, specificReturn := fake.requestRedeemReturnsOnCall[len(fake.requestRedeemArgsForCall)]
	fake.requestRedeemArgsForCall = append(fake.requestRedeemArgsForCall, struct {
		tokenIDs        []*token.TokenId
		quantity        string
		signingIdentity tk.SigningIdentity
	}{tokenIDsCopy, quantity, signingIdentity})
	fake.recordInvocation("RequestRedeem", []interface{}{tokenIDsCopy, quantity, signingIdentity})
	fake.requestRedeemMutex.Unlock()
	if fake.RequestRedeemStub != nil {
		return fake.RequestRedeemStub(tokenIDs, quantity, signingIdentity)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.requestRedeemReturns.result1, fake.requestRedeemReturns.result2
}

func (fake *Prover) RequestRedeemCallCount() int {
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	return len(fake.requestRedeemArgsForCall)
}

func (fake *Prover) RequestRedeemArgsForCall(i int) ([]*token.TokenId, string, tk.SigningIdentity) {
	fake.requestRedeemMutex.RLock()
	defer fake.requestRedeemMutex.RUnlock()
	return fake.requestRedeemArgsForCall[i].tokenIDs, fake.requestRedeemArgsForCall[i].quantity, fake.requestRedeemArgsForCall[i].signingIdentity
}

func (fake *Prover) RequestRedeemReturns(result1 []byte, result2 error) {
	fake.RequestRedeemStub = nil
	fake.requestRedeemReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Prover) RequestRedeemReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.RequestRedeemStub = nil
	if fake.requestRedeemReturnsOnCall == nil {
		fake.requestRedeemReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.requestRedeemReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Prover) ListTokens(signingIdentity tk.SigningIdentity) ([]*token.UnspentToken, error) {
	fake.listTokensMutex.Lock()
	ret, specificReturn := fake.listTokensReturnsOnCall[len(fake.listTokensArgsForCall)]
	fake.listTokensArgsForCall = append(fake.listTokensArgsForCall, struct {
		signingIdentity tk.SigningIdentity
	}{signingIdentity})
	fake.recordInvocation("ListTokens", []interface{}{signingIdentity})
	fake.listTokensMutex.Unlock()
	if fake.ListTokensStub != nil {
		return fake.ListTokensStub(signingIdentity)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.listTokensReturns.result1, fake.listTokensReturns.result2
}

func (fake *Prover) ListTokensCallCount() int {
	fake.listTokensMutex.RLock()
	defer fake.listTokensMutex.RUnlock()
	return len(fake.listTokensArgsForCall)
}

func (fake *Prover) ListTokensArgsForCall(i int) tk.SigningIdentity {
	fake.listTokensMutex.RLock()
	defer fake.listTokensMutex.RUnlock()
	return fake.listTokensArgsForCall[i].signingIdentity
}

func (fake *Prover) ListTokensReturns(result1 []*token.UnspentToken, result2 error) {
	fake.ListTokensStub = nil
	fake.listTokensReturns = struct {
		result1 []*token.UnspentToken
		result2 error
	}{result1, result2}
}

func (fake *Prover) ListTokensReturnsOnCall(i int, result1 []*token.UnspentToken, result2 error) {
	fake.ListTokensStub = nil
	if fake.listTokensReturnsOnCall == nil {
		fake.listTokensReturnsOnCall = make(map[int]struct {
			result1 []*token.UnspentToken
			result2 error
		})
	}
	fake.listTokensReturnsOnCall[i] = struct {
		result1 []*token.UnspentToken
		result2 error
	}{result1, result2}
}

func (fake *Prover) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.requestImportMutex.RLock()
	defer fake.requestImportMutex.RUnlock()
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

func (fake *Prover) recordInvocation(key string, args []interface{}) {
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

var _ client.Prover = new(Prover)
