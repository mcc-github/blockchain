/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)



type VsccValidatorImpl struct {
	chainID         string
	cr              ChannelResources
	pluginValidator *PluginValidator
}


func newVSCCValidator(chainID string, cr ChannelResources, pluginValidator *PluginValidator) *VsccValidatorImpl {
	return &VsccValidatorImpl{
		chainID:         chainID,
		cr:              cr,
		pluginValidator: pluginValidator,
	}
}

func getChaincodeHeaderExtension(hdr *common.Header) (*peer.ChaincodeHeaderExtension, error) {
	chdr, err := protoutil.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return nil, err
	}

	chaincodeHdrExt := &peer.ChaincodeHeaderExtension{}
	err = proto.Unmarshal(chdr.Extension, chaincodeHdrExt)
	return chaincodeHdrExt, errors.Wrap(err, "error unmarshaling ChaincodeHeaderExtension")
}


func (v *VsccValidatorImpl) VSCCValidateTx(seq int, payload *common.Payload, envBytes []byte, block *common.Block) (error, peer.TxValidationCode) {
	chainID := v.chainID
	logger.Debugf("[%s] VSCCValidateTx starts for bytes %p", chainID, envBytes)

	
	hdrExt, err := getChaincodeHeaderExtension(payload.Header)
	if err != nil {
		return err, peer.TxValidationCode_BAD_HEADER_EXTENSION
	}

	
	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return err, peer.TxValidationCode_BAD_CHANNEL_HEADER
	}

	
	writesToLSCC := false
	writesToNonInvokableSCC := false
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
		return err, peer.TxValidationCode_INVALID_OTHER_REASON
	}
	if ccID != respPayload.ChaincodeId.Name {
		err = errors.Errorf("inconsistent ccid info (%s/%s)", ccID, respPayload.ChaincodeId.Name)
		logger.Errorf("%+v", err)
		return err, peer.TxValidationCode_INVALID_OTHER_REASON
	}
	
	if ccVer == "" {
		err = errors.New("invalid chaincode version")
		logger.Errorf("%+v", err)
		return err, peer.TxValidationCode_INVALID_OTHER_REASON
	}

	var wrNamespace []string
	alwaysEnforceOriginalNamespace := v.cr.Capabilities().V1_2Validation()
	if alwaysEnforceOriginalNamespace {
		wrNamespace = append(wrNamespace, ccID)
		if respPayload.Events != nil {
			ccEvent := &peer.ChaincodeEvent{}
			if err = proto.Unmarshal(respPayload.Events, ccEvent); err != nil {
				return errors.Wrapf(err, "invalid chaincode event"), peer.TxValidationCode_INVALID_OTHER_REASON
			}
			if ccEvent.ChaincodeId != ccID {
				return errors.Errorf("chaincode event chaincode id does not match chaincode action chaincode id"), peer.TxValidationCode_INVALID_OTHER_REASON
			}
		}
	}

	namespaces := make(map[string]struct{})
	for _, ns := range txRWSet.NsRwSets {
		
		if _, ok := namespaces[ns.NameSpace]; ok {
			return errors.Errorf("duplicate namespace '%s' in txRWSet", ns.NameSpace),
				peer.TxValidationCode_ILLEGAL_WRITESET
		}
		namespaces[ns.NameSpace] = struct{}{}

		if !v.txWritesToNamespace(ns) {
			continue
		}

		
		
		if ns.NameSpace != ccID || !alwaysEnforceOriginalNamespace {
			wrNamespace = append(wrNamespace, ns.NameSpace)
		}

		if !writesToLSCC && ns.NameSpace == "lscc" {
			writesToLSCC = true
		}

		if !writesToNonInvokableSCC && IsSysCCAndNotInvokableCC2CC(ns.NameSpace) {
			writesToNonInvokableSCC = true
		}

		if !writesToNonInvokableSCC && IsSysCCAndNotInvokableExternal(ns.NameSpace) {
			writesToNonInvokableSCC = true
		}
	}

	
	
	

	if !IsSysCC(ccID) {
		
		
		
		
		
		
		
		if writesToLSCC {
			return errors.Errorf("chaincode %s attempted to write to the namespace of LSCC", ccID),
				peer.TxValidationCode_ILLEGAL_WRITESET
		}
		
		
		
		
		
		
		
		if writesToNonInvokableSCC {
			return errors.Errorf("chaincode %s attempted to write to the namespace of a system chaincode that cannot be invoked", ccID),
				peer.TxValidationCode_ILLEGAL_WRITESET
		}

		
		for _, ns := range wrNamespace {
			
			txcc, vscc, policy, err := v.GetInfoForValidate(chdr, ns)
			if err != nil {
				logger.Errorf("GetInfoForValidate for txId = %s returned error: %+v", chdr.TxId, err)
				return err, peer.TxValidationCode_INVALID_OTHER_REASON
			}

			
			
			
			if ns == ccID && txcc.ChaincodeVersion != ccVer {
				err = errors.Errorf("chaincode %s:%s/%s didn't match %s:%s/%s in lscc", ccID, ccVer, chdr.ChannelId, txcc.ChaincodeName, txcc.ChaincodeVersion, chdr.ChannelId)
				logger.Errorf("%+v", err)
				return err, peer.TxValidationCode_EXPIRED_CHAINCODE
			}

			
			ctx := &Context{
				Seq:       seq,
				Envelope:  envBytes,
				Block:     block,
				TxID:      chdr.TxId,
				Channel:   chdr.ChannelId,
				Namespace: ns,
				Policy:    policy,
				VSCCName:  vscc.ChaincodeName,
			}
			if err = v.VSCCValidateTxForCC(ctx); err != nil {
				switch err.(type) {
				case *commonerrors.VSCCEndorsementPolicyError:
					return err, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE
				default:
					return err, peer.TxValidationCode_INVALID_OTHER_REASON
				}
			}
		}
	} else {
		
		
		
		
		if IsSysCCAndNotInvokableExternal(ccID) {
			return errors.Errorf("committing an invocation of cc %s is illegal", ccID),
				peer.TxValidationCode_ILLEGAL_WRITESET
		}

		
		_, vscc, policy, err := v.GetInfoForValidate(chdr, ccID)
		if err != nil {
			logger.Errorf("GetInfoForValidate for txId = %s returned error: %+v", chdr.TxId, err)
			return err, peer.TxValidationCode_INVALID_OTHER_REASON
		}

		
		
		
		
		
		ctx := &Context{
			Seq:       seq,
			Envelope:  envBytes,
			Block:     block,
			TxID:      chdr.TxId,
			Channel:   chdr.ChannelId,
			Namespace: ccID,
			Policy:    policy,
			VSCCName:  vscc.ChaincodeName,
		}
		if err = v.VSCCValidateTxForCC(ctx); err != nil {
			switch err.(type) {
			case *commonerrors.VSCCEndorsementPolicyError:
				return err, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE
			default:
				return err, peer.TxValidationCode_INVALID_OTHER_REASON
			}
		}
	}
	logger.Debugf("[%s] VSCCValidateTx completes env bytes %p", chainID, envBytes)
	return nil, peer.TxValidationCode_VALID
}

