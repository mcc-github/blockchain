/*
Copyright 2017 Hitachi America, Ltd.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"context"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/admin"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/comm/testpb"
	"github.com/mcc-github/blockchain/core/peer"
	common2 "github.com/mcc-github/blockchain/internal/peer/common"
	"github.com/mcc-github/blockchain/internal/peer/mocks"
	"github.com/mcc-github/blockchain/msp"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testServiceServer struct{}

func (tss *testServiceServer) EmptyCall(context.Context, *testpb.Empty) (*testpb.Empty, error) {
	return new(testpb.Empty), nil
}

type mockEvaluator struct {
}

func (*mockEvaluator) Evaluate(signatureSet []*protoutil.SignedData) error {
	return nil
}

func TestStatusCmd(t *testing.T) {
	signer := &mocks.Signer{}
	common2.GetDefaultSignerFnc = func() (msp.SigningIdentity, error) {
		return signer, nil
	}
	viper.Set("peer.address", "localhost:7070")
	peerServer, err := peer.NewPeerServer("localhost:7070", comm.ServerConfig{})
	if err != nil {
		t.Fatalf("Failed to create peer server (%s)", err)
	} else {
		pb.RegisterAdminServer(peerServer.Server(), admin.NewAdminServer(&mockEvaluator{}))
		go peerServer.Start()
		defer peerServer.Stop()

		cmd := statusCmd()
		if err := cmd.Execute(); err != nil {
			t.Fail()
			t.Errorf("expected status command to succeed")
		}
	}
}

func TestStatus(t *testing.T) {
	defer viper.Reset()

	signer := &mocks.Signer{}
	common2.GetDefaultSignerFnc = func() (msp.SigningIdentity, error) {
		return signer, nil
	}
	var tests = []struct {
		name          string
		peerAddress   string
		listenAddress string
		timeout       time.Duration
		shouldSucceed bool
	}{
		{
			name:          "status function to success",
			peerAddress:   "localhost:7071",
			listenAddress: "localhost:7071",
			timeout:       time.Second,
			shouldSucceed: true,
		},
		{
			name:          "admin client error",
			peerAddress:   "",
			listenAddress: "localhost:7072",
			timeout:       100 * time.Millisecond,
			shouldSucceed: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Running test: %s", test.name)
			viper.Set("peer.address", test.peerAddress)
			viper.Set("peer.client.connTimeout", test.timeout)
			peerServer, err := peer.NewPeerServer(test.listenAddress, comm.ServerConfig{})
			if err != nil {
				t.Fatalf("Failed to create peer server (%s)", err)
			} else {
				pb.RegisterAdminServer(peerServer.Server(), admin.NewAdminServer(&mockEvaluator{}))
				go peerServer.Start()
				defer peerServer.Stop()
				if test.shouldSucceed {
					assert.NoError(t, status())
				} else {
					assert.Error(t, status())
				}
			}
		})
	}
}

func TestStatusWithGetStatusError(t *testing.T) {
	defer viper.Reset()

	viper.Set("peer.address", "localhost:7073")
	peerServer, err := peer.NewPeerServer(":7073", comm.ServerConfig{})
	if err != nil {
		t.Fatalf("Failed to create peer server (%s)", err)
	}
	testpb.RegisterTestServiceServer(peerServer.Server(), &testServiceServer{})
	go peerServer.Start()
	defer peerServer.Stop()
	assert.Error(t, status())
}
