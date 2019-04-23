/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugindispatcher

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/common/flogging"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	s "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)



type ChannelResources interface {
	
	
	GetMSPIDs(cid string) []string
}



type LedgerResources interface {
	
	
	
	NewQueryExecutor() (ledger.QueryExecutor, error)
}



type LifecycleResources interface {
	
	
	
	
	
	
	ValidationInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (plugin string, args []byte, unexpectedErr error, validationErr error)
}


type CollectionResources interface {
	
	
	
	
	
	
	CollectionValidationInfo(chaincodeName, collectionName string, state s.State) (args []byte, unexpectedErr error, validationErr error)
}




type CollectionAndLifecycleResources interface {
	LifecycleResources

	
	
	CollectionValidationInfo(channelID, chaincodeName, collectionName string, state s.State) (args []byte, unexpectedErr error, validationErr error)
}



var logger = flogging.MustGetLogger("committer.txvalidator")



type dispatcherImpl struct {
	chainID         string
	cr              ChannelResources
	ler             LedgerResources
	lcr             LifecycleResources
	pluginValidator *PluginValidator
}


func New(chainID string, cr ChannelResources, ler LedgerResources, lcr LifecycleResources, pluginValidator *PluginValidator) *dispatcherImpl {
	return &dispatcherImpl{
		chainID:         chainID,
		cr:              cr,
		ler:             ler,
		lcr:             lcr,
		pluginValidator: pluginValidator,
	}
}


func (v *dispatcherImpl) Dispatch(seq int, payload *common.Payload, envBytes []byte, block *common.Block) (error, peer.TxValidationCode) {
	chainID := v.chainID
	logger.Debugf("[%s] Dispatch starts for bytes %p", chainID, envBytes)

	
	hdrExt, err := protoutil.GetChaincodeHeaderExtension(payload.Header)
	if err != nil {
		return err, peer.TxValidationCode_BAD_HEADER_EXTENSION
	}

	
	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return err, peer.TxValidationCode_BAD_CHANNEL_HEADER
	}

	
	respPayload, err := protoutil.GetActionFromEnvelope(envBytes)
	if err != nil {
		return errors.WithMessage(err, "GetActionFromEnvelope failed"), peer.TxValidationCode_BAD_RESPONSE_PAYLOAD
	}
	txRWSet := &rwsetutil.TxRwSet{}
	if err = txRWSet.FromProtoBytes(respPayload.Results); err != nil {
		return errors.WithMessage(err, "txRWSet.FromProtoBytes failed"), peer.TxValidationCode_BAD_RWSET
	}

	
	if hdrExt.ChaincodeId == nil {
		return errors.New("nil ChaincodeId in header extension"), peer.TxValidationCode_INVALID_OTHER_REASON
	}

	if respPayload.ChaincodeId == nil {
		return errors.New("nil ChaincodeId in ChaincodeAction"), peer.TxValidationCode_INVALID_OTHER_REASON
	}

	
	ccID := hdrExt.ChaincodeId.Name
	ccVer := respPayload.ChaincodeId.Version

	
	if ccID == "" {
		err = errors.New("invalid chaincode ID")
		logger.Errorf("%+v", err)
		return err, peer.TxValidationCode_INVALID_CHAINCODE
	}
	if ccID != respPayload.ChaincodeId.Name {
		err = errors.Errorf("inconsistent ccid info (%s/%s)", ccID, respPayload.ChaincodeId.Name)
		logger.Errorf("%+v", err)
		return err, peer.TxValidationCode_INVALID_CHAINCODE
	}
	
	if ccVer == "" {
		err = errors.New("invalid chaincode version")
		logger.Errorf("%+v", err)
		return err, peer.TxValidationCode_INVALID_CHAINCODE
	}

	wrNamespace := map[string]bool{}
	wrNamespace[ccID] = true
	if respPayload.Events != nil {
		ccEvent := &peer.ChaincodeEvent{}
		if err = proto.Unmarshal(respPayload.Events, ccEvent); err != nil {
			return errors.Wrapf(err, "invalid chaincode event"), peer.TxValidationCode_INVALID_OTHER_REASON
		}
		if ccEvent.ChaincodeId != ccID {
			return errors.Errorf("chaincode event chaincode id does not match chaincode action chaincode id"), peer.TxValidationCode_INVALID_OTHER_REASON
		}
	}

	for _, ns := range txRWSet.NsRwSets {
		if v.txWritesToNamespace(ns) {
			wrNamespace[ns.NameSpace] = true
		}
	}

	
	

	
	for ns := range wrNamespace {
		
		validationPlugin, args, err := v.GetInfoForValidate(chdr, ns)
		if err != nil {
			logger.Errorf("GetInfoForValidate for txId = %s returned error: %+v", chdr.TxId, err)
			return err, peer.TxValidationCode_INVALID_CHAINCODE
		}

		
		ctx := &Context{
			Seq:        seq,
			Envelope:   envBytes,
			Block:      block,
			TxID:       chdr.TxId,
			Channel:    chdr.ChannelId,
			Namespace:  ns,
			Policy:     args,
			PluginName: validationPlugin,
		}
		if err = v.invokeValidationPlugin(ctx); err != nil {
			switch err.(type) {
			case *commonerrors.VSCCEndorsementPolicyError:
				return err, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE
			default:
				return err, peer.TxValidationCode_INVALID_OTHER_REASON
			}
		}
	}

	logger.Debugf("[%s] Dispatch completes env bytes %p", chainID, envBytes)
	return nil, peer.TxValidationCode_VALID
}

