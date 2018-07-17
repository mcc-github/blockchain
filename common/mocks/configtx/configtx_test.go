/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	"testing"

	"github.com/mcc-github/blockchain/common/configtx"
)

func TestConfigtxValidatorInterface(t *testing.T) {
	_ = configtx.Validator(&Validator{})
}
