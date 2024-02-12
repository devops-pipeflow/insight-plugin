package sights

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/gpt"
	"github.com/devops-pipeflow/insight-plugin/repo"
	"github.com/devops-pipeflow/insight-plugin/review"
)

type CodeSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type CodeSightConfig struct {
	Config config.Config
	Logger hclog.Logger
	Gpt    gpt.Gpt
	Repo   repo.Repo
	Review review.Review
}

type codesight struct {
	cfg *CodeSightConfig
}

func CodeSightNew(_ context.Context, cfg *CodeSightConfig) CodeSight {
	return &codesight{
		cfg: cfg,
	}
}

func DefaultCodeSightConfig() *CodeSightConfig {
	return &CodeSightConfig{}
}

func (cs *codesight) Init(ctx context.Context) error {
	cs.cfg.Logger.Debug("codesight: Init")

	return nil
}

func (cs *codesight) Deinit(ctx context.Context) error {
	cs.cfg.Logger.Debug("codesight: Deinit")

	return nil
}

func (cs *codesight) Run(ctx context.Context) error {
	cs.cfg.Logger.Debug("codesight: Run")

	return nil
}
