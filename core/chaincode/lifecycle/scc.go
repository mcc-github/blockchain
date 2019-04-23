/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"fmt"
	"regexp"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	persistenceintf "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/dispatcher"
	cb "github.com/mcc-github/blockchain/protos/common"
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

	
	
	QueryApprovalStatusFuncName = "QueryApprovalStatus"

	
	
	CommitChaincodeDefinitionFuncName = "CommitChaincodeDefinition"

	
	
	QueryChaincodeDefinitionFuncName = "QueryChaincodeDefinition"

	
	
	
	QueryNamespaceDefinitionsFuncName = "QueryNamespaceDefinitions"
)



type SCCFunctions interface {
	
	InstallChaincode([]byte) (*chaincode.InstalledChaincode, error)

	
	QueryInstalledChaincode(packageID persistenceintf.PackageID) (*chaincode.InstalledChaincode, error)

	
	QueryInstalledChaincodes() (chaincodes []chaincode.InstalledChaincode, err error)

	
	ApproveChaincodeDefinitionForOrg(chname, ccname string, cd *ChaincodeDefinition, packageID persistenceintf.PackageID, publicState ReadableState, orgState ReadWritableState) error

	
	
	QueryApprovalStatus(chname, ccname string, cd *ChaincodeDefinition, publicState ReadWritableState, orgStates []OpaqueState) ([]bool, error)

	
	CommitChaincodeDefinition(chname, ccname string, cd *ChaincodeDefinition, publicState ReadWritableState, orgStates []OpaqueState) ([]bool, error)

	
	QueryChaincodeDefinition(name string, publicState ReadableState) (*ChaincodeDefinition, error)

	
	QueryNamespaceDefinitions(publicState RangeableState) (map[string]string, error)
}





type ChannelConfigSource interface {
	
	
	
	GetStableChannelConfig(channelID string) channelconfig.Resources
}



