
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type CredRequest struct {
	SignStub        func(sk handlers.Big, ipk handlers.IssuerPublicKey) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		sk  handlers.Big
		ipk handlers.IssuerPublicKey
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyStub        func(credRequest []byte, ipk handlers.IssuerPublicKey) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		credRequest []byte
		ipk         handlers.IssuerPublicKey
	}
	verifyReturns struct {
		result1 error
	}
	verifyReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CredRequest) Sign(sk handlers.Big, ipk handlers.IssuerPublicKey) ([]byte, error) {
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		sk  handlers.Big
		ipk handlers.IssuerPublicKey
	}{sk, ipk})
	fake.recordInvocation("Sign", []interface{}{sk, ipk})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(sk, ipk)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.signReturns.result1, fake.signReturns.result2
}

func (fake *CredRequest) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *CredRequest) SignArgsForCall(i int) (handlers.Big, handlers.IssuerPublicKey) {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return fake.signArgsForCall[i].sk, fake.signArgsForCall[i].ipk
}

func (fake *CredRequest) SignReturns(result1 []byte, result2 error) {
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *CredRequest) SignReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.SignStub = nil
	if fake.signReturnsOnCall == nil {
		fake.signReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.signReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *CredRequest) Verify(credRequest []byte, ipk handlers.IssuerPublicKey) error {
	var credRequestCopy []byte
	if credRequest != nil {
		credRequestCopy = make([]byte, len(credRequest))
		copy(credRequestCopy, credRequest)
	}
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		credRequest []byte
		ipk         handlers.IssuerPublicKey
	}{credRequestCopy, ipk})
	fake.recordInvocation("Verify", []interface{}{credRequestCopy, ipk})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(credRequest, ipk)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyReturns.result1
}

func (fake *CredRequest) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *CredRequest) VerifyArgsForCall(i int) ([]byte, handlers.IssuerPublicKey) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].credRequest, fake.verifyArgsForCall[i].ipk
}

func (fake *CredRequest) VerifyReturns(result1 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *CredRequest) VerifyReturnsOnCall(i int, result1 error) {
	fake.VerifyStub = nil
	if fake.verifyReturnsOnCall == nil {
		fake.verifyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.verifyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CredRequest) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CredRequest) recordInvocation(key string, args []interface{}) {
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

var _ handlers.CredRequest = new(CredRequest)
