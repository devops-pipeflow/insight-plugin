package linters

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type GptLinter interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, string, []string) ([]string, error)
}

type GptLinterConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type gptlinter struct {
	cfg *GptLinterConfig
}

func GptLinterNew(_ context.Context, cfg *GptLinterConfig) GptLinter {
	return &gptlinter{
		cfg: cfg,
	}
}

func DefaultGptLinterConfig() *GptLinterConfig {
	return &GptLinterConfig{}
}

func (gl *gptlinter) Init(_ context.Context) error {
	gl.cfg.Logger.Debug("gptlinter: Init")

	// TBD: FIXME

	return nil
}

func (gl *gptlinter) Deinit(_ context.Context) error {
	gl.cfg.Logger.Debug("gptlinter: Deinit")

	// TBD: FIXME

	return nil
}

func (gl *gptlinter) Run(ctx context.Context, path string, files []string) ([]string, error) {
	gl.cfg.Logger.Debug("gptlinter: Run")

	// TBD: FIXME

	return nil, nil
}
