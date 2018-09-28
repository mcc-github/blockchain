/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft_test

import (
	"github.com/mcc-github/blockchain/common/flogging"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
	consensusmocks "github.com/mcc-github/blockchain/orderer/consensus/mocks"
	etcdraftproto "github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	"github.com/mcc-github/blockchain/protos/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Consenter", func() {
	var (
		support *consensusmocks.FakeConsenterSupport
		logger  *flogging.FabricLogger
	)

	BeforeEach(func() {
		logger = flogging.NewFabricLogger(zap.NewNop())
		support = &consensusmocks.FakeConsenterSupport{}
	})

	It("successfully constructs a chain", func() {
		certBytes := []byte("cert.orderer0.org0")
		m := &etcdraftproto.Metadata{
			Consenters: []*etcdraftproto.Consenter{
				{ServerTlsCert: certBytes},
			},
		}
		metadata := utils.MarshalOrPanic(m)
		support.SharedConfigReturns(&mockconfig.Orderer{ConsensusMetadataVal: metadata})

		consenter := &etcdraft.Consenter{
			Logger: logger,
			Cert:   certBytes,
		}

		chain, err := consenter.HandleChain(support, nil)
		Expect(chain).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		Expect(chain.Start).NotTo(Panic())
	})

	It("fails to handle chain if no matching cert found", func() {
		m := &etcdraftproto.Metadata{
			Consenters: []*etcdraftproto.Consenter{
				{ServerTlsCert: []byte("cert.orderer1.org1")},
			},
		}
		metadata := utils.MarshalOrPanic(m)
		support.SharedConfigReturns(&mockconfig.Orderer{ConsensusMetadataVal: metadata})

		consenter := &etcdraft.Consenter{
			Logger: logger,
			Cert:   []byte("cert.orderer2.org1"),
		}

		chain, err := consenter.HandleChain(support, nil)
		Expect(chain).To(BeNil())
		Expect(err).To(MatchError("failed to detect Raft ID because no matching certificate found"))
	})
})
