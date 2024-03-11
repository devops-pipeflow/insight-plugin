//go:build ssh_test

// go test -cover -covermode=atomic -parallel 2 -tags=ssh_test -v github.com/devops-pipeflow/insight-plugin/ssh

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

const (
	sshHost    = "127.0.0.1"
	sshPort    = 22
	sshUser    = "user"
	sshPass    = "pass"
	sshKey     = ""
	sshTimeout = "10s"
)

func initSsh() ssh {
	s := ssh{
		cfg:     DefaultConfig(),
		client:  nil,
		session: nil,
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

	err := s.initSession(ctx, sshHost, sshPort, sshUser, sshPass, sshKey, sshTimeout)
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

	_ = s.initSession(ctx, sshHost, sshPort, sshUser, sshPass, sshKey, sshTimeout)

	cmds := []string{
		"echo \"Hello\"",
		"echo \"World!\"",
	}

	out, err := s.runSession(ctx, cmds)
	assert.Equal(t, nil, err)

	fmt.Println(out)

	_ = s.deinitSession(ctx)
}

func TestSshSetAuth(t *testing.T) {
	ctx := context.Background()
	s := initSsh()

	auth, err := s.setAuth(ctx, sshPass, sshKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(auth))
}

func TestSshSetTimeout(t *testing.T) {
	ctx := context.Background()
	s := initSsh()

	timeout, err := s.setTimeout(ctx, "")
	assert.Equal(t, nil, err)
	assert.Equal(t, connTimeout, timeout)

	timeout, err = s.setTimeout(ctx, "1s")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1*time.Second, timeout)

	timeout, err = s.setTimeout(ctx, "10m")
	assert.Equal(t, nil, err)
	assert.Equal(t, 10*time.Minute, timeout)

	timeout, err = s.setTimeout(ctx, "100h")
	assert.Equal(t, nil, err)
	assert.Equal(t, 100*time.Hour, timeout)
}
