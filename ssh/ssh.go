package ssh

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

type Ssh interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type SshConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type ssh struct {
	cfg *SshConfig
}

func New(_ context.Context, cfg *SshConfig) Ssh {
	return &ssh{
		cfg: cfg,
	}
}

func DefaultConfig() *SshConfig {
	return &SshConfig{}
}

func (s *ssh) Init(ctx context.Context) error {
	s.cfg.Logger.Debug("ssh: Init")

	return nil
}

func (s *ssh) Deinit(ctx context.Context) error {
	s.cfg.Logger.Debug("ssh: Deinit")

	return nil
}

func (s *ssh) Run(ctx context.Context) error {
	s.cfg.Logger.Debug("ssh: Run")

	return nil
}
