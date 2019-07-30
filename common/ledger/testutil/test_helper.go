/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/crypto"
	mmsp "github.com/mcc-github/blockchain/common/mocks/msp"
	"github.com/mcc-github/blockchain/common/util"
	lutils "github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

var signer msp.SigningIdentity

func init() {
	var err error
	
	err = msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		panic(fmt.Errorf("Could not load msp config, err %s", err))
	}
	signer, err = mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		panic(fmt.Errorf("Could not initialize msp/signer"))
	}
}


type BlockGenerator struct {
	blockNum     uint64
	previousHash []byte
	signTxs      bool
	t            *testing.T
}

type TxDetails struct {
	TxID                            string
	ChaincodeName, ChaincodeVersion string
	SimulationResults               []byte
	Type                            common.HeaderType
}

type BlockDetails struct {
	BlockNum     uint64
	PreviousHash []byte
	Txs          []*TxDetails
}


func NewBlockGenerator(t *testing.T, ledgerID string, signTxs bool) (*BlockGenerator, *common.Block) {
	gb, err := test.MakeGenesisBlock(ledgerID)
	assert.NoError(t, err)
	gb.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = lutils.NewTxValidationFlagsSetValue(len(gb.Data.Data), pb.TxValidationCode_VALID)
	return &BlockGenerator{1, protoutil.BlockHeaderHash(gb.GetHeader()), signTxs, t}, gb
}


func (bg *BlockGenerator) NextBlock(simulationResults [][]byte) *common.Block {
	block := ConstructBlock(bg.t, bg.blockNum, bg.previousHash, simulationResults, bg.signTxs)
	bg.blockNum++
	bg.previousHash = protoutil.BlockHeaderHash(block.Header)
	return block
}


func (bg *BlockGenerator) NextBlockWithTxid(simulationResults [][]byte, txids []string) *common.Block {
	
	if len(simulationResults) != len(txids) {
		return nil
	}
	block := ConstructBlockWithTxid(bg.t, bg.blockNum, bg.previousHash, simulationResults, txids, bg.signTxs)
	bg.blockNum++
	bg.previousHash = protoutil.BlockHeaderHash(block.Header)
	return block
}


func (bg *BlockGenerator) NextTestBlock(numTx int, txSize int) *common.Block {
	simulationResults := [][]byte{}
	for i := 0; i < numTx; i++ {
		simulationResults = append(simulationResults, ConstructRandomBytes(bg.t, txSize))
	}
	return bg.NextBlock(simulationResults)
}


func (bg *BlockGenerator) NextTestBlocks(numBlocks int) []*common.Block {
	blocks := []*common.Block{}
	numTx := 10
	for i := 0; i < numBlocks; i++ {
		block := bg.NextTestBlock(numTx, 100)
		block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = lutils.NewTxValidationFlagsSetValue(numTx, pb.TxValidationCode_VALID)
		blocks = append(blocks, block)
	}
	return blocks
}


func ConstructTransaction(
	t *testing.T,
	simulationResults []byte,
	txid string,
	sign bool,
) (*common.Envelope, string, error) {
	return ConstructTransactionWithHeaderType(
		t,
		simulationResults,
		txid,
		sign,
		common.HeaderType_ENDORSER_TRANSACTION,
	)
}


func ConstructTransactionWithHeaderType(
	t *testing.T,
	simulationResults []byte,
	txid string,
	sign bool,
	headerType common.HeaderType,
) (*common.Envelope, string, error) {
	return ConstructTransactionFromTxDetails(
		&TxDetails{
			ChaincodeName:     "foo",
			ChaincodeVersion:  "v1",
			TxID:              txid,
			SimulationResults: simulationResults,
			Type:              headerType,
		},
		sign,
	)
}

func ConstructTransactionFromTxDetails(txDetails *TxDetails, sign bool) (*common.Envelope, string, error) {
	ccid := &pb.ChaincodeID{
		Name:    txDetails.ChaincodeName,
		Version: txDetails.ChaincodeVersion,
	}
	var txEnv *common.Envelope
	var err error
	var txID string
	if sign {
		txEnv, txID, err = ConstructSignedTxEnvWithDefaultSigner(
			util.GetTestChainID(),
			ccid,
			nil,
			txDetails.SimulationResults,
			txDetails.TxID,
			nil,
			nil,
			txDetails.Type,
		)
	} else {
		txEnv, txID, err = ConstructUnsignedTxEnv(
			util.GetTestChainID(),
			ccid,
			nil,
			txDetails.SimulationResults,
			txDetails.TxID,
			nil,
			nil,
			txDetails.Type,
		)
	}
	return txEnv, txID, err
}

