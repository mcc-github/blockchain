/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"

	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/protoutil"
)

func TestValidateConfigTx(t *testing.T) {
	channelID := util.GetTestChannelID()
	profile := configtxgentest.Load(genesisconfig.SampleSingleMSPChannelProfile)
	chCrtEnv, err := encoder.MakeChannelCreationTransaction(genesisconfig.SampleConsortiumName, nil, profile)
	if err != nil {
		t.Fatalf("MakeChannelCreationTransaction failed, err %s", err)
		return
	}

	updateResult := &cb.Envelope{
		Payload: protoutil.MarshalOrPanic(&cb.Payload{Header: &cb.Header{
			ChannelHeader: protoutil.MarshalOrPanic(&cb.ChannelHeader{
				Type:      int32(cb.HeaderType_CONFIG),
				ChannelId: channelID,
			}),
			SignatureHeader: protoutil.MarshalOrPanic(&cb.SignatureHeader{
				Creator: signerSerialized,
				Nonce:   protoutil.CreateNonceOrPanic(),
			}),
		},
			Data: protoutil.MarshalOrPanic(&cb.ConfigEnvelope{
				LastUpdate: chCrtEnv,
			}),
		}),
	}
	updateResult.Signature, _ = signer.Sign(updateResult.Payload)
	_, txResult := ValidateTransaction(updateResult)
	if txResult != peer.TxValidationCode_VALID {
		t.Fatalf("ValidateTransaction failed, err %s", err)
		return
	}
}
