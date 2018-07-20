/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package localconfig

import (
	"fmt"

	"strings"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/viperutil"
	logging "github.com/op/go-logging"

	"github.com/spf13/viper"

	"path/filepath"

	cf "github.com/mcc-github/blockchain/core/config"
	"github.com/mcc-github/blockchain/msp"
)

const (
	pkgLogID = "common/tools/configtxgen/localconfig"

	
	Prefix string = "CONFIGTX"
)

var (
	logger *logging.Logger

	configName string
)

func init() {
	logger = flogging.MustGetLogger(pkgLogID)
	flogging.SetModuleLevel(pkgLogID, "error")

	configName = strings.ToLower(Prefix)
}

const (
	
	
	TestChainID = "testchainid"

	
	
	SampleInsecureSoloProfile = "SampleInsecureSolo"
	
	
	SampleDevModeSoloProfile = "SampleDevModeSolo"
	
	
	SampleSingleMSPSoloProfile = "SampleSingleMSPSolo"

	
	
	SampleInsecureKafkaProfile = "SampleInsecureKafka"
	
	
	SampleDevModeKafkaProfile = "SampleDevModeKafka"
	
	
	SampleSingleMSPKafkaProfile = "SampleSingleMSPKafka"

	
	
	SampleSingleMSPChannelProfile = "SampleSingleMSPChannel"

	
	
	SampleConsortiumName = "SampleConsortium"
	
	SampleOrgName = "SampleOrg"

	
	
	AdminRoleAdminPrincipal = "Role.ADMIN"
	
	
	MemberRoleAdminPrincipal = "Role.MEMBER"
)


type TopLevel struct {
	Profiles      map[string]*Profile        `yaml:"Profiles"`
	Organizations []*Organization            `yaml:"Organizations"`
	Channel       *Profile                   `yaml:"Channel"`
	Application   *Application               `yaml:"Application"`
	Orderer       *Orderer                   `yaml:"Orderer"`
	Capabilities  map[string]map[string]bool `yaml:"Capabilities"`
	Resources     *Resources                 `yaml:"Resources"`
}



type Profile struct {
	Consortium   string                 `yaml:"Consortium"`
	Application  *Application           `yaml:"Application"`
	Orderer      *Orderer               `yaml:"Orderer"`
	Consortiums  map[string]*Consortium `yaml:"Consortiums"`
	Capabilities map[string]bool        `yaml:"Capabilities"`
	Policies     map[string]*Policy     `yaml:"Policies"`
}


type Policy struct {
	Type string `yaml:"Type"`
	Rule string `yaml:"Rule"`
}



type Consortium struct {
	Organizations []*Organization `yaml:"Organizations"`
}



type Application struct {
	Organizations []*Organization    `yaml:"Organizations"`
	Capabilities  map[string]bool    `yaml:"Capabilities"`
	Resources     *Resources         `yaml:"Resources"`
	Policies      map[string]*Policy `yaml:"Policies"`
	ACLs          map[string]string  `yaml:"ACLs"`
}



type Resources struct {
	DefaultModPolicy string
}



type Organization struct {
	Name     string             `yaml:"Name"`
	ID       string             `yaml:"ID"`
	MSPDir   string             `yaml:"MSPDir"`
	MSPType  string             `yaml:"MSPType"`
	Policies map[string]*Policy `yaml:"Policies"`

	
	
	
	AnchorPeers []*AnchorPeer `yaml:"AnchorPeers"`

	
	
	
	AdminPrincipal string `yaml:"AdminPrincipal"`
}


type AnchorPeer struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
}



type Orderer struct {
	OrdererType   string             `yaml:"OrdererType"`
	Addresses     []string           `yaml:"Addresses"`
	BatchTimeout  time.Duration      `yaml:"BatchTimeout"`
	BatchSize     BatchSize          `yaml:"BatchSize"`
	Kafka         Kafka              `yaml:"Kafka"`
	Organizations []*Organization    `yaml:"Organizations"`
	MaxChannels   uint64             `yaml:"MaxChannels"`
	Capabilities  map[string]bool    `yaml:"Capabilities"`
	Policies      map[string]*Policy `yaml:"Policies"`
}


type BatchSize struct {
	MaxMessageCount   uint32 `yaml:"MaxMessageCount"`
	AbsoluteMaxBytes  uint32 `yaml:"AbsoluteMaxBytes"`
	PreferredMaxBytes uint32 `yaml:"PreferredMaxBytes"`
}


type Kafka struct {
	Brokers []string `yaml:"Brokers"`
}

var genesisDefaults = TopLevel{
	Orderer: &Orderer{
		OrdererType:  "solo",
		Addresses:    []string{"127.0.0.1:7050"},
		BatchTimeout: 2 * time.Second,
		BatchSize: BatchSize{
			MaxMessageCount:   10,
			AbsoluteMaxBytes:  10 * 1024 * 1024,
			PreferredMaxBytes: 512 * 1024,
		},
		Kafka: Kafka{
			Brokers: []string{"127.0.0.1:9092"},
		},
	},
}







