package sarama

import "errors"





type ClusterAdmin interface {
	
	
	
	
	CreateTopic(topic string, detail *TopicDetail, validateOnly bool) error

	
	
	
	
	
	
	DeleteTopic(topic string) error

	
	
	
	
	
	
	CreatePartitions(topic string, count int32, assignment [][]int32, validateOnly bool) error

	
	
	DeleteRecords(topic string, partitionOffsets map[int32]int64) error

	
	
	
	
	
	
	
	DescribeConfig(resource ConfigResource) ([]ConfigEntry, error)

	
	
	
	
	
	AlterConfig(resourceType ConfigResourceType, name string, entries map[string]*string, validateOnly bool) error

	
	
	
	
	CreateACL(resource Resource, acl Acl) error

	
	
	
	ListAcls(filter AclFilter) ([]ResourceAcls, error)

	
	
	
	DeleteACL(filter AclFilter, validateOnly bool) ([]MatchingAcl, error)

	
	Close() error
}

type clusterAdmin struct {
	client Client
	conf   *Config
}


func NewClusterAdmin(addrs []string, conf *Config) (ClusterAdmin, error) {
	client, err := NewClient(addrs, conf)
	if err != nil {
		return nil, err
	}

	
	_, err = client.Controller()
	if err != nil {
		return nil, err
	}

	ca := &clusterAdmin{
		client: client,
		conf:   client.Config(),
	}
	return ca, nil
}

func (ca *clusterAdmin) Close() error {
	return ca.client.Close()
}

func (ca *clusterAdmin) Controller() (*Broker, error) {
	return ca.client.Controller()
}

func (ca *clusterAdmin) CreateTopic(topic string, detail *TopicDetail, validateOnly bool) error {

	if topic == "" {
		return ErrInvalidTopic
	}

	if detail == nil {
		return errors.New("You must specify topic details")
	}

	topicDetails := make(map[string]*TopicDetail)
	topicDetails[topic] = detail

	request := &CreateTopicsRequest{
		TopicDetails: topicDetails,
		ValidateOnly: validateOnly,
		Timeout:      ca.conf.Admin.Timeout,
	}

	if ca.conf.Version.IsAtLeast(V0_11_0_0) {
		request.Version = 1
	}
	if ca.conf.Version.IsAtLeast(V1_0_0_0) {
		request.Version = 2
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.CreateTopics(request)
	if err != nil {
		return err
	}

	topicErr, ok := rsp.TopicErrors[topic]
	if !ok {
		return ErrIncompleteResponse
	}

	if topicErr.Err != ErrNoError {
		return topicErr.Err
	}

	return nil
}

func (ca *clusterAdmin) DeleteTopic(topic string) error {

	if topic == "" {
		return ErrInvalidTopic
	}

	request := &DeleteTopicsRequest{
		Topics:  []string{topic},
		Timeout: ca.conf.Admin.Timeout,
	}

	if ca.conf.Version.IsAtLeast(V0_11_0_0) {
		request.Version = 1
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.DeleteTopics(request)
	if err != nil {
		return err
	}

	topicErr, ok := rsp.TopicErrorCodes[topic]
	if !ok {
		return ErrIncompleteResponse
	}

	if topicErr != ErrNoError {
		return topicErr
	}
	return nil
}

func (ca *clusterAdmin) CreatePartitions(topic string, count int32, assignment [][]int32, validateOnly bool) error {
	if topic == "" {
		return ErrInvalidTopic
	}

	topicPartitions := make(map[string]*TopicPartition)
	topicPartitions[topic] = &TopicPartition{Count: count, Assignment: assignment}

	request := &CreatePartitionsRequest{
		TopicPartitions: topicPartitions,
		Timeout:         ca.conf.Admin.Timeout,
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.CreatePartitions(request)
	if err != nil {
		return err
	}

	topicErr, ok := rsp.TopicPartitionErrors[topic]
	if !ok {
		return ErrIncompleteResponse
	}

	if topicErr.Err != ErrNoError {
		return topicErr.Err
	}

	return nil
}

func (ca *clusterAdmin) DeleteRecords(topic string, partitionOffsets map[int32]int64) error {

	if topic == "" {
		return ErrInvalidTopic
	}

	topics := make(map[string]*DeleteRecordsRequestTopic)
	topics[topic] = &DeleteRecordsRequestTopic{PartitionOffsets: partitionOffsets}
	request := &DeleteRecordsRequest{
		Topics:  topics,
		Timeout: ca.conf.Admin.Timeout,
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.DeleteRecords(request)
	if err != nil {
		return err
	}

	_, ok := rsp.Topics[topic]
	if !ok {
		return ErrIncompleteResponse
	}

	
	
	return nil
}

func (ca *clusterAdmin) DescribeConfig(resource ConfigResource) ([]ConfigEntry, error) {

	var entries []ConfigEntry
	var resources []*ConfigResource
	resources = append(resources, &resource)

	request := &DescribeConfigsRequest{
		Resources: resources,
	}

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	rsp, err := b.DescribeConfigs(request)
	if err != nil {
		return nil, err
	}

	for _, rspResource := range rsp.Resources {
		if rspResource.Name == resource.Name {
			if rspResource.ErrorMsg != "" {
				return nil, errors.New(rspResource.ErrorMsg)
			}
			for _, cfgEntry := range rspResource.Configs {
				entries = append(entries, *cfgEntry)
			}
		}
	}
	return entries, nil
}

func (ca *clusterAdmin) AlterConfig(resourceType ConfigResourceType, name string, entries map[string]*string, validateOnly bool) error {

	var resources []*AlterConfigsResource
	resources = append(resources, &AlterConfigsResource{
		Type:          resourceType,
		Name:          name,
		ConfigEntries: entries,
	})

	request := &AlterConfigsRequest{
		Resources:    resources,
		ValidateOnly: validateOnly,
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.AlterConfigs(request)
	if err != nil {
		return err
	}

	for _, rspResource := range rsp.Resources {
		if rspResource.Name == name {
			if rspResource.ErrorMsg != "" {
				return errors.New(rspResource.ErrorMsg)
			}
		}
	}
	return nil
}

func (ca *clusterAdmin) CreateACL(resource Resource, acl Acl) error {
	var acls []*AclCreation
	acls = append(acls, &AclCreation{resource, acl})
	request := &CreateAclsRequest{AclCreations: acls}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	_, err = b.CreateAcls(request)
	return err
}

func (ca *clusterAdmin) ListAcls(filter AclFilter) ([]ResourceAcls, error) {

	request := &DescribeAclsRequest{AclFilter: filter}

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	rsp, err := b.DescribeAcls(request)
	if err != nil {
		return nil, err
	}

	var lAcls []ResourceAcls
	for _, rAcl := range rsp.ResourceAcls {
		lAcls = append(lAcls, *rAcl)
	}
	return lAcls, nil
}

func (ca *clusterAdmin) DeleteACL(filter AclFilter, validateOnly bool) ([]MatchingAcl, error) {
	var filters []*AclFilter
	filters = append(filters, &filter)
	request := &DeleteAclsRequest{Filters: filters}

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	rsp, err := b.DeleteAcls(request)
	if err != nil {
		return nil, err
	}

	var mAcls []MatchingAcl
	for _, fr := range rsp.FilterResponses {
		for _, mACL := range fr.MatchingAcls {
			mAcls = append(mAcls, *mACL)
		}

	}
	return mAcls, nil
}
