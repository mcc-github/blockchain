/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccintf





import (
	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type ChaincodeStream interface {
	Send(*pb.ChaincodeMessage) error
	Recv() (*pb.ChaincodeMessage, error)
}


type CCID string


func (c CCID) String() string {
	return string(c)
}


func New(packageID persistence.PackageID) CCID {
	return CCID(packageID.String())
}


type PeerConnection struct {
	Address   string
	TLSConfig *TLSConfig
}


type TLSConfig struct {
	ClientCert []byte
	ClientKey  []byte
	RootCert   []byte
}
