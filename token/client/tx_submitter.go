/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	peercommon "github.com/mcc-github/blockchain/peer/common"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("token/client")

type TxSubmitter struct {
	Config        *ClientConfig
	Signer        SignerIdentity
	Creator       []byte
	OrdererClient OrdererClient
	DeliverClient DeliverClient
}







type TxEvent struct {
	Txid       string
	Committed  bool
	CommitPeer string
	Err        error
}


func NewTxSubmitter(config *ClientConfig) (*TxSubmitter, error) {
	err := ValidateClientConfig(config)
	if err != nil {
		return nil, err
	}

	
	mspType := "bccsp"
	peercommon.InitCrypto(config.MspDir, config.MspId, mspType)

	Signer, err := mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		return nil, err
	}
	creator, err := Signer.Serialize()
	if err != nil {
		return nil, err
	}

	ordererClient, err := NewOrdererClient(config)
	if err != nil {
		return nil, err
	}

	deliverClient, err := NewDeliverClient(config)
	if err != nil {
		return nil, err
	}

	return &TxSubmitter{
		Config:        config,
		Signer:        Signer,
		Creator:       creator,
		OrdererClient: ordererClient,
		DeliverClient: deliverClient,
	}, nil
}






func (s *TxSubmitter) SubmitTransaction(txEnvelope *common.Envelope, waitTimeInSeconds int) (committed bool, txId string, err error) {
	if waitTimeInSeconds > 0 {
		waitTime := time.Second * time.Duration(waitTimeInSeconds)
		ctx, cancelFunc := context.WithTimeout(context.Background(), waitTime)
		defer cancelFunc()
		localCh := make(chan TxEvent, 1)
		committed, txId, err = s.sendTransactionInternal(txEnvelope, ctx, localCh, true)
		close(localCh)
		return
	} else {
		committed, txId, err = s.sendTransactionInternal(txEnvelope, context.Background(), nil, false)
		return
	}
}






func (s *TxSubmitter) SubmitTransactionWithChan(txEnvelope *common.Envelope, eventCh chan TxEvent) (committed bool, txId string, err error) {
	committed, txId, err = s.sendTransactionInternal(txEnvelope, context.Background(), eventCh, false)
	return
}

func (s *TxSubmitter) sendTransactionInternal(txEnvelope *common.Envelope, ctx context.Context, eventCh chan TxEvent, waitForCommit bool) (bool, string, error) {
	if eventCh != nil && cap(eventCh) == 0 {
		return false, "", errors.New("eventCh buffer size must be greater than 0")
	}
	if eventCh != nil && len(eventCh) == cap(eventCh) {
		return false, "", errors.New("eventCh buffer is full. Read events and try again")
	}

	txid, err := getTransactionId(txEnvelope)
	if err != nil {
		return false, "", err
	}

	broadcast, err := s.OrdererClient.NewBroadcast(context.Background())
	if err != nil {
		return false, "", err
	}

	committed := false
	if eventCh != nil {
		deliverFiltered, err := s.DeliverClient.NewDeliverFiltered(ctx)
		if err != nil {
			return false, "", err
		}
		blockEnvelope, err := CreateDeliverEnvelope(s.Config.ChannelId, s.Creator, s.Signer, s.DeliverClient.Certificate())
		if err != nil {
			return false, "", err
		}
		err = DeliverSend(deliverFiltered, s.Config.CommitPeerCfg.Address, blockEnvelope)
		if err != nil {
			return false, "", err
		}
		go DeliverReceive(deliverFiltered, s.Config.CommitPeerCfg.Address, txid, eventCh)
	}

	err = BroadcastSend(broadcast, s.Config.OrdererCfg.Address, txEnvelope)
	if err != nil {
		return false, txid, err
	}

	
	responses := make(chan common.Status)
	errs := make(chan error, 1)
	go BroadcastReceive(broadcast, s.Config.OrdererCfg.Address, responses, errs)
	_, err = BroadcastWaitForResponse(responses, errs)

	
	if eventCh != nil && waitForCommit {
		committed, err = DeliverWaitForResponse(ctx, eventCh, txid)
	}

	return committed, txid, err
}

func (s *TxSubmitter) CreateTxEnvelope(txBytes []byte) (string, *common.Envelope, error) {
	
	
	var tlsCertHash []byte
	var err error
	
	cert := s.OrdererClient.Certificate()
	if cert != nil && len(cert.Certificate) > 0 {
		tlsCertHash, err = factory.GetDefault().Hash(cert.Certificate[0], &bccsp.SHA256Opts{})
		if err != nil {
			err = errors.New("failed to compute SHA256 on client certificate")
			logger.Errorf("%s", err)
			return "", nil, err
		}
	}

	txid, header, err := CreateHeader(common.HeaderType_TOKEN_TRANSACTION, s.Config.ChannelId, s.Creator, tlsCertHash)
	if err != nil {
		return txid, nil, err
	}

	txEnvelope, err := CreateEnvelope(txBytes, header, s.Signer)
	return txid, txEnvelope, err
}



func CreateHeader(txType common.HeaderType, channelId string, creator []byte, tlsCertHash []byte) (string, *common.Header, error) {
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return "", nil, err
	}

	nonce, err := crypto.GetRandomNonce()
	if err != nil {
		return "", nil, err
	}

	txId, err := utils.ComputeTxID(nonce, creator)
	if err != nil {
		return "", nil, err
	}

	chdr := &common.ChannelHeader{
		Type:        int32(txType),
		ChannelId:   channelId,
		TxId:        txId,
		Epoch:       uint64(0),
		Timestamp:   ts,
		TlsCertHash: tlsCertHash,
	}
	chdrBytes, err := proto.Marshal(chdr)
	if err != nil {
		return "", nil, err
	}

	shdr := &common.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}
	shdrBytes, err := proto.Marshal(shdr)
	if err != nil {
		return "", nil, err
	}

	header := &common.Header{
		ChannelHeader:   chdrBytes,
		SignatureHeader: shdrBytes,
	}

	return txId, header, nil
}


func CreateEnvelope(data []byte, header *common.Header, signer SignerIdentity) (*common.Envelope, error) {
	payload := &common.Payload{
		Header: header,
		Data:   data,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal common.Payload")
	}

	signature, err := signer.Sign(payloadBytes)
	if err != nil {
		return nil, err
	}

	txEnvelope := &common.Envelope{
		Payload:   payloadBytes,
		Signature: signature,
	}

	return txEnvelope, nil
}

func getTransactionId(txEnvelope *common.Envelope) (string, error) {
	payload := common.Payload{}
	err := proto.Unmarshal(txEnvelope.Payload, &payload)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal envelope payload")
	}

	channelHeader := common.ChannelHeader{}
	err = proto.Unmarshal(payload.Header.ChannelHeader, &channelHeader)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal channel header")
	}

	return channelHeader.TxId, nil
}


func createGrpcClient(cfg *ConnectionConfig, tlsEnabled bool) (*comm.GRPCClient, error) {
	clientConfig := comm.ClientConfig{Timeout: time.Second}

	if tlsEnabled {
		if cfg.TlsRootCertFile == "" {
			return nil, errors.New("missing TlsRootCertFile in client config")
		}
		caPEM, err := ioutil.ReadFile(cfg.TlsRootCertFile)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("unable to load TLS cert from %s", cfg.TlsRootCertFile))
		}
		secOpts := &comm.SecureOptions{
			UseTLS:            true,
			ServerRootCAs:     [][]byte{caPEM},
			RequireClientCert: false,
		}
		clientConfig.SecOpts = secOpts
	}

	return comm.NewGRPCClient(clientConfig)
}
