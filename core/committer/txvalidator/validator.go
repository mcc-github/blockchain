/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/common/validation"
	"github.com/mcc-github/blockchain/core/ledger"
	ledgerUtil "github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	mspprotos "github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)


type Support interface {
	
	Acquire(ctx context.Context, n int64) error

	
	Release(n int64)

	
	Ledger() ledger.PeerLedger

	
	MSPManager() msp.MSPManager

	
	Apply(configtx *common.ConfigEnvelope) error

	
	
	GetMSPIDs(cid string) []string

	
	Capabilities() channelconfig.ApplicationCapabilities
}




type Validator interface {
	Validate(block *common.Block) error
}




type vsccValidator interface {
	VSCCValidateTx(seq int, payload *common.Payload, envBytes []byte, block *common.Block) (error, peer.TxValidationCode)
}




type TxValidator struct {
	Support Support
	Vscc    vsccValidator
}

var logger *logging.Logger 

func init() {
	
	logger = flogging.MustGetLogger("committer/txvalidator")
}

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


func NewTxValidator(support Support, sccp sysccprovider.SystemChaincodeProvider, pm PluginMapper) *TxValidator {
	
	pluginValidator := NewPluginValidator(pm, support.Ledger(), &dynamicDeserializer{support: support}, &dynamicCapabilities{support: support})
	return &TxValidator{
		Support: support,
		Vscc:    newVSCCValidator(support, sccp, pluginValidator)}
}

func (v *TxValidator) chainExists(chain string) bool {
	
	return true
}





















