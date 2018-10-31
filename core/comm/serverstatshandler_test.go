/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/metrics/metricsfakes"
	"github.com/mcc-github/blockchain/core/comm"
	testpb "github.com/mcc-github/blockchain/core/comm/testdata/grpc"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func TestConnectionCounters(t *testing.T) {
	t.Parallel()
	gt := NewGomegaWithT(t)

	openConn := &metricsfakes.Counter{}
	closedConn := &metricsfakes.Counter{}
	sh := &comm.ServerStatsHandler{
		OpenConnCounter:   openConn,
		ClosedConnCounter: closedConn,
	}

	for i := 1; i <= 10; i++ {
		sh.HandleConn(context.Background(), &stats.ConnBegin{})
		gt.Expect(openConn.AddCallCount()).To(Equal(i))
	}

	for i := 1; i <= 5; i++ {
		sh.HandleConn(context.Background(), &stats.ConnEnd{})
		gt.Expect(closedConn.AddCallCount()).To(Equal(i))
	}
}

func TestConnMetricsGRPCServer(t *testing.T) {
	t.Parallel()
	gt := NewGomegaWithT(t)

	openConn := &metricsfakes.Counter{}
	closedConn := &metricsfakes.Counter{}

	listener, err := net.Listen("tcp", "localhost:0")
	gt.Expect(err).NotTo(HaveOccurred())
	srv, err := comm.NewGRPCServerFromListener(
		listener,
		comm.ServerConfig{
			SecOpts: &comm.SecureOptions{UseTLS: false},
			Metrics: &comm.Metrics{
				OpenConnCounter:   openConn,
				ClosedConnCounter: closedConn,
			},
		},
	)
	gt.Expect(err).NotTo(HaveOccurred())

	
	testpb.RegisterEmptyServiceServer(srv.Server(), &emptyServiceServer{})

	
	go srv.Start()
	defer srv.Stop()

	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	gt.Expect(openConn.AddCallCount()).To(Equal(0))
	gt.Expect(closedConn.AddCallCount()).To(Equal(0))

	
	var clientConns []*grpc.ClientConn
	for i := 1; i <= 3; i++ {
		clientConn, err := grpc.DialContext(ctx, listener.Addr().String(), grpc.WithInsecure())
		gt.Expect(err).NotTo(HaveOccurred())
		clientConns = append(clientConns, clientConn)

		
		client := testpb.NewEmptyServiceClient(clientConn)
		_, err = client.EmptyCall(context.Background(), &testpb.Empty{})
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(openConn.AddCallCount()).To(Equal(i))
	}

	for i, conn := range clientConns {
		gt.Expect(closedConn.AddCallCount()).Should(Equal(i))
		conn.Close()
		gt.Eventually(closedConn.AddCallCount, time.Second).Should(Equal(i + 1))
	}
}