type SCC struct {
	OrgMSPID string

	ACLProvider aclmgmt.ACLProvider

	ChannelConfigSource ChannelConfigSource

	
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
	if channelID := stub.GetChannelID(); channelID != "" {
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



func (i *Invocation) QueryInstalledChaincodes(input *lb.QueryInstalledChaincodesArgs) (proto.Message, error) {

	logger.Debugf("received invocation of QueryInstalledChaincodes")

	chaincodes, err := i.SCC.Functions.QueryInstalledChaincodes()
	if err != nil {
		return nil, err
	}

	result := &lb.QueryInstalledChaincodesResult{}
	for _, chaincode := range chaincodes {
		result.InstalledChaincodes = append(
			result.InstalledChaincodes,
			&lb.QueryInstalledChaincodesResult_InstalledChaincode{
				Label:     chaincode.Label,
				PackageId: chaincode.PackageID.String(),
			})
	}
	return result, nil
}



func (i *Invocation) ApproveChaincodeDefinitionForMyOrg(input *lb.ApproveChaincodeDefinitionForMyOrgArgs) (proto.Message, error) {
	if err := validateNameVersion(input.Name, input.Version); err != nil {
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



func (i *Invocation) QueryApprovalStatus(input *lb.QueryApprovalStatusArgs) (proto.Message, error) {
	if i.ApplicationConfig == nil {
		return nil, errors.Errorf("no application config for channel '%s'", i.Stub.GetChannelID())
	}

	orgs := i.ApplicationConfig.Organizations()
	opaqueStates := make([]OpaqueState, 0, len(orgs))
	orgNames := make([]string, 0, len(orgs))
	for _, org := range orgs {
		orgNames = append(orgNames, org.MSPID())
		opaqueStates = append(opaqueStates, &ChaincodePrivateLedgerShim{
			Collection: ImplicitCollectionNameForOrg(org.MSPID()),
			Stub:       i.Stub,
		})
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

	logger.Debugf("received invocation of QueryApprovalStatus on channel '%s' for definition '%s'",
		i.Stub.GetChannelID(),
		cd,
	)

	approved, err := i.SCC.Functions.QueryApprovalStatus(
		i.Stub.GetChannelID(),
		input.Name,
		cd,
		i.Stub,
		opaqueStates,
	)
	if err != nil {
		return nil, err
	}

	orgApproval := make(map[string]bool)
	for i, org := range orgNames {
		orgApproval[org] = approved[i]
	}

	return &lb.QueryApprovalStatusResults{
		Approved: orgApproval,
	}, nil
}



func (i *Invocation) CommitChaincodeDefinition(input *lb.CommitChaincodeDefinitionArgs) (proto.Message, error) {
	if err := validateNameVersion(input.Name, input.Version); err != nil {
		return nil, err
	}

	if i.ApplicationConfig == nil {
		return nil, errors.Errorf("no application config for channel '%s'", i.Stub.GetChannelID())
	}

	orgs := i.ApplicationConfig.Organizations()
	opaqueStates := make([]OpaqueState, 0, len(orgs))
	myOrgIndex := -1
	for _, org := range orgs {
		opaqueStates = append(opaqueStates, &ChaincodePrivateLedgerShim{
			Collection: ImplicitCollectionNameForOrg(org.MSPID()),
			Stub:       i.Stub,
		})
		if org.MSPID() == i.SCC.OrgMSPID {
			myOrgIndex = len(opaqueStates) - 1
		}
	}

	if myOrgIndex == -1 {
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

	agreement, err := i.SCC.Functions.CommitChaincodeDefinition(
		i.Stub.GetChannelID(),
		input.Name,
		cd,
		i.Stub,
		opaqueStates,
	)

	if err != nil {
		return nil, err
	}

	if !agreement[myOrgIndex] {
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

	return &lb.QueryChaincodeDefinitionResult{
		Sequence:            definedChaincode.Sequence,
		Version:             definedChaincode.EndorsementInfo.Version,
		EndorsementPlugin:   definedChaincode.EndorsementInfo.EndorsementPlugin,
		ValidationPlugin:    definedChaincode.ValidationInfo.ValidationPlugin,
		ValidationParameter: definedChaincode.ValidationInfo.ValidationParameter,
		InitRequired:        definedChaincode.EndorsementInfo.InitRequired,
		Collections:         definedChaincode.Collections,
	}, nil
}



func (i *Invocation) QueryNamespaceDefinitions(input *lb.QueryNamespaceDefinitionsArgs) (proto.Message, error) {

	logger.Debugf("received invocation of QueryNamespaceDefinitions on channel '%s'",
		i.Stub.GetChannelID(),
	)

	namespaces, err := i.SCC.Functions.QueryNamespaceDefinitions(&ChaincodePublicLedgerShim{ChaincodeStubInterface: i.Stub})
	if err != nil {
		return nil, err
	}
	result := map[string]*lb.QueryNamespaceDefinitionsResult_Namespace{}
	for namespace, nType := range namespaces {
		result[namespace] = &lb.QueryNamespaceDefinitionsResult_Namespace{
			Type: nType,
		}
	}
	return &lb.QueryNamespaceDefinitionsResult{
		Namespaces: result,
	}, nil
}

var (
	
	
	chaincodeNameRegExp    = regexp.MustCompile("^[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*$")
	chaincodeVersionRegExp = regexp.MustCompile("^[A-Za-z0-9_.+-]+$")
)

func validateNameVersion(name, version string) error {
	if !chaincodeNameRegExp.MatchString(name) {
		return errors.Errorf("invalid chaincode name '%s'. Names can only consist of alphanumerics, '_', and '-'", name)
	}
	if !chaincodeVersionRegExp.MatchString(version) {
		return errors.Errorf("invalid chaincode version '%s'. Versions can only consist of alphanumerics, '_', '-', '+', and '.'", version)
	}

	systemChaincodeNames := map[string]struct{}{
		"cscc": {},
		"escc": {},
		"lscc": {},
		"qscc": {},
		"vscc": {},
	}
	if _, ok := systemChaincodeNames[name]; ok {
		return errors.Errorf("chaincode name '%s' is the name of a system chaincode", name)
	}

	return nil
}
