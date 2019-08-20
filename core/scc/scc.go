/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	pb "github.com/mcc-github/blockchain/protos/peer"
)






const SysCCVersion = "syscc"


func ChaincodeID(ccName string) string {
	return ccName + "." + SysCCVersion
}


var sysccLogger = flogging.MustGetLogger("sccapi")







type BuiltinSCCs map[string]struct{}

func (bccs BuiltinSCCs) IsSysCC(name string) bool {
	_, ok := bccs[name]
	return ok
}



type ChaincodeStreamHandler interface {
	HandleChaincodeStream(ccintf.ChaincodeStream) error
	LaunchInProc(packageID string) <-chan struct{}
}

type SelfDescribingSysCC interface {
	
	Name() string

	
	Chaincode() shim.Chaincode
}



func DeploySysCC(sysCC SelfDescribingSysCC, chaincodeStreamHandler ChaincodeStreamHandler) {
	sysccLogger.Infof("deploying system chaincode '%s'", sysCC.Name())

	ccid := ChaincodeID(sysCC.Name())

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
		err := shim.StartInProc(ccid, newInProcStream(ccRcvPeerSend, peerRcvCCSend), sysCC.Chaincode())
		sysccLogger.Criticalf("system chaincode ended with err: %v", err)
	}(sysCC)
	<-done
}
