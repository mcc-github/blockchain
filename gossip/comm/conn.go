/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"context"
	"sync"

	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type handler func(message *protoext.SignedGossipMessage)

type blockingBehavior bool

const (
	blockingSend    = blockingBehavior(true)
	nonBlockingSend = blockingBehavior(false)
)

type connFactory interface {
	createConnection(endpoint string, pkiID common.PKIidType) (*connection, error)
}

type connectionStore struct {
	config           ConnConfig
	logger           util.Logger            
	isClosing        bool                   
	shutdownOnce     sync.Once              
	connFactory      connFactory            
	sync.RWMutex                            
	pki2Conn         map[string]*connection 
	destinationLocks map[string]*sync.Mutex 
	
}

func newConnStore(connFactory connFactory, logger util.Logger, config ConnConfig) *connectionStore {
	return &connectionStore{
		connFactory:      connFactory,
		isClosing:        false,
		pki2Conn:         make(map[string]*connection),
		destinationLocks: make(map[string]*sync.Mutex),
		logger:           logger,
		config:           config,
	}
}

func (cs *connectionStore) getConnection(peer *RemotePeer) (*connection, error) {
	cs.RLock()
	isClosing := cs.isClosing
	cs.RUnlock()

	if isClosing {
		return nil, errors.Errorf("conn store is closing")
	}

	pkiID := peer.PKIID
	endpoint := peer.Endpoint

	cs.Lock()
	destinationLock, hasConnected := cs.destinationLocks[string(pkiID)]
	if !hasConnected {
		destinationLock = &sync.Mutex{}
		cs.destinationLocks[string(pkiID)] = destinationLock
	}
	cs.Unlock()

	destinationLock.Lock()

	cs.RLock()
	conn, exists := cs.pki2Conn[string(pkiID)]
	if exists {
		cs.RUnlock()
		destinationLock.Unlock()
		return conn, nil
	}
	cs.RUnlock()

	createdConnection, err := cs.connFactory.createConnection(endpoint, pkiID)

	destinationLock.Unlock()

	cs.RLock()
	isClosing = cs.isClosing
	cs.RUnlock()
	if isClosing {
		return nil, errors.Errorf("conn store is closing")
	}

	cs.Lock()
	delete(cs.destinationLocks, string(pkiID))
	defer cs.Unlock()

	
	conn, exists = cs.pki2Conn[string(pkiID)]

	if exists {
		if createdConnection != nil {
			createdConnection.close()
		}
		return conn, nil
	}

	
	if err != nil {
		return nil, err
	}

	
	conn = createdConnection
	cs.pki2Conn[string(createdConnection.pkiID)] = conn

	go conn.serviceConnection()

	return conn, nil
}

func (cs *connectionStore) connNum() int {
	cs.RLock()
	defer cs.RUnlock()
	return len(cs.pki2Conn)
}

func (cs *connectionStore) shutdown() {
	cs.shutdownOnce.Do(func() {
		cs.Lock()
		cs.isClosing = true

		for _, conn := range cs.pki2Conn {
			conn.close()
		}
		cs.pki2Conn = make(map[string]*connection)

		cs.Unlock()
	})
}



func (cs *connectionStore) onConnected(serverStream proto.Gossip_GossipStreamServer,
	connInfo *protoext.ConnectionInfo, metrics *metrics.CommMetrics) *connection {
	cs.Lock()
	defer cs.Unlock()

	if c, exists := cs.pki2Conn[string(connInfo.ID)]; exists {
		c.close()
	}

	conn := newConnection(nil, nil, serverStream, metrics, cs.config)
	conn.pkiID = connInfo.ID
	conn.info = connInfo
	conn.logger = cs.logger
	cs.pki2Conn[string(connInfo.ID)] = conn
	return conn
}

