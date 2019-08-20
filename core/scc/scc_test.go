/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc_test

import (
	"testing"

	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/scc/mock"
	"github.com/onsi/gomega"
)


type chaincodeStreamHandler interface {
	scc.ChaincodeStreamHandler
}

func TestDeploy(t *testing.T) {
	gt := gomega.NewGomegaWithT(t)

	csh := &mock.ChaincodeStreamHandler{}
	doneC := make(chan struct{})
	close(doneC)
	csh.LaunchInProcReturns(doneC)
	scc.DeploySysCC(&lifecycle.SCC{}, csh)
	gt.Expect(csh.LaunchInProcCallCount()).To(gomega.Equal(1))
	gt.Expect(csh.LaunchInProcArgsForCall(0)).To(gomega.Equal("_lifecycle.syscc"))
	gt.Eventually(csh.HandleChaincodeStreamCallCount).Should(gomega.Equal(1))
}
