/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package orderers

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"math/rand"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"

	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("core.orderers")

type ConnectionSource struct {
	mutex              sync.RWMutex
	allEndpoints       []*Endpoint
	orgToEndpointsHash map[string][]byte
}

type Endpoint struct {
	Address   string
	CertPool  *x509.CertPool
	Refreshed chan struct{}
}

type OrdererOrg struct {
	Addresses []string
	RootCerts [][]byte
}

func NewConnectionSource() *ConnectionSource {
	return &ConnectionSource{
		orgToEndpointsHash: map[string][]byte{},
	}
}

func (cs *ConnectionSource) RandomEndpoint() (*Endpoint, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	if len(cs.allEndpoints) == 0 {
		return nil, errors.Errorf("no endpoints currently defined")
	}
	return cs.allEndpoints[rand.Intn(len(cs.allEndpoints))], nil
}

func (cs *ConnectionSource) Update(globalAddrs []string, orgs map[string]OrdererOrg) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	anyChange := false
	hasOrgEndpoints := false
	for orgName, org := range orgs {
		hasher := sha256.New()
		for _, cert := range org.RootCerts {
			hasher.Write(cert)
		}
		for _, address := range org.Addresses {
			hasOrgEndpoints = true
			hasher.Write([]byte(address))
		}
		hash := hasher.Sum(nil)

		lastHash, ok := cs.orgToEndpointsHash[orgName]
		cs.orgToEndpointsHash[orgName] = hash
		if ok && bytes.Equal(hash, lastHash) {
			continue
		}

		anyChange = true
	}

	for orgName := range cs.orgToEndpointsHash {
		if _, ok := orgs[orgName]; !ok {
			
			anyChange = true
		}
	}

	if hasOrgEndpoints && len(globalAddrs) > 0 {
		logger.Warning("Config defines both orderer org specific endpoints and global endpoints, global endpoints will be ignored")
	}

	if !hasOrgEndpoints && len(globalAddrs) != len(cs.allEndpoints) {
		
		anyChange = true
	}

	if !hasOrgEndpoints && !anyChange && len(globalAddrs) == len(cs.allEndpoints) {
		
		
		

		newAddresses := map[string]struct{}{}
		for _, address := range globalAddrs {
			newAddresses[address] = struct{}{}
		}

		for _, endpoint := range cs.allEndpoints {
			delete(newAddresses, endpoint.Address)
		}

		
		
		anyChange = len(newAddresses) != 0
	}

	if !anyChange {
		
		
		
		return
	}

	for _, endpoint := range cs.allEndpoints {
		
		
		
		
		
		
		
		close(endpoint.Refreshed)
	}

	cs.allEndpoints = nil

	globalCertPool := x509.NewCertPool()

	for orgName, org := range orgs {
		certPool := x509.NewCertPool()
		for _, rootCert := range org.RootCerts {
			if hasOrgEndpoints {
				if err := comm.AddPemToCertPool(rootCert, certPool); err != nil {
					logger.Warningf("Could not add orderer cert for org '%s' to cert pool: %s", orgName, err)
				}
			} else {
				if err := comm.AddPemToCertPool(rootCert, globalCertPool); err != nil {
					logger.Warningf("Could not add orderer cert for org '%s' to global cert pool: %s", orgName, err)

				}
			}
		}

		
		
		for _, address := range org.Addresses {
			cs.allEndpoints = append(cs.allEndpoints, &Endpoint{
				Address:   address,
				CertPool:  certPool,
				Refreshed: make(chan struct{}),
			})
		}
	}

	if len(cs.allEndpoints) != 0 {
		
		
		return
	}

	for _, address := range globalAddrs {
		cs.allEndpoints = append(cs.allEndpoints, &Endpoint{
			Address:   address,
			CertPool:  globalCertPool,
			Refreshed: make(chan struct{}),
		})
	}
}