func (v *dispatcherImpl) invokeValidationPlugin(ctx *Context) error {
	logger.Debug("Validating", ctx, "with plugin")
	err := v.pluginValidator.ValidateWithPlugin(ctx)
	if err == nil {
		return nil
	}
	
	if e, isExecutionError := err.(*validation.ExecutionFailureError); isExecutionError {
		return &commonerrors.VSCCExecutionFailureError{Err: e}
	}
	
	return &commonerrors.VSCCEndorsementPolicyError{Err: err}
}

func (v *dispatcherImpl) getCDataForCC(channelID, ccid string) (string, []byte, error) {
	qe, err := v.ler.NewQueryExecutor()
	if err != nil {
		return "", nil, errors.WithMessage(err, "could not retrieve QueryExecutor")
	}
	defer qe.Done()

	plugin, args, unexpectedErr, validationErr := v.lcr.ValidationInfo(channelID, ccid, qe)
	if unexpectedErr != nil {
		return "", nil, &commonerrors.VSCCInfoLookupFailureError{
			Reason: fmt.Sprintf("Could not retrieve state for chaincode %s, error %s", ccid, unexpectedErr),
		}
	}
	if validationErr != nil {
		return "", nil, validationErr
	}

	if plugin == "" {
		return "", nil, errors.Errorf("chaincode definition for [%s] is invalid, plugin field must be set", ccid)
	}

	if len(args) == 0 {
		return "", nil, errors.Errorf("chaincode definition for [%s] is invalid, policy field must be set", ccid)
	}

	return plugin, args, nil
}


func (v *dispatcherImpl) GetInfoForValidate(chdr *common.ChannelHeader, ccID string) (string, []byte, error) {
	
	plugin, args, err := v.getCDataForCC(chdr.ChannelId, ccID)
	if err != nil {
		msg := fmt.Sprintf("Unable to get chaincode data from ledger for txid %s, due to %s", chdr.TxId, err)
		logger.Errorf(msg)
		return "", nil, err
	}
	return plugin, args, nil
}



func (v *dispatcherImpl) txWritesToNamespace(ns *rwsetutil.NsRwSet) bool {
	
	if ns.KvRwSet != nil && len(ns.KvRwSet.Writes) > 0 {
		return true
	}

	
	for _, c := range ns.CollHashedRwSets {
		if c.HashedRwSet != nil && len(c.HashedRwSet.HashedWrites) > 0 {
			return true
		}

		
		if c.HashedRwSet != nil && len(c.HashedRwSet.MetadataWrites) > 0 {
			return true
		}
	}

	if ns.KvRwSet != nil && len(ns.KvRwSet.MetadataWrites) > 0 {
		return true
	}

	return false
}
