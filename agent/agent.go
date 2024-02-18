package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/devops-pipeflow/insight-plugin/sights"
)

const (
	agentDuration = 10 * time.Second
	agentLevel    = "INFO"
	agentName     = "name"
)

var (
	app          = kingpin.New(agentName, "insight agent")
	durationTime = app.Flag("duration-time", "Duration time ((h:hour, m:minute, s:second)").Required().String()
	logLevel     = app.Flag("log-level", "Log level (DEBUG|INFO|WARN|ERROR)").Default(agentLevel).String()
)

func main() {
	ctx := context.Background()

	if err := Run(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func Run(ctx context.Context) error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger, err := initLogger(ctx, *logLevel)
	if err != nil {
		return errors.Wrap(err, "failed to init logger")
	}

	duration, err := initDuration(ctx, logger, *durationTime)
	if err != nil {
		return errors.Wrap(err, "failed to init duration")
	}

	if err := runAgent(ctx, logger, duration); err != nil {
		return errors.Wrap(err, "failed to run agent")
	}

	return nil
}

func initLogger(_ context.Context, level string) (hclog.Logger, error) {
	return hclog.New(&hclog.LoggerOptions{
		Name:  agentName,
		Level: hclog.LevelFromString(level),
	}), nil
}

func initDuration(_ context.Context, logger hclog.Logger, duration string) (time.Duration, error) {
	logger.Debug("agent: initDuration")

	var d time.Duration
	var err error

	if duration != "" {
		d, err = time.ParseDuration(duration)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse duration")
		}
	} else {
		d = agentDuration
	}

	return d, nil
}

func runAgent(ctx context.Context, logger hclog.Logger, duration time.Duration) error {
	logger.Debug("agent: runAgent")

	var err error
	var nodeStat sights.NodeStat

	nodeStat.CpuStat, err = fetchCpuStat(ctx, logger, duration)
	if err != nil {
		return errors.Wrap(err, "failed to fetch cpu stat")
	}

	nodeStat.DiskStat, err = fetchDiskStat(ctx, logger, duration)
	if err != nil {
		return errors.Wrap(err, "failed to fetch disk stat")
	}

	nodeStat.DockerStat, _ = fetchDockerStat(ctx, logger, duration)

	nodeStat.HostStat, err = fetchHostStat(ctx, logger, duration)
	if err != nil {
		return errors.Wrap(err, "failed to fetch host stat")
	}

	nodeStat.LoadStat, err = fetchLoadStat(ctx, logger, duration)
	if err != nil {
		return errors.Wrap(err, "failed to fetch load stat")
	}

	nodeStat.MemStat, err = fetchMemStat(ctx, logger, duration)
	if err != nil {
		return errors.Wrap(err, "failed to fetch mem stat")
	}

	nodeStat.NetStat, err = fetchNetStat(ctx, logger, duration)
	if err != nil {
		return errors.Wrap(err, "failed to fetch net stat")
	}

	nodeStat.ProcessStat, err = fetchProcessStat(ctx, logger, duration)
	if err != nil {
		return errors.Wrap(err, "failed to fetch process stat")
	}

	buf, err := json.Marshal(nodeStat)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}

	fmt.Println(string(buf))

	return nil
}

func fetchCpuStat(ctx context.Context, logger hclog.Logger, duration time.Duration) (sights.CpuStat, error) {
	logger.Debug("agent: fetchCpuStat")

	helper := func(times []cpu.TimesStat) []sights.CpuTime {
		var buf []sights.CpuTime
		for i := range times {
			buf = append(buf, sights.CpuTime{
				Cpu:       times[i].CPU,
				User:      times[i].User,
				System:    times[i].System,
				Idle:      times[i].Idle,
				Nice:      times[i].Nice,
				Iowait:    times[i].Iowait,
				Irq:       times[i].Irq,
				Softirq:   times[i].Softirq,
				Steal:     times[i].Steal,
				Guest:     times[i].Guest,
				GuestNice: times[i].GuestNice,
			})
		}
		return buf
	}

	physicalCount, err := cpu.CountsWithContext(ctx, false)
	if err != nil {
		return sights.CpuStat{}, errors.Wrap(err, "failed to fetch physical count")
	}

	logicalCount, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		return sights.CpuStat{}, errors.Wrap(err, "failed to fetch logical count")
	}

	cpuPercents, err := cpu.PercentWithContext(ctx, duration, true)
	if err != nil {
		return sights.CpuStat{}, errors.Wrap(err, "failed to fetch cpu percent")
	}

	cpuTimes, err := cpu.TimesWithContext(ctx, true)
	if err != nil {
		return sights.CpuStat{}, errors.Wrap(err, "failed to fetch cpu times")
	}

	return sights.CpuStat{
		PhysicalCount: int64(physicalCount),
		LogicalCount:  int64(logicalCount),
		CpuPercents:   cpuPercents,
		CpuTimes:      helper(cpuTimes),
	}, nil
}

