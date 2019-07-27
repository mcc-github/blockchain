/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valinforetriever

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)




type LifecycleResources interface {
	plugindispatcher.LifecycleResources
}






type ValidationInfoRetrieveShim struct {
	Legacy plugindispatcher.LifecycleResources
	New    plugindispatcher.LifecycleResources
}









func (v *ValidationInfoRetrieveShim) ValidationInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (plugin string, args []byte, unexpectedErr error, validationErr error) {
	plugin, args, unexpectedErr, validationErr = v.New.ValidationInfo(channelID, chaincodeName, qe)
	if unexpectedErr != nil || validationErr != nil || plugin != "" || len(args) != 0 {
		
		
		
		return
	}

	plugin, args, unexpectedErr, validationErr = v.Legacy.ValidationInfo(channelID, chaincodeName, qe)
	if unexpectedErr != nil || validationErr != nil || plugin != "vscc" || len(args) == 0 {
		
		
		
		return
	}

	spe := &common.SignaturePolicyEnvelope{}
	err := proto.Unmarshal(args, spe)
	if err != nil {
		
		
		return
	}

	newArgs, err := proto.Marshal(
		&peer.ApplicationPolicy{
			Type: &peer.ApplicationPolicy_SignaturePolicy{
				SignaturePolicy: spe,
			},
		},
	)
	if err != nil {
		
		
		return
	}

	
	args = newArgs
	return
}
