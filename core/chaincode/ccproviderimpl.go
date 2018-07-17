/*
Copyright IBM Corp. 2017 All Rights Reserved.

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

package chaincode

import (
	"context"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type CCProviderImpl struct {
	cs *ChaincodeSupport
}

func NewProvider(cs *ChaincodeSupport) *CCProviderImpl {
	return &CCProviderImpl{cs: cs}
}


func (c *CCProviderImpl) GetContext(ledger ledger.PeerLedger, txid string) (context.Context, ledger.TxSimulator, error) {
	
	txsim, err := ledger.NewTxSimulator(txid)
	if err != nil {
		return nil, nil, err
	}
	ctxt := context.WithValue(context.Background(), TXSimulatorKey, txsim)
	return ctxt, txsim, nil
}


func (c *CCProviderImpl) ExecuteChaincode(ctxt context.Context, cccid *ccprovider.CCContext, args [][]byte) (*pb.Response, *pb.ChaincodeEvent, error) {
	invocationSpec := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type:        pb.ChaincodeSpec_GOLANG,
			ChaincodeId: &pb.ChaincodeID{Name: cccid.Name},
			Input:       &pb.ChaincodeInput{Args: args},
		},
	}
	return c.cs.Execute(ctxt, cccid, invocationSpec)
}


func (c *CCProviderImpl) Execute(ctxt context.Context, cccid *ccprovider.CCContext, spec ccprovider.ChaincodeSpecGetter) (*pb.Response, *pb.ChaincodeEvent, error) {
	return c.cs.Execute(ctxt, cccid, spec)
}


func (c *CCProviderImpl) Stop(ctxt context.Context, cccid *ccprovider.CCContext, spec *pb.ChaincodeDeploymentSpec) error {
	return c.cs.Stop(ctxt, cccid, spec)
}
