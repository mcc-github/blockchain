

package endorser

import (
	"fmt"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/mcc-github/blockchain-protos-go/transientstore"
	"github.com/mcc-github/blockchain/core/ledger"
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





func AssemblePvtRWSet(channelName string,
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
				return nil, errors.WithMessagef(err, "error while retrieving collection config for chaincode %#v", namespace)
			}
			colCP := ccInfo.AllCollectionsConfigPkg()
			if colCP == nil {
				return nil, errors.New(fmt.Sprintf("no collection config for chaincode %#v", namespace))
			}
			txPvtRwSetWithConfig.CollectionConfigs[namespace] = colCP
		}
	}
	trimCollectionConfigs(txPvtRwSetWithConfig)
	return txPvtRwSetWithConfig, nil
}

func trimCollectionConfigs(pvtData *transientstore.TxPvtReadWriteSetWithConfigInfo) {
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
