/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package valinforetriever

import (
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher"
	"github.com/mcc-github/blockchain/core/ledger"
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
	if unexpectedErr != nil || validationErr != nil || plugin != "" || args != nil {
		return
	}

	return v.Legacy.ValidationInfo(channelID, chaincodeName, qe)
}
