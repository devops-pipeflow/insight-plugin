# insight-plugin

[![Build Status](https://github.com/devops-pipeflow/insight-plugin/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/devops-pipeflow/insight-plugin/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/devops-pipeflow/insight-plugin/branch/main/graph/badge.svg?token=y5anikgcTz)](https://codecov.io/gh/devops-pipeflow/insight-plugin)
[![Go Report Card](https://goreportcard.com/badge/github.com/devops-pipeflow/insight-plugin)](https://goreportcard.com/report/github.com/devops-pipeflow/insight-plugin)
[![License](https://img.shields.io/github/license/devops-pipeflow/insight-plugin.svg)](https://github.com/devops-pipeflow/insight-plugin/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/devops-pipeflow/insight-plugin.svg)](https://github.com/devops-pipeflow/insight-plugin/tags)



## Introduction

*insight-plugin* is the insight plugin of [devops-pipeflow](https://github.com/devops-pipeflow) written in Go.



## Prerequisites

- Go >= 1.21.0



## Run

```bash
version=latest make build
./bin/insight --config-file="$PWD"/config/config.yml
```



## Usage

```
usage: insight --config-file=CONFIG-FILE [<flags>]

insight plugin


Flags:
  --[no-]help                Show context-sensitive help (also try --help-long
                             and --help-man).
  --[no-]version             Show application version.
  --config-file=CONFIG-FILE  Config file (.yml)
  --log-level="INFO"         Log level (DEBUG|INFO|WARN|ERROR)
```



## Settings

*insight-plugin* parameters can be set in the directory [config](https://github.com/devops-pipeflow/insight-plugin/blob/main/config).

An example of configuration in [config.yml](https://github.com/devops-pipeflow/insight-plugin/blob/main/config/config.yml):

```yaml
apiVersion: v1
kind: insight
metadata:
  name: insight
spec:
  envVariables:
    - key: env
      value: val
  buildConfig:
    loggingConfig:
      start: 1
      len: 2
      count: 3
    repoConfig:
      url: 127.0.0.1:8080
      user: user
      pass: pass
    reviewConfig:
      url: 127.0.0.1:8081
      user: user
      pass: pass
  codeConfig:
  gptConfig:
  nodeConfig:
    duration: 1s
    interval: 2s
```

> `nodeConfig`: Node config
> > `duration`: Node sight duration (h:hour, m:minute, s:second)
> >
> > `interval`: Node sight interval (h:hour, m:minute, s:second)



## Proto

```
syntax = "proto3";

package server.plugins.insight.proto;

option go_package = "github.com/devops-pipeflow/server/plugins/insight/proto";

service Insight {
  rpc Config(ConfigRequest) returns (ConfigResponse) {};
  rpc Trigger(TriggerRequest) returns (TriggerResponse) {};
}

message ConfigRequest {
  string pluginName = 1;
  repeated EnvVariable envVariables = 2;
  BuildConfig buildConfig = 3;
  CodeConfig codeConfig = 4;
  GptConfig gptConfig = 5;
  NodeConfig nodeConfig = 6;
  RepoConfig repoConfig = 7;
  ReviewConfig reviewConfig = 8;
}

message EnvVariable {
  string key = 1;
  string value = 2;
}

message BuildConfig {
  LoggingConfig loggingConfig = 1;
}

message CodeConfig {}

message GptConfig {}

message NodeConfig {
  int64 duration = 1;
  int64 interval = 2;
}

message RepoConfig {
  string url = 1;
  string user = 2;
  string pass = 3;
}

message ReviewConfig {
  string url = 1;
  string user = 2;
  string pass = 3;
}

message LoggingConfig {
  int64 start = 1;
  int64 len = 2;
  int64 count = 3;
}

message ConfigResponse {}

message TriggerRequest {
  BuildTrigger buildTrigger = 1;
  CodeTrigger codeTrigger = 2;
  GptTrigger gptTrigger = 3;
  NodeTrigger nodeTrigger = 4;
}

message BuildTrigger {
  LoggingTrigger loggingTrigger = 1;
}

message CodeTrigger {}

message GptTrigger {}

message NodeTrigger {
  repeated NodeConnect nodeConnects = 1;
}

message LoggingTrigger {
  repeated string lines = 1;
  int64 start = 2;
  int64 len = 3;
}

message NodeConnect {
  string host = 1;
  int64 port = 2;
  NodeSsh nodeSsh = 3;
}

message NodeSsh {
  string user = 1;
  string pass = 2;
  string key = 3;
}

message TriggerResponse {
  BuildInfo buildInfo = 1;
  CodeInfo codeInfo = 2;
  GptInfo gptInfo = 3;
  NodeInfo nodeInfo = 4;
}

message BuildInfo {
  LoggingInfo loggingInfo = 1;
  RepoInfo repoInfo = 2;
  ReviewInfo reviewInfo =3;
}

message CodeInfo {}

message GptInfo {}

message NodeInfo {
  repeated NodeStat nodeStats = 1;
}

message LoggingInfo {
  string file = 1;
  int64 line = 2;
  string type = 3;
  string detail = 4;
}

message RepoInfo {
  string project = 1;
  string branch = 2;
  string commit = 3;
  string committer = 4;
  string author = 5;
  string message = 6;
  string date = 7;
}

message ReviewInfo {
  string project = 1;
  string branch = 2;
  int64 change = 3;
  string owner = 4;
  string author = 5;
  string message = 6;
  string date = 7;
}

message NodeStat {
  string host = 1;
  CpuStat cpuStat = 2;
  DiskStat diskStat = 3;
  DockerStat dockerStat = 4;
  HostStat hostStat = 5;
  LoadStat loadStat = 6;
  MemStat memStat = 7;
  NetStat netStat = 8;
  ProcessStat processStat = 9;
}

message CpuStat {
  int64 physicalCount = 1;
  int64 logicalCount = 2;
  repeated float64 cpuPercents = 3;
  repeated CpuTime cpuTimes = 4;
}

message DiskStat {
  repeated DiskPartition diskPartitions = 1;
  repeated DiskUsage diskUsages = 2;
}

message DockerStat {
  repeated string dockerIds = 1;
  repeated float64 cgroupCpuDockerUsages = 2;
  repeated float64 cgroupCpuUsages = 3;
  repeated CgroupDocker cgroupDockers = 4;
  repeated CgroupMem cgroupMems = 5;
}

message HostStat {
  string hostname = 1;
  uint64 procs = 2;
  string os = 3;
  string platform = 4;
  string platformFamily = 5;
  string platformVersion = 6;
  string kernelVersion = 7;
  string kernelArch = 8;
  string hostID = 9;
}

message LoadStat {
  LoadAvg loadAvg = 1;
  LoadMisc loadMisc = 2;
}

message MemStat {
  repeated MemSwapDevice memSwapDevices = 1;
  MemSwapMemory memSwapMemory = 2;
  MemVirtual memVirtual = 3;
}

message NetStat {
  repeated NetIo netIos = 1;
  repeated NetInterface netInterfaces = 2;
}

message ProcessStat {
  repeated ProcessInfo processInfos = 1;
}

message CpuTime {
  string cpu = 1;
  float64 user = 2;
  float64 system = 3;
  float64 idle = 4;
  float64 nice = 5;
  float64 iowait = 6;
  float64 irq = 7;
  float64 softirq = 8;
  float64 steal = 9;
  float64 guest = 10;
  float64 guestNice = 11;
}

message DiskPartition {
  string device = 1;
  string mount = 2;
  string fstype = 3;
}

message DiskUsage {
  string path = 1;
  string fstype = 2;
  uint64 total = 3;
  uint64 free = 4;
  uint64 used = 5;
  float64 usedPercent = 6;
}

message CgroupDocker {
  string containerId = 1;
  string name = 2;
  string image = 3;
  string status = 4;
  bool running = 5;
}

message CgroupMem {
  uint64 cache = 1;
  uint64 rss = 2;
  uint64 rssHuge = 3;
  uint64 mappedFile = 4;
  uint64 totalCache = 5;
  uint64 totalRss = 6;
  uint64 totalRssHuge = 7;
  uint64 totalMappedFile = 8;
  uint64 memUsageInBytes = 9;
  uint64 memMaxUsageInBytes = 10;
  uint64 memLimitInBytes = 11;
}

message LoadAvg {
  float64 load1 = 1;
  float64 load5 = 2;
  float64 load15 = 3;
}

message LoadMisc {
  int64 procsTotal = 1;
  int64 procsCreated = 2;
  int64 procsRunning = 3;
  int64 procsBlocked = 4;
  int64 ctxt = 5;
}

message MemSwapDevice {
  string name = 1;
  uint64 usedBytes = 2;
  uint64 freeBytes = 3;
}

message MemSwapMemory {
  uint64 total = 1;
  uint64 used = 2;
  uint64 free = 3;
  float64 usedPercent = 4;
}

message MemVirtual {
  uint64 total = 1;
  uint64 available = 2;
  uint64 used = 3;
  float64 usedPercent = 4;
  uint64 free = 5;
  uint64 buffer = 6;
  uint64 cached = 7;
  uint64 swapCached = 8;
  uint64 swapTotal = 9;
  uint64 swapFree = 10;
  uint64 mapped = 11;
  uint64 vmallocTotal = 12;
  uint64 vmallocUsed = 13;
  uint64 vmallocChunk = 14;
  uint64 hugePagesTotal = 15;
  uint64 hugePagesFree = 16;
  uint64 hugePagesRsvd = 17;
  uint64 hugePagesSurp = 18;
  uint64 hugePageSize = 19;
  uint64 anonHugePage = 20;
}

message NetIo {
  string name = 1;
  uint64 bytesSent = 2;
  uint64 bytesRecv = 3;
  uint64 packetsSent = 4;
  uint64 packetsRecv = 5;
}

message NetInterface {
  int64 index = 1;
  int64 mtu = 2;
  string name = 3;
  string hardwareAddr = 4;
  repeated string flags = 5;
  repeated string addrs = 6;
}

message ProcessInfo {
  bool background = 1;
  float64 cpuPercent = 2;
  repeated int32 children = 3;
  string cmdline = 4;
  repeated string environs = 5;
  int32 ionice = 6;
  bool isRunning = 7;
  ProcessMemoryInfo processMemoryInfo = 8;
  float32 memoryPercent = 9;
  string name = 10;
  int32 numFd = 11;
  int32 numThread = 12;
  int32 parent = 13;
  float64 percent = 14;
  int32 ppid = 15;
  repeated ProcessRlimit processRlimit = 16;
  repeated string statuss = 17;
  repeated int32 uids = 18;
  string username = 19;
}

message ProcessMemoryInfo {
  uint64 rss = 1;
  uint64 vms = 2;
  uint64 hwm = 3;
  uint64 data = 4;
  uint64 stack = 5;
  uint64 locked = 6;
  uint64 swap = 7;
}

message ProcessRlimit {
  int32 resource = 1;
  uint64 soft = 2;
  uint64 hard = 3;
  uint64 used = 4;
}
```

> `LoggingConfig`: Logging config
> > `start`: Logging lines start
> >
> > `len`: Logging lines length
> >
> > `count`: Logging lines count
> >
> Total size: length*count

> `LoggingInfo.type`: Logging info type
> > `error`: Logging error type
> >
> > `warn`: Logging warning type
> >
> > `info`: Logging info type



## License

Project License can be found [here](LICENSE).



## Reference

- [Build AI App on Milvus, Xinference, LangChain and Llama 2-70B](https://mp.weixin.qq.com/s?__biz=MzUzMDI5OTA5NQ==&mid=2247498399&idx=1&sn=e6646dadd9a0d5b4979472e3b41749a0&chksm=fa515b27cd26d23185bf878532bff961f4d579719c47d3fc4e584325752d0806715cb4e5f7e9&xtrack=1&scene=90&subscene=93&sessionid=1693801894&flutter_pos=26&clicktime=1693801963&enterid=1693801963&finder_biz_enter_id=4&ascene=56&fasttmpl_type=0&fasttmpl_fullversion=6837651-zh_CN-zip&fasttmpl_flag=0&realreporttime=1693801963657#rd)

- [Gerrit in Go](https://github.com/devops-lintflow/lintflow/blob/main/review/gerrit.go)

- [Gitiles in Go](https://github.com/craftslab/gorepo/blob/master/gitiles/gitiles.go)
