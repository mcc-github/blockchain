/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/integration/helpers"
	"github.com/mcc-github/blockchain/integration/nwo/commands"
	"github.com/mcc-github/blockchain/integration/nwo/blockchainconfig"
	"github.com/mcc-github/blockchain/integration/runner"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
	"github.com/tedsuo/ifrit/grouper"
	"gopkg.in/yaml.v2"
)



type Organization struct {
	MSPID         string `yaml:"msp_id,omitempty"`
	Name          string `yaml:"name,omitempty"`
	Domain        string `yaml:"domain,omitempty"`
	EnableNodeOUs bool   `yaml:"enable_node_organizational_units"`
	Users         int    `yaml:"users,omitempty"`
	CA            *CA    `yaml:"ca,omitempty"`
}

type CA struct {
	Hostname string `yaml:"hostname,omitempty"`
}



type Consortium struct {
	Name          string   `yaml:"name,omitempty"`
	Organizations []string `yaml:"organizations,omitempty"`
}



type Consensus struct {
	Type       string `yaml:"type,omitempty"`
	Brokers    int    `yaml:"brokers,omitempty"`
	ZooKeepers int    `yaml:"zookeepers,omitempty"`
}



type SystemChannel struct {
	Name    string `yaml:"name,omitempty"`
	Profile string `yaml:"profile,omitempty"`
}


type Channel struct {
	Name    string `yaml:"name,omitempty"`
	Profile string `yaml:"profile,omitempty"`
}


type Orderer struct {
	Name         string `yaml:"name,omitempty"`
	Organization string `yaml:"organization,omitempty"`
}


func (o Orderer) ID() string {
	return fmt.Sprintf("%s.%s", o.Organization, o.Name)
}



type Peer struct {
	Name         string         `yaml:"name,omitempty"`
	Organization string         `yaml:"organization,omitempty"`
	Channels     []*PeerChannel `yaml:"channels,omitempty"`
}



type PeerChannel struct {
	Name   string `yaml:"name,omitempty"`
	Anchor bool   `yaml:"anchor"`
}


func (p *Peer) ID() string {
	return fmt.Sprintf("%s.%s", p.Organization, p.Name)
}


func (p *Peer) Anchor() bool {
	for _, c := range p.Channels {
		if c.Anchor {
			return true
		}
	}
	return false
}


type Profile struct {
	Name          string   `yaml:"name,omitempty"`
	Orderers      []string `yaml:"orderers,omitempty"`
	Consortium    string   `yaml:"consortium,omitempty"`
	Organizations []string `yaml:"organizations,omitempty"`
}


type Network struct {
	RootDir           string
	StartPort         uint16
	Components        *Components
	DockerClient      *docker.Client
	NetworkID         string
	EventuallyTimeout time.Duration
	MetricsProvider   string
	StatsdEndpoint    string

	PortsByBrokerID  map[string]Ports
	PortsByOrdererID map[string]Ports
	PortsByPeerID    map[string]Ports
	Organizations    []*Organization
	SystemChannel    *SystemChannel
	Channels         []*Channel
	Consensus        *Consensus
	Orderers         []*Orderer
	Peers            []*Peer
	Profiles         []*Profile
	Consortiums      []*Consortium
	Templates        *Templates

	colorIndex uint
}




func New(c *Config, rootDir string, client *docker.Client, startPort int, components *Components) *Network {
	network := &Network{
		StartPort:    uint16(startPort),
		RootDir:      rootDir,
		Components:   components,
		DockerClient: client,

		NetworkID:         helpers.UniqueName(),
		EventuallyTimeout: time.Minute,
		MetricsProvider:   "prometheus",
		PortsByBrokerID:   map[string]Ports{},
		PortsByOrdererID:  map[string]Ports{},
		PortsByPeerID:     map[string]Ports{},

		Organizations: c.Organizations,
		Consensus:     c.Consensus,
		Orderers:      c.Orderers,
		Peers:         c.Peers,
		SystemChannel: c.SystemChannel,
		Channels:      c.Channels,
		Profiles:      c.Profiles,
		Consortiums:   c.Consortiums,
		Templates:     c.Templates,
	}

	if network.Templates == nil {
		network.Templates = &Templates{}
	}

	for i := 0; i < network.Consensus.Brokers; i++ {
		ports := Ports{}
		for _, portName := range BrokerPortNames() {
			ports[portName] = network.ReservePort()
		}
		network.PortsByBrokerID[strconv.Itoa(i)] = ports
	}

	for _, o := range c.Orderers {
		ports := Ports{}
		for _, portName := range OrdererPortNames() {
			ports[portName] = network.ReservePort()
		}
		network.PortsByOrdererID[o.ID()] = ports
	}

	for _, p := range c.Peers {
		ports := Ports{}
		for _, portName := range PeerPortNames() {
			ports[portName] = network.ReservePort()
		}
		network.PortsByPeerID[p.ID()] = ports
	}
	return network
}