func (v *VsccValidatorImpl) VSCCValidateTxForCC(ctx *Context) error {
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

func (v *VsccValidatorImpl) getCDataForCC(chid, ccid string) (*ccprovider.ChaincodeData, error) {
	l := v.cr.Ledger()
	if l == nil {
		return nil, errors.New("nil ledger instance")
	}

	qe, err := l.NewQueryExecutor()
	if err != nil {
		return nil, errors.WithMessage(err, "could not retrieve QueryExecutor")
	}
	defer qe.Done()

	bytes, err := qe.GetState("lscc", ccid)
	if err != nil {
		return nil, &commonerrors.VSCCInfoLookupFailureError{
			Reason: fmt.Sprintf("Could not retrieve state for chaincode %s, error %s", ccid, err),
		}
	}

	if bytes == nil {
		return nil, errors.Errorf("lscc's state for [%s] not found.", ccid)
	}

	cd := &ccprovider.ChaincodeData{}
	err = proto.Unmarshal(bytes, cd)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling ChaincodeQueryResponse failed")
	}

	if cd.Vscc == "" {
		return nil, errors.Errorf("lscc's state for [%s] is invalid, vscc field must be set", ccid)
	}

	if len(cd.Policy) == 0 {
		return nil, errors.Errorf("lscc's state for [%s] is invalid, policy field must be set", ccid)
	}

	return cd, err
}


func (v *VsccValidatorImpl) GetInfoForValidate(chdr *common.ChannelHeader, ccID string) (*sysccprovider.ChaincodeInstance, *sysccprovider.ChaincodeInstance, []byte, error) {
	cc := &sysccprovider.ChaincodeInstance{
		ChainID:       chdr.ChannelId,
		ChaincodeName: ccID,
	}
	vscc := &sysccprovider.ChaincodeInstance{
		ChainID:       chdr.ChannelId,
		ChaincodeName: "vscc", 
	}
	var policy []byte
	var err error
	if !IsSysCC(ccID) {
		
		
		

		
		cd, err := v.getCDataForCC(chdr.ChannelId, ccID)
		if err != nil {
			msg := fmt.Sprintf("Unable to get chaincode data from ledger for txid %s, due to %s", chdr.TxId, err)
			logger.Errorf(msg)
			return nil, nil, nil, err
		}
		cc.ChaincodeName = cd.Name
		cc.ChaincodeVersion = cd.Version
		vscc.ChaincodeName, policy = cd.Vscc, cd.Policy
	} else {
		
		
		
		p := cauthdsl.SignedByAnyMember(v.cr.GetMSPIDs())
		policy, err = protoutil.Marshal(p)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return cc, vscc, policy, nil
}



func (v *VsccValidatorImpl) txWritesToNamespace(ns *rwsetutil.NsRwSet) bool {
	
	if ns.KvRwSet != nil && len(ns.KvRwSet.Writes) > 0 {
		return true
	}

	
	if v.cr.Capabilities().PrivateChannelData() {
		
		for _, c := range ns.CollHashedRwSets {
			if c.HashedRwSet != nil && len(c.HashedRwSet.HashedWrites) > 0 {
				return true
			}

			
			if v.cr.Capabilities().KeyLevelEndorsement() {
				
				if c.HashedRwSet != nil && len(c.HashedRwSet.MetadataWrites) > 0 {
					return true
				}
			}
		}
	}

	
	if v.cr.Capabilities().KeyLevelEndorsement() {
		
		if ns.KvRwSet != nil && len(ns.KvRwSet.MetadataWrites) > 0 {
			return true
		}
	}

	return false
}

func IsSysCCAndNotInvokableExternal(name string) bool {
	return name == "vscc" || name == "escc"
}

func IsSysCC(name string) bool {
	return name == "vscc" || name == "escc" || name == "lscc" || name == "qscc" || name == "cscc"
}

func IsSysCCAndNotInvokableCC2CC(name string) bool {
	return name == "vscc" || name == "escc" || name == "cscc"
}
