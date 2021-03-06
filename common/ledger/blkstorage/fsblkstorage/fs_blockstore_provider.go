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

package fsblkstorage

import (
	"github.com/mcc-github/blockchain/common/ledger/blkstorage"
	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
	"github.com/mcc-github/blockchain/common/metrics"
)

const (
	dataFormatVersion20 = "2.0"
)

func dataFormatVersion(indexConfig *blkstorage.IndexConfig) string {
	
	if indexConfig.Contains(blkstorage.IndexableAttrTxID) {
		return dataFormatVersion20
	}
	return ""
}


type FsBlockstoreProvider struct {
	conf            *Conf
	indexConfig     *blkstorage.IndexConfig
	leveldbProvider *leveldbhelper.Provider
	stats           *stats
}


func NewProvider(conf *Conf, indexConfig *blkstorage.IndexConfig, metricsProvider metrics.Provider) (blkstorage.BlockStoreProvider, error) {
	dbConf := &leveldbhelper.Conf{
		DBPath:                conf.getIndexDir(),
		ExpectedFormatVersion: dataFormatVersion(indexConfig),
	}

	p, err := leveldbhelper.NewProvider(dbConf)
	if err != nil {
		return nil, err
	}
	
	stats := newStats(metricsProvider)
	return &FsBlockstoreProvider{conf, indexConfig, p, stats}, nil
}


func (p *FsBlockstoreProvider) CreateBlockStore(ledgerid string) (blkstorage.BlockStore, error) {
	return p.OpenBlockStore(ledgerid)
}




func (p *FsBlockstoreProvider) OpenBlockStore(ledgerid string) (blkstorage.BlockStore, error) {
	indexStoreHandle := p.leveldbProvider.GetDBHandle(ledgerid)
	return newFsBlockStore(ledgerid, p.conf, p.indexConfig, indexStoreHandle, p.stats), nil
}


func (p *FsBlockstoreProvider) Exists(ledgerid string) (bool, error) {
	exists, _, err := util.FileExists(p.conf.getLedgerBlockDir(ledgerid))
	return exists, err
}


func (p *FsBlockstoreProvider) List() ([]string, error) {
	return util.ListSubdirs(p.conf.getChainsDir())
}


func (p *FsBlockstoreProvider) Close() {
	p.leveldbProvider.Close()
}
