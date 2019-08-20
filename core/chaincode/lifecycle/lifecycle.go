/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"bytes"
	"fmt"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"
	"github.com/mcc-github/blockchain/protoutil"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("lifecycle")

const (
	
	
	
	
	NamespacesName = "namespaces"

	
	
	
	ChaincodeSourcesName = "chaincode-sources"

	
	ChaincodeDefinitionType = "ChaincodeDefinition"

	
	FriendlyChaincodeDefinitionType = "Chaincode"

	
	DefaultEndorsementPolicyRef = "/Channel/Application/Endorsement"
)

var (
	DefaultEndorsementPolicyBytes = protoutil.MarshalOrPanic(&pb.ApplicationPolicy{
		Type: &pb.ApplicationPolicy_ChannelConfigPolicyReference{
			ChannelConfigPolicyReference: DefaultEndorsementPolicyRef,
		},
	})
)









































type ChaincodeLocalPackage struct {
	PackageID string
}





type ChaincodeParameters struct {
	EndorsementInfo *lb.ChaincodeEndorsementInfo
	ValidationInfo  *lb.ChaincodeValidationInfo
	Collections     *cb.CollectionConfigPackage
}

func (cp *ChaincodeParameters) Equal(ocp *ChaincodeParameters) error {
	switch {
	case cp.EndorsementInfo.Version != ocp.EndorsementInfo.Version:
		return errors.Errorf("Version '%s' != '%s'", cp.EndorsementInfo.Version, ocp.EndorsementInfo.Version)
	case cp.EndorsementInfo.EndorsementPlugin != ocp.EndorsementInfo.EndorsementPlugin:
		return errors.Errorf("EndorsementPlugin '%s' != '%s'", cp.EndorsementInfo.EndorsementPlugin, ocp.EndorsementInfo.EndorsementPlugin)
	case cp.ValidationInfo.ValidationPlugin != ocp.ValidationInfo.ValidationPlugin:
		return errors.Errorf("ValidationPlugin '%s' != '%s'", cp.ValidationInfo.ValidationPlugin, ocp.ValidationInfo.ValidationPlugin)
	case !bytes.Equal(cp.ValidationInfo.ValidationParameter, ocp.ValidationInfo.ValidationParameter):
		return errors.Errorf("ValidationParameter '%x' != '%x'", cp.ValidationInfo.ValidationParameter, ocp.ValidationInfo.ValidationParameter)
	case !proto.Equal(cp.Collections, ocp.Collections):
		return errors.Errorf("Collections do not match")
	default:
	}
	return nil
}






type ChaincodeDefinition struct {
	Sequence        int64
	EndorsementInfo *lb.ChaincodeEndorsementInfo
	ValidationInfo  *lb.ChaincodeValidationInfo
	Collections     *cb.CollectionConfigPackage
}


func (cd *ChaincodeDefinition) Parameters() *ChaincodeParameters {
	return &ChaincodeParameters{
		EndorsementInfo: cd.EndorsementInfo,
		ValidationInfo:  cd.ValidationInfo,
		Collections:     cd.Collections,
	}
}

func (cd *ChaincodeDefinition) String() string {
	endorsementInfo := "endorsement info: <EMPTY>"
	if cd.EndorsementInfo != nil {
		endorsementInfo = fmt.Sprintf("endorsement info: (version: '%s', plugin: '%s', init required: %t)",
			cd.EndorsementInfo.Version,
			cd.EndorsementInfo.EndorsementPlugin,
			cd.EndorsementInfo.InitRequired,
		)
	}

	validationInfo := "validation info: <EMPTY>"
	if cd.ValidationInfo != nil {
		validationInfo = fmt.Sprintf("validation info: (plugin: '%s', policy: '%x')",
			cd.ValidationInfo.ValidationPlugin,
			cd.ValidationInfo.ValidationParameter,
		)
	}

	return fmt.Sprintf("sequence: %d, %s, %s, collections: (%+v)",
		cd.Sequence,
		endorsementInfo,
		validationInfo,
		cd.Collections,
	)
}


