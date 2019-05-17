/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package runner

import (
	"time"

	"github.com/mcc-github/blockchain/integration/helpers"
)

const DefaultStartTimeout = 45 * time.Second


var DefaultNamer NameFunc = helpers.UniqueName


type NameFunc func() string
