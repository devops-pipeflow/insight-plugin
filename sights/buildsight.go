package sights

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/gpt"
	"github.com/devops-pipeflow/insight-plugin/repo"
	"github.com/devops-pipeflow/insight-plugin/review"
	pluginsInsight "github.com/devops-pipeflow/server/plugins/insight"
)

type BuildSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) (pluginsInsight.BuildInfo, error)
}

type BuildSightConfig struct {
	Config config.Config
	Logger hclog.Logger
	Gpt    gpt.Gpt
	Repo   repo.Repo
	Review review.Review
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

	// TBD: FIXME

	return nil
}

func (bs *buildsight) Deinit(ctx context.Context) error {
	bs.cfg.Logger.Debug("buildsight: Deinit")

	// TBD: FIXME

	return nil
}

func (bs *buildsight) Run(context.Context) (pluginsInsight.BuildInfo, error) {
	bs.cfg.Logger.Debug("buildsight: Run")

	// TBD: FIXME

	return pluginsInsight.BuildInfo{}, nil
}