type ChaincodeStore interface {
	Save(label string, ccInstallPkg []byte) (string, error)
	ListInstalledChaincodes() ([]chaincode.InstalledChaincode, error)
	Load(packageID string) (ccInstallPkg []byte, err error)
}

type PackageParser interface {
	Parse(data []byte) (*persistence.ChaincodePackage, error)
}


type InstallListener interface {
	HandleChaincodeInstalled(md *persistence.ChaincodePackageMetadata, packageID string)
}


type InstalledChaincodesLister interface {
	ListInstalledChaincodes() []*chaincode.InstalledChaincode
	GetInstalledChaincode(packageID string) (*chaincode.InstalledChaincode, error)
}




type Resources struct {
	ChannelConfigSource ChannelConfigSource
	ChaincodeStore      ChaincodeStore
	PackageParser       PackageParser
	Serializer          *Serializer
}




func (r *Resources) ChaincodeDefinitionIfDefined(chaincodeName string, state ReadableState) (bool, *ChaincodeDefinition, error) {
	if chaincodeName == LifecycleNamespace {
		return true, &ChaincodeDefinition{
			EndorsementInfo: &lb.ChaincodeEndorsementInfo{
				InitRequired: false,
			},
			ValidationInfo: &lb.ChaincodeValidationInfo{},
		}, nil
	}

	metadata, ok, err := r.Serializer.DeserializeMetadata(NamespacesName, chaincodeName, state)
	if err != nil {
		return false, nil, errors.WithMessagef(err, "could not deserialize metadata for chaincode %s", chaincodeName)
	}

	if !ok {
		return false, nil, nil
	}

	if metadata.Datatype != ChaincodeDefinitionType {
		return false, nil, errors.Errorf("not a chaincode type: %s", metadata.Datatype)
	}

	definedChaincode := &ChaincodeDefinition{}
	err = r.Serializer.Deserialize(NamespacesName, chaincodeName, metadata, definedChaincode, state)
	if err != nil {
		return false, nil, errors.WithMessagef(err, "could not deserialize chaincode definition for chaincode %s", chaincodeName)
	}

	return true, definedChaincode, nil
}







type ExternalFunctions struct {
	Resources                 *Resources
	InstallListener           InstallListener
	InstalledChaincodesLister InstalledChaincodesLister
}




func (ef *ExternalFunctions) CheckCommitReadiness(chname, ccname string, cd *ChaincodeDefinition, publicState ReadWritableState, orgStates []OpaqueState) (map[string]bool, error) {
	currentSequence, err := ef.Resources.Serializer.DeserializeFieldAsInt64(NamespacesName, ccname, "Sequence", publicState)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get current sequence")
	}

	if cd.Sequence != currentSequence+1 {
		return nil, errors.Errorf("requested sequence is %d, but new definition must be sequence %d", cd.Sequence, currentSequence+1)
	}

	if err := ef.SetChaincodeDefinitionDefaults(chname, cd); err != nil {
		return nil, errors.WithMessagef(err, "could not set defaults for chaincode definition in channel %s", chname)
	}

	var approvals map[string]bool
	if approvals, err = ef.QueryOrgApprovals(ccname, cd, orgStates); err != nil {
		return nil, err
	}

	logger.Infof("successfully simulated committing chaincode definition %s, name '%s' on channel '%s'", cd, ccname, chname)

	return approvals, nil
}







func (ef *ExternalFunctions) CommitChaincodeDefinition(chname, ccname string, cd *ChaincodeDefinition, publicState ReadWritableState, orgStates []OpaqueState) (map[string]bool, error) {
	approvals, err := ef.CheckCommitReadiness(chname, ccname, cd, publicState, orgStates)
	if err != nil {
		return nil, err
	}

	if err = ef.Resources.Serializer.Serialize(NamespacesName, ccname, cd, publicState); err != nil {
		return nil, errors.WithMessage(err, "could not serialize chaincode definition")
	}

	logger.Infof("successfully committed definition %s, name '%s' on channel '%s'", cd, ccname, chname)

	return approvals, nil
}



