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
	"path/filepath"
	"syscall"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	ccdef "github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/deliver"
	"github.com/mcc-github/blockchain/common/flogging"
	floggingmetrics "github.com/mcc-github/blockchain/common/flogging/metrics"
	"github.com/mcc-github/blockchain/common/grpclogging"
	"github.com/mcc-github/blockchain/common/grpcmetrics"
	"github.com/mcc-github/blockchain/common/metadata"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/viperutil"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/admin"
	"github.com/mcc-github/blockchain/core/cclifecycle"
	"github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	coreconfig "github.com/mcc-github/blockchain/core/config"
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
	peergossip "github.com/mcc-github/blockchain/internal/peer/gossip"
	"github.com/mcc-github/blockchain/internal/peer/version"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	cb "github.com/mcc-github/blockchain/protos/common"
	discprotos "github.com/mcc-github/blockchain/protos/discovery"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/mcc-github/blockchain/token/server"
	"github.com/mcc-github/blockchain/token/tms/manager"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	chaincodeAddrKey       = "peer.chaincodeAddress"
	chaincodeListenAddrKey = "peer.chaincodeListenAddress"
	defaultChaincodePort   = 7052
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

	
	coreConfig, err := peer.GlobalConfig()
	if err != nil {
		return err
	}

	platformRegistry := platforms.NewRegistry(platforms.SupportedPlatforms...)

	identityDeserializerFactory := func(chainID string) msp.IdentityDeserializer {
		return mgmt.GetManagerForChain(chainID)
	}

	opsSystem := newOperationsSystem(coreConfig)
	err = opsSystem.Start()
	if err != nil {
		return errors.WithMessage(err, "failed to initialize operations subystems")
	}
	defer opsSystem.Stop()

	metricsProvider := opsSystem.Provider
	logObserver := floggingmetrics.NewObserver(metricsProvider)
	flogging.Global.SetObserver(logObserver)

	membershipInfoProvider := privdata.NewMembershipInfoProvider(createSelfSignedData(), identityDeserializerFactory)

	mspID := coreConfig.LocalMSPID

	chaincodeInstallPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "lifecycle", "chaincodes")
	ccStore := persistence.NewStore(chaincodeInstallPath)
	ccPackageParser := &persistence.ChaincodePackageParser{}

	
	
	
	
	
	
	
	lifecycleResources := &lifecycle.Resources{
		Serializer:          &lifecycle.Serializer{},
		ChannelConfigSource: peer.Default,
		ChaincodeStore:      ccStore,
		PackageParser:       ccPackageParser,
	}

	lifecycleValidatorCommitter := &lifecycle.ValidatorCommitter{
		Resources:                    lifecycleResources,
		LegacyDeployedCCInfoProvider: &lscc.DeployedCCInfoProvider{},
	}

	lifecycleCache := lifecycle.NewCache(lifecycleResources, mspID)

	
	ledgermgmt.Initialize(
		&ledgermgmt.Initializer{
			CustomTxProcessors:            peer.ConfigTxProcessors,
			PlatformRegistry:              platformRegistry,
			DeployedChaincodeInfoProvider: lifecycleValidatorCommitter,
			MembershipInfoProvider:        membershipInfoProvider,
			MetricsProvider:               metricsProvider,
			HealthCheckRegistry:           opsSystem,
			StateListeners:                []ledger.StateListener{lifecycleCache},
			Config:                        ledgerConfig(),
		},
	)

	
	lsccInstallPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "chaincodes")
	ccprovider.SetChaincodesPath(lsccInstallPath)

	if err := lifecycleCache.InitializeLocalChaincodes(); err != nil {
		return errors.WithMessage(err, "could not initialize local chaincodes")
	}

	packageProvider := &persistence.PackageProvider{
		LegacyPP: &ccprovider.CCInfoFSImpl{},
		Store:    ccStore,
		Parser:   ccPackageParser,
	}

	
	
	if chaincodeDevMode {
		logger.Info("Running in chaincode development mode")
		logger.Info("Disable loading validity system chaincode")

		viper.Set("chaincode.mode", chaincode.DevModeUserRunsChaincode)
	}

	peerHost, _, err := net.SplitHostPort(coreConfig.PeerAddress)
	if err != nil {
		return fmt.Errorf("peer address is not in the format of host:port: %v", err)
	}

	listenAddr := coreConfig.ListenAddress
	serverConfig, err := peer.GetServerConfig()
	if err != nil {
		logger.Fatalf("Error loading secure config for peer (%s)", err)
	}

	serverConfig.Logger = flogging.MustGetLogger("core.comm").With("server", "PeerServer")
	serverConfig.MetricsProvider = metricsProvider
	serverConfig.UnaryInterceptors = append(
		serverConfig.UnaryInterceptors,
		grpcmetrics.UnaryServerInterceptor(grpcmetrics.NewUnaryMetrics(metricsProvider)),
		grpclogging.UnaryServerInterceptor(flogging.MustGetLogger("comm.grpc.server").Zap()),
	)
	serverConfig.StreamInterceptors = append(
		serverConfig.StreamInterceptors,
		grpcmetrics.StreamServerInterceptor(grpcmetrics.NewStreamMetrics(metricsProvider)),
		grpclogging.StreamServerInterceptor(flogging.MustGetLogger("comm.grpc.server").Zap()),
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

	metrics := deliver.NewMetrics(metricsProvider)
	abServer := &peer.DeliverServer{
		DeliverHandler:        deliver.NewHandler(&peer.DeliverChainManager{}, coreConfig.AuthenticationTimeWindow, mutualTLS, metrics),
		PolicyCheckerProvider: policyCheckerProvider,
	}
	pb.RegisterDeliverServer(peerServer.Server(), abServer)

	
	ca, err := tlsgen.NewCA()
	if err != nil {
		logger.Panic("Failed creating authentication layer:", err)
	}
	ccSrv, ccEndpoint, err := createChaincodeServer(coreConfig, ca, peerHost)
	if err != nil {
		logger.Panicf("Failed to create chaincode server: %s", err)
	}

	
	userRunsCC := chaincode.IsDevMode()
	tlsEnabled := coreConfig.PeerTLSEnabled

	
	authenticator := accesscontrol.NewAuthenticator(ca)
	ipRegistry := inproccontroller.NewRegistry()

	sccp := &scc.Provider{
		Peer:        peer.Default,
		PeerSupport: peer.DefaultSupport,
		Registrar:   ipRegistry,
		Whitelist:   scc.GlobalWhitelist(),
	}
	lsccInst := lscc.New(sccp, aclProvider, platformRegistry)

	chaincodeEndorsementInfo := &lifecycle.ChaincodeEndorsementInfo{
		LegacyImpl: lsccInst,
		Resources:  lifecycleResources,
		Cache:      lifecycleCache,
	}

	lifecycleFunctions := &lifecycle.ExternalFunctions{
		Resources:       lifecycleResources,
		InstallListener: lifecycleCache,
	}

	lifecycleSCC := &lifecycle.SCC{
		Dispatcher: &dispatcher.Dispatcher{
			Protobuf: &dispatcher.ProtobufImpl{},
		},
		Functions:           lifecycleFunctions,
		OrgMSPID:            mspID,
		ChannelConfigSource: peer.Default,
		ACLProvider:         aclProvider,
	}

	var client *docker.Client
	if coreConfig.VMDockerTLSEnabled {
		client, err = docker.NewTLSClient(coreConfig.VMEndpoint, coreConfig.DockerCert, coreConfig.DockerKey, coreConfig.DockerCA)
	} else {
		client, err = docker.NewClient(coreConfig.VMEndpoint)
	}
	if err != nil {
		logger.Panicf("cannot create docker client: %s", err)
	}

	dockerProvider := &dockercontroller.Provider{
		PeerID:        coreConfig.PeerID,
		NetworkID:     coreConfig.NetworkID,
		BuildMetrics:  dockercontroller.NewBuildMetrics(opsSystem.Provider),
		Client:        client,
		AttachStdOut:  coreConfig.VMDockerAttachStdout,
		HostConfig:    getDockerHostConfig(),
		ChaincodePull: coreConfig.ChaincodePull,
	}
	dockerVM := dockerProvider.NewVM()

	err = opsSystem.RegisterChecker("docker", dockerVM)
	if err != nil {
		logger.Panicf("failed to register docker health check: %s", err)
	}

	chaincodeConfig := chaincode.GlobalConfig()

	chaincodeVMController := container.NewVMController(
		map[string]container.VMProvider{
			dockercontroller.ContainerType: dockerProvider,
			inproccontroller.ContainerType: ipRegistry,
		},
	)

	containerRuntime := &chaincode.ContainerRuntime{
		CACert:        ca.CertBytes(),
		CertGenerator: authenticator,
		CommonEnv: []string{
			"CORE_CHAINCODE_LOGGING_LEVEL=" + chaincodeConfig.LogLevel,
			"CORE_CHAINCODE_LOGGING_SHIM=" + chaincodeConfig.ShimLogLevel,
			"CORE_CHAINCODE_LOGGING_FORMAT=" + chaincodeConfig.LogFormat,
		},
		DockerClient:     client,
		PeerAddress:      ccEndpoint,
		Processor:        chaincodeVMController,
		PlatformRegistry: platformRegistry,
	}

	
	if !chaincodeConfig.TLSEnabled {
		containerRuntime.CertGenerator = nil
	}

	globalConfig := chaincode.GlobalConfig()
	chaincodeHandlerRegistry := chaincode.NewHandlerRegistry(userRunsCC)
	chaincodeLauncher := &chaincode.RuntimeLauncher{
		Metrics:         chaincode.NewLaunchMetrics(opsSystem.Provider),
		PackageProvider: packageProvider,
		Registry:        chaincodeHandlerRegistry,
		Runtime:         containerRuntime,
		StartupTimeout:  globalConfig.StartupTimeout,
	}

	chaincodeSupport := &chaincode.ChaincodeSupport{
		ACLProvider:            aclProvider,
		AppConfig:              peer.DefaultSupport,
		DeployedCCInfoProvider: lifecycleValidatorCommitter,
		ExecuteTimeout:         globalConfig.ExecuteTimeout,
		HandlerRegistry:        chaincodeHandlerRegistry,
		HandlerMetrics:         chaincode.NewHandlerMetrics(opsSystem.Provider),
		Keepalive:              globalConfig.Keepalive,
		Launcher:               chaincodeLauncher,
		Lifecycle:              chaincodeEndorsementInfo,
		Runtime:                containerRuntime,
		SystemCCProvider:       sccp,
		TotalQueryLimit:        globalConfig.TotalQueryLimit,
		UserRunsCC:             userRunsCC,
	}

	ipRegistry.ChaincodeSupport = chaincodeSupport

	ccSupSrv := pb.ChaincodeSupportServer(chaincodeSupport)
	if tlsEnabled {
		ccSupSrv = authenticator.Wrap(ccSupSrv)
	}

	csccInst := cscc.New(sccp, aclProvider, lifecycleValidatorCommitter, lsccInst, lifecycleValidatorCommitter)
	qsccInst := scc.SelfDescribingSysCC(qscc.New(aclProvider))
	if maxConcurrency := coreConfig.LimitsConcurrencyQSCC; maxConcurrency != 0 {
		qsccInst = scc.Throttle(maxConcurrency, qsccInst)
	}

	
	sccs := scc.CreatePluginSysCCs(sccp)
	for _, cc := range append([]scc.SelfDescribingSysCC{lsccInst, csccInst, qsccInst, lifecycleSCC}, sccs...) {
		sccp.RegisterSysCC(cc)
	}
	pb.RegisterChaincodeSupportServer(ccSrv.Server(), ccSupSrv)

	
	go ccSrv.Start()

	logger.Debugf("Running peer")

	
	startAdminServer(coreConfig, listenAddr, peerServer.Server(), metricsProvider)

	privDataDist := func(channel string, txID string, privateData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error {
		return service.GetGossipService().DistributePrivateData(channel, txID, privateData, blkHt)
	}

	signingIdentity := mgmt.GetLocalSigningIdentityOrPanic()

	libConf := library.Config{}
	if err = viperutil.EnhancedExactUnmarshalKey("peer.handlers", &libConf); err != nil {
		return errors.WithMessage(err, "could not load YAML config")
	}
	reg := library.InitRegistry(libConf)

	authFilters := reg.Lookup(library.Auth).([]authHandler.Filter)
	endorserSupport := &endorser.SupportImpl{
		SignerSerializer: signingIdentity,
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
	serverEndorser := endorser.NewEndorserServer(privDataDist, endorserSupport, platformRegistry, metricsProvider)
	auth := authHandler.ChainFilters(serverEndorser, authFilters...)
	
	pb.RegisterEndorserServer(peerServer.Server(), auth)

	policyMgr := peer.NewChannelPolicyManagerGetter()

	
	err = initGossipService(policyMgr, metricsProvider, peerServer, signingIdentity, coreConfig.PeerAddress)
	if err != nil {
		return err
	}
	defer service.GetGossipService().Stop()

	
	err = registerProverService(peerServer, aclProvider, signingIdentity)
	if err != nil {
		return err
	}

	

	
	ccp := chaincode.NewProvider(chaincodeSupport)
	sccp.DeploySysCCs("", ccp)
	logger.Infof("Deployed system chaincodes")

	installedCCs := func() ([]ccdef.InstalledChaincode, error) {
		return packageProvider.ListInstalledChaincodes()
	}
	lifecycle, err := cclifecycle.NewLifecycle(cclifecycle.EnumerateFunc(installedCCs))
	if err != nil {
		logger.Panicf("Failed creating lifecycle: +%v", err)
	}
	onUpdate := cclifecycle.HandleMetadataUpdateFunc(func(channel string, chaincodes ccdef.MetadataSet) {
		service.GetGossipService().UpdateChaincodes(chaincodes.AsChaincodes(), gossipcommon.ChainID(channel))
	})
	lifecycle.AddListener(onUpdate)

	
	peer.Initialize(
		func(cid string) {
			logger.Debugf("Deploying system CC, for channel <%s>", cid)
			sccp.DeploySysCCs(cid, ccp)
			sub, err := lifecycle.NewChannelSubscription(cid, cclifecycle.QueryCreatorFunc(func() (cclifecycle.Query, error) {
				return peer.GetLedger(cid).NewQueryExecutor()
			}))
			if err != nil {
				logger.Panicf("Failed subscribing to chaincode lifecycle updates")
			}
			cceventmgmt.GetMgr().Register(cid, sub)
		},
		sccp,
		plugin.MapBasedMapper(validationPluginsByName),
		platformRegistry,
		lifecycleValidatorCommitter,
		membershipInfoProvider,
		metricsProvider,
		lsccInst,
		lifecycleValidatorCommitter,
		ledgerConfig(),
		coreConfig.ValidatorPoolSize,
	)

	if coreConfig.DiscoveryEnabled {
		registerDiscoveryService(coreConfig, peerServer, policyMgr, lifecycle)
	}

	logger.Infof("Starting peer with ID=[%s], network ID=[%s], address=[%s]", coreConfig.PeerID, coreConfig.NetworkID, coreConfig.PeerAddress)

	
	
	profileEnabled := coreConfig.ProfileEnabled
	profileListenAddress := coreConfig.ProfileListenAddress

	
	
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

	logger.Infof("Started peer with ID=[%s], network ID=[%s], address=[%s]", coreConfig.PeerID, coreConfig.NetworkID, coreConfig.PeerAddress)

	
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
	policy, _, err := pp.NewPolicy(protoutil.MarshalOrPanic(policyObject))
	if err != nil {
		logger.Panicf("Failed creating local policy: +%v", err)
	}
	return policy
}

func createSelfSignedData() protoutil.SignedData {
	sID := mgmt.GetLocalSigningIdentityOrPanic()
	msg := make([]byte, 32)
	sig, err := sID.Sign(msg)
	if err != nil {
		logger.Panicf("Failed creating self signed data because message signing failed: %v", err)
	}
	peerIdentity, err := sID.Serialize()
	if err != nil {
		logger.Panicf("Failed creating self signed data because peer identity couldn't be serialized: %v", err)
	}
	return protoutil.SignedData{
		Data:      msg,
		Signature: sig,
		Identity:  peerIdentity,
	}
}

func registerDiscoveryService(coreConfig *peer.Config, peerServer *comm.GRPCServer, polMgr policies.ChannelPolicyManagerGetter, lc *cclifecycle.Lifecycle) {
	mspID := coreConfig.LocalMSPID
	localAccessPolicy := localPolicy(cauthdsl.SignedByAnyAdmin([]string{mspID}))
	if coreConfig.DiscoveryOrgMembersAllowed {
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
		AuthCacheEnabled:             coreConfig.DiscoveryAuthCacheEnabled,
		AuthCacheMaxSize:             coreConfig.DiscoveryAuthCacheMaxSize,
		AuthCachePurgeRetentionRatio: coreConfig.DiscoveryAuthCachePurgeRetentionRatio,
	}, support)
	logger.Info("Discovery service activated")
	discprotos.RegisterDiscoveryServer(peerServer.Server(), svc)
}


func createChaincodeServer(coreConfig *peer.Config, ca tlsgen.CA, peerHostname string) (srv *comm.GRPCServer, ccEndpoint string, err error) {
	
	ccEndpoint, err = computeChaincodeEndpoint(coreConfig.ChaincodeAddress, coreConfig.ChaincodeListenAddress, peerHostname)
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

	cclistenAddress := coreConfig.ChaincodeListenAddress
	if cclistenAddress == "" {
		cclistenAddress = fmt.Sprintf("%s:%d", peerHostname, defaultChaincodePort)
		logger.Warningf("%s is not set, using %s", chaincodeListenAddrKey, cclistenAddress)
		coreConfig.ChaincodeListenAddress = cclistenAddress
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
	config.HealthCheckEnabled = true

	srv, err = comm.NewGRPCServer(cclistenAddress, config)
	if err != nil {
		logger.Errorf("Error creating GRPC server: %s", err)
		return nil, "", err
	}

	return srv, ccEndpoint, nil
}








func computeChaincodeEndpoint(chaincodeAddress string, chaincodeListenAddress string, peerHostname string) (ccEndpoint string, err error) {
	logger.Infof("Entering computeChaincodeEndpoint with peerHostname: %s", peerHostname)
	
	if chaincodeAddress != "" {
		host, _, err := net.SplitHostPort(chaincodeAddress)
		if err != nil {
			logger.Errorf("Fail to split chaincodeAddress: %s", err)
			return "", err
		}
		ccIP := net.ParseIP(host)
		if ccIP != nil && ccIP.IsUnspecified() {
			logger.Errorf("ChaincodeAddress' IP cannot be %s in non-dev mode", ccIP)
			return "", errors.New("invalid endpoint for chaincode to connect")
		}
		logger.Infof("Exit with ccEndpoint: %s", chaincodeAddress)
		return chaincodeAddress, nil
	}

	
	if chaincodeListenAddress != "" {
		ccEndpoint = chaincodeListenAddress
		host, port, err := net.SplitHostPort(ccEndpoint)
		if err != nil {
			logger.Errorf("ChaincodeAddress is nil and fail to split chaincodeListenAddress: %s", err)
			return "", err
		}

		ccListenerIP := net.ParseIP(host)
		
		
		if ccListenerIP != nil && ccListenerIP.IsUnspecified() {
			
			peerIP := net.ParseIP(peerHostname)
			if peerIP != nil && peerIP.IsUnspecified() {
				
				logger.Error("ChaincodeAddress is nil while both chaincodeListenAddressIP and peerIP are 0.0.0.0")
				return "", errors.New("invalid endpoint for chaincode to connect")
			}
			ccEndpoint = fmt.Sprintf("%s:%s", peerHostname, port)
		}
		logger.Infof("Exit with ccEndpoint: %s", ccEndpoint)
		return ccEndpoint, nil
	}

	
	peerIP := net.ParseIP(peerHostname)
	if peerIP != nil && peerIP.IsUnspecified() {
		
		logger.Errorf("ChaincodeAddress and chaincodeListenAddress are nil and peerIP is %s", peerIP)
		return "", errors.New("invalid endpoint for chaincode to connect")
	}

	
	ccEndpoint = fmt.Sprintf("%s:%d", peerHostname, defaultChaincodePort)

	logger.Infof("Exit with ccEndpoint: %s", ccEndpoint)
	return ccEndpoint, nil
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

func startAdminServer(coreConfig *peer.Config, peerListenAddr string, peerServer *grpc.Server, metricsProvider metrics.Provider) {
	adminListenAddress := coreConfig.AdminListenAddress
	separateLsnrForAdmin := adminHasSeparateListener(peerListenAddr, adminListenAddress)
	mspID := coreConfig.LocalMSPID
	adminPolicy := localPolicy(cauthdsl.SignedByAnyAdmin([]string{mspID}))
	gRPCService := peerServer
	if separateLsnrForAdmin {
		logger.Info("Creating gRPC server for admin service on", adminListenAddress)
		serverConfig, err := peer.GetServerConfig()
		if err != nil {
			logger.Fatalf("Error loading secure config for admin service (%s)", err)
		}
		serverConfig.Logger = flogging.MustGetLogger("core.comm").With("server", "AdminServer")
		serverConfig.MetricsProvider = metricsProvider
		serverConfig.UnaryInterceptors = append(
			serverConfig.UnaryInterceptors,
			grpcmetrics.UnaryServerInterceptor(grpcmetrics.NewUnaryMetrics(metricsProvider)),
			grpclogging.UnaryServerInterceptor(flogging.MustGetLogger("comm.grpc.server").Zap()),
		)
		serverConfig.StreamInterceptors = append(
			serverConfig.StreamInterceptors,
			grpcmetrics.StreamServerInterceptor(grpcmetrics.NewStreamMetrics(metricsProvider)),
			grpclogging.StreamServerInterceptor(flogging.MustGetLogger("comm.grpc.server").Zap()),
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






func initGossipService(
	policyMgr policies.ChannelPolicyManagerGetter,
	metricsProvider metrics.Provider,
	peerServer *comm.GRPCServer,
	signer msp.SigningIdentity,
	peerAddr string,
) error {

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
		signer,
		mgmt.NewDeserializersManager(),
	)
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	bootstrap := viper.GetStringSlice("peer.gossip.bootstrap")

	return service.InitGossipService(
		signer,
		metricsProvider,
		peerAddr,
		peerServer.Server(),
		certs,
		messageCryptoService,
		secAdv,
		secureDialOpts,
		bootstrap...,
	)
}

func newOperationsSystem(coreConfig *peer.Config) *operations.System {
	return operations.NewSystem(operations.Options{
		Logger:        flogging.MustGetLogger("peer.operations"),
		ListenAddress: coreConfig.OperationsListenAddress,
		Metrics: operations.MetricsOptions{
			Provider: coreConfig.MetricsProvider,
			Statsd: &operations.Statsd{
				Network:       coreConfig.StatsdNetwork,
				Address:       coreConfig.StatsdAaddress,
				WriteInterval: coreConfig.StatsdWriteInterval,
				Prefix:        coreConfig.StatsdPrefix,
			},
		},
		TLS: operations.TLS{
			Enabled:            coreConfig.OperationsTLSEnabled,
			CertFile:           coreConfig.OperationsTLSCertFile,
			KeyFile:            coreConfig.OperationsTLSKeyFile,
			ClientCertRequired: coreConfig.OperationsTLSClientAuthRequired,
			ClientCACertFiles:  coreConfig.OperationsTLSClientRootCAs,
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
			TokenOwnerValidatorManager: &server.PeerTokenOwnerValidatorManager{
				IdentityDeserializerManager: &manager.FabricIdentityDeserializerManager{},
			},
		},
	}
	token.RegisterProverServer(peerServer.Server(), prover)
	return nil
}

func getDockerHostConfig() *docker.HostConfig {
	dockerKey := func(key string) string { return "vm.docker.hostConfig." + key }
	getInt64 := func(key string) int64 { return int64(viper.GetInt(dockerKey(key))) }

	var logConfig docker.LogConfig
	err := viper.UnmarshalKey(dockerKey("LogConfig"), &logConfig)
	if err != nil {
		logger.Panicf("unable to parse Docker LogConfig: %s", err)
	}

	networkMode := viper.GetString(dockerKey("NetworkMode"))
	if networkMode == "" {
		networkMode = "host"
	}

	return &docker.HostConfig{
		CapAdd:  viper.GetStringSlice(dockerKey("CapAdd")),
		CapDrop: viper.GetStringSlice(dockerKey("CapDrop")),

		DNS:         viper.GetStringSlice(dockerKey("Dns")),
		DNSSearch:   viper.GetStringSlice(dockerKey("DnsSearch")),
		ExtraHosts:  viper.GetStringSlice(dockerKey("ExtraHosts")),
		NetworkMode: networkMode,
		IpcMode:     viper.GetString(dockerKey("IpcMode")),
		PidMode:     viper.GetString(dockerKey("PidMode")),
		UTSMode:     viper.GetString(dockerKey("UTSMode")),
		LogConfig:   logConfig,

		ReadonlyRootfs:   viper.GetBool(dockerKey("ReadonlyRootfs")),
		SecurityOpt:      viper.GetStringSlice(dockerKey("SecurityOpt")),
		CgroupParent:     viper.GetString(dockerKey("CgroupParent")),
		Memory:           getInt64("Memory"),
		MemorySwap:       getInt64("MemorySwap"),
		MemorySwappiness: getInt64("MemorySwappiness"),
		OOMKillDisable:   viper.GetBool(dockerKey("OomKillDisable")),
		CPUShares:        getInt64("CpuShares"),
		CPUSet:           viper.GetString(dockerKey("Cpuset")),
		CPUSetCPUs:       viper.GetString(dockerKey("CpusetCPUs")),
		CPUSetMEMs:       viper.GetString(dockerKey("CpusetMEMs")),
		CPUQuota:         getInt64("CpuQuota"),
		CPUPeriod:        getInt64("CpuPeriod"),
		BlkioWeight:      getInt64("BlkioWeight"),
	}
}
