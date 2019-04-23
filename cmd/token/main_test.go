/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"os/exec"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

func TestMissingArguments(t *testing.T) {
	gt := NewGomegaWithT(t)
	token, err := Build("github.com/mcc-github/blockchain/cmd/token")
	gt.Expect(err).NotTo(HaveOccurred())
	defer CleanupBuildArtifacts()

	
	cmd := exec.Command(token, "--configFile", "conf.yaml", "--MSP", "SampleOrg", "saveConfig")
	process, err := Start(cmd, nil, nil)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Eventually(process).Should(Exit(1))
	gt.Expect(process.Err).To(gbytes.Say("empty string that is mandatory"))
}
