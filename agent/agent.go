//go:build linux

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

	"github.com/devops-pipeflow/insight-plugin/proto"
)

const (
	agentLevel = "INFO"
	agentName  = "name"
)

var (
	app      = kingpin.New(agentName, "insight agent")
	logLevel = app.Flag("log-level", "Log level (DEBUG|INFO|WARN|ERROR)").Default(agentLevel).String()
)

func main() {
	ctx := context.Background()

	if err := Run(ctx); err != nil {
		fmt.Printf("%+v\n", err)
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

	if err := runAgent(ctx, logger); err != nil {
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

func runAgent(ctx context.Context, logger hclog.Logger) error {
	logger.Debug("agent: runAgent")

	var err error
	var nodeStat proto.NodeStat

	nodeStat.CpuStat, err = fetchCpuStat(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "failed to fetch cpu stat")
	}

	nodeStat.DiskStat, err = fetchDiskStat(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "failed to fetch disk stat")
	}

	nodeStat.DockerStat, _ = fetchDockerStat(ctx, logger)

	nodeStat.HostStat, err = fetchHostStat(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "failed to fetch host stat")
	}

	nodeStat.LoadStat, err = fetchLoadStat(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "failed to fetch load stat")
	}

	nodeStat.MemStat, err = fetchMemStat(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "failed to fetch mem stat")
	}

	nodeStat.NetStat, err = fetchNetStat(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "failed to fetch net stat")
	}

	nodeStat.ProcessStat, err = fetchProcessStat(ctx, logger)
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

func fetchCpuStat(ctx context.Context, logger hclog.Logger) (proto.CpuStat, error) {
	logger.Debug("agent: fetchCpuStat")

	helper := func(times []cpu.TimesStat) []proto.CpuTime {
		var buf []proto.CpuTime
		for i := range times {
			buf = append(buf, proto.CpuTime{
				Cpu:       times[i].CPU,
				User:      times[i].User,
				System:    times[i].System,
				Idle:      times[i].Idle,
				Nice:      times[i].Nice,
				IoWait:    times[i].Iowait,
				Irq:       times[i].Irq,
				SoftIrq:   times[i].Softirq,
				Steal:     times[i].Steal,
				Guest:     times[i].Guest,
				GuestNice: times[i].GuestNice,
			})
		}
		return buf
	}

	physicalCount, err := cpu.CountsWithContext(ctx, false)
	if err != nil {
		return proto.CpuStat{}, errors.Wrap(err, "failed to fetch physical count")
	}

	logicalCount, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		return proto.CpuStat{}, errors.Wrap(err, "failed to fetch logical count")
	}

	cpuTimes, err := cpu.TimesWithContext(ctx, false)
	if err != nil {
		return proto.CpuStat{}, errors.Wrap(err, "failed to fetch cpu times")
	}

	return proto.CpuStat{
		PhysicalCount: int64(physicalCount),
		LogicalCount:  int64(logicalCount),
		CpuTimes:      helper(cpuTimes),
	}, nil
}

func fetchDiskStat(ctx context.Context, logger hclog.Logger) (proto.DiskStat, error) {
	logger.Debug("agent: fetchDiskStat")

	partitionsHelper := func(partitions []disk.PartitionStat) []proto.DiskPartition {
		var buf []proto.DiskPartition
		for i := range partitions {
			buf = append(buf, proto.DiskPartition{
				Device:     partitions[i].Device,
				MountPoint: partitions[i].Mountpoint,
				FsType:     partitions[i].Fstype,
				Opts:       partitions[i].Opts,
			})
		}
		return buf
	}

	usageHelper := func(usage *disk.UsageStat) *proto.DiskUsage {
		if usage == nil {
			return &proto.DiskUsage{}
		}
		return &proto.DiskUsage{
			Path:        usage.Path,
			FsType:      usage.Fstype,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		}
	}

	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return proto.DiskStat{}, errors.Wrap(err, "failed to fetch disk partitions")
	}

	usage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return proto.DiskStat{}, errors.Wrap(err, "failed to fetch disk usage")
	}

	return proto.DiskStat{
		DiskPartitions: partitionsHelper(partitions),
		DiskUsage:      *usageHelper(usage),
	}, nil
}

