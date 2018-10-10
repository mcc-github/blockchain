/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
)

const (
	
	
	DefaultRPCTimeout = time.Second * 5
)



type ChannelExtractor interface {
	TargetChannel(message proto.Message) string
}




type Handler interface {
	OnStep(channel string, sender uint64, req *orderer.StepRequest) (*orderer.StepResponse, error)
	OnSubmit(channel string, sender uint64, req *orderer.SubmitRequest) (*orderer.SubmitResponse, error)
}


type RemoteNode struct {
	
	ID uint64
	
	Endpoint string
	
	ServerTLSCert []byte
	
	ClientTLSCert []byte
}


func (rm RemoteNode) String() string {
	return fmt.Sprintf("ID: %d\nEndpoint: %s\nServerTLSCert:%s ClientTLSCert:%s",
		rm.ID, rm.Endpoint, DERtoPEM(rm.ServerTLSCert), DERtoPEM(rm.ClientTLSCert))
}




type Communicator interface {
	
	
	
	Remote(channel string, id uint64) (*RemoteContext, error)
	
	
	
	Configure(channel string, members []RemoteNode)
	
	Shutdown()
}



type MembersByChannel map[string]MemberMapping


type Comm struct {
	shutdown     bool
	Lock         sync.RWMutex
	Logger       *flogging.FabricLogger
	ChanExt      ChannelExtractor
	H            Handler
	Connections  *ConnectionStore
	Chan2Members MembersByChannel
	RPCTimeout   time.Duration
}

type requestContext struct {
	channel string
	sender  uint64
}



