/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"time"

	"github.com/pkg/errors"
)


type ConnectionConfig struct {
	Address            string
	ConnectionTimeout  time.Duration
	TLSEnabled         bool
	TLSRootCertFile    string
	ServerNameOverride string
}

type MSPInfo struct {
	MSPConfigPath string `yaml:"mspConfigPath"`
	MSPID         string `yaml:"mspId"`
	MSPType       string `yaml:"mspType"`
}



type ClientConfig struct {
	ChannelID     string
	MSPInfo       MSPInfo
	Orderer       ConnectionConfig
	CommitterPeer ConnectionConfig
	ProverPeer    ConnectionConfig
}

func ValidateClientConfig(config ClientConfig) error {
	if config.ChannelID == "" {
		return errors.New("missing channel id")
	}

	if config.MSPInfo.MSPConfigPath == "" {
		return errors.New("missing MSP config path")
	}
	if config.MSPInfo.MSPID == "" {
		return errors.New("missing MSP ID")
	}

	if config.Orderer.Address == "" {
		return errors.New("missing orderer address")
	}
	if config.Orderer.TLSEnabled && config.Orderer.TLSRootCertFile == "" {
		return errors.New("missing orderer TLSRootCertFile")
	}

	if config.CommitterPeer.Address == "" {
		return errors.New("missing committer peer address")
	}
	if config.CommitterPeer.TLSEnabled && config.CommitterPeer.TLSRootCertFile == "" {
		return errors.New("missing committer peer TLSRootCertFile")
	}

	if config.ProverPeer.Address == "" {
		return errors.New("missing prover peer address")
	}
	if config.ProverPeer.TLSEnabled && config.ProverPeer.TLSRootCertFile == "" {
		return errors.New("missing prover peer TLSRootCertFile")
	}

	return nil
}
