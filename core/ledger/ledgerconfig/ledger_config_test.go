

package ledgerconfig

import (
	"testing"

	ledgertestutil "github.com/mcc-github/blockchain/core/ledger/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestIsCouchDBEnabledDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	

	
	
	if IsCouchDBEnabled() == true {
		ledgertestutil.ResetConfigToDefaultValues()
		defer viper.Set("ledger.state.stateDatabase", "CouchDB")
	}
	defaultValue := IsCouchDBEnabled()
	assert.False(t, defaultValue) 
}

func TestIsCouchDBEnabled(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.stateDatabase", "CouchDB")
	updatedValue := IsCouchDBEnabled()
	assert.True(t, updatedValue) 
}

func TestLedgerConfigPathDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	assert.Equal(t, "/var/mcc-github/production/ledgersData", GetRootPath())
	assert.Equal(t, "/var/mcc-github/production/ledgersData/ledgerProvider", GetLedgerProviderPath())
	assert.Equal(t, "/var/mcc-github/production/ledgersData/stateLeveldb", GetStateLevelDBPath())
	assert.Equal(t, "/var/mcc-github/production/ledgersData/historyLeveldb", GetHistoryLevelDBPath())
	assert.Equal(t, "/var/mcc-github/production/ledgersData/chains", GetBlockStorePath())
	assert.Equal(t, "/var/mcc-github/production/ledgersData/pvtdataStore", GetPvtdataStorePath())
	assert.Equal(t, "/var/mcc-github/production/ledgersData/bookkeeper", GetInternalBookkeeperPath())
}

func TestLedgerConfigPath(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("peer.fileSystemPath", "/tmp/mcc-github/production")
	assert.Equal(t, "/tmp/mcc-github/production/ledgersData", GetRootPath())
	assert.Equal(t, "/tmp/mcc-github/production/ledgersData/ledgerProvider", GetLedgerProviderPath())
	assert.Equal(t, "/tmp/mcc-github/production/ledgersData/stateLeveldb", GetStateLevelDBPath())
	assert.Equal(t, "/tmp/mcc-github/production/ledgersData/historyLeveldb", GetHistoryLevelDBPath())
	assert.Equal(t, "/tmp/mcc-github/production/ledgersData/chains", GetBlockStorePath())
	assert.Equal(t, "/tmp/mcc-github/production/ledgersData/pvtdataStore", GetPvtdataStorePath())
	assert.Equal(t, "/tmp/mcc-github/production/ledgersData/bookkeeper", GetInternalBookkeeperPath())
}

func TestGetTotalLimitDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetTotalQueryLimit()
	assert.Equal(t, 10000, defaultValue) 
}

func TestGetTotalLimitUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetTotalQueryLimit()
	assert.Equal(t, 10000, defaultValue) 
}

func TestGetTotalLimit(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.totalQueryLimit", 5000)
	updatedValue := GetTotalQueryLimit()
	assert.Equal(t, 5000, updatedValue) 
}

func TestGetQueryLimitDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetInternalQueryLimit()
	assert.Equal(t, 1000, defaultValue) 
}

func TestGetQueryLimitUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetInternalQueryLimit()
	assert.Equal(t, 1000, defaultValue) 
}

func TestGetQueryLimit(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.internalQueryLimit", 5000)
	updatedValue := GetInternalQueryLimit()
	assert.Equal(t, 5000, updatedValue) 
}

func TestMaxBatchUpdateSizeDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetMaxBatchUpdateSize()
	assert.Equal(t, 1000, defaultValue) 
}

func TestMaxBatchUpdateSizeUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetMaxBatchUpdateSize()
	assert.Equal(t, 500, defaultValue) 
}

func TestMaxBatchUpdateSize(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.maxBatchUpdateSize", 2000)
	updatedValue := GetMaxBatchUpdateSize()
	assert.Equal(t, 2000, updatedValue) 
}

func TestPvtdataStorePurgeIntervalDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetPvtdataStorePurgeInterval()
	assert.Equal(t, uint64(100), defaultValue) 
}

func TestPvtdataStorePurgeIntervalUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetPvtdataStorePurgeInterval()
	assert.Equal(t, uint64(100), defaultValue) 
}

func TestIsQueryReadHasingEnabled(t *testing.T) {
	assert.True(t, IsQueryReadsHashingEnabled())
}

func TestGetMaxDegreeQueryReadsHashing(t *testing.T) {
	assert.Equal(t, uint32(50), GetMaxDegreeQueryReadsHashing())
}

func TestPvtdataStorePurgeInterval(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.pvtdataStore.purgeInterval", 1000)
	updatedValue := GetPvtdataStorePurgeInterval()
	assert.Equal(t, uint64(1000), updatedValue) 
}

func TestIsHistoryDBEnabledDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := IsHistoryDBEnabled()
	assert.False(t, defaultValue) 
}

func TestIsHistoryDBEnabledTrue(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.history.enableHistoryDatabase", true)
	updatedValue := IsHistoryDBEnabled()
	assert.True(t, updatedValue) 
}

func TestIsHistoryDBEnabledFalse(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.history.enableHistoryDatabase", false)
	updatedValue := IsHistoryDBEnabled()
	assert.False(t, updatedValue) 
}

func TestIsAutoWarmIndexesEnabledDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := IsAutoWarmIndexesEnabled()
	assert.True(t, defaultValue) 
}

func TestIsAutoWarmIndexesEnabledUnset(t *testing.T) {
	viper.Reset()
	defaultValue := IsAutoWarmIndexesEnabled()
	assert.True(t, defaultValue) 
}

func TestIsAutoWarmIndexesEnabledTrue(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.autoWarmIndexes", true)
	updatedValue := IsAutoWarmIndexesEnabled()
	assert.True(t, updatedValue) 
}

func TestIsAutoWarmIndexesEnabledFalse(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.autoWarmIndexes", false)
	updatedValue := IsAutoWarmIndexesEnabled()
	assert.False(t, updatedValue) 
}

func TestGetWarmIndexesAfterNBlocksDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetWarmIndexesAfterNBlocks()
	assert.Equal(t, 1, defaultValue) 
}

func TestGetWarmIndexesAfterNBlocksUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetWarmIndexesAfterNBlocks()
	assert.Equal(t, 1, defaultValue) 
}

func TestGetWarmIndexesAfterNBlocks(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.warmIndexesAfterNBlocks", 10)
	updatedValue := GetWarmIndexesAfterNBlocks()
	assert.Equal(t, 10, updatedValue)
}

func TestGetMaxBlockfileSize(t *testing.T) {
	assert.Equal(t, 67108864, GetMaxBlockfileSize())
}

func setUpCoreYAMLConfig() {
	
	ledgertestutil.SetupCoreYAMLConfig()
}
