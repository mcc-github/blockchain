/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package accesscontrol

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/flogging"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ccSrv struct {
	l              net.Listener
	grpcSrv        *grpc.Server
	t              *testing.T
	cert           []byte
	expectedCCname string
}

func (cs *ccSrv) Register(stream pb.ChaincodeSupport_RegisterServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	
	assert.Equal(cs.t, pb.ChaincodeMessage_REGISTER.String(), msg.Type.String())
	
	chaincodeID := &pb.ChaincodeID{}
	err = proto.Unmarshal(msg.Payload, chaincodeID)
	if err != nil {
		return err
	}
	assert.Equal(cs.t, cs.expectedCCname, chaincodeID.Name)
	
	for {
		msg, _ = stream.Recv()
		if err != nil {
			return err
		}
		err = stream.Send(msg)
		if err != nil {
			return err
		}
	}
}

func (cs *ccSrv) stop() {
	cs.grpcSrv.Stop()
	cs.l.Close()
}

func createTLSService(t *testing.T, ca tlsgen.CA, host string) *grpc.Server {
	keyPair, err := ca.NewServerCertKeyPair(host)
	cert, err := tls.X509KeyPair(keyPair.Cert, keyPair.Key)
	assert.NoError(t, err)
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    x509.NewCertPool(),
	}
	tlsConf.ClientCAs.AppendCertsFromPEM(ca.CertBytes())
	return grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConf)))
}

func newCCServer(t *testing.T, port int, expectedCCname string, withTLS bool, ca tlsgen.CA) *ccSrv {
	var s *grpc.Server
	if withTLS {
		s = createTLSService(t, ca, "localhost")
	} else {
		s = grpc.NewServer()
	}

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "", port))
	assert.NoError(t, err, "%v", err)
	return &ccSrv{
		expectedCCname: expectedCCname,
		l:              l,
		grpcSrv:        s,
	}
}

type ccClient struct {
	conn   *grpc.ClientConn
	stream pb.ChaincodeSupport_RegisterClient
}

func newClient(t *testing.T, port int, cert *tls.Certificate, peerCACert []byte) (*ccClient, error) {
	tlsCfg := &tls.Config{
		RootCAs: x509.NewCertPool(),
	}

	tlsCfg.RootCAs.AppendCertsFromPEM(peerCACert)
	if cert != nil {
		tlsCfg.Certificates = []tls.Certificate{*cert}
	}
	tlsOpts := grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second)
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("localhost:%d", port), tlsOpts, grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	chaincodeSupportClient := pb.NewChaincodeSupportClient(conn)
	stream, err := chaincodeSupportClient.Register(context.Background())
	assert.NoError(t, err)
	return &ccClient{
		conn:   conn,
		stream: stream,
	}, nil
}

func (c *ccClient) close() {
	c.conn.Close()
}

func (c *ccClient) sendMsg(msg *pb.ChaincodeMessage) {
	c.stream.Send(msg)
}

func (c *ccClient) recv() *pb.ChaincodeMessage {
	msgs := make(chan *pb.ChaincodeMessage, 1)
	go func() {
		msg, _ := c.stream.Recv()
		if msg != nil {
			msgs <- msg
		}
	}()
	select {
	case <-time.After(time.Second):
		return nil
	case msg := <-msgs:
		return msg
	}
}

