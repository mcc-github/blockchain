

package ledgerconfig

import (
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	ledgertestutil "github.com/mcc-github/blockchain/core/ledger/testutil"
	"github.com/spf13/viper"
)

func TestIsCouchDBEnabledDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	

	
	
	if IsCouchDBEnabled() == true {
		ledgertestutil.ResetConfigToDefaultValues()
		defer viper.Set("ledger.state.stateDatabase", "CouchDB")
	}
	defaultValue := IsCouchDBEnabled()
	testutil.AssertEquals(t, defaultValue, false) 
}

func TestIsCouchDBEnabled(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.stateDatabase", "CouchDB")
	updatedValue := IsCouchDBEnabled()
	testutil.AssertEquals(t, updatedValue, true) 
}

func TestLedgerConfigPathDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	testutil.AssertEquals(t,
		GetRootPath(),
		"/var/mcc-github/production/ledgersData")
	testutil.AssertEquals(t,
		GetLedgerProviderPath(),
		"/var/mcc-github/production/ledgersData/ledgerProvider")
	testutil.AssertEquals(t,
		GetStateLevelDBPath(),
		"/var/mcc-github/production/ledgersData/stateLeveldb")
	testutil.AssertEquals(t,
		GetHistoryLevelDBPath(),
		"/var/mcc-github/production/ledgersData/historyLeveldb")
	testutil.AssertEquals(t,
		GetBlockStorePath(),
		"/var/mcc-github/production/ledgersData/chains")
	testutil.AssertEquals(t,
		GetPvtdataStorePath(),
		"/var/mcc-github/production/ledgersData/pvtdataStore")
	testutil.AssertEquals(t,
		GetInternalBookkeeperPath(),
		"/var/mcc-github/production/ledgersData/bookkeeper")

}

func TestLedgerConfigPath(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("peer.fileSystemPath", "/tmp/mcc-github/production")
	testutil.AssertEquals(t,
		GetRootPath(),
		"/tmp/mcc-github/production/ledgersData")
	testutil.AssertEquals(t,
		GetLedgerProviderPath(),
		"/tmp/mcc-github/production/ledgersData/ledgerProvider")
	testutil.AssertEquals(t,
		GetStateLevelDBPath(),
		"/tmp/mcc-github/production/ledgersData/stateLeveldb")
	testutil.AssertEquals(t,
		GetHistoryLevelDBPath(),
		"/tmp/mcc-github/production/ledgersData/historyLeveldb")
	testutil.AssertEquals(t,
		GetBlockStorePath(),
		"/tmp/mcc-github/production/ledgersData/chains")
	testutil.AssertEquals(t,
		GetPvtdataStorePath(),
		"/tmp/mcc-github/production/ledgersData/pvtdataStore")
	testutil.AssertEquals(t,
		GetInternalBookkeeperPath(),
		"/tmp/mcc-github/production/ledgersData/bookkeeper")
}

func TestGetQueryLimitDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetQueryLimit()
	testutil.AssertEquals(t, defaultValue, 10000) 
}

func TestGetQueryLimitUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetQueryLimit()
	testutil.AssertEquals(t, defaultValue, 10000) 
}

func TestGetQueryLimit(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.queryLimit", 5000)
	updatedValue := GetQueryLimit()
	testutil.AssertEquals(t, updatedValue, 5000) 
}

func TestMaxBatchUpdateSizeDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetMaxBatchUpdateSize()
	testutil.AssertEquals(t, defaultValue, 1000) 
}

func TestMaxBatchUpdateSizeUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetMaxBatchUpdateSize()
	testutil.AssertEquals(t, defaultValue, 500) 
}

func TestMaxBatchUpdateSize(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.maxBatchUpdateSize", 2000)
	updatedValue := GetMaxBatchUpdateSize()
	testutil.AssertEquals(t, updatedValue, 2000) 
}

func TestPvtdataStorePurgeIntervalDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetPvtdataStorePurgeInterval()
	testutil.AssertEquals(t, defaultValue, uint64(100)) 
}

func TestPvtdataStorePurgeIntervalUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetPvtdataStorePurgeInterval()
	testutil.AssertEquals(t, defaultValue, uint64(100)) 
}

func TestIsQueryReadHasingEnabled(t *testing.T) {
	testutil.AssertEquals(t, IsQueryReadsHashingEnabled(), true)
}

func TestGetMaxDegreeQueryReadsHashing(t *testing.T) {
	testutil.AssertEquals(t, GetMaxDegreeQueryReadsHashing(), uint32(50))
}

func TestPvtdataStorePurgeInterval(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.pvtdataStore.purgeInterval", 1000)
	updatedValue := GetPvtdataStorePurgeInterval()
	testutil.AssertEquals(t, updatedValue, uint64(1000)) 
}

func TestIsHistoryDBEnabledDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := IsHistoryDBEnabled()
	testutil.AssertEquals(t, defaultValue, false) 
}

func TestIsHistoryDBEnabledTrue(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.history.enableHistoryDatabase", true)
	updatedValue := IsHistoryDBEnabled()
	testutil.AssertEquals(t, updatedValue, true) 
}

func TestIsHistoryDBEnabledFalse(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.history.enableHistoryDatabase", false)
	updatedValue := IsHistoryDBEnabled()
	testutil.AssertEquals(t, updatedValue, false) 
}

func TestIsAutoWarmIndexesEnabledDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := IsAutoWarmIndexesEnabled()
	testutil.AssertEquals(t, defaultValue, true) 
}

func TestIsAutoWarmIndexesEnabledUnset(t *testing.T) {
	viper.Reset()
	defaultValue := IsAutoWarmIndexesEnabled()
	testutil.AssertEquals(t, defaultValue, true) 
}

func TestIsAutoWarmIndexesEnabledTrue(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.autoWarmIndexes", true)
	updatedValue := IsAutoWarmIndexesEnabled()
	testutil.AssertEquals(t, updatedValue, true) 
}

func TestIsAutoWarmIndexesEnabledFalse(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.autoWarmIndexes", false)
	updatedValue := IsAutoWarmIndexesEnabled()
	testutil.AssertEquals(t, updatedValue, false) 
}

func TestGetWarmIndexesAfterNBlocksDefault(t *testing.T) {
	setUpCoreYAMLConfig()
	defaultValue := GetWarmIndexesAfterNBlocks()
	testutil.AssertEquals(t, defaultValue, 1) 
}

func TestGetWarmIndexesAfterNBlocksUnset(t *testing.T) {
	viper.Reset()
	defaultValue := GetWarmIndexesAfterNBlocks()
	testutil.AssertEquals(t, defaultValue, 1) 
}

func TestGetWarmIndexesAfterNBlocks(t *testing.T) {
	setUpCoreYAMLConfig()
	defer ledgertestutil.ResetConfigToDefaultValues()
	viper.Set("ledger.state.couchDBConfig.warmIndexesAfterNBlocks", 10)
	updatedValue := GetWarmIndexesAfterNBlocks()
	testutil.AssertEquals(t, updatedValue, 10)
}

func TestGetMaxBlockfileSize(t *testing.T) {
	testutil.AssertEquals(t, GetMaxBlockfileSize(), 67108864)
}

func setUpCoreYAMLConfig() {
	
	ledgertestutil.SetupCoreYAMLConfig()
}
