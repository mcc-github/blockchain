
package mock

import (
	"crypto/ecdsa"
	"sync"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type Revocation struct {
	NewKeyStub        func() (*ecdsa.PrivateKey, error)
	newKeyMutex       sync.RWMutex
	newKeyArgsForCall []struct{}
	newKeyReturns     struct {
		result1 *ecdsa.PrivateKey
		result2 error
	}
	newKeyReturnsOnCall map[int]struct {
		result1 *ecdsa.PrivateKey
		result2 error
	}
	SignStub        func(key *ecdsa.PrivateKey, unrevokedHandles [][]byte, epoch int, alg bccsp.RevocationAlgorithm) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		key              *ecdsa.PrivateKey
		unrevokedHandles [][]byte
		epoch            int
		alg              bccsp.RevocationAlgorithm
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyStub        func(pk *ecdsa.PublicKey, cri []byte, epoch int, alg bccsp.RevocationAlgorithm) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		pk    *ecdsa.PublicKey
		cri   []byte
		epoch int
		alg   bccsp.RevocationAlgorithm
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

func (fake *Revocation) NewKey() (*ecdsa.PrivateKey, error) {
	fake.newKeyMutex.Lock()
	ret, specificReturn := fake.newKeyReturnsOnCall[len(fake.newKeyArgsForCall)]
	fake.newKeyArgsForCall = append(fake.newKeyArgsForCall, struct{}{})
	fake.recordInvocation("NewKey", []interface{}{})
	fake.newKeyMutex.Unlock()
	if fake.NewKeyStub != nil {
		return fake.NewKeyStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newKeyReturns.result1, fake.newKeyReturns.result2
}

func (fake *Revocation) NewKeyCallCount() int {
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	return len(fake.newKeyArgsForCall)
}

func (fake *Revocation) NewKeyReturns(result1 *ecdsa.PrivateKey, result2 error) {
	fake.NewKeyStub = nil
	fake.newKeyReturns = struct {
		result1 *ecdsa.PrivateKey
		result2 error
	}{result1, result2}
}

func (fake *Revocation) NewKeyReturnsOnCall(i int, result1 *ecdsa.PrivateKey, result2 error) {
	fake.NewKeyStub = nil
	if fake.newKeyReturnsOnCall == nil {
		fake.newKeyReturnsOnCall = make(map[int]struct {
			result1 *ecdsa.PrivateKey
			result2 error
		})
	}
	fake.newKeyReturnsOnCall[i] = struct {
		result1 *ecdsa.PrivateKey
		result2 error
	}{result1, result2}
}

func (fake *Revocation) Sign(key *ecdsa.PrivateKey, unrevokedHandles [][]byte, epoch int, alg bccsp.RevocationAlgorithm) ([]byte, error) {
	var unrevokedHandlesCopy [][]byte
	if unrevokedHandles != nil {
		unrevokedHandlesCopy = make([][]byte, len(unrevokedHandles))
		copy(unrevokedHandlesCopy, unrevokedHandles)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		key              *ecdsa.PrivateKey
		unrevokedHandles [][]byte
		epoch            int
		alg              bccsp.RevocationAlgorithm
	}{key, unrevokedHandlesCopy, epoch, alg})
	fake.recordInvocation("Sign", []interface{}{key, unrevokedHandlesCopy, epoch, alg})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(key, unrevokedHandles, epoch, alg)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.signReturns.result1, fake.signReturns.result2
}

func (fake *Revocation) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *Revocation) SignArgsForCall(i int) (*ecdsa.PrivateKey, [][]byte, int, bccsp.RevocationAlgorithm) {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return fake.signArgsForCall[i].key, fake.signArgsForCall[i].unrevokedHandles, fake.signArgsForCall[i].epoch, fake.signArgsForCall[i].alg
}

func (fake *Revocation) SignReturns(result1 []byte, result2 error) {
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Revocation) SignReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *Revocation) Verify(pk *ecdsa.PublicKey, cri []byte, epoch int, alg bccsp.RevocationAlgorithm) error {
	var criCopy []byte
	if cri != nil {
		criCopy = make([]byte, len(cri))
		copy(criCopy, cri)
	}
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		pk    *ecdsa.PublicKey
		cri   []byte
		epoch int
		alg   bccsp.RevocationAlgorithm
	}{pk, criCopy, epoch, alg})
	fake.recordInvocation("Verify", []interface{}{pk, criCopy, epoch, alg})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(pk, cri, epoch, alg)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyReturns.result1
}

func (fake *Revocation) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *Revocation) VerifyArgsForCall(i int) (*ecdsa.PublicKey, []byte, int, bccsp.RevocationAlgorithm) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].pk, fake.verifyArgsForCall[i].cri, fake.verifyArgsForCall[i].epoch, fake.verifyArgsForCall[i].alg
}

func (fake *Revocation) VerifyReturns(result1 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *Revocation) VerifyReturnsOnCall(i int, result1 error) {
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

func (fake *Revocation) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
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

func (fake *Revocation) recordInvocation(key string, args []interface{}) {
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

var _ handlers.Revocation = new(Revocation)