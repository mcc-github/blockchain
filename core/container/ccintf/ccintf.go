/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccintf





import (
	"fmt"

	pb "github.com/mcc-github/blockchain/protos/peer"
	"golang.org/x/net/context"
)


type ChaincodeStream interface {
	Send(*pb.ChaincodeMessage) error
	Recv() (*pb.ChaincodeMessage, error)
}



type CCSupport interface {
	HandleChaincodeStream(context.Context, ChaincodeStream) error
}


func GetCCHandlerKey() string {
	return "CCHANDLER"
}


type CCID struct {
	Name    string
	Version string
}


func (ccid *CCID) GetName() string {
	if ccid.Version != "" {
		return fmt.Sprintf("%s-%s", ccid.Name, ccid.Version)
	}
	return ccid.Name
}