func (ef *ExternalFunctions) DefaultEndorsementPolicyAsBytes(channelID string) ([]byte, error) {
	channelConfig := ef.Resources.ChannelConfigSource.GetStableChannelConfig(channelID)
	if channelConfig == nil {
		return nil, errors.Errorf("could not get channel config for channel '%s'", channelID)
	}

	
	if _, ok := channelConfig.PolicyManager().GetPolicy(DefaultEndorsementPolicyRef); ok {
		return DefaultEndorsementPolicyBytes, nil
	}

	return nil, errors.Errorf(
		"Policy '%s' must be defined for channel '%s' before chaincode operations can be attempted",
		DefaultEndorsementPolicyRef,
		channelID,
	)
}



func (ef *ExternalFunctions) SetChaincodeDefinitionDefaults(chname string, cd *ChaincodeDefinition) error {
	if cd.EndorsementInfo.EndorsementPlugin == "" {
		
		
		
		cd.EndorsementInfo.EndorsementPlugin = "escc"
	}

	if cd.ValidationInfo.ValidationPlugin == "" {
		
		
		
		cd.ValidationInfo.ValidationPlugin = "vscc"
	}

	if len(cd.ValidationInfo.ValidationParameter) == 0 {
		policyBytes, err := ef.DefaultEndorsementPolicyAsBytes(chname)
		if err != nil {
			return err
		}

		cd.ValidationInfo.ValidationParameter = policyBytes
	}

	return nil
}




func (ef *ExternalFunctions) ApproveChaincodeDefinitionForOrg(chname, ccname string, cd *ChaincodeDefinition, packageID string, publicState ReadableState, orgState ReadWritableState) error {
	
	currentSequence, err := ef.Resources.Serializer.DeserializeFieldAsInt64(NamespacesName, ccname, "Sequence", publicState)
	if err != nil {
		return errors.WithMessage(err, "could not get current sequence")
	}

	requestedSequence := cd.Sequence

	if currentSequence == requestedSequence && requestedSequence == 0 {
		return errors.Errorf("requested sequence is 0, but first definable sequence number is 1")
	}

	if requestedSequence < currentSequence {
		return errors.Errorf("currently defined sequence %d is larger than requested sequence %d", currentSequence, requestedSequence)
	}

	if requestedSequence > currentSequence+1 {
		return errors.Errorf("requested sequence %d is larger than the next available sequence number %d", requestedSequence, currentSequence+1)
	}

	if err := ef.SetChaincodeDefinitionDefaults(chname, cd); err != nil {
		return errors.WithMessagef(err, "could not set defaults for chaincode definition in channel %s", chname)
	}

	if requestedSequence == currentSequence {
		metadata, ok, err := ef.Resources.Serializer.DeserializeMetadata(NamespacesName, ccname, publicState)
		if err != nil {
			return errors.WithMessage(err, "could not fetch metadata for current definition")
		}
		if !ok {
			return errors.Errorf("missing metadata for currently committed sequence number (%d)", currentSequence)
		}

		definedChaincode := &ChaincodeDefinition{}
		if err := ef.Resources.Serializer.Deserialize(NamespacesName, ccname, metadata, definedChaincode, publicState); err != nil {
			return errors.WithMessagef(err, "could not deserialize namespace %s as chaincode", ccname)
		}

		if err := definedChaincode.Parameters().Equal(cd.Parameters()); err != nil {
			return errors.WithMessagef(err, "attempted to define the current sequence (%d) for namespace %s, but", currentSequence, ccname)
		}
	}

	privateName := fmt.Sprintf("%s#%d", ccname, requestedSequence)
	if err := ef.Resources.Serializer.Serialize(NamespacesName, privateName, cd.Parameters(), orgState); err != nil {
		return errors.WithMessage(err, "could not serialize chaincode parameters to state")
	}

	
	
	
	
	
	
	if err := ef.Resources.Serializer.Serialize(ChaincodeSourcesName, privateName, &ChaincodeLocalPackage{
		PackageID: packageID,
	}, orgState); err != nil {
		return errors.WithMessage(err, "could not serialize chaincode package info to state")
	}

	logger.Infof("successfully approved definition %s, name '%s' on channel '%s'", cd, ccname, chname)

	return nil
}




