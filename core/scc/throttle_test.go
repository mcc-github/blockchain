/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc_test

import (
	"testing"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/scc/mock"
	pb "github.com/mcc-github/blockchain/protos/peer"
	. "github.com/onsi/gomega"
)




type chaincode interface{ shim.Chaincode }
type selfDescribingSysCC interface{ scc.SelfDescribingSysCC } 

func TestThrottle(t *testing.T) {
	gt := NewGomegaWithT(t)

	runningCh := make(chan struct{}, 1)
	doneCh := make(chan struct{})
	chaincode := &mock.Chaincode{}
	chaincode.InvokeStub = func(shim.ChaincodeStubInterface) pb.Response {
		gt.Eventually(runningCh).Should(BeSent(struct{}{}))
		<-doneCh
		return pb.Response{}
	}

	sysCC := &mock.SelfDescribingSysCC{}
	sysCC.ChaincodeReturns(chaincode)

	
	throttled := scc.Throttle(5, sysCC).Chaincode()
	for i := 0; i < 5; i++ {
		go throttled.Invoke(nil)
		gt.Eventually(runningCh).Should(Receive())
	}
	
	go throttled.Invoke(nil)
	gt.Consistently(runningCh).ShouldNot(Receive())

	
	
	gt.Eventually(doneCh).Should(BeSent(struct{}{}))
	gt.Eventually(runningCh).Should(Receive())

	
	close(doneCh)
}

func TestThrottledChaincode(t *testing.T) {
	gt := NewGomegaWithT(t)

	chaincode := &mock.Chaincode{}
	chaincode.InitReturns(pb.Response{Message: "init-returns"})
	chaincode.InvokeReturns(pb.Response{Message: "invoke-returns"})

	sysCC := &mock.SelfDescribingSysCC{}
	sysCC.ChaincodeReturns(chaincode)
	throttled := scc.Throttle(5, sysCC).Chaincode()

	initResp := throttled.Init(nil)
	gt.Expect(initResp.Message).To(Equal("init-returns"))
	invokeResp := throttled.Invoke(nil)
	gt.Expect(invokeResp.Message).To(Equal("invoke-returns"))
}
