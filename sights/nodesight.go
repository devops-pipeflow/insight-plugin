package sights

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	nodeDuration = 30 * time.Second
	nodeInterval = 1

	routineNum = 10
)

type NodeSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, []NodeConnect) (NodeInfo, error)
}

type NodeSightConfig struct {
	Config config.Config
	Logger hclog.Logger
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

func (ns *nodesight) Run(ctx context.Context, connects []NodeConnect) (NodeInfo, error) {
	ns.cfg.Logger.Debug("nodesight: Run")

	conns, err := ns.runDetect(ctx, connects)
	if err != nil {
		return NodeInfo{}, errors.Wrap(err, "failed to run detect")
	}

	info, err := ns.buildInfo(ctx, conns)
	if err != nil {
		return NodeInfo{}, errors.Wrap(err, "failed to build info")
	}

	if err := ns.runStat(ctx, info); err != nil {
		return NodeInfo{}, errors.Wrap(err, "failed to run stat")
	}

	if err := ns.runReport(ctx, info); err != nil {
		return NodeInfo{}, errors.Wrap(err, "failed to run report")
	}

	if err := ns.runGpt(ctx, info); err != nil {
		return NodeInfo{}, errors.Wrap(err, "failed to run gpt")
	}

	return info, nil
}

func (ns *nodesight) runDetect(ctx context.Context, conns []NodeConnect) ([]NodeConnect, error) {
	ns.cfg.Logger.Debug("nodesight: runDetect")

	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(routineNum)

	// TBD: FIXME

	return nil, nil
}

func (ns *nodesight) buildInfo(ctx context.Context, conns []NodeConnect) (NodeInfo, error) {
	ns.cfg.Logger.Debug("nodesight: buildInfo")

	var info NodeInfo

	// TBD: FIXME

	return info, nil
}

func (ns *nodesight) runStat(_ context.Context, info NodeInfo) error {
	ns.cfg.Logger.Debug("nodesight: runStat")

	// TBD: FIXME

	return nil
}

func (ns *nodesight) runReport(_ context.Context, info NodeInfo) error {
	ns.cfg.Logger.Debug("nodesight: runReport")

	// TBD: FIXME

	return nil
}

func (ns *nodesight) runGpt(_ context.Context, info NodeInfo) error {
	ns.cfg.Logger.Debug("nodesight: runGpt")

	// TBD: FIXME

	return nil
}

func (ns *nodesight) setTimeout(_ context.Context) (time.Duration, error) {
	ns.cfg.Logger.Debug("nodesight: setTimeout")

	var err error
	timeout := nodeDuration

	if ns.cfg.Config.Spec.NodeConfig.Duration != "" {
		timeout, err = time.ParseDuration(ns.cfg.Config.Spec.NodeConfig.Duration)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse duration")
		}
	}

	return timeout, nil
}
