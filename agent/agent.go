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

	nodeStat.DockerStat, err = fetchDockerStat(ctx, logger, duration)
	if err != nil {
		return errors.Wrap(err, "failed to fetch docker stat")
	}

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

	var cpuStat sights.CpuStat

	physicalCount, err := cpu.CountsWithContext(ctx, false)
	if err != nil {
		return cpuStat, errors.Wrap(err, "failed to fetch physical count")
	}

	logicalCount, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		return cpuStat, errors.Wrap(err, "failed to fetch logical count")
	}

	cpuPercents, err := cpu.PercentWithContext(ctx, duration, true)
	if err != nil {
		return cpuStat, errors.Wrap(err, "failed to fetch cpu percent")
	}

	cpuTimes, err := cpu.TimesWithContext(ctx, true)
	if err != nil {
		return cpuStat, errors.Wrap(err, "failed to fetch cpu times")
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

	var diskStat sights.DiskStat

	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return diskStat, errors.Wrap(err, "failed to fetch disk partitions")
	}

	usage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return diskStat, errors.Wrap(err, "failed to fetch disk usage")
	}

	return sights.DiskStat{
		DiskPartitions: partitionsHelper(partitions),
		DiskUsage:      *usageHelper(usage),
	}, nil
}

func fetchDockerStat(_ context.Context, logger hclog.Logger, _ time.Duration) (sights.DockerStat, error) {
	logger.Debug("agent: fetchDockerStat")

	// TBD: FIXME

	return sights.DockerStat{}, nil
}

func fetchHostStat(_ context.Context, logger hclog.Logger, _ time.Duration) (sights.HostStat, error) {
	logger.Debug("agent: fetchHostStat")

	// TBD: FIXME

	return sights.HostStat{}, nil
}

func fetchLoadStat(_ context.Context, logger hclog.Logger, _ time.Duration) (sights.LoadStat, error) {
	logger.Debug("agent: fetchLoadStat")

	// TBD: FIXME

	return sights.LoadStat{}, nil
}

func fetchMemStat(_ context.Context, logger hclog.Logger, _ time.Duration) (sights.MemStat, error) {
	logger.Debug("agent: fetchMemStat")

	// TBD: FIXME

	return sights.MemStat{}, nil
}

func fetchNetStat(_ context.Context, logger hclog.Logger, _ time.Duration) (sights.NetStat, error) {
	logger.Debug("agent: fetchNetStat")

	// TBD: FIXME

	return sights.NetStat{}, nil
}

func fetchProcessStat(_ context.Context, logger hclog.Logger, _ time.Duration) (sights.ProcessStat, error) {
	logger.Debug("agent: fetchProcessStat")

	// TBD: FIXME

	return sights.ProcessStat{}, nil
}
