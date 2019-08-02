/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc_test

import (
	"os"
	"testing"

	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/scc/mock"
	"github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


type chaincodeStreamHandler interface {
	scc.ChaincodeStreamHandler
}

func init() {
	viper.Set("peer.fileSystemPath", os.TempDir())
}

func newTestProvider() *scc.Provider {
	p := &scc.Provider{
		Peer: &peer.Peer{},
		Whitelist: map[string]bool{
			"invokableExternalButNotCC2CC": true,
			"invokableCC2CCButNotExternal": true,
			"disabled":                     true,
		},
	}
	for _, cc := range []scc.SelfDescribingSysCC{
		&scc.SysCCWrapper{
			SCC: &scc.SystemChaincode{
				Name:              "invokableExternalButNotCC2CC",
				InvokableExternal: true,
				InvokableCC2CC:    false,
				Enabled:           true,
			},
		},
		&scc.SysCCWrapper{
			SCC: &scc.SystemChaincode{
				Name:              "invokableCC2CCButNotExternal",
				InvokableExternal: false,
				InvokableCC2CC:    true,
				Enabled:           true,
			},
		},
		&scc.SysCCWrapper{
			SCC: &scc.SystemChaincode{
				Name:    "disabled",
				Enabled: false,
			},
		},
	} {
		p.RegisterSysCC(cc)
	}
	return p
}

func TestDeploy(t *testing.T) {
	gt := gomega.NewGomegaWithT(t)

	p := newTestProvider()
	csh := &mock.ChaincodeStreamHandler{}
	doneC := make(chan struct{})
	close(doneC)
	csh.LaunchInProcReturns(doneC)
	p.DeploySysCCs(csh)
	gt.Expect(csh.LaunchInProcCallCount()).To(gomega.Equal(2))
	gt.Expect(csh.LaunchInProcArgsForCall(0)).To(gomega.Equal(ccintf.CCID("invokableExternalButNotCC2CC:latest")))
	gt.Expect(csh.LaunchInProcArgsForCall(1)).To(gomega.Equal(ccintf.CCID("invokableCC2CCButNotExternal:latest")))
	gt.Eventually(csh.HandleChaincodeStreamCallCount).Should(gomega.Equal(2))
}

func TestIsSysCC(t *testing.T) {
	assert.True(t, (newTestProvider()).IsSysCC("invokableExternalButNotCC2CC"))
	assert.False(t, (newTestProvider()).IsSysCC("noSCC"))
	assert.True(t, (newTestProvider()).IsSysCC("invokableCC2CCButNotExternal"))
	assert.True(t, (newTestProvider()).IsSysCC("disabled"))
}

func TestSccProviderImpl_GetQueryExecutorForLedger(t *testing.T) {
	p := &scc.Provider{
		Peer: &peer.Peer{},
	}
	qe, err := p.GetQueryExecutorForLedger("")
	assert.Nil(t, qe)
	assert.Error(t, err)
}

func TestCreatePluginSysCCs(t *testing.T) {
	assert.NotPanics(t, func() { scc.CreatePluginSysCCs(nil) }, "expected successful init")
}

func TestRegisterSysCC(t *testing.T) {
	p := &scc.Provider{
		Whitelist: map[string]bool{
			"invokableExternalButNotCC2CC": true,
			"invokableCC2CCButNotExternal": true,
		},
	}
	err := p.RegisterSysCC(&scc.SysCCWrapper{
		SCC: &scc.SystemChaincode{
			Name:      "invokableExternalButNotCC2CC",
			Path:      "path",
			Enabled:   true,
			Chaincode: nil,
		},
	})
	assert.NoError(t, err)
	err = p.RegisterSysCC(&scc.SysCCWrapper{
		SCC: &scc.SystemChaincode{
			Name:      "invokableExternalButNotCC2CC",
			Path:      "path",
			Enabled:   true,
			Chaincode: nil,
		},
	})
	assert.EqualError(t, err, "chaincode with name 'invokableExternalButNotCC2CC' already registered")
}
