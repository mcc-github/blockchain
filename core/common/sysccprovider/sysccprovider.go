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

package sysccprovider

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/ledger"
)





type SystemChaincodeProvider interface {
	
	
	
	
	GetQueryExecutorForLedger(cid string) (ledger.QueryExecutor, error)

	
	
	GetApplicationConfig(cid string) (channelconfig.Application, bool)

	
	
	PolicyManager(channelID string) (policies.Manager, bool)
}


type ChaincodeInstance struct {
	ChannelID        string
	ChaincodeName    string
	ChaincodeVersion string
}

func (ci *ChaincodeInstance) String() string {
	return ci.ChannelID + "." + ci.ChaincodeName + "#" + ci.ChaincodeVersion
}
