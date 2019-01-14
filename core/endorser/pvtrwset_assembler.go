

package endorser

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/pkg/errors"
)



type PvtRWSetAssembler interface {
	
	
	
	
	AssemblePvtRWSet(channelName string,
		privData *rwset.TxPvtReadWriteSet,
		txsim ledger.SimpleQueryExecutor,
		deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider) (
		*transientstore.TxPvtReadWriteSetWithConfigInfo, error,
	)
}



type CollectionConfigRetriever interface {
	
	GetState(namespace string, key string) ([]byte, error)
}

type rwSetAssembler struct {
}





func (as *rwSetAssembler) AssemblePvtRWSet(channelName string,
	privData *rwset.TxPvtReadWriteSet,
	txsim ledger.SimpleQueryExecutor,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider) (
	*transientstore.TxPvtReadWriteSetWithConfigInfo, error,
) {
	txPvtRwSetWithConfig := &transientstore.TxPvtReadWriteSetWithConfigInfo{
		PvtRwset:          privData,
		CollectionConfigs: make(map[string]*common.CollectionConfigPackage),
	}

	for _, pvtRwset := range privData.NsPvtRwset {
		namespace := pvtRwset.Namespace
		if _, found := txPvtRwSetWithConfig.CollectionConfigs[namespace]; !found {
			ccInfo, err := deployedCCInfoProvider.ChaincodeInfo(channelName, namespace, txsim)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf("error while retrieving collection config for chaincode %#v", namespace))
			}
			colCP := ccInfo.CollectionConfigPkg
			if colCP == nil {
				return nil, errors.New(fmt.Sprintf("no collection config for chaincode %#v", namespace))
			}
			txPvtRwSetWithConfig.CollectionConfigs[namespace] = colCP
		}
	}
	as.trimCollectionConfigs(txPvtRwSetWithConfig)
	return txPvtRwSetWithConfig, nil
}

func (as *rwSetAssembler) trimCollectionConfigs(pvtData *transientstore.TxPvtReadWriteSetWithConfigInfo) {
	flags := make(map[string]map[string]struct{})
	for _, pvtRWset := range pvtData.PvtRwset.NsPvtRwset {
		namespace := pvtRWset.Namespace
		for _, col := range pvtRWset.CollectionPvtRwset {
			if _, found := flags[namespace]; !found {
				flags[namespace] = make(map[string]struct{})
			}
			flags[namespace][col.CollectionName] = struct{}{}
		}
	}

	filteredConfigs := make(map[string]*common.CollectionConfigPackage)
	for namespace, configs := range pvtData.CollectionConfigs {
		filteredConfigs[namespace] = &common.CollectionConfigPackage{}
		for _, conf := range configs.Config {
			if colConf := conf.GetStaticCollectionConfig(); colConf != nil {
				if _, found := flags[namespace][colConf.Name]; found {
					filteredConfigs[namespace].Config = append(filteredConfigs[namespace].Config, conf)
				}
			}
		}
	}
	pvtData.CollectionConfigs = filteredConfigs
}
