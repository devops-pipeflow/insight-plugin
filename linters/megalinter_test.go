//go:build linux

package linters

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/devops-pipeflow/insight-plugin/config"
)

func initMegaLinter() megalinter {
	ml := megalinter{
		cfg: DefaultMegaLinterConfig(),
	}

	ml.cfg.Config = config.Config{}
	ml.cfg.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "megalinter",
		Level: hclog.LevelFromString("INFO"),
	})

	return ml
}

func TestMegaLinterRun(t *testing.T) {
	_ = initMegaLinter()

	// TBD: FIXME

	assert.Equal(t, nil, nil)
}
