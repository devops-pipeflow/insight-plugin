//go:build linux

package linters

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/devops-pipeflow/insight-plugin/config"
)

// nolint:misspell
var (
	kernelLinter  = filepath.Join("..", "ubuntu", checkPatchName)
	kernelOptions = []string{
		"--max-line-length=120",
		"--no-signoff",
		"--no-summary",
		"--no-tree",
		"--terse",
	}
	kernelPath = filepath.Join("..", "test", "linters")
)

func initKernelLinter() kernellinter {
	kl := kernellinter{
		cfg: DefaultKernelLinterConfig(),
	}

	kl.cfg.Config = config.Config{}
	kl.cfg.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "kernellinter",
		Level: hclog.LevelFromString("INFO"),
	})
	kl.cfg.Options = kernelOptions
	kl.linter = kernelLinter

	return kl
}

func TestLintPatch(t *testing.T) {
	ctx := context.Background()

	linter := initKernelLinter()
	files := []string{"kernel.c"}

	ret, err := linter.lintPatch(ctx, kernelPath, files)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(ret))

	buf := strings.Split(ret[0], checkPatchSep)
	assert.Equal(t, "kernel.c", buf[0])
	line, _ := strconv.Atoi(buf[1])
	assert.Equal(t, 1, line)
	assert.Equal(t, "Warn", buf[2])
	assert.NotEqual(t, 0, len(buf[3]))

	buf = strings.Split(ret[1], checkPatchSep)
	assert.Equal(t, "kernel.c", buf[0])
	line, _ = strconv.Atoi(buf[1])
	assert.Equal(t, 7, line)
	assert.Equal(t, "Error", buf[2])
	assert.NotEqual(t, 0, len(buf[3]))
}
