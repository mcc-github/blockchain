/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"fmt"
	"regexp"

	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	p "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	persistenceintf "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/dispatcher"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	cb "github.com/mcc-github/blockchain/protos/common"
	mspprotos "github.com/mcc-github/blockchain/protos/msp"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

const (
	
	
	LifecycleNamespace = "_lifecycle"

	
	
	InstallChaincodeFuncName = "InstallChaincode"

	
	
	QueryInstalledChaincodeFuncName = "QueryInstalledChaincode"

	
	
	QueryInstalledChaincodesFuncName = "QueryInstalledChaincodes"

	
	
	ApproveChaincodeDefinitionForMyOrgFuncName = "ApproveChaincodeDefinitionForMyOrg"

	
	
	
	CheckCommitReadinessFuncName = "CheckCommitReadiness"

	
	
	CommitChaincodeDefinitionFuncName = "CommitChaincodeDefinition"

	
	
	QueryChaincodeDefinitionFuncName = "QueryChaincodeDefinition"

	
	
	QueryChaincodeDefinitionsFuncName = "QueryChaincodeDefinitions"
)



type SCCFunctions interface {
	
	InstallChaincode([]byte) (*chaincode.InstalledChaincode, error)

	
	QueryInstalledChaincode(packageID persistenceintf.PackageID) (*chaincode.InstalledChaincode, error)

	
	
	GetInstalledChaincodePackage(packageID p.PackageID) ([]byte, error)

	
	QueryInstalledChaincodes() []*chaincode.InstalledChaincode

	
	ApproveChaincodeDefinitionForOrg(chname, ccname string, cd *ChaincodeDefinition, packageID persistenceintf.PackageID, publicState ReadableState, orgState ReadWritableState) error

	
	
	
	CheckCommitReadiness(chname, ccname string, cd *ChaincodeDefinition, publicState ReadWritableState, orgStates []OpaqueState) (map[string]bool, error)

	
	
	
	CommitChaincodeDefinition(chname, ccname string, cd *ChaincodeDefinition, publicState ReadWritableState, orgStates []OpaqueState) (map[string]bool, error)

	
	
	QueryChaincodeDefinition(name string, publicState ReadableState) (*ChaincodeDefinition, error)

	
	
	
	QueryOrgApprovals(name string, cd *ChaincodeDefinition, orgStates []OpaqueState) (map[string]bool, error)

	
	QueryNamespaceDefinitions(publicState RangeableState) (map[string]string, error)
}





type ChannelConfigSource interface {
	
	
	
	GetStableChannelConfig(channelID string) channelconfig.Resources
}




type QueryExecutorProvider interface {
	TxQueryExecutor(channelID, txID string) ledger.SimpleQueryExecutor
}



type SCC struct {
	OrgMSPID string

	ACLProvider aclmgmt.ACLProvider

	ChannelConfigSource ChannelConfigSource

	DeployedCCInfoProvider ledger.DeployedChaincodeInfoProvider
	QueryExecutorProvider  QueryExecutorProvider

	
	Functions SCCFunctions

	
	
	Dispatcher *dispatcher.Dispatcher
}


func (scc *SCC) Name() string {
	return LifecycleNamespace
}


func (scc *SCC) Path() string {
	return "github.com/mcc-github/blockchain/core/chaincode/lifecycle"
}


func (scc *SCC) InitArgs() [][]byte {
	return nil
}


func (scc *SCC) Chaincode() shim.Chaincode {
	return scc
}


func (scc *SCC) InvokableExternal() bool {
	return true
}


func (scc *SCC) InvokableCC2CC() bool {
	return true
}


func (scc *SCC) Enabled() bool {
	return true
}


func (scc *SCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}




