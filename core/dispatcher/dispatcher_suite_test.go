/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package dispatcher_test

import (
	"testing"

	"github.com/mcc-github/blockchain/core/dispatcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type protobuf interface {
	dispatcher.Protobuf
}

func TestDispatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dispatcher Suite")
}
