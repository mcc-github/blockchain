/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	mspprotos "github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/common/validation"
	"github.com/mcc-github/blockchain/core/ledger"
	ledgerUtil "github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)


type Semaphore interface {
	
	Acquire(ctx context.Context) error

	
	Release()
}



type ChannelResources interface {
	
	Ledger() ledger.PeerLedger

	
	MSPManager() msp.MSPManager

	
	Apply(configtx *common.ConfigEnvelope) error

	
	
	GetMSPIDs() []string

	
	Capabilities() channelconfig.ApplicationCapabilities
}




type vsccValidator interface {
	VSCCValidateTx(seq int, payload *common.Payload, envBytes []byte, block *common.Block) (error, peer.TxValidationCode)
}




type TxValidator struct {
	ChannelID        string
	Semaphore        Semaphore
	ChannelResources ChannelResources
	Vscc             vsccValidator
}

var logger = flogging.MustGetLogger("committer.txvalidator")

type blockValidationRequest struct {
	block *common.Block
	d     []byte
	tIdx  int
}

type blockValidationResult struct {
	tIdx                 int
	validationCode       peer.TxValidationCode
	txsChaincodeName     *sysccprovider.ChaincodeInstance
	txsUpgradedChaincode *sysccprovider.ChaincodeInstance
	err                  error
	txid                 string
}


func NewTxValidator(channelID string, sem Semaphore, cr ChannelResources, pm plugin.Mapper) *TxValidator {
	
	pluginValidator := NewPluginValidator(pm, cr.Ledger(), &dynamicDeserializer{cr: cr}, &dynamicCapabilities{cr: cr})
	return &TxValidator{
		ChannelID:        channelID,
		Semaphore:        sem,
		ChannelResources: cr,
		Vscc:             newVSCCValidator(channelID, cr, pluginValidator)}
}

func (v *TxValidator) chainExists(chain string) bool {
	
	return true
}





















func (v *TxValidator) Validate(block *common.Block) error {
	var err error
	var errPos int

	startValidation := time.Now() 
	logger.Debugf("[%s] START Block Validation for block [%d]", v.ChannelID, block.Header.Number)

	
	txsfltr := ledgerUtil.NewTxValidationFlags(len(block.Data.Data))
	
	txsChaincodeNames := make(map[int]*sysccprovider.ChaincodeInstance)
	
	txsUpgradedChaincodes := make(map[int]*sysccprovider.ChaincodeInstance)
	
	txidArray := make([]string, len(block.Data.Data))

	results := make(chan *blockValidationResult)
	go func() {
		for tIdx, d := range block.Data.Data {
			
			v.Semaphore.Acquire(context.Background())

			go func(index int, data []byte) {
				defer v.Semaphore.Release()

				v.validateTx(&blockValidationRequest{
					d:     data,
					block: block,
					tIdx:  index,
				}, results)
			}(tIdx, d)
		}
	}()

	logger.Debugf("expecting %d block validation responses", len(block.Data.Data))

	
	for i := 0; i < len(block.Data.Data); i++ {
		res := <-results

		if res.err != nil {
			
			
			
			logger.Debugf("got terminal error %s for idx %d", res.err, res.tIdx)

			if err == nil || res.tIdx < errPos {
				err = res.err
				errPos = res.tIdx
			}
		} else {
			
			
			logger.Debugf("got result for idx %d, code %d", res.tIdx, res.validationCode)

			txsfltr.SetFlag(res.tIdx, res.validationCode)

			if res.validationCode == peer.TxValidationCode_VALID {
				if res.txsChaincodeName != nil {
					txsChaincodeNames[res.tIdx] = res.txsChaincodeName
				}
				if res.txsUpgradedChaincode != nil {
					txsUpgradedChaincodes[res.tIdx] = res.txsUpgradedChaincode
				}
				txidArray[res.tIdx] = res.txid
			}
		}
	}

	
	
	
	if err != nil {
		return err
	}

	
	
	if v.ChannelResources.Capabilities().ForbidDuplicateTXIdInBlock() {
		markTXIdDuplicates(txidArray, txsfltr)
	}

	
	
	
	v.invalidTXsForUpgradeCC(txsChaincodeNames, txsUpgradedChaincodes, txsfltr)

	
	err = v.allValidated(txsfltr, block)
	if err != nil {
		return err
	}

	
	protoutil.InitBlockMetadata(block)

	block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsfltr

	elapsedValidation := time.Since(startValidation) / time.Millisecond 
	logger.Infof("[%s] Validated block [%d] in %dms", v.ChannelID, block.Header.Number, elapsedValidation)

	return nil
}



