package report

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type Report interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) (map[string]string, error)
}

type Config struct {
	Config config.Config
	Logger hclog.Logger
}

type report struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Report {
	return &report{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *report) Init(_ context.Context) error {
	r.cfg.Logger.Debug("report: Init")

	return nil
}

func (r *report) Deinit(_ context.Context) error {
	r.cfg.Logger.Debug("report: Deinit")

	return nil
}

func (r *report) Run(ctx context.Context) (map[string]string, error) {
	r.cfg.Logger.Debug("report: Run")

	buf := map[string]string{}

	return buf, nil
}
