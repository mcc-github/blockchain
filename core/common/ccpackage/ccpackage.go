/*
Copyright IBM Corp. 2016-2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ccpackage

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
)


func ExtractSignedCCDepSpec(env *common.Envelope) (*common.ChannelHeader, *peer.SignedChaincodeDeploymentSpec, error) {
	p := &common.Payload{}
	err := proto.Unmarshal(env.Payload, p)
	if err != nil {
		return nil, nil, err
	}
	ch := &common.ChannelHeader{}
	err = proto.Unmarshal(p.Header.ChannelHeader, ch)
	if err != nil {
		return nil, nil, err
	}

	sp := &peer.SignedChaincodeDeploymentSpec{}
	err = proto.Unmarshal(p.Data, sp)
	if err != nil {
		return nil, nil, err
	}

	return ch, sp, nil
}








func ValidateCip(baseCip, otherCip *peer.SignedChaincodeDeploymentSpec) error {
	if baseCip == nil || otherCip == nil {
		panic("do not call with nil parameters")
	}

	if (baseCip.OwnerEndorsements == nil && otherCip.OwnerEndorsements != nil) || (baseCip.OwnerEndorsements != nil && otherCip.OwnerEndorsements == nil) {
		return fmt.Errorf("endorsements should either be both nil or not nil")
	}

	bN := len(baseCip.OwnerEndorsements)
	oN := len(otherCip.OwnerEndorsements)
	if bN > 1 || oN > 1 {
		return fmt.Errorf("expect utmost 1 endorsement from a owner")
	}

	if bN != oN {
		return fmt.Errorf("Rule-all packages should be endorsed or none should be endorsed failed for (%d, %d)", bN, oN)
	}

	if !bytes.Equal(baseCip.ChaincodeDeploymentSpec, otherCip.ChaincodeDeploymentSpec) {
		return fmt.Errorf("Rule-all deployment specs should match(%d, %d)", len(baseCip.ChaincodeDeploymentSpec), len(otherCip.ChaincodeDeploymentSpec))
	}

	if !bytes.Equal(baseCip.InstantiationPolicy, otherCip.InstantiationPolicy) {
		return fmt.Errorf("Rule-all instantiation policies should match(%d, %d)", len(baseCip.InstantiationPolicy), len(otherCip.InstantiationPolicy))
	}

	return nil
}

func createSignedCCDepSpec(cdsbytes []byte, instpolicybytes []byte, endorsements []*peer.Endorsement) (*common.Envelope, error) {
	if cdsbytes == nil {
		return nil, fmt.Errorf("nil chaincode deployment spec")
	}

	if instpolicybytes == nil {
		return nil, fmt.Errorf("nil instantiation policy")
	}

	
	cip := &peer.SignedChaincodeDeploymentSpec{ChaincodeDeploymentSpec: cdsbytes, InstantiationPolicy: instpolicybytes, OwnerEndorsements: endorsements}

	
	cipbytes := utils.MarshalOrPanic(cip)

	
	msgVersion := int32(0)
	epoch := uint64(0)
	chdr := utils.MakeChannelHeader(common.HeaderType_CHAINCODE_PACKAGE, msgVersion, "", epoch)

	
	payl := &common.Payload{Header: &common.Header{ChannelHeader: utils.MarshalOrPanic(chdr)}, Data: cipbytes}
	paylBytes, err := utils.GetBytesPayload(payl)
	if err != nil {
		return nil, err
	}

	
	return &common.Envelope{Payload: paylBytes}, nil
}




func CreateSignedCCDepSpecForInstall(pack []*common.Envelope) (*common.Envelope, error) {
	if len(pack) == 0 {
		return nil, errors.New("no packages provided to collate")
	}

	
	
	
	var baseCip *peer.SignedChaincodeDeploymentSpec
	var err error
	var endorsementExists bool
	var endorsements []*peer.Endorsement
	for n, r := range pack {
		p := &common.Payload{}
		if err = proto.Unmarshal(r.Payload, p); err != nil {
			return nil, err
		}

		cip := &peer.SignedChaincodeDeploymentSpec{}
		if err = proto.Unmarshal(p.Data, cip); err != nil {
			return nil, err
		}

		
		
		if n == 0 {
			baseCip = cip
			
			if len(cip.OwnerEndorsements) > 0 {
				endorsements = make([]*peer.Endorsement, len(pack))
			}

		} else if err = ValidateCip(baseCip, cip); err != nil {
			return nil, err
		}

		if endorsementExists {
			endorsements[n] = cip.OwnerEndorsements[0]
		}
	}

	return createSignedCCDepSpec(baseCip.ChaincodeDeploymentSpec, baseCip.InstantiationPolicy, endorsements)
}



func OwnerCreateSignedCCDepSpec(cds *peer.ChaincodeDeploymentSpec, instPolicy *common.SignaturePolicyEnvelope, owner msp.SigningIdentity) (*common.Envelope, error) {
	if cds == nil {
		return nil, fmt.Errorf("invalid chaincode deployment spec")
	}

	if instPolicy == nil {
		return nil, fmt.Errorf("must provide an instantiation policy")
	}

	cdsbytes := utils.MarshalOrPanic(cds)

	instpolicybytes := utils.MarshalOrPanic(instPolicy)

	var endorsements []*peer.Endorsement
	
	
	
	if owner != nil {
		
		endorser, err := owner.Serialize()
		if err != nil {
			return nil, fmt.Errorf("Could not serialize the signing identity for %s, err %s", owner.GetIdentifier(), err)
		}

		
		signature, err := owner.Sign(append(cdsbytes, append(instpolicybytes, endorser...)...))
		if err != nil {
			return nil, fmt.Errorf("Could not sign the ccpackage, err %s", err)
		}

		
		
		
		endorsements = make([]*peer.Endorsement, 1)

		endorsements[0] = &peer.Endorsement{Signature: signature, Endorser: endorser}
	}

	return createSignedCCDepSpec(cdsbytes, instpolicybytes, endorsements)
}


func SignExistingPackage(env *common.Envelope, owner msp.SigningIdentity) (*common.Envelope, error) {
	if owner == nil {
		return nil, fmt.Errorf("owner not provided")
	}

	ch, sdepspec, err := ExtractSignedCCDepSpec(env)
	if err != nil {
		return nil, err
	}

	if ch == nil {
		return nil, fmt.Errorf("channel header not found in the envelope")
	}

	if sdepspec == nil || sdepspec.ChaincodeDeploymentSpec == nil || sdepspec.InstantiationPolicy == nil || sdepspec.OwnerEndorsements == nil {
		return nil, fmt.Errorf("invalid signed deployment spec")
	}

	
	endorser, err := owner.Serialize()
	if err != nil {
		return nil, fmt.Errorf("Could not serialize the signing identity for %s, err %s", owner.GetIdentifier(), err)
	}

	
	signature, err := owner.Sign(append(sdepspec.ChaincodeDeploymentSpec, append(sdepspec.InstantiationPolicy, endorser...)...))
	if err != nil {
		return nil, fmt.Errorf("Could not sign the ccpackage, err %s", err)
	}

	endorsements := append(sdepspec.OwnerEndorsements, &peer.Endorsement{Signature: signature, Endorser: endorser})

	return createSignedCCDepSpec(sdepspec.ChaincodeDeploymentSpec, sdepspec.InstantiationPolicy, endorsements)
}
