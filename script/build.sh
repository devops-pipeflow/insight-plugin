#!/bin/bash

ldflags="-s -w"

go env -w GOPROXY=https://goproxy.cn,direct

CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=$(go env GOOS) go build -ldflags "$ldflags" -o bin/agent agent/agent.go
CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=$(go env GOOS) go build -ldflags "$ldflags" -o bin/insight example/example.go

pushd bin || exit
checksum=$(shasum -a 256 agent); sed "s:^CHECKSUM=$:CHECKSUM=\"$checksum\":g" ../agent/agent.sh > agent.sh
chmod +x agent.sh
popd || exit
