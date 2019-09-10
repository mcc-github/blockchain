/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	discprotos "github.com/mcc-github/blockchain-protos-go/discovery"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/bccsp/factory"
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
	"github.com/mcc-github/blockchain/core/aclmgmt"
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
	"github.com/mcc-github/blockchain/core/container/externalbuilders"
	"github.com/mcc-github/blockchain/core/deliverservice"
	"github.com/mcc-github/blockchain/core/dispatcher"
	"github.com/mcc-github/blockchain/core/endorser"
	authHandler "github.com/mcc-github/blockchain/core/handlers/auth"
	endorsement2 "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	endorsement3 "github.com/mcc-github/blockchain/core/handlers/endorsement/api/identities"
	"github.com/mcc-github/blockchain/core/handlers/library"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/core/ledger/kvledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/operations"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/scc/cscc"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	"github.com/mcc-github/blockchain/core/scc/qscc"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/discovery"
	"github.com/mcc-github/blockchain/discovery/endorsement"
	discsupport "github.com/mcc-github/blockchain/discovery/support"
	discacl "github.com/mcc-github/blockchain/discovery/support/acl"
	ccsupport "github.com/mcc-github/blockchain/discovery/support/chaincode"
	"github.com/mcc-github/blockchain/discovery/support/config"
	"github.com/mcc-github/blockchain/discovery/support/gossip"
	gossipcommon "github.com/mcc-github/blockchain/gossip/common"
	gossipgossip "github.com/mcc-github/blockchain/gossip/gossip"
	gossipmetrics "github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/service"
	gossipservice "github.com/mcc-github/blockchain/gossip/service"
	peergossip "github.com/mcc-github/blockchain/internal/peer/gossip"
	"github.com/mcc-github/blockchain/internal/peer/version"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protoutil"
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
	flags.BoolVarP(&chaincodeDevMode, "peer-chaincodedev", "", false, "start peer in chaincode development mode")
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



type externalVMAdapter struct {
	detector *externalbuilders.Detector
}

func (e externalVMAdapter) Build(
	ccid string,
	metadata *persistence.ChaincodePackageMetadata,
	codePackage io.Reader,
) (container.Instance, error) {
	return e.detector.Build(ccid, metadata, codePackage)
}

type endorserChannelAdapter struct {
	peer *peer.Peer
}

func (e endorserChannelAdapter) Channel(channelID string) *endorser.Channel {
	if peerChannel := e.peer.Channel(channelID); peerChannel != nil {
		return &endorser.Channel{
			IdentityDeserializer: peerChannel.MSPManager(),
		}
	}

	return nil
}

