/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package dockercontroller_test

import (
	"github.com/mcc-github/blockchain/core/container/dockercontroller"
)


type platformBuilder interface {
	dockercontroller.PlatformBuilder
}
