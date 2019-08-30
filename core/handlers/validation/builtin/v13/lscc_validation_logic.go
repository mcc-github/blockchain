/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package v13

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset/kvrwset"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/java"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/node"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	vc "github.com/mcc-github/blockchain/core/handlers/validation/api/capabilities"
	vs "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/core/handlers/validation/builtin/internal/car"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	"github.com/mcc-github/blockchain/internal/peer/packaging"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)



var systemChaincodeNames = map[string]struct{}{
	"cscc": {},
	"escc": {},
	"lscc": {},
	"qscc": {},
	"vscc": {},
}


func (vscc *Validator) checkInstantiationPolicy(chainName string, env *common.Envelope, instantiationPolicy []byte, payl *common.Payload) commonerrors.TxValidationError {
	
	shdr, err := protoutil.UnmarshalSignatureHeader(payl.Header.SignatureHeader)
	if err != nil {
		return policyErr(err)
	}

	
	sd := []*protoutil.SignedData{{
		Data:      env.Payload,
		Identity:  shdr.Creator,
		Signature: env.Signature,
	}}
	err = vscc.policyEvaluator.Evaluate(instantiationPolicy, sd)
	if err != nil {
		return policyErr(fmt.Errorf("chaincode instantiation policy violated, error %s", err))
	}
	return nil
}

func validateNewCollectionConfigs(newCollectionConfigs []*common.CollectionConfig) error {
	newCollectionsMap := make(map[string]bool, len(newCollectionConfigs))
	
	for _, newCollectionConfig := range newCollectionConfigs {

		newCollection := newCollectionConfig.GetStaticCollectionConfig()
		if newCollection == nil {
			return errors.New("unknown collection configuration type")
		}

		
		collectionName := newCollection.GetName()

		if err := validateCollectionName(collectionName); err != nil {
			return err
		}

		if _, ok := newCollectionsMap[collectionName]; !ok {
			newCollectionsMap[collectionName] = true
		} else {
			return fmt.Errorf("collection-name: %s -- found duplicate collection configuration", collectionName)
		}

		
		maximumPeerCount := newCollection.GetMaximumPeerCount()
		requiredPeerCount := newCollection.GetRequiredPeerCount()
		if maximumPeerCount < requiredPeerCount {
			return fmt.Errorf("collection-name: %s -- maximum peer count (%d) cannot be greater than the required peer count (%d)",
				collectionName, maximumPeerCount, requiredPeerCount)

		}
		if requiredPeerCount < 0 {
			return fmt.Errorf("collection-name: %s -- requiredPeerCount (%d) cannot be less than zero (%d)",
				collectionName, maximumPeerCount, requiredPeerCount)

		}

		
		err := validateSpOrConcat(newCollection.MemberOrgsPolicy.GetSignaturePolicy().Rule)
		if err != nil {
			return errors.WithMessagef(err, "collection-name: %s -- error in member org policy", collectionName)
		}
	}
	return nil
}