func (scc *SCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) == 0 {
		return shim.Error("lifecycle scc must be invoked with arguments")
	}

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("lifecycle scc operations require exactly two arguments but received %d", len(args)))
	}

	var ac channelconfig.Application
	var channelID string
	if channelID = stub.GetChannelID(); channelID != "" {
		channelConfig := scc.ChannelConfigSource.GetStableChannelConfig(channelID)
		if channelConfig == nil {
			return shim.Error(fmt.Sprintf("could not get channelconfig for channel '%s'", channelID))
		}
		var ok bool
		ac, ok = channelConfig.ApplicationConfig()
		if !ok {
			return shim.Error(fmt.Sprintf("could not get application config for channel '%s'", channelID))
		}
		if !ac.Capabilities().LifecycleV20() {
			return shim.Error(fmt.Sprintf("cannot use new lifecycle for channel '%s' as it does not have the required capabilities enabled", channelID))
		}
	}

	
	sp, err := stub.GetSignedProposal()
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed getting signed proposal from stub: [%s]", err))
	}

	err = scc.ACLProvider.CheckACL(fmt.Sprintf("%s/%s", LifecycleNamespace, args[0]), stub.GetChannelID(), sp)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to authorize invocation due to failed ACL check: %s", err))
	}

	outputBytes, err := scc.Dispatcher.Dispatch(
		args[1],
		string(args[0]),
		&Invocation{
			ChannelID:         channelID,
			ApplicationConfig: ac,
			SCC:               scc,
			Stub:              stub,
		},
	)
	if err != nil {
		switch err.(type) {
		case ErrNamespaceNotDefined, persistence.CodePackageNotFoundErr:
			return pb.Response{
				Status:  404,
				Message: err.Error(),
			}
		default:
			return shim.Error(fmt.Sprintf("failed to invoke backing implementation of '%s': %s", string(args[0]), err.Error()))
		}
	}

	return shim.Success(outputBytes)
}

type Invocation struct {
	ChannelID         string
	ApplicationConfig channelconfig.Application 
	Stub              shim.ChaincodeStubInterface
	SCC               *SCC
}



func (i *Invocation) InstallChaincode(input *lb.InstallChaincodeArgs) (proto.Message, error) {
	if logger.IsEnabledFor(zapcore.DebugLevel) {
		end := 35
		if len(input.ChaincodeInstallPackage) < end {
			end = len(input.ChaincodeInstallPackage)
		}

		
		
		
		
		packageFingerprint := input.ChaincodeInstallPackage[0:end]
		logger.Debugf("received invocation of InstallChaincode for install package %x...",
			packageFingerprint,
		)
	}

	installedCC, err := i.SCC.Functions.InstallChaincode(input.ChaincodeInstallPackage)
	if err != nil {
		return nil, err
	}

	return &lb.InstallChaincodeResult{
		Label:     installedCC.Label,
		PackageId: installedCC.PackageID.String(),
	}, nil
}



func (i *Invocation) QueryInstalledChaincode(input *lb.QueryInstalledChaincodeArgs) (proto.Message, error) {
	logger.Debugf("received invocation of QueryInstalledChaincode for install package ID '%s'",
		input.PackageId,
	)

	chaincode, err := i.SCC.Functions.QueryInstalledChaincode(persistenceintf.PackageID(input.PackageId))
	if err != nil {
		return nil, err
	}
	return &lb.QueryInstalledChaincodeResult{
		Label:     chaincode.Label,
		PackageId: chaincode.PackageID.String(),
	}, nil
}



func (i *Invocation) GetInstalledChaincodePackage(input *lb.GetInstalledChaincodePackageArgs) (proto.Message, error) {
	logger.Debugf("received invocation of GetInstalledChaincodePackage")

	pkgBytes, err := i.SCC.Functions.GetInstalledChaincodePackage(p.PackageID(input.PackageId))
	if err != nil {
		return nil, err
	}

	return &lb.GetInstalledChaincodePackageResult{
		ChaincodeInstallPackage: pkgBytes,
	}, nil
}



