


package server

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/flogging/floggingtest"
	"github.com/mcc-github/blockchain/common/localmsp"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/common/metrics/prometheus"
	"github.com/mcc-github/blockchain/common/tools/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/common/tools/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitializeLogging(t *testing.T) {
	origEnvValue := os.Getenv("FABRIC_LOGGING_SPEC")
	os.Setenv("FABRIC_LOGGING_SPEC", "foo=debug")
	initializeLogging()
	assert.Equal(t, "debug", flogging.Global.Level("foo").String())
	os.Setenv("FABRIC_LOGGING_SPEC", origEnvValue)
}

func TestInitializeProfilingService(t *testing.T) {
	origEnvValue := os.Getenv("FABRIC_LOGGING_SPEC")
	defer os.Setenv("FABRIC_LOGGING_SPEC", origEnvValue)
	os.Setenv("FABRIC_LOGGING_SPEC", "debug")
	
	listenAddr := func() string {
		l, _ := net.Listen("tcp", "localhost:0")
		l.Close()
		return l.Addr().String()
	}()
	initializeProfilingService(
		&localconfig.TopLevel{
			General: localconfig.General{
				Profile: localconfig.Profile{
					Enabled: true,
					Address: listenAddr,
				}},
			Kafka: localconfig.Kafka{Verbose: true},
		},
	)
	time.Sleep(500 * time.Millisecond)
	if _, err := http.Get("http://" + listenAddr + "/" + "/debug/"); err != nil {
		t.Logf("Expected pprof to be up (will retry again in 3 seconds): %s", err)
		time.Sleep(3 * time.Second)
		if _, err := http.Get("http://" + listenAddr + "/" + "/debug/"); err != nil {
			t.Fatalf("Expected pprof to be up: %s", err)
		}
	}
}

func TestInitializeServerConfig(t *testing.T) {
	conf := &localconfig.TopLevel{
		General: localconfig.General{
			TLS: localconfig.TLS{
				Enabled:            true,
				ClientAuthRequired: true,
				Certificate:        "main.go",
				PrivateKey:         "main.go",
				RootCAs:            []string{"main.go"},
				ClientRootCAs:      []string{"main.go"},
			},
		},
	}
	sc := initializeServerConfig(conf, nil)
	defaultOpts := comm.DefaultKeepaliveOptions
	assert.Equal(t, defaultOpts.ServerMinInterval, sc.KaOpts.ServerMinInterval)
	assert.Equal(t, time.Duration(0), sc.KaOpts.ServerInterval)
	assert.Equal(t, time.Duration(0), sc.KaOpts.ServerTimeout)
	testDuration := 10 * time.Second
	conf.General.Keepalive = localconfig.Keepalive{
		ServerMinInterval: testDuration,
		ServerInterval:    testDuration,
		ServerTimeout:     testDuration,
	}
	sc = initializeServerConfig(conf, nil)
	assert.Equal(t, testDuration, sc.KaOpts.ServerMinInterval)
	assert.Equal(t, testDuration, sc.KaOpts.ServerInterval)
	assert.Equal(t, testDuration, sc.KaOpts.ServerTimeout)

	sc = initializeServerConfig(conf, nil)
	assert.NotNil(t, sc.Logger)
	assert.Equal(t, &disabled.Provider{}, sc.MetricsProvider)
	assert.Len(t, sc.UnaryInterceptors, 2)
	assert.Len(t, sc.StreamInterceptors, 2)

	sc = initializeServerConfig(conf, &prometheus.Provider{})
	assert.Equal(t, &prometheus.Provider{}, sc.MetricsProvider)

	goodFile := "main.go"
	badFile := "does_not_exist"

	oldLogger := logger
	defer func() { logger = oldLogger }()
	logger, _ = floggingtest.NewTestLogger(t)

	testCases := []struct {
		name           string
		certificate    string
		privateKey     string
		rootCA         string
		clientRootCert string
		clusterCert    string
		clusterKey     string
		clusterCA      string
	}{
		{"BadCertificate", badFile, goodFile, goodFile, goodFile, "", "", ""},
		{"BadPrivateKey", goodFile, badFile, goodFile, goodFile, "", "", ""},
		{"BadRootCA", goodFile, goodFile, badFile, goodFile, "", "", ""},
		{"BadClientRootCertificate", goodFile, goodFile, goodFile, badFile, "", "", ""},
		{"ClusterBadCertificate", goodFile, goodFile, goodFile, goodFile, badFile, goodFile, goodFile},
		{"ClusterBadPrivateKey", goodFile, goodFile, goodFile, goodFile, goodFile, badFile, goodFile},
		{"ClusterBadRootCA", goodFile, goodFile, goodFile, goodFile, goodFile, goodFile, badFile},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := &localconfig.TopLevel{
				General: localconfig.General{
					TLS: localconfig.TLS{
						Enabled:            true,
						ClientAuthRequired: true,
						Certificate:        tc.certificate,
						PrivateKey:         tc.privateKey,
						RootCAs:            []string{tc.rootCA},
						ClientRootCAs:      []string{tc.clientRootCert},
					},
					Cluster: localconfig.Cluster{
						ClientCertificate: tc.clusterCert,
						ClientPrivateKey:  tc.clusterKey,
						RootCAs:           []string{tc.clusterCA},
					},
				},
			}
			assert.Panics(t, func() {
				if tc.clusterCert == "" {
					initializeServerConfig(conf, nil)
				} else {
					initializeClusterClientConfig(conf)
				}
			},
			)
		})
	}
}