func TestAccessControl(t *testing.T) {
	backupTTL := ttl
	defer func() {
		ttl = backupTTL
	}()
	ttl = time.Second * 3

	logAsserter := &logBackend{
		logEntries: make(chan string, 1),
	}
	logger.SetBackend(logAsserter)
	defer func() {
		logger = flogging.MustGetLogger("accessControl")
	}()

	chaincodeID := &pb.ChaincodeID{Name: "example02"}
	payload, err := proto.Marshal(chaincodeID)
	registerMsg := &pb.ChaincodeMessage{
		Type:    pb.ChaincodeMessage_REGISTER,
		Payload: payload,
	}
	putStateMsg := &pb.ChaincodeMessage{
		Type: pb.ChaincodeMessage_PUT_STATE,
	}

	ca, _ := tlsgen.NewCA()
	srv := newCCServer(t, 7052, "example02", true, ca)
	auth := NewAuthenticator(ca)
	pb.RegisterChaincodeSupportServer(srv.grpcSrv, auth.Wrap(srv))
	go srv.grpcSrv.Serve(srv.l)
	defer srv.stop()

	
	_, err = newClient(t, 7052, nil, ca.CertBytes())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")

	
	maliciousCA, _ := tlsgen.NewCA()
	keyPair, err := maliciousCA.NewClientCertKeyPair()
	cert, err := tls.X509KeyPair(keyPair.Cert, keyPair.Key)
	assert.NoError(t, err)
	_, err = newClient(t, 7052, &cert, ca.CertBytes())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")

	
	kp, err := auth.Generate("example01")
	assert.NoError(t, err)
	keyBytes, err := base64.StdEncoding.DecodeString(kp.Key)
	assert.NoError(t, err)
	certBytes, err := base64.StdEncoding.DecodeString(kp.Cert)
	assert.NoError(t, err)
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	assert.NoError(t, err)
	mismatchedShim, err := newClient(t, 7052, &cert, ca.CertBytes())
	assert.NoError(t, err)
	defer mismatchedShim.close()
	mismatchedShim.sendMsg(registerMsg)
	mismatchedShim.sendMsg(putStateMsg)
	
	assert.Nil(t, mismatchedShim.recv())
	logAsserter.assertLastLogContains(t, "with given certificate hash", "belongs to a different chaincode")

	
	kp, err = auth.Generate("example02")
	assert.NoError(t, err)
	keyBytes, err = base64.StdEncoding.DecodeString(kp.Key)
	assert.NoError(t, err)
	certBytes, err = base64.StdEncoding.DecodeString(kp.Cert)
	assert.NoError(t, err)
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	assert.NoError(t, err)
	realCC, err := newClient(t, 7052, &cert, ca.CertBytes())
	assert.NoError(t, err)
	defer realCC.close()
	realCC.sendMsg(registerMsg)
	realCC.sendMsg(putStateMsg)
	echoMsg := realCC.recv()
	
	assert.NotNil(t, echoMsg)
	assert.Equal(t, pb.ChaincodeMessage_PUT_STATE, echoMsg.Type)
	
	assert.Empty(t, logAsserter.logEntries)

	
	
	
	
	kp, err = auth.Generate("example02")
	assert.NoError(t, err)
	keyBytes, err = base64.StdEncoding.DecodeString(kp.Key)
	assert.NoError(t, err)
	certBytes, err = base64.StdEncoding.DecodeString(kp.Cert)
	assert.NoError(t, err)
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	assert.NoError(t, err)
	confusedCC, err := newClient(t, 7052, &cert, ca.CertBytes())
	assert.NoError(t, err)
	defer confusedCC.close()
	confusedCC.sendMsg(putStateMsg)
	confusedCC.sendMsg(registerMsg)
	confusedCC.sendMsg(putStateMsg)
	assert.Nil(t, confusedCC.recv())
	logAsserter.assertLastLogContains(t, "expected a ChaincodeMessage_REGISTER message")

	
	
	kp, err = auth.Generate("example02")
	assert.NoError(t, err)
	keyBytes, err = base64.StdEncoding.DecodeString(kp.Key)
	assert.NoError(t, err)
	certBytes, err = base64.StdEncoding.DecodeString(kp.Cert)
	assert.NoError(t, err)
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	assert.NoError(t, err)
	malformedMessageCC, err := newClient(t, 7052, &cert, ca.CertBytes())
	assert.NoError(t, err)
	defer malformedMessageCC.close()
	
	originalPayload := registerMsg.Payload
	registerMsg.Payload = append(registerMsg.Payload, 0)
	malformedMessageCC.sendMsg(registerMsg)
	malformedMessageCC.sendMsg(putStateMsg)
	assert.Nil(t, malformedMessageCC.recv())
	logAsserter.assertLastLogContains(t, "Failed unmarshaling message")
	
	registerMsg.Payload = originalPayload

	
	
	
	
	
	kp, err = auth.Generate("example02")
	assert.NoError(t, err)
	keyBytes, err = base64.StdEncoding.DecodeString(kp.Key)
	assert.NoError(t, err)
	certBytes, err = base64.StdEncoding.DecodeString(kp.Cert)
	assert.NoError(t, err)
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	assert.NoError(t, err)
	lateCC, err := newClient(t, 7052, &cert, ca.CertBytes())
	assert.NoError(t, err)
	defer realCC.close()
	time.Sleep(ttl + time.Second*2)
	lateCC.sendMsg(registerMsg)
	lateCC.sendMsg(putStateMsg)
	echoMsg = lateCC.recv()
	assert.Nil(t, echoMsg)
	logAsserter.assertLastLogContains(t, "with given certificate hash", "not found in registry")
}

type logBackend struct {
	logEntries chan string
}

func (l *logBackend) assertLastLogContains(t *testing.T, ss ...string) {
	lastLogMsg := <-l.logEntries
	for _, s := range ss {
		assert.Contains(t, lastLogMsg, s)
	}
}

func (l *logBackend) Log(lvl logging.Level, n int, r *logging.Record) error {
	if lvl.String() != logging.WARNING.String() {
		return nil
	}
	l.logEntries <- fmt.Sprint(r.Args)
	return nil
}

func (*logBackend) GetLevel(string) logging.Level {
	return logging.DEBUG
}

func (*logBackend) SetLevel(logging.Level, string) {
	panic("implement me")
}

func (*logBackend) IsEnabledFor(logging.Level, string) bool {
	return true
}
