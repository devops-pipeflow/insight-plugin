package sights

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/gpt"
	"github.com/devops-pipeflow/insight-plugin/ssh"
)

func initNodeSight() nodesight {
	ctx := context.Background()

	ns := nodesight{
		cfg: DefaultNodeSightConfig(),
	}

	ns.cfg.Config = config.Config{}
	ns.cfg.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "nodesight",
		Level: hclog.LevelFromString("INFO"),
	})
	ns.cfg.Gpt = gpt.New(ctx, gpt.DefaultConfig())
	ns.cfg.Ssh = ssh.New(ctx, ssh.DefaultConfig())

	return ns
}

func TestNodeSightInit(t *testing.T) {
	t.Skip("Skipping TestNodeSightInit.")
}

func TestNodeSightDeinit(t *testing.T) {
	t.Skip("Skipping TestNodeSightDeinit.")
}

func TestNodeSightRun(t *testing.T) {
	t.Skip("Skipping TestNodeSightRun.")
}

func TestNodeSightSetDuration(t *testing.T) {
	ctx := context.Background()
	ns := initNodeSight()

	duration, err := ns.setDuration(ctx, "")
	assert.Equal(t, nil, err)
	assert.Equal(t, nodeDuration, duration)

	duration, err = ns.setDuration(ctx, "1s")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1*time.Second, duration)

	duration, err = ns.setDuration(ctx, "10m")
	assert.Equal(t, nil, err)
	assert.Equal(t, 10*time.Minute, duration)

	duration, err = ns.setDuration(ctx, "100h")
	assert.Equal(t, nil, err)
	assert.Equal(t, 100*time.Hour, duration)
}

func TestNodeSightRunDetect(t *testing.T) {
	// TBD: FIXME
}

func TestNodeSightRunStat(t *testing.T) {
	// TBD: FIXME
}

func TestNodeSightRunReport(t *testing.T) {
	// TBD: FIXME
}
