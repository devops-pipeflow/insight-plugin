package sights

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/gpt"
	"github.com/devops-pipeflow/insight-plugin/proto"
	"github.com/devops-pipeflow/insight-plugin/repo"
	"github.com/devops-pipeflow/insight-plugin/review"
)

type CodeSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, *proto.CodeTrigger) (proto.CodeInfo, error)
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

	// TBD: FIXME

	return nil
}

func (cs *codesight) Deinit(ctx context.Context) error {
	cs.cfg.Logger.Debug("codesight: Deinit")

	// TBD: FIXME

	return nil
}

func (cs *codesight) Run(ctx context.Context, trigger *proto.CodeTrigger) (proto.CodeInfo, error) {
	cs.cfg.Logger.Debug("codesight: Run")

	var info proto.CodeInfo

	// TBD: FIXME

	return info, nil
}
