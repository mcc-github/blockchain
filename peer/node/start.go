/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	ccdef "github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/deliver"
	"github.com/mcc-github/blockchain/common/flogging"
	floggingmetrics "github.com/mcc-github/blockchain/common/flogging/metrics"
	"github.com/mcc-github/blockchain/common/grpclogging"
	"github.com/mcc-github/blockchain/common/grpcmetrics"
	"github.com/mcc-github/blockchain/common/localmsp"
	"github.com/mcc-github/blockchain/common/metadata"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/viperutil"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/admin"
	cc "github.com/mcc-github/blockchain/core/cclifecycle"
	"github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/car"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/java"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/node"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/dockercontroller"
	"github.com/mcc-github/blockchain/core/container/inproccontroller"
	"github.com/mcc-github/blockchain/core/dispatcher"
	"github.com/mcc-github/blockchain/core/endorser"
	authHandler "github.com/mcc-github/blockchain/core/handlers/auth"
	endorsement2 "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	endorsement3 "github.com/mcc-github/blockchain/core/handlers/endorsement/api/identities"
	"github.com/mcc-github/blockchain/core/handlers/library"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/operations"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/scc/cscc"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	"github.com/mcc-github/blockchain/core/scc/qscc"
	"github.com/mcc-github/blockchain/discovery"
	"github.com/mcc-github/blockchain/discovery/endorsement"
	discsupport "github.com/mcc-github/blockchain/discovery/support"
	discacl "github.com/mcc-github/blockchain/discovery/support/acl"
	ccsupport "github.com/mcc-github/blockchain/discovery/support/chaincode"
	"github.com/mcc-github/blockchain/discovery/support/config"
	"github.com/mcc-github/blockchain/discovery/support/gossip"
	gossipcommon "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/service"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	peergossip "github.com/mcc-github/blockchain/peer/gossip"
	"github.com/mcc-github/blockchain/peer/version"
	cb "github.com/mcc-github/blockchain/protos/common"
	common2 "github.com/mcc-github/blockchain/protos/common"
	discprotos "github.com/mcc-github/blockchain/protos/discovery"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/mcc-github/blockchain/token/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	chaincodeAddrKey       = "peer.chaincodeAddress"
	chaincodeListenAddrKey = "peer.chaincodeListenAddress"
	defaultChaincodePort   = 7052
	grpcMaxConcurrency     = 2500
)

var chaincodeDevMode bool

func startCmd() *cobra.Command {
	
	flags := nodeStartCmd.Flags()
	flags.BoolVarP(&chaincodeDevMode, "peer-chaincodedev", "", false,
		"Whether peer in chaincode development mode")

	return nodeStartCmd
}

var nodeStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the node.",
	Long:  `Starts a node that interacts with the network.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("trailing args detected")
		}
		
		cmd.SilenceUsage = true
		return serve(args)
	},
}

func serve(args []string) error {
	
	
	
	
	
	
	mspType := mgmt.GetLocalMSP().GetType()
	if mspType != msp.FABRIC {
		panic("Unsupported msp type " + msp.ProviderTypeToString(mspType))
	}

	
	
	
	grpc.EnableTracing = true

	logger.Infof("Starting %s", version.GetInfo())

	
	
	aclProvider := aclmgmt.NewACLProvider(
		aclmgmt.ResourceGetter(peer.GetStableChannelConfig),
	)

	pr := platforms.NewRegistry(
		&golang.Platform{},
		&node.Platform{},
		&java.Platform{},
		&car.Platform{},
	)

	deployedCCInfoProvider := &lscc.DeployedCCInfoProvider{}

	identityDeserializerFactory := func(chainID string) msp.IdentityDeserializer {
		return mgmt.GetManagerForChain(chainID)
	}

	opsSystem := newOperationsSystem()
	err := opsSystem.Start()
	if err != nil {
		return errors.WithMessage(err, "failed to initialize operations subystems")
	}
	defer opsSystem.Stop()

	metricsProvider := opsSystem.Provider
	logObserver := floggingmetrics.NewObserver(metricsProvider)
	flogging.Global.SetObserver(logObserver)

	membershipInfoProvider := privdata.NewMembershipInfoProvider(createSelfSignedData(), identityDeserializerFactory)
	
	ledgermgmt.Initialize(
		&ledgermgmt.Initializer{
			CustomTxProcessors:            peer.ConfigTxProcessors,
			PlatformRegistry:              pr,
			DeployedChaincodeInfoProvider: deployedCCInfoProvider,
			MembershipInfoProvider:        membershipInfoProvider,
			MetricsProvider:               metricsProvider,
		},
	)

	
	
	if chaincodeDevMode {
		logger.Info("Running in chaincode development mode")
		logger.Info("Disable loading validity system chaincode")

		viper.Set("chaincode.mode", chaincode.DevModeUserRunsChaincode)
	}

	if err := peer.CacheConfiguration(); err != nil {
		return err
	}

	peerEndpoint, err := peer.GetPeerEndpoint()
	if err != nil {
		err = fmt.Errorf("Failed to get Peer Endpoint: %s", err)
		return err
	}

	peerHost, _, err := net.SplitHostPort(peerEndpoint.Address)
	if err != nil {
		return fmt.Errorf("peer address is not in the format of host:port: %v", err)
	}

	listenAddr := viper.GetString("peer.listenAddress")
	serverConfig, err := peer.GetServerConfig()
	if err != nil {
		logger.Fatalf("Error loading secure config for peer (%s)", err)
	}

	throttle := comm.NewThrottle(grpcMaxConcurrency)
	serverConfig.Logger = flogging.MustGetLogger("core.comm").With("server", "PeerServer")
	serverConfig.MetricsProvider = metricsProvider
	serverConfig.UnaryInterceptors = append(
		serverConfig.UnaryInterceptors,
		grpcmetrics.UnaryServerInterceptor(grpcmetrics.NewUnaryMetrics(metricsProvider)),
		grpclogging.UnaryServerInterceptor(flogging.MustGetLogger("comm.grpc.server").Zap()),
		throttle.UnaryServerIntercptor,
	)
	serverConfig.StreamInterceptors = append(
		serverConfig.StreamInterceptors,
		grpcmetrics.StreamServerInterceptor(grpcmetrics.NewStreamMetrics(metricsProvider)),
		grpclogging.StreamServerInterceptor(flogging.MustGetLogger("comm.grpc.server").Zap()),
		throttle.StreamServerInterceptor,
	)

	peerServer, err := peer.NewPeerServer(listenAddr, serverConfig)
	if err != nil {
		logger.Fatalf("Failed to create peer server (%s)", err)
	}

	if serverConfig.SecOpts.UseTLS {
		logger.Info("Starting peer with TLS enabled")
		
		cs := comm.GetCredentialSupport()
		cs.ServerRootCAs = serverConfig.SecOpts.ServerRootCAs

		
		clientCert, err := peer.GetClientCertificate()
		if err != nil {
			logger.Fatalf("Failed to set TLS client certificate (%s)", err)
		}
		comm.GetCredentialSupport().SetClientCertificate(clientCert)
	}

	mutualTLS := serverConfig.SecOpts.UseTLS && serverConfig.SecOpts.RequireClientCert
	policyCheckerProvider := func(resourceName string) deliver.PolicyCheckerFunc {
		return func(env *cb.Envelope, channelID string) error {
			return aclProvider.CheckACL(resourceName, channelID, env)
		}
	}

	abServer := peer.NewDeliverEventsServer(mutualTLS, policyCheckerProvider, &peer.DeliverChainManager{}, metricsProvider)
	pb.RegisterDeliverServer(peerServer.Server(), abServer)

	
	chaincodeSupport, ccp, sccp, packageProvider := startChaincodeServer(peerHost, aclProvider, pr, opsSystem, deployedCCInfoProvider)

	logger.Debugf("Running peer")

	
	startAdminServer(listenAddr, peerServer.Server(), metricsProvider)

	privDataDist := func(channel string, txID string, privateData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error {
		return service.GetGossipService().DistributePrivateData(channel, txID, privateData, blkHt)
	}

	signingIdentity := mgmt.GetLocalSigningIdentityOrPanic()
	serializedIdentity, err := signingIdentity.Serialize()
	if err != nil {
		logger.Panicf("Failed serializing self identity: %v", err)
	}

	libConf := library.Config{}
	if err = viperutil.EnhancedExactUnmarshalKey("peer.handlers", &libConf); err != nil {
		return errors.WithMessage(err, "could not load YAML config")
	}
	reg := library.InitRegistry(libConf)

	authFilters := reg.Lookup(library.Auth).([]authHandler.Filter)
	endorserSupport := &endorser.SupportImpl{
		SignerSupport:    signingIdentity,
		Peer:             peer.Default,
		PeerSupport:      peer.DefaultSupport,
		ChaincodeSupport: chaincodeSupport,
		SysCCProvider:    sccp,
		ACLProvider:      aclProvider,
	}
	endorsementPluginsByName := reg.Lookup(library.Endorsement).(map[string]endorsement2.PluginFactory)
	validationPluginsByName := reg.Lookup(library.Validation).(map[string]validation.PluginFactory)
	signingIdentityFetcher := (endorsement3.SigningIdentityFetcher)(endorserSupport)
	channelStateRetriever := endorser.ChannelStateRetriever(endorserSupport)
	pluginMapper := endorser.MapBasedPluginMapper(endorsementPluginsByName)
	pluginEndorser := endorser.NewPluginEndorser(&endorser.PluginSupport{
		ChannelStateRetriever:   channelStateRetriever,
		TransientStoreRetriever: peer.TransientStoreFactory,
		PluginMapper:            pluginMapper,
		SigningIdentityFetcher:  signingIdentityFetcher,
	})
	endorserSupport.PluginEndorser = pluginEndorser
	serverEndorser := endorser.NewEndorserServer(privDataDist, endorserSupport, pr, metricsProvider)
	auth := authHandler.ChainFilters(serverEndorser, authFilters...)
	
	pb.RegisterEndorserServer(peerServer.Server(), auth)

	policyMgr := peer.NewChannelPolicyManagerGetter()

	
	err = initGossipService(policyMgr, peerServer, serializedIdentity, peerEndpoint.Address)
	if err != nil {
		return err
	}
	defer service.GetGossipService().Stop()

	
	err = registerProverService(peerServer, aclProvider, signingIdentity)
	if err != nil {
		return err
	}

	

	
	sccp.DeploySysCCs("", ccp)
	logger.Infof("Deployed system chaincodes")

	installedCCs := func() ([]ccdef.InstalledChaincode, error) {
		return packageProvider.ListInstalledChaincodes()
	}
	lifecycle, err := cc.NewLifeCycle(cc.Enumerate(installedCCs))
	if err != nil {
		logger.Panicf("Failed creating lifecycle: +%v", err)
	}
	onUpdate := cc.HandleMetadataUpdate(func(channel string, chaincodes ccdef.MetadataSet) {
		service.GetGossipService().UpdateChaincodes(chaincodes.AsChaincodes(), gossipcommon.ChainID(channel))
	})
	lifecycle.AddListener(onUpdate)

	
	peer.Initialize(func(cid string) {
		logger.Debugf("Deploying system CC, for channel <%s>", cid)
		sccp.DeploySysCCs(cid, ccp)
		sub, err := lifecycle.NewChannelSubscription(cid, cc.QueryCreatorFunc(func() (cc.Query, error) {
			return peer.GetLedger(cid).NewQueryExecutor()
		}))
		if err != nil {
			logger.Panicf("Failed subscribing to chaincode lifecycle updates")
		}
		cceventmgmt.GetMgr().Register(cid, sub)
	}, ccp, sccp, txvalidator.MapBasedPluginMapper(validationPluginsByName),
		pr, deployedCCInfoProvider, membershipInfoProvider, metricsProvider)

	if viper.GetBool("peer.discovery.enabled") {
		registerDiscoveryService(peerServer, policyMgr, lifecycle)
	}

	networkID := viper.GetString("peer.networkId")

	logger.Infof("Starting peer with ID=[%s], network ID=[%s], address=[%s]", peerEndpoint.Id, networkID, peerEndpoint.Address)

	
	
	profileEnabled := viper.GetBool("peer.profile.enabled")
	profileListenAddress := viper.GetString("peer.profile.listenAddress")

	
	
	serve := make(chan error)

	go func() {
		var grpcErr error
		if grpcErr = peerServer.Start(); grpcErr != nil {
			grpcErr = fmt.Errorf("grpc server exited with error: %s", grpcErr)
		} else {
			logger.Info("peer server exited")
		}
		serve <- grpcErr
	}()

	
	if profileEnabled {
		go func() {
			logger.Infof("Starting profiling server with listenAddress = %s", profileListenAddress)
			if profileErr := http.ListenAndServe(profileListenAddress, nil); profileErr != nil {
				logger.Errorf("Error starting profiler: %s", profileErr)
			}
		}()
	}

	go handleSignals(addPlatformSignals(map[os.Signal]func(){
		syscall.SIGINT:  func() { serve <- nil },
		syscall.SIGTERM: func() { serve <- nil },
	}))

	logger.Infof("Started peer with ID=[%s], network ID=[%s], address=[%s]", peerEndpoint.Id, networkID, peerEndpoint.Address)

	
	return <-serve
}

func handleSignals(handlers map[os.Signal]func()) {
	var signals []os.Signal
	for sig := range handlers {
		signals = append(signals, sig)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)

	for sig := range signalChan {
		logger.Infof("Received signal: %d (%s)", sig, sig)
		handlers[sig]()
	}
}

func localPolicy(policyObject proto.Message) policies.Policy {
	localMSP := mgmt.GetLocalMSP()
	pp := cauthdsl.NewPolicyProvider(localMSP)
	policy, _, err := pp.NewPolicy(utils.MarshalOrPanic(policyObject))
	if err != nil {
		logger.Panicf("Failed creating local policy: +%v", err)
	}
	return policy
}

func createSelfSignedData() common2.SignedData {
	sId := mgmt.GetLocalSigningIdentityOrPanic()
	msg := make([]byte, 32)
	sig, err := sId.Sign(msg)
	if err != nil {
		logger.Panicf("Failed creating self signed data because message signing failed: %v", err)
	}
	peerIdentity, err := sId.Serialize()
	if err != nil {
		logger.Panicf("Failed creating self signed data because peer identity couldn't be serialized: %v", err)
	}
	return common2.SignedData{
		Data:      msg,
		Signature: sig,
		Identity:  peerIdentity,
	}
}

func registerDiscoveryService(peerServer *comm.GRPCServer, polMgr policies.ChannelPolicyManagerGetter, lc *cc.Lifecycle) {
	mspID := viper.GetString("peer.localMspId")
	localAccessPolicy := localPolicy(cauthdsl.SignedByAnyAdmin([]string{mspID}))
	if viper.GetBool("peer.discovery.orgMembersAllowedAccess") {
		localAccessPolicy = localPolicy(cauthdsl.SignedByAnyMember([]string{mspID}))
	}
	channelVerifier := discacl.NewChannelVerifier(policies.ChannelApplicationWriters, polMgr)
	acl := discacl.NewDiscoverySupport(channelVerifier, localAccessPolicy, discacl.ChannelConfigGetterFunc(peer.GetStableChannelConfig))
	gSup := gossip.NewDiscoverySupport(service.GetGossipService())
	ccSup := ccsupport.NewDiscoverySupport(lc)
	ea := endorsement.NewEndorsementAnalyzer(gSup, ccSup, acl, lc)
	confSup := config.NewDiscoverySupport(config.CurrentConfigBlockGetterFunc(peer.GetCurrConfigBlock))
	support := discsupport.NewDiscoverySupport(acl, gSup, ea, confSup, acl)
	svc := discovery.NewService(discovery.Config{
		TLS:                          peerServer.TLSEnabled(),
		AuthCacheEnabled:             viper.GetBool("peer.discovery.authCacheEnabled"),
		AuthCacheMaxSize:             viper.GetInt("peer.discovery.authCacheMaxSize"),
		AuthCachePurgeRetentionRatio: viper.GetFloat64("peer.discovery.authCachePurgeRetentionRatio"),
	}, support)
	logger.Info("Discovery service activated")
	discprotos.RegisterDiscoveryServer(peerServer.Server(), svc)
}


func createChaincodeServer(ca tlsgen.CA, peerHostname string) (srv *comm.GRPCServer, ccEndpoint string, err error) {
	
	ccEndpoint, err = computeChaincodeEndpoint(peerHostname)
	if err != nil {
		if chaincode.IsDevMode() {
			
			ccEndpoint = fmt.Sprintf("%s:%d", "0.0.0.0", defaultChaincodePort)
			logger.Warningf("use %s as chaincode endpoint because of error in computeChaincodeEndpoint: %s", ccEndpoint, err)
		} else {
			
			logger.Errorf("Error computing chaincode endpoint: %s", err)
			return nil, "", err
		}
	}

	host, _, err := net.SplitHostPort(ccEndpoint)
	if err != nil {
		logger.Panic("Chaincode service host", ccEndpoint, "isn't a valid hostname:", err)
	}

	cclistenAddress := viper.GetString(chaincodeListenAddrKey)
	if cclistenAddress == "" {
		cclistenAddress = fmt.Sprintf("%s:%d", peerHostname, defaultChaincodePort)
		logger.Warningf("%s is not set, using %s", chaincodeListenAddrKey, cclistenAddress)
		viper.Set(chaincodeListenAddrKey, cclistenAddress)
	}

	config, err := peer.GetServerConfig()
	if err != nil {
		logger.Errorf("Error getting server config: %s", err)
		return nil, "", err
	}

	
	config.Logger = flogging.MustGetLogger("core.comm").With("server", "ChaincodeServer")

	
	if config.SecOpts.UseTLS {
		
		certKeyPair, err := ca.NewServerCertKeyPair(host)
		if err != nil {
			logger.Panicf("Failed generating TLS certificate for chaincode service: +%v", err)
		}
		config.SecOpts = &comm.SecureOptions{
			UseTLS: true,
			
			RequireClientCert: true,
			
			ClientRootCAs: [][]byte{ca.CertBytes()},
			
			Certificate: certKeyPair.Cert,
			Key:         certKeyPair.Key,
			
			
			ServerRootCAs: nil,
		}
	}

	
	chaincodeKeepaliveOptions := &comm.KeepaliveOptions{
		ServerInterval:    time.Duration(2) * time.Hour,    
		ServerTimeout:     time.Duration(20) * time.Second, 
		ServerMinInterval: time.Duration(1) * time.Minute,  
	}
	config.KaOpts = chaincodeKeepaliveOptions

	srv, err = comm.NewGRPCServer(cclistenAddress, config)
	if err != nil {
		logger.Errorf("Error creating GRPC server: %s", err)
		return nil, "", err
	}

	return srv, ccEndpoint, nil
}








func computeChaincodeEndpoint(peerHostname string) (ccEndpoint string, err error) {
	logger.Infof("Entering computeChaincodeEndpoint with peerHostname: %s", peerHostname)
	
	
	
	ccEndpoint = viper.GetString(chaincodeAddrKey)
	if ccEndpoint == "" {
		
		
		ccEndpoint = viper.GetString(chaincodeListenAddrKey)
		if ccEndpoint == "" {
			
			peerIp := net.ParseIP(peerHostname)
			if peerIp != nil && peerIp.IsUnspecified() {
				
				logger.Errorf("ChaincodeAddress and chaincodeListenAddress are nil and peerIP is %s", peerIp)
				return "", errors.New("invalid endpoint for chaincode to connect")
			}

			
			ccEndpoint = fmt.Sprintf("%s:%d", peerHostname, defaultChaincodePort)

		} else {
			
			host, port, err := net.SplitHostPort(ccEndpoint)
			if err != nil {
				logger.Errorf("ChaincodeAddress is nil and fail to split chaincodeListenAddress: %s", err)
				return "", err
			}

			ccListenerIp := net.ParseIP(host)
			
			
			if ccListenerIp != nil && ccListenerIp.IsUnspecified() {
				
				peerIp := net.ParseIP(peerHostname)
				if peerIp != nil && peerIp.IsUnspecified() {
					
					logger.Error("ChaincodeAddress is nil while both chaincodeListenAddressIP and peerIP are 0.0.0.0")
					return "", errors.New("invalid endpoint for chaincode to connect")
				}
				ccEndpoint = fmt.Sprintf("%s:%s", peerHostname, port)
			}

		}

	} else {
		
		if host, _, err := net.SplitHostPort(ccEndpoint); err != nil {
			logger.Errorf("Fail to split chaincodeAddress: %s", err)
			return "", err
		} else {
			ccIP := net.ParseIP(host)
			if ccIP != nil && ccIP.IsUnspecified() {
				logger.Errorf("ChaincodeAddress' IP cannot be %s in non-dev mode", ccIP)
				return "", errors.New("invalid endpoint for chaincode to connect")
			}
		}
	}

	logger.Infof("Exit with ccEndpoint: %s", ccEndpoint)
	return ccEndpoint, nil
}




func registerChaincodeSupport(
	grpcServer *comm.GRPCServer,
	ccEndpoint string,
	ca tlsgen.CA,
	packageProvider *persistence.PackageProvider,
	aclProvider aclmgmt.ACLProvider,
	pr *platforms.Registry,
	lifecycleSCC *lifecycle.SCC,
	ops *operations.System,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
) (*chaincode.ChaincodeSupport, ccprovider.ChaincodeProvider, *scc.Provider) {
	
	userRunsCC := chaincode.IsDevMode()
	tlsEnabled := viper.GetBool("peer.tls.enabled")

	authenticator := accesscontrol.NewAuthenticator(ca)
	ipRegistry := inproccontroller.NewRegistry()

	sccp := scc.NewProvider(peer.Default, peer.DefaultSupport, ipRegistry)
	lsccInst := lscc.New(sccp, aclProvider, pr)

	dockerProvider := dockercontroller.NewProvider(
		viper.GetString("peer.id"),
		viper.GetString("peer.networkId"),
		ops.Provider,
	)
	dockerVM := dockercontroller.NewDockerVM(
		dockerProvider.PeerID,
		dockerProvider.NetworkID,
		dockerProvider.BuildMetrics,
	)

	err := ops.RegisterChecker("docker", dockerVM)
	if err != nil {
		logger.Panicf("failed to register docker health check: %s", err)
	}

	chaincodeSupport := chaincode.NewChaincodeSupport(
		chaincode.GlobalConfig(),
		ccEndpoint,
		userRunsCC,
		ca.CertBytes(),
		authenticator,
		packageProvider,
		lsccInst,
		aclProvider,
		container.NewVMController(
			map[string]container.VMProvider{
				dockercontroller.ContainerType: dockerProvider,
				inproccontroller.ContainerType: ipRegistry,
			},
		),
		sccp,
		pr,
		peer.DefaultSupport,
		ops.Provider,
		deployedCCInfoProvider,
	)
	ipRegistry.ChaincodeSupport = chaincodeSupport
	ccp := chaincode.NewProvider(chaincodeSupport)

	ccSrv := pb.ChaincodeSupportServer(chaincodeSupport)
	if tlsEnabled {
		ccSrv = authenticator.Wrap(ccSrv)
	}

	csccInst := cscc.New(ccp, sccp, aclProvider, deployedCCInfoProvider)
	qsccInst := qscc.New(aclProvider)

	
	sccs := scc.CreatePluginSysCCs(sccp)
	for _, cc := range append([]scc.SelfDescribingSysCC{lsccInst, csccInst, qsccInst, lifecycleSCC}, sccs...) {
		sccp.RegisterSysCC(cc)
	}
	pb.RegisterChaincodeSupportServer(grpcServer.Server(), ccSrv)

	return chaincodeSupport, ccp, sccp
}





func startChaincodeServer(
	peerHost string,
	aclProvider aclmgmt.ACLProvider,
	pr *platforms.Registry,
	ops *operations.System,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
) (*chaincode.ChaincodeSupport, ccprovider.ChaincodeProvider, *scc.Provider, *persistence.PackageProvider) {
	
	chaincodeInstallPath := ccprovider.GetChaincodeInstallPathFromViper()
	ccprovider.SetChaincodesPath(chaincodeInstallPath)

	ccPackageParser := &persistence.ChaincodePackageParser{}
	ccStore := &persistence.Store{
		Path:       chaincodeInstallPath,
		ReadWriter: &persistence.FilesystemIO{},
	}

	packageProvider := &persistence.PackageProvider{
		LegacyPP: &ccprovider.CCInfoFSImpl{},
		Store:    ccStore,
	}

	lifecycleSCC := &lifecycle.SCC{
		Dispatcher: &dispatcher.Dispatcher{
			Protobuf: &dispatcher.ProtobufImpl{},
		},
		Functions: &lifecycle.Lifecycle{
			PackageParser:  ccPackageParser,
			ChaincodeStore: ccStore,
		},
	}

	
	ca, err := tlsgen.NewCA()
	if err != nil {
		logger.Panic("Failed creating authentication layer:", err)
	}
	ccSrv, ccEndpoint, err := createChaincodeServer(ca, peerHost)
	if err != nil {
		logger.Panicf("Failed to create chaincode server: %s", err)
	}
	chaincodeSupport, ccp, sccp := registerChaincodeSupport(
		ccSrv,
		ccEndpoint,
		ca,
		packageProvider,
		aclProvider,
		pr,
		lifecycleSCC,
		ops,
		deployedCCInfoProvider,
	)
	go ccSrv.Start()
	return chaincodeSupport, ccp, sccp, packageProvider
}

func adminHasSeparateListener(peerListenAddr string, adminListenAddress string) bool {
	
	if adminListenAddress == "" {
		return false
	}
	_, peerPort, err := net.SplitHostPort(peerListenAddr)
	if err != nil {
		logger.Panicf("Failed parsing peer listen address")
	}

	_, adminPort, err := net.SplitHostPort(adminListenAddress)
	if err != nil {
		logger.Panicf("Failed parsing admin listen address")
	}
	
	
	return adminPort != peerPort
}

func startAdminServer(peerListenAddr string, peerServer *grpc.Server, metricsProvider metrics.Provider) {
	adminListenAddress := viper.GetString("peer.adminService.listenAddress")
	separateLsnrForAdmin := adminHasSeparateListener(peerListenAddr, adminListenAddress)
	mspID := viper.GetString("peer.localMspId")
	adminPolicy := localPolicy(cauthdsl.SignedByAnyAdmin([]string{mspID}))
	gRPCService := peerServer
	if separateLsnrForAdmin {
		logger.Info("Creating gRPC server for admin service on", adminListenAddress)
		serverConfig, err := peer.GetServerConfig()
		if err != nil {
			logger.Fatalf("Error loading secure config for admin service (%s)", err)
		}
		throttle := comm.NewThrottle(grpcMaxConcurrency)
		serverConfig.Logger = flogging.MustGetLogger("core.comm").With("server", "AdminServer")
		serverConfig.MetricsProvider = metricsProvider
		serverConfig.UnaryInterceptors = append(
			serverConfig.UnaryInterceptors,
			grpcmetrics.UnaryServerInterceptor(grpcmetrics.NewUnaryMetrics(metricsProvider)),
			grpclogging.UnaryServerInterceptor(flogging.MustGetLogger("comm.grpc.server").Zap()),
			throttle.UnaryServerIntercptor,
		)
		serverConfig.StreamInterceptors = append(
			serverConfig.StreamInterceptors,
			grpcmetrics.StreamServerInterceptor(grpcmetrics.NewStreamMetrics(metricsProvider)),
			grpclogging.StreamServerInterceptor(flogging.MustGetLogger("comm.grpc.server").Zap()),
			throttle.StreamServerInterceptor,
		)
		adminServer, err := peer.NewPeerServer(adminListenAddress, serverConfig)
		if err != nil {
			logger.Fatalf("Failed to create admin server (%s)", err)
		}
		gRPCService = adminServer.Server()
		defer func() {
			go adminServer.Start()
		}()
	}

	pb.RegisterAdminServer(gRPCService, admin.NewAdminServer(adminPolicy))
}


func secureDialOpts() []grpc.DialOption {
	var dialOpts []grpc.DialOption
	
	dialOpts = append(
		dialOpts,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(comm.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(comm.MaxSendMsgSize)))
	
	kaOpts := comm.DefaultKeepaliveOptions
	if viper.IsSet("peer.keepalive.client.interval") {
		kaOpts.ClientInterval = viper.GetDuration("peer.keepalive.client.interval")
	}
	if viper.IsSet("peer.keepalive.client.timeout") {
		kaOpts.ClientTimeout = viper.GetDuration("peer.keepalive.client.timeout")
	}
	dialOpts = append(dialOpts, comm.ClientKeepaliveOptions(kaOpts)...)

	if viper.GetBool("peer.tls.enabled") {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(comm.GetCredentialSupport().GetPeerCredentials()))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	return dialOpts
}






func initGossipService(policyMgr policies.ChannelPolicyManagerGetter, peerServer *comm.GRPCServer, serializedIdentity []byte, peerAddr string) error {
	var certs *gossipcommon.TLSCertificates
	if peerServer.TLSEnabled() {
		serverCert := peerServer.ServerCertificate()
		clientCert, err := peer.GetClientCertificate()
		if err != nil {
			return errors.Wrap(err, "failed obtaining client certificates")
		}
		certs = &gossipcommon.TLSCertificates{}
		certs.TLSServerCert.Store(&serverCert)
		certs.TLSClientCert.Store(&clientCert)
	}

	messageCryptoService := peergossip.NewMCS(
		policyMgr,
		localmsp.NewSigner(),
		mgmt.NewDeserializersManager(),
	)
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	bootstrap := viper.GetStringSlice("peer.gossip.bootstrap")

	return service.InitGossipService(
		serializedIdentity,
		peerAddr,
		peerServer.Server(),
		certs,
		messageCryptoService,
		secAdv,
		secureDialOpts,
		bootstrap...,
	)
}

func newOperationsSystem() *operations.System {
	return operations.NewSystem(operations.Options{
		Logger:        flogging.MustGetLogger("peer.operations"),
		ListenAddress: viper.GetString("operations.listenAddress"),
		Metrics: operations.MetricsOptions{
			Provider: viper.GetString("metrics.provider"),
			Statsd: &operations.Statsd{
				Network:       viper.GetString("metrics.statsd.network"),
				Address:       viper.GetString("metrics.statsd.address"),
				WriteInterval: viper.GetDuration("metrics.statsd.writeInterval"),
				Prefix:        viper.GetString("metrics.statsd.prefix"),
			},
		},
		TLS: operations.TLS{
			Enabled:            viper.GetBool("operations.tls.enabled"),
			CertFile:           viper.GetString("operations.tls.cert.file"),
			KeyFile:            viper.GetString("operations.tls.key.file"),
			ClientCertRequired: viper.GetBool("operations.tls.clientAuthRequired"),
			ClientCACertFiles:  viper.GetStringSlice("operations.tls.clientRootCAs.files"),
		},
		Version: metadata.Version,
	})
}

func registerProverService(peerServer *comm.GRPCServer, aclProvider aclmgmt.ACLProvider, signingIdentity msp.SigningIdentity) error {
	policyChecker := &server.PolicyBasedAccessControl{
		ACLProvider: aclProvider,
		ACLResources: &server.ACLResources{
			IssueTokens:    resources.Token_Issue,
			TransferTokens: resources.Token_Transfer,
			ListTokens:     resources.Token_List,
		},
	}

	responseMarshaler, err := server.NewResponseMarshaler(signingIdentity)
	if err != nil {
		logger.Errorf("Failed to create prover service: %s", err)
		return err
	}

	prover := &server.Prover{
		CapabilityChecker: &server.TokenCapabilityChecker{
			PeerOps: peer.Default,
		},
		Marshaler:     responseMarshaler,
		PolicyChecker: policyChecker,
		TMSManager: &server.Manager{
			LedgerManager: &server.PeerLedgerManager{},
		},
	}
	token.RegisterProverServer(peerServer.Server(), prover)
	return nil
}
