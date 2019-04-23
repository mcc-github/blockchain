/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mspext_test

import (
	"github.com/mcc-github/blockchain/common/tools/protolator"
	"github.com/mcc-github/blockchain/common/tools/protolator/protoext/mspext"
)


var (
	_ protolator.VariablyOpaqueFieldProto = &mspext.MSPConfig{}
	_ protolator.DecoratedProto           = &mspext.MSPConfig{}

	_ protolator.VariablyOpaqueFieldProto = &mspext.MSPPrincipal{}
	_ protolator.DecoratedProto           = &mspext.MSPPrincipal{}
)
