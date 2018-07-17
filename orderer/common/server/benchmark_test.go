


package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	genesisconfig "github.com/mcc-github/blockchain/common/tools/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	perf "github.com/mcc-github/blockchain/orderer/common/performance"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/stretchr/testify/assert"
)
















































































const (
	MaxMessageCount = 10

	
	
	
	AbsoluteMaxBytes  = 16 
	PreferredMaxBytes = 10 
	ChannelProfile    = genesisconfig.SampleSingleMSPChannelProfile
)

var envvars = map[string]string{
	"ORDERER_GENERAL_GENESISPROFILE":                              genesisconfig.SampleDevModeSoloProfile,
	"ORDERER_GENERAL_LEDGERTYPE":                                  "file",
	"ORDERER_GENERAL_LOGLEVEL":                                    "error",
	"ORDERER_KAFKA_VERBOSE":                                       "false",
	genesisconfig.Prefix + "_ORDERER_BATCHSIZE_MAXMESSAGECOUNT":   strconv.Itoa(MaxMessageCount),
	genesisconfig.Prefix + "_ORDERER_BATCHSIZE_ABSOLUTEMAXBYTES":  strconv.Itoa(AbsoluteMaxBytes) + " KB",
	genesisconfig.Prefix + "_ORDERER_BATCHSIZE_PREFERREDMAXBYTES": strconv.Itoa(PreferredMaxBytes) + " KB",
	genesisconfig.Prefix + "_ORDERER_KAFKA_BROKERS":               "[localhost:9092]",
}

type factors struct {
	numOfChannels             int 
	totalTx                   int 
	messageSize               int 
	broadcastClientPerChannel int 
	deliverClientPerChannel   int 
	numOfOrderer              int 
}



func (f factors) String() string {
	return fmt.Sprintf(
		"%dch/%dtx/%dkb/%dbc/%ddc/%dord",
		f.numOfChannels,
		f.totalTx,
		f.messageSize,
		f.broadcastClientPerChannel,
		f.deliverClientPerChannel,
		f.numOfOrderer,
	)
}




func TestOrdererBenchmarkSolo(t *testing.T) {
	for key, value := range envvars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	t.Run("Benchmark Sample Test (Solo)", func(t *testing.T) {
		benchmarkOrderer(t, 1, 5, PreferredMaxBytes, 1, 0, 1, true)
	})
}


func TestOrdererBenchmarkSoloBroadcast(t *testing.T) {
	if os.Getenv("BENCHMARK") == "" {
		t.Skip("Skipping benchmark test")
	}

	for key, value := range envvars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	var (
		channelCounts             = []int{1, 10, 50}
		totalTx                   = []int{10000}
		messagesSizes             = []int{1, 2, 10}
		broadcastClientPerChannel = []int{1, 10, 50}
		deliverClientPerChannel   = []int{0} 
		numOfOrderer              = []int{1}

		args = [][]int{
			channelCounts,
			totalTx,
			messagesSizes,
			broadcastClientPerChannel,
			deliverClientPerChannel,
			numOfOrderer,
		}
	)

	for factors := range combinations(args) {
		t.Run(factors.String(), func(t *testing.T) {
			benchmarkOrderer(
				t,
				factors.numOfChannels,
				factors.totalTx,
				factors.messageSize,
				factors.broadcastClientPerChannel,
				factors.deliverClientPerChannel,
				1, 
				true,
			)
		})
	}
}


func TestOrdererBenchmarkSoloDeliver(t *testing.T) {
	if os.Getenv("BENCHMARK") == "" {
		t.Skip("Skipping benchmark test")
	}

	for key, value := range envvars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	var (
		channelCounts             = []int{1, 10, 50}
		totalTx                   = []int{10000}
		messagesSizes             = []int{1, 2, 10}
		broadcastClientPerChannel = []int{100}
		deliverClientPerChannel   = []int{1, 10, 50}
		numOfOrderer              = []int{1}

		args = [][]int{
			channelCounts,
			totalTx,
			messagesSizes,
			broadcastClientPerChannel,
			deliverClientPerChannel,
			numOfOrderer,
		}
	)

	for factors := range combinations(args) {
		t.Run(factors.String(), func(t *testing.T) {
			benchmarkOrderer(
				t,
				factors.numOfChannels,
				factors.totalTx,
				factors.messageSize,
				factors.broadcastClientPerChannel,
				factors.deliverClientPerChannel,
				1, 
				false,
			)
		})
	}
}


