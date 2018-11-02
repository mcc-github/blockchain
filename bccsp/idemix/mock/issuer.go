
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/bccsp/idemix"
)

type Issuer struct {
	NewKeyStub        func(AttributeNames []string) (idemix.IssuerSecretKey, error)
	newKeyMutex       sync.RWMutex
	newKeyArgsForCall []struct {
		AttributeNames []string
	}
	newKeyReturns struct {
		result1 idemix.IssuerSecretKey
		result2 error
	}
	newKeyReturnsOnCall map[int]struct {
		result1 idemix.IssuerSecretKey
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Issuer) NewKey(AttributeNames []string) (idemix.IssuerSecretKey, error) {
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

func (fake *Issuer) NewKeyReturns(result1 idemix.IssuerSecretKey, result2 error) {
	fake.NewKeyStub = nil
	fake.newKeyReturns = struct {
		result1 idemix.IssuerSecretKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) NewKeyReturnsOnCall(i int, result1 idemix.IssuerSecretKey, result2 error) {
	fake.NewKeyStub = nil
	if fake.newKeyReturnsOnCall == nil {
		fake.newKeyReturnsOnCall = make(map[int]struct {
			result1 idemix.IssuerSecretKey
			result2 error
		})
	}
	fake.newKeyReturnsOnCall[i] = struct {
		result1 idemix.IssuerSecretKey
		result2 error
	}{result1, result2}
}

func (fake *Issuer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newKeyMutex.RLock()
	defer fake.newKeyMutex.RUnlock()
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

var _ idemix.Issuer = new(Issuer)
