













package base

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)


type PickerBuilder interface {
	
	
	Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker
}



func NewBalancerBuilder(name string, pb PickerBuilder) balancer.Builder {
	return NewBalancerBuilderWithConfig(name, pb, Config{})
}


type Config struct {
	
	HealthCheck bool
}


func NewBalancerBuilderWithConfig(name string, pb PickerBuilder, config Config) balancer.Builder {
	return &baseBuilder{
		name:          name,
		pickerBuilder: pb,
		config:        config,
	}
}
