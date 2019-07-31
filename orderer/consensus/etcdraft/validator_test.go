

package etcdraft_test

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/orderer/consensus/etcdraft"
	etcdraftproto "github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata Validation", func() {
	var (
		chain *etcdraft.Chain
	)

	BeforeEach(func() {
		chain = &etcdraft.Chain{}
	})

	When("determining parameter well-formedness", func() {
		It("succeeds when new consensus metadata is nil", func() {
			Expect(chain.ValidateConsensusMetadata(nil, nil, false)).To(Succeed())
		})

		It("fails when new consensus metadata is not nil while old consensus metadata is nil", func() {
			Expect(func() {
				chain.ValidateConsensusMetadata(nil, []byte("test"), false)
			}).To(Panic())
		})

		It("fails when old consensus metadata is not well-formed", func() {
			Expect(func() {
				chain.ValidateConsensusMetadata([]byte("test"), []byte("test"), false)
			}).To(Panic())
		})

		It("fails when new consensus metadata is not well-formed", func() {
			oldBytes, _ := proto.Marshal(&etcdraftproto.ConfigMetadata{})
			Expect(chain.ValidateConsensusMetadata(oldBytes, []byte("test"), false)).NotTo(Succeed())
		})
	})

	Context("valid old consensus metadata", func() {
		var (
			oldBytes    []byte
			oldMetadata *etcdraftproto.ConfigMetadata
			newMetadata *etcdraftproto.ConfigMetadata
			tlsCA       tlsgen.CA
			newChannel  bool
		)

		BeforeEach(func() {
			tlsCA, _ = tlsgen.NewCA()
			oldMetadata = &etcdraftproto.ConfigMetadata{
				Options: &etcdraftproto.Options{
					TickInterval:         "500ms",
					ElectionTick:         10,
					HeartbeatTick:        1,
					MaxInflightBlocks:    5,
					SnapshotIntervalSize: 20 * 1024 * 1024, 
				},
				Consenters: []*etcdraftproto.Consenter{
					{
						Host:          "host1",
						Port:          10001,
						ClientTlsCert: clientTLSCert(tlsCA),
						ServerTlsCert: serverTLSCert(tlsCA),
					},
					{
						Host:          "host2",
						Port:          10002,
						ClientTlsCert: clientTLSCert(tlsCA),
						ServerTlsCert: serverTLSCert(tlsCA),
					},
					{
						Host:          "host3",
						Port:          10003,
						ClientTlsCert: clientTLSCert(tlsCA),
						ServerTlsCert: serverTLSCert(tlsCA),
					},
				},
			}
			newMetadata = oldMetadata
			oldBytes, _ = proto.Marshal(oldMetadata)
			newChannel = false
		})

		It("fails when new consensus metadata has invalid options", func() {
			
			newMetadata.Options.TickInterval = ""
			newBytes, _ := proto.Marshal(newMetadata)
			Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).NotTo(Succeed())
		})

		Context("new channel creation", func() {

			BeforeEach(func() {
				newChannel = true
			})

			It("fails when the new consenters are an empty set", func() {
				newMetadata.Consenters = []*etcdraftproto.Consenter{}
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).NotTo(Succeed())
			})

			It("succeeds when the new consenters are the same as the existing consenters", func() {
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).To(Succeed())
			})

			It("succeeds when the new consenters are a subset of the existing consenters", func() {
				newMetadata.Consenters = newMetadata.Consenters[:2]
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).To(Succeed())
			})

			It("fails when the new consenters are not a subset of the existing consenters", func() {
				newMetadata.Consenters[2].ClientTlsCert = clientTLSCert(tlsCA)
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).NotTo(Succeed())
			})

		})

		Context("config update on a channel", func() {

			BeforeEach(func() {
				newChannel = false
			})

			It("fails when the new consenters are an empty set", func() {
				
				newMetadata.Consenters = []*etcdraftproto.Consenter{}
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).NotTo(Succeed())
			})

			It("succeeds when the new consenters are the same as the existing consenters", func() {
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).To(Succeed())
			})

			It("succeeds on addition of a single consenter", func() {
				newMetadata.Consenters = append(newMetadata.Consenters, &etcdraftproto.Consenter{
					Host:          "host4",
					Port:          10004,
					ClientTlsCert: clientTLSCert(tlsCA),
					ServerTlsCert: serverTLSCert(tlsCA),
				})
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).To(Succeed())
			})

			It("fails on addition of more than one consenter", func() {
				newMetadata.Consenters = append(newMetadata.Consenters,
					&etcdraftproto.Consenter{
						Host:          "host4",
						Port:          10004,
						ClientTlsCert: clientTLSCert(tlsCA),
						ServerTlsCert: serverTLSCert(tlsCA),
					},
					&etcdraftproto.Consenter{
						Host:          "host5",
						Port:          10005,
						ClientTlsCert: clientTLSCert(tlsCA),
						ServerTlsCert: serverTLSCert(tlsCA),
					},
				)
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).NotTo(Succeed())
			})

			It("succeeds on removal of a single consenter", func() {
				newMetadata.Consenters = newMetadata.Consenters[:2]
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).To(Succeed())
			})

			It("fails on removal of more than one consenter", func() {
				newMetadata.Consenters = newMetadata.Consenters[:1]
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).NotTo(Succeed())
			})

			It("succeeds on rotating certs in case of both addition and removal of a node each to reuse the raft NodeId", func() {
				newMetadata.Consenters = append(newMetadata.Consenters[:2], &etcdraftproto.Consenter{
					Host:          "host4",
					Port:          10004,
					ClientTlsCert: clientTLSCert(tlsCA),
					ServerTlsCert: serverTLSCert(tlsCA),
				})
				newBytes, _ := proto.Marshal(newMetadata)
				Expect(chain.ValidateConsensusMetadata(oldBytes, newBytes, newChannel)).To(Succeed())
			})

		})
	})
})
