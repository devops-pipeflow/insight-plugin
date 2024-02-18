package sights

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/gpt"
	"github.com/devops-pipeflow/insight-plugin/ssh"
)

const (
	agentExec   = "agent"
	agentPath   = "/tmp/"
	agentScript = agentExec + ".sh"

	argDurationTime = "--duration-time"
	argLogLevel     = "--log-level"
	argSep          = "="

	artifactPath = "/devops-pipeflow/plugins/"

	routineNum = -1
)

type NodeSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) (NodeInfo, error)
}

type NodeSightConfig struct {
	Config config.Config
	Logger hclog.Logger
	Gpt    gpt.Gpt
	Ssh    ssh.Ssh
}

type NodeInfo struct {
	NodeStat   NodeStat
	NodeReport NodeReport
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
	CgroupCpuDockerUsages []float64          `json:"cgroupCpuDockerUsages"`
	CgroupDockerStats     []CgroupDockerStat `json:"cgroupDockerStats"`
	CgroupMemDockers      []CgroupMemDocker  `json:"cgroupMemDockers"`
}

type HostStat struct {
	Hostname        string `json:"hostname"`
	Procs           uint64 `json:"procs"`
	Os              string `json:"os"`
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
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}

type DiskPartition struct {
	Device     string   `json:"device"`
	Mountpoint string   `json:"mountpoint"`
	Fstype     string   `json:"fstype"`
	Opts       []string `json:"opts"`
}

type DiskUsage struct {
	Path        string  `json:"path"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type CgroupDockerStat struct {
	ContainerId string `json:"containerId"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	Running     bool   `json:"running"`
}

type CgroupMemDocker struct {
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
	VmallocTotal   uint64  `json:"vmallocTotal"`
	VmallocUsed    uint64  `json:"vmallocUsed"`
	VmallocChunk   uint64  `json:"vmallocChunk"`
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
	Addrs        []string `json:"addrs"`
}

type ProcessInfo struct {
	Background        bool              `json:"background"`
	CpuPercent        float64           `json:"cpuPercent"`
	Children          []int32           `json:"children"`
	Cmdline           string            `json:"cmdline"`
	Environs          []string          `json:"environs"`
	Ionice            int32             `json:"ionice"`
	IsRunning         bool              `json:"isRunning"`
	ProcessMemoryInfo ProcessMemoryInfo `json:"processMemoryInfo"`
	MemoryPercent     float32           `json:"memoryPercent"`
	Name              string            `json:"name"`
	NumFd             int32             `json:"numFd"`
	NumThread         int32             `json:"numThread"`
	Parent            int32             `json:"parent"`
	Ppid              int32             `json:"ppid"`
	ProcessRlimits    []ProcessRlimit   `json:"processRlimits"`
	Statuses          []string          `json:"statuses"`
	Uids              []int32           `json:"uids"`
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

type ProcessRlimit struct {
	Resource int32  `json:"resource"`
	Soft     uint64 `json:"soft"`
	Hard     uint64 `json:"hard"`
	Used     uint64 `json:"used"`
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

func (ns *nodesight) Init(ctx context.Context) error {
	ns.cfg.Logger.Debug("nodesight: Init")

	return nil
}

func (ns *nodesight) Deinit(_ context.Context) error {
	ns.cfg.Logger.Debug("nodesight: Deinit")

	return nil
}

func (ns *nodesight) Run(ctx context.Context) (NodeInfo, error) {
	ns.cfg.Logger.Debug("nodesight: Run")

	var info NodeInfo

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(routineNum)

	g.Go(func() error {
		if err := ns.runDetect(ctx); err != nil {
			return errors.Wrap(err, "failed to run detect")
		}
		stat, err := ns.runStat(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to run stat")
		}
		info.NodeStat = *stat
		report, err := ns.runReport(ctx, stat)
		if err != nil {
			return errors.Wrap(err, "failed to run report")
		}
		info.NodeReport = *report
		return nil
	})

	if err := g.Wait(); err != nil {
		return info, errors.Wrap(err, "failed to wait routine")
	}

	return info, nil
}

func (ns *nodesight) runDetect(ctx context.Context) error {
	ns.cfg.Logger.Debug("nodesight: runDetect")

	if err := ns.cfg.Ssh.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init ssh")
	}

	defer func() {
		_ = ns.cfg.Ssh.Deinit(ctx)
	}()

	cmds := []string{
		fmt.Sprintf("curl -s -u%s:%s -L %s -o %s",
			ns.cfg.Config.Spec.ArtifactConfig.User,
			ns.cfg.Config.Spec.ArtifactConfig.Pass,
			ns.cfg.Config.Spec.ArtifactConfig.Url+artifactPath+agentScript,
			agentPath+agentScript),
		fmt.Sprintf("cd %s; bash %s %s %s %s %s",
			agentPath,
			agentScript,
			ns.cfg.Config.Spec.ArtifactConfig.User,
			ns.cfg.Config.Spec.ArtifactConfig.Pass,
			ns.cfg.Config.Spec.ArtifactConfig.Url+artifactPath+agentExec,
			agentPath+agentExec),
		fmt.Sprintf("rm -f %s", agentPath+agentScript),
	}

	for i := range cmds {
		out, err := ns.cfg.Ssh.Run(ctx, cmds[i])
		if err != nil {
			return errors.Wrap(err, "failed to run ssh")
		}
		if out != "" {
			return errors.Wrap(errors.New(out), "failed to deploy agent")
		}
	}

	return nil
}

func (ns *nodesight) runStat(ctx context.Context) (*NodeStat, error) {
	ns.cfg.Logger.Debug("nodesight: runStat")

	var stat NodeStat

	if err := ns.cfg.Ssh.Init(ctx); err != nil {
		return &stat, errors.Wrap(err, "failed to init ssh")
	}

	defer func() {
		_ = ns.cfg.Ssh.Deinit(ctx)
	}()

	cmd := fmt.Sprintf("%s %s %s",
		agentPath+agentExec,
		argDurationTime+argSep+ns.cfg.Config.Spec.NodeConfig.Duration,
		argLogLevel+argSep+"ERROR")

	out, err := ns.cfg.Ssh.Run(ctx, cmd)
	if err != nil {
		return &stat, errors.Wrap(err, "failed to run ssh")
	}

	if err := json.Unmarshal([]byte(out), &stat); err != nil {
		return &stat, errors.Wrap(err, "failed to unmarshal json")
	}

	return &stat, nil
}

func (ns *nodesight) runReport(_ context.Context, stat *NodeStat) (*NodeReport, error) {
	ns.cfg.Logger.Debug("nodesight: runReport")

	var report NodeReport

	// TBD: FIXME

	return &report, nil
}
