/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliver_test

import (
	"testing"

	"github.com/mcc-github/blockchain/common/deliver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)



type filteredResponseSender interface {
	deliver.ResponseSender
	deliver.Filtered
}



type privateDataResponseSender interface {
	deliver.ResponseSender
}

func TestDeliver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deliver Suite")
}
