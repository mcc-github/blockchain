
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix/handlers"
)

type Issuer struct {
	NewKeyStub        func(AttributeNames []string) (handlers.IssuerSecretKey, error)
	newKeyMutex       sync.RWMutex
	newKeyArgsForCall []struct {
		AttributeNames []string
	}
	newKeyReturns struct {
		result1 handlers.IssuerSecretKey
		result2 error
	}
	newKeyReturnsOnCall map[int]struct {
		result1 handlers.IssuerSecretKey
		result2 error
	}
	NewPublicKeyFromBytesStub        func(raw []byte, attributes []string) (handlers.IssuerPublicKey, error)
	newPublicKeyFromBytesMutex       sync.RWMutex
	newPublicKeyFromBytesArgsForCall []struct {
		raw        []byte
		attributes []string
	}
	newPublicKeyFromBytesReturns struct {
		result1 handlers.IssuerPublicKey
		result2 error
	}
	newPublicKeyFromBytesReturnsOnCall map[int]struct {
		result1 handlers.IssuerPublicKey
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Issuer) NewKey(AttributeNames []string) (handlers.IssuerSecretKey, error) {
	var AttributeNamesCopy []string
	if AttributeNames != nil {
		AttributeNamesCopy = make([]string, len(AttributeNames))
		copy(AttributeNamesCopy, AttributeNames)
	}
	fake.newKeyMutex.Lock()
	ret, specificReturn := fake.newKeyReturnsOnCall[len(fake.newKeyArgsForCall)]
	fake.newKeyArgsForCall = append(fake.newKeyArgsForCall, struct {
		AttributeNames []string
	}{AttributeNamesCopy})
	fake.recordInvocation("NewKey", []interface{}{AttributeNamesCopy})
	fake.newKeyMutex.Unlock()
	if fake.NewKeyStub != nil {
		return fake.NewKeyStub(AttributeNames)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newKeyReturns.result1, fake.newKeyReturns.result2
}

func (fake *Issuer) NewKeyCallCount() int {
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	return len(fake.newKeyArgsForCall)
}

func (fake *Issuer) NewKeyArgsForCall(i int) []string {
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	return fake.newKeyArgsForCall[i].AttributeNames
}

func (fake *Issuer) NewKeyReturns(result1 handlers.IssuerSecretKey, result2 error) {
	fake.NewKeyStub = nil
	fake.newKeyReturns = struct {
		result1 handlers.IssuerSecretKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) NewKeyReturnsOnCall(i int, result1 handlers.IssuerSecretKey, result2 error) {
	fake.NewKeyStub = nil
	if fake.newKeyReturnsOnCall == nil {
		fake.newKeyReturnsOnCall = make(map[int]struct {
			result1 handlers.IssuerSecretKey
			result2 error
		})
	}
	fake.newKeyReturnsOnCall[i] = struct {
		result1 handlers.IssuerSecretKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) NewPublicKeyFromBytes(raw []byte, attributes []string) (handlers.IssuerPublicKey, error) {
	var rawCopy []byte
	if raw != nil {
		rawCopy = make([]byte, len(raw))
		copy(rawCopy, raw)
	}
	var attributesCopy []string
	if attributes != nil {
		attributesCopy = make([]string, len(attributes))
		copy(attributesCopy, attributes)
	}
	fake.newPublicKeyFromBytesMutex.Lock()
	ret, specificReturn := fake.newPublicKeyFromBytesReturnsOnCall[len(fake.newPublicKeyFromBytesArgsForCall)]
	fake.newPublicKeyFromBytesArgsForCall = append(fake.newPublicKeyFromBytesArgsForCall, struct {
		raw        []byte
		attributes []string
	}{rawCopy, attributesCopy})
	fake.recordInvocation("NewPublicKeyFromBytes", []interface{}{rawCopy, attributesCopy})
	fake.newPublicKeyFromBytesMutex.Unlock()
	if fake.NewPublicKeyFromBytesStub != nil {
		return fake.NewPublicKeyFromBytesStub(raw, attributes)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newPublicKeyFromBytesReturns.result1, fake.newPublicKeyFromBytesReturns.result2
}

func (fake *Issuer) NewPublicKeyFromBytesCallCount() int {
	fake.newPublicKeyFromBytesMutex.RLock()
	defer fake.newPublicKeyFromBytesMutex.RUnlock()
	return len(fake.newPublicKeyFromBytesArgsForCall)
}

func (fake *Issuer) NewPublicKeyFromBytesArgsForCall(i int) ([]byte, []string) {
	fake.newPublicKeyFromBytesMutex.RLock()
	defer fake.newPublicKeyFromBytesMutex.RUnlock()
	return fake.newPublicKeyFromBytesArgsForCall[i].raw, fake.newPublicKeyFromBytesArgsForCall[i].attributes
}

func (fake *Issuer) NewPublicKeyFromBytesReturns(result1 handlers.IssuerPublicKey, result2 error) {
	fake.NewPublicKeyFromBytesStub = nil
	fake.newPublicKeyFromBytesReturns = struct {
		result1 handlers.IssuerPublicKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) NewPublicKeyFromBytesReturnsOnCall(i int, result1 handlers.IssuerPublicKey, result2 error) {
	fake.NewPublicKeyFromBytesStub = nil
	if fake.newPublicKeyFromBytesReturnsOnCall == nil {
		fake.newPublicKeyFromBytesReturnsOnCall = make(map[int]struct {
			result1 handlers.IssuerPublicKey
			result2 error
		})
	}
	fake.newPublicKeyFromBytesReturnsOnCall[i] = struct {
		result1 handlers.IssuerPublicKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
	fake.newPublicKeyFromBytesMutex.RLock()
	defer fake.newPublicKeyFromBytesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Issuer) recordInvocation(key string, args []interface{}) {
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

var _ handlers.Issuer = new(Issuer)
