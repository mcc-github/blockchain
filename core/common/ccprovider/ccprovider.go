/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccprovider

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

var ccproviderLogger = flogging.MustGetLogger("ccprovider")

var chaincodeInstallPath string






type CCPackage interface {
	
	InitFromBuffer(buf []byte) (*ChaincodeData, error)

	
	InitFromFS(ccname string, ccversion string) ([]byte, *pb.ChaincodeDeploymentSpec, error)

	
	PutChaincodeToFS() error

	
	GetDepSpec() *pb.ChaincodeDeploymentSpec

	
	GetDepSpecBytes() []byte

	
	
	
	ValidateCC(ccdata *ChaincodeData) error

	
	GetPackageObject() proto.Message

	
	GetChaincodeData() *ChaincodeData

	
	GetId() []byte
}


func SetChaincodesPath(path string) {
	if s, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path, 0755); err != nil {
				panic(fmt.Sprintf("Could not create chaincodes install path: %s", err))
			}
		} else {
			panic(fmt.Sprintf("Could not stat chaincodes install path: %s", err))
		}
	} else if !s.IsDir() {
		panic(fmt.Errorf("chaincode path exists but not a dir: %s", path))
	}

	chaincodeInstallPath = path
}

func GetChaincodePackage(ccname string, ccversion string) ([]byte, error) {
	return GetChaincodePackageFromPath(ccname, ccversion, chaincodeInstallPath)
}


func GetChaincodePackageFromPath(ccname string, ccversion string, ccInstallPath string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s.%s", ccInstallPath, ccname, ccversion)
	var ccbytes []byte
	var err error
	if ccbytes, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}
	return ccbytes, nil
}


func ChaincodePackageExists(ccname string, ccversion string) (bool, error) {
	path := filepath.Join(chaincodeInstallPath, ccname+"."+ccversion)
	_, err := os.Stat(path)
	if err == nil {
		
		return true, nil
	}
	return false, err
}

type CCCacheSupport interface {
	
	GetChaincode(ccname string, ccversion string) (CCPackage, error)
}



type CCInfoFSImpl struct{}



func (cifs *CCInfoFSImpl) GetChaincode(ccname string, ccversion string) (CCPackage, error) {
	return cifs.GetChaincodeFromPath(ccname, ccversion, chaincodeInstallPath)
}


func (*CCInfoFSImpl) GetChaincodeFromPath(ccname string, ccversion string, path string) (CCPackage, error) {
	
	cccdspack := &CDSPackage{}
	_, _, err := cccdspack.InitFromPath(ccname, ccversion, path)
	if err != nil {
		
		ccscdspack := &SignedCDSPackage{}
		_, _, err = ccscdspack.InitFromPath(ccname, ccversion, path)
		if err != nil {
			return nil, err
		}
		return ccscdspack, nil
	}
	return cccdspack, nil
}



func (*CCInfoFSImpl) PutChaincode(depSpec *pb.ChaincodeDeploymentSpec) (CCPackage, error) {
	buf, err := proto.Marshal(depSpec)
	if err != nil {
		return nil, err
	}
	cccdspack := &CDSPackage{}
	if _, err := cccdspack.InitFromBuffer(buf); err != nil {
		return nil, err
	}
	err = cccdspack.PutChaincodeToFS()
	if err != nil {
		return nil, err
	}

	return cccdspack, nil
}


















var ccInfoFSProvider = &CCInfoFSImpl{}


var ccInfoCache = NewCCInfoCache(ccInfoFSProvider)



var ccInfoCacheEnabled bool


func EnableCCInfoCache() {
	ccInfoCacheEnabled = true
}


func GetChaincodeFromFS(ccname string, ccversion string) (CCPackage, error) {
	return ccInfoFSProvider.GetChaincode(ccname, ccversion)
}




func PutChaincodeIntoFS(depSpec *pb.ChaincodeDeploymentSpec) error {
	_, err := ccInfoFSProvider.PutChaincode(depSpec)
	return err
}


func GetChaincodeData(ccname string, ccversion string) (*ChaincodeData, error) {
	if ccInfoCacheEnabled {
		ccproviderLogger.Debugf("Getting chaincode data for <%s, %s> from cache", ccname, ccversion)
		return ccInfoCache.GetChaincodeData(ccname, ccversion)
	}
	if ccpack, err := ccInfoFSProvider.GetChaincode(ccname, ccversion); err != nil {
		return nil, err
	} else {
		ccproviderLogger.Infof("Putting chaincode data for <%s, %s> into cache", ccname, ccversion)
		return ccpack.GetChaincodeData(), nil
	}
}

func CheckInstantiationPolicy(name, version string, cdLedger *ChaincodeData) error {
	ccdata, err := GetChaincodeData(name, version)
	if err != nil {
		return err
	}

	
	
	
	
	
	
	
	
	
	
	
	
	if ccdata.InstantiationPolicy != nil {
		if !bytes.Equal(ccdata.InstantiationPolicy, cdLedger.InstantiationPolicy) {
			return fmt.Errorf("Instantiation policy mismatch for cc %s/%s", name, version)
		}
	}

	return nil
}



func GetCCPackage(buf []byte) (CCPackage, error) {
	
	cccdspack := &CDSPackage{}
	if _, err := cccdspack.InitFromBuffer(buf); err != nil {
		
		ccscdspack := &SignedCDSPackage{}
		if _, err := ccscdspack.InitFromBuffer(buf); err != nil {
			return nil, err
		}
		return ccscdspack, nil
	}
	return cccdspack, nil
}





