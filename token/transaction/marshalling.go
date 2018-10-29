/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transaction

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/protos/common"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/mcc-github/blockchain/token/identity"
	"github.com/pkg/errors"
)

func UnmarshalTokenTransaction(raw []byte) (*cb.ChannelHeader, *token.TokenTransaction, identity.PublicInfo, error) {
	
	payload := &common.Payload{}
	err := proto.Unmarshal(raw, payload)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error unmarshaling Payload")
	}

	
	sh, err := utils.GetSignatureHeader(payload.Header.SignatureHeader)
	if err != nil {
		return nil, nil, nil, err
	}
	creatorInfo := &TxCreatorInfo{public: sh.Creator}

	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, nil, nil, err
	}

	
	if common.HeaderType(chdr.Type) != common.HeaderType_TOKEN_TRANSACTION {
		return nil, nil, nil, errors.Errorf("only token transactions are supported, provided type: %d", chdr.Type)
	}

	ttx := &token.TokenTransaction{}
	err = proto.Unmarshal(payload.Data, ttx)
	if err != nil {
		return nil, nil, nil, errors.Errorf("failed getting token token transaction, %s", err)
	}

	return chdr, ttx, creatorInfo, nil
}
