/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo



type Config struct {
	Organizations []*Organization `yaml:"organizations,omitempty"`
	Consortiums   []*Consortium   `yaml:"consortiums,omitempty"`
	SystemChannel *SystemChannel  `yaml:"system_channel,omitempty"`
	Channels      []*Channel      `yaml:"channels,omitempty"`
	Consensus     *Consensus      `yaml:"consensus,omitempty"`
	Orderers      []*Orderer      `yaml:"orderers,omitempty"`
	Peers         []*Peer         `yaml:"peers,omitempty"`
	Profiles      []*Profile      `yaml:"profiles,omitempty"`
}

func (c *Config) RemovePeer(orgName, peerName string) {
	peers := []*Peer{}
	for _, p := range c.Peers {
		if p.Organization != orgName || p.Name != peerName {
			peers = append(peers, p)
		}
	}
	c.Peers = peers
}
