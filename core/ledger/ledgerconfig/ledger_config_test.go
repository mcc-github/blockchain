

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

func TestPvtdataStoreCollElgProcMaxDbBatchSize(t *testing.T) {
	defaultVal := confCollElgProcMaxDbBatchSize.DefaultVal
	testVal := defaultVal + 1
	assert.Equal(t, defaultVal, GetPvtdataStoreCollElgProcMaxDbBatchSize())
	viper.Set("ledger.pvtdataStore.collElgProcMaxDbBatchSize", testVal)
	assert.Equal(t, testVal, GetPvtdataStoreCollElgProcMaxDbBatchSize())
}

func TestCollElgProcDbBatchesInterval(t *testing.T) {
	defaultVal := confCollElgProcDbBatchesInterval.DefaultVal
	testVal := defaultVal + 1
	assert.Equal(t, defaultVal, GetPvtdataStoreCollElgProcDbBatchesInterval())
	viper.Set("ledger.pvtdataStore.collElgProcDbBatchesInterval", testVal)
	assert.Equal(t, testVal, GetPvtdataStoreCollElgProcDbBatchesInterval())
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

func TestGetMaxBlockfileSize(t *testing.T) {
	assert.Equal(t, 67108864, GetMaxBlockfileSize())
}

func setUpCoreYAMLConfig() {
	
	ledgertestutil.SetupCoreYAMLConfig()
}
