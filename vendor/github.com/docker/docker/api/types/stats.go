

package types 

import "time"



type ThrottlingData struct {
	
	Periods uint64 `json:"periods"`
	
	ThrottledPeriods uint64 `json:"throttled_periods"`
	
	ThrottledTime uint64 `json:"throttled_time"`
}


type CPUUsage struct {
	
	
	
	TotalUsage uint64 `json:"total_usage"`

	
	
	PercpuUsage []uint64 `json:"percpu_usage,omitempty"`

	
	
	
	
	UsageInKernelmode uint64 `json:"usage_in_kernelmode"`

	
	
	
	
	UsageInUsermode uint64 `json:"usage_in_usermode"`
}


type CPUStats struct {
	
	CPUUsage CPUUsage `json:"cpu_usage"`

	
	SystemUsage uint64 `json:"system_cpu_usage,omitempty"`

	
	OnlineCPUs uint32 `json:"online_cpus,omitempty"`

	
	ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
}



type MemoryStats struct {
	

	
	Usage uint64 `json:"usage,omitempty"`
	
	MaxUsage uint64 `json:"max_usage,omitempty"`
	
	
	Stats map[string]uint64 `json:"stats,omitempty"`
	
	Failcnt uint64 `json:"failcnt,omitempty"`
	Limit   uint64 `json:"limit,omitempty"`

	
	

	
	Commit uint64 `json:"commitbytes,omitempty"`
	
	CommitPeak uint64 `json:"commitpeakbytes,omitempty"`
	
	PrivateWorkingSet uint64 `json:"privateworkingset,omitempty"`
}



type BlkioStatEntry struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Op    string `json:"op"`
	Value uint64 `json:"value"`
}





type BlkioStats struct {
	
	IoServiceBytesRecursive []BlkioStatEntry `json:"io_service_bytes_recursive"`
	IoServicedRecursive     []BlkioStatEntry `json:"io_serviced_recursive"`
	IoQueuedRecursive       []BlkioStatEntry `json:"io_queue_recursive"`
	IoServiceTimeRecursive  []BlkioStatEntry `json:"io_service_time_recursive"`
	IoWaitTimeRecursive     []BlkioStatEntry `json:"io_wait_time_recursive"`
	IoMergedRecursive       []BlkioStatEntry `json:"io_merged_recursive"`
	IoTimeRecursive         []BlkioStatEntry `json:"io_time_recursive"`
	SectorsRecursive        []BlkioStatEntry `json:"sectors_recursive"`
}


type StorageStats struct {
	ReadCountNormalized  uint64 `json:"read_count_normalized,omitempty"`
	ReadSizeBytes        uint64 `json:"read_size_bytes,omitempty"`
	WriteCountNormalized uint64 `json:"write_count_normalized,omitempty"`
	WriteSizeBytes       uint64 `json:"write_size_bytes,omitempty"`
}


type NetworkStats struct {
	
	RxBytes uint64 `json:"rx_bytes"`
	
	RxPackets uint64 `json:"rx_packets"`
	
	
	RxErrors uint64 `json:"rx_errors"`
	
	RxDropped uint64 `json:"rx_dropped"`
	
	TxBytes uint64 `json:"tx_bytes"`
	
	TxPackets uint64 `json:"tx_packets"`
	
	
	TxErrors uint64 `json:"tx_errors"`
	
	TxDropped uint64 `json:"tx_dropped"`
	
	EndpointID string `json:"endpoint_id,omitempty"`
	
	InstanceID string `json:"instance_id,omitempty"`
}


type PidsStats struct {
	
	Current uint64 `json:"current,omitempty"`
	
	
	Limit uint64 `json:"limit,omitempty"`
}


type Stats struct {
	
	Read    time.Time `json:"read"`
	PreRead time.Time `json:"preread"`

	
	PidsStats  PidsStats  `json:"pids_stats,omitempty"`
	BlkioStats BlkioStats `json:"blkio_stats,omitempty"`

	
	NumProcs     uint32       `json:"num_procs"`
	StorageStats StorageStats `json:"storage_stats,omitempty"`

	
	CPUStats    CPUStats    `json:"cpu_stats,omitempty"`
	PreCPUStats CPUStats    `json:"precpu_stats,omitempty"` 
	MemoryStats MemoryStats `json:"memory_stats,omitempty"`
}


type StatsJSON struct {
	Stats

	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`

	
	Networks map[string]NetworkStats `json:"networks,omitempty"`
}
