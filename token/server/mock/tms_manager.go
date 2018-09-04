
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/token/server"
)

type TMSManager struct {
	GetIssuerStub        func(channel string, privateCredential, publicCredential []byte) (server.Issuer, error)
	getIssuerMutex       sync.RWMutex
	getIssuerArgsForCall []struct {
		channel           string
		privateCredential []byte
		publicCredential  []byte
	}
	getIssuerReturns struct {
		result1 server.Issuer
		result2 error
	}
	getIssuerReturnsOnCall map[int]struct {
		result1 server.Issuer
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TMSManager) GetIssuer(channel string, privateCredential []byte, publicCredential []byte) (server.Issuer, error) {
	var privateCredentialCopy []byte
	if privateCredential != nil {
		privateCredentialCopy = make([]byte, len(privateCredential))
		copy(privateCredentialCopy, privateCredential)
	}
	var publicCredentialCopy []byte
	if publicCredential != nil {
		publicCredentialCopy = make([]byte, len(publicCredential))
		copy(publicCredentialCopy, publicCredential)
	}
	fake.getIssuerMutex.Lock()
	ret, specificReturn := fake.getIssuerReturnsOnCall[len(fake.getIssuerArgsForCall)]
	fake.getIssuerArgsForCall = append(fake.getIssuerArgsForCall, struct {
		channel           string
		privateCredential []byte
		publicCredential  []byte
	}{channel, privateCredentialCopy, publicCredentialCopy})
	fake.recordInvocation("GetIssuer", []interface{}{channel, privateCredentialCopy, publicCredentialCopy})
	fake.getIssuerMutex.Unlock()
	if fake.GetIssuerStub != nil {
		return fake.GetIssuerStub(channel, privateCredential, publicCredential)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getIssuerReturns.result1, fake.getIssuerReturns.result2
}

func (fake *TMSManager) GetIssuerCallCount() int {
	fake.getIssuerMutex.RLock()
	defer fake.getIssuerMutex.RUnlock()
	return len(fake.getIssuerArgsForCall)
}

func (fake *TMSManager) GetIssuerArgsForCall(i int) (string, []byte, []byte) {
	fake.getIssuerMutex.RLock()
	defer fake.getIssuerMutex.RUnlock()
	return fake.getIssuerArgsForCall[i].channel, fake.getIssuerArgsForCall[i].privateCredential, fake.getIssuerArgsForCall[i].publicCredential
}

func (fake *TMSManager) GetIssuerReturns(result1 server.Issuer, result2 error) {
	fake.GetIssuerStub = nil
	fake.getIssuerReturns = struct {
		result1 server.Issuer
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) GetIssuerReturnsOnCall(i int, result1 server.Issuer, result2 error) {
	fake.GetIssuerStub = nil
	if fake.getIssuerReturnsOnCall == nil {
		fake.getIssuerReturnsOnCall = make(map[int]struct {
			result1 server.Issuer
			result2 error
		})
	}
	fake.getIssuerReturnsOnCall[i] = struct {
		result1 server.Issuer
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getIssuerMutex.RLock()
	defer fake.getIssuerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TMSManager) recordInvocation(key string, args []interface{}) {
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

var _ server.TMSManager = new(TMSManager)
