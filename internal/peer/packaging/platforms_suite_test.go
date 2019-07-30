/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package packaging_test

import (
	"testing"

	"github.com/mcc-github/blockchain/internal/peer/packaging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type platform interface {
	packaging.Platform
}

func TestPackaging(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Platforms Suite")
}
