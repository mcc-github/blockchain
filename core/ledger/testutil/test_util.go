/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	"flag"
	"fmt"
	mathRand "math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)


type TestRandomNumberGenerator struct {
	rand      *mathRand.Rand
	maxNumber int
}


func NewTestRandomNumberGenerator(maxNumber int) *TestRandomNumberGenerator {
	return &TestRandomNumberGenerator{
		mathRand.New(mathRand.NewSource(time.Now().UnixNano())),
		maxNumber,
	}
}


func (randNumGenerator *TestRandomNumberGenerator) Next() int {
	return randNumGenerator.rand.Intn(randNumGenerator.maxNumber)
}


func SetupTestConfig() {
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("peer.ledger.test.loadYAML", true)
	loadYAML := viper.GetBool("peer.ledger.test.loadYAML")
	if loadYAML {
		viper.SetConfigName("test")
		err := viper.ReadInConfig()
		if err != nil { 
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	var formatter = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} [%{module}] %{shortfunc} [%{shortfile}] -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(formatter)
}


func SetupCoreYAMLConfig() {
	viper.SetConfigName("core")
	viper.SetEnvPrefix("CORE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := configtest.AddDevConfigPath(nil)
	if err != nil {
		panic(fmt.Errorf("Fatal error adding dev dir: %s \n", err))
	}

	err = viper.ReadInConfig()
	if err != nil { 
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}


func ResetConfigToDefaultValues() {
	
	viper.Set("ledger.state.couchDBConfig.queryLimit", 10000)
	viper.Set("ledger.state.stateDatabase", "goleveldb")
	viper.Set("ledger.history.enableHistoryDatabase", false)
	viper.Set("ledger.state.couchDBConfig.autoWarmIndexes", true)
	viper.Set("ledger.state.couchDBConfig.warmIndexesAfterNBlocks", 1)
	viper.Set("peer.fileSystemPath", "/var/mcc-github/production")
}


func SetLogLevel(level logging.Level, module string) {
	logging.SetLevel(level, module)
}


func ParseTestParams() []string {
	testParams := flag.String("testParams", "", "Test specific parameters")
	flag.Parse()
	regex, err := regexp.Compile(",(\\s+)?")
	if err != nil {
		panic(fmt.Errorf("err = %s\n", err))
	}
	paramsArray := regex.Split(*testParams, -1)
	return paramsArray
}