func (n *Network) ConfigTxConfigPath() string {
	return filepath.Join(n.RootDir, "configtx.yaml")
}



func (n *Network) CryptoPath() string {
	return filepath.Join(n.RootDir, "crypto")
}



func (n *Network) CryptoConfigPath() string {
	return filepath.Join(n.RootDir, "crypto-config.yaml")
}



func (n *Network) OutputBlockPath(channelName string) string {
	return filepath.Join(n.RootDir, fmt.Sprintf("%s_block.pb", channelName))
}



func (n *Network) CreateChannelTxPath(channelName string) string {
	return filepath.Join(n.RootDir, fmt.Sprintf("%s_tx.pb", channelName))
}



func (n *Network) OrdererDir(o *Orderer) string {
	return filepath.Join(n.RootDir, "orderers", o.ID())
}



func (n *Network) OrdererConfigPath(o *Orderer) string {
	return filepath.Join(n.OrdererDir(o), "orderer.yaml")
}



func (n *Network) ReadOrdererConfig(o *Orderer) *blockchainconfig.Orderer {
	var orderer blockchainconfig.Orderer
	ordererBytes, err := ioutil.ReadFile(n.OrdererConfigPath(o))
	Expect(err).NotTo(HaveOccurred())

	err = yaml.Unmarshal(ordererBytes, &orderer)
	Expect(err).NotTo(HaveOccurred())

	return &orderer
}



func (n *Network) WriteOrdererConfig(o *Orderer, config *blockchainconfig.Orderer) {
	ordererBytes, err := yaml.Marshal(config)
	Expect(err).NotTo(HaveOccurred())

	err = ioutil.WriteFile(n.OrdererConfigPath(o), ordererBytes, 0644)
	Expect(err).NotTo(HaveOccurred())
}



func (n *Network) PeerDir(p *Peer) string {
	return filepath.Join(n.RootDir, "peers", p.ID())
}



func (n *Network) PeerConfigPath(p *Peer) string {
	return filepath.Join(n.PeerDir(p), "core.yaml")
}



func (n *Network) ReadPeerConfig(p *Peer) *blockchainconfig.Core {
	var core blockchainconfig.Core
	coreBytes, err := ioutil.ReadFile(n.PeerConfigPath(p))
	Expect(err).NotTo(HaveOccurred())

	err = yaml.Unmarshal(coreBytes, &core)
	Expect(err).NotTo(HaveOccurred())

	return &core
}



func (n *Network) WritePeerConfig(p *Peer, config *blockchainconfig.Core) {
	coreBytes, err := yaml.Marshal(config)
	Expect(err).NotTo(HaveOccurred())

	err = ioutil.WriteFile(n.PeerConfigPath(p), coreBytes, 0644)
	Expect(err).NotTo(HaveOccurred())
}



func (n *Network) peerUserCryptoDir(p *Peer, user, cryptoMaterialType string) string {
	org := n.Organization(p.Organization)
	Expect(org).NotTo(BeNil())

	return n.userCryptoDir(org, "peerOrganizations", user, cryptoMaterialType)
}



func (n *Network) ordererUserCryptoDir(o *Orderer, user, cryptoMaterialType string) string {
	org := n.Organization(o.Organization)
	Expect(org).NotTo(BeNil())

	return n.userCryptoDir(org, "ordererOrganizations", user, cryptoMaterialType)
}



func (n *Network) userCryptoDir(org *Organization, nodeOrganizationType, user, cryptoMaterialType string) string {
	return filepath.Join(
		n.RootDir,
		"crypto",
		nodeOrganizationType,
		org.Domain,
		"users",
		fmt.Sprintf("%s@%s", user, org.Domain),
		cryptoMaterialType,
	)
}



func (n *Network) PeerUserMSPDir(p *Peer, user string) string {
	return n.peerUserCryptoDir(p, user, "msp")
}



func (n *Network) OrdererUserMSPDir(o *Orderer, user string) string {
	return n.ordererUserCryptoDir(o, user, "msp")
}