func (v *TxValidator) Validate(block *common.Block) error {
	var err error
	var errPos int

	logger.Debug("START Block Validation")
	defer logger.Debug("END Block Validation")
	
	txsfltr := ledgerUtil.NewTxValidationFlags(len(block.Data.Data))
	
	txsChaincodeNames := make(map[int]*sysccprovider.ChaincodeInstance)
	
	txsUpgradedChaincodes := make(map[int]*sysccprovider.ChaincodeInstance)
	
	txidArray := make([]string, len(block.Data.Data))

	results := make(chan *blockValidationResult)
	go func() {
		for tIdx, d := range block.Data.Data {
			
			v.Support.Acquire(context.Background(), 1)

			go func(index int, data []byte) {
				defer v.Support.Release(1)

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

	
	
	if v.Support.Capabilities().ForbidDuplicateTXIdInBlock() {
		markTXIdDuplicates(txidArray, txsfltr)
	}

	
	
	
	v.invalidTXsForUpgradeCC(txsChaincodeNames, txsUpgradedChaincodes, txsfltr)

	
	err = v.allValidated(txsfltr, block)
	if err != nil {
		return err
	}

	
	utils.InitBlockMetadata(block)

	block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsfltr

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

	if env, err := utils.GetEnvelopeFromBlock(d); err != nil {
		logger.Warningf("Error getting tx from block: %+v", err)
		results <- &blockValidationResult{
			tIdx:           tIdx,
			validationCode: peer.TxValidationCode_INVALID_OTHER_REASON,
		}
		return
	} else if env != nil {
		
		
		
		
		
		logger.Debugf("validateTx starts for block %p env %p txn %d", block, env, tIdx)
		defer logger.Debugf("validateTx completes for block %p env %p txn %d", block, env, tIdx)
		var payload *common.Payload
		var err error
		var txResult peer.TxValidationCode
		var txsChaincodeName *sysccprovider.ChaincodeInstance
		var txsUpgradedChaincode *sysccprovider.ChaincodeInstance

		if payload, txResult = validation.ValidateTransaction(env, v.Support.Capabilities()); txResult != peer.TxValidationCode_VALID {
			logger.Errorf("Invalid transaction with index %d", tIdx)
			results <- &blockValidationResult{
				tIdx:           tIdx,
				validationCode: txResult,
			}
			return
		}

		chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
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
			
			_, err := v.Support.Ledger().GetTransactionByID(txID)
			
			if err == nil {
				logger.Error("Duplicate transaction found, ", txID, ", skipping")
				results <- &blockValidationResult{
					tIdx:           tIdx,
					validationCode: peer.TxValidationCode_DUPLICATE_TXID,
				}
				return
			}
			
			if _, isNotFoundInIndexErrType := err.(ledger.NotFoundInIndexErr); !isNotFoundInIndexErrType {
				logger.Errorf("Ledger failure while attempting to detect duplicate status for txid %s, err '%s'. Aborting", txID, err)
				results <- &blockValidationResult{
					tIdx: tIdx,
					err:  err,
				}
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
				logger.Infof("Find chaincode upgrade transaction for chaincode %s on channel %s with new version %s", upgradeCC.ChaincodeName, upgradeCC.ChainID, upgradeCC.ChaincodeVersion)
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

			if err := v.Support.Apply(configEnvelope); err != nil {
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
		upgradedCCKey := v.generateCCKey(cc.ChaincodeName, cc.ChainID)

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
		ccKey := v.generateCCKey(cc.ChaincodeName, cc.ChainID)
		if _, exist := upgradedChaincodes[ccKey]; exist {
			if txsfltr.IsValid(tIdx) {
				logger.Infof("Invalid transaction with index %d: chaincode was upgraded in the same block", tIdx)
				txsfltr.SetFlag(tIdx, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)
			}
		}
	}
}

func (v *TxValidator) getTxCCInstance(payload *common.Payload) (invokeCCIns, upgradeCCIns *sysccprovider.ChaincodeInstance, err error) {
	
	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, nil, err
	}

	
	chainID := chdr.ChannelId 

	
	hdrExt, err := utils.GetChaincodeHeaderExtension(payload.Header)
	if err != nil {
		return nil, nil, err
	}
	invokeCC := hdrExt.ChaincodeId
	invokeIns := &sysccprovider.ChaincodeInstance{ChainID: chainID, ChaincodeName: invokeCC.Name, ChaincodeVersion: invokeCC.Version}

	
	tx, err := utils.GetTransaction(payload.Data)
	if err != nil {
		logger.Errorf("GetTransaction failed: %+v", err)
		return invokeIns, nil, nil
	}

	
	cap, err := utils.GetChaincodeActionPayload(tx.Actions[0].Payload)
	if err != nil {
		logger.Errorf("GetChaincodeActionPayload failed: %+v", err)
		return invokeIns, nil, nil
	}

	
	cpp, err := utils.GetChaincodeProposalPayload(cap.ChaincodeProposalPayload)
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
			upgradeIns, err := v.getUpgradeTxInstance(chainID, cis.ChaincodeSpec.Input.Args[2])
			if err != nil {
				return invokeIns, nil, nil
			}
			return invokeIns, upgradeIns, nil
		}
	}

	return invokeIns, nil, nil
}

func (v *TxValidator) getUpgradeTxInstance(chainID string, cdsBytes []byte) (*sysccprovider.ChaincodeInstance, error) {
	cds, err := utils.GetChaincodeDeploymentSpec(cdsBytes, platforms.NewRegistry(&golang.Platform{}))
	if err != nil {
		return nil, err
	}

	return &sysccprovider.ChaincodeInstance{
		ChainID:          chainID,
		ChaincodeName:    cds.ChaincodeSpec.ChaincodeId.Name,
		ChaincodeVersion: cds.ChaincodeSpec.ChaincodeId.Version,
	}, nil
}

type dynamicDeserializer struct {
	support Support
}

func (ds *dynamicDeserializer) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	return ds.support.MSPManager().DeserializeIdentity(serializedIdentity)
}

func (ds *dynamicDeserializer) IsWellFormed(identity *mspprotos.SerializedIdentity) error {
	return ds.support.MSPManager().IsWellFormed(identity)
}

type dynamicCapabilities struct {
	support Support
}

func (ds *dynamicCapabilities) ACLs() bool {
	return ds.support.Capabilities().ACLs()
}

func (ds *dynamicCapabilities) CollectionUpgrade() bool {
	return ds.support.Capabilities().CollectionUpgrade()
}

func (ds *dynamicCapabilities) ForbidDuplicateTXIdInBlock() bool {
	return ds.support.Capabilities().ForbidDuplicateTXIdInBlock()
}

func (ds *dynamicCapabilities) KeyLevelEndorsement() bool {
	return ds.support.Capabilities().KeyLevelEndorsement()
}

func (ds *dynamicCapabilities) MetadataLifecycle() bool {
	return ds.support.Capabilities().MetadataLifecycle()
}

func (ds *dynamicCapabilities) PrivateChannelData() bool {
	return ds.support.Capabilities().PrivateChannelData()
}

func (ds *dynamicCapabilities) Supported() error {
	return ds.support.Capabilities().Supported()
}

func (ds *dynamicCapabilities) V1_1Validation() bool {
	return ds.support.Capabilities().V1_1Validation()
}

func (ds *dynamicCapabilities) V1_2Validation() bool {
	return ds.support.Capabilities().V1_2Validation()
}

func (ds *dynamicCapabilities) V1_3Validation() bool {
	return ds.support.Capabilities().V1_3Validation()
}