func fetchDockerStat(ctx context.Context, logger hclog.Logger) (proto.DockerStat, error) {
	logger.Debug("agent: fetchDockerStat")

	dockerStatHelper := func(stats []docker.CgroupDockerStat) []proto.CGroupDockerStat {
		var buf []proto.CGroupDockerStat
		for i := range stats {
			buf = append(buf, proto.CGroupDockerStat{
				ContainerId: stats[i].ContainerID,
				Name:        stats[i].Name,
				Image:       stats[i].Image,
				Status:      stats[i].Status,
				Running:     stats[i].Running,
			})
		}
		return buf
	}

	memDockerHelper := func(stat *docker.CgroupMemStat) proto.CGroupMemDocker {
		if stat == nil {
			return proto.CGroupMemDocker{}
		}
		return proto.CGroupMemDocker{
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

	var dockerStat proto.DockerStat

	stat, err := docker.GetDockerStatWithContext(ctx)
	if err != nil {
		return dockerStat, errors.Wrap(err, "failed to fetch docker stat")
	}

	for i := range stat {
		if cpuDockerUsage, e := docker.CgroupCPUDockerUsageWithContext(ctx, stat[i].ContainerID); e == nil {
			dockerStat.CGroupCpuDockerUsages = append(dockerStat.CGroupCpuDockerUsages, cpuDockerUsage)
		}
	}

	for i := range stat {
		if memDocker, e := docker.CgroupMemDockerWithContext(ctx, stat[i].ContainerID); e == nil {
			dockerStat.CGroupMemDockers = append(dockerStat.CGroupMemDockers, memDockerHelper(memDocker))
		}
	}

	return proto.DockerStat{
		CGroupCpuDockerUsages: dockerStat.CGroupCpuDockerUsages,
		CGroupDockerStats:     dockerStatHelper(stat),
		CGroupMemDockers:      dockerStat.CGroupMemDockers,
	}, nil
}

func fetchHostStat(ctx context.Context, logger hclog.Logger) (proto.HostStat, error) {
	logger.Debug("agent: fetchHostStat")

	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return proto.HostStat{}, errors.Wrap(err, "failed to fetch host info")
	}

	return proto.HostStat{
		Hostname:        info.Hostname,
		Procs:           info.Procs,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformFamily:  info.PlatformFamily,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
		HostID:          info.HostID,
	}, nil
}

func fetchLoadStat(ctx context.Context, logger hclog.Logger) (proto.LoadStat, error) {
	logger.Debug("agent: fetchLoadStat")

	avg, err := load.AvgWithContext(ctx)
	if err != nil {
		return proto.LoadStat{}, errors.Wrap(err, "failed to fetch load avg")
	}

	misc, err := load.MiscWithContext(ctx)
	if err != nil {
		return proto.LoadStat{}, errors.Wrap(err, "failed to fetch load misc")
	}

	return proto.LoadStat{
		LoadAvg: proto.LoadAvg{
			Load1:  avg.Load1,
			Load5:  avg.Load5,
			Load15: avg.Load15,
		},
		LoadMisc: proto.LoadMisc{
			ProcsTotal:   int64(misc.ProcsTotal),
			ProcsCreated: int64(misc.ProcsCreated),
			ProcsRunning: int64(misc.ProcsRunning),
			ProcsBlocked: int64(misc.ProcsBlocked),
			Ctxt:         int64(misc.Ctxt),
		},
	}, nil
}

func fetchMemStat(ctx context.Context, logger hclog.Logger) (proto.MemStat, error) {
	logger.Debug("agent: fetchMemStat")

	helper := func(devices []*mem.SwapDevice) []proto.MemSwapDevice {
		var buf []proto.MemSwapDevice
		for i := range devices {
			if devices[i] == nil {
				continue
			}
			buf = append(buf, proto.MemSwapDevice{
				Name:      devices[i].Name,
				UsedBytes: devices[i].UsedBytes,
				FreeBytes: devices[i].FreeBytes,
			})
		}
		return buf
	}

	swapDevices, err := mem.SwapDevicesWithContext(ctx)
	if err != nil {
		return proto.MemStat{}, errors.Wrap(err, "failed to fetch swap devices")
	}

	swapMemory, err := mem.SwapMemoryWithContext(ctx)
	if err != nil {
		return proto.MemStat{}, errors.Wrap(err, "failed to fetch swap memory")
	}

	virtualMemory, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return proto.MemStat{}, errors.Wrap(err, "failed to fetch virtual memory")
	}

	return proto.MemStat{
		MemSwapDevices: helper(swapDevices),
		MemSwapMemory: proto.MemSwapMemory{
			Total:       swapMemory.Total,
			Used:        swapMemory.Used,
			Free:        swapMemory.Free,
			UsedPercent: swapMemory.UsedPercent,
		},
		MemVirtual: proto.MemVirtual{
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
			VMallocTotal:   virtualMemory.VmallocTotal,
			VMallocUsed:    virtualMemory.VmallocUsed,
			VMallocChunk:   virtualMemory.VmallocChunk,
			HugePagesTotal: virtualMemory.HugePagesTotal,
			HugePagesFree:  virtualMemory.HugePagesFree,
			HugePagesRsvd:  virtualMemory.HugePagesRsvd,
			HugePagesSurp:  virtualMemory.HugePagesSurp,
			HugePageSize:   virtualMemory.HugePageSize,
			AnonHugePage:   virtualMemory.AnonHugePages,
		},
	}, nil
}

