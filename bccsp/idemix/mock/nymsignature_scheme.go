
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix"
)

type NymSignatureScheme struct {
	SignStub        func(sk idemix.Big, Nym idemix.Ecp, RNym idemix.Big, ipk idemix.IssuerPublicKey, digest []byte) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		sk     idemix.Big
		Nym    idemix.Ecp
		RNym   idemix.Big
		ipk    idemix.IssuerPublicKey
		digest []byte
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyStub        func(pk idemix.IssuerPublicKey, Nym idemix.Ecp, signature, digest []byte) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		pk        idemix.IssuerPublicKey
		Nym       idemix.Ecp
		signature []byte
		digest    []byte
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

func (fake *NymSignatureScheme) Sign(sk idemix.Big, Nym idemix.Ecp, RNym idemix.Big, ipk idemix.IssuerPublicKey, digest []byte) ([]byte, error) {
	var digestCopy []byte
	if digest != nil {
		digestCopy = make([]byte, len(digest))
		copy(digestCopy, digest)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		sk     idemix.Big
		Nym    idemix.Ecp
		RNym   idemix.Big
		ipk    idemix.IssuerPublicKey
		digest []byte
	}{sk, Nym, RNym, ipk, digestCopy})
	fake.recordInvocation("Sign", []interface{}{sk, Nym, RNym, ipk, digestCopy})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(sk, Nym, RNym, ipk, digest)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.signReturns.result1, fake.signReturns.result2
}

func (fake *NymSignatureScheme) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *NymSignatureScheme) SignArgsForCall(i int) (idemix.Big, idemix.Ecp, idemix.Big, idemix.IssuerPublicKey, []byte) {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return fake.signArgsForCall[i].sk, fake.signArgsForCall[i].Nym, fake.signArgsForCall[i].RNym, fake.signArgsForCall[i].ipk, fake.signArgsForCall[i].digest
}

func (fake *NymSignatureScheme) SignReturns(result1 []byte, result2 error) {
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *NymSignatureScheme) SignReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *NymSignatureScheme) Verify(pk idemix.IssuerPublicKey, Nym idemix.Ecp, signature []byte, digest []byte) error {
	var signatureCopy []byte
	if signature != nil {
		signatureCopy = make([]byte, len(signature))
		copy(signatureCopy, signature)
	}
	var digestCopy []byte
	if digest != nil {
		digestCopy = make([]byte, len(digest))
		copy(digestCopy, digest)
	}
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		pk        idemix.IssuerPublicKey
		Nym       idemix.Ecp
		signature []byte
		digest    []byte
	}{pk, Nym, signatureCopy, digestCopy})
	fake.recordInvocation("Verify", []interface{}{pk, Nym, signatureCopy, digestCopy})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(pk, Nym, signature, digest)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyReturns.result1
}

func (fake *NymSignatureScheme) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *NymSignatureScheme) VerifyArgsForCall(i int) (idemix.IssuerPublicKey, idemix.Ecp, []byte, []byte) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].pk, fake.verifyArgsForCall[i].Nym, fake.verifyArgsForCall[i].signature, fake.verifyArgsForCall[i].digest
}

func (fake *NymSignatureScheme) VerifyReturns(result1 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *NymSignatureScheme) VerifyReturnsOnCall(i int, result1 error) {
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

func (fake *NymSignatureScheme) Invocations() map[string][][]interface{} {
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

func (fake *NymSignatureScheme) recordInvocation(key string, args []interface{}) {
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

var _ idemix.NymSignatureScheme = new(NymSignatureScheme)
