













package xfs








type Stats struct {
	
	
	
	Name string

	ExtentAllocation   ExtentAllocationStats
	AllocationBTree    BTreeStats
	BlockMapping       BlockMappingStats
	BlockMapBTree      BTreeStats
	DirectoryOperation DirectoryOperationStats
	Transaction        TransactionStats
	InodeOperation     InodeOperationStats
	LogOperation       LogOperationStats
	ReadWrite          ReadWriteStats
	AttributeOperation AttributeOperationStats
	InodeClustering    InodeClusteringStats
	Vnode              VnodeStats
	Buffer             BufferStats
	ExtendedPrecision  ExtendedPrecisionStats
}


type ExtentAllocationStats struct {
	ExtentsAllocated uint32
	BlocksAllocated  uint32
	ExtentsFreed     uint32
	BlocksFreed      uint32
}


type BTreeStats struct {
	Lookups         uint32
	Compares        uint32
	RecordsInserted uint32
	RecordsDeleted  uint32
}


type BlockMappingStats struct {
	Reads                uint32
	Writes               uint32
	Unmaps               uint32
	ExtentListInsertions uint32
	ExtentListDeletions  uint32
	ExtentListLookups    uint32
	ExtentListCompares   uint32
}


type DirectoryOperationStats struct {
	Lookups  uint32
	Creates  uint32
	Removes  uint32
	Getdents uint32
}


type TransactionStats struct {
	Sync  uint32
	Async uint32
	Empty uint32
}


type InodeOperationStats struct {
	Attempts        uint32
	Found           uint32
	Recycle         uint32
	Missed          uint32
	Duplicate       uint32
	Reclaims        uint32
	AttributeChange uint32
}


type LogOperationStats struct {
	Writes            uint32
	Blocks            uint32
	NoInternalBuffers uint32
	Force             uint32
	ForceSleep        uint32
}



type ReadWriteStats struct {
	Read  uint32
	Write uint32
}



type AttributeOperationStats struct {
	Get    uint32
	Set    uint32
	Remove uint32
	List   uint32
}



type InodeClusteringStats struct {
	Iflush     uint32
	Flush      uint32
	FlushInode uint32
}


type VnodeStats struct {
	Active   uint32
	Allocate uint32
	Get      uint32
	Hold     uint32
	Release  uint32
	Reclaim  uint32
	Remove   uint32
	Free     uint32
}


type BufferStats struct {
	Get             uint32
	Create          uint32
	GetLocked       uint32
	GetLockedWaited uint32
	BusyLocked      uint32
	MissLocked      uint32
	PageRetries     uint32
	PageFound       uint32
	GetRead         uint32
}



type ExtendedPrecisionStats struct {
	FlushBytes uint64
	WriteBytes uint64
	ReadBytes  uint64
}
