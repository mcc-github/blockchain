/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
	lc "github.com/mcc-github/blockchain/protos/peer/lifecycle"

	"github.com/golang/protobuf/proto"
)

var _ = Describe("ProtobufImpl", func() {
	var (
		pi        *lifecycle.ProtobufImpl
		sampleMsg *lc.InstallChaincodeArgs
	)

	BeforeEach(func() {
		pi = &lifecycle.ProtobufImpl{}
		sampleMsg = &lc.InstallChaincodeArgs{
			Name:                    "name",
			Version:                 "version",
			ChaincodeInstallPackage: []byte("install-package"),
		}
	})

	Describe("Marshal", func() {
		It("passes through to the proto implementation", func() {
			res, err := pi.Marshal(sampleMsg)
			Expect(err).NotTo(HaveOccurred())

			msg := &lc.InstallChaincodeArgs{}
			err = proto.Unmarshal(res, msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(proto.Equal(msg, sampleMsg)).To(BeTrue())
		})
	})

	Describe("Unmarshal", func() {
		It("passes through to the proto implementation", func() {
			res, err := proto.Marshal(sampleMsg)
			Expect(err).NotTo(HaveOccurred())

			msg := &lc.InstallChaincodeArgs{}
			err = pi.Unmarshal(res, msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(proto.Equal(msg, sampleMsg)).To(BeTrue())
		})
	})
})