func LoadTopLevel(configPaths ...string) *TopLevel {
	config := viper.New()
	if len(configPaths) > 0 {
		for _, p := range configPaths {
			config.AddConfigPath(p)
		}
		config.SetConfigName(configName)
	} else {
		cf.InitViper(config, configName)
	}

	
	config.SetEnvPrefix(Prefix)
	config.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)

	err := config.ReadInConfig()
	if err != nil {
		logger.Panic("Error reading configuration: ", err)
	}
	logger.Debugf("Using config file: %s", config.ConfigFileUsed())

	var uconf TopLevel
	err = viperutil.EnhancedExactUnmarshal(config, &uconf)
	if err != nil {
		logger.Panic("Error unmarshaling config into struct: ", err)
	}

	(&uconf).completeInitialization(filepath.Dir(config.ConfigFileUsed()))

	logger.Infof("Loaded configuration: %s", config.ConfigFileUsed())

	return &uconf
}




func Load(profile string, configPaths ...string) *Profile {
	config := viper.New()
	if len(configPaths) > 0 {
		for _, p := range configPaths {
			config.AddConfigPath(p)
		}
		config.SetConfigName(configName)
	} else {
		cf.InitViper(config, configName)
	}

	
	config.SetEnvPrefix(Prefix)
	config.AutomaticEnv()

	
	
	replacer := strings.NewReplacer(strings.ToUpper(fmt.Sprintf("profiles.%s.", profile)), "", ".", "_")
	config.SetEnvKeyReplacer(replacer)

	err := config.ReadInConfig()
	if err != nil {
		logger.Panic("Error reading configuration: ", err)
	}
	logger.Debugf("Using config file: %s", config.ConfigFileUsed())

	var uconf TopLevel
	err = viperutil.EnhancedExactUnmarshal(config, &uconf)
	if err != nil {
		logger.Panic("Error unmarshaling config into struct: ", err)
	}

	result, ok := uconf.Profiles[profile]
	if !ok {
		logger.Panic("Could not find profile: ", profile)
	}

	result.completeInitialization(filepath.Dir(config.ConfigFileUsed()))

	logger.Infof("Loaded configuration: %s", config.ConfigFileUsed())

	return result
}

func (t *TopLevel) completeInitialization(configDir string) {
	for _, org := range t.Organizations {
		org.completeInitialization(configDir)
	}

	if t.Orderer != nil {
		t.Orderer.completeInitialization()
	}
}

func (p *Profile) completeInitialization(configDir string) {
	if p.Orderer != nil {
		for _, org := range p.Orderer.Organizations {
			org.completeInitialization(configDir)
		}
	}

	if p.Application != nil {
		for _, org := range p.Application.Organizations {
			org.completeInitialization(configDir)
		}
		if p.Application.Resources != nil {
			p.Application.Resources.completeInitialization()
		}
	}

	if p.Consortiums != nil {
		for _, consortium := range p.Consortiums {
			for _, org := range consortium.Organizations {
				org.completeInitialization(configDir)
			}
		}
	}

	
	if p.Orderer != nil {
		p.Orderer.completeInitialization()
	}
}

func (r *Resources) completeInitialization() {
	for {
		switch {
		case r.DefaultModPolicy == "":
			r.DefaultModPolicy = policies.ChannelApplicationAdmins
		default:
			return
		}
	}
}

func (org *Organization) completeInitialization(configDir string) {
	
	if org.MSPType == "" {
		org.MSPType = msp.ProviderTypeToString(msp.FABRIC)
	}

	if org.AdminPrincipal == "" {
		org.AdminPrincipal = AdminRoleAdminPrincipal
	}
	translatePaths(configDir, org)
}

func (oc *Orderer) completeInitialization() {
	for {
		switch {
		case oc.OrdererType == "":
			logger.Infof("Orderer.OrdererType unset, setting to %v", genesisDefaults.Orderer.OrdererType)
			oc.OrdererType = genesisDefaults.Orderer.OrdererType
		case oc.Addresses == nil:
			logger.Infof("Orderer.Addresses unset, setting to %s", genesisDefaults.Orderer.Addresses)
			oc.Addresses = genesisDefaults.Orderer.Addresses
		case oc.BatchTimeout == 0:
			logger.Infof("Orderer.BatchTimeout unset, setting to %s", genesisDefaults.Orderer.BatchTimeout)
			oc.BatchTimeout = genesisDefaults.Orderer.BatchTimeout
		case oc.BatchSize.MaxMessageCount == 0:
			logger.Infof("Orderer.BatchSize.MaxMessageCount unset, setting to %v", genesisDefaults.Orderer.BatchSize.MaxMessageCount)
			oc.BatchSize.MaxMessageCount = genesisDefaults.Orderer.BatchSize.MaxMessageCount
		case oc.BatchSize.AbsoluteMaxBytes == 0:
			logger.Infof("Orderer.BatchSize.AbsoluteMaxBytes unset, setting to %v", genesisDefaults.Orderer.BatchSize.AbsoluteMaxBytes)
			oc.BatchSize.AbsoluteMaxBytes = genesisDefaults.Orderer.BatchSize.AbsoluteMaxBytes
		case oc.BatchSize.PreferredMaxBytes == 0:
			logger.Infof("Orderer.BatchSize.PreferredMaxBytes unset, setting to %v", genesisDefaults.Orderer.BatchSize.PreferredMaxBytes)
			oc.BatchSize.PreferredMaxBytes = genesisDefaults.Orderer.BatchSize.PreferredMaxBytes
		case oc.Kafka.Brokers == nil:
			logger.Infof("Orderer.Kafka.Brokers unset, setting to %v", genesisDefaults.Orderer.Kafka.Brokers)
			oc.Kafka.Brokers = genesisDefaults.Orderer.Kafka.Brokers
		default:
			return
		}
	}
}

func translatePaths(configDir string, org *Organization) {
	cf.TranslatePathInPlace(configDir, &org.MSPDir)
}