func ConstructBlockFromBlockDetails(t *testing.T, blockDetails *BlockDetails, sign bool) *common.Block {
	var envs []*common.Envelope
	for _, txDetails := range blockDetails.Txs {
		env, _, err := ConstructTransactionFromTxDetails(txDetails, sign)
		if err != nil {
			t.Fatalf("ConstructTestTransaction failed, err %s", err)
		}
		envs = append(envs, env)
	}
	return NewBlock(envs, blockDetails.BlockNum, blockDetails.PreviousHash)
}

func ConstructBlockWithTxid(
	t *testing.T,
	blockNum uint64,
	previousHash []byte,
	simulationResults [][]byte,
	txids []string,
	sign bool,
) *common.Block {
	return ConstructBlockWithTxidHeaderType(
		t,
		blockNum,
		previousHash,
		simulationResults,
		txids,
		sign,
		common.HeaderType_ENDORSER_TRANSACTION,
	)
}

func ConstructBlockWithTxidHeaderType(
	t *testing.T,
	blockNum uint64,
	previousHash []byte,
	simulationResults [][]byte,
	txids []string,
	sign bool,
	headerType common.HeaderType,
) *common.Block {
	envs := []*common.Envelope{}
	for i := 0; i < len(simulationResults); i++ {
		env, _, err := ConstructTransactionWithHeaderType(
			t,
			simulationResults[i],
			txids[i],
			sign,
			headerType,
		)
		if err != nil {
			t.Fatalf("ConstructTestTransaction failed, err %s", err)
		}
		envs = append(envs, env)
	}
	return NewBlock(envs, blockNum, previousHash)
}


func ConstructBlock(
	t *testing.T,
	blockNum uint64,
	previousHash []byte,
	simulationResults [][]byte,
	sign bool,
) *common.Block {
	envs := []*common.Envelope{}
	for i := 0; i < len(simulationResults); i++ {
		env, _, err := ConstructTransaction(t, simulationResults[i], "", sign)
		if err != nil {
			t.Fatalf("ConstructTestTransaction failed, err %s", err)
		}
		envs = append(envs, env)
	}
	return NewBlock(envs, blockNum, previousHash)
}


func ConstructTestBlock(t *testing.T, blockNum uint64, numTx int, txSize int) *common.Block {
	simulationResults := [][]byte{}
	for i := 0; i < numTx; i++ {
		simulationResults = append(simulationResults, ConstructRandomBytes(t, txSize))
	}
	return ConstructBlock(t, blockNum, ConstructRandomBytes(t, 32), simulationResults, false)
}




func ConstructTestBlocks(t *testing.T, numBlocks int) []*common.Block {
	bg, gb := NewBlockGenerator(t, util.GetTestChainID(), false)
	blocks := []*common.Block{}
	if numBlocks != 0 {
		blocks = append(blocks, gb)
	}
	return append(blocks, bg.NextTestBlocks(numBlocks-1)...)
}


func ConstructBytesProposalResponsePayload(version string, simulationResults []byte) ([]byte, error) {
	ccid := &pb.ChaincodeID{
		Name:    "foo",
		Version: version,
	}
	return constructBytesProposalResponsePayload(util.GetTestChainID(), ccid, nil, simulationResults)
}

func NewBlock(env []*common.Envelope, blockNum uint64, previousHash []byte) *common.Block {
	block := protoutil.NewBlock(blockNum, previousHash)
	for i := 0; i < len(env); i++ {
		txEnvBytes, _ := proto.Marshal(env[i])
		block.Data.Data = append(block.Data.Data, txEnvBytes)
	}
	block.Header.DataHash = protoutil.BlockDataHash(block.Data)
	protoutil.InitBlockMetadata(block)

	block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = lutils.NewTxValidationFlagsSetValue(len(env), pb.TxValidationCode_VALID)

	return block
}


