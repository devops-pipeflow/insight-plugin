package proto

type Config struct {
	EnvVariables   []EnvVariable  `json:"envVariables"`
	BuildConfig    BuildConfig    `json:"buildConfig"`
	CodeConfig     CodeConfig     `json:"codeConfig"`
	NodeConfig     NodeConfig     `json:"nodeConfig"`
	ArtifactConfig ArtifactConfig `json:"artifactConfig"`
	GptConfig      GptConfig      `json:"gptConfig"`
	RepoConfig     RepoConfig     `json:"repoConfig"`
	ReviewConfig   ReviewConfig   `json:"reviewConfig"`
}

type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BuildConfig struct {
	LoggingConfig LoggingConfig `json:"loggingConfig"`
}

type CodeConfig struct{}

type NodeConfig struct {
	Duration string `json:"duration"`
}

type ArtifactConfig struct {
	Url  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type GptConfig struct{}

type RepoConfig struct {
	Url  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type ReviewConfig struct {
	Url  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoggingConfig struct {
	Start int64 `json:"start"`
	Len   int64 `json:"len"`
	Count int64 `json:"count"`
}

type ConfigResponse struct{}

type TriggerRequest struct {
	BuildTrigger BuildTrigger `json:"buildTrigger"`
	CodeTrigger  CodeTrigger  `json:"codeTrigger"`
	NodeTrigger  NodeTrigger  `json:"nodeTrigger"`
}

type BuildTrigger struct {
	LoggingTrigger LoggingTrigger `json:"loggingTrigger"`
	GerritTrigger  GerritTrigger  `json:"gerritTrigger"`
}

type GerritTrigger struct {
	Host                  string `json:"host"`
	Port                  string `json:"port"`
	Project               string `json:"project"`
	Topic                 string `json:"topic"`
	Branch                string `json:"branch"`
	EventType             string `json:"eventType"`
	Scheme                string `json:"scheme"`
	Refspec               string `json:"refspec"`
	ChangeID              string `json:"changeID"`
	ChangeUrl             string `json:"changeUrl"`
	ChangeNumber          string `json:"changeNumber"`
	ChangeSubject         string `json:"changeSubject"`
	ChangeOwner           string `json:"changeOwner"`
	ChangeOwnerName       string `json:"changeOwnerName"`
	ChangeOwnerEmail      string `json:"changeOwnerEmail"`
	ChangeWIPState        string `json:"changeWIPState"`
	ChangePrivateState    string `json:"changePrivateState"`
	ChangeCommitMessage   string `json:"changeCommitMessage"`
	PatchsetNumber        string `json:"patchsetNumber"`
	PatchsetRevision      string `json:"patchsetRevision"`
	PatchsetUploader      string `json:"patchsetUploader"`
	PatchsetUploaderName  string `json:"patchsetUploaderName"`
	PatchsetUploaderEmail string `json:"patchsetUploaderEmail"`
}

type CodeTrigger struct{}

type NodeTrigger struct {
	SshConfig SshConfig `json:"sshConfig"`
}

type LoggingTrigger struct {
	Lines []string `json:"lines"`
	Start int64    `json:"start"`
	Len   int64    `json:"len"`
}

type SshConfig struct {
	Host    string `json:"host"`
	Port    int64  `json:"port"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	Key     string `json:"key"`
	Timeout string `json:"timeout"`
}

type TriggerResponse struct {
	BuildInfo BuildInfo `json:"buildInfo"`
	CodeInfo  CodeInfo  `json:"codeInfo"`
	NodeInfo  NodeInfo  `json:"nodeInfo"`
}

type BuildInfo struct {
	LoggingInfo LoggingInfo `json:"loggingInfo"`
	RepoInfo    RepoInfo    `json:"repoInfo"`
	ReviewInfo  ReviewInfo  `json:"reviewInfo"`
}

type CodeInfo struct{}

type NodeInfo struct {
	NodeStat   NodeStat   `json:"nodeStat"`
	NodeReport NodeReport `json:"nodeReport"`
}

type LoggingInfo struct {
	File   string `json:"file"`
	Line   int64  `json:"line"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

type RepoInfo struct {
	Project   string `json:"project"`
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	Committer string `json:"committer"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	Date      string `json:"date"`
}

type ReviewInfo struct {
	Project string `json:"project"`
	Branch  string `json:"branch"`
	Change  int64  `json:"change"`
	Owner   string `json:"owner"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

type NodeStat struct {
	CpuStat     CpuStat     `json:"cpuStat"`
	DiskStat    DiskStat    `json:"diskStat"`
	DockerStat  DockerStat  `json:"dockerStat"`
	HostStat    HostStat    `json:"hostStat"`
	LoadStat    LoadStat    `json:"loadStat"`
	MemStat     MemStat     `json:"memStat"`
	NetStat     NetStat     `json:"netStat"`
	ProcessStat ProcessStat `json:"processStat"`
}

type NodeReport struct {
	CpuReport     string `json:"cpuReport"`
	DiskReport    string `json:"diskReport"`
	DockerReport  string `json:"dockerReport"`
	HealthReport  string `json:"healthReport"`
	HostReport    string `json:"hostReport"`
	LoadReport    string `json:"loadReport"`
	MemReport     string `json:"memReport"`
	NetReport     string `json:"netReport"`
	ProcessReport string `json:"processReport"`
}

type CpuStat struct {
	PhysicalCount int64     `json:"physicalCount"`
	LogicalCount  int64     `json:"logicalCount"`
	CpuPercents   []float64 `json:"cpuPercents"`
	CpuTimes      []CpuTime `json:"cpuTimes"`
}

type DiskStat struct {
	DiskPartitions []DiskPartition `json:"diskPartitions"`
	DiskUsage      DiskUsage       `json:"diskUsage"`
}

type DockerStat struct {
	CGroupCpuDockerUsages []float64          `json:"cgroupCpuDockerUsages"`
	CGroupDockerStats     []CGroupDockerStat `json:"cgroupDockerStats"`
	CGroupMemDockers      []CGroupMemDocker  `json:"cgroupMemDockers"`
}

type HostStat struct {
	Hostname        string `json:"hostname"`
	Procs           uint64 `json:"procs"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformFamily  string `json:"platformFamily"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
	KernelArch      string `json:"kernelArch"`
	HostID          string `json:"hostID"`
}

type LoadStat struct {
	LoadAvg  LoadAvg  `json:"loadAvg"`
	LoadMisc LoadMisc `json:"loadMisc"`
}

type MemStat struct {
	MemSwapDevices []MemSwapDevice `json:"memSwapDevices"`
	MemSwapMemory  MemSwapMemory   `json:"memSwapMemory"`
	MemVirtual     MemVirtual      `json:"memVirtual"`
}

type NetStat struct {
	NetIos        []NetIo        `json:"netIos"`
	NetInterfaces []NetInterface `json:"netInterfaces"`
}

type ProcessStat struct {
	ProcessInfos []ProcessInfo `json:"processInfos"`
}

type CpuTime struct {
	Cpu       string  `json:"cpu"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	IoWait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	SoftIrq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}

type DiskPartition struct {
	Device     string   `json:"device"`
	MountPoint string   `json:"mountpoint"`
	FsType     string   `json:"fstype"`
	Opts       []string `json:"opts"`
}

type DiskUsage struct {
	Path        string  `json:"path"`
	FsType      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type CGroupDockerStat struct {
	ContainerId string `json:"containerId"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	Running     bool   `json:"running"`
}

type CGroupMemDocker struct {
	Cache              uint64 `json:"cache"`
	Rss                uint64 `json:"rss"`
	RssHuge            uint64 `json:"rssHuge"`
	MappedFile         uint64 `json:"mappedFile"`
	TotalCache         uint64 `json:"totalCache"`
	TotalRss           uint64 `json:"totalRss"`
	TotalRssHuge       uint64 `json:"totalRssHuge"`
	TotalMappedFile    uint64 `json:"totalMappedFile"`
	MemUsageInBytes    uint64 `json:"memUsageInBytes"`
	MemMaxUsageInBytes uint64 `json:"memMaxUsageInBytes"`
	MemLimitInBytes    uint64 `json:"memLimitInBytes"`
}

type LoadAvg struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type LoadMisc struct {
	ProcsTotal   int64 `json:"procsTotal"`
	ProcsCreated int64 `json:"procsCreated"`
	ProcsRunning int64 `json:"procsRunning"`
	ProcsBlocked int64 `json:"procsBlocked"`
	Ctxt         int64 `json:"ctxt"`
}

type MemSwapDevice struct {
	Name      string `json:"name"`
	UsedBytes uint64 `json:"usedBytes"`
	FreeBytes uint64 `json:"freeBytes"`
}

type MemSwapMemory struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type MemVirtual struct {
	Total          uint64  `json:"total"`
	Available      uint64  `json:"available"`
	Used           uint64  `json:"used"`
	UsedPercent    float64 `json:"usedPercent"`
	Free           uint64  `json:"free"`
	Buffer         uint64  `json:"buffer"`
	Cached         uint64  `json:"cached"`
	SwapCached     uint64  `json:"swapCached"`
	SwapTotal      uint64  `json:"swapTotal"`
	SwapFree       uint64  `json:"swapFree"`
	Mapped         uint64  `json:"mapped"`
	VMallocTotal   uint64  `json:"vmallocTotal"`
	VMallocUsed    uint64  `json:"vmallocUsed"`
	VMallocChunk   uint64  `json:"vmallocChunk"`
	HugePagesTotal uint64  `json:"hugePagesTotal"`
	HugePagesFree  uint64  `json:"hugePagesFree"`
	HugePagesRsvd  uint64  `json:"hugePagesRsvd"`
	HugePagesSurp  uint64  `json:"hugePagesSurp"`
	HugePageSize   uint64  `json:"hugePageSize"`
	AnonHugePage   uint64  `json:"anonHugePage"`
}

type NetIo struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
}

type NetInterface struct {
	Index        int64    `json:"index"`
	Mtu          int64    `json:"mtu"`
	Name         string   `json:"name"`
	HardwareAddr string   `json:"hardwareAddr"`
	Flags        []string `json:"flags"`
	Addresses    []string `json:"addrs"`
}

type ProcessInfo struct {
	Background        bool              `json:"background"`
	CpuPercent        float64           `json:"cpuPercent"`
	Children          []int32           `json:"children"`
	Cmdline           string            `json:"cmdline"`
	Environs          []string          `json:"environs"`
	IoNice            int32             `json:"ionice"`
	IsRunning         bool              `json:"isRunning"`
	ProcessMemoryInfo ProcessMemoryInfo `json:"processMemoryInfo"`
	MemoryPercent     float32           `json:"memoryPercent"`
	Name              string            `json:"name"`
	NumFd             int32             `json:"numFd"`
	NumThread         int32             `json:"numThread"`
	Parent            int32             `json:"parent"`
	Ppid              int32             `json:"ppid"`
	ProcessRLimits    []ProcessRLimit   `json:"processRlimit"`
	Statuses          []string          `json:"statuss"`
	UIDs              []int32           `json:"uids"`
	Username          string            `json:"username"`
}

type ProcessMemoryInfo struct {
	Rss    uint64 `json:"rss"`
	Vms    uint64 `json:"vms"`
	Hwm    uint64 `json:"hwm"`
	Data   uint64 `json:"data"`
	Stack  uint64 `json:"stack"`
	Locked uint64 `json:"locked"`
	Swap   uint64 `json:"swap"`
}

type ProcessRLimit struct {
	Resource int32  `json:"resource"`
	Soft     uint64 `json:"soft"`
	Hard     uint64 `json:"hard"`
	Used     uint64 `json:"used"`
}
