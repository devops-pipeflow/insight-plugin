package linters

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	commitSep = ":"
)

var (
	commitPath = filepath.Join("..", "test", "linters")
)

func initCommitLinter() commitlinter {
	cl := commitlinter{
		cfg: DefaultCommitLinterConfig(),
	}

	cl.cfg.Config = config.Config{}
	cl.cfg.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "commitlinter",
		Level: hclog.LevelFromString("INFO"),
	})

	return cl
}

func TestLintConflict(t *testing.T) {
	ctx := context.Background()

	linter := initCommitLinter()
	files := []string{"commit.conflict"}

	ret, err := linter.lintConflict(ctx, commitPath, files)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(ret))

	buf := strings.Split(ret[0], commitSep)
	assert.Equal(t, "commit.conflict", buf[0])
	line, _ := strconv.Atoi(buf[1])
	assert.Equal(t, 0, line)
	assert.Equal(t, "Error", buf[2])
	assert.Equal(t, "Conflict character found", buf[3])
}

func TestLintJson(t *testing.T) {
	ctx := context.Background()

	linter := initCommitLinter()
	files := []string{"commit.json"}

	ret, err := linter.lintJson(ctx, commitPath, files)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(ret))

	buf := strings.Split(ret[0], commitSep)
	assert.Equal(t, "commit.json", buf[0])
	line, _ := strconv.Atoi(buf[1])
	assert.Equal(t, 0, line)
	assert.Equal(t, "Error", buf[2])
	assert.NotEqual(t, 0, len(buf[3]))
}

func TestLintMessage(t *testing.T) {
	ctx := context.Background()

	linter := initCommitLinter()
	files := []string{"commit.message"}

	ret, err := linter.lintMessage(ctx, commitPath, files)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(ret))

	buf := strings.Split(ret[0], commitSep)
	assert.Equal(t, messageName, buf[0])
	line, _ := strconv.Atoi(buf[1])
	assert.Equal(t, 0, line)
	assert.Equal(t, "Error", buf[2])
	r := strings.Contains(buf[3], fmt.Sprintf("Subject shorter than %d characters", subjectMin))
	assert.Equal(t, true, r)

	buf = strings.Split(ret[1], commitSep)
	assert.Equal(t, messageName, buf[0])
	line, _ = strconv.Atoi(buf[1])
	assert.Equal(t, 0, line)
	assert.Equal(t, "Error", buf[2])
	r = strings.Contains(buf[3], fmt.Sprintf("Description longer than %d characters", subjectMax))
	assert.Equal(t, true, r)
}

func TestLintNewline(t *testing.T) {
	ctx := context.Background()

	linter := initCommitLinter()
	files := []string{"commit.te"}

	ret, err := linter.lintNewline(ctx, commitPath, files)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(ret))

	buf := strings.Split(ret[0], commitSep)
	assert.Equal(t, "commit.te", buf[0])
	line, _ := strconv.Atoi(buf[1])
	assert.Equal(t, 0, line)
	assert.Equal(t, "Error", buf[2])
	assert.Equal(t, "No newline at end of file", buf[3])
}

func TestLintXml(t *testing.T) {
	ctx := context.Background()

	linter := initCommitLinter()
	files := []string{"commit.xml"}

	ret, err := linter.lintXml(ctx, commitPath, files)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(ret))

	buf := strings.Split(ret[0], commitSep)
	assert.Equal(t, "commit.xml", buf[0])
	line, _ := strconv.Atoi(buf[1])
	assert.Equal(t, 0, line)
	assert.Equal(t, "Error", buf[2])
	assert.NotEqual(t, 0, len(buf[3]))
}