func TestInitializeBootstrapChannel(t *testing.T) {
	cleanup := configtest.SetDevFabricConfigPath(t)
	defer cleanup()

	testCases := []struct {
		genesisMethod string
		ledgerType    string
		panics        bool
	}{
		{"provisional", "ram", false},
		{"provisional", "file", false},
		{"provisional", "json", false},
		{"invalid", "ram", true},
		{"file", "ram", true},
	}

	for _, tc := range testCases {

		t.Run(tc.genesisMethod+"/"+tc.ledgerType, func(t *testing.T) {

			fileLedgerLocation, _ := ioutil.TempDir("", "test-ledger")
			ledgerFactory, _ := createLedgerFactory(
				&localconfig.TopLevel{
					General: localconfig.General{LedgerType: tc.ledgerType},
					FileLedger: localconfig.FileLedger{
						Location: fileLedgerLocation,
					},
				},
			)

			bootstrapConfig := &localconfig.TopLevel{
				General: localconfig.General{
					GenesisMethod:  tc.genesisMethod,
					GenesisProfile: "SampleSingleMSPSolo",
					GenesisFile:    "genesisblock",
					SystemChannel:  genesisconfig.TestChainID,
				},
			}

			if tc.panics {
				assert.Panics(t, func() {
					genesisBlock := extractBootstrapBlock(bootstrapConfig)
					initializeBootstrapChannel(genesisBlock, ledgerFactory)
				})
			} else {
				assert.NotPanics(t, func() {
					genesisBlock := extractBootstrapBlock(bootstrapConfig)
					initializeBootstrapChannel(genesisBlock, ledgerFactory)
				})
			}
		})
	}
}

func TestInitializeLocalMsp(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		assert.NotPanics(t, func() {
			localMSPDir, _ := configtest.GetDevMspDir()
			initializeLocalMsp(
				&localconfig.TopLevel{
					General: localconfig.General{
						LocalMSPDir: localMSPDir,
						LocalMSPID:  "SampleOrg",
						BCCSP: &factory.FactoryOpts{
							ProviderName: "SW",
							SwOpts: &factory.SwOpts{
								HashFamily: "SHA2",
								SecLevel:   256,
								Ephemeral:  true,
							},
						},
					},
				})
		})
	})
	t.Run("Error", func(t *testing.T) {
		oldLogger := logger
		defer func() { logger = oldLogger }()
		logger, _ = floggingtest.NewTestLogger(t)

		assert.Panics(t, func() {
			initializeLocalMsp(
				&localconfig.TopLevel{
					General: localconfig.General{
						LocalMSPDir: "",
						LocalMSPID:  "",
					},
				})
		})
	})
}