func fetchDiskStat(ctx context.Context, logger hclog.Logger, _ time.Duration) (sights.DiskStat, error) {
	logger.Debug("agent: fetchDiskStat")

	partitionsHelper := func(partitions []disk.PartitionStat) []sights.DiskPartition {
		var buf []sights.DiskPartition
		for i := range partitions {
			buf = append(buf, sights.DiskPartition{
				Device:     partitions[i].Device,
				Mountpoint: partitions[i].Mountpoint,
				Fstype:     partitions[i].Fstype,
				Opts:       partitions[i].Opts,
			})
		}
		return buf
	}

	usageHelper := func(usage *disk.UsageStat) *sights.DiskUsage {
		return &sights.DiskUsage{
			Path:        usage.Path,
			Fstype:      usage.Fstype,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		}
	}

	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return sights.DiskStat{}, errors.Wrap(err, "failed to fetch disk partitions")
	}

	usage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return sights.DiskStat{}, errors.Wrap(err, "failed to fetch disk usage")
	}

	return sights.DiskStat{
		DiskPartitions: partitionsHelper(partitions),
		DiskUsage:      *usageHelper(usage),
	}, nil
}

func fetchDockerStat(ctx context.Context, logger hclog.Logger, _ time.Duration) (sights.DockerStat, error) {
	logger.Debug("agent: fetchDockerStat")

	dockerStatHelper := func(stats []docker.CgroupDockerStat) []sights.CgroupDockerStat {
		var buf []sights.CgroupDockerStat
		for i := range stats {
			buf = append(buf, sights.CgroupDockerStat{
				ContainerId: stats[i].ContainerID,
				Name:        stats[i].Name,
				Image:       stats[i].Image,
				Status:      stats[i].Status,
				Running:     stats[i].Running,
			})
		}
		return buf
	}

	memDockerHelper := func(stat *docker.CgroupMemStat) sights.CgroupMemDocker {
		return sights.CgroupMemDocker{
			Cache:              stat.Cache,
			Rss:                stat.RSS,
			RssHuge:            stat.RSSHuge,
			MappedFile:         stat.MappedFile,
			TotalCache:         stat.TotalCache,
			TotalRss:           stat.TotalRSS,
			TotalRssHuge:       stat.TotalRSSHuge,
			TotalMappedFile:    stat.TotalMappedFile,
			MemUsageInBytes:    stat.MemUsageInBytes,
			MemMaxUsageInBytes: stat.MemMaxUsageInBytes,
			MemLimitInBytes:    stat.MemLimitInBytes,
		}
	}

	var dockerStat sights.DockerStat

	stat, err := docker.GetDockerStatWithContext(ctx)
	if err != nil {
		return dockerStat, errors.Wrap(err, "failed to fetch docker stat")
	}

	for i := range stat {
		if cpuDockerUsage, e := docker.CgroupCPUDockerUsageWithContext(ctx, stat[i].ContainerID); e == nil {
			dockerStat.CgroupCpuDockerUsages = append(dockerStat.CgroupCpuDockerUsages, cpuDockerUsage)
		}
	}

	for i := range stat {
		if memDocker, e := docker.CgroupMemDockerWithContext(ctx, stat[i].ContainerID); e == nil {
			dockerStat.CgroupMemDockers = append(dockerStat.CgroupMemDockers, memDockerHelper(memDocker))
		}
	}

	return sights.DockerStat{
		CgroupCpuDockerUsages: dockerStat.CgroupCpuDockerUsages,
		CgroupDockerStats:     dockerStatHelper(stat),
		CgroupMemDockers:      dockerStat.CgroupMemDockers,
	}, nil
}

func fetchHostStat(ctx context.Context, logger hclog.Logger, _ time.Duration) (sights.HostStat, error) {
	logger.Debug("agent: fetchHostStat")

	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return sights.HostStat{}, errors.Wrap(err, "failed to fetch host info")
	}

	return sights.HostStat{
		Hostname:        info.Hostname,
		Procs:           info.Procs,
		Os:              info.OS,
		Platform:        info.Platform,
		PlatformFamily:  info.PlatformFamily,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
		HostID:          info.HostID,
	}, nil
}

