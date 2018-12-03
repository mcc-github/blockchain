
package mock

import (
	"crypto/ecdsa"
	"sync"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type SignatureScheme struct {
	SignStub        func(cred []byte, sk handlers.Big, Nym handlers.Ecp, RNym handlers.Big, ipk handlers.IssuerPublicKey, attributes []bccsp.IdemixAttribute, msg []byte, rhIndex int, cri []byte) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		cred       []byte
		sk         handlers.Big
		Nym        handlers.Ecp
		RNym       handlers.Big
		ipk        handlers.IssuerPublicKey
		attributes []bccsp.IdemixAttribute
		msg        []byte
		rhIndex    int
		cri        []byte
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyStub        func(pk handlers.IssuerPublicKey, signature, digest []byte, attributes []bccsp.IdemixAttribute, hIndex int, revocationPublicKey *ecdsa.PublicKey, epoch int) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		pk                  handlers.IssuerPublicKey
		signature           []byte
		digest              []byte
		attributes          []bccsp.IdemixAttribute
		hIndex              int
		revocationPublicKey *ecdsa.PublicKey
		epoch               int
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

func (fake *SignatureScheme) Sign(cred []byte, sk handlers.Big, Nym handlers.Ecp, RNym handlers.Big, ipk handlers.IssuerPublicKey, attributes []bccsp.IdemixAttribute, msg []byte, rhIndex int, cri []byte) ([]byte, error) {
	var credCopy []byte
	if cred != nil {
		credCopy = make([]byte, len(cred))
		copy(credCopy, cred)
	}
	var attributesCopy []bccsp.IdemixAttribute
	if attributes != nil {
		attributesCopy = make([]bccsp.IdemixAttribute, len(attributes))
		copy(attributesCopy, attributes)
	}
	var msgCopy []byte
	if msg != nil {
		msgCopy = make([]byte, len(msg))
		copy(msgCopy, msg)
	}
	var criCopy []byte
	if cri != nil {
		criCopy = make([]byte, len(cri))
		copy(criCopy, cri)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		cred       []byte
		sk         handlers.Big
		Nym        handlers.Ecp
		RNym       handlers.Big
		ipk        handlers.IssuerPublicKey
		attributes []bccsp.IdemixAttribute
		msg        []byte
		rhIndex    int
		cri        []byte
	}{credCopy, sk, Nym, RNym, ipk, attributesCopy, msgCopy, rhIndex, criCopy})
	fake.recordInvocation("Sign", []interface{}{credCopy, sk, Nym, RNym, ipk, attributesCopy, msgCopy, rhIndex, criCopy})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(cred, sk, Nym, RNym, ipk, attributes, msg, rhIndex, cri)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.signReturns.result1, fake.signReturns.result2
}

func (fake *SignatureScheme) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *SignatureScheme) SignArgsForCall(i int) ([]byte, handlers.Big, handlers.Ecp, handlers.Big, handlers.IssuerPublicKey, []bccsp.IdemixAttribute, []byte, int, []byte) {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return fake.signArgsForCall[i].cred, fake.signArgsForCall[i].sk, fake.signArgsForCall[i].Nym, fake.signArgsForCall[i].RNym, fake.signArgsForCall[i].ipk, fake.signArgsForCall[i].attributes, fake.signArgsForCall[i].msg, fake.signArgsForCall[i].rhIndex, fake.signArgsForCall[i].cri
}

func (fake *SignatureScheme) SignReturns(result1 []byte, result2 error) {
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *SignatureScheme) SignReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *SignatureScheme) Verify(pk handlers.IssuerPublicKey, signature []byte, digest []byte, attributes []bccsp.IdemixAttribute, hIndex int, revocationPublicKey *ecdsa.PublicKey, epoch int) error {
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
	var attributesCopy []bccsp.IdemixAttribute
	if attributes != nil {
		attributesCopy = make([]bccsp.IdemixAttribute, len(attributes))
		copy(attributesCopy, attributes)
	}
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		pk                  handlers.IssuerPublicKey
		signature           []byte
		digest              []byte
		attributes          []bccsp.IdemixAttribute
		hIndex              int
		revocationPublicKey *ecdsa.PublicKey
		epoch               int
	}{pk, signatureCopy, digestCopy, attributesCopy, hIndex, revocationPublicKey, epoch})
	fake.recordInvocation("Verify", []interface{}{pk, signatureCopy, digestCopy, attributesCopy, hIndex, revocationPublicKey, epoch})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(pk, signature, digest, attributes, hIndex, revocationPublicKey, epoch)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyReturns.result1
}

func (fake *SignatureScheme) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *SignatureScheme) VerifyArgsForCall(i int) (handlers.IssuerPublicKey, []byte, []byte, []bccsp.IdemixAttribute, int, *ecdsa.PublicKey, int) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].pk, fake.verifyArgsForCall[i].signature, fake.verifyArgsForCall[i].digest, fake.verifyArgsForCall[i].attributes, fake.verifyArgsForCall[i].hIndex, fake.verifyArgsForCall[i].revocationPublicKey, fake.verifyArgsForCall[i].epoch
}

func (fake *SignatureScheme) VerifyReturns(result1 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *SignatureScheme) VerifyReturnsOnCall(i int, result1 error) {
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

func (fake *SignatureScheme) Invocations() map[string][][]interface{} {
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

func (fake *SignatureScheme) recordInvocation(key string, args []interface{}) {
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

var _ handlers.SignatureScheme = new(SignatureScheme)