func (n *Network) PeerUserTLSDir(p *Peer, user string) string {
	return n.peerUserCryptoDir(p, user, "tls")
}



func (n *Network) PeerUserCert(p *Peer, user string) string {
	org := n.Organization(p.Organization)
	Expect(org).NotTo(BeNil())

	return filepath.Join(
		n.PeerUserMSPDir(p, user),
		"signcerts",
		fmt.Sprintf("%s@%s-cert.pem", user, org.Domain),
	)
}



func (n *Network) PeerUserKey(p *Peer, user string) string {
	org := n.Organization(p.Organization)
	Expect(org).NotTo(BeNil())

	keystore := filepath.Join(
		n.PeerUserMSPDir(p, user),
		"keystore",
	)

	
	keys, err := ioutil.ReadDir(keystore)
	Expect(err).NotTo(HaveOccurred())
	Expect(keys).To(HaveLen(1))

	return filepath.Join(keystore, keys[0].Name())
}


func (n *Network) peerLocalCryptoDir(p *Peer, cryptoType string) string {
	org := n.Organization(p.Organization)
	Expect(org).NotTo(BeNil())

	return filepath.Join(
		n.RootDir,
		"crypto",
		"peerOrganizations",
		org.Domain,
		"peers",
		fmt.Sprintf("%s.%s", p.Name, org.Domain),
		cryptoType,
	)
}


func (n *Network) PeerLocalMSPDir(p *Peer) string {
	return n.peerLocalCryptoDir(p, "msp")
}


func (n *Network) PeerLocalTLSDir(p *Peer) string {
	return n.peerLocalCryptoDir(p, "tls")
}


func (n *Network) PeerCert(p *Peer) string {
	org := n.Organization(p.Organization)
	Expect(org).NotTo(BeNil())

	return filepath.Join(
		n.PeerLocalMSPDir(p),
		"signcerts",
		fmt.Sprintf("%s.%s-cert.pem", p.Name, org.Domain),
	)
}


func (n *Network) PeerOrgMSPDir(org *Organization) string {
	return filepath.Join(
		n.RootDir,
		"crypto",
		"peerOrganizations",
		org.Domain,
		"msp",
	)
}



func (n *Network) OrdererOrgMSPDir(o *Organization) string {
	return filepath.Join(
		n.RootDir,
		"crypto",
		"ordererOrganizations",
		o.Domain,
		"msp",
	)
}



func (n *Network) OrdererLocalCryptoDir(o *Orderer, cryptoType string) string {
	org := n.Organization(o.Organization)
	Expect(org).NotTo(BeNil())

	return filepath.Join(
		n.RootDir,
		"crypto",
		"ordererOrganizations",
		org.Domain,
		"orderers",
		fmt.Sprintf("%s.%s", o.Name, org.Domain),
		cryptoType,
	)
}



func (n *Network) OrdererLocalMSPDir(o *Orderer) string {
	return n.OrdererLocalCryptoDir(o, "msp")
}



func (n *Network) OrdererLocalTLSDir(o *Orderer) string {
	return n.OrdererLocalCryptoDir(o, "tls")
}



func (n *Network) ProfileForChannel(channelName string) string {
	for _, ch := range n.Channels {
		if ch.Name == channelName {
			return ch.Profile
		}
	}
	return ""
}



func (n *Network) CACertsBundlePath() string {
	return filepath.Join(
		n.RootDir,
		"crypto",
		"ca-certs.pem",
	)
}

















func (n *Network) GenerateConfigTree() {
	n.GenerateCryptoConfig()
	n.GenerateConfigTxConfig()
	for _, o := range n.Orderers {
		n.GenerateOrdererConfig(o)
	}
	for _, p := range n.Peers {
		n.GenerateCoreConfig(p)
	}
}
















