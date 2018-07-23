
package mock

import (
	"archive/tar"
	"sync"
)

type PackageWriter struct {
	WriteStub        func(name string, payload []byte, tw *tar.Writer) error
	writeMutex       sync.RWMutex
	writeArgsForCall []struct {
		name    string
		payload []byte
		tw      *tar.Writer
	}
	writeReturns struct {
		result1 error
	}
	writeReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PackageWriter) Write(name string, payload []byte, tw *tar.Writer) error {
	var payloadCopy []byte
	if payload != nil {
		payloadCopy = make([]byte, len(payload))
		copy(payloadCopy, payload)
	}
	fake.writeMutex.Lock()
	ret, specificReturn := fake.writeReturnsOnCall[len(fake.writeArgsForCall)]
	fake.writeArgsForCall = append(fake.writeArgsForCall, struct {
		name    string
		payload []byte
		tw      *tar.Writer
	}{name, payloadCopy, tw})
	fake.recordInvocation("Write", []interface{}{name, payloadCopy, tw})
	fake.writeMutex.Unlock()
	if fake.WriteStub != nil {
		return fake.WriteStub(name, payload, tw)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.writeReturns.result1
}

func (fake *PackageWriter) WriteCallCount() int {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return len(fake.writeArgsForCall)
}

func (fake *PackageWriter) WriteArgsForCall(i int) (string, []byte, *tar.Writer) {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return fake.writeArgsForCall[i].name, fake.writeArgsForCall[i].payload, fake.writeArgsForCall[i].tw
}

func (fake *PackageWriter) WriteReturns(result1 error) {
	fake.WriteStub = nil
	fake.writeReturns = struct {
		result1 error
	}{result1}
}

func (fake *PackageWriter) WriteReturnsOnCall(i int, result1 error) {
	fake.WriteStub = nil
	if fake.writeReturnsOnCall == nil {
		fake.writeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.writeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PackageWriter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PackageWriter) recordInvocation(key string, args []interface{}) {
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