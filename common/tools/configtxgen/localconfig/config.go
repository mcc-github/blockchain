/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package localconfig

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/viperutil"
	cf "github.com/mcc-github/blockchain/core/config"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	"github.com/spf13/viper"
)

const (
	pkgLogID = "common.tools.configtxgen.localconfig"

	
	Prefix string = "CONFIGTX"
)

var logger = flogging.MustGetLogger(pkgLogID)
var configName = strings.ToLower(Prefix)

func init() {
	flogging.InitFromSpec(pkgLogID + "=error")
}

const (
	
	
	TestChainID = "testchainid"

	
	
	SampleInsecureSoloProfile = "SampleInsecureSolo"
	
	
	SampleDevModeSoloProfile = "SampleDevModeSolo"
	
	
	SampleSingleMSPSoloProfile = "SampleSingleMSPSolo"

	
	
	SampleInsecureKafkaProfile = "SampleInsecureKafka"
	
	
	SampleDevModeKafkaProfile = "SampleDevModeKafka"
	
	
	SampleSingleMSPKafkaProfile = "SampleSingleMSPKafka"

	
	
	SampleDevModeEtcdRaftProfile = "SampleDevModeEtcdRaft"

	
	
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
	EtcdRaft      *etcdraft.Metadata `yaml:"EtcdRaft"`
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
		EtcdRaft: &etcdraft.Metadata{
			Options: &etcdraft.Options{
				TickInterval:    100,
				ElectionTick:    10,
				HeartbeatTick:   1,
				MaxInflightMsgs: 256,
				MaxSizePerMsg:   1048576,
			},
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
		t.Orderer.completeInitialization(configDir)
	}
}

