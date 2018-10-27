/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

package file

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/orderer/common/bootstrap"
	cb "github.com/mcc-github/blockchain/protos/common"
)

type fileBootstrapper struct {
	GenesisBlockFile string
}


func New(fileName string) bootstrap.Helper {
	return &fileBootstrapper{
		GenesisBlockFile: fileName,
	}
}


func (b *fileBootstrapper) GenesisBlock() *cb.Block {
	bootstrapFile, fileErr := ioutil.ReadFile(b.GenesisBlockFile)
	if fileErr != nil {
		panic(fmt.Errorf("Unable to bootstrap orderer. Error reading genesis block file: %v", fileErr))
	}
	genesisBlock := &cb.Block{}
	unmarshallErr := proto.Unmarshal(bootstrapFile, genesisBlock)
	if unmarshallErr != nil {
		panic(fmt.Errorf("Unable to bootstrap orderer. Error unmarshalling genesis block: %v", unmarshallErr))

	}
	return genesisBlock
} 
