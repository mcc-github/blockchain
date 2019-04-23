/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)


type NetworkConfig struct {
	Name                   string                          `yaml:"name"`
	Xtype                  string                          `yaml:"x-type"`
	Description            string                          `yaml:"description"`
	Version                string                          `yaml:"version"`
	Channels               map[string]ChannelNetworkConfig `yaml:"channels"`
	Organizations          map[string]OrganizationConfig   `yaml:"organizations"`
	Peers                  map[string]PeerConfig           `yaml:"peers"`
	Client                 ClientConfig                    `yaml:"client"`
	Orderers               map[string]OrdererConfig        `yaml:"orderers"`
	CertificateAuthorities map[string]CAConfig             `yaml:"certificateAuthorities"`
}


type ClientConfig struct {
	Organization    string              `yaml:"organization"`
	Logging         LoggingType         `yaml:"logging"`
	CryptoConfig    CCType              `yaml:"cryptoconfig"`
	TLS             TLSType             `yaml:"tls"`
	CredentialStore CredentialStoreType `yaml:"credentialStore"`
}


type LoggingType struct {
	Level string `yaml:"level"`
}


type CCType struct {
	Path string `yaml:"path"`
}


type TLSType struct {
	Enabled bool `yaml:"enabled"`
}


type CredentialStoreType struct {
	Path        string `yaml:"path"`
	CryptoStore struct {
		Path string `yaml:"path"`
	}
	Wallet string `yaml:"wallet"`
}


type ChannelNetworkConfig struct {
	
	Orderers []string `yaml:"orderers"`
	
	
	Peers map[string]PeerChannelConfig `yaml:"peers"`
	
	Chaincodes []string `yaml:"chaincodes"`
}


type PeerChannelConfig struct {
	EndorsingPeer  bool `yaml:"endorsingPeer"`
	ChaincodeQuery bool `yaml:"chaincodeQuery"`
	LedgerQuery    bool `yaml:"ledgerQuery"`
	EventSource    bool `yaml:"eventSource"`
}



type OrganizationConfig struct {
	MspID                  string    `yaml:"mspid"`
	Peers                  []string  `yaml:"peers"`
	CryptoPath             string    `yaml:"cryptoPath"`
	CertificateAuthorities []string  `yaml:"certificateAuthorities"`
	AdminPrivateKey        TLSConfig `yaml:"adminPrivateKey"`
	SignedCert             TLSConfig `yaml:"signedCert"`
}



type OrdererConfig struct {
	URL         string                 `yaml:"url"`
	GrpcOptions map[string]interface{} `yaml:"grpcOptions"`
	TLSCACerts  TLSConfig              `yaml:"tlsCACerts"`
}


type PeerConfig struct {
	URL         string                 `yaml:"url"`
	EventURL    string                 `yaml:"eventUrl"`
	GRPCOptions map[string]interface{} `yaml:"grpcOptions"`
	TLSCACerts  TLSConfig              `yaml:"tlsCACerts"`
}



type CAConfig struct {
	URL         string                 `yaml:"url"`
	HTTPOptions map[string]interface{} `yaml:"httpOptions"`
	TLSCACerts  MutualTLSConfig        `yaml:"tlsCACerts"`
	Registrar   EnrollCredentials      `yaml:"registrar"`
	CaName      string                 `yaml:"caName"`
}



type EnrollCredentials struct {
	EnrollID     string `yaml:"enrollId"`
	EnrollSecret string `yaml:"enrollSecret"`
}


type TLSConfig struct {
	
	
	
	
	Path string `yaml:"path"`
	
	Pem string `yaml:"pem"`
}



type MutualTLSConfig struct {
	Pem []string `yaml:"pem"`

	
	Path string `yaml:"path"`

	
	Client TLSKeyPair `yaml:"client"`
}



type TLSKeyPair struct {
	Key  TLSConfig `yaml:"key"`
	Cert TLSConfig `yaml:"cert"`
}



func GetConfig(fileName string) (*NetworkConfig, error) {
	if fileName == "" {
		return nil, errors.New("filename cannot be empty")
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "error reading connection profile")
	}

	configData := string(data)
	config := &NetworkConfig{}
	err = yaml.Unmarshal([]byte(configData), &config)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling YAML")
	}

	return config, nil
}
