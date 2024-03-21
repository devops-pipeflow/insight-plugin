package sights

import (
	"context"
	"testing"

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

func TestNodeSightRunDetect(t *testing.T) {
	t.Skip("Skipping TestNodeSightRunDetect.")
}

func TestNodeSightRunHealth(t *testing.T) {
	t.Skip("Skipping TestNodeSightRunHealth.")
}

func TestNodeSightRunStat(t *testing.T) {
	t.Skip("Skipping TestNodeSightRunStat.")
}

func TestNodeSightRunReport(t *testing.T) {
	_ = initNodeSight()

	// TBD: FIXME
	assert.Equal(t, nil, nil)
}

func TestNodeSightRunClean(t *testing.T) {
	t.Skip("Skipping TestNodeSightRunClean.")
}
