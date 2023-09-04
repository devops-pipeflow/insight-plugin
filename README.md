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
./bin/example --config-file="$PWD"/config/config.yml
```



## Usage

```
devops-pipeflow insight-plugin

Usage:
  insight-plugin [flags]

Flags:
  -c, --config-file string   config file (.yml)
  -h, --help                 help for insight-plugin
  -v, --version              version for insight-plugin
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
  buildsight:
    file:
      - name: string
        content: base64
    path:
      - name: string
        path: /path/to/file
  codesight:
    file:
      - name: string
        content: base64
    path:
      - name: string
        path: /path/to/file
  gptsight:
    file:
      - name: string
        content: base64
    path:
      - name: string
        path: /path/to/file
    gpt:
      - name: string
        url: 127.0.0.1:8080
        user: user
        pass: pass
```



## License

Project License can be found [here](LICENSE).



## Reference

- [Build AI App on Milvus, Xinference, LangChain and Llama 2-70B](https://mp.weixin.qq.com/s?__biz=MzUzMDI5OTA5NQ==&mid=2247498399&idx=1&sn=e6646dadd9a0d5b4979472e3b41749a0&chksm=fa515b27cd26d23185bf878532bff961f4d579719c47d3fc4e584325752d0806715cb4e5f7e9&xtrack=1&scene=90&subscene=93&sessionid=1693801894&flutter_pos=26&clicktime=1693801963&enterid=1693801963&finder_biz_enter_id=4&ascene=56&fasttmpl_type=0&fasttmpl_fullversion=6837651-zh_CN-zip&fasttmpl_flag=0&realreporttime=1693801963657#rd)
