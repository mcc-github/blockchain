/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lscc

import (
	"fmt"
	"regexp"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/ccmetadata"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/core/policyprovider"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protos/common"
	mb "github.com/mcc-github/blockchain/protos/msp"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)








var (
	logger = flogging.MustGetLogger("lscc")
	
	
	chaincodeNameRegExp    = regexp.MustCompile("^[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*$")
	chaincodeVersionRegExp = regexp.MustCompile("^[A-Za-z0-9_.+-]+$")
)

const (
	

	
	INSTALL = "install"

	
	DEPLOY = "deploy"

	
	UPGRADE = "upgrade"

	
	CCEXISTS = "getid"

	
	CHAINCODEEXISTS = "ChaincodeExists"

	
	GETDEPSPEC = "getdepspec"

	
	GETDEPLOYMENTSPEC = "GetDeploymentSpec"

	
	GETCCDATA = "getccdata"

	
	GETCHAINCODEDATA = "GetChaincodeData"

	
	GETCHAINCODES = "getchaincodes"

	
	GETCHAINCODESALIAS = "GetChaincodes"

	
	GETINSTALLEDCHAINCODES = "getinstalledchaincodes"

	
	GETINSTALLEDCHAINCODESALIAS = "GetInstalledChaincodes"

	
	GETCOLLECTIONSCONFIG = "GetCollectionsConfig"

	
	GETCOLLECTIONSCONFIGALIAS = "getcollectionsconfig"
)


type FilesystemSupport interface {
	
	
	PutChaincodeToLocalStorage(ccprovider.CCPackage) error

	
	
	GetChaincodeFromLocalStorage(ccname string, ccversion string) (ccprovider.CCPackage, error)

	
	
	GetChaincodesFromLocalStorage() (*pb.ChaincodeQueryResponse, error)

	
	
	GetInstantiationPolicy(channel string, ccpack ccprovider.CCPackage) ([]byte, error)

	
	
	CheckInstantiationPolicy(signedProposal *pb.SignedProposal, chainName string, instantiationPolicy []byte) error
}




type LifeCycleSysCC struct {
	
	ACLProvider aclmgmt.ACLProvider

	
	
	SCCProvider sysccprovider.SystemChaincodeProvider

	
	
	PolicyChecker policy.PolicyChecker

	
	
	Support FilesystemSupport

	PlatformRegistry *platforms.Registry
}



func New(sccp sysccprovider.SystemChaincodeProvider, ACLProvider aclmgmt.ACLProvider, platformRegistry *platforms.Registry) *LifeCycleSysCC {
	return &LifeCycleSysCC{
		Support:          &supportImpl{},
		PolicyChecker:    policyprovider.GetPolicyChecker(),
		SCCProvider:      sccp,
		ACLProvider:      ACLProvider,
		PlatformRegistry: platformRegistry,
	}
}

func (lscc *LifeCycleSysCC) Name() string              { return "lscc" }
func (lscc *LifeCycleSysCC) Path() string              { return "github.com/mcc-github/blockchain/core/scc/lscc" }
func (lscc *LifeCycleSysCC) InitArgs() [][]byte        { return nil }
func (lscc *LifeCycleSysCC) Chaincode() shim.Chaincode { return lscc }
func (lscc *LifeCycleSysCC) InvokableExternal() bool   { return true }
func (lscc *LifeCycleSysCC) InvokableCC2CC() bool      { return true }
func (lscc *LifeCycleSysCC) Enabled() bool             { return true }

func (lscc *LifeCycleSysCC) ChaincodeContainerInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (*ccprovider.ChaincodeContainerInfo, error) {
	chaincodeDataBytes, err := qe.GetState("lscc", chaincodeName)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve state for chaincode %s", chaincodeName)
	}

	if chaincodeDataBytes == nil {
		return nil, errors.Errorf("chaincode %s not found", chaincodeName)
	}

	
	
	
	cds, _, err := lscc.getCCCode(chaincodeName, chaincodeDataBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get chaincode code")
	}

	return ccprovider.DeploymentSpecToChaincodeContainerInfo(cds), nil
}

