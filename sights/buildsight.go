package sights

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type BuildSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type BuildSightConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type buildsight struct {
	cfg *BuildSightConfig
}

func BuildSightNew(_ context.Context, cfg *BuildSightConfig) BuildSight {
	return &buildsight{
		cfg: cfg,
	}
}

func DefaultBuildSightConfig() *BuildSightConfig {
	return &BuildSightConfig{}
}

func (bs *buildsight) Init(ctx context.Context) error {
	bs.cfg.Logger.Debug("buildsight: Init")

	return nil
}

func (bs *buildsight) Deinit(ctx context.Context) error {
	bs.cfg.Logger.Debug("buildsight: Deinit")

	return nil
}

func (bs *buildsight) Run(ctx context.Context) error {
	bs.cfg.Logger.Debug("buildsight: Run")

	return nil
}
