package linters

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/devops-pipeflow/insight-plugin/config"
)

func initGptLinter() gptlinter {
	gl := gptlinter{
		cfg: DefaultGptLinterConfig(),
	}

	gl.cfg.Config = config.Config{}
	gl.cfg.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "gptlinter",
		Level: hclog.LevelFromString("INFO"),
	})

	return gl
}

func TestGptLinterRun(t *testing.T) {
	_ = initGptLinter()

	// TBD: FIXME

	assert.Equal(t, nil, nil)
}
