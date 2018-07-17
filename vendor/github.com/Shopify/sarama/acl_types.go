package sarama

type AclOperation int


const (
	AclOperationUnknown         AclOperation = 0
	AclOperationAny             AclOperation = 1
	AclOperationAll             AclOperation = 2
	AclOperationRead            AclOperation = 3
	AclOperationWrite           AclOperation = 4
	AclOperationCreate          AclOperation = 5
	AclOperationDelete          AclOperation = 6
	AclOperationAlter           AclOperation = 7
	AclOperationDescribe        AclOperation = 8
	AclOperationClusterAction   AclOperation = 9
	AclOperationDescribeConfigs AclOperation = 10
	AclOperationAlterConfigs    AclOperation = 11
	AclOperationIdempotentWrite AclOperation = 12
)

type AclPermissionType int


const (
	AclPermissionUnknown AclPermissionType = 0
	AclPermissionAny     AclPermissionType = 1
	AclPermissionDeny    AclPermissionType = 2
	AclPermissionAllow   AclPermissionType = 3
)

type AclResourceType int


const (
	AclResourceUnknown         AclResourceType = 0
	AclResourceAny             AclResourceType = 1
	AclResourceTopic           AclResourceType = 2
	AclResourceGroup           AclResourceType = 3
	AclResourceCluster         AclResourceType = 4
	AclResourceTransactionalID AclResourceType = 5
)
