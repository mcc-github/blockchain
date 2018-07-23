/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package resources





const (
	
	Lscc_Install                   = "lscc/Install"
	Lscc_Deploy                    = "lscc/Deploy"
	Lscc_Upgrade                   = "lscc/Upgrade"
	Lscc_ChaincodeExists           = "lscc/ChaincodeExists"
	Lscc_GetDeploymentSpec         = "lscc/GetDeploymentSpec"
	Lscc_GetChaincodeData          = "lscc/GetChaincodeData"
	Lscc_GetInstantiatedChaincodes = "lscc/GetInstantiatedChaincodes"
	Lscc_GetInstalledChaincodes    = "lscc/GetInstalledChaincodes"

	
	Qscc_GetChainInfo       = "qscc/GetChainInfo"
	Qscc_GetBlockByNumber   = "qscc/GetBlockByNumber"
	Qscc_GetBlockByHash     = "qscc/GetBlockByHash"
	Qscc_GetTransactionByID = "qscc/GetTransactionByID"
	Qscc_GetBlockByTxID     = "qscc/GetBlockByTxID"

	
	Cscc_JoinChain                = "cscc/JoinChain"
	Cscc_GetConfigBlock           = "cscc/GetConfigBlock"
	Cscc_GetChannels              = "cscc/GetChannels"
	Cscc_GetConfigTree            = "cscc/GetConfigTree"
	Cscc_SimulateConfigTreeUpdate = "cscc/SimulateConfigTreeUpdate"

	
	Peer_Propose              = "peer/Propose"
	Peer_ChaincodeToChaincode = "peer/ChaincodeToChaincode"

	
	Event_Block         = "event/Block"
	Event_FilteredBlock = "event/FilteredBlock"
)