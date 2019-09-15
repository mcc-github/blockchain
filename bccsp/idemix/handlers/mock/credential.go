
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type Credential struct {
	SignStub        func(handlers.IssuerSecretKey, []byte, []bccsp.IdemixAttribute) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		arg1 handlers.IssuerSecretKey
		arg2 []byte
		arg3 []bccsp.IdemixAttribute
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyStub        func(handlers.Big, handlers.IssuerPublicKey, []byte, []bccsp.IdemixAttribute) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		arg1 handlers.Big
		arg2 handlers.IssuerPublicKey
		arg3 []byte
		arg4 []bccsp.IdemixAttribute
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

func (fake *Credential) Sign(arg1 handlers.IssuerSecretKey, arg2 []byte, arg3 []bccsp.IdemixAttribute) ([]byte, error) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	var arg3Copy []bccsp.IdemixAttribute
	if arg3 != nil {
		arg3Copy = make([]bccsp.IdemixAttribute, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		arg1 handlers.IssuerSecretKey
		arg2 []byte
		arg3 []bccsp.IdemixAttribute
	}{arg1, arg2Copy, arg3Copy})
	fake.recordInvocation("Sign", []interface{}{arg1, arg2Copy, arg3Copy})
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

func (fake *Credential) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *Credential) SignCalls(stub func(handlers.IssuerSecretKey, []byte, []bccsp.IdemixAttribute) ([]byte, error)) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = stub
}

func (fake *Credential) SignArgsForCall(i int) (handlers.IssuerSecretKey, []byte, []bccsp.IdemixAttribute) {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	argsForCall := fake.signArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *Credential) SignReturns(result1 []byte, result2 error) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Credential) SignReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *Credential) Verify(arg1 handlers.Big, arg2 handlers.IssuerPublicKey, arg3 []byte, arg4 []bccsp.IdemixAttribute) error {
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	var arg4Copy []bccsp.IdemixAttribute
	if arg4 != nil {
		arg4Copy = make([]bccsp.IdemixAttribute, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		arg1 handlers.Big
		arg2 handlers.IssuerPublicKey
		arg3 []byte
		arg4 []bccsp.IdemixAttribute
	}{arg1, arg2, arg3Copy, arg4Copy})
	fake.recordInvocation("Verify", []interface{}{arg1, arg2, arg3Copy, arg4Copy})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.verifyReturns
	return fakeReturns.result1
}

func (fake *Credential) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *Credential) VerifyCalls(stub func(handlers.Big, handlers.IssuerPublicKey, []byte, []bccsp.IdemixAttribute) error) {
	fake.verifyMutex.Lock()
	defer fake.verifyMutex.Unlock()
	fake.VerifyStub = stub
}

func (fake *Credential) VerifyArgsForCall(i int) (handlers.Big, handlers.IssuerPublicKey, []byte, []bccsp.IdemixAttribute) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	argsForCall := fake.verifyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *Credential) VerifyReturns(result1 error) {
	fake.verifyMutex.Lock()
	defer fake.verifyMutex.Unlock()
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *Credential) VerifyReturnsOnCall(i int, result1 error) {
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

func (fake *Credential) Invocations() map[string][][]interface{} {
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

func (fake *Credential) recordInvocation(key string, args []interface{}) {
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

var _ handlers.Credential = new(Credential)
