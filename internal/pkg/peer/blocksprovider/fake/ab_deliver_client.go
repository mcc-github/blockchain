
package fake

import (
	"context"
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/orderer"
	"google.golang.org/grpc/metadata"
)

type DeliverClient struct {
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
	SendStub        func(*common.Envelope) error
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		arg1 *common.Envelope
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

func (fake *DeliverClient) CloseSend() error {
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

func (fake *DeliverClient) CloseSendCallCount() int {
	fake.closeSendMutex.RLock()
	defer fake.closeSendMutex.RUnlock()
	return len(fake.closeSendArgsForCall)
}

func (fake *DeliverClient) CloseSendCalls(stub func() error) {
	fake.closeSendMutex.Lock()
	defer fake.closeSendMutex.Unlock()
	fake.CloseSendStub = stub
}

func (fake *DeliverClient) CloseSendReturns(result1 error) {
	fake.closeSendMutex.Lock()
	defer fake.closeSendMutex.Unlock()
	fake.CloseSendStub = nil
	fake.closeSendReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverClient) CloseSendReturnsOnCall(i int, result1 error) {
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

func (fake *DeliverClient) Context() context.Context {
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

func (fake *DeliverClient) ContextCallCount() int {
	fake.contextMutex.RLock()
	defer fake.contextMutex.RUnlock()
	return len(fake.contextArgsForCall)
}

func (fake *DeliverClient) ContextCalls(stub func() context.Context) {
	fake.contextMutex.Lock()
	defer fake.contextMutex.Unlock()
	fake.ContextStub = stub
}

func (fake *DeliverClient) ContextReturns(result1 context.Context) {
	fake.contextMutex.Lock()
	defer fake.contextMutex.Unlock()
	fake.ContextStub = nil
	fake.contextReturns = struct {
		result1 context.Context
	}{result1}
}

func (fake *DeliverClient) ContextReturnsOnCall(i int, result1 context.Context) {
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

func (fake *DeliverClient) Header() (metadata.MD, error) {
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

func (fake *DeliverClient) HeaderCallCount() int {
	fake.headerMutex.RLock()
	defer fake.headerMutex.RUnlock()
	return len(fake.headerArgsForCall)
}

func (fake *DeliverClient) HeaderCalls(stub func() (metadata.MD, error)) {
	fake.headerMutex.Lock()
	defer fake.headerMutex.Unlock()
	fake.HeaderStub = stub
}

func (fake *DeliverClient) HeaderReturns(result1 metadata.MD, result2 error) {
	fake.headerMutex.Lock()
	defer fake.headerMutex.Unlock()
	fake.HeaderStub = nil
	fake.headerReturns = struct {
		result1 metadata.MD
		result2 error
	}{result1, result2}
}

func (fake *DeliverClient) HeaderReturnsOnCall(i int, result1 metadata.MD, result2 error) {
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

func (fake *DeliverClient) Recv() (*orderer.DeliverResponse, error) {
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

func (fake *DeliverClient) RecvCallCount() int {
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	return len(fake.recvArgsForCall)
}

func (fake *DeliverClient) RecvCalls(stub func() (*orderer.DeliverResponse, error)) {
	fake.recvMutex.Lock()
	defer fake.recvMutex.Unlock()
	fake.RecvStub = stub
}

func (fake *DeliverClient) RecvReturns(result1 *orderer.DeliverResponse, result2 error) {
	fake.recvMutex.Lock()
	defer fake.recvMutex.Unlock()
	fake.RecvStub = nil
	fake.recvReturns = struct {
		result1 *orderer.DeliverResponse
		result2 error
	}{result1, result2}
}

func (fake *DeliverClient) RecvReturnsOnCall(i int, result1 *orderer.DeliverResponse, result2 error) {
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

func (fake *DeliverClient) RecvMsg(arg1 interface{}) error {
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

func (fake *DeliverClient) RecvMsgCallCount() int {
	fake.recvMsgMutex.RLock()
	defer fake.recvMsgMutex.RUnlock()
	return len(fake.recvMsgArgsForCall)
}

func (fake *DeliverClient) RecvMsgCalls(stub func(interface{}) error) {
	fake.recvMsgMutex.Lock()
	defer fake.recvMsgMutex.Unlock()
	fake.RecvMsgStub = stub
}

func (fake *DeliverClient) RecvMsgArgsForCall(i int) interface{} {
	fake.recvMsgMutex.RLock()
	defer fake.recvMsgMutex.RUnlock()
	argsForCall := fake.recvMsgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeliverClient) RecvMsgReturns(result1 error) {
	fake.recvMsgMutex.Lock()
	defer fake.recvMsgMutex.Unlock()
	fake.RecvMsgStub = nil
	fake.recvMsgReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverClient) RecvMsgReturnsOnCall(i int, result1 error) {
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

func (fake *DeliverClient) Send(arg1 *common.Envelope) error {
	fake.sendMutex.Lock()
	ret, specificReturn := fake.sendReturnsOnCall[len(fake.sendArgsForCall)]
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		arg1 *common.Envelope
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

func (fake *DeliverClient) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *DeliverClient) SendCalls(stub func(*common.Envelope) error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = stub
}

func (fake *DeliverClient) SendArgsForCall(i int) *common.Envelope {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	argsForCall := fake.sendArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeliverClient) SendReturns(result1 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverClient) SendReturnsOnCall(i int, result1 error) {
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

func (fake *DeliverClient) SendMsg(arg1 interface{}) error {
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

func (fake *DeliverClient) SendMsgCallCount() int {
	fake.sendMsgMutex.RLock()
	defer fake.sendMsgMutex.RUnlock()
	return len(fake.sendMsgArgsForCall)
}

func (fake *DeliverClient) SendMsgCalls(stub func(interface{}) error) {
	fake.sendMsgMutex.Lock()
	defer fake.sendMsgMutex.Unlock()
	fake.SendMsgStub = stub
}

func (fake *DeliverClient) SendMsgArgsForCall(i int) interface{} {
	fake.sendMsgMutex.RLock()
	defer fake.sendMsgMutex.RUnlock()
	argsForCall := fake.sendMsgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeliverClient) SendMsgReturns(result1 error) {
	fake.sendMsgMutex.Lock()
	defer fake.sendMsgMutex.Unlock()
	fake.SendMsgStub = nil
	fake.sendMsgReturns = struct {
		result1 error
	}{result1}
}

func (fake *DeliverClient) SendMsgReturnsOnCall(i int, result1 error) {
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

func (fake *DeliverClient) Trailer() metadata.MD {
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

func (fake *DeliverClient) TrailerCallCount() int {
	fake.trailerMutex.RLock()
	defer fake.trailerMutex.RUnlock()
	return len(fake.trailerArgsForCall)
}

func (fake *DeliverClient) TrailerCalls(stub func() metadata.MD) {
	fake.trailerMutex.Lock()
	defer fake.trailerMutex.Unlock()
	fake.TrailerStub = stub
}

func (fake *DeliverClient) TrailerReturns(result1 metadata.MD) {
	fake.trailerMutex.Lock()
	defer fake.trailerMutex.Unlock()
	fake.TrailerStub = nil
	fake.trailerReturns = struct {
		result1 metadata.MD
	}{result1}
}

func (fake *DeliverClient) TrailerReturnsOnCall(i int, result1 metadata.MD) {
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

func (fake *DeliverClient) Invocations() map[string][][]interface{} {
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

func (fake *DeliverClient) recordInvocation(key string, args []interface{}) {
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