func GetInstalledChaincodes() (*pb.ChaincodeQueryResponse, error) {
	files, err := ioutil.ReadDir(chaincodeInstallPath)
	if err != nil {
		return nil, err
	}

	
	var ccInfoArray []*pb.ChaincodeInfo

	for _, file := range files {
		
		
		fileNameArray := strings.SplitN(file.Name(), ".", 2)

		
		if len(fileNameArray) == 2 {
			ccname := fileNameArray[0]
			ccversion := fileNameArray[1]
			ccpack, err := GetChaincodeFromFS(ccname, ccversion)
			if err != nil {
				
				
				ccproviderLogger.Errorf("Unreadable chaincode file found on filesystem: %s", file.Name())
				continue
			}

			cdsfs := ccpack.GetDepSpec()

			name := cdsfs.GetChaincodeSpec().GetChaincodeId().Name
			version := cdsfs.GetChaincodeSpec().GetChaincodeId().Version
			if name != ccname || version != ccversion {
				
				
				ccproviderLogger.Errorf("Chaincode file's name/version has been modified on the filesystem: %s", file.Name())
				continue
			}

			path := cdsfs.GetChaincodeSpec().ChaincodeId.Path
			
			input, escc, vscc := "", "", ""

			ccInfo := &pb.ChaincodeInfo{Name: name, Version: version, Path: path, Input: input, Escc: escc, Vscc: vscc, Id: ccpack.GetId()}

			
			ccInfoArray = append(ccInfoArray, ccInfo)
		}
	}
	
	
	cqr := &pb.ChaincodeQueryResponse{Chaincodes: ccInfoArray}

	return cqr, nil
}


type CCContext struct {
	
	ChainID string

	
	Name string

	
	Version string

	
	TxID string

	
	Syscc bool

	
	
	SignedProposal *pb.SignedProposal

	
	
	Proposal *pb.Proposal

	
	canonicalName string

	
	ProposalDecorations map[string][]byte
}


func NewCCContext(cname, name, version, txid string, syscc bool, signedProp *pb.SignedProposal, prop *pb.Proposal) *CCContext {
	cccid := &CCContext{
		ChainID:             cname,
		Name:                name,
		Version:             version,
		TxID:                txid,
		Syscc:               syscc,
		SignedProposal:      signedProp,
		Proposal:            prop,
		canonicalName:       name + ":" + version,
		ProposalDecorations: nil,
	}

	
	
	if version == "" {
		panic(fmt.Sprintf("---empty version---(%s)", cccid))
	}

	ccproviderLogger.Debugf("NewCCCC(%s)", cccid)
	return cccid
}

func (cccid *CCContext) String() string {
	return fmt.Sprintf("chain=%s,chaincode=%s,version=%s,txid=%s,syscc=%t,proposal=%p,canname=%s",
		cccid.ChainID, cccid.Name, cccid.Version, cccid.TxID, cccid.Syscc, cccid.Proposal, cccid.canonicalName)
}


func (cccid *CCContext) GetCanonicalName() string {
	if cccid.canonicalName == "" {
		panic(fmt.Sprintf("missing canonical name: %s", cccid))
	}

	return cccid.canonicalName
}




type ChaincodeDefinition interface {
	
	CCName() string

	
	Hash() []byte

	
	CCVersion() string

	
	
	
	
	Validation() (string, []byte)

	
	
	Endorsement() string
}






type ChaincodeData struct {
	
	Name string `protobuf:"bytes,1,opt,name=name"`

	
	Version string `protobuf:"bytes,2,opt,name=version"`

	
	Escc string `protobuf:"bytes,3,opt,name=escc"`

	
	Vscc string `protobuf:"bytes,4,opt,name=vscc"`

	
	Policy []byte `protobuf:"bytes,5,opt,name=policy,proto3"`

	
	Data []byte `protobuf:"bytes,6,opt,name=data,proto3"`

	
	
	Id []byte `protobuf:"bytes,7,opt,name=id,proto3"`

	
	InstantiationPolicy []byte `protobuf:"bytes,8,opt,name=instantiation_policy,proto3"`
}


func (cd *ChaincodeData) CCName() string {
	return cd.Name
}


func (cd *ChaincodeData) Hash() []byte {
	return cd.Id
}


func (cd *ChaincodeData) CCVersion() string {
	return cd.Version
}





func (cd *ChaincodeData) Validation() (string, []byte) {
	return cd.Vscc, cd.Policy
}



func (cd *ChaincodeData) Endorsement() string {
	return cd.Escc
}




func (cd *ChaincodeData) Reset() { *cd = ChaincodeData{} }


func (cd *ChaincodeData) String() string { return proto.CompactTextString(cd) }


func (*ChaincodeData) ProtoMessage() {}



type ChaincodeSpecGetter interface {
	GetChaincodeSpec() *pb.ChaincodeSpec
}





type ChaincodeProvider interface {
	
	
	
	GetContext(ledger ledger.PeerLedger, txid string) (context.Context, ledger.TxSimulator, error)
	
	ExecuteChaincode(ctxt context.Context, cccid *CCContext, args [][]byte) (*pb.Response, *pb.ChaincodeEvent, error)
	
	Execute(ctxt context.Context, cccid *CCContext, spec ChaincodeSpecGetter) (*pb.Response, *pb.ChaincodeEvent, error)
	
	Stop(ctxt context.Context, cccid *CCContext, spec *pb.ChaincodeDeploymentSpec) error
}
