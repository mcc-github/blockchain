/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccprovider

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/golang/protobuf/proto"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/pkg/errors"
)

var ccproviderLogger = flogging.MustGetLogger("ccprovider")

var chaincodeInstallPath string






type CCPackage interface {
	
	InitFromBuffer(buf []byte) (*ChaincodeData, error)

	
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




func isPrintable(name string) bool {
	notASCII := func(r rune) bool {
		return !unicode.IsPrint(r)
	}
	return strings.IndexFunc(name, notASCII) == -1
}


func GetChaincodePackageFromPath(ccNameVersion string, ccInstallPath string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s", ccInstallPath, strings.ReplaceAll(ccNameVersion, ":", "."))
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
	
	GetChaincode(ccNameVersion string) (CCPackage, error)
}



type CCInfoFSImpl struct {
	GetHasher GetHasher
}



func (cifs *CCInfoFSImpl) GetChaincode(ccNameVersion string) (CCPackage, error) {
	return cifs.GetChaincodeFromPath(ccNameVersion, chaincodeInstallPath)
}

func (cifs *CCInfoFSImpl) GetChaincodeCodePackage(ccNameVersion string) ([]byte, error) {
	ccpack, err := cifs.GetChaincode(ccNameVersion)
	if err != nil {
		return nil, err
	}
	return ccpack.GetDepSpec().CodePackage, nil
}

func (cifs *CCInfoFSImpl) GetChaincodeDepSpec(ccNameVersion string) (*pb.ChaincodeDeploymentSpec, error) {
	ccpack, err := cifs.GetChaincode(ccNameVersion)
	if err != nil {
		return nil, err
	}
	return ccpack.GetDepSpec(), nil
}


func (cifs *CCInfoFSImpl) GetChaincodeFromPath(ccNameVersion string, path string) (CCPackage, error) {
	
	cccdspack := &CDSPackage{GetHasher: cifs.GetHasher}
	_, _, err := cccdspack.InitFromPath(ccNameVersion, path)
	if err != nil {
		
		ccscdspack := &SignedCDSPackage{GetHasher: cifs.GetHasher}
		_, _, err = ccscdspack.InitFromPath(ccNameVersion, path)
		if err != nil {
			return nil, err
		}
		return ccscdspack, nil
	}
	return cccdspack, nil
}


func (*CCInfoFSImpl) GetChaincodeInstallPath() string {
	return chaincodeInstallPath
}



func (cifs *CCInfoFSImpl) PutChaincode(depSpec *pb.ChaincodeDeploymentSpec) (CCPackage, error) {
	buf, err := proto.Marshal(depSpec)
	if err != nil {
		return nil, err
	}
	cccdspack := &CDSPackage{GetHasher: cifs.GetHasher}
	if _, err := cccdspack.InitFromBuffer(buf); err != nil {
		return nil, err
	}
	err = cccdspack.PutChaincodeToFS()
	if err != nil {
		return nil, err
	}

	return cccdspack, nil
}


type DirEnumerator func(string) ([]os.FileInfo, error)


type ChaincodeExtractor func(ccNameVersion string, path string, getHasher GetHasher) (CCPackage, error)


func (cifs *CCInfoFSImpl) ListInstalledChaincodes(dir string, ls DirEnumerator, ccFromPath ChaincodeExtractor) ([]chaincode.InstalledChaincode, error) {
	var chaincodes []chaincode.InstalledChaincode
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return nil, nil
	}
	files, err := ls(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading directory %s", dir)
	}

	for _, f := range files {
		
		if f.IsDir() {
			continue
		}
		
		
		
		i := strings.Index(f.Name(), ".")
		if i == -1 {
			ccproviderLogger.Info("Skipping", f.Name(), "because of missing separator '.'")
			continue
		}
		ccName := f.Name()[:i]      
		ccVersion := f.Name()[i+1:] 

		ccPackage, err := ccFromPath(ccName+":"+ccVersion, dir, cifs.GetHasher)
		if err != nil {
			ccproviderLogger.Warning("Failed obtaining chaincode information about", ccName, ccVersion, ":", err)
			return nil, errors.Wrapf(err, "failed obtaining information about %s, version %s", ccName, ccVersion)
		}

		chaincodes = append(chaincodes, chaincode.InstalledChaincode{
			Name:    ccName,
			Version: ccVersion,
			Hash:    ccPackage.GetId(),
		})
	}
	ccproviderLogger.Debug("Returning", chaincodes)
	return chaincodes, nil
}



var ccInfoFSProvider = &CCInfoFSImpl{GetHasher: factory.GetDefault()}


var ccInfoCache = NewCCInfoCache(ccInfoFSProvider)


func GetChaincodeFromFS(ccNameVersion string) (CCPackage, error) {
	return ccInfoFSProvider.GetChaincode(ccNameVersion)
}


func GetChaincodeData(ccNameVersion string) (*ChaincodeData, error) {
	ccproviderLogger.Debugf("Getting chaincode data for <%s> from cache", ccNameVersion)
	return ccInfoCache.GetChaincodeData(ccNameVersion)
}



func GetCCPackage(buf []byte, bccsp bccsp.BCCSP) (CCPackage, error) {
	
	cds := &CDSPackage{GetHasher: bccsp}
	if ccdata, err := cds.InitFromBuffer(buf); err != nil {
		cds = nil
	} else {
		err = cds.ValidateCC(ccdata)
		if err != nil {
			cds = nil
		}
	}

	
	scds := &SignedCDSPackage{GetHasher: bccsp}
	if ccdata, err := scds.InitFromBuffer(buf); err != nil {
		scds = nil
	} else {
		err = scds.ValidateCC(ccdata)
		if err != nil {
			scds = nil
		}
	}

	if cds != nil && scds != nil {
		
		
		ccproviderLogger.Errorf("Could not determine chaincode package type, guessing SignedCDS")
		return scds, nil
	}

	if cds != nil {
		return cds, nil
	}

	if scds != nil {
		return scds, nil
	}

	return nil, errors.New("could not unmarshal chaincode package to CDS or SignedCDS")
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
			ccpack, err := GetChaincodeFromFS(ccname + ":" + ccversion)
			if err != nil {
				
				
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


func (cd *ChaincodeData) ChaincodeID() string {
	return cd.Name + ":" + cd.Version
}




func (cd *ChaincodeData) Reset() { *cd = ChaincodeData{} }


func (cd *ChaincodeData) String() string { return proto.CompactTextString(cd) }


func (*ChaincodeData) ProtoMessage() {}



type TransactionParams struct {
	TxID                 string
	ChannelID            string
	NamespaceID          string
	SignedProp           *pb.SignedProposal
	Proposal             *pb.Proposal
	TXSimulator          ledger.TxSimulator
	HistoryQueryExecutor ledger.HistoryQueryExecutor
	CollectionStore      privdata.CollectionStore
	IsInitTransaction    bool

	
	ProposalDecorations map[string][]byte
}
