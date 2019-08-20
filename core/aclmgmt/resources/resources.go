/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/





package resources

const (
	
	Lifecycle_InstallChaincode                   = "_lifecycle/InstallChaincode"
	Lifecycle_QueryInstalledChaincode            = "_lifecycle/QueryInstalledChaincode"
	Lifecycle_GetInstalledChaincodePackage       = "_lifecycle/GetInstalledChaincodePackage"
	Lifecycle_QueryInstalledChaincodes           = "_lifecycle/QueryInstalledChaincodes"
	Lifecycle_ApproveChaincodeDefinitionForMyOrg = "_lifecycle/ApproveChaincodeDefinitionForMyOrg"
	Lifecycle_CommitChaincodeDefinition          = "_lifecycle/CommitChaincodeDefinition"
	Lifecycle_QueryChaincodeDefinition           = "_lifecycle/QueryChaincodeDefinition"
	Lifecycle_QueryChaincodeDefinitions          = "_lifecycle/QueryChaincodeDefinitions"
	Lifecycle_CheckCommitReadiness               = "_lifecycle/CheckCommitReadiness"

	
	Lscc_Install                   = "lscc/Install"
	Lscc_Deploy                    = "lscc/Deploy"
	Lscc_Upgrade                   = "lscc/Upgrade"
	Lscc_ChaincodeExists           = "lscc/ChaincodeExists"
	Lscc_GetDeploymentSpec         = "lscc/GetDeploymentSpec"
	Lscc_GetChaincodeData          = "lscc/GetChaincodeData"
	Lscc_GetInstantiatedChaincodes = "lscc/GetInstantiatedChaincodes"
	Lscc_GetInstalledChaincodes    = "lscc/GetInstalledChaincodes"
	Lscc_GetCollectionsConfig      = "lscc/GetCollectionsConfig"

	
	Qscc_GetChainInfo       = "qscc/GetChainInfo"
	Qscc_GetBlockByNumber   = "qscc/GetBlockByNumber"
	Qscc_GetBlockByHash     = "qscc/GetBlockByHash"
	Qscc_GetTransactionByID = "qscc/GetTransactionByID"
	Qscc_GetBlockByTxID     = "qscc/GetBlockByTxID"

	
	Cscc_JoinChain      = "cscc/JoinChain"
	Cscc_GetConfigBlock = "cscc/GetConfigBlock"
	Cscc_GetChannels    = "cscc/GetChannels"

	
	Peer_Propose              = "peer/Propose"
	Peer_ChaincodeToChaincode = "peer/ChaincodeToChaincode"

	
	Event_Block         = "event/Block"
	Event_FilteredBlock = "event/FilteredBlock"
)