func (lscc *LifeCycleSysCC) ChaincodeDefinition(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	chaincodeDataBytes, err := qe.GetState("lscc", chaincodeName)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve state for chaincode %s", chaincodeName)
	}

	if chaincodeDataBytes == nil {
		return nil, errors.Errorf("chaincode %s not found", chaincodeName)
	}

	chaincodeData := &ccprovider.ChaincodeData{}
	err = proto.Unmarshal(chaincodeDataBytes, chaincodeData)
	if err != nil {
		return nil, errors.Wrapf(err, "chaincode %s has bad definition", chaincodeName)
	}

	return chaincodeData, nil
}







func (lscc *LifeCycleSysCC) ValidationInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (plugin string, args []byte, unexpectedErr error, validationErr error) {
	chaincodeDataBytes, err := qe.GetState("lscc", chaincodeName)
	if err != nil {
		
		
		unexpectedErr = errors.Wrapf(err, "could not retrieve state for chaincode %s", chaincodeName)
		return
	}

	if chaincodeDataBytes == nil {
		
		
		validationErr = errors.Errorf("chaincode %s not found", chaincodeName)
		return
	}

	chaincodeData := &ccprovider.ChaincodeData{}
	err = proto.Unmarshal(chaincodeDataBytes, chaincodeData)
	if err != nil {
		
		
		unexpectedErr = errors.Wrapf(err, "chaincode %s has bad definition", chaincodeName)
		return
	}

	plugin = chaincodeData.Vscc
	args = chaincodeData.Policy
	return
}


func (lscc *LifeCycleSysCC) putChaincodeData(stub shim.ChaincodeStubInterface, cd *ccprovider.ChaincodeData) error {
	cdbytes, err := proto.Marshal(cd)
	if err != nil {
		return err
	}

	if cdbytes == nil {
		return MarshallErr(cd.Name)
	}

	err = stub.PutState(cd.Name, cdbytes)

	return err
}



func checkCollectionMemberPolicy(collectionConfig *common.CollectionConfig, mspmgr msp.MSPManager) error {
	if mspmgr == nil {
		return fmt.Errorf("msp manager not set")
	}
	msps, err := mspmgr.GetMSPs()
	if err != nil {
		return errors.Wrapf(err, "error getting channel msp")
	}
	if collectionConfig == nil {
		return fmt.Errorf("collection configuration is not set")
	}
	coll := collectionConfig.GetStaticCollectionConfig()
	if coll == nil {
		return fmt.Errorf("collection configuration is empty")
	}
	if coll.MemberOrgsPolicy == nil {
		return fmt.Errorf("collection member policy is not set")
	}
	if coll.MemberOrgsPolicy.GetSignaturePolicy() == nil {
		return fmt.Errorf("collection member org policy is empty")
	}
	
	
	for _, principal := range coll.MemberOrgsPolicy.GetSignaturePolicy().Identities {
		found := false
		var orgID string
		
		switch principal.PrincipalClassification {

		case mb.MSPPrincipal_ROLE:
			msprole := &mb.MSPRole{}
			err := proto.Unmarshal(principal.Principal, msprole)
			if err != nil {
				return errors.Wrapf(err, "collection-name: %s -- cannot unmarshal identities", coll.GetName())
			}
			orgID = msprole.MspIdentifier
			
			for mspid := range msps {
				if mspid == orgID {
					found = true
					break
				}
			}

		case mb.MSPPrincipal_ORGANIZATION_UNIT:
			mspou := &mb.OrganizationUnit{}
			err := proto.Unmarshal(principal.Principal, mspou)
			if err != nil {
				return errors.Wrapf(err, "collection-name: %s -- cannot unmarshal identities", coll.GetName())
			}
			orgID = mspou.MspIdentifier
			
			for mspid := range msps {
				if mspid == orgID {
					found = true
					break
				}
			}

		case mb.MSPPrincipal_IDENTITY:
			orgID = "identity principal"
			for _, msp := range msps {
				_, err := msp.DeserializeIdentity(principal.Principal)
				if err == nil {
					found = true
					break
				}
			}
		default:
			return fmt.Errorf("collection-name: %s -- principal type %v is not supported", coll.GetName(), principal.PrincipalClassification)
		}
		if !found {
			logger.Warningf("collection-name: %s collection member %s is not part of the channel", coll.GetName(), orgID)
		}
	}

	return nil
}


