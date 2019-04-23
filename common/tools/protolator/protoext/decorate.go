/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoext

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/tools/protolator/protoext/commonext"
	"github.com/mcc-github/blockchain/common/tools/protolator/protoext/ledger/rwsetext"
	"github.com/mcc-github/blockchain/common/tools/protolator/protoext/mspext"
	"github.com/mcc-github/blockchain/common/tools/protolator/protoext/ordererext"
	"github.com/mcc-github/blockchain/common/tools/protolator/protoext/peerext"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/peer"
)



func Decorate(msg proto.Message) proto.Message {
	switch m := msg.(type) {
	case *common.BlockData:
		return &commonext.BlockData{BlockData: m}
	case *common.Config:
		return &commonext.Config{Config: m}
	case *common.ConfigSignature:
		return &commonext.ConfigSignature{ConfigSignature: m}
	case *common.ConfigUpdate:
		return &commonext.ConfigUpdate{ConfigUpdate: m}
	case *common.ConfigUpdateEnvelope:
		return &commonext.ConfigUpdateEnvelope{ConfigUpdateEnvelope: m}
	case *common.Envelope:
		return &commonext.Envelope{Envelope: m}
	case *common.Header:
		return &commonext.Header{Header: m}
	case *common.SignatureHeader:
		return &commonext.SignatureHeader{SignatureHeader: m}
	case *common.Payload:
		return &commonext.Payload{Payload: m}
	case *common.Policy:
		return &commonext.Policy{Policy: m}

	case *msp.MSPConfig:
		return &mspext.MSPConfig{MSPConfig: m}
	case *msp.MSPPrincipal:
		return &mspext.MSPPrincipal{MSPPrincipal: m}

	case *orderer.ConsensusType:
		return &ordererext.ConsensusType{ConsensusType: m}

	case *peer.ChaincodeAction:
		return &peerext.ChaincodeAction{ChaincodeAction: m}
	case *peer.ChaincodeActionPayload:
		return &peerext.ChaincodeActionPayload{ChaincodeActionPayload: m}
	case *peer.ChaincodeEndorsedAction:
		return &peerext.ChaincodeEndorsedAction{ChaincodeEndorsedAction: m}
	case *peer.ChaincodeProposalPayload:
		return &peerext.ChaincodeProposalPayload{ChaincodeProposalPayload: m}
	case *peer.ProposalResponsePayload:
		return &peerext.ProposalResponsePayload{ProposalResponsePayload: m}
	case *peer.TransactionAction:
		return &peerext.TransactionAction{TransactionAction: m}

	case *rwset.TxReadWriteSet:
		return &rwsetext.TxReadWriteSet{TxReadWriteSet: m}

	default:
		return msg
	}
}
