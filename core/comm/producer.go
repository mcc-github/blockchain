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


type ConnectionFactory func(endpoint string) (*grpc.ClientConn, error)



type EndpointCriteria struct {
	Endpoint      string
	Organizations []string
}


func (ec EndpointCriteria) Equals(other EndpointCriteria) bool {
	return ec.Endpoint == other.Endpoint && equalStringSets(ec.Organizations, other.Organizations)
}

func equalStringSets(s1, s2 []string) bool {
	
	if len(s1) != len(s2) {
		return false
	}

	return reflect.DeepEqual(stringSliceToSet(s1), stringSliceToSet(s2))
}

func stringSliceToSet(set []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, s := range set {
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
	endpoints         []string
	connect           ConnectionFactory
	nextEndpointIndex int
}



func NewConnectionProducer(factory ConnectionFactory, endpoints []string) ConnectionProducer {
	if len(endpoints) == 0 {
		return nil
	}
	return &connProducer{endpoints: shuffle(endpoints), connect: factory}
}




func (cp *connProducer) NewConnection() (*grpc.ClientConn, string, error) {
	cp.Lock()
	defer cp.Unlock()

	logger.Debugf("Creating a new connection")

	for i := 0; i < len(cp.endpoints); i++ {
		currentEndpoint := cp.endpoints[cp.nextEndpointIndex]
		conn, err := cp.connect(currentEndpoint)
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