func constructBytesProposalResponsePayload(chainID string, ccid *pb.ChaincodeID, pResponse *pb.Response, simulationResults []byte) ([]byte, error) {
	ss, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	prop, _, err := protoutil.CreateChaincodeProposal(common.HeaderType_ENDORSER_TRANSACTION, chainID, &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: ccid}}, ss)
	if err != nil {
		return nil, err
	}

	presp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, pResponse, simulationResults, nil, ccid, nil, signer)
	if err != nil {
		return nil, err
	}

	return presp.Payload, nil
}



func ConstructSignedTxEnvWithDefaultSigner(
	chainID string,
	ccid *pb.ChaincodeID,
	response *pb.Response,
	simulationResults []byte,
	txid string,
	events []byte,
	visibility []byte,
	headerType common.HeaderType,
) (*common.Envelope, string, error) {

	return ConstructSignedTxEnv(
		chainID,
		ccid,
		response,
		simulationResults,
		txid,
		events,
		visibility,
		signer,
		headerType,
	)
}


func ConstructUnsignedTxEnv(
	chainID string,
	ccid *pb.ChaincodeID,
	response *pb.Response,
	simulationResults []byte,
	txid string,
	events []byte,
	visibility []byte,
	headerType common.HeaderType,
) (*common.Envelope, string, error) {
	mspLcl := mmsp.NewNoopMsp()
	sigId, _ := mspLcl.GetDefaultSigningIdentity()

	return ConstructSignedTxEnv(
		chainID,
		ccid,
		response,
		simulationResults,
		txid,
		events,
		visibility,
		sigId,
		headerType,
	)
}


func ConstructSignedTxEnv(
	chainID string,
	ccid *pb.ChaincodeID,
	pResponse *pb.Response,
	simulationResults []byte,
	txid string,
	events []byte,
	visibility []byte,
	signer msp.SigningIdentity,
	headerType common.HeaderType,
) (*common.Envelope, string, error) {
	ss, err := signer.Serialize()
	if err != nil {
		return nil, "", err
	}

	var prop *pb.Proposal
	if txid == "" {
		
		prop, txid, err = protoutil.CreateChaincodeProposal(
			common.HeaderType_ENDORSER_TRANSACTION,
			chainID,
			&pb.ChaincodeInvocationSpec{
				ChaincodeSpec: &pb.ChaincodeSpec{
					ChaincodeId: ccid,
				},
			},
			ss,
		)

	} else {
		
		nonce, err := crypto.GetRandomNonce()
		if err != nil {
			return nil, "", err
		}
		prop, txid, err = protoutil.CreateChaincodeProposalWithTxIDNonceAndTransient(
			txid,
			headerType,
			chainID,
			&pb.ChaincodeInvocationSpec{
				ChaincodeSpec: &pb.ChaincodeSpec{
					ChaincodeId: ccid,
				},
			},
			nonce,
			ss,
			nil,
		)
	}
	if err != nil {
		return nil, "", err
	}

	presp, err := protoutil.CreateProposalResponse(
		prop.Header,
		prop.Payload,
		pResponse,
		simulationResults,
		nil,
		ccid,
		nil,
		signer,
	)
	if err != nil {
		return nil, "", err
	}

	env, err := protoutil.CreateSignedTx(prop, signer, presp)
	if err != nil {
		return nil, "", err
	}
	return env, txid, nil
}

func SetTxID(t *testing.T, block *common.Block, txNum int, txID string) {
	envelopeBytes := block.Data.Data[txNum]
	envelope, err := protoutil.UnmarshalEnvelope(envelopeBytes)
	if err != nil {
		t.Fatalf("error unmarshaling envelope: %s", err)
	}

	payload, err := protoutil.GetPayload(envelope)
	if err != nil {
		t.Fatalf("error getting payload from envelope: %s", err)
	}

	channelHeader, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		t.Fatalf("error unmarshaling channel header: %s", err)
	}
	channelHeader.TxId = txID
	channelHeaderBytes, err := proto.Marshal(channelHeader)
	if err != nil {
		t.Fatalf("error marshaling channel header: %s", err)
	}
	payload.Header.ChannelHeader = channelHeaderBytes

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		t.Fatalf("error marshaling payload: %s", err)
	}

	envelope.Payload = payloadBytes
	envelopeBytes, err = proto.Marshal(envelope)
	if err != nil {
		t.Fatalf("error marshaling envelope: %s", err)
	}
	block.Data.Data[txNum] = envelopeBytes
}