func (v *TxValidator) allValidated(txsfltr ledgerUtil.TxValidationFlags, block *common.Block) error {
	for id, f := range txsfltr {
		if peer.TxValidationCode(f) == peer.TxValidationCode_NOT_VALIDATED {
			return errors.Errorf("transaction %d in block %d has skipped validation", id, block.Header.Number)
		}
	}

	return nil
}

func markTXIdDuplicates(txids []string, txsfltr ledgerUtil.TxValidationFlags) {
	txidMap := make(map[string]struct{})

	for id, txid := range txids {
		if txid == "" {
			continue
		}

		_, in := txidMap[txid]
		if in {
			logger.Error("Duplicate txid", txid, "found, skipping")
			txsfltr.SetFlag(id, peer.TxValidationCode_DUPLICATE_TXID)
		} else {
			txidMap[txid] = struct{}{}
		}
	}
}

func (v *TxValidator) validateTx(req *blockValidationRequest, results chan<- *blockValidationResult) {
	block := req.block
	d := req.d
	tIdx := req.tIdx
	txID := ""

	if d == nil {
		results <- &blockValidationResult{
			tIdx: tIdx,
		}
		return
	}

	if env, err := protoutil.GetEnvelopeFromBlock(d); err != nil {
		logger.Warningf("Error getting tx from block: %+v", err)
		results <- &blockValidationResult{
			tIdx:           tIdx,
			validationCode: peer.TxValidationCode_INVALID_OTHER_REASON,
		}
		return
	} else if env != nil {
		
		
		
		
		
		logger.Debugf("[%s] validateTx starts for block %p env %p txn %d", v.ChannelID, block, env, tIdx)
		defer logger.Debugf("[%s] validateTx completes for block %p env %p txn %d", v.ChannelID, block, env, tIdx)
		var payload *common.Payload
		var err error
		var txResult peer.TxValidationCode
		var txsChaincodeName *sysccprovider.ChaincodeInstance
		var txsUpgradedChaincode *sysccprovider.ChaincodeInstance

		if payload, txResult = validation.ValidateTransaction(env); txResult != peer.TxValidationCode_VALID {
			logger.Errorf("Invalid transaction with index %d", tIdx)
			results <- &blockValidationResult{
				tIdx:           tIdx,
				validationCode: txResult,
			}
			return
		}

		chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			logger.Warningf("Could not unmarshal channel header, err %s, skipping", err)
			results <- &blockValidationResult{
				tIdx:           tIdx,
				validationCode: peer.TxValidationCode_INVALID_OTHER_REASON,
			}
			return
		}

		channel := chdr.ChannelId
		logger.Debugf("Transaction is for channel %s", channel)

		if !v.chainExists(channel) {
			logger.Errorf("Dropping transaction for non-existent channel %s", channel)
			results <- &blockValidationResult{
				tIdx:           tIdx,
				validationCode: peer.TxValidationCode_TARGET_CHAIN_NOT_FOUND,
			}
			return
		}

		if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

			txID = chdr.TxId

			
			erroneousResultEntry := v.checkTxIdDupsLedger(tIdx, chdr, v.ChannelResources.Ledger())
			if erroneousResultEntry != nil {
				results <- erroneousResultEntry
				return
			}

			
			logger.Debug("Validating transaction vscc tx validate")
			err, cde := v.Vscc.VSCCValidateTx(tIdx, payload, d, block)
			if err != nil {
				logger.Errorf("VSCCValidateTx for transaction txId = %s returned error: %s", txID, err)
				switch err.(type) {
				case *commonerrors.VSCCExecutionFailureError:
					results <- &blockValidationResult{
						tIdx: tIdx,
						err:  err,
					}
					return
				case *commonerrors.VSCCInfoLookupFailureError:
					results <- &blockValidationResult{
						tIdx: tIdx,
						err:  err,
					}
					return
				default:
					results <- &blockValidationResult{
						tIdx:           tIdx,
						validationCode: cde,
					}
					return
				}
			}

			invokeCC, upgradeCC, err := v.getTxCCInstance(payload)
			if err != nil {
				logger.Errorf("Get chaincode instance from transaction txId = %s returned error: %+v", txID, err)
				results <- &blockValidationResult{
					tIdx:           tIdx,
					validationCode: peer.TxValidationCode_INVALID_OTHER_REASON,
				}
				return
			}
			txsChaincodeName = invokeCC
			if upgradeCC != nil {
				logger.Infof("Find chaincode upgrade transaction for chaincode %s on channel %s with new version %s", upgradeCC.ChaincodeName, upgradeCC.ChannelID, upgradeCC.ChaincodeVersion)
				txsUpgradedChaincode = upgradeCC
			}
		} else if common.HeaderType(chdr.Type) == common.HeaderType_CONFIG {
			configEnvelope, err := configtx.UnmarshalConfigEnvelope(payload.Data)
			if err != nil {
				err = errors.WithMessage(err, "error unmarshalling config which passed initial validity checks")
				logger.Criticalf("%+v", err)
				results <- &blockValidationResult{
					tIdx: tIdx,
					err:  err,
				}
				return
			}

			if err := v.ChannelResources.Apply(configEnvelope); err != nil {
				err = errors.WithMessage(err, "error validating config which passed initial validity checks")
				logger.Criticalf("%+v", err)
				results <- &blockValidationResult{
					tIdx: tIdx,
					err:  err,
				}
				return
			}
			logger.Debugf("config transaction received for chain %s", channel)
		} else {
			logger.Warningf("Unknown transaction type [%s] in block number [%d] transaction index [%d]",
				common.HeaderType(chdr.Type), block.Header.Number, tIdx)
			results <- &blockValidationResult{
				tIdx:           tIdx,
				validationCode: peer.TxValidationCode_UNKNOWN_TX_TYPE,
			}
			return
		}

		if _, err := proto.Marshal(env); err != nil {
			logger.Warningf("Cannot marshal transaction: %s", err)
			results <- &blockValidationResult{
				tIdx:           tIdx,
				validationCode: peer.TxValidationCode_MARSHAL_TX_ERROR,
			}
			return
		}
		
		results <- &blockValidationResult{
			tIdx:                 tIdx,
			txsChaincodeName:     txsChaincodeName,
			txsUpgradedChaincode: txsUpgradedChaincode,
			validationCode:       peer.TxValidationCode_VALID,
			txid:                 txID,
		}
		return
	} else {
		logger.Warning("Nil tx from block")
		results <- &blockValidationResult{
			tIdx:           tIdx,
			validationCode: peer.TxValidationCode_NIL_ENVELOPE,
		}
		return
	}
}






