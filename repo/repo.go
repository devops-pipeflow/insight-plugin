package repo

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type Repo interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Config config.Config
	Logger hclog.Logger
}

type repo struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Repo {
	return &repo{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *repo) Init(ctx context.Context) error {
	r.cfg.Logger.Debug("repo: Init")

	return nil
}

func (r *repo) Deinit(ctx context.Context) error {
	r.cfg.Logger.Debug("repo: Deinit")

	return nil
}

func (r *repo) Run(ctx context.Context) error {
	r.cfg.Logger.Debug("repo: Run")

	return nil
}
