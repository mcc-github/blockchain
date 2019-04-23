




package internal

import (
	"context"
	"time"
)

var (
	
	WithResolverBuilder interface{} 
	
	WithHealthCheckFunc interface{} 
	
	HealthCheckFunc HealthChecker
	
	BalancerUnregister func(name string)
	
	
	KeepaliveMinPingTime = 10 * time.Second
)


type HealthChecker func(ctx context.Context, newStream func() (interface{}, error), reportHealth func(bool), serviceName string) error

const (
	
	CredsBundleModeFallback = "fallback"
	
	
	CredsBundleModeBalancer = "balancer"
	
	
	CredsBundleModeBackendFromBalancer = "backend-from-balancer"
)