func (v *TxValidator) checkTxIdDupsLedger(tIdx int, chdr *common.ChannelHeader, ldgr ledger.PeerLedger) *blockValidationResult {

	
	txID := chdr.TxId

	
	_, err := ldgr.GetTransactionByID(txID)

	switch err.(type) {
	case nil:
		
		logger.Error("Duplicate transaction found, ", txID, ", skipping")
		return &blockValidationResult{
			tIdx:           tIdx,
			validationCode: peer.TxValidationCode_DUPLICATE_TXID,
		}
	case ledger.NotFoundInIndexErr:
		
		
		return nil
	default:
		
		
		logger.Errorf("Ledger failure while attempting to detect duplicate status for txid %s: %s", txID, err)
		return &blockValidationResult{
			tIdx: tIdx,
			err:  err,
		}
	}
}


func (v *TxValidator) generateCCKey(ccName, chainID string) string {
	return fmt.Sprintf("%s/%s", ccName, chainID)
}


func (v *TxValidator) invalidTXsForUpgradeCC(txsChaincodeNames map[int]*sysccprovider.ChaincodeInstance, txsUpgradedChaincodes map[int]*sysccprovider.ChaincodeInstance, txsfltr ledgerUtil.TxValidationFlags) {
	if len(txsUpgradedChaincodes) == 0 {
		return
	}

	
	finalValidUpgradeTXs := make(map[string]int)
	upgradedChaincodes := make(map[string]*sysccprovider.ChaincodeInstance)
	for tIdx, cc := range txsUpgradedChaincodes {
		if cc == nil {
			continue
		}
		upgradedCCKey := v.generateCCKey(cc.ChaincodeName, cc.ChannelID)

		if finalIdx, exist := finalValidUpgradeTXs[upgradedCCKey]; !exist {
			finalValidUpgradeTXs[upgradedCCKey] = tIdx
			upgradedChaincodes[upgradedCCKey] = cc
		} else if finalIdx < tIdx {
			logger.Infof("Invalid transaction with index %d: chaincode was upgraded by latter tx", finalIdx)
			txsfltr.SetFlag(finalIdx, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)

			
			finalValidUpgradeTXs[upgradedCCKey] = tIdx
			upgradedChaincodes[upgradedCCKey] = cc
		} else {
			logger.Infof("Invalid transaction with index %d: chaincode was upgraded by latter tx", tIdx)
			txsfltr.SetFlag(tIdx, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)
		}
	}

	
	for tIdx, cc := range txsChaincodeNames {
		if cc == nil {
			continue
		}
		ccKey := v.generateCCKey(cc.ChaincodeName, cc.ChannelID)
		if _, exist := upgradedChaincodes[ccKey]; exist {
			if txsfltr.IsValid(tIdx) {
				logger.Infof("Invalid transaction with index %d: chaincode was upgraded in the same block", tIdx)
				txsfltr.SetFlag(tIdx, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)
			}
		}
	}
}

