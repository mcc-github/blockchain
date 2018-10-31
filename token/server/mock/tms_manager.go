
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
	GetTransactorStub        func(channel string, privateCredential, publicCredential []byte) (server.Transactor, error)
	getTransactorMutex       sync.RWMutex
	getTransactorArgsForCall []struct {
		channel           string
		privateCredential []byte
		publicCredential  []byte
	}
	getTransactorReturns struct {
		result1 server.Transactor
		result2 error
	}
	getTransactorReturnsOnCall map[int]struct {
		result1 server.Transactor
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

func (fake *TMSManager) GetTransactor(channel string, privateCredential []byte, publicCredential []byte) (server.Transactor, error) {
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
	fake.getTransactorMutex.Lock()
	ret, specificReturn := fake.getTransactorReturnsOnCall[len(fake.getTransactorArgsForCall)]
	fake.getTransactorArgsForCall = append(fake.getTransactorArgsForCall, struct {
		channel           string
		privateCredential []byte
		publicCredential  []byte
	}{channel, privateCredentialCopy, publicCredentialCopy})
	fake.recordInvocation("GetTransactor", []interface{}{channel, privateCredentialCopy, publicCredentialCopy})
	fake.getTransactorMutex.Unlock()
	if fake.GetTransactorStub != nil {
		return fake.GetTransactorStub(channel, privateCredential, publicCredential)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTransactorReturns.result1, fake.getTransactorReturns.result2
}

func (fake *TMSManager) GetTransactorCallCount() int {
	fake.getTransactorMutex.RLock()
	defer fake.getTransactorMutex.RUnlock()
	return len(fake.getTransactorArgsForCall)
}

func (fake *TMSManager) GetTransactorArgsForCall(i int) (string, []byte, []byte) {
	fake.getTransactorMutex.RLock()
	defer fake.getTransactorMutex.RUnlock()
	return fake.getTransactorArgsForCall[i].channel, fake.getTransactorArgsForCall[i].privateCredential, fake.getTransactorArgsForCall[i].publicCredential
}

func (fake *TMSManager) GetTransactorReturns(result1 server.Transactor, result2 error) {
	fake.GetTransactorStub = nil
	fake.getTransactorReturns = struct {
		result1 server.Transactor
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) GetTransactorReturnsOnCall(i int, result1 server.Transactor, result2 error) {
	fake.GetTransactorStub = nil
	if fake.getTransactorReturnsOnCall == nil {
		fake.getTransactorReturnsOnCall = make(map[int]struct {
			result1 server.Transactor
			result2 error
		})
	}
	fake.getTransactorReturnsOnCall[i] = struct {
		result1 server.Transactor
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getIssuerMutex.RLock()
	defer fake.getIssuerMutex.RUnlock()
	fake.getTransactorMutex.RLock()
	defer fake.getTransactorMutex.RUnlock()
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