func (i *Invocation) QueryInstalledChaincodes(input *lb.QueryInstalledChaincodesArgs) (proto.Message, error) {
	logger.Debugf("received invocation of QueryInstalledChaincodes")

	chaincodes := i.SCC.Functions.QueryInstalledChaincodes()

	result := &lb.QueryInstalledChaincodesResult{}
	for _, chaincode := range chaincodes {
		references := map[string]*lb.QueryInstalledChaincodesResult_References{}
		for channel, chaincodeMetadata := range chaincode.References {
			chaincodes := make([]*lb.QueryInstalledChaincodesResult_Chaincode, len(chaincodeMetadata))
			for i, metadata := range chaincodeMetadata {
				chaincodes[i] = &lb.QueryInstalledChaincodesResult_Chaincode{
					Name:    metadata.Name,
					Version: metadata.Version,
				}
			}

			references[channel] = &lb.QueryInstalledChaincodesResult_References{
				Chaincodes: chaincodes,
			}
		}

		result.InstalledChaincodes = append(result.InstalledChaincodes,
			&lb.QueryInstalledChaincodesResult_InstalledChaincode{
				Label:      chaincode.Label,
				PackageId:  chaincode.PackageID.String(),
				References: references,
			})
	}

	return result, nil
}



func (i *Invocation) ApproveChaincodeDefinitionForMyOrg(input *lb.ApproveChaincodeDefinitionForMyOrgArgs) (proto.Message, error) {
	if err := i.validateInput(input.Name, input.Version, input.Collections); err != nil {
		return nil, err
	}
	collectionName := ImplicitCollectionNameForOrg(i.SCC.OrgMSPID)
	var collectionConfig []*cb.CollectionConfig
	if input.Collections != nil {
		collectionConfig = input.Collections.Config
	}

	var packageID persistenceintf.PackageID
	if input.Source != nil {
		switch source := input.Source.Type.(type) {
		case *lb.ChaincodeSource_LocalPackage:
			packageID = persistenceintf.PackageID(source.LocalPackage.PackageId)
		case *lb.ChaincodeSource_Unavailable_:
		default:
		}
	}

	cd := &ChaincodeDefinition{
		Sequence: input.Sequence,
		EndorsementInfo: &lb.ChaincodeEndorsementInfo{
			Version:           input.Version,
			EndorsementPlugin: input.EndorsementPlugin,
			InitRequired:      input.InitRequired,
		},
		ValidationInfo: &lb.ChaincodeValidationInfo{
			ValidationPlugin:    input.ValidationPlugin,
			ValidationParameter: input.ValidationParameter,
		},
		Collections: &cb.CollectionConfigPackage{
			Config: collectionConfig,
		},
	}

	logger.Debugf("received invocation of ApproveChaincodeDefinitionForMyOrg on channel '%s' for definition '%s'",
		i.Stub.GetChannelID(),
		cd,
	)

	if err := i.SCC.Functions.ApproveChaincodeDefinitionForOrg(
		i.Stub.GetChannelID(),
		input.Name,
		cd,
		packageID,
		i.Stub,
		&ChaincodePrivateLedgerShim{
			Collection: collectionName,
			Stub:       i.Stub,
		},
	); err != nil {
		return nil, err
	}
	return &lb.ApproveChaincodeDefinitionForMyOrgResult{}, nil
}



func (i *Invocation) CheckCommitReadiness(input *lb.CheckCommitReadinessArgs) (proto.Message, error) {
	opaqueStates, err := i.createOpaqueStates()
	if err != nil {
		return nil, err
	}

	cd := &ChaincodeDefinition{
		Sequence: input.Sequence,
		EndorsementInfo: &lb.ChaincodeEndorsementInfo{
			Version:           input.Version,
			EndorsementPlugin: input.EndorsementPlugin,
			InitRequired:      input.InitRequired,
		},
		ValidationInfo: &lb.ChaincodeValidationInfo{
			ValidationPlugin:    input.ValidationPlugin,
			ValidationParameter: input.ValidationParameter,
		},
		Collections: input.Collections,
	}

	logger.Debugf("received invocation of CheckCommitReadiness on channel '%s' for definition '%s'",
		i.Stub.GetChannelID(),
		cd,
	)

	approvals, err := i.SCC.Functions.CheckCommitReadiness(
		i.Stub.GetChannelID(),
		input.Name,
		cd,
		i.Stub,
		opaqueStates,
	)
	if err != nil {
		return nil, err
	}

	return &lb.CheckCommitReadinessResult{
		Approvals: approvals,
	}, nil
}



