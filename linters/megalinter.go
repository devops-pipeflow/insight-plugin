package linters

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type MegaLinter interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, string, []string) ([]string, error)
}

type MegaLinterConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type megalinter struct {
	cfg *MegaLinterConfig
}

func MegaLinterNew(_ context.Context, cfg *MegaLinterConfig) MegaLinter {
	return &megalinter{
		cfg: cfg,
	}
}

func DefaultMegaLinterConfig() *MegaLinterConfig {
	return &MegaLinterConfig{}
}

func (ml *megalinter) Init(_ context.Context) error {
	ml.cfg.Logger.Debug("megalinter: Init")

	// TBD: FIXME

	return nil
}

func (ml *megalinter) Deinit(_ context.Context) error {
	ml.cfg.Logger.Debug("megalinter: Deinit")

	// TBD: FIXME

	return nil
}

func (ml *megalinter) Run(ctx context.Context, path string, files []string) ([]string, error) {
	ml.cfg.Logger.Debug("megalinter: Run")

	// TBD: FIXME

	return nil, nil
}
