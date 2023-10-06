package sights

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type GptSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type GptSightConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type gptsight struct {
	cfg *GptSightConfig
}

func GptSightNew(_ context.Context, cfg *GptSightConfig) GptSight {
	return &gptsight{
		cfg: cfg,
	}
}

func DefaultGptSightConfig() *GptSightConfig {
	return &GptSightConfig{}
}

func (gs *gptsight) Init(ctx context.Context) error {
	gs.cfg.Logger.Debug("gptsight: Init")

	return nil
}

func (gs *gptsight) Deinit(ctx context.Context) error {
	gs.cfg.Logger.Debug("gptsight: Deinit")

	return nil
}

func (gs *gptsight) Run(ctx context.Context) error {
	gs.cfg.Logger.Debug("gptsight: Run")

	return nil
}
