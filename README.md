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
