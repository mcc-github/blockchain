
package mock

import (
	"context"
	"sync"

	commona "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"google.golang.org/grpc/metadata"
)

type DeliverService struct {
	CloseSendStub        func() error
	closeSendMutex       sync.RWMutex
	closeSendArgsForCall []struct {
	}
	closeSendReturns struct {
		result1 error
	}
	closeSendReturnsOnCall map[int]struct {
		result1 error
	}
	ContextStub        func() context.Context
	contextMutex       sync.RWMutex
	contextArgsForCall []struct {
	}
	contextReturns struct {
		result1 context.Context
	}
	contextReturnsOnCall map[int]struct {
		result1 context.Context
	}
	HeaderStub        func() (metadata.MD, error)
	headerMutex       sync.RWMutex
	headerArgsForCall []struct {
	}
	headerReturns struct {
		result1 metadata.MD
		result2 error
	}
	headerReturnsOnCall map[int]struct {
		result1 metadata.MD
		result2 error
	}
	RecvStub        func() (*orderer.DeliverResponse, error)
	recvMutex       sync.RWMutex
	recvArgsForCall []struct {
	}
	recvReturns struct {
		result1 *orderer.DeliverResponse
		result2 error
	}
	recvReturnsOnCall map[int]struct {
		result1 *orderer.DeliverResponse
		result2 error
	}
	RecvMsgStub        func(interface{}) error
	recvMsgMutex       sync.RWMutex
	recvMsgArgsForCall []struct {
		arg1 interface{}
	}
	recvMsgReturns struct {
		result1 error
	}
	recvMsgReturnsOnCall map[int]struct {
		result1 error
	}
	SendStub        func(*commona.Envelope) error
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		arg1 *commona.Envelope
	}
	sendReturns struct {
		result1 error
	}
	sendReturnsOnCall map[int]struct {
		result1 error
	}
	SendMsgStub        func(interface{}) error
	sendMsgMutex       sync.RWMutex
	sendMsgArgsForCall []struct {
		arg1 interface{}
	}
	sendMsgReturns struct {
		result1 error
	}
	sendMsgReturnsOnCall map[int]struct {
		result1 error
	}
	TrailerStub        func() metadata.MD
	trailerMutex       sync.RWMutex
	trailerArgsForCall []struct {
	}
	trailerReturns struct {
		result1 metadata.MD
	}
	trailerReturnsOnCall map[int]struct {
		result1 metadata.MD
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DeliverService) CloseSend() error {
	fake.closeSendMutex.Lock()
	ret, specificReturn := fake.closeSendReturnsOnCall[len(fake.closeSendArgsForCall)]
	fake.closeSendArgsForCall = append(fake.closeSendArgsForCall, struct {
	}{})
	fake.recordInvocation("CloseSend", []interface{}{})
	fake.closeSendMutex.Unlock()
	if fake.CloseSendStub != nil {
		return fake.CloseSendStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.closeSendReturns
	return fakeReturns.result1
}

func (fake *DeliverService) CloseSendCallCount() int {
	fake.closeSendMutex.RLock()
	defer fake.closeSendMutex.RUnlock()
	return len(fake.closeSendArgsForCall)
}

func (fake *DeliverService) CloseSendCalls(stub func() error) {
	fake.closeSendMutex.Lock()
	defer fake.closeSendMutex.Unlock()
	fake.CloseSendStub = stub
}

func (fake *DeliverService) CloseSendReturns(result1 error) {
	fake.closeSendMutex.Lock()
	defer fake.closeSendMutex.Unlock()
	fake.CloseSendStub = nil
	fake.closeSendReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverService) CloseSendReturnsOnCall(i int, result1 error) {
	fake.closeSendMutex.Lock()
	defer fake.closeSendMutex.Unlock()
	fake.CloseSendStub = nil
	if fake.closeSendReturnsOnCall == nil {
		fake.closeSendReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeSendReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DeliverService) Context() context.Context {
	fake.contextMutex.Lock()
	ret, specificReturn := fake.contextReturnsOnCall[len(fake.contextArgsForCall)]
	fake.contextArgsForCall = append(fake.contextArgsForCall, struct {
	}{})
	fake.recordInvocation("Context", []interface{}{})
	fake.contextMutex.Unlock()
	if fake.ContextStub != nil {
		return fake.ContextStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.contextReturns
	return fakeReturns.result1
}

func (fake *DeliverService) ContextCallCount() int {
	fake.contextMutex.RLock()
	defer fake.contextMutex.RUnlock()
	return len(fake.contextArgsForCall)
}

func (fake *DeliverService) ContextCalls(stub func() context.Context) {
	fake.contextMutex.Lock()
	defer fake.contextMutex.Unlock()
	fake.ContextStub = stub
}

func (fake *DeliverService) ContextReturns(result1 context.Context) {
	fake.contextMutex.Lock()
	defer fake.contextMutex.Unlock()
	fake.ContextStub = nil
	fake.contextReturns = struct {
		result1 context.Context
	}{result1}
}

func (fake *DeliverService) ContextReturnsOnCall(i int, result1 context.Context) {
	fake.contextMutex.Lock()
	defer fake.contextMutex.Unlock()
	fake.ContextStub = nil
	if fake.contextReturnsOnCall == nil {
		fake.contextReturnsOnCall = make(map[int]struct {
			result1 context.Context
		})
	}
	fake.contextReturnsOnCall[i] = struct {
		result1 context.Context
	}{result1}
}

func (fake *DeliverService) Header() (metadata.MD, error) {
	fake.headerMutex.Lock()
	ret, specificReturn := fake.headerReturnsOnCall[len(fake.headerArgsForCall)]
	fake.headerArgsForCall = append(fake.headerArgsForCall, struct {
	}{})
	fake.recordInvocation("Header", []interface{}{})
	fake.headerMutex.Unlock()
	if fake.HeaderStub != nil {
		return fake.HeaderStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.headerReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DeliverService) HeaderCallCount() int {
	fake.headerMutex.RLock()
	defer fake.headerMutex.RUnlock()
	return len(fake.headerArgsForCall)
}

func (fake *DeliverService) HeaderCalls(stub func() (metadata.MD, error)) {
	fake.headerMutex.Lock()
	defer fake.headerMutex.Unlock()
	fake.HeaderStub = stub
}

func (fake *DeliverService) HeaderReturns(result1 metadata.MD, result2 error) {
	fake.headerMutex.Lock()
	defer fake.headerMutex.Unlock()
	fake.HeaderStub = nil
	fake.headerReturns = struct {
		result1 metadata.MD
		result2 error
	}{result1, result2}
}

func (fake *DeliverService) HeaderReturnsOnCall(i int, result1 metadata.MD, result2 error) {
	fake.headerMutex.Lock()
	defer fake.headerMutex.Unlock()
	fake.HeaderStub = nil
	if fake.headerReturnsOnCall == nil {
		fake.headerReturnsOnCall = make(map[int]struct {
			result1 metadata.MD
			result2 error
		})
	}
	fake.headerReturnsOnCall[i] = struct {
		result1 metadata.MD
		result2 error
	}{result1, result2}
}

func (fake *DeliverService) Recv() (*orderer.DeliverResponse, error) {
	fake.recvMutex.Lock()
	ret, specificReturn := fake.recvReturnsOnCall[len(fake.recvArgsForCall)]
	fake.recvArgsForCall = append(fake.recvArgsForCall, struct {
	}{})
	fake.recordInvocation("Recv", []interface{}{})
	fake.recvMutex.Unlock()
	if fake.RecvStub != nil {
		return fake.RecvStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.recvReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DeliverService) RecvCallCount() int {
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	return len(fake.recvArgsForCall)
}

func (fake *DeliverService) RecvCalls(stub func() (*orderer.DeliverResponse, error)) {
	fake.recvMutex.Lock()
	defer fake.recvMutex.Unlock()
	fake.RecvStub = stub
}

func (fake *DeliverService) RecvReturns(result1 *orderer.DeliverResponse, result2 error) {
	fake.recvMutex.Lock()
	defer fake.recvMutex.Unlock()
	fake.RecvStub = nil
	fake.recvReturns = struct {
		result1 *orderer.DeliverResponse
		result2 error
	}{result1, result2}
}

func (fake *DeliverService) RecvReturnsOnCall(i int, result1 *orderer.DeliverResponse, result2 error) {
	fake.recvMutex.Lock()
	defer fake.recvMutex.Unlock()
	fake.RecvStub = nil
	if fake.recvReturnsOnCall == nil {
		fake.recvReturnsOnCall = make(map[int]struct {
			result1 *orderer.DeliverResponse
			result2 error
		})
	}
	fake.recvReturnsOnCall[i] = struct {
		result1 *orderer.DeliverResponse
		result2 error
	}{result1, result2}
}

func (fake *DeliverService) RecvMsg(arg1 interface{}) error {
	fake.recvMsgMutex.Lock()
	ret, specificReturn := fake.recvMsgReturnsOnCall[len(fake.recvMsgArgsForCall)]
	fake.recvMsgArgsForCall = append(fake.recvMsgArgsForCall, struct {
		arg1 interface{}
	}{arg1})
	fake.recordInvocation("RecvMsg", []interface{}{arg1})
	fake.recvMsgMutex.Unlock()
	if fake.RecvMsgStub != nil {
		return fake.RecvMsgStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.recvMsgReturns
	return fakeReturns.result1
}

func (fake *DeliverService) RecvMsgCallCount() int {
	fake.recvMsgMutex.RLock()
	defer fake.recvMsgMutex.RUnlock()
	return len(fake.recvMsgArgsForCall)
}

func (fake *DeliverService) RecvMsgCalls(stub func(interface{}) error) {
	fake.recvMsgMutex.Lock()
	defer fake.recvMsgMutex.Unlock()
	fake.RecvMsgStub = stub
}

func (fake *DeliverService) RecvMsgArgsForCall(i int) interface{} {
	fake.recvMsgMutex.RLock()
	defer fake.recvMsgMutex.RUnlock()
	argsForCall := fake.recvMsgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeliverService) RecvMsgReturns(result1 error) {
	fake.recvMsgMutex.Lock()
	defer fake.recvMsgMutex.Unlock()
	fake.RecvMsgStub = nil
	fake.recvMsgReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverService) RecvMsgReturnsOnCall(i int, result1 error) {
	fake.recvMsgMutex.Lock()
	defer fake.recvMsgMutex.Unlock()
	fake.RecvMsgStub = nil
	if fake.recvMsgReturnsOnCall == nil {
		fake.recvMsgReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.recvMsgReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DeliverService) Send(arg1 *commona.Envelope) error {
	fake.sendMutex.Lock()
	ret, specificReturn := fake.sendReturnsOnCall[len(fake.sendArgsForCall)]
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		arg1 *commona.Envelope
	}{arg1})
	fake.recordInvocation("Send", []interface{}{arg1})
	fake.sendMutex.Unlock()
	if fake.SendStub != nil {
		return fake.SendStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sendReturns
	return fakeReturns.result1
}

func (fake *DeliverService) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *DeliverService) SendCalls(stub func(*commona.Envelope) error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = stub
}

func (fake *DeliverService) SendArgsForCall(i int) *commona.Envelope {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	argsForCall := fake.sendArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeliverService) SendReturns(result1 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverService) SendReturnsOnCall(i int, result1 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	if fake.sendReturnsOnCall == nil {
		fake.sendReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DeliverService) SendMsg(arg1 interface{}) error {
	fake.sendMsgMutex.Lock()
	ret, specificReturn := fake.sendMsgReturnsOnCall[len(fake.sendMsgArgsForCall)]
	fake.sendMsgArgsForCall = append(fake.sendMsgArgsForCall, struct {
		arg1 interface{}
	}{arg1})
	fake.recordInvocation("SendMsg", []interface{}{arg1})
	fake.sendMsgMutex.Unlock()
	if fake.SendMsgStub != nil {
		return fake.SendMsgStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sendMsgReturns
	return fakeReturns.result1
}

func (fake *DeliverService) SendMsgCallCount() int {
	fake.sendMsgMutex.RLock()
	defer fake.sendMsgMutex.RUnlock()
	return len(fake.sendMsgArgsForCall)
}

func (fake *DeliverService) SendMsgCalls(stub func(interface{}) error) {
	fake.sendMsgMutex.Lock()
	defer fake.sendMsgMutex.Unlock()
	fake.SendMsgStub = stub
}

func (fake *DeliverService) SendMsgArgsForCall(i int) interface{} {
	fake.sendMsgMutex.RLock()
	defer fake.sendMsgMutex.RUnlock()
	argsForCall := fake.sendMsgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeliverService) SendMsgReturns(result1 error) {
	fake.sendMsgMutex.Lock()
	defer fake.sendMsgMutex.Unlock()
	fake.SendMsgStub = nil
	fake.sendMsgReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverService) SendMsgReturnsOnCall(i int, result1 error) {
	fake.sendMsgMutex.Lock()
	defer fake.sendMsgMutex.Unlock()
	fake.SendMsgStub = nil
	if fake.sendMsgReturnsOnCall == nil {
		fake.sendMsgReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendMsgReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DeliverService) Trailer() metadata.MD {
	fake.trailerMutex.Lock()
	ret, specificReturn := fake.trailerReturnsOnCall[len(fake.trailerArgsForCall)]
	fake.trailerArgsForCall = append(fake.trailerArgsForCall, struct {
	}{})
	fake.recordInvocation("Trailer", []interface{}{})
	fake.trailerMutex.Unlock()
	if fake.TrailerStub != nil {
		return fake.TrailerStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.trailerReturns
	return fakeReturns.result1
}

func (fake *DeliverService) TrailerCallCount() int {
	fake.trailerMutex.RLock()
	defer fake.trailerMutex.RUnlock()
	return len(fake.trailerArgsForCall)
}

func (fake *DeliverService) TrailerCalls(stub func() metadata.MD) {
	fake.trailerMutex.Lock()
	defer fake.trailerMutex.Unlock()
	fake.TrailerStub = stub
}

func (fake *DeliverService) TrailerReturns(result1 metadata.MD) {
	fake.trailerMutex.Lock()
	defer fake.trailerMutex.Unlock()
	fake.TrailerStub = nil
	fake.trailerReturns = struct {
		result1 metadata.MD
	}{result1}
}

func (fake *DeliverService) TrailerReturnsOnCall(i int, result1 metadata.MD) {
	fake.trailerMutex.Lock()
	defer fake.trailerMutex.Unlock()
	fake.TrailerStub = nil
	if fake.trailerReturnsOnCall == nil {
		fake.trailerReturnsOnCall = make(map[int]struct {
			result1 metadata.MD
		})
	}
	fake.trailerReturnsOnCall[i] = struct {
		result1 metadata.MD
	}{result1}
}

func (fake *DeliverService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeSendMutex.RLock()
	defer fake.closeSendMutex.RUnlock()
	fake.contextMutex.RLock()
	defer fake.contextMutex.RUnlock()
	fake.headerMutex.RLock()
	defer fake.headerMutex.RUnlock()
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	fake.recvMsgMutex.RLock()
	defer fake.recvMsgMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	fake.sendMsgMutex.RLock()
	defer fake.sendMsgMutex.RUnlock()
	fake.trailerMutex.RLock()
	defer fake.trailerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DeliverService) recordInvocation(key string, args []interface{}) {
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
