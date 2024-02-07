//go:build ssh_test

//go test -cover -covermode=atomic -parallel 2 -tags=ssh_test -v github.com/devops-pipeflow/insight-plugin/ssh

package ssh

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/devops-pipeflow/insight-plugin/config"
)

func initSsh() ssh {
	s := ssh{
		cfg:     DefaultConfig(),
		client:  nil,
		session: nil,
		host:    "127.0.0.1",
		port:    22,
		user:    "user",
		pass:    "pass",
		key:     "",
	}

	s.cfg.Config = config.Config{}

	s.cfg.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "ssh",
		Level: hclog.LevelFromString("INFO"),
	})

	return s
}

func TestSshInit(t *testing.T) {
	t.Skip("Skipping TestSshInit.")
}

func TestSshDeinit(t *testing.T) {
	t.Skip("Skipping TestSshDeinit.")
}

func TestSshRun(t *testing.T) {
	t.Skip("Skipping TestSshRun.")
}

func TestSshInitSession(t *testing.T) {
	ctx := context.Background()
	s := initSsh()

	err := s.initSession(ctx)
	assert.Equal(t, nil, err)

	_ = s.deinitSession(ctx)
}

func TestSshDeinitSession(t *testing.T) {
	ctx := context.Background()
	s := initSsh()

	err := s.deinitSession(ctx)
	assert.Equal(t, nil, err)
}

func TestSshRunSession(t *testing.T) {
	ctx := context.Background()
	s := initSsh()

	_ = s.initSession(ctx)

	out, err := s.runSession(ctx, "echo \"Hello World!\"")
	assert.Equal(t, nil, err)

	fmt.Println(out)

	_ = s.deinitSession(ctx)
}

func TestSshSetAuth(t *testing.T) {
	ctx := context.Background()
	s := initSsh()

	auth, err := s.setAuth(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(auth))
}

func TestSshSetTimeout(t *testing.T) {
	ctx := context.Background()
	s := initSsh()

	s.cfg.Config.Spec.NodeConfig.Duration = ""

	duration, err := s.setTimeout(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, connDuration, duration)

	s.cfg.Config.Spec.NodeConfig.Duration = "1s"

	duration, err = s.setTimeout(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1*time.Second, duration)

	s.cfg.Config.Spec.NodeConfig.Duration = "10m"

	duration, err = s.setTimeout(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 10*time.Minute, duration)

	s.cfg.Config.Spec.NodeConfig.Duration = "100h"

	duration, err = s.setTimeout(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 100*time.Hour, duration)
}