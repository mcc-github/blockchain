/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"context"

	"github.com/pkg/errors"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)


func NewHealthCheckClient(config ClientConfig, address, service string) (HealthCheckClient, error) {
	hcc := HealthCheckClient{
		Address: address,
		Service: service,
	}
	client, err := NewGRPCClient(config)
	if err != nil {
		return hcc, errors.Wrapf(err, "failed to create health check client")
	}
	hcc.Client = client
	return hcc, nil
}





type HealthCheckClient struct {
	Client  *GRPCClient
	Address string
	Service string
}




func (hcc HealthCheckClient) HealthCheck(ctx context.Context) error {
	conn, err := hcc.Client.NewConnection(hcc.Address, "")
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to connect to service '%s' at '%s'",
			hcc.Service,
			hcc.Address)
	}
	defer conn.Close()
	h := healthpb.NewHealthClient(conn)
	req := &healthpb.HealthCheckRequest{
		Service: hcc.Service,
	}
	res, err := h.Check(ctx, req)
	if res.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		return errors.Wrapf(
			err,
			"failed to connect to service '%s' at '%s'",
			hcc.Service,
			hcc.Address)
	}
	return nil
}
