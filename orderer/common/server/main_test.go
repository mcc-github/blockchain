


package server

import (
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/flogging/floggingtest"
	"github.com/mcc-github/blockchain/common/localmsp"
	genesisconfig "github.com/mcc-github/blockchain/common/tools/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/stretchr/testify/assert"
)

func TestInitializeLoggingLevel(t *testing.T) {
	initializeLoggingLevel(
		&localconfig.TopLevel{
			
			
			
			
			General: localconfig.General{LogLevel: "foo=debug"},
		},
	)
	assert.Equal(t, flogging.GetModuleLevel("foo"), "DEBUG")
}

func TestInitializeProfilingService(t *testing.T) {
	
	listenAddr := func() string {
		l, _ := net.Listen("tcp", "localhost:0")
		l.Close()
		return l.Addr().String()
	}()
	initializeProfilingService(
		&localconfig.TopLevel{
			General: localconfig.General{
				LogLevel: "debug",
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
	sc := initializeServerConfig(conf)
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
	sc = initializeServerConfig(conf)
	assert.Equal(t, testDuration, sc.KaOpts.ServerMinInterval)
	assert.Equal(t, testDuration, sc.KaOpts.ServerInterval)
	assert.Equal(t, testDuration, sc.KaOpts.ServerTimeout)

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
					initializeServerConfig(conf)
				} else {
					initializeClusterConfig(conf)
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
					initializeBootstrapChannel(bootstrapConfig, ledgerFactory)
				})
			} else {
				assert.NotPanics(t, func() {
					initializeBootstrapChannel(bootstrapConfig, ledgerFactory)
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
		initializeMultichannelRegistrar(conf, localmsp.NewSigner())
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
		grpcServer := initializeGrpcServer(conf, initializeServerConfig(conf))
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
	grpcServer := initializeGrpcServer(conf, initializeServerConfig(conf))
	caSupport := &comm.CASupport{
		AppRootCAsByChain:     make(map[string][][]byte),
		OrdererRootCAsByChain: make(map[string][][]byte),
	}
	callback := func(bundle *channelconfig.Bundle) {
		if grpcServer.MutualTLSRequired() {
			t.Log("callback called")
			updateTrustedRoots(grpcServer, caSupport, bundle)
		}
	}
	initializeMultichannelRegistrar(genesisConfig(t), localmsp.NewSigner(), callback)
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
	grpcServer = initializeGrpcServer(conf, initializeServerConfig(conf))
	caSupport = &comm.CASupport{
		AppRootCAsByChain:     make(map[string][][]byte),
		OrdererRootCAsByChain: make(map[string][][]byte),
	}
	callback = func(bundle *channelconfig.Bundle) {
		if grpcServer.MutualTLSRequired() {
			t.Log("callback called")
			updateTrustedRoots(grpcServer, caSupport, bundle)
		}
	}
	initializeMultichannelRegistrar(genesisConfig(t), localmsp.NewSigner(), callback)
	t.Logf("# app CAs: %d", len(caSupport.AppRootCAsByChain[genesisconfig.TestChainID]))
	t.Logf("# orderer CAs: %d", len(caSupport.OrdererRootCAsByChain[genesisconfig.TestChainID]))
	
	
	assert.Equal(t, 2, len(caSupport.AppRootCAsByChain[genesisconfig.TestChainID]))
	assert.Equal(t, 2, len(caSupport.OrdererRootCAsByChain[genesisconfig.TestChainID]))
	grpcServer.Listener().Close()
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