func (i *Invocation) CommitChaincodeDefinition(input *lb.CommitChaincodeDefinitionArgs) (proto.Message, error) {
	if err := i.validateInput(input.Name, input.Version, input.Collections); err != nil {
		return nil, err
	}

	if i.ApplicationConfig == nil {
		return nil, errors.Errorf("no application config for channel '%s'", i.Stub.GetChannelID())
	}

	orgs := i.ApplicationConfig.Organizations()
	opaqueStates := make([]OpaqueState, 0, len(orgs))
	var myOrg string
	for _, org := range orgs {
		opaqueStates = append(opaqueStates, &ChaincodePrivateLedgerShim{
			Collection: ImplicitCollectionNameForOrg(org.MSPID()),
			Stub:       i.Stub,
		})
		if org.MSPID() == i.SCC.OrgMSPID {
			myOrg = i.SCC.OrgMSPID
		}
	}

	if myOrg == "" {
		return nil, errors.Errorf("impossibly, this peer's org is processing requests for a channel it is not a member of")
	}

	cd := &ChaincodeDefinition{
		Sequence: input.Sequence,
		EndorsementInfo: &lb.ChaincodeEndorsementInfo{
			Version:           input.Version,
			EndorsementPlugin: input.EndorsementPlugin,
			InitRequired:      input.InitRequired,
		},
		ValidationInfo: &lb.ChaincodeValidationInfo{
			ValidationPlugin:    input.ValidationPlugin,
			ValidationParameter: input.ValidationParameter,
		},
		Collections: input.Collections,
	}

	logger.Debugf("received invocation of CommitChaincodeDefinition on channel '%s' for definition '%s'",
		i.Stub.GetChannelID(),
		cd,
	)

	approvals, err := i.SCC.Functions.CommitChaincodeDefinition(
		i.Stub.GetChannelID(),
		input.Name,
		cd,
		i.Stub,
		opaqueStates,
	)
	if err != nil {
		return nil, err
	}

	if !approvals[myOrg] {
		return nil, errors.Errorf("chaincode definition not agreed to by this org (%s)", i.SCC.OrgMSPID)
	}

	return &lb.CommitChaincodeDefinitionResult{}, nil
}



func (i *Invocation) QueryChaincodeDefinition(input *lb.QueryChaincodeDefinitionArgs) (proto.Message, error) {
	logger.Debugf("received invocation of QueryChaincodeDefinition on channel '%s' for chaincode '%s'",
		i.Stub.GetChannelID(),
		input.Name,
	)

	definedChaincode, err := i.SCC.Functions.QueryChaincodeDefinition(input.Name, i.Stub)
	if err != nil {
		return nil, err
	}

	opaqueStates, err := i.createOpaqueStates()
	if err != nil {
		return nil, err
	}

	var approvals map[string]bool
	if approvals, err = i.SCC.Functions.QueryOrgApprovals(input.Name, definedChaincode, opaqueStates); err != nil {
		return nil, err
	}

	return &lb.QueryChaincodeDefinitionResult{
		Sequence:            definedChaincode.Sequence,
		Version:             definedChaincode.EndorsementInfo.Version,
		EndorsementPlugin:   definedChaincode.EndorsementInfo.EndorsementPlugin,
		ValidationPlugin:    definedChaincode.ValidationInfo.ValidationPlugin,
		ValidationParameter: definedChaincode.ValidationInfo.ValidationParameter,
		InitRequired:        definedChaincode.EndorsementInfo.InitRequired,
		Collections:         definedChaincode.Collections,
		Approvals:           approvals,
	}, nil
}



