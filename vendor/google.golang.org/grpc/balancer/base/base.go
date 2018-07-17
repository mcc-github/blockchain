













package base

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)


type PickerBuilder interface {
	
	
	Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker
}



func NewBalancerBuilder(name string, pb PickerBuilder) balancer.Builder {
	return &baseBuilder{
		name:          name,
		pickerBuilder: pb,
	}
}
