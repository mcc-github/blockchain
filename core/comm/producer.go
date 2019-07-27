/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"google.golang.org/grpc"
)

var logger = flogging.MustGetLogger("ConnProducer")


type ConnectionFactory func(endpoint string, connectionTimeout time.Duration) (*grpc.ClientConn, error)



type EndpointCriteria struct {
	Endpoint      string
	Organizations []string
}


func (ec EndpointCriteria) Equals(other EndpointCriteria) bool {
	ss1 := stringSet(ec.Organizations)
	ss2 := stringSet(other.Organizations)
	return ec.Endpoint == other.Endpoint && ss1.equals(ss2)
}



type stringSet []string

func (ss stringSet) equals(ss2 stringSet) bool {
	
	if len(ss) != len(ss2) {
		return false
	}
	return reflect.DeepEqual(ss.toMap(), ss2.toMap())
}

func (ss stringSet) toMap() map[string]struct{} {
	m := make(map[string]struct{})
	for _, s := range ss {
		m[s] = struct{}{}
	}
	return m
}



type ConnectionProducer interface {
	
	
	
	NewConnection() (*grpc.ClientConn, string, error)
	
	
	UpdateEndpoints(endpoints []string)
	
	GetEndpoints() []string
}

type connProducer struct {
	sync.RWMutex
	endpoints             []string
	connect               ConnectionFactory
	nextEndpointIndex     int
	deliverClientDialOpts []grpc.DialOption
	peerTLSEnabled        bool
	connectionTimeout     time.Duration
}



func NewConnectionProducer(
	factory ConnectionFactory,
	endpoints []string,
	deliverClientDialOpts []grpc.DialOption,
	peerTLSEnabled bool,
	connectionTimeout time.Duration,
) ConnectionProducer {
	if len(endpoints) == 0 {
		return nil
	}
	return &connProducer{
		endpoints:             shuffle(endpoints),
		connect:               factory,
		deliverClientDialOpts: deliverClientDialOpts,
		peerTLSEnabled:        peerTLSEnabled,
		connectionTimeout:     connectionTimeout,
	}
}




func (cp *connProducer) NewConnection() (*grpc.ClientConn, string, error) {
	cp.Lock()
	defer cp.Unlock()

	logger.Debugf("Creating a new connection")

	for i := 0; i < len(cp.endpoints); i++ {
		currentEndpoint := cp.endpoints[cp.nextEndpointIndex]
		conn, err := cp.connect(currentEndpoint, cp.connectionTimeout)
		cp.nextEndpointIndex = (cp.nextEndpointIndex + 1) % len(cp.endpoints)
		if err != nil {
			logger.Error("Failed connecting to", currentEndpoint, ", error:", err)
			continue
		}
		logger.Debugf("Connected to %s", currentEndpoint)
		return conn, currentEndpoint, nil
	}

	logger.Errorf("Could not connect to any of the endpoints: %v", cp.endpoints)

	return nil, "", fmt.Errorf("could not connect to any of the endpoints: %v", cp.endpoints)
}



func (cp *connProducer) UpdateEndpoints(endpoints []string) {
	if len(endpoints) == 0 {
		
		return
	}
	cp.Lock()
	defer cp.Unlock()

	cp.nextEndpointIndex = 0
	cp.endpoints = endpoints
}

func shuffle(a []string) []string {
	n := len(a)
	returnedSlice := make([]string, n)
	rand.Seed(time.Now().UnixNano())
	indices := rand.Perm(n)
	for i, idx := range indices {
		returnedSlice[i] = a[idx]
	}
	return returnedSlice
}


func (cp *connProducer) GetEndpoints() []string {
	cp.RLock()
	defer cp.RUnlock()
	return cp.endpoints
}
