/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/msp"
	"google.golang.org/grpc/credentials"
)

var commLogger = flogging.MustGetLogger("comm")


type CredentialSupport struct {
	mutex                 sync.RWMutex
	appRootCAsByChain     map[string][][]byte
	ordererRootCAsByChain map[string][][]byte
	serverRootCAs         [][]byte
	clientCert            tls.Certificate
}


func NewCredentialSupport(rootCAs ...[]byte) *CredentialSupport {
	return &CredentialSupport{
		appRootCAsByChain:     make(map[string][][]byte),
		ordererRootCAsByChain: make(map[string][][]byte),
		serverRootCAs:         rootCAs,
	}
}



func (cs *CredentialSupport) SetClientCertificate(cert tls.Certificate) {
	cs.mutex.Lock()
	cs.clientCert = cert
	cs.mutex.Unlock()
}


func (cs *CredentialSupport) GetClientCertificate() tls.Certificate {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	return cs.clientCert
}




func (cs *CredentialSupport) GetDeliverServiceCredentials(channelID string) (credentials.TransportCredentials, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	rootCACerts, exists := cs.ordererRootCAsByChain[channelID]
	if !exists {
		commLogger.Errorf("Attempted to obtain root CA certs of an unknown channel: %s", channelID)
		return nil, fmt.Errorf("didn't find any root CA certs for channel %s", channelID)
	}

	certPool := x509.NewCertPool()
	for _, cert := range rootCACerts {
		err := AddPemToCertPool(cert, certPool)
		if err != nil {
			commLogger.Warningf("Failed to add root cert to credentials (%s)", err)
			continue
		}
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cs.clientCert},
		RootCAs:      certPool,
	}), nil
}



func (cs *CredentialSupport) GetPeerCredentials() credentials.TransportCredentials {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	var appRootCAs [][]byte
	appRootCAs = append(appRootCAs, cs.serverRootCAs...)
	for _, appRootCA := range cs.appRootCAsByChain {
		appRootCAs = append(appRootCAs, appRootCA...)
	}

	certPool := x509.NewCertPool()
	for _, appRootCA := range appRootCAs {
		err := AddPemToCertPool(appRootCA, certPool)
		if err != nil {
			commLogger.Warningf("Failed adding certificates to peer's client TLS trust pool: %s", err)
		}
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cs.clientCert},
		RootCAs:      certPool,
	})
}

func (cs *CredentialSupport) AppRootCAsByChain() map[string][][]byte {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	return cs.appRootCAsByChain
}




func (cs *CredentialSupport) BuildTrustedRootsForChain(cm channelconfig.Resources) {
	appOrgMSPs := make(map[string]struct{})
	if ac, ok := cm.ApplicationConfig(); ok {
		for _, appOrg := range ac.Organizations() {
			appOrgMSPs[appOrg.MSPID()] = struct{}{}
		}
	}

	ordOrgMSPs := make(map[string]struct{})
	if ac, ok := cm.OrdererConfig(); ok {
		for _, ordOrg := range ac.Organizations() {
			ordOrgMSPs[ordOrg.MSPID()] = struct{}{}
		}
	}

	cid := cm.ConfigtxValidator().ChannelID()
	msps, err := cm.MSPManager().GetMSPs()
	if err != nil {
		commLogger.Errorf("Error getting root CAs for channel %s (%s)", cid, err)
		return
	}

	var appRootCAs, ordererRootCAs [][]byte
	for k, v := range msps {
		
		if v.GetType() != msp.FABRIC {
			continue
		}

		for _, root := range v.GetTLSRootCerts() {
			
			if _, ok := appOrgMSPs[k]; ok {
				commLogger.Debugf("adding app root CAs for MSP [%s]", k)
				appRootCAs = append(appRootCAs, root)
			}
			
			if _, ok := ordOrgMSPs[k]; ok {
				commLogger.Debugf("adding orderer root CAs for MSP [%s]", k)
				ordererRootCAs = append(ordererRootCAs, root)
			}
		}
		for _, intermediate := range v.GetTLSIntermediateCerts() {
			
			if _, ok := appOrgMSPs[k]; ok {
				commLogger.Debugf("adding app root CAs for MSP [%s]", k)
				appRootCAs = append(appRootCAs, intermediate)
			}
			
			if _, ok := ordOrgMSPs[k]; ok {
				commLogger.Debugf("adding orderer root CAs for MSP [%s]", k)
				ordererRootCAs = append(ordererRootCAs, intermediate)
			}
		}
	}

	cs.mutex.Lock()
	cs.appRootCAsByChain[cid] = appRootCAs
	cs.ordererRootCAsByChain[cid] = ordererRootCAs
	cs.mutex.Unlock()
}
