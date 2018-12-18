/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package client

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	peercommon "github.com/mcc-github/blockchain/peer/common"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	tk "github.com/mcc-github/blockchain/token"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("token.client")



type TxSubmitter struct {
	Config          *ClientConfig
	SigningIdentity tk.SigningIdentity
	Creator         []byte
	OrdererClient   OrdererClient
	DeliverClient   DeliverClient
}


type TxEvent struct {
	Txid       string
	Committed  bool
	CommitPeer string
	Err        error
}


func NewTxSubmitter(config ClientConfig) (*TxSubmitter, error) {
	err := ValidateClientConfig(config)
	if err != nil {
		return nil, err
	}

	if config.MSPInfo.MSPType == "" {
		config.MSPInfo.MSPType = "bccsp"
	}
	peercommon.InitCrypto(config.MSPInfo.MSPConfigPath, config.MSPInfo.MSPID, config.MSPInfo.MSPType)

	signingIdentity, err := mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		return nil, err
	}
	creator, err := signingIdentity.Serialize()
	if err != nil {
		return nil, err
	}

	ordererClient, err := NewOrdererClient(&config.Orderer)
	if err != nil {
		return nil, err
	}

	deliverClient, err := NewDeliverClient(&config.CommitterPeer)
	if err != nil {
		return nil, err
	}

	return &TxSubmitter{
		Config:          &config,
		SigningIdentity: signingIdentity,
		Creator:         creator,
		OrdererClient:   ordererClient,
		DeliverClient:   deliverClient,
	}, nil
}







func (s *TxSubmitter) Submit(txEnvelope *common.Envelope, waitTimeout time.Duration) (*common.Status, bool, error) {
	if txEnvelope == nil {
		return nil, false, errors.New("envelope is nil")
	}

	txid, err := getTransactionId(txEnvelope)
	if err != nil {
		return nil, false, err
	}

	broadcast, err := s.OrdererClient.NewBroadcast(context.Background())
	if err != nil {
		return nil, false, err
	}

	var eventCh chan TxEvent
	var ctx context.Context
	var cancelFunc context.CancelFunc
	if waitTimeout > 0 {
		ctx, cancelFunc = context.WithTimeout(context.Background(), waitTimeout)
		defer cancelFunc()
		deliverFiltered, err := s.DeliverClient.NewDeliverFiltered(ctx)
		if err != nil {
			return nil, false, err
		}
		blockEnvelope, err := CreateDeliverEnvelope(s.Config.ChannelID, s.Creator, s.SigningIdentity, s.DeliverClient.Certificate())
		if err != nil {
			return nil, false, err
		}
		err = DeliverSend(deliverFiltered, s.Config.CommitterPeer.Address, blockEnvelope)
		if err != nil {
			return nil, false, err
		}
		eventCh = make(chan TxEvent, 1)
		go DeliverReceive(deliverFiltered, s.Config.CommitterPeer.Address, txid, eventCh)
	}

	err = BroadcastSend(broadcast, s.Config.Orderer.Address, txEnvelope)
	if err != nil {
		return nil, false, err
	}

	
	responses := make(chan common.Status)
	errs := make(chan error, 1)
	go BroadcastReceive(broadcast, s.Config.Orderer.Address, responses, errs)
	status, err := BroadcastWaitForResponse(responses, errs)

	
	committed := false
	if err == nil && waitTimeout > 0 {
		committed, err = DeliverWaitForResponse(ctx, eventCh, txid)
	}

	return &status, committed, err
}



func (s *TxSubmitter) CreateTxEnvelope(txBytes []byte) (*common.Envelope, string, error) {
	var tlsCertHash []byte
	var err error
	
	cert := s.OrdererClient.Certificate()
	if cert != nil && len(cert.Certificate) > 0 {
		tlsCertHash, err = factory.GetDefault().Hash(cert.Certificate[0], &bccsp.SHA256Opts{})
		if err != nil {
			err = errors.New("failed to compute SHA256 on client certificate")
			logger.Errorf("%s", err)
			return nil, "", err
		}
	}

	txid, header, err := CreateHeader(common.HeaderType_TOKEN_TRANSACTION, s.Config.ChannelID, s.Creator, tlsCertHash)
	if err != nil {
		return nil, txid, err
	}

	txEnvelope, err := CreateEnvelope(txBytes, header, s.SigningIdentity)
	return txEnvelope, txid, err
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
		Epoch:       0,
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


func CreateEnvelope(data []byte, header *common.Header, signingIdentity tk.SigningIdentity) (*common.Envelope, error) {
	payload := &common.Payload{
		Header: header,
		Data:   data,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal common.Payload")
	}

	signature, err := signingIdentity.Sign(payloadBytes)
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
