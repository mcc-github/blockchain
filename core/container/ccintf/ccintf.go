/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccintf

import (
	pb "github.com/mcc-github/blockchain-protos-go/peer"
)




type ChaincodeStream interface {
	Send(*pb.ChaincodeMessage) error
	Recv() (*pb.ChaincodeMessage, error)
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
