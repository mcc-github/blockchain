/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ccprovider

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	pb "github.com/mcc-github/blockchain/protos/peer"
)






type CDSData struct {
	
	CodeHash []byte `protobuf:"bytes,1,opt,name=codehash,proto3"`

	
	MetaDataHash []byte `protobuf:"bytes,2,opt,name=metadatahash,proto3"`
}




func (data *CDSData) Reset() { *data = CDSData{} }


func (data *CDSData) String() string { return proto.CompactTextString(data) }


func (*CDSData) ProtoMessage() {}


func (data *CDSData) Equals(other *CDSData) bool {
	return other != nil && bytes.Equal(data.CodeHash, other.CodeHash) && bytes.Equal(data.MetaDataHash, other.MetaDataHash)
}




type CDSPackage struct {
	buf     []byte
	depSpec *pb.ChaincodeDeploymentSpec
	data    *CDSData
	datab   []byte
	id      []byte
}


func (ccpack *CDSPackage) reset() {
	*ccpack = CDSPackage{}
}


func (ccpack *CDSPackage) GetId() []byte {
	
	
	if ccpack.id == nil {
		panic("GetId called on uninitialized package")
	}
	return ccpack.id
}


func (ccpack *CDSPackage) GetDepSpec() *pb.ChaincodeDeploymentSpec {
	
	
	if ccpack.depSpec == nil {
		panic("GetDepSpec called on uninitialized package")
	}
	return ccpack.depSpec
}


func (ccpack *CDSPackage) GetDepSpecBytes() []byte {
	
	
	if ccpack.buf == nil {
		panic("GetDepSpecBytes called on uninitialized package")
	}
	return ccpack.buf
}


func (ccpack *CDSPackage) GetPackageObject() proto.Message {
	return ccpack.depSpec
}


func (ccpack *CDSPackage) GetChaincodeData() *ChaincodeData {
	
	
	if ccpack.depSpec == nil || ccpack.datab == nil || ccpack.id == nil {
		panic("GetChaincodeData called on uninitialized package")
	}
	return &ChaincodeData{Name: ccpack.depSpec.ChaincodeSpec.ChaincodeId.Name, Version: ccpack.depSpec.ChaincodeSpec.ChaincodeId.Version, Data: ccpack.datab, Id: ccpack.id}
}

func (ccpack *CDSPackage) getCDSData(cds *pb.ChaincodeDeploymentSpec) ([]byte, []byte, *CDSData, error) {
	
	
	
	if cds == nil {
		panic("nil cds")
	}

	b, err := proto.Marshal(cds)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = factory.InitFactories(nil); err != nil {
		return nil, nil, nil, fmt.Errorf("Internal error, BCCSP could not be initialized : %s", err)
	}

	
	hash, err := factory.GetDefault().GetHash(&bccsp.SHAOpts{})
	if err != nil {
		return nil, nil, nil, err
	}

	cdsdata := &CDSData{}

	
	hash.Write(cds.CodePackage)
	cdsdata.CodeHash = hash.Sum(nil)

	hash.Reset()

	
	hash.Write([]byte(cds.ChaincodeSpec.ChaincodeId.Name))
	hash.Write([]byte(cds.ChaincodeSpec.ChaincodeId.Version))

	cdsdata.MetaDataHash = hash.Sum(nil)

	b, err = proto.Marshal(cdsdata)
	if err != nil {
		return nil, nil, nil, err
	}

	hash.Reset()

	
	hash.Write(cdsdata.CodeHash)
	hash.Write(cdsdata.MetaDataHash)

	id := hash.Sum(nil)

	return b, id, cdsdata, nil
}



func (ccpack *CDSPackage) ValidateCC(ccdata *ChaincodeData) error {
	if ccpack.depSpec == nil {
		return fmt.Errorf("uninitialized package")
	}

	if ccpack.data == nil {
		return fmt.Errorf("nil data")
	}

	
	
	
	
	
	
	if !isPrintable(ccdata.Name) {
		return fmt.Errorf("invalid chaincode name: %q", ccdata.Name)
	}

	if ccdata.Name != ccpack.depSpec.ChaincodeSpec.ChaincodeId.Name || ccdata.Version != ccpack.depSpec.ChaincodeSpec.ChaincodeId.Version {
		return fmt.Errorf("invalid chaincode data %v (%v)", ccdata, ccpack.depSpec.ChaincodeSpec.ChaincodeId)
	}

	otherdata := &CDSData{}
	err := proto.Unmarshal(ccdata.Data, otherdata)
	if err != nil {
		return err
	}

	if !ccpack.data.Equals(otherdata) {
		return fmt.Errorf("data mismatch")
	}

	return nil
}


func (ccpack *CDSPackage) InitFromBuffer(buf []byte) (*ChaincodeData, error) {
	
	ccpack.reset()

	depSpec := &pb.ChaincodeDeploymentSpec{}
	err := proto.Unmarshal(buf, depSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deployment spec from bytes")
	}

	databytes, id, data, err := ccpack.getCDSData(depSpec)
	if err != nil {
		return nil, err
	}

	ccpack.buf = buf
	ccpack.depSpec = depSpec
	ccpack.data = data
	ccpack.datab = databytes
	ccpack.id = id

	return ccpack.GetChaincodeData(), nil
}


func (ccpack *CDSPackage) InitFromPath(ccNameVersion string, path string) ([]byte, *pb.ChaincodeDeploymentSpec, error) {
	
	ccpack.reset()

	buf, err := GetChaincodePackageFromPath(ccNameVersion, path)
	if err != nil {
		return nil, nil, err
	}

	ccdata, err := ccpack.InitFromBuffer(buf)
	if err != nil {
		return nil, nil, err
	}

	if err := ccpack.ValidateCC(ccdata); err != nil {
		return nil, nil, err
	}

	return ccpack.buf, ccpack.depSpec, nil
}


func (ccpack *CDSPackage) InitFromFS(ccNameVersion string) ([]byte, *pb.ChaincodeDeploymentSpec, error) {
	return ccpack.InitFromPath(ccNameVersion, chaincodeInstallPath)
}


func (ccpack *CDSPackage) PutChaincodeToFS() error {
	if ccpack.buf == nil {
		return fmt.Errorf("uninitialized package")
	}

	if ccpack.id == nil {
		return fmt.Errorf("id cannot be nil if buf is not nil")
	}

	if ccpack.depSpec == nil {
		return fmt.Errorf("depspec cannot be nil if buf is not nil")
	}

	if ccpack.data == nil {
		return fmt.Errorf("nil data")
	}

	if ccpack.datab == nil {
		return fmt.Errorf("nil data bytes")
	}

	ccname := ccpack.depSpec.ChaincodeSpec.ChaincodeId.Name
	ccversion := ccpack.depSpec.ChaincodeSpec.ChaincodeId.Version

	
	path := fmt.Sprintf("%s/%s.%s", chaincodeInstallPath, ccname, ccversion)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("chaincode %s exists", path)
	}

	if err := ioutil.WriteFile(path, ccpack.buf, 0644); err != nil {
		return err
	}

	return nil
}
