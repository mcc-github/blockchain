/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package accesscontrol

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/flogging"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"google.golang.org/grpc"
)

var logger = flogging.MustGetLogger("accessControl")



type CertAndPrivKeyPair struct {
	
	Cert string
	
	Key string
}

type Authenticator struct {
	mapper *certMapper
}

func (auth *Authenticator) Wrap(srv pb.ChaincodeSupportServer) pb.ChaincodeSupportServer {
	return newInterceptor(srv, auth.authenticate)
}


func NewAuthenticator(ca tlsgen.CA) *Authenticator {
	return &Authenticator{
		mapper: newCertMapper(ca.NewClientCertKeyPair),
	}
}




func (ac *Authenticator) Generate(ccName string) (*CertAndPrivKeyPair, error) {
	cert, err := ac.mapper.genCert(ccName)
	if err != nil {
		return nil, err
	}
	return &CertAndPrivKeyPair{
		Key:  cert.PrivKeyString(),
		Cert: cert.PubKeyString(),
	}, nil
}

func (ac *Authenticator) authenticate(msg *pb.ChaincodeMessage, stream grpc.ServerStream) error {
	if msg.Type != pb.ChaincodeMessage_REGISTER {
		logger.Warning("Got message", msg, "but expected a ChaincodeMessage_REGISTER message")
		return errors.New("First message needs to be a register")
	}

	chaincodeID := &pb.ChaincodeID{}
	err := proto.Unmarshal(msg.Payload, chaincodeID)
	if err != nil {
		logger.Warning("Failed unmarshaling message:", err)
		return err
	}
	ccName := chaincodeID.Name
	
	hash := extractCertificateHashFromContext(stream.Context())
	if len(hash) == 0 {
		errMsg := fmt.Sprintf("TLS is active but chaincode %s didn't send certificate", ccName)
		logger.Warning(errMsg)
		return errors.New(errMsg)
	}
	
	registeredName := ac.mapper.lookup(certHash(hash))
	if registeredName == "" {
		errMsg := fmt.Sprintf("Chaincode %s with given certificate hash %v not found in registry", ccName, hash)
		logger.Warning(errMsg)
		return errors.New(errMsg)
	}
	if registeredName != ccName {
		errMsg := fmt.Sprintf("Chaincode %s with given certificate hash %v belongs to a different chaincode", ccName, hash)
		logger.Warning(errMsg)
		return fmt.Errorf(errMsg)
	}

	logger.Debug("Chaincode", ccName, "'s authentication is authorized")
	return nil
}
