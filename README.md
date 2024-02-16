# insight-plugin

[![Build Status](https://github.com/devops-pipeflow/insight-plugin/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/devops-pipeflow/insight-plugin/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/devops-pipeflow/insight-plugin/branch/main/graph/badge.svg?token=y5anikgcTz)](https://codecov.io/gh/devops-pipeflow/insight-plugin)
[![Go Report Card](https://goreportcard.com/badge/github.com/devops-pipeflow/insight-plugin)](https://goreportcard.com/report/github.com/devops-pipeflow/insight-plugin)
[![License](https://img.shields.io/github/license/devops-pipeflow/insight-plugin.svg)](https://github.com/devops-pipeflow/insight-plugin/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/devops-pipeflow/insight-plugin.svg)](https://github.com/devops-pipeflow/insight-plugin/tags)



## Introduction

*insight-plugin* is the insight plugin of [devops-pipeflow](https://github.com/devops-pipeflow) written in Go.



## Prerequisites

- Go >= 1.22.0



## Run

```bash
# Run agent
version=latest make build
./bin/agent --duration-time=10s
```

```bash
# Run insight
version=latest make build
./bin/insight --config-file="$PWD"/config/config.yml
```



## Usage

```
usage: agent --duration-time=DURATION-TIME [<flags>]

insight agent


Flags:
  --[no-]help                    Show context-sensitive help (also try --help-long and --help-man).
  --duration-time=DURATION-TIME  Duration time ((h:hour, m:minute, s:second)
  --log-level="INFO"             Log level (DEBUG|INFO|WARN|ERROR)
```

```
usage: insight --config-file=CONFIG-FILE [<flags>]

insight plugin


Flags:
  --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
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
    - name: env
      value: val
  buildConfig:
    loggingConfig:
      start: 1
      len: 2
      count: 3
  codeConfig:
  gptConfig:
  artifactConfig:
    url: 127.0.0.1:8080
    user: user
    pass: pass
  nodeConfig:
    duration: 10s
  repoConfig:
    url: 127.0.0.1:8081
    user: user
    pass: pass
  reviewConfig:
    url: 127.0.0.1:8082
    user: user
    pass: pass
  sshConfig:
    host: 127.0.0.1
    port: 22
    user: user
    pass: pass
    key: key
    timeout: 10s
```

> `nodeConfig`: Node config
> > `duration`: Node sight duration (h:hour, m:minute, s:second)

> `sshConfig`: SSH config
> > `timeout`: SSH connection timeout (h:hour, m:minute, s:second)



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
  string pluginName = 1;  // plugin name
  repeated EnvVariable envVariables = 2;  // environment variables in list
  BuildConfig buildConfig = 3;  // buildsight config
  CodeConfig codeConfig = 4;  // codesight config
  GptConfig gptConfig = 5;  // gptsight config
  ArtifactConfig artifactConfig = 6;  // artifactory config
  NodeConfig nodeConfig = 7;  // nodesight config
  RepoConfig repoConfig = 8;  // repo config (Gitiles)
  ReviewConfig reviewConfig = 9;  // review config (Gerrit, pingview)
  SshConfig sshConfig = 10;  // ssh config
}

message EnvVariable {
  string name = 1;  // variable name
  string value = 2;  // variable value
}

message BuildConfig {
  LoggingConfig loggingConfig = 1;
}

message CodeConfig {}

message GptConfig {}

message ArtifactConfig {
  string url = 1;  // artifactory url
  string user = 2;  // artifactory user
  string pass = 3;  // artifactory pass
}

message NodeConfig {
  string duration = 1;  // node duration time in string (h:hour, m:minute, s:second)
}

message RepoConfig {
  string url = 1;  // repo url (Gitiles)
  string user = 2;  // repo user (Gitiles)
  string pass = 3;  // repo pass (Gitiles)
}

message ReviewConfig {
  string url = 1;  // review url (Gerrit, pingview)
  string user = 2;  // review user (Gerrit, pingview)
  string pass = 3;  // review pass (Gerrit, pingview)
}

message SshConfig {
  string host = 1;  // ssh host
  int64 port = 2;  // ssh port
  string user = 3;  // ssh user
  string pass = 4;  // ssh pass
  string key = 5;  // ssh private key
  string timeout = 6; // ssh timeout time in string (h:hour, m:minute, s:second)
}

message LoggingConfig {
  int64 start = 1;  // logging lines start (>=1)
  int64 len = 2;  // logging lines length
  int64 count = 3;  // logging lines count (total size: len*count)
}

message ConfigResponse {}

message TriggerRequest {
  BuildTrigger buildTrigger = 1;  // buildsight trigger
  CodeTrigger codeTrigger = 2;  // codesight trigger
  NodeTrigger nodeTrigger = 3;  // nodesight trigger
}

message BuildTrigger {
  LoggingTrigger loggingTrigger = 1;
}

message CodeTrigger {}

message NodeTrigger {}

message LoggingTrigger {
  repeated string lines = 1;  // logging lines in list
  int64 start = 2;  // logging lines start (>=1)
  int64 len = 3;  // logging lines length
}

message TriggerResponse {
  BuildInfo buildInfo = 1;  // buildsight info
  CodeInfo codeInfo = 2;  // codesight info
  NodeInfo nodeInfo = 3;  // nodesight info
}

message BuildInfo {
  LoggingInfo loggingInfo = 1;  // logging info
  RepoInfo repoInfo = 2;  // repo info (Gitiles)
  ReviewInfo reviewInfo =3;  // review info (Gerrit, pingview)
}

message CodeInfo {}

message NodeInfo {
  NodeStat nodeStat = 1;  // node statistic
  NodeReport nodeReport = 2;  // node report
}

message LoggingInfo {
  string file = 1;  // file name
  int64 line = 2;  // file line
  string type = 3;  // error type (info, warn, error)
  string detail = 4;  // error details
}

message RepoInfo {
  string project = 1;  // project name in repo
  string branch = 2;  // branch name in repo
  string commit = 3;  // commit id in repo
  string committer = 4;  // committer name in repo
  string author = 5;  // author name in repo
  string message = 6;  // commit message in repo
  string date = 7;  // commit date in repo
}

message ReviewInfo {
  string project = 1;  // project name in review
  string branch = 2;  // branch name in review
  int64 change = 3;  // change id in review
  string owner = 4;  // owner name in review
  string author = 5;  // author name in review
  string message = 6;  // commit message in review
  string date = 7;  // commit date in review
}

message NodeStat {
  CpuStat cpuStat = 1;  // cpu statistic
  DiskStat diskStat = 2;  // dist statistic
  DockerStat dockerStat = 3;  // docker statistic
  HostStat hostStat = 4;  // host statistic
  LoadStat loadStat = 5;  // load statistic
  MemStat memStat = 6;  // memory statistic
  NetStat netStat = 7;  // net statistic
  ProcessStat processStat = 8;  // process statistic
}

message NodeReport {
  string cpuReport = 1; // cpu report
  string diskReport = 2; // disk report
  string dockerReport = 3; // docker report
  string hostReport = 4; // host report
  string loadReport = 5; // load report
  string memReport = 6; // memory report
  string netReport = 7; // net report
  string processReport = 8; // process report
}

message CpuStat {
  int64 physicalCount = 1; // physical cores
  int64 logicalCount = 2;  // logical cores
  repeated double cpuPercents = 3;  // the percentage of cpu used per cpu in list
  repeated CpuTime cpuTimes = 4;  // the time of cpu used per cpu in list
}

message DiskStat {
  repeated DiskPartition diskPartitions = 1;  // disk partitions in list (for physical devices only)
  DiskUsage diskUsage = 2;  // file system usage
}

message DockerStat {
  repeated double cgroupCpuDockerUsages = 1;  // cpu usage for docker in list
  repeated CgroupDockerStat cgroupDockerStats = 2; // cgroup docker stat in list
  repeated CgroupMemDocker cgroupMemDockers = 3; // cgroup memory stat in list
}

message HostStat {
  string hostname = 1;  // host name
  uint64 procs = 2;  // number of processes
  string os = 3;  // OS name (linux)
  string platform = 4;  // platform name (ubuntu)
  string platformFamily = 5;  // platform family (debian)
  string platformVersion = 6;  // the complete OS version
  string kernelVersion = 7;  // the kernel version
  string kernelArch = 8;  // native cpu architecture (`uname -r`)
  string hostID = 9;  // host id (uuid)
}

message LoadStat {
  LoadAvg loadAvg = 1;  // load average
  LoadMisc loadMisc = 2;  // load misc
}

message MemStat {
  repeated MemSwapDevice memSwapDevices = 1;  // swap device in list
  MemSwapMemory memSwapMemory = 2;  // swap memory
  MemVirtual memVirtual = 3;  // virtual memory
}

message NetStat {
  repeated NetIo netIos = 1;  // network I/O statistics in list
  repeated NetInterface netInterfaces = 2;  // network interface in list
}

message ProcessStat {
  repeated ProcessInfo processInfos = 1;  // process info in list
}

message CpuTime {
  string cpu = 1;
  double user = 2;
  double system = 3;
  double idle = 4;
  double nice = 5;
  double iowait = 6;
  double irq = 7;
  double softirq = 8;
  double steal = 9;
  double guest = 10;
  double guestNice = 11;
}

message DiskPartition {
  string device = 1;
  string mountpoint = 2;
  string fstype = 3;
  repeated string opts = 4;
}

message DiskUsage {
  string path = 1;
  string fstype = 2;
  uint64 total = 3;
  uint64 free = 4;
  uint64 used = 5;
  double usedPercent = 6;
}

message CgroupDockerStat {
  string containerId = 1;
  string name = 2;
  string image = 3;
  string status = 4;
  bool running = 5;
}

message CgroupMemDocker {
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
  double load1 = 1;
  double load5 = 2;
  double load15 = 3;
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
  double usedPercent = 4;
}

message MemVirtual {
  uint64 total = 1;
  uint64 available = 2;
  uint64 used = 3;
  double usedPercent = 4;
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
  double cpuPercent = 2;
  repeated int32 children = 3;
  string cmdline = 4;
  repeated string environs = 5;
  int32 ionice = 6;
  bool isRunning = 7;
  ProcessMemoryInfo processMemoryInfo = 8;
  float memoryPercent = 9;
  string name = 10;
  int32 numFd = 11;
  int32 numThread = 12;
  int32 parent = 13;
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



## License

Project License can be found [here](LICENSE).



## Reference

- [Build AI App on Milvus, Xinference, LangChain and Llama 2-70B](https://mp.weixin.qq.com/s?__biz=MzUzMDI5OTA5NQ==&mid=2247498399&idx=1&sn=e6646dadd9a0d5b4979472e3b41749a0&chksm=fa515b27cd26d23185bf878532bff961f4d579719c47d3fc4e584325752d0806715cb4e5f7e9&xtrack=1&scene=90&subscene=93&sessionid=1693801894&flutter_pos=26&clicktime=1693801963&enterid=1693801963&finder_biz_enter_id=4&ascene=56&fasttmpl_type=0&fasttmpl_fullversion=6837651-zh_CN-zip&fasttmpl_flag=0&realreporttime=1693801963657#rd)
- [go-routine](https://gist.github.com/craftslab/ed14cc36bd0cd313040299722e819273)
- [gopsutil](https://github.com/shirou/gopsutil)
- [gorepo-gitiles](https://github.com/craftslab/gorepo/blob/master/gitiles/gitiles.go)
- [lintflow-gerrit](https://github.com/devops-lintflow/lintflow/blob/main/review/gerrit.go)
- [multissh](https://github.com/shanghai-edu/multissh)
