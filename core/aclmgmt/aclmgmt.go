/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package aclmgmt

import (
	"github.com/mcc-github/blockchain/common/flogging"
)

var aclLogger = flogging.MustGetLogger("aclmgmt")

type ACLProvider interface {
	
	
	
	CheckACL(resName string, channelID string, idinfo interface{}) error
}