func (cs *connectionStore) closeConnByPKIid(pkiID common.PKIidType) {
	cs.Lock()
	defer cs.Unlock()
	if conn, exists := cs.pki2Conn[string(pkiID)]; exists {
		conn.close()
		delete(cs.pki2Conn, string(pkiID))
	}
}

func newConnection(cl proto.GossipClient, c *grpc.ClientConn, s stream, metrics *metrics.CommMetrics, config ConnConfig) *connection {
	connection := &connection{
		metrics:      metrics,
		outBuff:      make(chan *msgSending, config.SendBuffSize),
		cl:           cl,
		conn:         c,
		gossipStream: s,
		stopChan:     make(chan struct{}, 1),
		recvBuffSize: config.RecvBuffSize,
	}
	return connection
}


type ConnConfig struct {
	RecvBuffSize int
	SendBuffSize int
}

type connection struct {
	recvBuffSize int
	metrics      *metrics.CommMetrics
	cancel       context.CancelFunc
	info         *protoext.ConnectionInfo
	outBuff      chan *msgSending
	logger       util.Logger        
	pkiID        common.PKIidType   
	handler      handler            
	conn         *grpc.ClientConn   
	cl           proto.GossipClient 
	gossipStream stream             
	stopChan     chan struct{}      
	stopOnce     sync.Once          
}

func (conn *connection) close() {
	conn.stopOnce.Do(func() {
		close(conn.stopChan)

		if conn.conn != nil {
			conn.conn.Close()
		}

		if conn.cancel != nil {
			conn.cancel()
		}
	})
}

func (conn *connection) send(msg *protoext.SignedGossipMessage, onErr func(error), shouldBlock blockingBehavior) {
	m := &msgSending{
		envelope: msg.Envelope,
		onErr:    onErr,
	}

	select {
	case conn.outBuff <- m:
		
	case <-conn.stopChan:
		conn.logger.Debugf("Aborting send() to %s because connection is closing", conn.info.Endpoint)
	default: 
		if shouldBlock {
			select {
			case conn.outBuff <- m: 
			case <-conn.stopChan: 
			}
		} else {
			conn.metrics.BufferOverflow.Add(1)
			conn.logger.Debugf("Buffer to %s overflowed, dropping message %s", conn.info.Endpoint, msg)
		}
	}
}

func (conn *connection) serviceConnection() error {
	errChan := make(chan error, 1)
	msgChan := make(chan *protoext.SignedGossipMessage, conn.recvBuffSize)
	
	
	
	
	
	go conn.readFromStream(errChan, msgChan)

	go conn.writeToStream()

	for {
		select {
		case <-conn.stopChan:
			conn.logger.Debug("Closing reading from stream")
			return nil
		case err := <-errChan:
			return err
		case msg := <-msgChan:
			conn.handler(msg)
		}
	}
}

func (conn *connection) writeToStream() {
	stream := conn.gossipStream
	for {
		select {
		case m := <-conn.outBuff:
			err := stream.Send(m.envelope)
			if err != nil {
				go m.onErr(err)
				return
			}
			conn.metrics.SentMessages.Add(1)
		case <-conn.stopChan:
			conn.logger.Debug("Closing writing to stream")
			return
		}
	}
}

func (conn *connection) readFromStream(errChan chan error, msgChan chan *protoext.SignedGossipMessage) {
	stream := conn.gossipStream
	for {
		select {
		case <-conn.stopChan:
			return
		default:
			envelope, err := stream.Recv()
			if err != nil {
				errChan <- err
				conn.logger.Debugf("Got error, aborting: %v", err)
				return
			}
			conn.metrics.ReceivedMessages.Add(1)
			msg, err := protoext.EnvelopeToGossipMessage(envelope)
			if err != nil {
				errChan <- err
				conn.logger.Warningf("Got error, aborting: %v", err)
			}
			select {
			case <-conn.stopChan:
			case msgChan <- msg:
			}
		}
	}
}

type msgSending struct {
	envelope *proto.Envelope
	onErr    func(error)
}