func (lscc *LifeCycleSysCC) putChaincodeCollectionData(stub shim.ChaincodeStubInterface, cd *ccprovider.ChaincodeData, collectionConfigBytes []byte) error {
	if cd == nil {
		return errors.New("nil ChaincodeData")
	}

	if len(collectionConfigBytes) == 0 {
		logger.Debug("No collection configuration specified")
		return nil
	}

	collections := &common.CollectionConfigPackage{}
	err := proto.Unmarshal(collectionConfigBytes, collections)
	if err != nil {
		return errors.Errorf("invalid collection configuration supplied for chaincode %s:%s", cd.Name, cd.Version)
	}

	mspmgr := mgmt.GetManagerForChain(stub.GetChannelID())
	if mspmgr == nil {
		return fmt.Errorf("could not get MSP manager for channel %s", stub.GetChannelID())
	}
	for _, collectionConfig := range collections.Config {
		err = checkCollectionMemberPolicy(collectionConfig, mspmgr)
		if err != nil {
			return errors.Wrapf(err, "collection member policy check failed")
		}
	}

	key := privdata.BuildCollectionKVSKey(cd.Name)

	err = stub.PutState(key, collectionConfigBytes)
	if err != nil {
		return errors.WithMessagef(err, "error putting collection for chaincode %s:%s", cd.Name, cd.Version)
	}

	return nil
}


func (lscc *LifeCycleSysCC) getChaincodeCollectionData(stub shim.ChaincodeStubInterface, chaincodeName string) pb.Response {
	key := privdata.BuildCollectionKVSKey(chaincodeName)
	collectionsConfigBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(collectionsConfigBytes) == 0 {
		return shim.Error(fmt.Sprintf("collections config not defined for chaincode %s", chaincodeName))
	}
	return shim.Success(collectionsConfigBytes)
}


func (lscc *LifeCycleSysCC) getCCInstance(stub shim.ChaincodeStubInterface, ccname string) ([]byte, error) {
	cdbytes, err := stub.GetState(ccname)
	if err != nil {
		return nil, TXNotFoundErr(err.Error())
	}
	if cdbytes == nil {
		return nil, NotFoundErr(ccname)
	}

	return cdbytes, nil
}


func (lscc *LifeCycleSysCC) getChaincodeData(ccname string, cdbytes []byte) (*ccprovider.ChaincodeData, error) {
	cd := &ccprovider.ChaincodeData{}
	err := proto.Unmarshal(cdbytes, cd)
	if err != nil {
		return nil, MarshallErr(ccname)
	}

	
	if cd.Name != ccname {
		return nil, ChaincodeMismatchErr(fmt.Sprintf("%s!=%s", ccname, cd.Name))
	}

	return cd, nil
}


func (lscc *LifeCycleSysCC) getCCCode(ccname string, cdbytes []byte) (*pb.ChaincodeDeploymentSpec, []byte, error) {
	cd, err := lscc.getChaincodeData(ccname, cdbytes)
	if err != nil {
		return nil, nil, err
	}

	ccpack, err := lscc.Support.GetChaincodeFromLocalStorage(ccname, cd.Version)
	if err != nil {
		return nil, nil, InvalidDeploymentSpecErr(err.Error())
	}

	
	
	
	if err = ccpack.ValidateCC(cd); err != nil {
		return nil, nil, InvalidCCOnFSError(err.Error())
	}

	
	depspec := ccpack.GetDepSpec()
	depspecbytes := ccpack.GetDepSpecBytes()

	return depspec, depspecbytes, nil
}


