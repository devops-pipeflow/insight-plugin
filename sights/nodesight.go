package sights

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type NodeSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type NodeSightConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type NodeSightTrigger struct {
	NodeConnects []NodeConnect
}

type NodeConnect struct {
	Host    string
	Port    int64
	NodeSsh NodeSsh
}

type NodeSsh struct {
	User string
	Pass string
	Key  string
}

type NodeInfo struct {
	NodeStats   []NodeStat
	NodeReports []NodeReport
}

type NodeStat struct {
	Host        string
	CpuStat     CpuStat
	DiskStat    DiskStat
	DockerStat  DockerStat
	HostStat    HostStat
	LoadStat    LoadStat
	MemStat     MemStat
	NetStat     NetStat
	ProcessStat ProcessStat
}

type NodeReport struct {
	Host          string
	CpuReport     string
	DiskReport    string
	DockerReport  string
	HostReport    string
	LoadReport    string
	MemReport     string
	NetReport     string
	ProcessReport string
}

type CpuStat struct {
	PhysicalCount int64
	LogicalCount  int64
	CpuPercents   []float64
	CpuTimes      []CpuTime
}

type DiskStat struct {
	DiskPartitions []DiskPartition
	DiskUsages     []DiskUsage
}

type DockerStat struct {
	ContainerIds          []string
	CgroupCpuDockerUsages []float64
	CgroupCpuUsages       []float64
	CgroupDockers         []CgroupDocker
	CgroupMems            []CgroupMem
}

type HostStat struct {
	Hostname        string
	Procs           uint64
	Os              string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
	HostID          string
}

type LoadStat struct {
	LoadAvg  LoadAvg
	LoadMisc LoadMisc
}

type MemStat struct {
	MemSwapDevices []MemSwapDevice
	MemSwapMemory  MemSwapMemory
	MemVirtual     MemVirtual
}

type NetStat struct {
	NetIos        []NetIo
	NetInterfaces []NetInterface
}

type ProcessStat struct {
	ProcessInfos []ProcessInfo
}

type CpuTime struct {
	Cpu       string
	User      float64
	System    float64
	Idle      float64
	Nice      float64
	Iowait    float64
	Irq       float64
	Softirq   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}

type DiskPartition struct {
	Device string
	Mount  string
	Fstype string
}

type DiskUsage struct {
	Path        string
	Fstype      string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

type CgroupDocker struct {
	ContainerId string
	Name        string
	Image       string
	Status      string
	Running     bool
}

type CgroupMem struct {
	Cache              uint64
	Rss                uint64
	RssHuge            uint64
	MappedFile         uint64
	TotalCache         uint64
	TotalRss           uint64
	TotalRssHuge       uint64
	TotalMappedFile    uint64
	MemUsageInBytes    uint64
	MemMaxUsageInBytes uint64
	MemLimitInBytes    uint64
}

type LoadAvg struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

type LoadMisc struct {
	ProcsTotal   int64
	ProcsCreated int64
	ProcsRunning int64
	ProcsBlocked int64
	Ctxt         int64
}

type MemSwapDevice struct {
	Name      string
	UsedBytes uint64
	FreeBytes uint64
}

type MemSwapMemory struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

type MemVirtual struct {
	Total          uint64
	Available      uint64
	Used           uint64
	UsedPercent    float64
	Free           uint64
	Buffer         uint64
	Cached         uint64
	SwapCached     uint64
	SwapTotal      uint64
	SwapFree       uint64
	Mapped         uint64
	VmallocTotal   uint64
	VmallocUsed    uint64
	VmallocChunk   uint64
	HugePagesTotal uint64
	HugePagesFree  uint64
	HugePagesRsvd  uint64
	HugePagesSurp  uint64
	HugePageSize   uint64
	AnonHugePage   uint64
}

type NetIo struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
}

type NetInterface struct {
	Index        int64
	Mtu          int64
	Name         string
	HardwareAddr string
	Flags        []string
	Addrs        []string
}

type ProcessInfo struct {
	Background        bool
	CpuPercent        float64
	Children          []int32
	Cmdline           string
	Environs          []string
	Ionice            int32
	IsRunning         bool
	ProcessMemoryInfo ProcessMemoryInfo
	MemoryPercent     float32
	Name              string
	NumFd             int32
	NumThread         int32
	Parent            int32
	Percent           float64
	Ppid              int32
	ProcessRlimits    []ProcessRlimit
	Statuss           []string
	Uids              []int32
	Username          string
}

type ProcessMemoryInfo struct {
	Rss    uint64
	Vms    uint64
	Hwm    uint64
	Data   uint64
	Stack  uint64
	Locked uint64
	Swap   uint64
}

type ProcessRlimit struct {
	Resource int32
	Soft     uint64
	Hard     uint64
	Used     uint64
}

type nodesight struct {
	cfg *NodeSightConfig
}

func NodeSightNew(_ context.Context, cfg *NodeSightConfig) NodeSight {
	return &nodesight{
		cfg: cfg,
	}
}

func DefaultNodeSightConfig() *NodeSightConfig {
	return &NodeSightConfig{}
}

func (ns *nodesight) Init(_ context.Context) error {
	ns.cfg.Logger.Debug("nodesight: Init")

	return nil
}

func (ns *nodesight) Deinit(_ context.Context) error {
	ns.cfg.Logger.Debug("nodesight: Deinit")

	return nil
}

func (ns *nodesight) Run(_ context.Context) error {
	ns.cfg.Logger.Debug("nodesight: Run")

	// Node test on ssh connection
	// TBD: FIXME

	// Node statistic on cpu/disk/docker/host/load/mem/net/process
	// TBD: FIXME

	// Node report on cpu/disk/docker/host/load/mem/net/process
	// TBD: FIXME

	// Node report via gpt
	// TBD: FIXME

	return nil
}