func (p *Profile) completeInitialization(configDir string) {
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
		for _, org := range p.Orderer.Organizations {
			org.completeInitialization(configDir)
		}
		
		p.Orderer.completeInitialization(configDir)
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

func (ord *Orderer) completeInitialization(configDir string) {
loop:
	for {
		switch {
		case ord.OrdererType == "":
			logger.Infof("Orderer.OrdererType unset, setting to %v", genesisDefaults.Orderer.OrdererType)
			ord.OrdererType = genesisDefaults.Orderer.OrdererType
		case ord.Addresses == nil:
			logger.Infof("Orderer.Addresses unset, setting to %s", genesisDefaults.Orderer.Addresses)
			ord.Addresses = genesisDefaults.Orderer.Addresses
		case ord.BatchTimeout == 0:
			logger.Infof("Orderer.BatchTimeout unset, setting to %s", genesisDefaults.Orderer.BatchTimeout)
			ord.BatchTimeout = genesisDefaults.Orderer.BatchTimeout
		case ord.BatchSize.MaxMessageCount == 0:
			logger.Infof("Orderer.BatchSize.MaxMessageCount unset, setting to %v", genesisDefaults.Orderer.BatchSize.MaxMessageCount)
			ord.BatchSize.MaxMessageCount = genesisDefaults.Orderer.BatchSize.MaxMessageCount
		case ord.BatchSize.AbsoluteMaxBytes == 0:
			logger.Infof("Orderer.BatchSize.AbsoluteMaxBytes unset, setting to %v", genesisDefaults.Orderer.BatchSize.AbsoluteMaxBytes)
			ord.BatchSize.AbsoluteMaxBytes = genesisDefaults.Orderer.BatchSize.AbsoluteMaxBytes
		case ord.BatchSize.PreferredMaxBytes == 0:
			logger.Infof("Orderer.BatchSize.PreferredMaxBytes unset, setting to %v", genesisDefaults.Orderer.BatchSize.PreferredMaxBytes)
			ord.BatchSize.PreferredMaxBytes = genesisDefaults.Orderer.BatchSize.PreferredMaxBytes
		default:
			break loop
		}
	}

	logger.Infof("orderer type: %s", ord.OrdererType)
	
	
	switch ord.OrdererType {
	case "solo":
		
	case "kafka":
		if ord.Kafka.Brokers == nil {
			logger.Infof("Orderer.Kafka unset, setting to %v", genesisDefaults.Orderer.Kafka.Brokers)
			ord.Kafka.Brokers = genesisDefaults.Orderer.Kafka.Brokers
		}
	case etcdraft.TypeKey:
		if ord.EtcdRaft == nil {
			logger.Panicf("%s raft configuration missing", etcdraft.TypeKey)
		}
		if ord.EtcdRaft.Options == nil {
			logger.Infof("Orderer.EtcdRaft.Options unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options)
			ord.EtcdRaft.Options = genesisDefaults.Orderer.EtcdRaft.Options
		}
	second_loop:
		for {
			switch {
			case ord.EtcdRaft.Options.TickInterval == 0:
				logger.Infof("Orderer.EtcdRaft.Options.TickInterval unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.TickInterval)
				ord.EtcdRaft.Options.TickInterval = genesisDefaults.Orderer.EtcdRaft.Options.TickInterval

			case ord.EtcdRaft.Options.ElectionTick == 0:
				logger.Infof("Orderer.EtcdRaft.Options.ElectionTick unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.ElectionTick)
				ord.EtcdRaft.Options.ElectionTick = genesisDefaults.Orderer.EtcdRaft.Options.ElectionTick

			case ord.EtcdRaft.Options.HeartbeatTick == 0:
				logger.Infof("Orderer.EtcdRaft.Options.HeartbeatTick unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.HeartbeatTick)
				ord.EtcdRaft.Options.HeartbeatTick = genesisDefaults.Orderer.EtcdRaft.Options.HeartbeatTick

			case ord.EtcdRaft.Options.MaxInflightMsgs == 0:
				logger.Infof("Orderer.EtcdRaft.Options.MaxInflightMsgs unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.MaxInflightMsgs)
				ord.EtcdRaft.Options.MaxInflightMsgs = genesisDefaults.Orderer.EtcdRaft.Options.MaxInflightMsgs

			case ord.EtcdRaft.Options.MaxSizePerMsg == 0:
				logger.Infof("Orderer.EtcdRaft.Options.MaxSizePerMsg unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.MaxSizePerMsg)
				ord.EtcdRaft.Options.MaxSizePerMsg = genesisDefaults.Orderer.EtcdRaft.Options.MaxSizePerMsg

			case len(ord.EtcdRaft.Consenters) == 0:
				logger.Panicf("%s configuration did not specify any consenter", etcdraft.TypeKey)

			default:
				break second_loop
			}
		}

		
		if ord.EtcdRaft.Options.ElectionTick <= ord.EtcdRaft.Options.HeartbeatTick {
			logger.Panicf("election tick must be greater than heartbeat tick")
		}

		for _, c := range ord.EtcdRaft.GetConsenters() {
			if c.Host == "" {
				logger.Panicf("consenter info in %s configuration did not specify host", etcdraft.TypeKey)
			}
			if c.Port == 0 {
				logger.Panicf("consenter info in %s configuration did not specify port", etcdraft.TypeKey)
			}
			if c.ClientTlsCert == nil {
				logger.Panicf("consenter info in %s configuration did not specify client TLS cert", etcdraft.TypeKey)
			}
			if c.ServerTlsCert == nil {
				logger.Panicf("consenter info in %s configuration did not specify server TLS cert", etcdraft.TypeKey)
			}
			clientCertPath := string(c.GetClientTlsCert())
			cf.TranslatePathInPlace(configDir, &clientCertPath)
			c.ClientTlsCert = []byte(clientCertPath)
			serverCertPath := string(c.GetServerTlsCert())
			cf.TranslatePathInPlace(configDir, &serverCertPath)
			c.ServerTlsCert = []byte(serverCertPath)
		}
	default:
		logger.Panicf("unknown orderer type: %s", ord.OrdererType)
	}
}

func translatePaths(configDir string, org *Organization) {
	cf.TranslatePathInPlace(configDir, &org.MSPDir)
}
