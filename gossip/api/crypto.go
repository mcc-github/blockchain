/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"time"

	"github.com/mcc-github/blockchain/gossip/common"
	"google.golang.org/grpc"
)





type MessageCryptoService interface {

	
	
	
	
	GetPKIidOfCert(peerIdentity PeerIdentityType) common.PKIidType

	
	
	
	VerifyBlock(chainID common.ChainID, seqNum uint64, signedBlock []byte) error

	
	
	Sign(msg []byte) ([]byte, error)

	
	
	
	Verify(peerIdentity PeerIdentityType, signature, message []byte) error

	
	
	
	
	VerifyByChannel(chainID common.ChainID, peerIdentity PeerIdentityType, signature, message []byte) error

	
	
	
	ValidateIdentity(peerIdentity PeerIdentityType) error

	
	
	
	
	
	
	
	Expiration(peerIdentity PeerIdentityType) (time.Time, error)
}



type PeerIdentityInfo struct {
	PKIId        common.PKIidType
	Identity     PeerIdentityType
	Organization OrgIdentityType
}


type PeerIdentitySet []PeerIdentityInfo


func (pis PeerIdentitySet) ByOrg() map[string]PeerIdentitySet {
	m := make(map[string]PeerIdentitySet)
	for _, id := range pis {
		m[string(id.Organization)] = append(m[string(id.Organization)], id)
	}
	return m
}


func (pis PeerIdentitySet) ByID() map[string]PeerIdentityInfo {
	m := make(map[string]PeerIdentityInfo)
	for _, id := range pis {
		m[string(id.PKIId)] = id
	}
	return m
}


type PeerIdentityType []byte



type PeerSuspector func(identity PeerIdentityType) bool



type PeerSecureDialOpts func() []grpc.DialOption



type PeerSignature struct {
	Signature    []byte
	Message      []byte
	PeerIdentity PeerIdentityType
}