func TestInitializeMultiChainManager(t *testing.T) {
	cleanup := configtest.SetDevFabricConfigPath(t)
	defer cleanup()
	conf := genesisConfig(t)
	assert.NotPanics(t, func() {
		initializeLocalMsp(conf)
		lf, _ := createLedgerFactory(conf)
		bootBlock := encoder.New(genesisconfig.Load(genesisconfig.SampleDevModeSoloProfile)).GenesisBlockForChannel("system")
		initializeMultichannelRegistrar(bootBlock, &cluster.PredicateDialer{}, comm.ServerConfig{}, nil, conf, localmsp.NewSigner(), &disabled.Provider{}, lf)
	})
}

func TestInitializeGrpcServer(t *testing.T) {
	
	listenAddr := func() string {
		l, _ := net.Listen("tcp", "localhost:0")
		l.Close()
		return l.Addr().String()
	}()
	host := strings.Split(listenAddr, ":")[0]
	port, _ := strconv.ParseUint(strings.Split(listenAddr, ":")[1], 10, 16)
	conf := &localconfig.TopLevel{
		General: localconfig.General{
			ListenAddress: host,
			ListenPort:    uint16(port),
			TLS: localconfig.TLS{
				Enabled:            false,
				ClientAuthRequired: false,
			},
		},
	}
	assert.NotPanics(t, func() {
		grpcServer := initializeGrpcServer(conf, initializeServerConfig(conf, nil))
		grpcServer.Listener().Close()
	})
}

func TestUpdateTrustedRoots(t *testing.T) {
	cleanup := configtest.SetDevFabricConfigPath(t)
	defer cleanup()
	initializeLocalMsp(genesisConfig(t))
	
	listenAddr := func() string {
		l, _ := net.Listen("tcp", "localhost:0")
		l.Close()
		return l.Addr().String()
	}()
	port, _ := strconv.ParseUint(strings.Split(listenAddr, ":")[1], 10, 16)
	conf := &localconfig.TopLevel{
		General: localconfig.General{
			ListenAddress: "localhost",
			ListenPort:    uint16(port),
			TLS: localconfig.TLS{
				Enabled:            false,
				ClientAuthRequired: false,
			},
		},
	}
	grpcServer := initializeGrpcServer(conf, initializeServerConfig(conf, nil))
	caSupport := &comm.CASupport{
		AppRootCAsByChain:     make(map[string][][]byte),
		OrdererRootCAsByChain: make(map[string][][]byte),
	}
	callback := func(bundle *channelconfig.Bundle) {
		if grpcServer.MutualTLSRequired() {
			t.Log("callback called")
			updateTrustedRoots(caSupport, bundle, grpcServer)
		}
	}
	lf, _ := createLedgerFactory(conf)
	bootBlock := encoder.New(genesisconfig.Load(genesisconfig.SampleDevModeSoloProfile)).GenesisBlockForChannel("system")
	initializeMultichannelRegistrar(bootBlock, &cluster.PredicateDialer{}, comm.ServerConfig{}, nil, genesisConfig(t), localmsp.NewSigner(), &disabled.Provider{}, lf, callback)
	t.Logf("# app CAs: %d", len(caSupport.AppRootCAsByChain[genesisconfig.TestChainID]))
	t.Logf("# orderer CAs: %d", len(caSupport.OrdererRootCAsByChain[genesisconfig.TestChainID]))
	
	assert.Equal(t, 0, len(caSupport.AppRootCAsByChain[genesisconfig.TestChainID]))
	assert.Equal(t, 0, len(caSupport.OrdererRootCAsByChain[genesisconfig.TestChainID]))
	grpcServer.Listener().Close()

	conf = &localconfig.TopLevel{
		General: localconfig.General{
			ListenAddress: "localhost",
			ListenPort:    uint16(port),
			TLS: localconfig.TLS{
				Enabled:            true,
				ClientAuthRequired: true,
				PrivateKey:         filepath.Join(".", "testdata", "tls", "server.key"),
				Certificate:        filepath.Join(".", "testdata", "tls", "server.crt"),
			},
		},
	}
	grpcServer = initializeGrpcServer(conf, initializeServerConfig(conf, nil))
	caSupport = &comm.CASupport{
		AppRootCAsByChain:     make(map[string][][]byte),
		OrdererRootCAsByChain: make(map[string][][]byte),
	}

	predDialer := &cluster.PredicateDialer{}
	clusterConf := initializeClusterClientConfig(conf)
	predDialer.SetConfig(clusterConf)

	callback = func(bundle *channelconfig.Bundle) {
		if grpcServer.MutualTLSRequired() {
			t.Log("callback called")
			updateTrustedRoots(caSupport, bundle, grpcServer)
			updateClusterDialer(caSupport, predDialer, clusterConf.SecOpts.ServerRootCAs)
		}
	}
	initializeMultichannelRegistrar(bootBlock, &cluster.PredicateDialer{}, comm.ServerConfig{}, nil, genesisConfig(t), localmsp.NewSigner(), &disabled.Provider{}, lf, callback)
	t.Logf("# app CAs: %d", len(caSupport.AppRootCAsByChain[genesisconfig.TestChainID]))
	t.Logf("# orderer CAs: %d", len(caSupport.OrdererRootCAsByChain[genesisconfig.TestChainID]))
	
	
	assert.Equal(t, 2, len(caSupport.AppRootCAsByChain[genesisconfig.TestChainID]))
	assert.Equal(t, 2, len(caSupport.OrdererRootCAsByChain[genesisconfig.TestChainID]))
	assert.Len(t, predDialer.Config.Load().(comm.ClientConfig).SecOpts.ServerRootCAs, 2)
	grpcServer.Listener().Close()
}

