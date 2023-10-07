package report

import (
	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

func initReport() report {
	r := report{
		cfg: DefaultConfig(),
	}

	r.cfg.Config = config.Config{}

	r.cfg.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "report",
		Level: hclog.LevelFromString("INFO"),
	})

	return r
}
