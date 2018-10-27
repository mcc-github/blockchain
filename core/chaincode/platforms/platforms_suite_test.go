/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package platforms_test

import (
	"testing"

	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type platform interface {
	platforms.Platform
}


type packageWriter interface {
	platforms.PackageWriter
}

func TestPlatforms(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Platforms Suite")
}
