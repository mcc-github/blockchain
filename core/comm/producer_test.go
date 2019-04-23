/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"fmt"
	"testing"

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
	assert.NoError(t, err)
	assert.Equal(t, "b", conn2Endpoint[fmt.Sprintf("%p", conn)])
}

func TestNewConnectionRoundRobin(t *testing.T) {
	t.Parallel()

	

	totalEndpoints := []string{"a", "b", "c", "d"}
	conn2Endpoint := make(map[string]string)
	connFactory := func(endpoint string) (*grpc.ClientConn, error) {
		conn := &grpc.ClientConn{}
		conn2Endpoint[fmt.Sprintf("%p", conn)] = endpoint
		return conn, nil
	}

	producer := NewConnectionProducer(connFactory, totalEndpoints)
	connectedEndpoints := make(map[string]struct{})

	assertAllEndpointsUsed := func() {
		for i := 0; i < len(totalEndpoints); i++ {
			_, endpoint, err := producer.NewConnection()
			assert.NoError(t, err)
			connectedEndpoints[endpoint] = struct{}{}
		}
		assert.Len(t, connectedEndpoints, 4)
	}

	assertAllEndpointsUsed()
	
	connectedEndpoints = make(map[string]struct{})
	
	assertAllEndpointsUsed()
}

func TestEndpointCriteria(t *testing.T) {
	endpointCriteria := EndpointCriteria{
		Endpoint:      "a",
		Organizations: []string{"o1", "o2"},
	}

	for _, testCase := range []struct {
		description           string
		expectedEqual         bool
		otherEndpointCriteria EndpointCriteria
	}{
		{
			description: "different endpoint",
			otherEndpointCriteria: EndpointCriteria{
				Endpoint:      "b",
				Organizations: []string{"o1", "o2"},
			},
		},
		{
			description: "more organizations",
			otherEndpointCriteria: EndpointCriteria{
				Endpoint:      "a",
				Organizations: []string{"o1", "o2", "o3"},
			},
		},
		{
			description: "less organizations",
			otherEndpointCriteria: EndpointCriteria{
				Endpoint:      "a",
				Organizations: []string{"o1"},
			},
		},
		{
			description: "different organizations",
			otherEndpointCriteria: EndpointCriteria{
				Endpoint:      "a",
				Organizations: []string{"o1", "o3"},
			},
		},
		{
			description:   "permuted organizations",
			expectedEqual: true,
			otherEndpointCriteria: EndpointCriteria{
				Endpoint:      "a",
				Organizations: []string{"o2", "o1"},
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			assert.Equal(t, testCase.expectedEqual, endpointCriteria.Equals(testCase.otherEndpointCriteria))
		})
	}
}