func TestOrdererBenchmarkKafkaBroadcast(t *testing.T) {
	if os.Getenv("BENCHMARK") == "" {
		t.Skip("Skipping benchmark test")
	}

	for key, value := range envvars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	os.Setenv("ORDERER_GENERAL_GENESISPROFILE", genesisconfig.SampleDevModeKafkaProfile)
	defer os.Unsetenv("ORDERER_GENERAL_GENESISPROFILE")

	var (
		channelCounts             = []int{1, 10}
		totalTx                   = []int{10000}
		messagesSizes             = []int{1, 2, 10}
		broadcastClientPerChannel = []int{1, 10, 50}
		deliverClientPerChannel   = []int{0} 
		numOfOrderer              = []int{1, 5, 10}

		args = [][]int{
			channelCounts,
			totalTx,
			messagesSizes,
			broadcastClientPerChannel,
			deliverClientPerChannel,
			numOfOrderer,
		}
	)

	for factors := range combinations(args) {
		t.Run(factors.String(), func(t *testing.T) {
			benchmarkOrderer(
				t,
				factors.numOfChannels,
				factors.totalTx,
				factors.messageSize,
				factors.broadcastClientPerChannel,
				factors.deliverClientPerChannel,
				factors.numOfOrderer,
				true,
			)
		})
	}
}


func TestOrdererBenchmarkKafkaDeliver(t *testing.T) {
	if os.Getenv("BENCHMARK") == "" {
		t.Skip("Skipping benchmark test")
	}

	for key, value := range envvars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	os.Setenv("ORDERER_GENERAL_GENESISPROFILE", genesisconfig.SampleDevModeKafkaProfile)
	defer os.Unsetenv("ORDERER_GENERAL_GENESISPROFILE")

	var (
		channelCounts             = []int{1, 10}
		totalTx                   = []int{10000}
		messagesSizes             = []int{1, 2, 10}
		broadcastClientPerChannel = []int{50}
		deliverClientPerChannel   = []int{1, 10, 50}
		numOfOrderer              = []int{1, 5, 10}

		args = [][]int{
			channelCounts,
			totalTx,
			messagesSizes,
			broadcastClientPerChannel,
			deliverClientPerChannel,
			numOfOrderer,
		}
	)

	for factors := range combinations(args) {
		t.Run(factors.String(), func(t *testing.T) {
			benchmarkOrderer(
				t,
				factors.numOfChannels,
				factors.totalTx,
				factors.messageSize,
				factors.broadcastClientPerChannel,
				factors.deliverClientPerChannel,
				factors.numOfOrderer,
				false,
			)
		})
	}
}