func (i *Invocation) QueryChaincodeDefinitions(input *lb.QueryChaincodeDefinitionsArgs) (proto.Message, error) {
	logger.Debugf("received invocation of QueryChaincodeDefinitions on channel '%s'",
		i.Stub.GetChannelID(),
	)

	namespaces, err := i.SCC.Functions.QueryNamespaceDefinitions(&ChaincodePublicLedgerShim{ChaincodeStubInterface: i.Stub})
	if err != nil {
		return nil, err
	}

	chaincodeDefinitions := []*lb.QueryChaincodeDefinitionsResult_ChaincodeDefinition{}
	for namespace, nType := range namespaces {
		if nType == FriendlyChaincodeDefinitionType {
			definedChaincode, err := i.SCC.Functions.QueryChaincodeDefinition(namespace, i.Stub)
			if err != nil {
				return nil, err
			}

			chaincodeDefinitions = append(chaincodeDefinitions, &lb.QueryChaincodeDefinitionsResult_ChaincodeDefinition{
				Name:                namespace,
				Sequence:            definedChaincode.Sequence,
				Version:             definedChaincode.EndorsementInfo.Version,
				EndorsementPlugin:   definedChaincode.EndorsementInfo.EndorsementPlugin,
				ValidationPlugin:    definedChaincode.ValidationInfo.ValidationPlugin,
				ValidationParameter: definedChaincode.ValidationInfo.ValidationParameter,
				InitRequired:        definedChaincode.EndorsementInfo.InitRequired,
				Collections:         definedChaincode.Collections,
			})
		}
	}

	return &lb.QueryChaincodeDefinitionsResult{
		ChaincodeDefinitions: chaincodeDefinitions,
	}, nil
}

var (
	
	
	ChaincodeNameRegExp    = regexp.MustCompile("^[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*$")
	ChaincodeVersionRegExp = regexp.MustCompile("^[A-Za-z0-9_.+-]+$")

	collectionNameRegExp = regexp.MustCompile("^[A-Za-z0-9-]+([A-Za-z0-9_-]+)*$")

	
	
	systemChaincodeNames = map[string]struct{}{
		"cscc": {},
		"escc": {},
		"lscc": {},
		"qscc": {},
		"vscc": {},
	}
)

func (i *Invocation) validateInput(name, version string, collections *cb.CollectionConfigPackage) error {
	if !ChaincodeNameRegExp.MatchString(name) {
		return errors.Errorf("invalid chaincode name '%s'. Names can only consist of alphanumerics, '_', and '-' and can only begin with alphanumerics", name)
	}
	if _, ok := systemChaincodeNames[name]; ok {
		return errors.Errorf("chaincode name '%s' is the name of a system chaincode", name)
	}

	if !ChaincodeVersionRegExp.MatchString(version) {
		return errors.Errorf("invalid chaincode version '%s'. Versions can only consist of alphanumerics, '_', '-', '+', and '.'", version)
	}

	collConfigs, err := extractStaticCollectionConfigs(collections)
	if err != nil {
		return err
	}
	
	
	
	
	
	
	channelConfig := i.SCC.ChannelConfigSource.GetStableChannelConfig(i.ChannelID)
	if channelConfig == nil {
		return errors.Errorf("could not get channelconfig for channel '%s'", i.ChannelID)
	}
	mspMgr := channelConfig.MSPManager()
	if mspMgr == nil {
		return errors.Errorf("could not get MSP manager for channel '%s'", i.ChannelID)
	}

	if err := validateCollectionConfigs(collConfigs, mspMgr); err != nil {
		return err
	}

	
	qe := i.SCC.QueryExecutorProvider.TxQueryExecutor(i.Stub.GetChannelID(), i.Stub.GetTxID())
	committedCCDef, err := i.SCC.DeployedCCInfoProvider.ChaincodeInfo(i.ChannelID, name, qe)
	if err != nil {
		return errors.Wrapf(err, "could not retrieve committed definition for chaincode '%s'", name)
	}
	if committedCCDef == nil {
		return nil
	}
	if err := validateCollConfigsAgainstCommittedDef(collConfigs, committedCCDef.ExplicitCollectionConfigPkg); err != nil {
		return err
	}
	return nil
}

