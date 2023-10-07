package review

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type Review interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Config config.Config
	Logger hclog.Logger
}

type review struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Review {
	return &review{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *review) Init(ctx context.Context) error {
	r.cfg.Logger.Debug("review: Init")

	return nil
}

func (r *review) Deinit(ctx context.Context) error {
	r.cfg.Logger.Debug("review: Deinit")

	return nil
}

func (r *review) Run(ctx context.Context) error {
	r.cfg.Logger.Debug("review: Run")

	return nil
}
