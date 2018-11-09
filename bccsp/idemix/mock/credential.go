
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/idemix"
)

type Credential struct {
	SignStub        func(key idemix.IssuerSecretKey, credentialRequest []byte, attributes []bccsp.IdemixAttribute) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		key               idemix.IssuerSecretKey
		credentialRequest []byte
		attributes        []bccsp.IdemixAttribute
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyStub        func(sk idemix.Big, ipk idemix.IssuerPublicKey, credential []byte, attributes []bccsp.IdemixAttribute) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		sk         idemix.Big
		ipk        idemix.IssuerPublicKey
		credential []byte
		attributes []bccsp.IdemixAttribute
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

func (fake *Credential) Sign(key idemix.IssuerSecretKey, credentialRequest []byte, attributes []bccsp.IdemixAttribute) ([]byte, error) {
	var credentialRequestCopy []byte
	if credentialRequest != nil {
		credentialRequestCopy = make([]byte, len(credentialRequest))
		copy(credentialRequestCopy, credentialRequest)
	}
	var attributesCopy []bccsp.IdemixAttribute
	if attributes != nil {
		attributesCopy = make([]bccsp.IdemixAttribute, len(attributes))
		copy(attributesCopy, attributes)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		key               idemix.IssuerSecretKey
		credentialRequest []byte
		attributes        []bccsp.IdemixAttribute
	}{key, credentialRequestCopy, attributesCopy})
	fake.recordInvocation("Sign", []interface{}{key, credentialRequestCopy, attributesCopy})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(key, credentialRequest, attributes)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.signReturns.result1, fake.signReturns.result2
}

func (fake *Credential) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *Credential) SignArgsForCall(i int) (idemix.IssuerSecretKey, []byte, []bccsp.IdemixAttribute) {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return fake.signArgsForCall[i].key, fake.signArgsForCall[i].credentialRequest, fake.signArgsForCall[i].attributes
}

func (fake *Credential) SignReturns(result1 []byte, result2 error) {
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Credential) SignReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *Credential) Verify(sk idemix.Big, ipk idemix.IssuerPublicKey, credential []byte, attributes []bccsp.IdemixAttribute) error {
	var credentialCopy []byte
	if credential != nil {
		credentialCopy = make([]byte, len(credential))
		copy(credentialCopy, credential)
	}
	var attributesCopy []bccsp.IdemixAttribute
	if attributes != nil {
		attributesCopy = make([]bccsp.IdemixAttribute, len(attributes))
		copy(attributesCopy, attributes)
	}
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		sk         idemix.Big
		ipk        idemix.IssuerPublicKey
		credential []byte
		attributes []bccsp.IdemixAttribute
	}{sk, ipk, credentialCopy, attributesCopy})
	fake.recordInvocation("Verify", []interface{}{sk, ipk, credentialCopy, attributesCopy})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(sk, ipk, credential, attributes)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyReturns.result1
}

func (fake *Credential) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *Credential) VerifyArgsForCall(i int) (idemix.Big, idemix.IssuerPublicKey, []byte, []bccsp.IdemixAttribute) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].sk, fake.verifyArgsForCall[i].ipk, fake.verifyArgsForCall[i].credential, fake.verifyArgsForCall[i].attributes
}

func (fake *Credential) VerifyReturns(result1 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *Credential) VerifyReturnsOnCall(i int, result1 error) {
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

var _ idemix.Credential = new(Credential)