func fetchLoadStat(ctx context.Context, logger hclog.Logger, _ time.Duration) (sights.LoadStat, error) {
	logger.Debug("agent: fetchLoadStat")

	avg, err := load.AvgWithContext(ctx)
	if err != nil {
		return sights.LoadStat{}, errors.Wrap(err, "failed to fetch load avg")
	}

	misc, err := load.MiscWithContext(ctx)
	if err != nil {
		return sights.LoadStat{}, errors.Wrap(err, "failed to fetch load misc")
	}

	return sights.LoadStat{
		LoadAvg: sights.LoadAvg{
			Load1:  avg.Load1,
			Load5:  avg.Load5,
			Load15: avg.Load15,
		},
		LoadMisc: sights.LoadMisc{
			ProcsTotal:   int64(misc.ProcsTotal),
			ProcsCreated: int64(misc.ProcsCreated),
			ProcsRunning: int64(misc.ProcsRunning),
			ProcsBlocked: int64(misc.ProcsBlocked),
			Ctxt:         int64(misc.Ctxt),
		},
	}, nil
}

func fetchMemStat(ctx context.Context, logger hclog.Logger, _ time.Duration) (sights.MemStat, error) {
	logger.Debug("agent: fetchMemStat")

	helper := func(devices []*mem.SwapDevice) []sights.MemSwapDevice {
		var buf []sights.MemSwapDevice
		for i := range devices {
			buf = append(buf, sights.MemSwapDevice{
				Name:      devices[i].Name,
				UsedBytes: devices[i].UsedBytes,
				FreeBytes: devices[i].FreeBytes,
			})
		}
		return buf
	}

	swapDevices, err := mem.SwapDevicesWithContext(ctx)
	if err != nil {
		return sights.MemStat{}, errors.Wrap(err, "failed to fetch swap devices")
	}

	swapMemory, err := mem.SwapMemoryWithContext(ctx)
	if err != nil {
		return sights.MemStat{}, errors.Wrap(err, "failed to fetch swap memory")
	}

	virtualMemory, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return sights.MemStat{}, errors.Wrap(err, "failed to fetch virtual memory")
	}

	return sights.MemStat{
		MemSwapDevices: helper(swapDevices),
		MemSwapMemory: sights.MemSwapMemory{
			Total:       swapMemory.Total,
			Used:        swapMemory.Used,
			Free:        swapMemory.Free,
			UsedPercent: swapMemory.UsedPercent,
		},
		MemVirtual: sights.MemVirtual{
			Total:          virtualMemory.Total,
			Available:      virtualMemory.Available,
			Used:           virtualMemory.Used,
			UsedPercent:    virtualMemory.UsedPercent,
			Free:           virtualMemory.Free,
			Buffer:         virtualMemory.Buffers,
			Cached:         virtualMemory.Cached,
			SwapCached:     virtualMemory.SwapCached,
			SwapTotal:      virtualMemory.SwapTotal,
			SwapFree:       virtualMemory.SwapFree,
			Mapped:         virtualMemory.Mapped,
			VmallocTotal:   virtualMemory.VmallocTotal,
			VmallocUsed:    virtualMemory.VmallocUsed,
			VmallocChunk:   virtualMemory.VmallocChunk,
			HugePagesTotal: virtualMemory.HugePagesTotal,
			HugePagesFree:  virtualMemory.HugePagesFree,
			HugePagesRsvd:  virtualMemory.HugePagesRsvd,
			HugePagesSurp:  virtualMemory.HugePagesSurp,
			HugePageSize:   virtualMemory.HugePageSize,
			AnonHugePage:   virtualMemory.AnonHugePages,
		},
	}, nil
}