type ErrNamespaceNotDefined struct {
	Namespace string
}

func (e ErrNamespaceNotDefined) Error() string {
	return fmt.Sprintf("namespace %s is not defined", e.Namespace)
}



func (ef *ExternalFunctions) QueryChaincodeDefinition(name string, publicState ReadableState) (*ChaincodeDefinition, error) {
	metadata, ok, err := ef.Resources.Serializer.DeserializeMetadata(NamespacesName, name, publicState)
	if err != nil {
		return nil, errors.WithMessagef(err, "could not fetch metadata for namespace %s", name)
	}
	if !ok {
		return nil, ErrNamespaceNotDefined{Namespace: name}
	}

	definedChaincode := &ChaincodeDefinition{}
	if err := ef.Resources.Serializer.Deserialize(NamespacesName, name, metadata, definedChaincode, publicState); err != nil {
		return nil, errors.WithMessagef(err, "could not deserialize namespace %s as chaincode", name)
	}

	logger.Infof("successfully queried definition %s, name '%s'", definedChaincode, name)

	return definedChaincode, nil
}




func (ef *ExternalFunctions) QueryOrgApprovals(name string, cd *ChaincodeDefinition, orgStates []OpaqueState) (map[string]bool, error) {
	approvals := map[string]bool{}
	privateName := fmt.Sprintf("%s#%d", name, cd.Sequence)
	for _, orgState := range orgStates {
		match, err := ef.Resources.Serializer.IsSerialized(NamespacesName, privateName, cd.Parameters(), orgState)
		if err != nil {
			return nil, errors.WithMessagef(err, "serialization check failed for key %s", privateName)
		}

		org := OrgFromImplicitCollectionName(orgState.CollectionName())
		approvals[org] = match
	}

	return approvals, nil
}



func (ef *ExternalFunctions) InstallChaincode(chaincodeInstallPackage []byte) (*chaincode.InstalledChaincode, error) {
	
	pkg, err := ef.Resources.PackageParser.Parse(chaincodeInstallPackage)
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse as a chaincode install package")
	}

	if pkg.Metadata == nil {
		return nil, errors.New("empty metadata for supplied chaincode")
	}

	packageID, err := ef.Resources.ChaincodeStore.Save(pkg.Metadata.Label, chaincodeInstallPackage)
	if err != nil {
		return nil, errors.WithMessage(err, "could not save cc install package")
	}

	if ef.InstallListener != nil {
		ef.InstallListener.HandleChaincodeInstalled(pkg.Metadata, packageID)
	}

	return &chaincode.InstalledChaincode{
		PackageID: packageID,
		Label:     pkg.Metadata.Label,
	}, nil
}



func (ef *ExternalFunctions) GetInstalledChaincodePackage(packageID string) ([]byte, error) {
	pkgBytes, err := ef.Resources.ChaincodeStore.Load(packageID)
	if err != nil {
		return nil, errors.WithMessage(err, "could not load cc install package")
	}

	return pkgBytes, nil
}



func (ef *ExternalFunctions) QueryNamespaceDefinitions(publicState RangeableState) (map[string]string, error) {
	metadatas, err := ef.Resources.Serializer.DeserializeAllMetadata(NamespacesName, publicState)
	if err != nil {
		return nil, errors.WithMessage(err, "could not query namespace metadata")
	}

	result := map[string]string{}
	for key, value := range metadatas {
		switch value.Datatype {
		case ChaincodeDefinitionType:
			result[key] = FriendlyChaincodeDefinitionType
		default:
			
			result[key] = value.Datatype
		}
	}
	return result, nil
}


func (ef *ExternalFunctions) QueryInstalledChaincode(packageID string) (*chaincode.InstalledChaincode, error) {
	return ef.InstalledChaincodesLister.GetInstalledChaincode(packageID)
}


func (ef *ExternalFunctions) QueryInstalledChaincodes() []*chaincode.InstalledChaincode {
	return ef.InstalledChaincodesLister.ListInstalledChaincodes()
}
