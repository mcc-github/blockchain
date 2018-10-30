
package mock

import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/mcc-github/blockchain/token/client"
	"google.golang.org/grpc"
)

type OrdererClient struct {
	NewBroadcastStub        func(ctx context.Context, opts ...grpc.CallOption) (client.Broadcast, error)
	newBroadcastMutex       sync.RWMutex
	newBroadcastArgsForCall []struct {
		ctx  context.Context
		opts []grpc.CallOption
	}
	newBroadcastReturns struct {
		result1 client.Broadcast
		result2 error
	}
	newBroadcastReturnsOnCall map[int]struct {
		result1 client.Broadcast
		result2 error
	}
	CertificateStub        func() *tls.Certificate
	certificateMutex       sync.RWMutex
	certificateArgsForCall []struct{}
	certificateReturns     struct {
		result1 *tls.Certificate
	}
	certificateReturnsOnCall map[int]struct {
		result1 *tls.Certificate
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *OrdererClient) NewBroadcast(ctx context.Context, opts ...grpc.CallOption) (client.Broadcast, error) {
	fake.newBroadcastMutex.Lock()
	ret, specificReturn := fake.newBroadcastReturnsOnCall[len(fake.newBroadcastArgsForCall)]
	fake.newBroadcastArgsForCall = append(fake.newBroadcastArgsForCall, struct {
		ctx  context.Context
		opts []grpc.CallOption
	}{ctx, opts})
	fake.recordInvocation("NewBroadcast", []interface{}{ctx, opts})
	fake.newBroadcastMutex.Unlock()
	if fake.NewBroadcastStub != nil {
		return fake.NewBroadcastStub(ctx, opts...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newBroadcastReturns.result1, fake.newBroadcastReturns.result2
}

func (fake *OrdererClient) NewBroadcastCallCount() int {
	fake.newBroadcastMutex.RLock()
	defer fake.newBroadcastMutex.RUnlock()
	return len(fake.newBroadcastArgsForCall)
}

func (fake *OrdererClient) NewBroadcastArgsForCall(i int) (context.Context, []grpc.CallOption) {
	fake.newBroadcastMutex.RLock()
	defer fake.newBroadcastMutex.RUnlock()
	return fake.newBroadcastArgsForCall[i].ctx, fake.newBroadcastArgsForCall[i].opts
}

func (fake *OrdererClient) NewBroadcastReturns(result1 client.Broadcast, result2 error) {
	fake.NewBroadcastStub = nil
	fake.newBroadcastReturns = struct {
		result1 client.Broadcast
		result2 error
	}{result1, result2}
}

func (fake *OrdererClient) NewBroadcastReturnsOnCall(i int, result1 client.Broadcast, result2 error) {
	fake.NewBroadcastStub = nil
	if fake.newBroadcastReturnsOnCall == nil {
		fake.newBroadcastReturnsOnCall = make(map[int]struct {
			result1 client.Broadcast
			result2 error
		})
	}
	fake.newBroadcastReturnsOnCall[i] = struct {
		result1 client.Broadcast
		result2 error
	}{result1, result2}
}

func (fake *OrdererClient) Certificate() *tls.Certificate {
	fake.certificateMutex.Lock()
	ret, specificReturn := fake.certificateReturnsOnCall[len(fake.certificateArgsForCall)]
	fake.certificateArgsForCall = append(fake.certificateArgsForCall, struct{}{})
	fake.recordInvocation("Certificate", []interface{}{})
	fake.certificateMutex.Unlock()
	if fake.CertificateStub != nil {
		return fake.CertificateStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.certificateReturns.result1
}

func (fake *OrdererClient) CertificateCallCount() int {
	fake.certificateMutex.RLock()
	defer fake.certificateMutex.RUnlock()
	return len(fake.certificateArgsForCall)
}

func (fake *OrdererClient) CertificateReturns(result1 *tls.Certificate) {
	fake.CertificateStub = nil
	fake.certificateReturns = struct {
		result1 *tls.Certificate
	}{result1}
}

func (fake *OrdererClient) CertificateReturnsOnCall(i int, result1 *tls.Certificate) {
	fake.CertificateStub = nil
	if fake.certificateReturnsOnCall == nil {
		fake.certificateReturnsOnCall = make(map[int]struct {
			result1 *tls.Certificate
		})
	}
	fake.certificateReturnsOnCall[i] = struct {
		result1 *tls.Certificate
	}{result1}
}

func (fake *OrdererClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newBroadcastMutex.RLock()
	defer fake.newBroadcastMutex.RUnlock()
	fake.certificateMutex.RLock()
	defer fake.certificateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *OrdererClient) recordInvocation(key string, args []interface{}) {
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

var _ client.OrdererClient = new(OrdererClient)
