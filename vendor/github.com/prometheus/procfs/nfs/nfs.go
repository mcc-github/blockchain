














package nfs


type ReplyCache struct {
	Hits    uint64
	Misses  uint64
	NoCache uint64
}


type FileHandles struct {
	Stale        uint64
	TotalLookups uint64
	AnonLookups  uint64
	DirNoCache   uint64
	NoDirNoCache uint64
}


type InputOutput struct {
	Read  uint64
	Write uint64
}


type Threads struct {
	Threads uint64
	FullCnt uint64
}


type ReadAheadCache struct {
	CacheSize      uint64
	CacheHistogram []uint64
	NotFound       uint64
}


type Network struct {
	NetCount   uint64
	UDPCount   uint64
	TCPCount   uint64
	TCPConnect uint64
}


type ClientRPC struct {
	RPCCount        uint64
	Retransmissions uint64
	AuthRefreshes   uint64
}


type ServerRPC struct {
	RPCCount uint64
	BadCnt   uint64
	BadFmt   uint64
	BadAuth  uint64
	BadcInt  uint64
}


type V2Stats struct {
	Null     uint64
	GetAttr  uint64
	SetAttr  uint64
	Root     uint64
	Lookup   uint64
	ReadLink uint64
	Read     uint64
	WrCache  uint64
	Write    uint64
	Create   uint64
	Remove   uint64
	Rename   uint64
	Link     uint64
	SymLink  uint64
	MkDir    uint64
	RmDir    uint64
	ReadDir  uint64
	FsStat   uint64
}


type V3Stats struct {
	Null        uint64
	GetAttr     uint64
	SetAttr     uint64
	Lookup      uint64
	Access      uint64
	ReadLink    uint64
	Read        uint64
	Write       uint64
	Create      uint64
	MkDir       uint64
	SymLink     uint64
	MkNod       uint64
	Remove      uint64
	RmDir       uint64
	Rename      uint64
	Link        uint64
	ReadDir     uint64
	ReadDirPlus uint64
	FsStat      uint64
	FsInfo      uint64
	PathConf    uint64
	Commit      uint64
}


type ClientV4Stats struct {
	Null               uint64
	Read               uint64
	Write              uint64
	Commit             uint64
	Open               uint64
	OpenConfirm        uint64
	OpenNoattr         uint64
	OpenDowngrade      uint64
	Close              uint64
	Setattr            uint64
	FsInfo             uint64
	Renew              uint64
	SetClientId        uint64
	SetClientIdConfirm uint64
	Lock               uint64
	Lockt              uint64
	Locku              uint64
	Access             uint64
	Getattr            uint64
	Lookup             uint64
	LookupRoot         uint64
	Remove             uint64
	Rename             uint64
	Link               uint64
	Symlink            uint64
	Create             uint64
	Pathconf           uint64
	StatFs             uint64
	ReadLink           uint64
	ReadDir            uint64
	ServerCaps         uint64
	DelegReturn        uint64
	GetAcl             uint64
	SetAcl             uint64
	FsLocations        uint64
	ReleaseLockowner   uint64
	Secinfo            uint64
	FsidPresent        uint64
	ExchangeId         uint64
	CreateSession      uint64
	DestroySession     uint64
	Sequence           uint64
	GetLeaseTime       uint64
	ReclaimComplete    uint64
	LayoutGet          uint64
	GetDeviceInfo      uint64
	LayoutCommit       uint64
	LayoutReturn       uint64
	SecinfoNoName      uint64
	TestStateId        uint64
	FreeStateId        uint64
	GetDeviceList      uint64
	BindConnToSession  uint64
	DestroyClientId    uint64
	Seek               uint64
	Allocate           uint64
	DeAllocate         uint64
	LayoutStats        uint64
	Clone              uint64
}


type ServerV4Stats struct {
	Null     uint64
	Compound uint64
}






type V4Ops struct {
	
	Op0Unused    uint64
	Op1Unused    uint64
	Op2Future    uint64
	Access       uint64
	Close        uint64
	Commit       uint64
	Create       uint64
	DelegPurge   uint64
	DelegReturn  uint64
	GetAttr      uint64
	GetFH        uint64
	Link         uint64
	Lock         uint64
	Lockt        uint64
	Locku        uint64
	Lookup       uint64
	LookupRoot   uint64
	Nverify      uint64
	Open         uint64
	OpenAttr     uint64
	OpenConfirm  uint64
	OpenDgrd     uint64
	PutFH        uint64
	PutPubFH     uint64
	PutRootFH    uint64
	Read         uint64
	ReadDir      uint64
	ReadLink     uint64
	Remove       uint64
	Rename       uint64
	Renew        uint64
	RestoreFH    uint64
	SaveFH       uint64
	SecInfo      uint64
	SetAttr      uint64
	Verify       uint64
	Write        uint64
	RelLockOwner uint64
}


type ClientRPCStats struct {
	Network       Network
	ClientRPC     ClientRPC
	V2Stats       V2Stats
	V3Stats       V3Stats
	ClientV4Stats ClientV4Stats
}


type ServerRPCStats struct {
	ReplyCache     ReplyCache
	FileHandles    FileHandles
	InputOutput    InputOutput
	Threads        Threads
	ReadAheadCache ReadAheadCache
	Network        Network
	ServerRPC      ServerRPC
	V2Stats        V2Stats
	V3Stats        V3Stats
	ServerV4Stats  ServerV4Stats
	V4Ops          V4Ops
}
