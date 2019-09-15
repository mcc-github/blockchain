
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type CredRequest struct {
	SignStub        func(handlers.Big, handlers.IssuerPublicKey, []byte) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		arg1 handlers.Big
		arg2 handlers.IssuerPublicKey
		arg3 []byte
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyStub        func([]byte, handlers.IssuerPublicKey, []byte) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		arg1 []byte
		arg2 handlers.IssuerPublicKey
		arg3 []byte
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

func (fake *CredRequest) Sign(arg1 handlers.Big, arg2 handlers.IssuerPublicKey, arg3 []byte) ([]byte, error) {
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		arg1 handlers.Big
		arg2 handlers.IssuerPublicKey
		arg3 []byte
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("Sign", []interface{}{arg1, arg2, arg3Copy})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.signReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CredRequest) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *CredRequest) SignCalls(stub func(handlers.Big, handlers.IssuerPublicKey, []byte) ([]byte, error)) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = stub
}

func (fake *CredRequest) SignArgsForCall(i int) (handlers.Big, handlers.IssuerPublicKey, []byte) {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	argsForCall := fake.signArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CredRequest) SignReturns(result1 []byte, result2 error) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *CredRequest) SignReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
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

func (fake *CredRequest) Verify(arg1 []byte, arg2 handlers.IssuerPublicKey, arg3 []byte) error {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		arg1 []byte
		arg2 handlers.IssuerPublicKey
		arg3 []byte
	}{arg1Copy, arg2, arg3Copy})
	fake.recordInvocation("Verify", []interface{}{arg1Copy, arg2, arg3Copy})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.verifyReturns
	return fakeReturns.result1
}

func (fake *CredRequest) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *CredRequest) VerifyCalls(stub func([]byte, handlers.IssuerPublicKey, []byte) error) {
	fake.verifyMutex.Lock()
	defer fake.verifyMutex.Unlock()
	fake.VerifyStub = stub
}

func (fake *CredRequest) VerifyArgsForCall(i int) ([]byte, handlers.IssuerPublicKey, []byte) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	argsForCall := fake.verifyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CredRequest) VerifyReturns(result1 error) {
	fake.verifyMutex.Lock()
	defer fake.verifyMutex.Unlock()
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *CredRequest) VerifyReturnsOnCall(i int, result1 error) {
	fake.verifyMutex.Lock()
	defer fake.verifyMutex.Unlock()
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
