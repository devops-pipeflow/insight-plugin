# insight-plugin

[![Build Status](https://github.com/devops-pipeflow/insight-plugin/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/devops-pipeflow/insight-plugin/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/devops-pipeflow/insight-plugin/branch/main/graph/badge.svg?token=y5anikgcTz)](https://codecov.io/gh/devops-pipeflow/insight-plugin)
[![Go Report Card](https://goreportcard.com/badge/github.com/devops-pipeflow/insight-plugin)](https://goreportcard.com/report/github.com/devops-pipeflow/insight-plugin)
[![License](https://img.shields.io/github/license/devops-pipeflow/insight-plugin.svg)](https://github.com/devops-pipeflow/insight-plugin/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/devops-pipeflow/insight-plugin.svg)](https://github.com/devops-pipeflow/insight-plugin/tags)



## Introduction

*insight-plugin* is the insight plugin of [devops-pipeflow](https://github.com/devops-pipeflow) written in Go.



## Prerequisites

- Go >= 1.18.0



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
  sights:
    - name: buildSight
      enable: true
    - name: codeSight
      enable: true
    - name: gptSight
      enable: true
    - name: nodeSight
      enable: true
  repo:
    url: 127.0.0.1:8080
    user: user
    pass: pass
  review:
    url: 127.0.0.1:8081
    user: user
    pass: pass
  gpt:
    url: 127.0.0.1:8082
    user: user
    pass: pass
```



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
  string plugin_name = 1;
  repeated EnvVariable envs = 2;
  repeated SightConfig sights = 3;
  Repo repo = 4;
  Review review = 5;
  GPT gpt = 6;
}

message EnvVariable {
  string key = 1;
  string value = 2;
}

message SightConfig {
  string name = 1;
  bool enable = 2;
  LoggingConfig = 3;
}

message LoggingConfig {
  int64 start = 1;
  int64 len = 2;
  int64 count = 3;
}

message Repo {
  string url = 1;
  string user = 2;
  string pass = 3;
}

message Review {
  string url = 1;
  string user = 2;
  string pass = 3;
}

message GPT {
  string url = 1;
  string user = 2;
  string pass = 3;
}

message ConfigResponse {}

message TriggerRequest {
  repeated string lines = 1;
  int64 start = 2;
  int64 len = 3;
}

message TriggerResponse {
  repeated SightOutput sights = 1;
}

message SightOutput {
  string name = 1;
  SightInfo sight = 2;
  RepoInfo repo = 3;
  ReviewInfo review =4;
}

message SightInfo {
  string file = 1;
  int64 line = 2;
  string type = 3;
  string details = 4;
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
```

> `LoggingConfig.start`: Logging lines start
>
> `LoggingConfig.len`: Logging lines length
>
> `LoggingConfig.count`: Logging lines count
>
> Logging lines size: length*count



## Output

```json
{
  "sights": [
    {
      "name": "buildSight",
      "sight": {
        "file": "name",
        "line": 1,
        "type": "error",
        "details": "text"
      },
      "repo": {
        "project": "name",
        "branch": "name",
        "commit": "hash",
        "committer": "name <name@example.com>",
        "author": "name <name@example.com>",
        "message": "base64",
        "date": "2023-01-01T12:34:56+0800"
      },
      "review": {
        "project": "name",
        "branch": "name",
        "change": 1,
        "owner": "name <name@example.com>",
        "author": "name <name@example.com>",
        "message": "base64",
        "date": "2023-01-01T12:34:56+0800"
      }
    }
  ]
}
```

> `sights.sight.type`: sight type
> > The sight type in `sights.sight.type` should be one of below:
> >
> > `error`
> >
> > `warn`
> >
> > `info`



## License

Project License can be found [here](LICENSE).



## Reference

- [Build AI App on Milvus, Xinference, LangChain and Llama 2-70B](https://mp.weixin.qq.com/s?__biz=MzUzMDI5OTA5NQ==&mid=2247498399&idx=1&sn=e6646dadd9a0d5b4979472e3b41749a0&chksm=fa515b27cd26d23185bf878532bff961f4d579719c47d3fc4e584325752d0806715cb4e5f7e9&xtrack=1&scene=90&subscene=93&sessionid=1693801894&flutter_pos=26&clicktime=1693801963&enterid=1693801963&finder_biz_enter_id=4&ascene=56&fasttmpl_type=0&fasttmpl_fullversion=6837651-zh_CN-zip&fasttmpl_flag=0&realreporttime=1693801963657#rd)

- [Gerrit in Go](https://github.com/devops-lintflow/lintflow/blob/main/review/gerrit.go)

- [Gitiles in Go](https://github.com/craftslab/gorepo/blob/master/gitiles/gitiles.go)
