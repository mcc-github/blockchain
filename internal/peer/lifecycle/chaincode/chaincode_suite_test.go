/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/internal/peer/common"
	"github.com/mcc-github/blockchain/internal/peer/lifecycle/chaincode"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	pb "github.com/mcc-github/blockchain/protos/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type writer interface {
	chaincode.Writer
}


type platformRegistryIntf interface {
	chaincode.PlatformRegistry
}


type reader interface {
	chaincode.Reader
}


type endorserClient interface {
	chaincode.EndorserClient
}


type signer interface {
	chaincode.Signer
}


type broadcastClient interface {
	common.BroadcastClient
}


type peerDeliverClient interface {
	pb.DeliverClient
}


type deliver interface {
	pb.Deliver_DeliverClient
}

func TestChaincode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chaincode Suite")
}


func TestMain(m *testing.M) {
	err := msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		panic(fmt.Sprintf("Fatal error when reading MSP config: %s", err))
	}
	os.Exit(m.Run())
}
