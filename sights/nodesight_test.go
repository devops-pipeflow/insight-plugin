package sights

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/devops-pipeflow/insight-plugin/config"
)

func initNodeSight() nodesight {
	ns := nodesight{
		cfg: DefaultNodeSightConfig(),
	}

	ns.cfg.Config = config.Config{}

	ns.cfg.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "nodesight",
		Level: hclog.LevelFromString("INFO"),
	})

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

	ns.cfg.Config.Spec.NodeConfig.Duration = ""

	duration, err := ns.setDuration(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, nodeDuration, duration)

	ns.cfg.Config.Spec.NodeConfig.Duration = "1s"

	duration, err = ns.setDuration(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1*time.Second, duration)

	ns.cfg.Config.Spec.NodeConfig.Duration = "10m"

	duration, err = ns.setDuration(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, 10*time.Minute, duration)

	ns.cfg.Config.Spec.NodeConfig.Duration = "100h"

	duration, err = ns.setDuration(ctx)
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
