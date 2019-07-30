/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import validation "github.com/mcc-github/blockchain/core/handlers/validation/api"



type Capabilities interface {
	validation.Dependency
	
	Supported() error

	
	
	ForbidDuplicateTXIdInBlock() bool

	
	ACLs() bool

	
	PrivateChannelData() bool

	
	
	CollectionUpgrade() bool

	
	
	V1_1Validation() bool

	
	
	V1_2Validation() bool

	
	
	
	
	V1_3Validation() bool

	
	
	StorePvtDataOfInvalidTx() bool

	
	
	
	
	V2_0Validation() bool

	
	MetadataLifecycle() bool

	
	
	KeyLevelEndorsement() bool

	
	FabToken() bool
}
