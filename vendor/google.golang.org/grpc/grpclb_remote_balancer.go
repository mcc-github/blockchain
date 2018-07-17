

package grpc

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	lbpb "google.golang.org/grpc/grpclb/grpc_lb_v1/messages"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
)



func (lb *lbBalancer) processServerList(l *lbpb.ServerList) {
	grpclog.Infof("lbBalancer: processing server list: %+v", l)
	lb.mu.Lock()
	defer lb.mu.Unlock()

	
	
	lb.serverListReceived = true

	
	if reflect.DeepEqual(lb.fullServerList, l.Servers) {
		grpclog.Infof("lbBalancer: new serverlist same as the previous one, ignoring")
		return
	}
	lb.fullServerList = l.Servers

	var backendAddrs []resolver.Address
	for _, s := range l.Servers {
		if s.DropForLoadBalancing || s.DropForRateLimiting {
			continue
		}

		md := metadata.Pairs(lbTokeyKey, s.LoadBalanceToken)
		ip := net.IP(s.IpAddress)
		ipStr := ip.String()
		if ip.To4() == nil {
			
			
			ipStr = fmt.Sprintf("[%s]", ipStr)
		}
		addr := resolver.Address{
			Addr:     fmt.Sprintf("%s:%d", ipStr, s.Port),
			Metadata: &md,
		}

		backendAddrs = append(backendAddrs, addr)
	}

	
	backendsUpdated := lb.refreshSubConns(backendAddrs)
	
	
	
	
	if !backendsUpdated {
		lb.regeneratePicker()
		lb.cc.UpdateBalancerState(lb.state, lb.picker)
	}
}





func (lb *lbBalancer) refreshSubConns(backendAddrs []resolver.Address) bool {
	lb.backendAddrs = nil
	var backendsUpdated bool
	
	
	addrsSet := make(map[resolver.Address]struct{})
	
	for _, addr := range backendAddrs {
		addrWithoutMD := addr
		addrWithoutMD.Metadata = nil
		addrsSet[addrWithoutMD] = struct{}{}
		lb.backendAddrs = append(lb.backendAddrs, addrWithoutMD)

		if _, ok := lb.subConns[addrWithoutMD]; !ok {
			backendsUpdated = true

			
			sc, err := lb.cc.NewSubConn([]resolver.Address{addr}, balancer.NewSubConnOptions{})
			if err != nil {
				grpclog.Warningf("roundrobinBalancer: failed to create new SubConn: %v", err)
				continue
			}
			lb.subConns[addrWithoutMD] = sc 
			lb.scStates[sc] = connectivity.Idle
			sc.Connect()
		}
	}

	for a, sc := range lb.subConns {
		
		if _, ok := addrsSet[a]; !ok {
			backendsUpdated = true

			lb.cc.RemoveSubConn(sc)
			delete(lb.subConns, a)
			
			
		}
	}

	return backendsUpdated
}

func (lb *lbBalancer) readServerList(s *balanceLoadClientStream) error {
	for {
		reply, err := s.Recv()
		if err != nil {
			return fmt.Errorf("grpclb: failed to recv server list: %v", err)
		}
		if serverList := reply.GetServerList(); serverList != nil {
			lb.processServerList(serverList)
		}
	}
}

func (lb *lbBalancer) sendLoadReport(s *balanceLoadClientStream, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
		case <-s.Context().Done():
			return
		}
		stats := lb.clientStats.toClientStats()
		t := time.Now()
		stats.Timestamp = &lbpb.Timestamp{
			Seconds: t.Unix(),
			Nanos:   int32(t.Nanosecond()),
		}
		if err := s.Send(&lbpb.LoadBalanceRequest{
			LoadBalanceRequestType: &lbpb.LoadBalanceRequest_ClientStats{
				ClientStats: stats,
			},
		}); err != nil {
			return
		}
	}
}
func (lb *lbBalancer) callRemoteBalancer() error {
	lbClient := &loadBalancerClient{cc: lb.ccRemoteLB}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := lbClient.BalanceLoad(ctx, FailFast(false))
	if err != nil {
		return fmt.Errorf("grpclb: failed to perform RPC to the remote balancer %v", err)
	}

	
	initReq := &lbpb.LoadBalanceRequest{
		LoadBalanceRequestType: &lbpb.LoadBalanceRequest_InitialRequest{
			InitialRequest: &lbpb.InitialLoadBalanceRequest{
				Name: lb.target,
			},
		},
	}
	if err := stream.Send(initReq); err != nil {
		return fmt.Errorf("grpclb: failed to send init request: %v", err)
	}
	reply, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("grpclb: failed to recv init response: %v", err)
	}
	initResp := reply.GetInitialResponse()
	if initResp == nil {
		return fmt.Errorf("grpclb: reply from remote balancer did not include initial response")
	}
	if initResp.LoadBalancerDelegate != "" {
		return fmt.Errorf("grpclb: Delegation is not supported")
	}

	go func() {
		if d := convertDuration(initResp.ClientStatsReportInterval); d > 0 {
			lb.sendLoadReport(stream, d)
		}
	}()
	return lb.readServerList(stream)
}

func (lb *lbBalancer) watchRemoteBalancer() {
	for {
		err := lb.callRemoteBalancer()
		select {
		case <-lb.doneCh:
			return
		default:
			if err != nil {
				grpclog.Error(err)
			}
		}

	}
}

func (lb *lbBalancer) dialRemoteLB(remoteLBName string) {
	var dopts []DialOption
	if creds := lb.opt.DialCreds; creds != nil {
		if err := creds.OverrideServerName(remoteLBName); err == nil {
			dopts = append(dopts, WithTransportCredentials(creds))
		} else {
			grpclog.Warningf("grpclb: failed to override the server name in the credentials: %v, using Insecure", err)
			dopts = append(dopts, WithInsecure())
		}
	} else {
		dopts = append(dopts, WithInsecure())
	}
	if lb.opt.Dialer != nil {
		
		
		dopts = append(dopts, withContextDialer(lb.opt.Dialer))
	}
	
	dopts = append(dopts, WithBalancerName(PickFirstBalancerName))
	dopts = append(dopts, withResolverBuilder(lb.manualResolver))
	
	
	cc, err := Dial("grpclb:/
	if err != nil {
		grpclog.Fatalf("failed to dial: %v", err)
	}
	lb.ccRemoteLB = cc
	go lb.watchRemoteBalancer()
}
