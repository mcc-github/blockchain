/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster_test

import (
	"sync"
	"testing"

	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/cluster/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestConcurrentConnections(t *testing.T) {
	t.Parallel()
	
	
	
	n := 100
	var wg sync.WaitGroup
	wg.Add(n)
	dialer := &mocks.SecureDialer{}
	conn := &grpc.ClientConn{}
	dialer.On("Dial", mock.Anything, mock.Anything).Return(conn, nil)
	connStore := cluster.NewConnectionStore(dialer, &disabled.Gauge{})
	connect := func() {
		defer wg.Done()
		conn2, err := connStore.Connection("", nil)
		assert.NoError(t, err)
		assert.True(t, conn2 == conn)
	}
	for i := 0; i < n; i++ {
		go connect()
	}
	wg.Wait()
	dialer.AssertNumberOfCalls(t, "Dial", 1)
}

type connectionMapperSpy struct {
	lookupDelay   chan struct{}
	lookupInvoked chan struct{}
	cluster.ConnectionMapper
}

func (cms *connectionMapperSpy) Lookup(cert []byte) (*grpc.ClientConn, bool) {
	
	cms.lookupInvoked <- struct{}{}
	
	
	
	<-cms.lookupDelay
	return cms.ConnectionMapper.Lookup(cert)
}

func TestConcurrentLookupMiss(t *testing.T) {
	t.Parallel()
	
	
	
	
	

	dialer := &mocks.SecureDialer{}
	conn := &grpc.ClientConn{}
	dialer.On("Dial", mock.Anything, mock.Anything).Return(conn, nil)

	connStore := cluster.NewConnectionStore(dialer, &disabled.Gauge{})
	
	spy := &connectionMapperSpy{
		ConnectionMapper: connStore.Connections,
		lookupDelay:      make(chan struct{}, 2),
		lookupInvoked:    make(chan struct{}, 2),
	}
	connStore.Connections = spy

	var goroutinesExited sync.WaitGroup
	goroutinesExited.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer goroutinesExited.Done()
			conn2, err := connStore.Connection("", nil)
			assert.NoError(t, err)
			
			
			assert.True(t, conn2 == conn)
		}()
	}
	
	
	<-spy.lookupInvoked
	<-spy.lookupInvoked
	
	spy.lookupDelay <- struct{}{}
	spy.lookupDelay <- struct{}{}
	
	close(spy.lookupDelay)
	
	goroutinesExited.Wait()
}