func extractStaticCollectionConfigs(collConfigPkg *common.CollectionConfigPackage) ([]*common.StaticCollectionConfig, error) {
	if collConfigPkg == nil || len(collConfigPkg.Config) == 0 {
		return nil, nil
	}
	collConfigs := make([]*common.StaticCollectionConfig, len(collConfigPkg.Config))
	for i, c := range collConfigPkg.Config {
		switch t := c.Payload.(type) {
		case *cb.CollectionConfig_StaticCollectionConfig:
			collConfig := t.StaticCollectionConfig
			if collConfig == nil {
				return nil, errors.Errorf("collection configuration is empty")
			}
			collConfigs[i] = collConfig
		default:
			
			
			return nil, errors.Errorf("collection config contains unexpected payload type: %T", t)
		}
	}
	return collConfigs, nil
}

func validateCollectionConfigs(collConfigs []*common.StaticCollectionConfig, mspMgr msp.MSPManager) error {
	if len(collConfigs) == 0 {
		return nil
	}
	collNamesMap := map[string]struct{}{}
	
	for _, c := range collConfigs {
		if !collectionNameRegExp.MatchString(c.Name) {
			return errors.Errorf("invalid collection name '%s'. Names can only consist of alphanumerics, '_', and '-' and cannot begin with '_'",
				c.Name)
		}
		
		if _, ok := collNamesMap[c.Name]; ok {
			return errors.Errorf("collection-name: %s -- found duplicate in collection configuration",
				c.Name)
		}
		collNamesMap[c.Name] = struct{}{}
		
		if c.MaximumPeerCount < c.RequiredPeerCount {
			return errors.Errorf("collection-name: %s -- maximum peer count (%d) cannot be greater than the required peer count (%d)",
				c.Name, c.MaximumPeerCount, c.RequiredPeerCount)
		}
		if c.RequiredPeerCount < 0 {
			return errors.Errorf("collection-name: %s -- requiredPeerCount (%d) cannot be less than zero",
				c.Name, c.RequiredPeerCount)
		}
		if err := validateCollectionConfigMemberOrgsPolicy(c, mspMgr); err != nil {
			return err
		}
	}
	return nil
}