func serve(args []string) error {
	
	
	
	
	
	
	mspType := mgmt.GetLocalMSP().GetType()
	if mspType != msp.FABRIC {
		panic("Unsupported msp type " + msp.ProviderTypeToString(mspType))
	}

	
	
	
	grpc.EnableTracing = true

	logger.Infof("Starting %s", version.GetInfo())

	
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
		return errors.WithMessage(err, "failed to initialize operations subsystems")
	}
	defer opsSystem.Stop()

	metricsProvider := opsSystem.Provider
	logObserver := floggingmetrics.NewObserver(metricsProvider)
	flogging.SetObserver(logObserver)

	membershipInfoProvider := privdata.NewMembershipInfoProvider(createSelfSignedData(), identityDeserializerFactory)

	mspID := coreConfig.LocalMSPID

	chaincodeInstallPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "lifecycle", "chaincodes")
	ccStore := persistence.NewStore(chaincodeInstallPath)
	ccPackageParser := &persistence.ChaincodePackageParser{
		MetadataProvider: ccprovider.PersistenceAdapter(ccprovider.MetadataAsTarEntries),
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

	cs := comm.NewCredentialSupport()
	if serverConfig.SecOpts.UseTLS {
		logger.Info("Starting peer with TLS enabled")
		cs = comm.NewCredentialSupport(serverConfig.SecOpts.ServerRootCAs...)

		
		clientCert, err := peer.GetClientCertificate()
		if err != nil {
			logger.Fatalf("Failed to set TLS client certificate (%s)", err)
		}
		cs.SetClientCertificate(clientCert)
	}

	peerServer, err := comm.NewGRPCServer(listenAddr, serverConfig)
	if err != nil {
		logger.Fatalf("Failed to create peer server (%s)", err)
	}

	transientStoreProvider, err := transientstore.NewStoreProvider(
		filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "transientstore"),
	)
	if err != nil {
		return errors.WithMessage(err, "failed to open transient store")
	}
	peerInstance := &peer.Peer{
		Server:            peerServer,
		ServerConfig:      serverConfig,
		CredentialSupport: cs,
		StoreProvider:     transientStoreProvider,
		CryptoProvider:    factory.GetDefault(),
	}

	localMSP := mgmt.GetLocalMSP()
	signingIdentity, err := localMSP.GetDefaultSigningIdentity()
	if err != nil {
		logger.Panicf("Could not get the default signing identity from the local MSP: [%+v]", err)
	}
	policyMgr := policies.PolicyManagerGetterFunc(peerInstance.GetPolicyManager)

	
	
	deliverClientDialOpts := deliverClientDialOpts(coreConfig)
	deliverServiceConfig := deliverservice.GlobalConfig()

	gossipService, err := initGossipService(
		policyMgr,
		metricsProvider,
		peerServer,
		signingIdentity,
		cs,
		coreConfig.PeerAddress,
		deliverClientDialOpts,
		deliverServiceConfig,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to initialize gossip service")
	}
	defer gossipService.Stop()

	peerInstance.GossipService = gossipService

	policyChecker := policy.NewPolicyChecker(
		policies.PolicyManagerGetterFunc(peerInstance.GetPolicyManager),
		mgmt.GetLocalMSP(),
		mgmt.NewLocalMSPPrincipalGetter(),
	)

	
	
	aclProvider := aclmgmt.NewACLProvider(
		aclmgmt.ResourceGetter(peerInstance.GetStableChannelConfig),
		policyChecker,
	)

	
	
	
	
	
	
	
	
	lifecycleResources := &lifecycle.Resources{
		Serializer:          &lifecycle.Serializer{},
		ChannelConfigSource: peerInstance,
		ChaincodeStore:      ccStore,
		PackageParser:       ccPackageParser,
	}

	lifecycleValidatorCommitter := &lifecycle.ValidatorCommitter{
		Resources:                    lifecycleResources,
		LegacyDeployedCCInfoProvider: &lscc.DeployedCCInfoProvider{},
	}

	packageProvider := &persistence.PackageProvider{
		LegacyPP: &ccprovider.CCInfoFSImpl{GetHasher: factory.GetDefault()},
	}

	
	
	legacyMetadataManager, err := cclifecycle.NewMetadataManager(
		cclifecycle.EnumerateFunc(
			func() ([]ccdef.InstalledChaincode, error) {
				return packageProvider.ListInstalledChaincodesLegacy()
			},
		),
	)
	if err != nil {
		logger.Panicf("Failed creating LegacyMetadataManager: +%v", err)
	}

	
	
	metadataManager := lifecycle.NewMetadataManager()

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	lifecycleCache := lifecycle.NewCache(lifecycleResources, mspID, metadataManager)

	txProcessors := map[common.HeaderType]ledger.CustomTxProcessor{
		common.HeaderType_CONFIG: &peer.ConfigTxProcessor{},
	}

	peerInstance.LedgerMgr = ledgermgmt.NewLedgerMgr(
		&ledgermgmt.Initializer{
			CustomTxProcessors:              txProcessors,
			DeployedChaincodeInfoProvider:   lifecycleValidatorCommitter,
			MembershipInfoProvider:          membershipInfoProvider,
			ChaincodeLifecycleEventProvider: lifecycleCache,
			MetricsProvider:                 metricsProvider,
			HealthCheckRegistry:             opsSystem,
			StateListeners:                  []ledger.StateListener{lifecycleCache},
			Config:                          ledgerConfig(),
		},
	)

	
	lsccInstallPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "chaincodes")
	ccprovider.SetChaincodesPath(lsccInstallPath)

	if err := lifecycleCache.InitializeLocalChaincodes(); err != nil {
		return errors.WithMessage(err, "could not initialize local chaincodes")
	}

	
	
	if chaincodeDevMode {
		logger.Info("Running in chaincode development mode")
		logger.Info("Disable loading validity system chaincode")

		viper.Set("chaincode.mode", chaincode.DevModeUserRunsChaincode)
	}

	mutualTLS := serverConfig.SecOpts.UseTLS && serverConfig.SecOpts.RequireClientCert
	policyCheckerProvider := func(resourceName string) deliver.PolicyCheckerFunc {
		return func(env *cb.Envelope, channelID string) error {
			return aclProvider.CheckACL(resourceName, channelID, env)
		}
	}

	metrics := deliver.NewMetrics(metricsProvider)
	abServer := &peer.DeliverServer{
		DeliverHandler: deliver.NewHandler(
			&peer.DeliverChainManager{Peer: peerInstance},
			coreConfig.AuthenticationTimeWindow,
			mutualTLS,
			metrics,
			false,
		),
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

	builtinSCCs := map[string]struct{}{
		"lscc":       {},
		"qscc":       {},
		"cscc":       {},
		"_lifecycle": {},
	}

	lsccInst := lscc.New(builtinSCCs, &lscc.PeerShim{Peer: peerInstance}, aclProvider, peerInstance.GetMSPIDs, policyChecker)

	chaincodeHandlerRegistry := chaincode.NewHandlerRegistry(userRunsCC)
	lifecycleTxQueryExecutorGetter := &chaincode.TxQueryExecutorGetter{
		CCID:            scc.ChaincodeID(lifecycle.LifecycleNamespace),
		HandlerRegistry: chaincodeHandlerRegistry,
	}
	chaincodeEndorsementInfo := &lifecycle.ChaincodeEndorsementInfoSource{
		LegacyImpl:  lsccInst,
		Resources:   lifecycleResources,
		Cache:       lifecycleCache,
		BuiltinSCCs: builtinSCCs,
	}

	lifecycleFunctions := &lifecycle.ExternalFunctions{
		Resources:                 lifecycleResources,
		InstallListener:           lifecycleCache,
		InstalledChaincodesLister: lifecycleCache,
	}

	lifecycleSCC := &lifecycle.SCC{
		Dispatcher: &dispatcher.Dispatcher{
			Protobuf: &dispatcher.ProtobufImpl{},
		},
		DeployedCCInfoProvider: lifecycleValidatorCommitter,
		QueryExecutorProvider:  lifecycleTxQueryExecutorGetter,
		Functions:              lifecycleFunctions,
		OrgMSPID:               mspID,
		ChannelConfigSource:    peerInstance,
		ACLProvider:            aclProvider,
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

	chaincodeConfig := chaincode.GlobalConfig()

	dockerVM := &dockercontroller.DockerVM{
		PeerID:        coreConfig.PeerID,
		NetworkID:     coreConfig.NetworkID,
		BuildMetrics:  dockercontroller.NewBuildMetrics(opsSystem.Provider),
		Client:        client,
		AttachStdOut:  coreConfig.VMDockerAttachStdout,
		HostConfig:    getDockerHostConfig(),
		ChaincodePull: coreConfig.ChaincodePull,
		NetworkMode:   coreConfig.VMNetworkMode,
		PlatformBuilder: &platforms.Builder{
			Registry: platformRegistry,
			Client:   client,
		},
		
		
		
		LoggingEnv: []string{
			"CORE_CHAINCODE_LOGGING_LEVEL=" + chaincodeConfig.LogLevel,
			"CORE_CHAINCODE_LOGGING_SHIM=" + chaincodeConfig.ShimLogLevel,
			"CORE_CHAINCODE_LOGGING_FORMAT=" + chaincodeConfig.LogFormat,
		},
	}
	if err := opsSystem.RegisterChecker("docker", dockerVM); err != nil {
		logger.Panicf("failed to register docker health check: %s", err)
	}

	externalVM := &externalbuilders.Detector{
		Builders: coreConfig.ExternalBuilders,
	}

	containerRuntime := &chaincode.ContainerRuntime{
		CACert:        ca.CertBytes(),
		CertGenerator: authenticator,
		PeerAddress:   ccEndpoint,
		ContainerRouter: &container.Router{
			DockerVM:   dockerVM,
			ExternalVM: externalVMAdapter{externalVM},
			PackageProvider: &persistence.FallbackPackageLocator{
				ChaincodePackageLocator: &persistence.ChaincodePackageLocator{
					ChaincodeDir: chaincodeInstallPath,
				},
				LegacyCCPackageLocator: &ccprovider.CCInfoFSImpl{GetHasher: factory.GetDefault()},
			},
		},
	}

	
	if !chaincodeConfig.TLSEnabled {
		containerRuntime.CertGenerator = nil
	}

	chaincodeLauncher := &chaincode.RuntimeLauncher{
		Metrics:        chaincode.NewLaunchMetrics(opsSystem.Provider),
		Registry:       chaincodeHandlerRegistry,
		Runtime:        containerRuntime,
		StartupTimeout: chaincodeConfig.StartupTimeout,
	}

	chaincodeSupport := &chaincode.ChaincodeSupport{
		ACLProvider:            aclProvider,
		AppConfig:              peerInstance,
		DeployedCCInfoProvider: lifecycleValidatorCommitter,
		ExecuteTimeout:         chaincodeConfig.ExecuteTimeout,
		HandlerRegistry:        chaincodeHandlerRegistry,
		HandlerMetrics:         chaincode.NewHandlerMetrics(opsSystem.Provider),
		Keepalive:              chaincodeConfig.Keepalive,
		Launcher:               chaincodeLauncher,
		Lifecycle:              chaincodeEndorsementInfo,
		Peer:                   peerInstance,
		Runtime:                containerRuntime,
		BuiltinSCCs:            builtinSCCs,
		TotalQueryLimit:        chaincodeConfig.TotalQueryLimit,
		UserRunsCC:             userRunsCC,
	}

	ccSupSrv := pb.ChaincodeSupportServer(chaincodeSupport)
	if tlsEnabled {
		ccSupSrv = authenticator.Wrap(ccSupSrv)
	}

	csccInst := cscc.New(
		aclProvider,
		lifecycleValidatorCommitter,
		lsccInst,
		lifecycleValidatorCommitter,
		policyChecker,
		peerInstance,
	)
	qsccInst := scc.SelfDescribingSysCC(qscc.New(aclProvider, peerInstance))
	if maxConcurrency := coreConfig.LimitsConcurrencyQSCC; maxConcurrency != 0 {
		qsccInst = scc.Throttle(maxConcurrency, qsccInst)
	}

	pb.RegisterChaincodeSupportServer(ccSrv.Server(), ccSupSrv)

	
	go ccSrv.Start()

	logger.Debugf("Running peer")

	libConf, err := library.LoadConfig()
	if err != nil {
		return errors.WithMessage(err, "could not decode peer handlers configuration")
	}

	reg := library.InitRegistry(libConf)

	authFilters := reg.Lookup(library.Auth).([]authHandler.Filter)
	endorserSupport := &endorser.SupportImpl{
		SignerSerializer: signingIdentity,
		Peer:             peerInstance,
		ChaincodeSupport: chaincodeSupport,
		ACLProvider:      aclProvider,
		BuiltinSCCs:      builtinSCCs,
	}
	endorsementPluginsByName := reg.Lookup(library.Endorsement).(map[string]endorsement2.PluginFactory)
	validationPluginsByName := reg.Lookup(library.Validation).(map[string]validation.PluginFactory)
	signingIdentityFetcher := (endorsement3.SigningIdentityFetcher)(endorserSupport)
	channelStateRetriever := endorser.ChannelStateRetriever(endorserSupport)
	pluginMapper := endorser.MapBasedPluginMapper(endorsementPluginsByName)
	pluginEndorser := endorser.NewPluginEndorser(&endorser.PluginSupport{
		ChannelStateRetriever:   channelStateRetriever,
		TransientStoreRetriever: peerInstance,
		PluginMapper:            pluginMapper,
		SigningIdentityFetcher:  signingIdentityFetcher,
	})
	endorserSupport.PluginEndorser = pluginEndorser
	channelFetcher := endorserChannelAdapter{
		peer: peerInstance,
	}
	serverEndorser := &endorser.Endorser{
		PrivateDataDistributor: gossipService,
		ChannelFetcher:         channelFetcher,
		LocalMSP:               localMSP,
		Support:                endorserSupport,
		Metrics:                endorser.NewMetrics(metricsProvider),
	}

	
	for _, cc := range []scc.SelfDescribingSysCC{lsccInst, csccInst, qsccInst, lifecycleSCC} {
		if enabled, ok := chaincodeConfig.SCCWhitelist[cc.Name()]; !ok || !enabled {
			logger.Infof("not deploying chaincode %s as it is not enabled", cc.Name())
			continue
		}
		scc.DeploySysCC(cc, chaincodeSupport)
	}

	logger.Infof("Deployed system chaincodes")

	
	
	
	
	legacyMetadataManager.AddListener(metadataManager)

	
	metadataManager.AddListener(lifecycle.HandleMetadataUpdateFunc(func(channel string, chaincodes ccdef.MetadataSet) {
		gossipService.UpdateChaincodes(chaincodes.AsChaincodes(), gossipcommon.ChannelID(channel))
	}))

	
	peerInstance.Initialize(
		func(cid string) {
			
			
			
			lifecycleCache.InitializeMetadata(cid)

			
			
			
			
			
			
			sub, err := legacyMetadataManager.NewChannelSubscription(cid, cclifecycle.QueryCreatorFunc(func() (cclifecycle.Query, error) {
				return peerInstance.GetLedger(cid).NewQueryExecutor()
			}))
			if err != nil {
				logger.Panicf("Failed subscribing to chaincode lifecycle updates")
			}

			
			
			cceventmgmt.GetMgr().Register(cid, sub)
		},
		plugin.MapBasedMapper(validationPluginsByName),
		lifecycleValidatorCommitter,
		lsccInst,
		lifecycleValidatorCommitter,
		coreConfig.ValidatorPoolSize,
	)

	if coreConfig.DiscoveryEnabled {
		registerDiscoveryService(
			coreConfig,
			peerInstance,
			peerServer,
			policyMgr,
			lifecycle.NewMetadataProvider(
				lifecycleCache,
				legacyMetadataManager,
				peerInstance,
			),
			gossipService,
		)
	}

	logger.Infof("Starting peer with ID=[%s], network ID=[%s], address=[%s]", coreConfig.PeerID, coreConfig.NetworkID, coreConfig.PeerAddress)

	
	
	profileEnabled := coreConfig.ProfileEnabled
	profileListenAddress := coreConfig.ProfileListenAddress

	
	
	serve := make(chan error)

	
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

	
	rootFSPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "ledgersData")
	preResetHeights, err := kvledger.LoadPreResetHeight(rootFSPath)
	if err != nil {
		return fmt.Errorf("error loading prereset height: %s", err)
	}

	for cid, height := range preResetHeights {
		logger.Infof("Ledger rebuild: channel [%s]: preresetHeight: [%d]", cid, height)
	}

	if len(preResetHeights) > 0 {
		logger.Info("Ledger rebuild: Entering loop to check if current ledger heights surpass prereset ledger heights. Endorsement request processing will be disabled.")
		resetFilter := &reset{
			reject: true,
		}
		authFilters = append(authFilters, resetFilter)
		go resetLoop(resetFilter, preResetHeights, peerInstance.GetLedger, 10*time.Second)
	}

	
	auth := authHandler.ChainFilters(serverEndorser, authFilters...)
	
	pb.RegisterEndorserServer(peerServer.Server(), auth)

	go func() {
		var grpcErr error
		if grpcErr = peerServer.Start(); grpcErr != nil {
			grpcErr = fmt.Errorf("grpc server exited with error: %s", grpcErr)
		}
		serve <- grpcErr
	}()

	
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

func registerDiscoveryService(
	coreConfig *peer.Config,
	peerInstance *peer.Peer,
	peerServer *comm.GRPCServer,
	polMgr policies.ChannelPolicyManagerGetter,
	metadataProvider *lifecycle.MetadataProvider,
	gossipService *gossipservice.GossipService,
) {
	mspID := coreConfig.LocalMSPID
	localAccessPolicy := localPolicy(cauthdsl.SignedByAnyAdmin([]string{mspID}))
	if coreConfig.DiscoveryOrgMembersAllowed {
		localAccessPolicy = localPolicy(cauthdsl.SignedByAnyMember([]string{mspID}))
	}
	channelVerifier := discacl.NewChannelVerifier(policies.ChannelApplicationWriters, polMgr)
	acl := discacl.NewDiscoverySupport(channelVerifier, localAccessPolicy, discacl.ChannelConfigGetterFunc(peerInstance.GetStableChannelConfig))
	gSup := gossip.NewDiscoverySupport(gossipService)
	ccSup := ccsupport.NewDiscoverySupport(metadataProvider)
	ea := endorsement.NewEndorsementAnalyzer(gSup, ccSup, acl, metadataProvider)
	confSup := config.NewDiscoverySupport(config.CurrentConfigBlockGetterFunc(func(channelID string) *common.Block {
		channel := peerInstance.Channel(channelID)
		if channel == nil {
			return nil
		}
		block, err := peer.ConfigBlockFromLedger(channel.Ledger())
		if err != nil {
			logger.Error("failed to get config block", err)
			return nil
		}
		return block
	}))
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
		config.SecOpts = comm.SecureOptions{
			UseTLS: true,
			
			RequireClientCert: true,
			
			ClientRootCAs: [][]byte{ca.CertBytes()},
			
			Certificate: certKeyPair.Cert,
			Key:         certKeyPair.Key,
			
			
			ServerRootCAs: nil,
		}
	}

	
	chaincodeKeepaliveOptions := comm.KeepaliveOptions{
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


func secureDialOpts(credSupport *comm.CredentialSupport) func() []grpc.DialOption {
	return func() []grpc.DialOption {
		var dialOpts []grpc.DialOption
		
		dialOpts = append(
			dialOpts,
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(comm.MaxRecvMsgSize), grpc.MaxCallSendMsgSize(comm.MaxSendMsgSize)),
		)
		
		kaOpts := comm.DefaultKeepaliveOptions
		if viper.IsSet("peer.keepalive.client.interval") {
			kaOpts.ClientInterval = viper.GetDuration("peer.keepalive.client.interval")
		}
		if viper.IsSet("peer.keepalive.client.timeout") {
			kaOpts.ClientTimeout = viper.GetDuration("peer.keepalive.client.timeout")
		}
		dialOpts = append(dialOpts, comm.ClientKeepaliveOptions(kaOpts)...)

		if viper.GetBool("peer.tls.enabled") {
			dialOpts = append(dialOpts, grpc.WithTransportCredentials(credSupport.GetPeerCredentials()))
		} else {
			dialOpts = append(dialOpts, grpc.WithInsecure())
		}
		return dialOpts
	}
}


func deliverClientDialOpts(coreConfig *peer.Config) []grpc.DialOption {
	dialOpts := []grpc.DialOption{grpc.WithBlock()}
	
	dialOpts = append(
		dialOpts,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(comm.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(comm.MaxSendMsgSize)))
	
	kaOpts := coreConfig.DeliverClientKeepaliveOptions
	dialOpts = append(dialOpts, comm.ClientKeepaliveOptions(kaOpts)...)

	return dialOpts
}






func initGossipService(
	policyMgr policies.ChannelPolicyManagerGetter,
	metricsProvider metrics.Provider,
	peerServer *comm.GRPCServer,
	signer msp.SigningIdentity,
	credSupport *comm.CredentialSupport,
	peerAddress string,
	deliverClientDialOpts []grpc.DialOption,
	deliverServiceConfig *deliverservice.DeliverServiceConfig,
) (*gossipservice.GossipService, error) {

	var certs *gossipcommon.TLSCertificates
	if peerServer.TLSEnabled() {
		serverCert := peerServer.ServerCertificate()
		clientCert, err := peer.GetClientCertificate()
		if err != nil {
			return nil, errors.Wrap(err, "failed obtaining client certificates")
		}
		certs = &gossipcommon.TLSCertificates{}
		certs.TLSServerCert.Store(&serverCert)
		certs.TLSClientCert.Store(&clientCert)
	}

	messageCryptoService := peergossip.NewMCS(
		policyMgr,
		signer,
		mgmt.NewDeserializersManager(),
		factory.GetDefault(),
	)
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	bootstrap := viper.GetStringSlice("peer.gossip.bootstrap")

	serviceConfig := service.GlobalConfig()
	if serviceConfig.Endpoint != "" {
		peerAddress = serviceConfig.Endpoint
	}
	gossipConfig, err := gossipgossip.GlobalConfig(peerAddress, certs, bootstrap...)
	if err != nil {
		return nil, errors.Wrap(err, "failed obtaining gossip config")
	}

	return gossipservice.New(
		signer,
		gossipmetrics.NewGossipMetrics(metricsProvider),
		peerAddress,
		peerServer.Server(),
		messageCryptoService,
		secAdv,
		secureDialOpts(credSupport),
		credSupport,
		deliverClientDialOpts,
		gossipConfig,
		serviceConfig,
		deliverServiceConfig,
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

	memorySwappiness := getInt64("MemorySwappiness")
	oomKillDisable := viper.GetBool(dockerKey("OomKillDisable"))

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
		MemorySwappiness: &memorySwappiness,
		OOMKillDisable:   &oomKillDisable,
		CPUShares:        getInt64("CpuShares"),
		CPUSet:           viper.GetString(dockerKey("Cpuset")),
		CPUSetCPUs:       viper.GetString(dockerKey("CpusetCPUs")),
		CPUSetMEMs:       viper.GetString(dockerKey("CpusetMEMs")),
		CPUQuota:         getInt64("CpuQuota"),
		CPUPeriod:        getInt64("CpuPeriod"),
		BlkioWeight:      getInt64("BlkioWeight"),
	}
}




type peerLedger interface {
	ledger.PeerLedger
}

type getLedger func(string) ledger.PeerLedger

func resetLoop(
	resetFilter *reset,
	preResetHeights map[string]uint64,
	pLedger getLedger,
	interval time.Duration,
) {
	
	ticker := time.NewTicker(interval)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			logger.Info("Ledger rebuild: Checking if current ledger heights surpass prereset ledger heights")
			logger.Debugf("Ledger rebuild: Number of ledgers still rebuilding before check: %d", len(preResetHeights))
			for cid, height := range preResetHeights {
				var l peerLedger
				l = pLedger(cid)
				if l != nil {
					bcInfo, err := l.GetBlockchainInfo()
					if bcInfo != nil {
						logger.Debugf("Ledger rebuild: channel [%s]: currentHeight [%d] : preresetHeight [%d]", cid, bcInfo.GetHeight(), height)
						if bcInfo.GetHeight() >= height {
							delete(preResetHeights, cid)
						} else {
							break
						}
					} else {
						if err != nil {
							logger.Warningf("Ledger rebuild: could not retrieve info for channel [%s]: %s", cid, err.Error())
						}
					}
				}
			}

			logger.Debugf("Ledger rebuild: Number of ledgers still rebuilding after check: %d", len(preResetHeights))
			if len(preResetHeights) == 0 {
				logger.Infof("Ledger rebuild: Complete, all ledgers surpass prereset heights. Endorsement request processing will be enabled.")
				rootFSPath := filepath.Join(coreconfig.GetPath("peer.fileSystemPath"), "ledgersData")
				err := kvledger.ClearPreResetHeight(rootFSPath)
				if err != nil {
					logger.Warningf("Ledger rebuild: could not clear off prerest files: error=%s", err)
				}
				resetFilter.setReject(false)
				return
			}
		}
	}
}


type reset struct {
	sync.RWMutex
	next   pb.EndorserServer
	reject bool
}

func (r *reset) setReject(reject bool) {
	r.Lock()
	defer r.Unlock()
	r.reject = reject
}


func (r *reset) Init(next pb.EndorserServer) {
	r.next = next
}


func (r *reset) ProcessProposal(ctx context.Context, signedProp *pb.SignedProposal) (*pb.ProposalResponse, error) {
	r.RLock()
	defer r.RUnlock()
	if r.reject {
		return nil, errors.New("endorse requests are blocked while ledgers are being rebuilt")
	}
	return r.next.ProcessProposal(ctx, signedProp)
}
