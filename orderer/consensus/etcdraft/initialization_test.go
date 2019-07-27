/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft_test

import (
	"testing"

	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/orderer/common/multichannel"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewEtcdRaftConsenter(t *testing.T) {
	srv, err := comm.NewGRPCServer("127.0.0.1:0", comm.ServerConfig{})
	assert.NoError(t, err)
	defer srv.Stop()
	dialer := &cluster.PredicateDialer{}
	consenter := etcdraft.New(dialer,
		&localconfig.TopLevel{},
		comm.ServerConfig{
			SecOpts: comm.SecureOptions{
				Certificate: []byte{1, 2, 3},
			},
		}, srv, &multichannel.Registrar{},
		&mocks.InactiveChainRegistry{}, &disabled.Provider{})

	
	assert.Equal(t, []byte{1, 2, 3}, consenter.Cert)
	
	assert.NotNil(t, consenter.Communication)
	assert.NotNil(t, consenter.Chains)
	assert.NotNil(t, consenter.ChainSelector)
	assert.NotNil(t, consenter.Dispatcher)
	assert.NotNil(t, consenter.Logger)
}
