package insight

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/sights"
)

type Insight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Config     config.Config
	Logger     hclog.Logger
	BuildSight sights.BuildSight
	CodeSight  sights.CodeSight
	GptSight   sights.GptSight
}

type insight struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Insight {
	return &insight{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (i *insight) Init(ctx context.Context) error {
	i.cfg.Logger.Debug("insight: Init")

	if err := i.cfg.BuildSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init buildsight")
	}

	if err := i.cfg.CodeSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init codesight")
	}

	if err := i.cfg.GptSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init gptsight")
	}

	return nil
}

func (i *insight) Deinit(ctx context.Context) error {
	i.cfg.Logger.Debug("insight: Deinit")

	_ = i.cfg.GptSight.Deinit(ctx)
	_ = i.cfg.CodeSight.Deinit(ctx)
	_ = i.cfg.BuildSight.Deinit(ctx)

	return nil
}

func (i *insight) Run(ctx context.Context) error {
	i.cfg.Logger.Debug("insight: Run")

	return nil
}