func (n *Network) Bootstrap() {
	_, err := n.DockerClient.CreateNetwork(
		docker.CreateNetworkOptions{
			Name:   n.NetworkID,
			Driver: "bridge",
		},
	)
	Expect(err).NotTo(HaveOccurred())

	sess, err := n.Cryptogen(commands.Generate{
		Config: n.CryptoConfigPath(),
		Output: n.CryptoPath(),
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

	sess, err = n.ConfigTxGen(commands.OutputBlock{
		ChannelID:   n.SystemChannel.Name,
		Profile:     n.SystemChannel.Profile,
		ConfigPath:  n.RootDir,
		OutputBlock: n.OutputBlockPath(n.SystemChannel.Name),
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

	for _, c := range n.Channels {
		sess, err := n.ConfigTxGen(commands.CreateChannelTx{
			ChannelID:             c.Name,
			Profile:               c.Profile,
			ConfigPath:            n.RootDir,
			OutputCreateChannelTx: n.CreateChannelTxPath(c.Name),
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	}

	n.concatenateTLSCACertificates()
}



func (n *Network) concatenateTLSCACertificates() {
	bundle := &bytes.Buffer{}
	for _, tlsCertPath := range n.listTLSCACertificates() {
		certBytes, err := ioutil.ReadFile(tlsCertPath)
		Expect(err).NotTo(HaveOccurred())
		bundle.Write(certBytes)
	}
	err := ioutil.WriteFile(n.CACertsBundlePath(), bundle.Bytes(), 0660)
	Expect(err).NotTo(HaveOccurred())
}



func (n *Network) listTLSCACertificates() []string {
	fileName2Path := make(map[string]string)
	filepath.Walk(filepath.Join(n.RootDir, "crypto"), func(path string, info os.FileInfo, err error) error {
		
		if strings.HasPrefix(info.Name(), "tlsca") && strings.Contains(info.Name(), "-cert.pem") {
			fileName2Path[info.Name()] = path
		}
		return nil
	})

	var tlsCACertificates []string
	for _, path := range fileName2Path {
		tlsCACertificates = append(tlsCACertificates, path)
	}
	return tlsCACertificates
}



func (n *Network) Cleanup() {
	nw, err := n.DockerClient.NetworkInfo(n.NetworkID)
	Expect(err).NotTo(HaveOccurred())

	err = n.DockerClient.RemoveNetwork(nw.ID)
	Expect(err).NotTo(HaveOccurred())

	containers, err := n.DockerClient.ListContainers(docker.ListContainersOptions{All: true})
	Expect(err).NotTo(HaveOccurred())
	for _, c := range containers {
		for _, name := range c.Names {
			if strings.HasPrefix(name, "/"+n.NetworkID) {
				err := n.DockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: c.ID, Force: true})
				Expect(err).NotTo(HaveOccurred())
				break
			}
		}
	}

	images, err := n.DockerClient.ListImages(docker.ListImagesOptions{All: true})
	Expect(err).NotTo(HaveOccurred())
	for _, i := range images {
		for _, tag := range i.RepoTags {
			if strings.HasPrefix(tag, n.NetworkID) {
				err := n.DockerClient.RemoveImage(i.ID)
				Expect(err).NotTo(HaveOccurred())
				break
			}
		}
	}
}






func (n *Network) CreateAndJoinChannels(o *Orderer) {
	for _, c := range n.Channels {
		n.CreateAndJoinChannel(o, c.Name)
	}
}





func (n *Network) CreateAndJoinChannel(o *Orderer, channelName string) {
	peers := n.PeersWithChannel(channelName)
	if len(peers) == 0 {
		return
	}

	n.CreateChannel(channelName, o, peers[0])
	n.JoinChannel(channelName, o, peers...)
}




func (n *Network) UpdateChannelAnchors(o *Orderer, channelName string) {
	tempFile, err := ioutil.TempFile("", "update-anchors")
	Expect(err).NotTo(HaveOccurred())
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	peersByOrg := map[string]*Peer{}
	for _, p := range n.AnchorsForChannel(channelName) {
		peersByOrg[p.Organization] = p
	}

	for orgName, p := range peersByOrg {
		anchorUpdate := commands.OutputAnchorPeersUpdate{
			OutputAnchorPeersUpdate: tempFile.Name(),
			ChannelID:               channelName,
			Profile:                 n.ProfileForChannel(channelName),
			ConfigPath:              n.RootDir,
			AsOrg:                   orgName,
		}
		sess, err := n.ConfigTxGen(anchorUpdate)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

		sess, err = n.PeerAdminSession(p, commands.ChannelUpdate{
			ChannelID: channelName,
			Orderer:   n.OrdererAddress(o, ListenPort),
			File:      tempFile.Name(),
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	}
}






func (n *Network) CreateChannel(channelName string, o *Orderer, p *Peer) {
	createChannel := func() int {
		sess, err := n.PeerAdminSession(p, commands.ChannelCreate{
			ChannelID:   channelName,
			Orderer:     n.OrdererAddress(o, ListenPort),
			File:        n.CreateChannelTxPath(channelName),
			OutputBlock: "/dev/null",
		})
		Expect(err).NotTo(HaveOccurred())
		return sess.Wait(n.EventuallyTimeout).ExitCode()
	}

	Eventually(createChannel, n.EventuallyTimeout).Should(Equal(0))
}





func (n *Network) JoinChannel(name string, o *Orderer, peers ...*Peer) {
	if len(peers) == 0 {
		return
	}

	tempFile, err := ioutil.TempFile("", "genesis-block")
	Expect(err).NotTo(HaveOccurred())
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	sess, err := n.PeerAdminSession(peers[0], commands.ChannelFetch{
		Block:      "0",
		ChannelID:  name,
		Orderer:    n.OrdererAddress(o, ListenPort),
		OutputFile: tempFile.Name(),
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))

	for _, p := range peers {
		sess, err := n.PeerAdminSession(p, commands.ChannelJoin{
			BlockPath: tempFile.Name(),
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	}
}


func (n *Network) Cryptogen(command Command) (*gexec.Session, error) {
	cmd := NewCommand(n.Components.Cryptogen(), command)
	return n.StartSession(cmd, command.SessionName())
}


func (n *Network) ConfigTxGen(command Command) (*gexec.Session, error) {
	cmd := NewCommand(n.Components.ConfigTxGen(), command)
	return n.StartSession(cmd, command.SessionName())
}


func (n *Network) Discover(command Command) (*gexec.Session, error) {
	cmd := NewCommand(n.Components.Discover(), command)
	cmd.Args = append(cmd.Args, "--peerTLSCA", n.CACertsBundlePath())
	return n.StartSession(cmd, command.SessionName())
}


func (n *Network) ZooKeeperRunner(idx int) *runner.ZooKeeper {
	colorCode := n.nextColor()
	name := fmt.Sprintf("zookeeper-%d-%s", idx, n.NetworkID)

	return &runner.ZooKeeper{
		ZooMyID:     idx + 1, 
		Client:      n.DockerClient,
		Name:        name,
		NetworkName: n.NetworkID,
		OutputStream: gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[32m[o]\x1b[%s[%s]\x1b[0m ", colorCode, name),
			ginkgo.GinkgoWriter,
		),
		ErrorStream: gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[91m[e]\x1b[%s[%s]\x1b[0m ", colorCode, name),
			ginkgo.GinkgoWriter,
		),
	}
}

func (n *Network) minBrokersInSync() int {
	if n.Consensus.Brokers < 2 {
		return n.Consensus.Brokers
	}
	return 2
}

func (n *Network) defaultBrokerReplication() int {
	if n.Consensus.Brokers < 3 {
		return n.Consensus.Brokers
	}
	return 3
}


func (n *Network) BrokerRunner(id int, zookeepers []string) *runner.Kafka {
	colorCode := n.nextColor()
	name := fmt.Sprintf("kafka-%d-%s", id, n.NetworkID)

	return &runner.Kafka{
		BrokerID:                 id + 1,
		Client:                   n.DockerClient,
		AdvertisedListeners:      "127.0.0.1",
		HostPort:                 int(n.PortsByBrokerID[strconv.Itoa(id)][HostPort]),
		Name:                     name,
		NetworkName:              n.NetworkID,
		MinInsyncReplicas:        n.minBrokersInSync(),
		DefaultReplicationFactor: n.defaultBrokerReplication(),
		ZooKeeperConnect:         strings.Join(zookeepers, ","),
		OutputStream: gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[32m[o]\x1b[%s[%s]\x1b[0m ", colorCode, name),
			ginkgo.GinkgoWriter,
		),
		ErrorStream: gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[91m[e]\x1b[%s[%s]\x1b[0m ", colorCode, name),
			ginkgo.GinkgoWriter,
		),
	}
}



func (n *Network) BrokerGroupRunner() ifrit.Runner {
	members := grouper.Members{}
	zookeepers := []string{}

	for i := 0; i < n.Consensus.ZooKeepers; i++ {
		zk := n.ZooKeeperRunner(i)
		zookeepers = append(zookeepers, fmt.Sprintf("%s:2181", zk.Name))
		members = append(members, grouper.Member{Name: zk.Name, Runner: zk})
	}

	for i := 0; i < n.Consensus.Brokers; i++ {
		kafka := n.BrokerRunner(i, zookeepers)
		members = append(members, grouper.Member{Name: kafka.Name, Runner: kafka})
	}

	return grouper.NewOrdered(syscall.SIGTERM, members)
}



func (n *Network) OrdererRunner(o *Orderer) *ginkgomon.Runner {
	cmd := exec.Command(n.Components.Orderer())
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("FABRIC_CFG_PATH=%s", n.OrdererDir(o)))

	config := ginkgomon.Config{
		AnsiColorCode:     n.nextColor(),
		Name:              o.ID(),
		Command:           cmd,
		StartCheck:        "Beginning to serve requests",
		StartCheckTimeout: 15 * time.Second,
	}

	if n.Consensus.Brokers != 0 {
		config.StartCheck = "Start phase completed successfully"
		config.StartCheckTimeout = 30 * time.Second
	}

	return ginkgomon.New(config)
}



func (n *Network) OrdererGroupRunner() ifrit.Runner {
	members := grouper.Members{}
	for _, o := range n.Orderers {
		members = append(members, grouper.Member{Name: o.ID(), Runner: n.OrdererRunner(o)})
	}
	return grouper.NewParallel(syscall.SIGTERM, members)
}



func (n *Network) PeerRunner(p *Peer) *ginkgomon.Runner {
	cmd := n.peerCommand(
		commands.NodeStart{PeerID: p.ID()},
		fmt.Sprintf("FABRIC_CFG_PATH=%s", n.PeerDir(p)),
	)

	return ginkgomon.New(ginkgomon.Config{
		AnsiColorCode:     n.nextColor(),
		Name:              p.ID(),
		Command:           cmd,
		StartCheck:        `Started peer with ID=.*, .*, address=`,
		StartCheckTimeout: 15 * time.Second,
	})
}



func (n *Network) PeerGroupRunner() ifrit.Runner {
	members := grouper.Members{}
	for _, p := range n.Peers {
		members = append(members, grouper.Member{Name: p.ID(), Runner: n.PeerRunner(p)})
	}
	return grouper.NewParallel(syscall.SIGTERM, members)
}



func (n *Network) NetworkGroupRunner() ifrit.Runner {
	members := grouper.Members{
		{Name: "brokers", Runner: n.BrokerGroupRunner()},
		{Name: "orderers", Runner: n.OrdererGroupRunner()},
		{Name: "peers", Runner: n.PeerGroupRunner()},
	}
	return grouper.NewOrdered(syscall.SIGTERM, members)
}

func (n *Network) peerCommand(command Command, env ...string) *exec.Cmd {
	cmd := NewCommand(n.Components.Peer(), command)
	cmd.Env = append(cmd.Env, env...)
	if ConnectsToOrderer(command) {
		cmd.Args = append(cmd.Args, "--tls")
		cmd.Args = append(cmd.Args, "--cafile", n.CACertsBundlePath())
	}

	
	
	
	
	
	requiredPeerAddresses := flagCount("--peerAddresses", cmd.Args)
	for i := 0; i < requiredPeerAddresses; i++ {
		cmd.Args = append(cmd.Args, "--tlsRootCertFiles")
		cmd.Args = append(cmd.Args, n.CACertsBundlePath())
	}
	return cmd
}

func flagCount(flag string, args []string) int {
	var c int
	for _, arg := range args {
		if arg == flag {
			c++
		}
	}
	return c
}




func (n *Network) PeerAdminSession(p *Peer, command Command) (*gexec.Session, error) {
	return n.PeerUserSession(p, "Admin", command)
}




func (n *Network) PeerUserSession(p *Peer, user string, command Command) (*gexec.Session, error) {
	cmd := n.peerCommand(
		command,
		fmt.Sprintf("FABRIC_CFG_PATH=%s", n.PeerDir(p)),
		fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%s", n.PeerUserMSPDir(p, user)),
	)
	return n.StartSession(cmd, command.SessionName())
}



func (n *Network) OrdererAdminSession(o *Orderer, p *Peer, command Command) (*gexec.Session, error) {
	cmd := n.peerCommand(
		command,
		"CORE_PEER_LOCALMSPID=OrdererMSP",
		fmt.Sprintf("FABRIC_CFG_PATH=%s", n.PeerDir(p)),
		fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%s", n.OrdererUserMSPDir(o, "Admin")),
	)
	return n.StartSession(cmd, command.SessionName())
}


func (n *Network) Peer(orgName, peerName string) *Peer {
	for _, p := range n.PeersInOrg(orgName) {
		if p.Name == peerName {
			return p
		}
	}
	return nil
}


func (n *Network) DiscoveredPeer(p *Peer, chaincodes ...string) DiscoveredPeer {
	peerCert, err := ioutil.ReadFile(n.PeerCert(p))
	Expect(err).NotTo(HaveOccurred())

	return DiscoveredPeer{
		MSPID:      n.Organization(p.Organization).MSPID,
		Endpoint:   fmt.Sprintf("127.0.0.1:%d", n.PeerPort(p, ListenPort)),
		Identity:   string(peerCert),
		Chaincodes: chaincodes,
	}
}


func (n *Network) Orderer(name string) *Orderer {
	for _, o := range n.Orderers {
		if o.Name == name {
			return o
		}
	}
	return nil
}


func (n *Network) Organization(orgName string) *Organization {
	for _, org := range n.Organizations {
		if org.Name == orgName {
			return org
		}
	}
	return nil
}


func (n *Network) Consortium(name string) *Consortium {
	for _, c := range n.Consortiums {
		if c.Name == name {
			return c
		}
	}
	return nil
}


func (n *Network) PeerOrgs() []*Organization {
	orgsByName := map[string]*Organization{}
	for _, p := range n.Peers {
		orgsByName[p.Organization] = n.Organization(p.Organization)
	}

	orgs := []*Organization{}
	for _, org := range orgsByName {
		orgs = append(orgs, org)
	}
	return orgs
}



func (n *Network) PeersWithChannel(chanName string) []*Peer {
	peers := []*Peer{}
	for _, p := range n.Peers {
		for _, c := range p.Channels {
			if c.Name == chanName {
				peers = append(peers, p)
			}
		}
	}
	return peers
}



func (n *Network) AnchorsForChannel(chanName string) []*Peer {
	anchors := []*Peer{}
	for _, p := range n.Peers {
		for _, pc := range p.Channels {
			if pc.Name == chanName && pc.Anchor {
				anchors = append(anchors, p)
			}
		}
	}
	return anchors
}



func (n *Network) AnchorsInOrg(orgName string) []*Peer {
	anchors := []*Peer{}
	for _, p := range n.PeersInOrg(orgName) {
		if p.Anchor() {
			anchors = append(anchors, p)
			break
		}
	}

	
	if len(anchors) == 0 {
		anchors = n.PeersInOrg(orgName)
	}

	return anchors
}


func (n *Network) OrderersInOrg(orgName string) []*Orderer {
	orderers := []*Orderer{}
	for _, o := range n.Orderers {
		if o.Organization == orgName {
			orderers = append(orderers, o)
		}
	}
	return orderers
}



func (n *Network) OrgsForOrderers(ordererNames []string) []*Organization {
	orgsByName := map[string]*Organization{}
	for _, name := range ordererNames {
		orgName := n.Orderer(name).Organization
		orgsByName[orgName] = n.Organization(orgName)
	}
	orgs := []*Organization{}
	for _, org := range orgsByName {
		orgs = append(orgs, org)
	}
	return orgs
}



func (n *Network) OrdererOrgs() []*Organization {
	orgsByName := map[string]*Organization{}
	for _, o := range n.Orderers {
		orgsByName[o.Organization] = n.Organization(o.Organization)
	}

	orgs := []*Organization{}
	for _, org := range orgsByName {
		orgs = append(orgs, org)
	}
	return orgs
}



func (n *Network) PeersInOrg(orgName string) []*Peer {
	peers := []*Peer{}
	for _, o := range n.Peers {
		if o.Organization == orgName {
			peers = append(peers, o)
		}
	}
	return peers
}


func (n *Network) ReservePort() uint16 {
	n.StartPort++
	return n.StartPort - 1
}

type PortName string
type Ports map[PortName]uint16

const (
	ChaincodePort  PortName = "Chaincode"
	EventsPort     PortName = "Events"
	HostPort       PortName = "HostPort"
	ListenPort     PortName = "Listen"
	ProfilePort    PortName = "Profile"
	OperationsPort PortName = "Operations"
)


func PeerPortNames() []PortName {
	return []PortName{ListenPort, ChaincodePort, EventsPort, ProfilePort, OperationsPort}
}



func OrdererPortNames() []PortName {
	return []PortName{ListenPort, ProfilePort, OperationsPort}
}



func BrokerPortNames() []PortName {
	return []PortName{HostPort}
}


func (n *Network) BrokerAddresses(portName PortName) []string {
	addresses := []string{}
	for _, ports := range n.PortsByBrokerID {
		addresses = append(addresses, fmt.Sprintf("127.0.0.1:%d", ports[portName]))
	}
	return addresses
}







func (n *Network) OrdererAddress(o *Orderer, portName PortName) string {
	return fmt.Sprintf("127.0.0.1:%d", n.OrdererPort(o, portName))
}


func (n *Network) OrdererPort(o *Orderer, portName PortName) uint16 {
	ordererPorts := n.PortsByOrdererID[o.ID()]
	Expect(ordererPorts).NotTo(BeNil())
	return ordererPorts[portName]
}







func (n *Network) PeerAddress(p *Peer, portName PortName) string {
	return fmt.Sprintf("127.0.0.1:%d", n.PeerPort(p, portName))
}


func (n *Network) PeerPort(p *Peer, portName PortName) uint16 {
	peerPorts := n.PortsByPeerID[p.ID()]
	Expect(peerPorts).NotTo(BeNil())
	return peerPorts[portName]
}

func (n *Network) nextColor() string {
	color := n.colorIndex%14 + 31
	if color > 37 {
		color = color + 90 - 37
	}

	n.colorIndex++
	return fmt.Sprintf("%dm", color)
}



func (n *Network) StartSession(cmd *exec.Cmd, name string) (*gexec.Session, error) {
	ansiColorCode := n.nextColor()
	return gexec.Start(
		cmd,
		gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[32m[o]\x1b[%s[%s]\x1b[0m ", ansiColorCode, name),
			ginkgo.GinkgoWriter,
		),
		gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[91m[e]\x1b[%s[%s]\x1b[0m ", ansiColorCode, name),
			ginkgo.GinkgoWriter,
		),
	)
}

func (n *Network) GenerateCryptoConfig() {
	crypto, err := os.Create(n.CryptoConfigPath())
	Expect(err).NotTo(HaveOccurred())
	defer crypto.Close()

	t, err := template.New("crypto").Parse(n.Templates.CryptoTemplate())
	Expect(err).NotTo(HaveOccurred())

	pw := gexec.NewPrefixedWriter("[crypto-config.yaml] ", ginkgo.GinkgoWriter)
	err = t.Execute(io.MultiWriter(crypto, pw), n)
	Expect(err).NotTo(HaveOccurred())
}

func (n *Network) GenerateConfigTxConfig() {
	config, err := os.Create(n.ConfigTxConfigPath())
	Expect(err).NotTo(HaveOccurred())
	defer config.Close()

	t, err := template.New("configtx").Parse(n.Templates.ConfigTxTemplate())
	Expect(err).NotTo(HaveOccurred())

	pw := gexec.NewPrefixedWriter("[configtx.yaml] ", ginkgo.GinkgoWriter)
	err = t.Execute(io.MultiWriter(config, pw), n)
	Expect(err).NotTo(HaveOccurred())
}

func (n *Network) GenerateOrdererConfig(o *Orderer) {
	err := os.MkdirAll(n.OrdererDir(o), 0755)
	Expect(err).NotTo(HaveOccurred())

	orderer, err := os.Create(n.OrdererConfigPath(o))
	Expect(err).NotTo(HaveOccurred())
	defer orderer.Close()

	t, err := template.New("orderer").Funcs(template.FuncMap{
		"Orderer":    func() *Orderer { return o },
		"ToLower":    func(s string) string { return strings.ToLower(s) },
		"ReplaceAll": func(s, old, new string) string { return strings.Replace(s, old, new, -1) },
	}).Parse(n.Templates.OrdererTemplate())
	Expect(err).NotTo(HaveOccurred())

	pw := gexec.NewPrefixedWriter(fmt.Sprintf("[%s#orderer.yaml] ", o.ID()), ginkgo.GinkgoWriter)
	err = t.Execute(io.MultiWriter(orderer, pw), n)
	Expect(err).NotTo(HaveOccurred())
}

func (n *Network) GenerateCoreConfig(p *Peer) {
	err := os.MkdirAll(n.PeerDir(p), 0755)
	Expect(err).NotTo(HaveOccurred())

	core, err := os.Create(n.PeerConfigPath(p))
	Expect(err).NotTo(HaveOccurred())
	defer core.Close()

	t, err := template.New("peer").Funcs(template.FuncMap{
		"Peer":       func() *Peer { return p },
		"ToLower":    func(s string) string { return strings.ToLower(s) },
		"ReplaceAll": func(s, old, new string) string { return strings.Replace(s, old, new, -1) },
	}).Parse(n.Templates.CoreTemplate())
	Expect(err).NotTo(HaveOccurred())

	pw := gexec.NewPrefixedWriter(fmt.Sprintf("[%s#core.yaml] ", p.ID()), ginkgo.GinkgoWriter)
	err = t.Execute(io.MultiWriter(core, pw), n)
	Expect(err).NotTo(HaveOccurred())
}