func (c *Comm) DispatchSubmit(ctx context.Context, request *orderer.SubmitRequest) (*orderer.SubmitResponse, error) {
	c.Logger.Debug(request.Channel)
	reqCtx, err := c.requestContext(ctx, request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return c.H.OnSubmit(reqCtx.channel, reqCtx.sender, request)
}



func (c *Comm) DispatchStep(ctx context.Context, request *orderer.StepRequest) (*orderer.StepResponse, error) {
	reqCtx, err := c.requestContext(ctx, request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return c.H.OnStep(reqCtx.channel, reqCtx.sender, request)
}



func (c *Comm) requestContext(ctx context.Context, msg proto.Message) (*requestContext, error) {
	channel := c.ChanExt.TargetChannel(msg)
	if channel == "" {
		return nil, errors.Errorf("badly formatted message, cannot extract channel")
	}
	c.Lock.RLock()
	mapping, exists := c.Chan2Members[channel]
	c.Lock.RUnlock()

	if !exists {
		return nil, errors.Errorf("channel %s doesn't exist", channel)
	}

	cert := comm.ExtractCertificateFromContext(ctx)
	if len(cert) == 0 {
		return nil, errors.Errorf("no TLS certificate sent")
	}
	stub := mapping.LookupByClientCert(cert)
	if stub == nil {
		return nil, errors.Errorf("certificate extracted from TLS connection isn't authorized")
	}
	return &requestContext{
		channel: channel,
		sender:  stub.ID,
	}, nil
}



func (c *Comm) Remote(channel string, id uint64) (*RemoteContext, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	if c.shutdown {
		return nil, errors.New("communication has been shut down")
	}

	mapping, exists := c.Chan2Members[channel]
	if !exists {
		return nil, errors.Errorf("channel %s doesn't exist", channel)
	}
	stub := mapping.ByID(id)
	if stub == nil {
		return nil, errors.Errorf("node %d doesn't exist in channel %s's membership", id, channel)
	}

	if stub.Active() {
		return stub.RemoteContext, nil
	}

	err := stub.Activate(c.createRemoteContext(stub))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return stub.RemoteContext, nil
}


func (c *Comm) Configure(channel string, newNodes []RemoteNode) {
	c.Logger.Infof("Entering, channel: %s, nodes: %v", channel, newNodes)
	defer c.Logger.Infof("Exiting")

	c.Lock.Lock()
	defer c.Lock.Unlock()

	if c.shutdown {
		return
	}

	beforeConfigChange := c.serverCertsInUse()
	
	c.applyMembershipConfig(channel, newNodes)
	
	c.cleanUnusedConnections(beforeConfigChange)
}


func (c *Comm) Shutdown() {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.shutdown = true
	for _, members := range c.Chan2Members {
		for _, member := range members {
			c.Connections.Disconnect(member.ServerTLSCert)
		}
	}
}



func (c *Comm) cleanUnusedConnections(serverCertsBeforeConfig StringSet) {
	
	serverCertsAfterConfig := c.serverCertsInUse()
	
	serverCertsBeforeConfig.subtract(serverCertsAfterConfig)
	
	for serverCertificate := range serverCertsBeforeConfig {
		c.Connections.Disconnect([]byte(serverCertificate))
	}
}



func (c *Comm) serverCertsInUse() StringSet {
	endpointsInUse := make(StringSet)
	for _, mapping := range c.Chan2Members {
		endpointsInUse.union(mapping.ServerCertificates())
	}
	return endpointsInUse
}


func (c *Comm) applyMembershipConfig(channel string, newNodes []RemoteNode) {
	mapping := c.getOrCreateMapping(channel)
	newNodeIDs := make(map[uint64]struct{})

	for _, node := range newNodes {
		newNodeIDs[node.ID] = struct{}{}
		c.updateStubInMapping(mapping, node)
	}

	
	
	for id, stub := range mapping {
		if _, exists := newNodeIDs[id]; exists {
			c.Logger.Info(id, "exists in both old and new membership, skipping its deactivation")
			continue
		}
		c.Logger.Info("Deactivated node", id, "who's endpoint is", stub.Endpoint, "as it's removed from membership")
		delete(mapping, id)
		stub.Deactivate()
	}
}


func (c *Comm) updateStubInMapping(mapping MemberMapping, node RemoteNode) {
	stub := mapping.ByID(node.ID)
	if stub == nil {
		c.Logger.Info("Allocating a new stub for node", node.ID, "with endpoint of", node.Endpoint)
		stub = &Stub{}
	}

	
	
	
	if !bytes.Equal(stub.ServerTLSCert, node.ServerTLSCert) {
		c.Logger.Info("Deactivating node", node.ID, "with endpoint of", node.Endpoint, "due to TLS certificate change")
		stub.Deactivate()
	}

	
	stub.RemoteNode = node

	
	mapping.Put(stub)

	
	if stub.Active() {
		return
	}

	
	stub.Activate(c.createRemoteContext(stub))
}




func (c *Comm) createRemoteContext(stub *Stub) func() (*RemoteContext, error) {
	return func() (*RemoteContext, error) {
		timeout := c.RPCTimeout
		if timeout == time.Duration(0) {
			timeout = DefaultRPCTimeout
		}

		c.Logger.Debug("Connecting to", stub.RemoteNode, "with gRPC timeout of", timeout)

		conn, err := c.Connections.Connection(stub.Endpoint, stub.ServerTLSCert)
		if err != nil {
			c.Logger.Warningf("Unable to obtain connection to %d(%s): %v", stub.ID, stub.Endpoint, err)
			return nil, err
		}

		clusterClient := orderer.NewClusterClient(conn)

		rc := &RemoteContext{
			RPCTimeout: timeout,
			Client:     clusterClient,
			onAbort: func() {
				c.Logger.Info("Aborted connection to", stub.ID, stub.Endpoint)
				stub.RemoteContext = nil
			},
		}
		return rc, nil
	}
}



func (c *Comm) getOrCreateMapping(channel string) MemberMapping {
	
	mapping, exists := c.Chan2Members[channel]
	if !exists {
		mapping = make(MemberMapping)
		c.Chan2Members[channel] = mapping
	}
	return mapping
}




type Stub struct {
	lock sync.RWMutex
	RemoteNode
	*RemoteContext
}



func (stub *Stub) Active() bool {
	stub.lock.RLock()
	defer stub.lock.RUnlock()
	return stub.isActive()
}



func (stub *Stub) isActive() bool {
	return stub.RemoteContext != nil
}




func (stub *Stub) Deactivate() {
	stub.lock.Lock()
	defer stub.lock.Unlock()
	if !stub.isActive() {
		return
	}
	stub.RemoteContext.Abort()
	stub.RemoteContext = nil
}




func (stub *Stub) Activate(createRemoteContext func() (*RemoteContext, error)) error {
	stub.lock.Lock()
	defer stub.lock.Unlock()
	
	if stub.isActive() {
		return nil
	}
	remoteStub, err := createRemoteContext()
	if err != nil {
		return errors.WithStack(err)
	}

	stub.RemoteContext = remoteStub
	return nil
}



type RemoteContext struct {
	RPCTimeout         time.Duration
	onAbort            func()
	Client             orderer.ClusterClient
	stepLock           sync.Mutex
	cancelStep         func()
	submitLock         sync.Mutex
	cancelSubmitStream func()
	submitStream       orderer.Cluster_SubmitClient
}


func (rc *RemoteContext) SubmitStream() (orderer.Cluster_SubmitClient, error) {
	rc.submitLock.Lock()
	defer rc.submitLock.Unlock()
	
	rc.closeSubmitStream()

	ctx, cancel := context.WithCancel(context.TODO())
	submitStream, err := rc.Client.Submit(ctx)
	if err != nil {
		cancel()
		return nil, errors.WithStack(err)
	}
	rc.submitStream = submitStream
	rc.cancelSubmitStream = cancel
	return rc.submitStream, nil
}


func (rc *RemoteContext) Step(req *orderer.StepRequest) (*orderer.StepResponse, error) {
	ctx, abort := context.WithCancel(context.TODO())
	ctx, cancel := context.WithTimeout(ctx, rc.RPCTimeout)
	defer cancel()

	rc.stepLock.Lock()
	rc.cancelStep = abort
	rc.stepLock.Unlock()

	return rc.Client.Step(ctx, req)
}




func (rc *RemoteContext) Abort() {
	rc.stepLock.Lock()
	defer rc.stepLock.Unlock()

	rc.submitLock.Lock()
	defer rc.submitLock.Unlock()

	if rc.cancelStep != nil {
		rc.cancelStep()
		rc.cancelStep = nil
	}

	rc.closeSubmitStream()
	rc.onAbort()
}



func (rc *RemoteContext) closeSubmitStream() {
	if rc.cancelSubmitStream != nil {
		rc.cancelSubmitStream()
		rc.cancelSubmitStream = nil
	}

	if rc.submitStream != nil {
		rc.submitStream.CloseSend()
		rc.submitStream = nil
	}
}