func (lscc *LifeCycleSysCC) getChaincodes(stub shim.ChaincodeStubInterface) pb.Response {
	
	itr, err := stub.GetStateByRange("", "")

	if err != nil {
		return shim.Error(err.Error())
	}
	defer itr.Close()

	
	var ccInfoArray []*pb.ChaincodeInfo

	for itr.HasNext() {
		response, err := itr.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		
		if privdata.IsCollectionConfigKey(response.Key) {
			continue
		}

		ccdata := &ccprovider.ChaincodeData{}
		if err = proto.Unmarshal(response.Value, ccdata); err != nil {
			return shim.Error(err.Error())
		}

		var path string
		var input string

		
		
		ccpack, err := lscc.Support.GetChaincodeFromLocalStorage(ccdata.Name, ccdata.Version)
		if err == nil {
			path = ccpack.GetDepSpec().GetChaincodeSpec().ChaincodeId.Path
			input = ccpack.GetDepSpec().GetChaincodeSpec().Input.String()
		}

		
		ccInfo := &pb.ChaincodeInfo{Name: ccdata.Name, Version: ccdata.Version, Path: path, Input: input, Escc: ccdata.Escc, Vscc: ccdata.Vscc}
		ccInfoArray = append(ccInfoArray, ccInfo)
	}
	
	
	cqr := &pb.ChaincodeQueryResponse{Chaincodes: ccInfoArray}

	cqrbytes, err := proto.Marshal(cqr)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(cqrbytes)
}


func (lscc *LifeCycleSysCC) getInstalledChaincodes() pb.Response {
	
	
	cqr, err := lscc.Support.GetChaincodesFromLocalStorage()
	if err != nil {
		return shim.Error(err.Error())
	}

	cqrbytes, err := proto.Marshal(cqr)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(cqrbytes)
}


func (lscc *LifeCycleSysCC) isValidChannelName(channel string) bool {
	
	if channel == "" {
		return false
	}
	return true
}



func (lscc *LifeCycleSysCC) isValidChaincodeName(chaincodeName string) error {
	if chaincodeName == "" {
		return EmptyChaincodeNameErr("")
	}

	if !chaincodeNameRegExp.MatchString(chaincodeName) {
		return InvalidChaincodeNameErr(chaincodeName)
	}

	return nil
}




func (lscc *LifeCycleSysCC) isValidChaincodeVersion(chaincodeName string, version string) error {
	if version == "" {
		return EmptyVersionErr(chaincodeName)
	}

	if !chaincodeVersionRegExp.MatchString(version) {
		return InvalidVersionErr(version)
	}

	return nil
}

