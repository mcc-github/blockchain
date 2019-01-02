




package internal

import "context"

var (
	
	WithContextDialer interface{} 
	
	WithResolverBuilder interface{} 
	
	HealthCheckFunc func(ctx context.Context, newStream func() (interface{}, error), reportHealth func(bool), serviceName string) error
)

const (
	
	CredsBundleModeFallback = "fallback"
	
	
	CredsBundleModeBalancer = "balancer"
	
	
	CredsBundleModeBackendFromBalancer = "backend-from-balancer"
)
