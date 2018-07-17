/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package performance

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mcc-github/blockchain/common/localmsp"
	"github.com/mcc-github/blockchain/common/tools/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/common/tools/configtxgen/localconfig"

	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	protosutils "github.com/mcc-github/blockchain/protos/utils"
)

const (
	
	Kilo = 1024 
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

var seedOnce sync.Once


func MakeNormalTx(channelID string, size int) *cb.Envelope {
	env, err := protosutils.CreateSignedEnvelope(
		cb.HeaderType_ENDORSER_TRANSACTION,
		channelID,
		localmsp.NewSigner(),
		&cb.Envelope{Payload: make([]byte, size*Kilo)},
		0,
		0,
	)
	if err != nil {
		panic(fmt.Errorf("Failed to create signed envelope because: %s", err))
	}

	return env
}


func OrdererExecWithArgs(f func(s *BenchmarkServer, i ...interface{}), i ...interface{}) {
	servers := GetBenchmarkServerPool()
	var wg sync.WaitGroup
	wg.Add(len(servers))
	for _, server := range servers {
		go func(server *BenchmarkServer) {
			f(server, i...)
			wg.Done()
		}(server)
	}
	wg.Wait()
}


func OrdererExec(f func(s *BenchmarkServer)) {
	servers := GetBenchmarkServerPool()
	var wg sync.WaitGroup
	wg.Add(len(servers))
	for _, server := range servers {
		go func(server *BenchmarkServer) {
			f(server)
			wg.Done()
		}(server)
	}
	wg.Wait()
}


func RandomID(num int) string {
	seedOnce.Do(func() { rand.Seed(time.Now().UnixNano()) })

	b := make([]rune, num)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


func CreateChannel(server *BenchmarkServer, channelProfile *genesisconfig.Profile) string {
	client := server.CreateBroadcastClient()
	defer client.Close()

	channelID := RandomID(10)
	createChannelTx, err := encoder.MakeChannelCreationTransaction(channelID, localmsp.NewSigner(), nil, channelProfile)
	if err != nil {
		logger.Panicf("Failed to create channel creation transaction: %s", err)
	}
	client.SendRequest(createChannelTx)
	response := client.GetResponse()
	if response.Status != cb.Status_SUCCESS {
		logger.Panicf("Failed to create channel: %s -- %v:%s", channelID, response.Status, response.Info)
	}
	return channelID
}


func WaitForChannels(server *BenchmarkServer, channelIDs ...interface{}) {
	var scoutWG sync.WaitGroup
	scoutWG.Add(len(channelIDs))
	for _, channelID := range channelIDs {
		id, ok := channelID.(string)
		if !ok {
			panic("Expect a string as channelID")
		}
		go func(channelID string) {
			logger.Infof("Scouting for channel: %s", channelID)
			for {
				status, err := SeekAllBlocks(server.CreateDeliverClient(), channelID, 0)
				if err != nil {
					panic(fmt.Errorf("Failed to call deliver because: %s", err))
				}

				switch status {
				case cb.Status_SUCCESS:
					logger.Infof("Channel '%s' is ready", channelID)
					scoutWG.Done()
					return
				case cb.Status_SERVICE_UNAVAILABLE:
					fallthrough
				case cb.Status_NOT_FOUND:
					logger.Debugf("Channel '%s' is not ready yet, keep scouting", channelID)
					time.Sleep(time.Second)
				default:
					logger.Fatalf("Unexpected reply status '%s' while scouting for channel %s, exit", status.String(), channelID)
				}
			}
		}(id)
	}
	scoutWG.Wait()
}

var seekOldest = &ab.SeekPosition{Type: &ab.SeekPosition_Oldest{Oldest: &ab.SeekOldest{}}}


func SeekAllBlocks(c *DeliverClient, channelID string, number uint64) (status cb.Status, err error) {
	env, err := protosutils.CreateSignedEnvelope(
		cb.HeaderType_DELIVER_SEEK_INFO,
		channelID,
		localmsp.NewSigner(),
		&ab.SeekInfo{Start: seekOldest, Stop: seekSpecified(number), Behavior: ab.SeekInfo_BLOCK_UNTIL_READY},
		0,
		0,
	)
	if err != nil {
		panic(fmt.Errorf("Failed to create signed envelope because: %s", err))
	}

	c.SendRequest(env)

	for {
		select {
		case reply := <-c.ResponseChan:
			if reply.GetBlock() == nil {
				status = reply.GetStatus()
				c.Close()
			}
		case err = <-c.ResultChan:
			return
		}
	}
}

func seekSpecified(number uint64) *ab.SeekPosition {
	return &ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: number}}}
}
