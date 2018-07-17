/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package runner

import (
	"fmt"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit/ginkgomon"
)


type Configtxgen struct {
	
	Path string
	
	ChannelID string
	
	Profile string
	
	AsOrg string
	
	ConfigDir string
	
	EnvConfigDir string
	
	Output string
}

func (c *Configtxgen) setupCommandEnv(cmd *exec.Cmd) {
	if c.EnvConfigDir != "" {
		configDir, err := filepath.Abs(c.EnvConfigDir)
		Expect(err).NotTo(HaveOccurred())
		cmd.Env = append(cmd.Env, fmt.Sprintf("FABRIC_CFG_PATH=%s", configDir))
	}
}


func (c *Configtxgen) OutputBlock(extraArgs ...string) *ginkgomon.Runner {
	cmd := exec.Command(
		c.Path,
		append([]string{
			"-outputBlock", c.Output,
			"-profile", c.Profile,
			"-channelID", c.ChannelID,
			"-configPath", c.ConfigDir,
		}, extraArgs...)...,
	)
	c.setupCommandEnv(cmd)
	return ginkgomon.New(ginkgomon.Config{
		Name:          "config output block",
		AnsiColorCode: "32m",
		Command:       cmd,
	})
}

func (c *Configtxgen) OutputCreateChannelTx(extraArgs ...string) *ginkgomon.Runner {
	cmd := exec.Command(
		c.Path,
		append([]string{
			"-channelID", c.ChannelID,
			"-outputCreateChannelTx", c.Output,
			"-profile", c.Profile,
			"-configPath", c.ConfigDir,
		}, extraArgs...)...,
	)
	c.setupCommandEnv(cmd)

	return ginkgomon.New(ginkgomon.Config{
		Name:          "config create channel",
		AnsiColorCode: "33m",
		Command:       cmd,
	})
}

func (c *Configtxgen) OutputAnchorPeersUpdate(extraArgs ...string) *ginkgomon.Runner {
	cmd := exec.Command(
		c.Path,
		append([]string{
			"-channelID", c.ChannelID,
			"-outputAnchorPeersUpdate", c.Output,
			"-profile", c.Profile,
			"-asOrg", c.AsOrg,
			"-configPath", c.ConfigDir,
		}, extraArgs...)...,
	)
	c.setupCommandEnv(cmd)

	return ginkgomon.New(ginkgomon.Config{
		Name:          "config update peer",
		AnsiColorCode: "34m",
		Command:       cmd,
	})
}