func validateSpOrConcat(sp *common.SignaturePolicy) error {
	if sp.GetNOutOf() == nil {
		return nil
	}
	
	if sp.GetNOutOf().N != 1 {
		return errors.New(fmt.Sprintf("signature policy is not an OR concatenation, NOutOf %d", sp.GetNOutOf().N))
	}
	
	for _, rule := range sp.GetNOutOf().Rules {
		err := validateSpOrConcat(rule)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkForMissingCollections(newCollectionsMap map[string]*common.StaticCollectionConfig, oldCollectionConfigs []*common.CollectionConfig,
) error {
	var missingCollections []string

	
	
	for _, oldCollectionConfig := range oldCollectionConfigs {

		oldCollection := oldCollectionConfig.GetStaticCollectionConfig()
		
		if oldCollection == nil {
			return policyErr(fmt.Errorf("unknown collection configuration type"))
		}

		
		oldCollectionName := oldCollection.GetName()
		_, ok := newCollectionsMap[oldCollectionName]
		if !ok {
			missingCollections = append(missingCollections, oldCollectionName)
		}
	}

	if len(missingCollections) > 0 {
		return policyErr(fmt.Errorf("the following existing collections are missing in the new collection configuration package: %v",
			missingCollections))
	}

	return nil
}

func checkForModifiedCollectionsBTL(newCollectionsMap map[string]*common.StaticCollectionConfig, oldCollectionConfigs []*common.CollectionConfig,
) error {
	var modifiedCollectionsBTL []string

	
	
	for _, oldCollectionConfig := range oldCollectionConfigs {

		oldCollection := oldCollectionConfig.GetStaticCollectionConfig()
		
		if oldCollection == nil {
			return policyErr(fmt.Errorf("unknown collection configuration type"))
		}

		oldCollectionName := oldCollection.GetName()
		newCollection, _ := newCollectionsMap[oldCollectionName]
		
		if newCollection.GetBlockToLive() != oldCollection.GetBlockToLive() {
			modifiedCollectionsBTL = append(modifiedCollectionsBTL, oldCollectionName)
		}
	}

	if len(modifiedCollectionsBTL) > 0 {
		return policyErr(fmt.Errorf("the BlockToLive in the following existing collections must not be modified: %v",
			modifiedCollectionsBTL))
	}

	return nil
}

func validateNewCollectionConfigsAgainstOld(newCollectionConfigs []*common.CollectionConfig, oldCollectionConfigs []*common.CollectionConfig,
) error {
	newCollectionsMap := make(map[string]*common.StaticCollectionConfig, len(newCollectionConfigs))

	for _, newCollectionConfig := range newCollectionConfigs {
		newCollection := newCollectionConfig.GetStaticCollectionConfig()
		
		
		newCollectionsMap[newCollection.GetName()] = newCollection
	}

	if err := checkForMissingCollections(newCollectionsMap, oldCollectionConfigs); err != nil {
		return err
	}

	if err := checkForModifiedCollectionsBTL(newCollectionsMap, oldCollectionConfigs); err != nil {
		return err
	}

	return nil
}

func validateCollectionName(collectionName string) error {
	if collectionName == "" {
		return fmt.Errorf("empty collection-name is not allowed")
	}
	match := validCollectionNameRegex.FindString(collectionName)
	if len(match) != len(collectionName) {
		return fmt.Errorf("collection-name: %s not allowed. A valid collection name follows the pattern: %s",
			collectionName, AllowedCharsCollectionName)
	}
	return nil
}




func (vscc *Validator) validateRWSetAndCollection(
	lsccrwset *kvrwset.KVRWSet,
	cdRWSet *ccprovider.ChaincodeData,
	lsccArgs [][]byte,
	lsccFunc string,
	ac vc.Capabilities,
	channelName string,
) commonerrors.TxValidationError {
	
	
	
	
	if len(lsccrwset.Writes) > 2 {
		return policyErr(fmt.Errorf("LSCC can only issue one or two putState upon deploy"))
	}

	
	
	
	var collectionsConfigArg []byte
	if len(lsccArgs) > 5 {
		collectionsConfigArg = lsccArgs[5]
	}

	var collectionsConfigLedger []byte
	if len(lsccrwset.Writes) == 2 {
		key := privdata.BuildCollectionKVSKey(cdRWSet.Name)
		if lsccrwset.Writes[1].Key != key {
			return policyErr(fmt.Errorf("invalid key for the collection of chaincode %s:%s; expected '%s', received '%s'",
				cdRWSet.Name, cdRWSet.Version, key, lsccrwset.Writes[1].Key))

		}

		collectionsConfigLedger = lsccrwset.Writes[1].Value
	}

	if !bytes.Equal(collectionsConfigArg, collectionsConfigLedger) {
		return policyErr(fmt.Errorf("collection configuration arguments supplied for chaincode %s:%s do not match the configuration in the lscc writeset",
			cdRWSet.Name, cdRWSet.Version))

	}

	channelState, err := vscc.stateFetcher.FetchState()
	if err != nil {
		return &commonerrors.VSCCExecutionFailureError{Err: fmt.Errorf("failed obtaining query executor: %v", err)}
	}
	defer channelState.Done()

	state := &state{channelState}

	
	
	if lsccFunc == lscc.DEPLOY {
		colCriteria := common.CollectionCriteria{Channel: channelName, Namespace: cdRWSet.Name}
		ccp, err := privdata.RetrieveCollectionConfigPackageFromState(colCriteria, state)
		if err != nil {
			
			
			
			if _, ok := err.(privdata.NoSuchCollectionError); !ok {
				return &commonerrors.VSCCExecutionFailureError{Err: fmt.Errorf("unable to check whether collection existed earlier for chaincode %s:%s",
					cdRWSet.Name, cdRWSet.Version),
				}
			}
		}
		if ccp != nil {
			return policyErr(fmt.Errorf("collection data should not exist for chaincode %s:%s", cdRWSet.Name, cdRWSet.Version))
		}
	}

	
	
	newCollectionConfigPackage := &common.CollectionConfigPackage{}

	if collectionsConfigArg != nil {
		err := proto.Unmarshal(collectionsConfigArg, newCollectionConfigPackage)
		if err != nil {
			return policyErr(fmt.Errorf("invalid collection configuration supplied for chaincode %s:%s",
				cdRWSet.Name, cdRWSet.Version))
		}
	} else {
		return nil
	}

	if ac.V1_2Validation() {
		newCollectionConfigs := newCollectionConfigPackage.GetConfig()
		if err := validateNewCollectionConfigs(newCollectionConfigs); err != nil {
			return policyErr(err)
		}

		if lsccFunc == lscc.UPGRADE {

			collectionCriteria := common.CollectionCriteria{Channel: channelName, Namespace: cdRWSet.Name}
			
			oldCollectionConfigPackage, err := privdata.RetrieveCollectionConfigPackageFromState(collectionCriteria, state)
			if err != nil {
				
				
				
				if _, ok := err.(privdata.NoSuchCollectionError); !ok {
					return &commonerrors.VSCCExecutionFailureError{Err: fmt.Errorf("unable to check whether collection existed earlier for chaincode %s:%s: %v",
						cdRWSet.Name, cdRWSet.Version, err),
					}
				}
			}

			
			if oldCollectionConfigPackage != nil {
				oldCollectionConfigs := oldCollectionConfigPackage.GetConfig()
				if err := validateNewCollectionConfigsAgainstOld(newCollectionConfigs, oldCollectionConfigs); err != nil {
					return policyErr(err)
				}

			}
		}
	}

	return nil
}

func (vscc *Validator) ValidateLSCCInvocation(
	chid string,
	env *common.Envelope,
	cap *pb.ChaincodeActionPayload,
	payl *common.Payload,
	ac vc.Capabilities,
) commonerrors.TxValidationError {
	cpp, err := protoutil.UnmarshalChaincodeProposalPayload(cap.ChaincodeProposalPayload)
	if err != nil {
		logger.Errorf("VSCC error: GetChaincodeProposalPayload failed, err %s", err)
		return policyErr(err)
	}

	cis := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(cpp.Input, cis)
	if err != nil {
		logger.Errorf("VSCC error: Unmarshal ChaincodeInvocationSpec failed, err %s", err)
		return policyErr(err)
	}

	if cis.ChaincodeSpec == nil ||
		cis.ChaincodeSpec.Input == nil ||
		cis.ChaincodeSpec.Input.Args == nil {
		logger.Errorf("VSCC error: committing invalid vscc invocation")
		return policyErr(fmt.Errorf("malformed chaincode invocation spec"))
	}

	lsccFunc := string(cis.ChaincodeSpec.Input.Args[0])
	lsccArgs := cis.ChaincodeSpec.Input.Args[1:]

	logger.Debugf("VSCC info: ValidateLSCCInvocation acting on %s %#v", lsccFunc, lsccArgs)

	switch lsccFunc {
	case lscc.UPGRADE, lscc.DEPLOY:
		logger.Debugf("VSCC info: validating invocation of lscc function %s on arguments %#v", lsccFunc, lsccArgs)

		if len(lsccArgs) < 2 {
			return policyErr(fmt.Errorf("Wrong number of arguments for invocation lscc(%s): expected at least 2, received %d", lsccFunc, len(lsccArgs)))
		}

		if (!ac.PrivateChannelData() && len(lsccArgs) > 5) ||
			(ac.PrivateChannelData() && len(lsccArgs) > 6) {
			return policyErr(fmt.Errorf("Wrong number of arguments for invocation lscc(%s): received %d", lsccFunc, len(lsccArgs)))
		}

		cdsArgs, err := protoutil.UnmarshalChaincodeDeploymentSpec(lsccArgs[1])
		if err != nil {
			return policyErr(fmt.Errorf("GetChaincodeDeploymentSpec error %s", err))
		}

		if cdsArgs == nil || cdsArgs.ChaincodeSpec == nil || cdsArgs.ChaincodeSpec.ChaincodeId == nil ||
			cap.Action == nil || cap.Action.ProposalResponsePayload == nil {
			return policyErr(fmt.Errorf("VSCC error: invocation of lscc(%s) does not have appropriate arguments", lsccFunc))
		}

		err = packaging.NewRegistry(
			
			
			
			
			&golang.Platform{},
			&node.Platform{},
			&java.Platform{},
			&car.Platform{},
		).ValidateDeploymentSpec(cdsArgs.ChaincodeSpec.Type.String(), cdsArgs.CodePackage)
		if err != nil {
			return policyErr(fmt.Errorf("failed to validate deployment spec: %s", err))
		}

		
		ccName := cdsArgs.ChaincodeSpec.ChaincodeId.Name
		
		if !lscc.ChaincodeNameRegExp.MatchString(ccName) {
			return policyErr(errors.Errorf("invalid chaincode name '%s'", ccName))
		}

		
		if _, in := systemChaincodeNames[ccName]; in {
			return policyErr(errors.Errorf("chaincode name '%s' is reserved for system chaincodes", ccName))
		}

		
		ccVersion := cdsArgs.ChaincodeSpec.ChaincodeId.Version
		
		if !lscc.ChaincodeVersionRegExp.MatchString(ccVersion) {
			return policyErr(errors.Errorf("invalid chaincode version '%s'", ccVersion))
		}

		
		pRespPayload, err := protoutil.UnmarshalProposalResponsePayload(cap.Action.ProposalResponsePayload)
		if err != nil {
			return policyErr(fmt.Errorf("GetProposalResponsePayload error %s", err))
		}
		if pRespPayload.Extension == nil {
			return policyErr(fmt.Errorf("nil pRespPayload.Extension"))
		}
		respPayload, err := protoutil.UnmarshalChaincodeAction(pRespPayload.Extension)
		if err != nil {
			return policyErr(fmt.Errorf("GetChaincodeAction error %s", err))
		}
		txRWSet := &rwsetutil.TxRwSet{}
		if err = txRWSet.FromProtoBytes(respPayload.Results); err != nil {
			return policyErr(fmt.Errorf("txRWSet.FromProtoBytes error %s", err))
		}

		
		var lsccrwset *kvrwset.KVRWSet
		for _, ns := range txRWSet.NsRwSets {
			logger.Debugf("Namespace %s", ns.NameSpace)
			if ns.NameSpace == "lscc" {
				lsccrwset = ns.KvRwSet
				break
			}
		}

		
		cdLedger, ccExistsOnLedger, err := vscc.getInstantiatedCC(chid, cdsArgs.ChaincodeSpec.ChaincodeId.Name)
		if err != nil {
			return &commonerrors.VSCCExecutionFailureError{Err: err}
		}

		
		
		
		
		if lsccrwset == nil {
			return policyErr(fmt.Errorf("No read write set for lscc was found"))
		}
		
		if len(lsccrwset.Writes) < 1 {
			return policyErr(fmt.Errorf("LSCC must issue at least one single putState upon deploy/upgrade"))
		}
		
		if lsccrwset.Writes[0].Key != cdsArgs.ChaincodeSpec.ChaincodeId.Name {
			return policyErr(fmt.Errorf("expected key %s, found %s", cdsArgs.ChaincodeSpec.ChaincodeId.Name, lsccrwset.Writes[0].Key))
		}
		
		cdRWSet := &ccprovider.ChaincodeData{}
		err = proto.Unmarshal(lsccrwset.Writes[0].Value, cdRWSet)
		if err != nil {
			return policyErr(fmt.Errorf("unmarhsalling of ChaincodeData failed, error %s", err))
		}
		
		if cdRWSet.Name != cdsArgs.ChaincodeSpec.ChaincodeId.Name {
			return policyErr(fmt.Errorf("expected cc name %s, found %s", cdsArgs.ChaincodeSpec.ChaincodeId.Name, cdRWSet.Name))
		}
		
		if cdRWSet.Version != cdsArgs.ChaincodeSpec.ChaincodeId.Version {
			return policyErr(fmt.Errorf("expected cc version %s, found %s", cdsArgs.ChaincodeSpec.ChaincodeId.Version, cdRWSet.Version))
		}
		
		for _, ns := range txRWSet.NsRwSets {
			if ns.NameSpace != "lscc" && ns.NameSpace != cdRWSet.Name && len(ns.KvRwSet.Writes) > 0 {
				return policyErr(fmt.Errorf("LSCC invocation is attempting to write to namespace %s", ns.NameSpace))
			}
		}

		logger.Debugf("Validating %s for cc %s version %s", lsccFunc, cdRWSet.Name, cdRWSet.Version)

		switch lsccFunc {
		case lscc.DEPLOY:

			
			
			
			if ccExistsOnLedger {
				return policyErr(fmt.Errorf("Chaincode %s is already instantiated", cdsArgs.ChaincodeSpec.ChaincodeId.Name))
			}

			
			
			
			if ac.PrivateChannelData() {
				
				err := vscc.validateRWSetAndCollection(lsccrwset, cdRWSet, lsccArgs, lsccFunc, ac, chid)
				if err != nil {
					return err
				}
			} else {
				
				if len(lsccrwset.Writes) != 1 {
					return policyErr(fmt.Errorf("LSCC can only issue a single putState upon deploy"))
				}
			}

			
			
			
			pol := cdRWSet.InstantiationPolicy
			if pol == nil {
				return policyErr(fmt.Errorf("no instantiation policy was specified"))
			}
			
			
			
			
			
			err := vscc.checkInstantiationPolicy(chid, env, pol, payl)
			if err != nil {
				return err
			}

		case lscc.UPGRADE:
			
			
			
			if !ccExistsOnLedger {
				return policyErr(fmt.Errorf("Upgrading non-existent chaincode %s", cdsArgs.ChaincodeSpec.ChaincodeId.Name))
			}

			
			
			
			if cdLedger.Version == cdsArgs.ChaincodeSpec.ChaincodeId.Version {
				return policyErr(fmt.Errorf("Existing version of the cc on the ledger (%s) should be different from the upgraded one", cdsArgs.ChaincodeSpec.ChaincodeId.Version))
			}

			
			
			
			
			if ac.V1_2Validation() {
				
				err := vscc.validateRWSetAndCollection(lsccrwset, cdRWSet, lsccArgs, lsccFunc, ac, chid)
				if err != nil {
					return err
				}
			} else {
				
				if len(lsccrwset.Writes) != 1 {
					return policyErr(fmt.Errorf("LSCC can only issue a single putState upon upgrade"))
				}
			}

			
			
			
			pol := cdLedger.InstantiationPolicy
			if pol == nil {
				return policyErr(fmt.Errorf("No instantiation policy was specified"))
			}
			
			
			
			
			
			err := vscc.checkInstantiationPolicy(chid, env, pol, payl)
			if err != nil {
				return err
			}

			
			
			
			if ac.V1_1Validation() {
				polNew := cdRWSet.InstantiationPolicy
				if polNew == nil {
					return policyErr(fmt.Errorf("No instantiation policy was specified"))
				}

				
				if !bytes.Equal(polNew, pol) {
					err = vscc.checkInstantiationPolicy(chid, env, polNew, payl)
					if err != nil {
						return err
					}
				}
			}
		}

		
		return nil
	default:
		return policyErr(fmt.Errorf("VSCC error: committing an invocation of function %s of lscc is invalid", lsccFunc))
	}
}

func (vscc *Validator) getInstantiatedCC(chid, ccid string) (cd *ccprovider.ChaincodeData, exists bool, err error) {
	qe, err := vscc.stateFetcher.FetchState()
	if err != nil {
		err = fmt.Errorf("could not retrieve QueryExecutor for channel %s, error %s", chid, err)
		return
	}
	defer qe.Done()
	channelState := &state{qe}
	bytes, err := channelState.GetState("lscc", ccid)
	if err != nil {
		err = fmt.Errorf("could not retrieve state for chaincode %s on channel %s, error %s", ccid, chid, err)
		return
	}

	if bytes == nil {
		return
	}

	cd = &ccprovider.ChaincodeData{}
	err = proto.Unmarshal(bytes, cd)
	if err != nil {
		err = fmt.Errorf("unmarshalling ChaincodeQueryResponse failed, error %s", err)
		return
	}

	exists = true
	return
}

type state struct {
	vs.State
}


func (s *state) GetState(namespace string, key string) ([]byte, error) {
	values, err := s.GetStateMultipleKeys(namespace, []string{key})
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}
	return values[0], nil
}