func benchmarkOrderer(
	t *testing.T,
	numOfChannels int,
	totalTx int,
	msgSize int,
	broadcastClientPerChannel int,
	deliverClientPerChannel int,
	numOfOrderer int,
	multiplex bool,
) {
	cleanup := configtest.SetDevFabricConfigPath(t)
	defer cleanup()

	
	conf, err := localconfig.Load()
	if err != nil {
		t.Fatal("failed to load config")
	}

	initializeLoggingLevel(conf)
	initializeLocalMsp(conf)
	perf.InitializeServerPool(numOfOrderer)

	
	channelProfile := genesisconfig.Load(ChannelProfile)

	
	
	txPerClient := totalTx / (broadcastClientPerChannel * numOfChannels * numOfOrderer)
	
	
	
	totalTx = txPerClient * broadcastClientPerChannel * numOfChannels * numOfOrderer

	
	txPerChannel := totalTx / numOfChannels
	
	msg := perf.MakeNormalTx("abcdefghij", msgSize)
	actualMsgSize := len(msg.Payload) + len(msg.Signature)
	
	txPerBlk := min(MaxMessageCount, max(1, PreferredMaxBytes*perf.Kilo/actualMsgSize))
	
	blkPerChannel := txPerChannel / txPerBlk

	var txCount uint64 

	
	
	systemchannel := "system-channel-" + perf.RandomID(5)
	conf.General.SystemChannel = systemchannel

	
	for i := 0; i < numOfOrderer; i++ {
		
		
		
		
		
		
		
		
		localConf := localconfig.TopLevel(*conf)
		if localConf.General.LedgerType != "ram" {
			tempDir, err := ioutil.TempDir("", "blockchain-benchmark-test-")
			assert.NoError(t, err, "Should be able to create temp dir")
			localConf.FileLedger.Location = tempDir
			defer os.RemoveAll(tempDir)
		}

		go Start("benchmark", &localConf)
	}

	defer perf.OrdererExec(perf.Halt)

	
	perf.OrdererExec(perf.WaitForService)
	perf.OrdererExecWithArgs(perf.WaitForChannels, systemchannel)

	
	benchmarkServers := perf.GetBenchmarkServerPool()
	channelIDs := make([]string, numOfChannels)
	txs := make(map[string]*cb.Envelope)
	for i := 0; i < numOfChannels; i++ {
		id := perf.CreateChannel(benchmarkServers[0], channelProfile) 
		channelIDs[i] = id
		txs[id] = perf.MakeNormalTx(id, msgSize)
	}

	
	perf.OrdererExecWithArgs(perf.WaitForChannels, stoi(channelIDs)...)

	
	broadcast := func(wg *sync.WaitGroup) {
		perf.OrdererExec(func(server *perf.BenchmarkServer) {
			var broadcastWG sync.WaitGroup
			
			
			
			
			broadcastWG.Add(numOfChannels * (broadcastClientPerChannel + 1))

			for _, channelID := range channelIDs {
				go func(channelID string) {
					
					go func() {
						deliverClient := server.CreateDeliverClient()
						status, err := perf.SeekAllBlocks(deliverClient, channelID, uint64(blkPerChannel))
						assert.Equal(t, cb.Status_SUCCESS, status, "Expect deliver reply to be SUCCESS")
						assert.NoError(t, err, "Expect deliver handler to exist normally")

						broadcastWG.Done()
					}()

					for c := 0; c < broadcastClientPerChannel; c++ {
						go func() {
							broadcastClient := server.CreateBroadcastClient()
							defer func() {
								broadcastClient.Close()
								err := <-broadcastClient.Errors()
								assert.NoError(t, err, "Expect broadcast handler to shutdown gracefully")
							}()

							for i := 0; i < txPerClient; i++ {
								atomic.AddUint64(&txCount, 1)
								broadcastClient.SendRequest(txs[channelID])
								assert.Equal(t, cb.Status_SUCCESS, broadcastClient.GetResponse().Status, "Expect enqueue to succeed")
							}
							broadcastWG.Done()
						}()
					}
				}(channelID)
			}

			broadcastWG.Wait()
		})

		if wg != nil {
			wg.Done()
		}
	}

	
	deliver := func(wg *sync.WaitGroup) {
		perf.OrdererExec(func(server *perf.BenchmarkServer) {
			var deliverWG sync.WaitGroup
			deliverWG.Add(deliverClientPerChannel * numOfChannels)
			for g := 0; g < deliverClientPerChannel; g++ {
				go func() {
					for _, channelID := range channelIDs {
						go func(channelID string) {
							deliverClient := server.CreateDeliverClient()
							status, err := perf.SeekAllBlocks(deliverClient, channelID, uint64(blkPerChannel))
							assert.Equal(t, cb.Status_SUCCESS, status, "Expect deliver reply to be SUCCESS")
							assert.NoError(t, err, "Expect deliver handler to exist normally")

							deliverWG.Done()
						}(channelID)
					}
				}()
			}
			deliverWG.Wait()
		})

		if wg != nil {
			wg.Done()
		}
	}

	var wg sync.WaitGroup
	var btime, dtime time.Duration

	if multiplex {
		
		start := time.Now()

		wg.Add(2)
		go broadcast(&wg)
		go deliver(&wg)
		wg.Wait()

		btime = time.Since(start)
		dtime = time.Since(start)
	} else {
		
		start := time.Now()
		broadcast(nil)
		btime = time.Since(start)

		start = time.Now()
		deliver(nil)
		dtime = time.Since(start)
	}

	
	
	assert.Equal(t, uint64(totalTx), txCount, "Expected to send %d msg, but actually sent %d", uint64(totalTx), txCount)

	ordererProfile := os.Getenv("ORDERER_GENERAL_GENESISPROFILE")

	fmt.Printf(
		"Messages: %6d  Message Size: %3dKB  Channels: %3d Orderer (%s): %2d | "+
			"Broadcast Clients: %3d  Write tps: %5.1f tx/s Elapsed Time: %0.2fs | "+
			"Deliver clients: %3d  Read tps: %8.1f blk/s Elapsed Time: %0.2fs\n",
		totalTx,
		msgSize,
		numOfChannels,
		ordererProfile,
		numOfOrderer,
		broadcastClientPerChannel*numOfChannels*numOfOrderer,
		float64(totalTx)/btime.Seconds(),
		btime.Seconds(),
		deliverClientPerChannel*numOfChannels*numOfOrderer,
		float64(blkPerChannel*deliverClientPerChannel*numOfChannels)/dtime.Seconds(),
		dtime.Seconds())
}

func combinations(args [][]int) <-chan factors {
	ch := make(chan factors)
	go func() {
		defer close(ch)
		for c := range combine(args) {
			ch <- factors{
				numOfChannels:             c[0],
				totalTx:                   c[1],
				messageSize:               c[2],
				broadcastClientPerChannel: c[3],
				deliverClientPerChannel:   c[4],
				numOfOrderer:              c[5],
			}
		}
	}()
	return ch
}











func combine(args [][]int) <-chan []int {
	ch := make(chan []int)
	go func() {
		defer close(ch)
		if len(args) == 1 {
			for _, i := range args[0] {
				ch <- []int{i}
			}
		} else {
			for _, i := range args[0] {
				for j := range combine(args[1:]) {
					ch <- append([]int{i}, j...)
				}
			}
		}
	}()
	return ch
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func stoi(s []string) (ret []interface{}) {
	ret = make([]interface{}, len(s))
	for i, d := range s {
		ret[i] = d
	}
	return
}
