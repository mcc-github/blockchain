/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/util"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type handler func(message *proto.SignedGossipMessage)

type blockingBehavior bool

const (
	blockingSend    = blockingBehavior(true)
	nonBlockingSend = blockingBehavior(false)
)

type connFactory interface {
	createConnection(endpoint string, pkiID common.PKIidType) (*connection, error)
}

type connectionStore struct {
	logger           *logging.Logger        
	isClosing        bool                   
	connFactory      connFactory            
	sync.RWMutex                            
	pki2Conn         map[string]*connection 
	destinationLocks map[string]*sync.Mutex 
	
}

func newConnStore(connFactory connFactory, logger *logging.Logger) *connectionStore {
	return &connectionStore{
		connFactory:      connFactory,
		isClosing:        false,
		pki2Conn:         make(map[string]*connection),
		destinationLocks: make(map[string]*sync.Mutex),
		logger:           logger,
	}
}

func (cs *connectionStore) getConnection(peer *RemotePeer) (*connection, error) {
	cs.RLock()
	isClosing := cs.isClosing
	cs.RUnlock()

	if isClosing {
		return nil, fmt.Errorf("Shutting down")
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
		return nil, fmt.Errorf("ConnStore is closing")
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

func (cs *connectionStore) closeConn(peer *RemotePeer) {
	cs.Lock()
	defer cs.Unlock()

	if conn, exists := cs.pki2Conn[string(peer.PKIID)]; exists {
		conn.close()
		delete(cs.pki2Conn, string(conn.pkiID))
	}
}

func (cs *connectionStore) shutdown() {
	cs.Lock()
	cs.isClosing = true
	pkiIds2conn := cs.pki2Conn

	var connections2Close []*connection
	for _, conn := range pkiIds2conn {
		connections2Close = append(connections2Close, conn)
	}
	cs.Unlock()

	wg := sync.WaitGroup{}
	for _, conn := range connections2Close {
		wg.Add(1)
		go func(conn *connection) {
			cs.closeByPKIid(conn.pkiID)
			wg.Done()
		}(conn)
	}
	wg.Wait()
}

func (cs *connectionStore) onConnected(serverStream proto.Gossip_GossipStreamServer, connInfo *proto.ConnectionInfo) *connection {
	cs.Lock()
	defer cs.Unlock()

	if c, exists := cs.pki2Conn[string(connInfo.ID)]; exists {
		c.close()
	}

	return cs.registerConn(connInfo, serverStream)
}

func (cs *connectionStore) registerConn(connInfo *proto.ConnectionInfo, serverStream proto.Gossip_GossipStreamServer) *connection {
	conn := newConnection(nil, nil, nil, serverStream)
	conn.pkiID = connInfo.ID
	conn.info = connInfo
	conn.logger = cs.logger
	cs.pki2Conn[string(connInfo.ID)] = conn
	return conn
}

func (cs *connectionStore) closeByPKIid(pkiID common.PKIidType) {
	cs.Lock()
	defer cs.Unlock()
	if conn, exists := cs.pki2Conn[string(pkiID)]; exists {
		conn.close()
		delete(cs.pki2Conn, string(pkiID))
	}
}

func newConnection(cl proto.GossipClient, c *grpc.ClientConn, cs proto.Gossip_GossipStreamClient, ss proto.Gossip_GossipStreamServer) *connection {
	connection := &connection{
		outBuff:      make(chan *msgSending, util.GetIntOrDefault("peer.gossip.sendBuffSize", defSendBuffSize)),
		cl:           cl,
		conn:         c,
		clientStream: cs,
		serverStream: ss,
		stopFlag:     int32(0),
		stopChan:     make(chan struct{}, 1),
	}
	return connection
}

type connection struct {
	cancel       context.CancelFunc
	info         *proto.ConnectionInfo
	outBuff      chan *msgSending
	logger       *logging.Logger                 
	pkiID        common.PKIidType                
	handler      handler                         
	conn         *grpc.ClientConn                
	cl           proto.GossipClient              
	clientStream proto.Gossip_GossipStreamClient 
	serverStream proto.Gossip_GossipStreamServer 
	stopFlag     int32                           
	stopChan     chan struct{}                   
	sync.RWMutex                                 
}

func (conn *connection) close() {
	if conn.toDie() {
		return
	}

	amIFirst := atomic.CompareAndSwapInt32(&conn.stopFlag, int32(0), int32(1))
	if !amIFirst {
		return
	}

	conn.stopChan <- struct{}{}

	conn.drainOutputBuffer()
	conn.Lock()
	defer conn.Unlock()

	if conn.clientStream != nil {
		conn.clientStream.CloseSend()
	}
	if conn.conn != nil {
		conn.conn.Close()
	}

	if conn.cancel != nil {
		conn.cancel()
	}
}

func (conn *connection) toDie() bool {
	return atomic.LoadInt32(&(conn.stopFlag)) == int32(1)
}

func (conn *connection) send(msg *proto.SignedGossipMessage, onErr func(error), shouldBlock blockingBehavior) {
	if conn.toDie() {
		conn.logger.Debug("Aborting send() to ", conn.info.Endpoint, "because connection is closing")
		return
	}

	m := &msgSending{
		envelope: msg.Envelope,
		onErr:    onErr,
	}

	if len(conn.outBuff) == cap(conn.outBuff) {
		if conn.logger.IsEnabledFor(logging.DEBUG) {
			conn.logger.Debug("Buffer to", conn.info.Endpoint, "overflowed, dropping message", msg.String())
		}
		if !shouldBlock {
			return
		}
	}

	conn.outBuff <- m
}

func (conn *connection) serviceConnection() error {
	errChan := make(chan error, 1)
	msgChan := make(chan *proto.SignedGossipMessage, util.GetIntOrDefault("peer.gossip.recvBuffSize", defRecvBuffSize))
	defer close(msgChan)

	
	
	
	
	
	go conn.readFromStream(errChan, msgChan)

	go conn.writeToStream()

	for !conn.toDie() {
		select {
		case stop := <-conn.stopChan:
			conn.logger.Debug("Closing reading from stream")
			conn.stopChan <- stop
			return nil
		case err := <-errChan:
			return err
		case msg := <-msgChan:
			conn.handler(msg)
		}
	}
	return nil
}

func (conn *connection) writeToStream() {
	for !conn.toDie() {
		stream := conn.getStream()
		if stream == nil {
			conn.logger.Error(conn.pkiID, "Stream is nil, aborting!")
			return
		}
		select {
		case m := <-conn.outBuff:
			err := stream.Send(m.envelope)
			if err != nil {
				go m.onErr(err)
				return
			}
		case stop := <-conn.stopChan:
			conn.logger.Debug("Closing writing to stream")
			conn.stopChan <- stop
			return
		}
	}
}

func (conn *connection) drainOutputBuffer() {
	
	for len(conn.outBuff) > 0 {
		<-conn.outBuff
	}
}

func (conn *connection) readFromStream(errChan chan error, msgChan chan *proto.SignedGossipMessage) {
	defer func() {
		recover()
	}() 
	for !conn.toDie() {
		stream := conn.getStream()
		if stream == nil {
			conn.logger.Error(conn.pkiID, "Stream is nil, aborting!")
			errChan <- fmt.Errorf("Stream is nil")
			return
		}
		envelope, err := stream.Recv()
		if conn.toDie() {
			conn.logger.Debug(conn.pkiID, "canceling read because closing")
			return
		}
		if err != nil {
			errChan <- err
			conn.logger.Debugf("%v Got error, aborting: %v", err)
			return
		}
		msg, err := envelope.ToGossipMessage()
		if err != nil {
			errChan <- err
			conn.logger.Warning("%v Got error, aborting: %v", err)
		}
		msgChan <- msg
	}
}

func (conn *connection) getStream() stream {
	conn.Lock()
	defer conn.Unlock()

	if conn.clientStream != nil && conn.serverStream != nil {
		e := errors.New("Both client and server stream are not nil, something went wrong")
		conn.logger.Errorf("%+v", e)
	}

	if conn.clientStream != nil {
		return conn.clientStream
	}

	if conn.serverStream != nil {
		return conn.serverStream
	}

	return nil
}

type msgSending struct {
	envelope *proto.Envelope
	onErr    func(error)
}