func fetchNetStat(ctx context.Context, logger hclog.Logger, _ time.Duration) (sights.NetStat, error) {
	logger.Debug("agent: fetchNetStat")

	netIoHelper := func(stats []net.IOCountersStat) []sights.NetIo {
		var buf []sights.NetIo
		for i := range stats {
			buf = append(buf, sights.NetIo{
				Name:        stats[i].Name,
				BytesSent:   stats[i].BytesSent,
				BytesRecv:   stats[i].BytesRecv,
				PacketsSent: stats[i].PacketsSent,
				PacketsRecv: stats[i].PacketsRecv,
			})
		}
		return buf
	}

	interfaceAddrHelper := func(addrs net.InterfaceAddrList) []string {
		var buf []string
		for i := range addrs {
			buf = append(buf, addrs[i].Addr)
		}
		return buf
	}

	netInterfaceHelper := func(stats []net.InterfaceStat) []sights.NetInterface {
		var buf []sights.NetInterface
		for i := range stats {
			buf = append(buf, sights.NetInterface{
				Index:        int64(stats[i].Index),
				Mtu:          int64(stats[i].MTU),
				Name:         stats[i].Name,
				HardwareAddr: stats[i].HardwareAddr,
				Flags:        stats[i].Flags,
				Addrs:        interfaceAddrHelper(stats[i].Addrs),
			})
		}
		return buf
	}

	ioCounters, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		return sights.NetStat{}, errors.Wrap(err, "failed to fetch io counters")
	}

	interfaces, err := net.InterfacesWithContext(ctx)
	if err != nil {
		return sights.NetStat{}, errors.Wrap(err, "failed to fetch net interfaces")
	}

	return sights.NetStat{
		NetIos:        netIoHelper(ioCounters),
		NetInterfaces: netInterfaceHelper(interfaces),
	}, nil
}

func fetchProcessStat(ctx context.Context, logger hclog.Logger, duration time.Duration) (sights.ProcessStat, error) {
	logger.Debug("agent: fetchProcessStat")

	processChilderHelper := func(processes []*process.Process) []int32 {
		var buf []int32
		for i := range processes {
			buf = append(buf, processes[i].Pid)
		}
		return buf
	}

	processMemoryInfoHelper := func(stat *process.MemoryInfoStat) sights.ProcessMemoryInfo {
		return sights.ProcessMemoryInfo{
			Rss:    stat.RSS,
			Vms:    stat.VMS,
			Hwm:    stat.HWM,
			Data:   stat.Data,
			Stack:  stat.Stack,
			Locked: stat.Locked,
			Swap:   stat.Swap,
		}
	}

	processRlimitHelper := func(stats []process.RlimitStat) []sights.ProcessRlimit {
		var buf []sights.ProcessRlimit
		for i := range stats {
			buf = append(buf, sights.ProcessRlimit{
				Resource: stats[i].Resource,
				Soft:     stats[i].Soft,
				Hard:     stats[i].Hard,
				Used:     stats[i].Used,
			})
		}
		return buf
	}

	processInfoHelper := func(ctx context.Context, processes []*process.Process) []sights.ProcessInfo {
		var buf []sights.ProcessInfo
		for i := range processes {
			background, _ := processes[i].BackgroundWithContext(ctx)
			cpuPercent, _ := processes[i].CPUPercentWithContext(ctx)
			procs, _ := processes[i].ChildrenWithContext(ctx)
			cmdline, _ := processes[i].CmdlineWithContext(ctx)
			envs, _ := processes[i].EnvironWithContext(ctx)
			ionice, _ := processes[i].IOniceWithContext(ctx)
			isRunning, _ := processes[i].IsRunningWithContext(ctx)
			memoryInfo, _ := processes[i].MemoryInfoWithContext(ctx)
			memoryPercent, _ := processes[i].MemoryPercentWithContext(ctx)
			name, _ := processes[i].NameWithContext(ctx)
			numFd, _ := processes[i].NumFDsWithContext(ctx)
			numThread, _ := processes[i].NumThreadsWithContext(ctx)
			parent, _ := processes[i].ParentWithContext(ctx)
			ppid, _ := processes[i].PpidWithContext(ctx)
			rlimits, _ := processes[i].RlimitWithContext(ctx)
			statuses, _ := processes[i].StatusWithContext(ctx)
			uids, _ := processes[i].UidsWithContext(ctx)
			username, _ := processes[i].UsernameWithContext(ctx)
			buf = append(buf, sights.ProcessInfo{
				Background:        background,
				CpuPercent:        cpuPercent,
				Children:          processChilderHelper(procs),
				Cmdline:           cmdline,
				Environs:          envs,
				Ionice:            ionice,
				IsRunning:         isRunning,
				ProcessMemoryInfo: processMemoryInfoHelper(memoryInfo),
				MemoryPercent:     memoryPercent,
				Name:              name,
				NumFd:             numFd,
				NumThread:         numThread,
				Parent:            parent.Pid,
				Ppid:              ppid,
				ProcessRlimits:    processRlimitHelper(rlimits),
				Statuses:          statuses,
				Uids:              uids,
				Username:          username,
			})
		}
		return buf
	}

	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return sights.ProcessStat{}, errors.Wrap(err, "failed to fetch processes")
	}

	return sights.ProcessStat{
		ProcessInfos: processInfoHelper(ctx, processes),
	}, nil
}
