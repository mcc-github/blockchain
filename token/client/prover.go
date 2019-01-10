/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"io"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/protos/token"
	tk "github.com/mcc-github/blockchain/token"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type TimeFunc func() time.Time




type ProverPeerClient interface {
	
	CreateProverClient() (*grpc.ClientConn, token.ProverClient, error)

	
	Certificate() *tls.Certificate
}


type ProverPeerClientImpl struct {
	Address            string
	ServerNameOverride string
	GRPCClient         *comm.GRPCClient
}


type ProverPeer struct {
	ChannelID        string
	ProverPeerClient ProverPeerClient
	RandomnessReader io.Reader
	Time             TimeFunc
}

func NewProverPeer(config *ClientConfig) (*ProverPeer, error) {
	
	grpcClient, err := CreateGRPCClient(&config.ProverPeer)
	if err != nil {
		return nil, err
	}

	return &ProverPeer{
		ChannelID:        config.ChannelID,
		RandomnessReader: rand.Reader,
		Time:             time.Now,
		ProverPeerClient: &ProverPeerClientImpl{
			Address:            config.ProverPeer.Address,
			ServerNameOverride: config.ProverPeer.ServerNameOverride,
			GRPCClient:         grpcClient,
		},
	}, nil
}

func (pc *ProverPeerClientImpl) CreateProverClient() (*grpc.ClientConn, token.ProverClient, error) {
	conn, err := pc.GRPCClient.NewConnection(pc.Address, pc.ServerNameOverride)
	if err != nil {
		return conn, nil, err
	}
	return conn, token.NewProverClient(conn), nil
}

func (pc *ProverPeerClientImpl) Certificate() *tls.Certificate {
	cert := pc.GRPCClient.Certificate()
	return &cert
}




func (prover *ProverPeer) RequestImport(tokensToIssue []*token.TokenToIssue, signingIdentity tk.SigningIdentity) ([]byte, error) {
	ir := &token.ImportRequest{
		TokensToIssue: tokensToIssue,
	}
	payload := &token.Command_ImportRequest{ImportRequest: ir}

	sc, err := prover.CreateSignedCommand(payload, signingIdentity)
	if err != nil {
		return nil, err
	}

	return prover.SendCommand(context.Background(), sc)
}






func (prover *ProverPeer) RequestTransfer(
	tokenIDs [][]byte,
	shares []*token.RecipientTransferShare,
	signingIdentity tk.SigningIdentity) ([]byte, error) {

	tr := &token.TransferRequest{
		Shares:   shares,
		TokenIds: tokenIDs,
	}
	payload := &token.Command_TransferRequest{TransferRequest: tr}

	sc, err := prover.CreateSignedCommand(payload, signingIdentity)
	if err != nil {
		return nil, err
	}

	return prover.SendCommand(context.Background(), sc)
}


func (prover *ProverPeer) SendCommand(ctx context.Context, sc *token.SignedCommand) ([]byte, error) {
	conn, proverClient, err := prover.ProverPeerClient.CreateProverClient()
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return nil, err
	}
	scr, err := proverClient.ProcessCommand(ctx, sc)
	if err != nil {
		return nil, err
	}

	commandResp := &token.CommandResponse{}
	err = proto.Unmarshal(scr.Response, commandResp)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal command response")
	}
	if commandResp.GetErr() != nil {
		return nil, errors.Errorf("error from prover: %s", commandResp.GetErr().GetMessage())
	}

	txBytes, err := proto.Marshal(commandResp.GetTokenTransaction())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal TokenTransaction")
	}
	return txBytes, nil
}

func (prover *ProverPeer) CreateSignedCommand(payload interface{}, signingIdentity tk.SigningIdentity) (*token.SignedCommand, error) {

	command, err := commandFromPayload(payload)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 32)
	_, err = io.ReadFull(prover.RandomnessReader, nonce)
	if err != nil {
		return nil, err
	}

	ts, err := ptypes.TimestampProto(prover.Time())
	if err != nil {
		return nil, err
	}

	creator, err := signingIdentity.Serialize()
	if err != nil {
		return nil, err
	}

	
	tlsCertHash, err := GetTLSCertHash(prover.ProverPeerClient.Certificate())
	if err != nil {
		return nil, err
	}
	command.Header = &token.Header{
		Timestamp:   ts,
		Nonce:       nonce,
		Creator:     creator,
		ChannelId:   prover.ChannelID,
		TlsCertHash: tlsCertHash,
	}

	raw, err := proto.Marshal(command)
	if err != nil {
		return nil, err
	}

	signature, err := signingIdentity.Sign(raw)
	if err != nil {
		return nil, err
	}

	sc := &token.SignedCommand{
		Command:   raw,
		Signature: signature,
	}
	return sc, nil
}

func commandFromPayload(payload interface{}) (*token.Command, error) {
	switch t := payload.(type) {
	case *token.Command_ImportRequest:
		return &token.Command{Payload: t}, nil
	case *token.Command_TransferRequest:
		return &token.Command{Payload: t}, nil
	default:
		return nil, errors.Errorf("command type not recognized: %T", t)
	}
}