func isValidStatedbArtifactsTar(statedbArtifactsTar []byte) error {
	
	
	
	archiveFiles, err := ccprovider.ExtractFileEntries(statedbArtifactsTar, "")
	if err != nil {
		return err
	}
	
	for _, archiveDirectoryFiles := range archiveFiles {
		for _, fileEntry := range archiveDirectoryFiles {
			indexData := fileEntry.FileContent
			
			err = ccmetadata.ValidateMetadataFile(fileEntry.FileHeader.Name, indexData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


func (lscc *LifeCycleSysCC) executeInstall(stub shim.ChaincodeStubInterface, ccbytes []byte) error {
	ccpack, err := ccprovider.GetCCPackage(ccbytes)
	if err != nil {
		return err
	}

	cds := ccpack.GetDepSpec()

	if cds == nil {
		return fmt.Errorf("nil deployment spec from from the CC package")
	}

	if err = lscc.isValidChaincodeName(cds.ChaincodeSpec.ChaincodeId.Name); err != nil {
		return err
	}

	if err = lscc.isValidChaincodeVersion(cds.ChaincodeSpec.ChaincodeId.Name, cds.ChaincodeSpec.ChaincodeId.Version); err != nil {
		return err
	}

	if lscc.SCCProvider.IsSysCC(cds.ChaincodeSpec.ChaincodeId.Name) {
		return errors.Errorf("cannot install: %s is the name of a system chaincode", cds.ChaincodeSpec.ChaincodeId.Name)
	}

	
	statedbArtifactsTar, err := ccprovider.ExtractStatedbArtifactsFromCCPackage(ccpack, lscc.PlatformRegistry)
	if err != nil {
		return err
	}

	if err = isValidStatedbArtifactsTar(statedbArtifactsTar); err != nil {
		return InvalidStatedbArtifactsErr(err.Error())
	}

	chaincodeDefinition := &cceventmgmt.ChaincodeDefinition{
		Name:    ccpack.GetChaincodeData().Name,
		Version: ccpack.GetChaincodeData().Version,
		Hash:    ccpack.GetId()} 

	
	
	
	
	err = cceventmgmt.GetMgr().HandleChaincodeInstall(chaincodeDefinition, statedbArtifactsTar)
	defer func() {
		cceventmgmt.GetMgr().ChaincodeInstallDone(err == nil)
	}()
	if err != nil {
		return err
	}

	
	if err = lscc.Support.PutChaincodeToLocalStorage(ccpack); err != nil {
		return err
	}

	logger.Infof("Installed Chaincode [%s] Version [%s] to peer", ccpack.GetChaincodeData().Name, ccpack.GetChaincodeData().Version)

	return nil
}



func (lscc *LifeCycleSysCC) executeDeployOrUpgrade(
	stub shim.ChaincodeStubInterface,
	chainname string,
	cds *pb.ChaincodeDeploymentSpec,
	policy, escc, vscc, collectionConfigBytes []byte,
	function string,
) (*ccprovider.ChaincodeData, error) {

	chaincodeName := cds.ChaincodeSpec.ChaincodeId.Name
	chaincodeVersion := cds.ChaincodeSpec.ChaincodeId.Version

	if err := lscc.isValidChaincodeName(chaincodeName); err != nil {
		return nil, err
	}

	if err := lscc.isValidChaincodeVersion(chaincodeName, chaincodeVersion); err != nil {
		return nil, err
	}

	ccpack, err := lscc.Support.GetChaincodeFromLocalStorage(chaincodeName, chaincodeVersion)
	if err != nil {
		retErrMsg := fmt.Sprintf("cannot get package for chaincode (%s:%s)", chaincodeName, chaincodeVersion)
		logger.Errorf("%s-err:%s", retErrMsg, err)
		return nil, fmt.Errorf("%s", retErrMsg)
	}
	cd := ccpack.GetChaincodeData()

	switch function {
	case DEPLOY:
		return lscc.executeDeploy(stub, chainname, cds, policy, escc, vscc, cd, ccpack, collectionConfigBytes)
	case UPGRADE:
		return lscc.executeUpgrade(stub, chainname, cds, policy, escc, vscc, cd, ccpack, collectionConfigBytes)
	default:
		logger.Panicf("Programming error, unexpected function '%s'", function)
		panic("") 
	}
}


func (lscc *LifeCycleSysCC) executeDeploy(
	stub shim.ChaincodeStubInterface,
	chainname string,
	cds *pb.ChaincodeDeploymentSpec,
	policy []byte,
	escc []byte,
	vscc []byte,
	cdfs *ccprovider.ChaincodeData,
	ccpackfs ccprovider.CCPackage,
	collectionConfigBytes []byte,
) (*ccprovider.ChaincodeData, error) {
	
	chaincodeName := cds.ChaincodeSpec.ChaincodeId.Name
	_, err := lscc.getCCInstance(stub, chaincodeName)
	if err == nil {
		return nil, ExistsErr(chaincodeName)
	}

	
	cdfs.Escc = string(escc)
	cdfs.Vscc = string(vscc)
	cdfs.Policy = policy

	
	cdfs.InstantiationPolicy, err = lscc.Support.GetInstantiationPolicy(chainname, ccpackfs)
	if err != nil {
		return nil, err
	}
	
	signedProp, err := stub.GetSignedProposal()
	if err != nil {
		return nil, err
	}
	err = lscc.Support.CheckInstantiationPolicy(signedProp, chainname, cdfs.InstantiationPolicy)
	if err != nil {
		return nil, err
	}

	err = lscc.putChaincodeData(stub, cdfs)
	if err != nil {
		return nil, err
	}

	err = lscc.putChaincodeCollectionData(stub, cdfs, collectionConfigBytes)
	if err != nil {
		return nil, err
	}

	return cdfs, nil
}


func (lscc *LifeCycleSysCC) executeUpgrade(stub shim.ChaincodeStubInterface, chainName string, cds *pb.ChaincodeDeploymentSpec, policy []byte, escc []byte, vscc []byte, cdfs *ccprovider.ChaincodeData, ccpackfs ccprovider.CCPackage, collectionConfigBytes []byte) (*ccprovider.ChaincodeData, error) {

	chaincodeName := cds.ChaincodeSpec.ChaincodeId.Name

	
	
	
	cdbytes, _ := lscc.getCCInstance(stub, chaincodeName)
	if cdbytes == nil {
		return nil, NotFoundErr(chaincodeName)
	}

	
	cdLedger, err := lscc.getChaincodeData(chaincodeName, cdbytes)
	if err != nil {
		return nil, err
	}

	
	if cdLedger.Version == cds.ChaincodeSpec.ChaincodeId.Version {
		return nil, IdenticalVersionErr(chaincodeName)
	}

	
	if cdLedger.InstantiationPolicy == nil {
		return nil, InstantiationPolicyMissing("")
	}
	
	signedProp, err := stub.GetSignedProposal()
	if err != nil {
		return nil, err
	}
	err = lscc.Support.CheckInstantiationPolicy(signedProp, chainName, cdLedger.InstantiationPolicy)
	if err != nil {
		return nil, err
	}

	
	cdfs.Escc = string(escc)
	cdfs.Vscc = string(vscc)
	cdfs.Policy = policy

	
	cdfs.InstantiationPolicy, err = lscc.Support.GetInstantiationPolicy(chainName, ccpackfs)
	if err != nil {
		return nil, err
	}
	err = lscc.Support.CheckInstantiationPolicy(signedProp, chainName, cdfs.InstantiationPolicy)
	if err != nil {
		return nil, err
	}

	err = lscc.putChaincodeData(stub, cdfs)
	if err != nil {
		return nil, err
	}

	ac, exists := lscc.SCCProvider.GetApplicationConfig(chainName)
	if !exists {
		logger.Panicf("programming error, non-existent appplication config for channel '%s'", chainName)
	}

	if ac.Capabilities().CollectionUpgrade() {
		err = lscc.putChaincodeCollectionData(stub, cdfs, collectionConfigBytes)
		if err != nil {
			return nil, err
		}
	} else {
		if collectionConfigBytes != nil {
			return nil, errors.New(CollectionsConfigUpgradesNotAllowed("").Error())
		}
	}

	lifecycleEvent := &pb.LifecycleEvent{ChaincodeName: chaincodeName}
	lifecycleEventBytes := protoutil.MarshalOrPanic(lifecycleEvent)
	stub.SetEvent(UPGRADE, lifecycleEventBytes)
	return cdfs, nil
}




func (lscc *LifeCycleSysCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}






func (lscc *LifeCycleSysCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) < 1 {
		return shim.Error(InvalidArgsLenErr(len(args)).Error())
	}

	function := string(args[0])

	
	
	sp, err := stub.GetSignedProposal()
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed retrieving signed proposal on executing %s with error %s", function, err))
	}

	switch function {
	case INSTALL:
		if len(args) < 2 {
			return shim.Error(InvalidArgsLenErr(len(args)).Error())
		}

		
		if err = lscc.ACLProvider.CheckACL(resources.Lscc_Install, "", sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s]: %s", function, err))
		}

		depSpec := args[1]

		err := lscc.executeInstall(stub, depSpec)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte("OK"))
	case DEPLOY, UPGRADE:
		
		
		if len(args) < 3 {
			return shim.Error(InvalidArgsLenErr(len(args)).Error())
		}

		
		
		channel := string(args[1])

		if !lscc.isValidChannelName(channel) {
			return shim.Error(InvalidChannelNameErr(channel).Error())
		}

		ac, exists := lscc.SCCProvider.GetApplicationConfig(channel)
		if !exists {
			logger.Panicf("programming error, non-existent appplication config for channel '%s'", channel)
		}

		if ac.Capabilities().LifecycleV20() {
			return shim.Error(fmt.Sprintf("Channel '%s' has been migrated to the new lifecycle, LSCC is now read-only", channel))
		}

		
		if !ac.Capabilities().PrivateChannelData() && len(args) > 6 {
			return shim.Error(PrivateChannelDataNotAvailable("").Error())
		}
		if ac.Capabilities().PrivateChannelData() && len(args) > 7 {
			return shim.Error(InvalidArgsLenErr(len(args)).Error())
		}

		depSpec := args[2]
		cds := &pb.ChaincodeDeploymentSpec{}
		err := proto.Unmarshal(depSpec, cds)
		if err != nil {
			return shim.Error(fmt.Sprintf("error unmarshaling ChaincodeDeploymentSpec: %s", err))
		}

		
		
		
		
		
		var EP []byte
		if len(args) > 3 && len(args[3]) > 0 {
			EP = args[3]
		} else {
			p := cauthdsl.SignedByAnyMember(peer.GetMSPIDs(channel))
			EP, err = protoutil.Marshal(p)
			if err != nil {
				return shim.Error(err.Error())
			}
		}

		var escc []byte
		if len(args) > 4 && len(args[4]) > 0 {
			escc = args[4]
		} else {
			escc = []byte("escc")
		}

		var vscc []byte
		if len(args) > 5 && len(args[5]) > 0 {
			vscc = args[5]
		} else {
			vscc = []byte("vscc")
		}

		var collectionsConfig []byte
		
		
		if ac.Capabilities().PrivateChannelData() && len(args) > 6 {
			collectionsConfig = args[6]
		}

		cd, err := lscc.executeDeployOrUpgrade(stub, channel, cds, EP, escc, vscc, collectionsConfig, function)
		if err != nil {
			return shim.Error(err.Error())
		}
		cdbytes, err := proto.Marshal(cd)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(cdbytes)
	case CCEXISTS, CHAINCODEEXISTS, GETDEPSPEC, GETDEPLOYMENTSPEC, GETCCDATA, GETCHAINCODEDATA:
		if len(args) != 3 {
			return shim.Error(InvalidArgsLenErr(len(args)).Error())
		}

		channel := string(args[1])
		ccname := string(args[2])

		
		var resource string
		switch function {
		case CCEXISTS, CHAINCODEEXISTS:
			resource = resources.Lscc_ChaincodeExists
		case GETDEPSPEC, GETDEPLOYMENTSPEC:
			resource = resources.Lscc_GetDeploymentSpec
		case GETCCDATA, GETCHAINCODEDATA:
			resource = resources.Lscc_GetChaincodeData
		}
		if err = lscc.ACLProvider.CheckACL(resource, channel, sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s][%s]: %s", function, channel, err))
		}

		cdbytes, err := lscc.getCCInstance(stub, ccname)
		if err != nil {
			logger.Errorf("error getting chaincode %s on channel [%s]: %s", ccname, channel, err)
			return shim.Error(err.Error())
		}

		switch function {
		case CCEXISTS, CHAINCODEEXISTS:
			cd, err := lscc.getChaincodeData(ccname, cdbytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			return shim.Success([]byte(cd.Name))
		case GETCCDATA, GETCHAINCODEDATA:
			return shim.Success(cdbytes)
		case GETDEPSPEC, GETDEPLOYMENTSPEC:
			_, depspecbytes, err := lscc.getCCCode(ccname, cdbytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			return shim.Success(depspecbytes)
		default:
			panic("unreachable")
		}
	case GETCHAINCODES, GETCHAINCODESALIAS:
		if len(args) != 1 {
			return shim.Error(InvalidArgsLenErr(len(args)).Error())
		}

		if err = lscc.ACLProvider.CheckACL(resources.Lscc_GetInstantiatedChaincodes, stub.GetChannelID(), sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s][%s]: %s", function, stub.GetChannelID(), err))
		}

		return lscc.getChaincodes(stub)
	case GETINSTALLEDCHAINCODES, GETINSTALLEDCHAINCODESALIAS:
		if len(args) != 1 {
			return shim.Error(InvalidArgsLenErr(len(args)).Error())
		}

		
		if err = lscc.ACLProvider.CheckACL(resources.Lscc_GetInstalledChaincodes, "", sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s]: %s", function, err))
		}

		return lscc.getInstalledChaincodes()
	case GETCOLLECTIONSCONFIG, GETCOLLECTIONSCONFIGALIAS:
		if len(args) != 2 {
			return shim.Error(InvalidArgsLenErr(len(args)).Error())
		}

		chaincodeName := string(args[1])

		logger.Debugf("GetCollectionsConfig, chaincodeName:%s, start to check ACL for current identity policy", chaincodeName)
		if err = lscc.ACLProvider.CheckACL(resources.Lscc_GetCollectionsConfig, stub.GetChannelID(), sp); err != nil {
			logger.Debugf("ACL Check Failed for channel:%s, chaincode:%s", stub.GetChannelID(), chaincodeName)
			return shim.Error(fmt.Sprintf("access denied for [%s]: %s", function, err))
		}

		return lscc.getChaincodeCollectionData(stub, chaincodeName)
	}

	return shim.Error(InvalidFunctionErr(function).Error())
}
