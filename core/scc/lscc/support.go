/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lscc

import (
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

type supportImpl struct {
	GetMSPIDs MSPIDsGetter
}



func (s *supportImpl) PutChaincodeToLocalStorage(ccpack ccprovider.CCPackage) error {
	if err := ccpack.PutChaincodeToFS(); err != nil {
		return errors.Errorf("error installing chaincode code %s:%s(%s)", ccpack.GetChaincodeData().Name, ccpack.GetChaincodeData().Version, err)
	}

	return nil
}



func (s *supportImpl) GetChaincodeFromLocalStorage(ccNameVersion string) (ccprovider.CCPackage, error) {
	return ccprovider.GetChaincodeFromFS(ccNameVersion)
}



func (s *supportImpl) GetChaincodesFromLocalStorage() (*pb.ChaincodeQueryResponse, error) {
	return ccprovider.GetInstalledChaincodes()
}



func (s *supportImpl) GetInstantiationPolicy(channel string, ccpack ccprovider.CCPackage) ([]byte, error) {
	var ip []byte
	var err error
	
	sccpack, isSccpack := ccpack.(*ccprovider.SignedCDSPackage)
	if isSccpack {
		ip = sccpack.GetInstantiationPolicy()
		if ip == nil {
			return nil, errors.Errorf("instantiation policy cannot be nil for a SignedCCDeploymentSpec")
		}
	} else {
		
		
		mspids := s.GetMSPIDs(channel)

		p := cauthdsl.SignedByAnyAdmin(mspids)
		ip, err = protoutil.Marshal(p)
		if err != nil {
			return nil, errors.Errorf("error marshalling default instantiation policy")
		}

	}
	return ip, nil
}



func (s *supportImpl) CheckInstantiationPolicy(signedProp *pb.SignedProposal, chainName string, instantiationPolicy []byte) error {
	
	mgr := mgmt.GetManagerForChain(chainName)
	if mgr == nil {
		return errors.Errorf("error checking chaincode instantiation policy: MSP manager for channel %s not found", chainName)
	}
	npp := cauthdsl.NewPolicyProvider(mgr)
	instPol, _, err := npp.NewPolicy(instantiationPolicy)
	if err != nil {
		return err
	}
	proposal, err := protoutil.UnmarshalProposal(signedProp.ProposalBytes)
	if err != nil {
		return err
	}
	
	header, err := protoutil.UnmarshalHeader(proposal.Header)
	if err != nil {
		return err
	}
	shdr, err := protoutil.UnmarshalSignatureHeader(header.SignatureHeader)
	if err != nil {
		return err
	}
	
	sd := []*protoutil.SignedData{{
		Data:      signedProp.ProposalBytes,
		Identity:  shdr.Creator,
		Signature: signedProp.Signature,
	}}
	err = instPol.Evaluate(sd)
	if err != nil {
		return errors.WithMessage(err, "instantiation policy violation")
	}
	return nil
}
