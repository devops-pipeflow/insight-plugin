package gpt

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type Gpt interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Config config.Config
	Logger hclog.Logger
}

type gpt struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Gpt {
	return &gpt{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (g *gpt) Init(ctx context.Context) error {
	g.cfg.Logger.Debug("gpt: Init")

	// TBD: FIXME

	return nil
}

func (g *gpt) Deinit(ctx context.Context) error {
	g.cfg.Logger.Debug("gpt: Deinit")

	// TBD: FIXME

	return nil
}

func (g *gpt) Run(ctx context.Context) error {
	g.cfg.Logger.Debug("gpt: Run")

	// TBD: FIXME

	return nil
}