func validateCollectionConfigMemberOrgsPolicy(coll *common.StaticCollectionConfig, mspMgr msp.MSPManager) error {
	if coll.MemberOrgsPolicy == nil {
		return errors.Errorf("collection member policy is not set for collection '%s'", coll.Name)
	}
	if coll.MemberOrgsPolicy.GetSignaturePolicy() == nil {
		return errors.Errorf("collection member org policy is empty for collection '%s'", coll.Name)
	}

	
	pp := &cauthdsl.EnvelopeBasedPolicyProvider{Deserializer: mspMgr}
	if _, err := pp.NewPolicy(coll.MemberOrgsPolicy.GetSignaturePolicy()); err != nil {
		return errors.WithMessagef(err, "invalid member org policy for collection '%s'", coll.Name)
	}

	
	if err := validateSpOrConcat(coll.MemberOrgsPolicy.GetSignaturePolicy().Rule); err != nil {
		return errors.WithMessagef(err, "collection-name: %s -- error in member org policy", coll.Name)
	}

	msps, err := mspMgr.GetMSPs()
	if err != nil {
		return errors.Wrapf(err, "could not get MSPs")
	}

	
	
	for _, principal := range coll.MemberOrgsPolicy.GetSignaturePolicy().Identities {
		var orgID string
		
		switch principal.PrincipalClassification {

		case mspprotos.MSPPrincipal_ROLE:
			msprole := &mspprotos.MSPRole{}
			err := proto.Unmarshal(principal.Principal, msprole)
			if err != nil {
				return errors.Wrapf(err, "collection-name: %s -- cannot unmarshal identity bytes into MSPRole", coll.GetName())
			}
			orgID = msprole.MspIdentifier
			
			_, ok := msps[orgID]
			if !ok {
				return errors.Errorf("collection-name: %s -- collection member '%s' is not part of the channel", coll.GetName(), orgID)
			}

		case mspprotos.MSPPrincipal_ORGANIZATION_UNIT:
			mspou := &mspprotos.OrganizationUnit{}
			err := proto.Unmarshal(principal.Principal, mspou)
			if err != nil {
				return errors.Wrapf(err, "collection-name: %s -- cannot unmarshal identity bytes into OrganizationUnit", coll.GetName())
			}
			orgID = mspou.MspIdentifier
			
			_, ok := msps[orgID]
			if !ok {
				return errors.Errorf("collection-name: %s -- collection member '%s' is not part of the channel", coll.GetName(), orgID)
			}

		case mspprotos.MSPPrincipal_IDENTITY:
			if _, err := mspMgr.DeserializeIdentity(principal.Principal); err != nil {
				return errors.Errorf("collection-name: %s -- contains an identity that is not part of the channel", coll.GetName())
			}

		default:
			return errors.Errorf("collection-name: %s -- principal type %v is not supported", coll.GetName(), principal.PrincipalClassification)
		}
	}
	return nil
}


func validateSpOrConcat(sp *common.SignaturePolicy) error {
	if sp.GetNOutOf() == nil {
		return nil
	}
	
	if sp.GetNOutOf().N != 1 {
		return errors.Errorf("signature policy is not an OR concatenation, NOutOf %d", sp.GetNOutOf().N)
	}
	
	for _, rule := range sp.GetNOutOf().Rules {
		err := validateSpOrConcat(rule)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateCollConfigsAgainstCommittedDef(
	proposedCollConfs []*common.StaticCollectionConfig,
	committedCollConfPkg *common.CollectionConfigPackage,
) error {
	if committedCollConfPkg == nil || len(committedCollConfPkg.Config) == 0 {
		return nil
	}

	if len(proposedCollConfs) == 0 {
		return errors.Errorf("the proposed collection config does not contain previously defined collections")
	}

	proposedCollsMap := map[string]*common.StaticCollectionConfig{}
	for _, c := range proposedCollConfs {
		proposedCollsMap[c.Name] = c
	}

	
	
	for _, committedCollConfig := range committedCollConfPkg.Config {
		committedColl := committedCollConfig.GetStaticCollectionConfig()
		
		if committedColl == nil {
			return errors.Errorf("unknown collection configuration type")
		}

		newCollection, ok := proposedCollsMap[committedColl.Name]
		if !ok {
			return errors.Errorf("existing collection [%s] missing in the proposed collection configuration", committedColl.Name)
		}

		if newCollection.BlockToLive != committedColl.BlockToLive {
			return errors.Errorf("the BlockToLive in an existing collection [%s] modified. Existing value [%d]", committedColl.Name, committedColl.BlockToLive)
		}
	}
	return nil
}

func (i *Invocation) createOpaqueStates() ([]OpaqueState, error) {
	if i.ApplicationConfig == nil {
		return nil, errors.Errorf("no application config for channel '%s'", i.Stub.GetChannelID())
	}
	orgs := i.ApplicationConfig.Organizations()
	opaqueStates := make([]OpaqueState, 0, len(orgs))
	for _, org := range orgs {
		opaqueStates = append(opaqueStates, &ChaincodePrivateLedgerShim{
			Collection: ImplicitCollectionNameForOrg(org.MSPID()),
			Stub:       i.Stub,
		})
	}
	return opaqueStates, nil
}
