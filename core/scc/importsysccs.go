/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	pb "github.com/mcc-github/blockchain/protos/peer"
)



type ChaincodeStreamHandler interface {
	HandleChaincodeStream(ccintf.ChaincodeStream) error
	LaunchInProc(packageID ccintf.CCID) <-chan struct{}
}



func (p *Provider) DeploySysCCs(chaincodeStreamHandler ChaincodeStreamHandler) {
	for _, sysCC := range p.SysCCs {
		if !sysCC.Enabled() || !p.isWhitelisted(sysCC) {
			sysccLogger.Infof("System chaincode '%s' is disabled", sysCC.Name())
			continue
		}
		sysccLogger.Infof("deploying system chaincode '%s'", sysCC.Name())

		
		version := util.GetSysCCVersion()
		ccid := ccintf.CCID(sysCC.Name() + ":" + version)

		done := chaincodeStreamHandler.LaunchInProc(ccid)

		peerRcvCCSend := make(chan *pb.ChaincodeMessage)
		ccRcvPeerSend := make(chan *pb.ChaincodeMessage)

		
		go func() {
			sysccLogger.Debugf("starting chaincode-support stream for  %s", ccid)
			err := chaincodeStreamHandler.HandleChaincodeStream(newInProcStream(peerRcvCCSend, ccRcvPeerSend))
			sysccLogger.Criticalf("shim stream ended with err: %v", err)
		}()

		go func(sysCC SelfDescribingSysCC) {
			sysccLogger.Debugf("chaincode started for %s", ccid)
			err := shim.StartInProc(ccid.String(), newInProcStream(ccRcvPeerSend, peerRcvCCSend), sysCC.Chaincode())
			sysccLogger.Criticalf("system chaincode ended with err: %v", err)
		}(sysCC)
		<-done
	}
}
