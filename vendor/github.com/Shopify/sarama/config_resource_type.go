package sarama

type ConfigResourceType int8




const (
	UnknownResource ConfigResourceType = 0
	AnyResource     ConfigResourceType = 1
	TopicResource   ConfigResourceType = 2
	GroupResource   ConfigResourceType = 3
	ClusterResource ConfigResourceType = 4
	BrokerResource  ConfigResourceType = 5
)
