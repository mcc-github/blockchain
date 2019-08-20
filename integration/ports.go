/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package integration

import (
	"fmt"
	"os"

	"github.com/onsi/ginkgo"
)


type TestPortRange int

const (
	basePort      = 20000
	portsPerNode  = 50
	portsPerSuite = 10 * portsPerNode
)

const (
	DiscoveryBasePort TestPortRange = basePort + portsPerSuite*iota
	E2EBasePort
	GossipBasePort
	LedgerPort
	LifecyclePort
	NWOBasePort
	PluggableBasePort
	PrivateDataBasePort
	RaftBasePort
	SBEBasePort
	IdemixBasePort
)





func (t TestPortRange) StartPortForNode() int {
	const startEphemeral, endEphemeral = 32768, 60999

	port := int(t) + portsPerNode*(ginkgo.GinkgoParallelNode()-1)
	if port >= startEphemeral-portsPerNode && port <= endEphemeral-portsPerNode {
		fmt.Fprintf(os.Stderr, "WARNING: port %d is part of the default ephemeral port range on linux", port)
	}
	return port
}
