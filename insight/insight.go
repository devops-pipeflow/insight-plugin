package insight

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/gpt"
	"github.com/devops-pipeflow/insight-plugin/repo"
	"github.com/devops-pipeflow/insight-plugin/report"
	"github.com/devops-pipeflow/insight-plugin/review"
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
	Gpt        gpt.Gpt
	Repo       repo.Repo
	Review     review.Review
	BuildSight sights.BuildSight
	CodeSight  sights.CodeSight
	GptSight   sights.GptSight
	Report     report.Report
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

	if err := i.cfg.Gpt.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init gpt")
	}

	if err := i.cfg.Repo.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init repo")
	}

	if err := i.cfg.Review.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init review")
	}

	if err := i.cfg.BuildSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init buildsight")
	}

	if err := i.cfg.CodeSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init codesight")
	}

	if err := i.cfg.GptSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init gptsight")
	}

	if err := i.cfg.Report.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init report")
	}

	return nil
}

func (i *insight) Deinit(ctx context.Context) error {
	i.cfg.Logger.Debug("insight: Deinit")

	_ = i.cfg.Report.Deinit(ctx)
	_ = i.cfg.GptSight.Deinit(ctx)
	_ = i.cfg.CodeSight.Deinit(ctx)
	_ = i.cfg.BuildSight.Deinit(ctx)
	_ = i.cfg.Review.Deinit(ctx)
	_ = i.cfg.Repo.Deinit(ctx)
	_ = i.cfg.Gpt.Deinit(ctx)

	return nil
}

func (i *insight) Run(ctx context.Context) error {
	i.cfg.Logger.Debug("insight: Run")

	return nil
}
