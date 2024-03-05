#!/bin/bash

ldflags="-s -w"
os=$(go env GOOS)

go env -w GOPROXY=https://goproxy.cn,direct

if [ "$os" = "windows" ]; then
  CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=windows go build -ldflags "$ldflags" -o bin/insight example/example.go
else
  CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=linux go build -ldflags "$ldflags" -o bin/agent agent/agent.go
  CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=linux go build -ldflags "$ldflags" -o bin/insight .
  pushd bin || exit
  checksum=$(shasum -a 256 agent); sed "s:^CHECKSUM=$:CHECKSUM=\"$checksum\":g" ../agent/agent.sh > agent.sh
  chmod +x agent.sh
  popd || exit
fi
