/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/msp/mgmt"
)

var saLogger = flogging.MustGetLogger("peer.gossip.sa")









type mspSecurityAdvisor struct {
	deserializer mgmt.DeserializersManager
}



func NewSecurityAdvisor(deserializer mgmt.DeserializersManager) api.SecurityAdvisor {
	return &mspSecurityAdvisor{deserializer: deserializer}
}






func (advisor *mspSecurityAdvisor) OrgByPeerIdentity(peerIdentity api.PeerIdentityType) api.OrgIdentityType {
	
	if len(peerIdentity) == 0 {
		saLogger.Error("Invalid Peer Identity. It must be different from nil.")

		return nil
	}

	
	

	
	
	
	
	

	
	identity, err := advisor.deserializer.GetLocalDeserializer().DeserializeIdentity([]byte(peerIdentity))
	if err == nil {
		return []byte(identity.GetMSPIdentifier())
	}

	
	for chainID, mspManager := range advisor.deserializer.GetChannelDeserializers() {
		
		identity, err := mspManager.DeserializeIdentity([]byte(peerIdentity))
		if err != nil {
			saLogger.Debugf("Failed deserialization identity [% x] on [%s]: [%s]", peerIdentity, chainID, err)
			continue
		}

		return []byte(identity.GetMSPIdentifier())
	}

	saLogger.Warningf("Peer Identity [% x] cannot be desirialized. No MSP found able to do that.", peerIdentity)

	return nil
}