func TestConfigureClusterListener(t *testing.T) {
	logEntries := make(chan string, 100)

	backupLogger := logger
	logger = logger.With(zap.Hooks(func(entry zapcore.Entry) error {
		logEntries <- entry.Message
		return nil
	}))

	defer func() {
		logger = backupLogger
	}()

	ca, err := tlsgen.NewCA()
	assert.NoError(t, err)
	serverKeyPair, err := ca.NewServerCertKeyPair("127.0.0.1")
	assert.NoError(t, err)

	loadPEM := func(fileName string) ([]byte, error) {
		switch fileName {
		case "cert":
			return serverKeyPair.Cert, nil
		case "key":
			return serverKeyPair.Key, nil
		case "ca":
			return ca.CertBytes(), nil
		default:
			return nil, errors.New("I/O error")
		}
	}

	for _, testCase := range []struct {
		name               string
		conf               *localconfig.TopLevel
		generalConf        comm.ServerConfig
		generalSrv         *comm.GRPCServer
		shouldBeEqual      bool
		expectedPanic      string
		expectedLogEntries []string
	}{
		{
			name:          "no separate listener",
			shouldBeEqual: true,
			generalConf:   comm.ServerConfig{},
			conf:          &localconfig.TopLevel{},
			generalSrv:    &comm.GRPCServer{},
		},
		{
			name:        "partial configuration",
			generalConf: comm.ServerConfig{},
			conf: &localconfig.TopLevel{
				General: localconfig.General{
					Cluster: localconfig.Cluster{
						ListenPort: 5000,
					},
				},
			},
			expectedPanic: "Options: General.Cluster.ListenPort, General.Cluster.ListenAddress, " +
				"General.Cluster.ServerCertificate, General.Cluster.ServerPrivateKey, should be defined altogether.",
			generalSrv: &comm.GRPCServer{},
			expectedLogEntries: []string{"Options: General.Cluster.ListenPort, General.Cluster.ListenAddress, " +
				"General.Cluster.ServerCertificate," +
				" General.Cluster.ServerPrivateKey, should be defined altogether."},
		},
		{
			name:        "invalid certificate",
			generalConf: comm.ServerConfig{},
			conf: &localconfig.TopLevel{
				General: localconfig.General{
					Cluster: localconfig.Cluster{
						ListenAddress:     "127.0.0.1",
						ListenPort:        5000,
						ServerPrivateKey:  "key",
						ServerCertificate: "bad",
						RootCAs:           []string{"ca"},
					},
				},
			},
			expectedPanic:      "Failed to load cluster server certificate from 'bad' (I/O error)",
			generalSrv:         &comm.GRPCServer{},
			expectedLogEntries: []string{"Failed to load cluster server certificate from 'bad' (I/O error)"},
		},
		{
			name:        "invalid key",
			generalConf: comm.ServerConfig{},
			conf: &localconfig.TopLevel{
				General: localconfig.General{
					Cluster: localconfig.Cluster{
						ListenAddress:     "127.0.0.1",
						ListenPort:        5000,
						ServerPrivateKey:  "bad",
						ServerCertificate: "cert",
						RootCAs:           []string{"ca"},
					},
				},
			},
			expectedPanic:      "Failed to load cluster server key from 'bad' (I/O error)",
			generalSrv:         &comm.GRPCServer{},
			expectedLogEntries: []string{"Failed to load cluster server certificate from 'bad' (I/O error)"},
		},
		{
			name:        "invalid ca cert",
			generalConf: comm.ServerConfig{},
			conf: &localconfig.TopLevel{
				General: localconfig.General{
					Cluster: localconfig.Cluster{
						ListenAddress:     "127.0.0.1",
						ListenPort:        5000,
						ServerPrivateKey:  "key",
						ServerCertificate: "cert",
						RootCAs:           []string{"bad"},
					},
				},
			},
			expectedPanic:      "Failed to load CA cert file 'I/O error' (bad)",
			generalSrv:         &comm.GRPCServer{},
			expectedLogEntries: []string{"Failed to load CA cert file 'I/O error' (bad)"},
		},
		{
			name:        "bad listen address",
			generalConf: comm.ServerConfig{},
			conf: &localconfig.TopLevel{
				General: localconfig.General{
					Cluster: localconfig.Cluster{
						ListenAddress:     "99.99.99.99",
						ListenPort:        5000,
						ServerPrivateKey:  "key",
						ServerCertificate: "cert",
						RootCAs:           []string{"ca"},
					},
				},
			},
			expectedPanic: "Failed creating gRPC server on 99.99.99.99:5000 due to " +
				"listen tcp 99.99.99.99:5000: bind: cannot assign requested address",
			generalSrv: &comm.GRPCServer{},
		},
		{
			name:        "green path",
			generalConf: comm.ServerConfig{},
			conf: &localconfig.TopLevel{
				General: localconfig.General{
					Cluster: localconfig.Cluster{
						ListenAddress:     "127.0.0.1",
						ListenPort:        5000,
						ServerPrivateKey:  "key",
						ServerCertificate: "cert",
						RootCAs:           []string{"ca"},
					},
				},
			},
			generalSrv: &comm.GRPCServer{},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.shouldBeEqual {
				conf, srv := configureClusterListener(testCase.conf, testCase.generalConf, testCase.generalSrv, loadPEM)
				assert.Equal(t, conf, testCase.generalConf)
				assert.Equal(t, srv, testCase.generalSrv)
			}

			if testCase.expectedPanic != "" {
				f := func() {
					configureClusterListener(testCase.conf, testCase.generalConf, testCase.generalSrv, loadPEM)
				}
				assert.PanicsWithValue(t, testCase.expectedPanic, f)
			} else {
				configureClusterListener(testCase.conf, testCase.generalConf, testCase.generalSrv, loadPEM)
			}
			
			var loggedMessages []string
			for len(logEntries) > 0 {
				logEntry := <-logEntries
				loggedMessages = append(loggedMessages, logEntry)
			}
			assert.Subset(t, testCase.expectedLogEntries, loggedMessages)
		})
	}
}

func genesisConfig(t *testing.T) *localconfig.TopLevel {
	t.Helper()
	localMSPDir, _ := configtest.GetDevMspDir()
	return &localconfig.TopLevel{
		General: localconfig.General{
			LedgerType:     "ram",
			GenesisMethod:  "provisional",
			GenesisProfile: "SampleDevModeSolo",
			SystemChannel:  genesisconfig.TestChainID,
			LocalMSPDir:    localMSPDir,
			LocalMSPID:     "SampleOrg",
			BCCSP: &factory.FactoryOpts{
				ProviderName: "SW",
				SwOpts: &factory.SwOpts{
					HashFamily: "SHA2",
					SecLevel:   256,
					Ephemeral:  true,
				},
			},
		},
	}
}