func fetchNetStat(ctx context.Context, logger hclog.Logger) (proto.NetStat, error) {
	logger.Debug("agent: fetchNetStat")

	netIoHelper := func(stats []net.IOCountersStat) []proto.NetIo {
		var buf []proto.NetIo
		for i := range stats {
			buf = append(buf, proto.NetIo{
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

	netInterfaceHelper := func(stats []net.InterfaceStat) []proto.NetInterface {
		var buf []proto.NetInterface
		for i := range stats {
			buf = append(buf, proto.NetInterface{
				Index:        int64(stats[i].Index),
				Mtu:          int64(stats[i].MTU),
				Name:         stats[i].Name,
				HardwareAddr: stats[i].HardwareAddr,
				Flags:        stats[i].Flags,
				Addresses:    interfaceAddrHelper(stats[i].Addrs),
			})
		}
		return buf
	}

	ioCounters, err := net.IOCountersWithContext(ctx, false)
	if err != nil {
		return proto.NetStat{}, errors.Wrap(err, "failed to fetch io counters")
	}

	interfaces, err := net.InterfacesWithContext(ctx)
	if err != nil {
		return proto.NetStat{}, errors.Wrap(err, "failed to fetch net interfaces")
	}

	return proto.NetStat{
		NetIos:        netIoHelper(ioCounters),
		NetInterfaces: netInterfaceHelper(interfaces),
	}, nil
}

func fetchProcessStat(ctx context.Context, logger hclog.Logger) (proto.ProcessStat, error) {
	logger.Debug("agent: fetchProcessStat")

	processChilderHelper := func(processes []*process.Process) []int32 {
		var buf []int32
		for i := range processes {
			if processes[i] == nil {
				continue
			}
			buf = append(buf, processes[i].Pid)
		}
		return buf
	}

	processMemoryInfoHelper := func(stat *process.MemoryInfoStat) proto.ProcessMemoryInfo {
		if stat == nil {
			return proto.ProcessMemoryInfo{}
		}
		return proto.ProcessMemoryInfo{
			Rss:    stat.RSS,
			Vms:    stat.VMS,
			Hwm:    stat.HWM,
			Data:   stat.Data,
			Stack:  stat.Stack,
			Locked: stat.Locked,
			Swap:   stat.Swap,
		}
	}

	processRlimitHelper := func(stats []process.RlimitStat) []proto.ProcessRLimit {
		var buf []proto.ProcessRLimit
		for i := range stats {
			buf = append(buf, proto.ProcessRLimit{
				Resource: stats[i].Resource,
				Soft:     stats[i].Soft,
				Hard:     stats[i].Hard,
				Used:     stats[i].Used,
			})
		}
		return buf
	}

	processInfoHelper := func(ctx context.Context, processes []*process.Process) []proto.ProcessInfo {
		var buf []proto.ProcessInfo
		for i := range processes {
			if processes[i] == nil {
				continue
			}
			background, _ := processes[i].BackgroundWithContext(ctx)
			cpuPercent, _ := processes[i].CPUPercentWithContext(ctx)
			procs, _ := processes[i].ChildrenWithContext(ctx)
			cmdline, _ := processes[i].CmdlineWithContext(ctx)
			ionice, _ := processes[i].IOniceWithContext(ctx)
			isRunning, _ := processes[i].IsRunningWithContext(ctx)
			memoryInfo, _ := processes[i].MemoryInfoWithContext(ctx)
			memoryPercent, _ := processes[i].MemoryPercentWithContext(ctx)
			name, _ := processes[i].NameWithContext(ctx)
			numFd, _ := processes[i].NumFDsWithContext(ctx)
			numThread, _ := processes[i].NumThreadsWithContext(ctx)
			ppid, _ := processes[i].PpidWithContext(ctx)
			rlimits, _ := processes[i].RlimitWithContext(ctx)
			statuses, _ := processes[i].StatusWithContext(ctx)
			uids, _ := processes[i].UidsWithContext(ctx)
			username, _ := processes[i].UsernameWithContext(ctx)
			buf = append(buf, proto.ProcessInfo{
				Background:        background,
				CpuPercent:        cpuPercent,
				Children:          processChilderHelper(procs),
				Cmdline:           cmdline,
				IoNice:            ionice,
				IsRunning:         isRunning,
				ProcessMemoryInfo: processMemoryInfoHelper(memoryInfo),
				MemoryPercent:     memoryPercent,
				Name:              name,
				NumFd:             numFd,
				NumThread:         numThread,
				Parent:            ppid,
				Ppid:              ppid,
				ProcessRLimits:    processRlimitHelper(rlimits),
				Statuses:          statuses,
				UIDs:              uids,
				Username:          username,
			})
		}
		return buf
	}

	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return proto.ProcessStat{}, errors.Wrap(err, "failed to fetch processes")
	}

	return proto.ProcessStat{
		ProcessInfos: processInfoHelper(ctx, processes),
	}, nil
}
