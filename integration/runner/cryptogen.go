/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package runner

import (
	"os/exec"

	"github.com/tedsuo/ifrit/ginkgomon"
)


type Cryptogen struct {
	
	Path string
	
	Config string
	
	Output string
}


func (c *Cryptogen) Generate(extraArgs ...string) *ginkgomon.Runner {
	return ginkgomon.New(ginkgomon.Config{
		Name:          "cryptogen generate",
		AnsiColorCode: "31m",
		Command: exec.Command(
			c.Path,
			append([]string{
				"generate",
				"--config", c.Config,
				"--output", c.Output,
			}, extraArgs...)...,
		),
	})
}
