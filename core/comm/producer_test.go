/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestEmptyEndpoints(t *testing.T) {
	t.Parallel()
	noopFactory := func(endpoint string) (*grpc.ClientConn, error) {
		return nil, nil
	}
	assert.Nil(t, NewConnectionProducer(noopFactory, []string{}))
}

func TestConnFailures(t *testing.T) {
	t.Parallel()
	conn2Endpoint := make(map[string]string)
	shouldConnFail := map[string]bool{
		"a": true,
		"b": false,
		"c": false,
	}
	connFactory := func(endpoint string) (*grpc.ClientConn, error) {
		conn := &grpc.ClientConn{}
		conn2Endpoint[fmt.Sprintf("%p", conn)] = endpoint
		if !shouldConnFail[endpoint] {
			return conn, nil
		}
		return nil, fmt.Errorf("Failed connecting to %s", endpoint)
	}
	
	producer := NewConnectionProducer(connFactory, []string{"a", "b", "c"})
	conn, _, err := producer.NewConnection()
	assert.NoError(t, err)
	
	assert.NotEqual(t, "a", conn2Endpoint[fmt.Sprintf("%p", conn)])
	
	shouldConnFail["a"] = false
	
	selected := make(map[string]struct{})
	for i := 0; i < 1000; i++ {
		conn, _, err := producer.NewConnection()
		assert.NoError(t, err)
		selected[conn2Endpoint[fmt.Sprintf("%p", conn)]] = struct{}{}
	}
	
	_, isAselected := selected["a"]
	_, isBselected := selected["b"]
	_, isCselected := selected["c"]
	assert.True(t, isBselected)
	assert.True(t, isCselected)
	assert.True(t, isAselected)

	
	shouldConnFail["a"] = true
	shouldConnFail["b"] = true
	shouldConnFail["c"] = true
	conn, _, err = producer.NewConnection()
	assert.Nil(t, conn)
	assert.Error(t, err)
}

func TestUpdateEndpoints(t *testing.T) {
	t.Parallel()
	conn2Endpoint := make(map[string]string)
	connFactory := func(endpoint string) (*grpc.ClientConn, error) {
		conn := &grpc.ClientConn{}
		conn2Endpoint[fmt.Sprintf("%p", conn)] = endpoint
		return conn, nil
	}
	
	producer := NewConnectionProducer(connFactory, []string{"a"})
	conn, a, err := producer.NewConnection()
	assert.NoError(t, err)
	assert.Equal(t, "a", conn2Endpoint[fmt.Sprintf("%p", conn)])
	assert.Equal(t, "a", a)
	
	
	producer.UpdateEndpoints([]string{"b"})
	conn, b, err := producer.NewConnection()
	assert.NoError(t, err)
	assert.NotEqual(t, "a", conn2Endpoint[fmt.Sprintf("%p", conn)])
	assert.Equal(t, "b", conn2Endpoint[fmt.Sprintf("%p", conn)])
	assert.Equal(t, "b", b)
	
	producer.UpdateEndpoints([]string{})
	conn, _, err = producer.NewConnection()
	assert.Equal(t, "b", conn2Endpoint[fmt.Sprintf("%p", conn)])
}

func TestDisableEndpoint(t *testing.T) {
	orgEndpointDisableInterval := EndpointDisableInterval
	EndpointDisableInterval = time.Millisecond * 100
	defer func() { EndpointDisableInterval = orgEndpointDisableInterval }()

	conn2Endpoint := make(map[string]string)
	connFactory := func(endpoint string) (*grpc.ClientConn, error) {
		conn := &grpc.ClientConn{}
		conn2Endpoint[fmt.Sprintf("%p", conn)] = endpoint
		return conn, nil
	}
	
	producer := NewConnectionProducer(connFactory, []string{"a"})
	conn, a, err := producer.NewConnection()
	assert.NoError(t, err)
	assert.Equal(t, "a", conn2Endpoint[fmt.Sprintf("%p", conn)])
	assert.Equal(t, "a", a)
	
	producer.DisableEndpoint("a")
	_, _, err = producer.NewConnection()
	
	assert.NoError(t, err)
	
	producer.UpdateEndpoints([]string{"a", "b"})
	
	producer.DisableEndpoint("a")
	conn, a, err = producer.NewConnection()
	assert.NoError(t, err)
	
	assert.Equal(t, "b", conn2Endpoint[fmt.Sprintf("%p", conn)])
	assert.Equal(t, "b", a)

}