func (v *TxValidator) getTxCCInstance(payload *common.Payload) (invokeCCIns, upgradeCCIns *sysccprovider.ChaincodeInstance, err error) {
	
	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, nil, err
	}

	
	channelID := chdr.ChannelId 

	
	hdrExt, err := protoutil.UnmarshalChaincodeHeaderExtension(chdr.Extension)
	if err != nil {
		return nil, nil, err
	}
	invokeCC := hdrExt.ChaincodeId
	invokeIns := &sysccprovider.ChaincodeInstance{ChannelID: channelID, ChaincodeName: invokeCC.Name, ChaincodeVersion: invokeCC.Version}

	
	tx, err := protoutil.UnmarshalTransaction(payload.Data)
	if err != nil {
		logger.Errorf("GetTransaction failed: %+v", err)
		return invokeIns, nil, nil
	}

	
	cap, err := protoutil.UnmarshalChaincodeActionPayload(tx.Actions[0].Payload)
	if err != nil {
		logger.Errorf("GetChaincodeActionPayload failed: %+v", err)
		return invokeIns, nil, nil
	}

	
	cpp, err := protoutil.UnmarshalChaincodeProposalPayload(cap.ChaincodeProposalPayload)
	if err != nil {
		logger.Errorf("GetChaincodeProposalPayload failed: %+v", err)
		return invokeIns, nil, nil
	}

	
	cis := &peer.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(cpp.Input, cis)
	if err != nil {
		logger.Errorf("GetChaincodeInvokeSpec failed: %+v", err)
		return invokeIns, nil, nil
	}

	if invokeCC.Name == "lscc" {
		if string(cis.ChaincodeSpec.Input.Args[0]) == "upgrade" {
			upgradeIns, err := v.getUpgradeTxInstance(channelID, cis.ChaincodeSpec.Input.Args[2])
			if err != nil {
				return invokeIns, nil, nil
			}
			return invokeIns, upgradeIns, nil
		}
	}

	return invokeIns, nil, nil
}

func (v *TxValidator) getUpgradeTxInstance(channelID string, cdsBytes []byte) (*sysccprovider.ChaincodeInstance, error) {
	cds, err := protoutil.UnmarshalChaincodeDeploymentSpec(cdsBytes)
	if err != nil {
		return nil, err
	}

	if cds.ChaincodeSpec.Type.String() != "GOLANG" {
		return nil, errors.Errorf("unexpected chaincode type: %s", cds.ChaincodeSpec.Type.String())
	}

	return &sysccprovider.ChaincodeInstance{
		ChannelID:        channelID,
		ChaincodeName:    cds.ChaincodeSpec.ChaincodeId.Name,
		ChaincodeVersion: cds.ChaincodeSpec.ChaincodeId.Version,
	}, nil
}

type dynamicDeserializer struct {
	cr ChannelResources
}

func (ds *dynamicDeserializer) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	return ds.cr.MSPManager().DeserializeIdentity(serializedIdentity)
}

func (ds *dynamicDeserializer) IsWellFormed(identity *mspprotos.SerializedIdentity) error {
	return ds.cr.MSPManager().IsWellFormed(identity)
}

type dynamicCapabilities struct {
	cr ChannelResources
}

func (ds *dynamicCapabilities) ACLs() bool {
	return ds.cr.Capabilities().ACLs()
}

func (ds *dynamicCapabilities) CollectionUpgrade() bool {
	return ds.cr.Capabilities().CollectionUpgrade()
}

func (ds *dynamicCapabilities) StorePvtDataOfInvalidTx() bool {
	return ds.cr.Capabilities().StorePvtDataOfInvalidTx()
}

func (ds *dynamicCapabilities) ForbidDuplicateTXIdInBlock() bool {
	return ds.cr.Capabilities().ForbidDuplicateTXIdInBlock()
}

func (ds *dynamicCapabilities) KeyLevelEndorsement() bool {
	return ds.cr.Capabilities().KeyLevelEndorsement()
}

func (ds *dynamicCapabilities) MetadataLifecycle() bool {
	
	return false
}

func (ds *dynamicCapabilities) PrivateChannelData() bool {
	return ds.cr.Capabilities().PrivateChannelData()
}

func (ds *dynamicCapabilities) Supported() error {
	return ds.cr.Capabilities().Supported()
}

func (ds *dynamicCapabilities) V1_1Validation() bool {
	return ds.cr.Capabilities().V1_1Validation()
}

func (ds *dynamicCapabilities) V1_2Validation() bool {
	return ds.cr.Capabilities().V1_2Validation()
}

func (ds *dynamicCapabilities) V1_3Validation() bool {
	return ds.cr.Capabilities().V1_3Validation()
}

func (ds *dynamicCapabilities) V2_0Validation() bool {
	return ds.cr.Capabilities().V2_0Validation()
}
