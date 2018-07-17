/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package runner

import (
	"os/exec"

	"github.com/tedsuo/ifrit/ginkgomon"
)


type Idemixgen struct {
	
	Path string
	
	Output string
	
	EnrollID string
	
	OrgUnit string
	
	IsAdmin bool
	
	RevocationHandle int
}


func (c *Idemixgen) CAKeyGen(extraArgs ...string) *ginkgomon.Runner {
	return ginkgomon.New(ginkgomon.Config{
		Name:          "idemix ca-keygen",
		AnsiColorCode: "38m",
		Command: exec.Command(
			c.Path,
			append([]string{
				"ca-keygen",
				"--output", c.Output,
			}, extraArgs...)...,
		),
	})
}


func (c *Idemixgen) SignerConfig(extraArgs ...string) *ginkgomon.Runner {
	return ginkgomon.New(ginkgomon.Config{
		Name:          "idemix signerconfig",
		AnsiColorCode: "38m",
		Command: exec.Command(
			c.Path,
			append([]string{
				"signerconfig",
				"-e", c.EnrollID,
				"-u", c.OrgUnit,
				"--output", c.Output,
			}, extraArgs...)...,
		),
	})
}
