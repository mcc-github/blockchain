/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser_test

import (
	"testing"

	"github.com/mcc-github/blockchain/core/endorser"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/msp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type endorserSupport interface {
	endorser.Support
}


type identityDeserializer interface {
	msp.IdentityDeserializer
}


type identity interface {
	msp.Identity
}


type txSimulator interface {
	ledger.TxSimulator
}

func TestEndorser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Endorser Suite")
}
