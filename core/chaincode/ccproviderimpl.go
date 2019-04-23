/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type CCProviderImpl struct {
	cs *ChaincodeSupport
}

func NewProvider(cs *ChaincodeSupport) *CCProviderImpl {
	return &CCProviderImpl{cs: cs}
}


func (c *CCProviderImpl) ExecuteLegacyInit(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, spec *pb.ChaincodeDeploymentSpec) (*pb.Response, *pb.ChaincodeEvent, error) {
	return c.cs.ExecuteLegacyInit(txParams, cccid, spec)
}


func (c *CCProviderImpl) Stop(ccci *ccprovider.ChaincodeContainerInfo) error {
	return c.cs.Stop(ccci)
}
