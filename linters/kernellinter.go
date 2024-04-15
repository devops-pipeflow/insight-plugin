package linters

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	checkPatchLen  = 4
	checkPatchName = "checkpatch.pl"
	checkPatchSep  = ":"
)

var (
	checkPatchTypes = []string{"Error", "Info", "Warn"}
)

type KernelLinter interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, string, []string) ([]string, error)
}

type KernelLinterConfig struct {
	Config  config.Config
	Logger  hclog.Logger
	Options []string
}

type kernellinter struct {
	cfg    *KernelLinterConfig
	linter string
}

func KernelLinterNew(_ context.Context, cfg *KernelLinterConfig) KernelLinter {
	ex, _ := os.Executable()

	return &kernellinter{
		cfg:    cfg,
		linter: filepath.Join(filepath.Dir(ex), checkPatchName),
	}
}

func DefaultKernelLinterConfig() *KernelLinterConfig {
	return &KernelLinterConfig{}
}

func (kl *kernellinter) Init(_ context.Context) error {
	kl.cfg.Logger.Debug("kernellinter: Init")

	return nil
}

func (kl *kernellinter) Deinit(_ context.Context) error {
	kl.cfg.Logger.Debug("kernellinter: Deinit")

	return nil
}

func (kl *kernellinter) Run(ctx context.Context, path string, files []string) ([]string, error) {
	kl.cfg.Logger.Debug("kernellinter: Run")

	return kl.lintPatch(ctx, path, files)
}

// nolint: gosec
func (kl *kernellinter) lintPatch(_ context.Context, path string, files []string) ([]string, error) {
	kl.cfg.Logger.Debug("kernellinter: lintPatch")

	parseType := func(name string) string {
		var buf string
		for _, item := range checkPatchTypes {
			if strings.Contains(strings.ToLower(name), strings.ToLower(item)) {
				buf = item
				break
			}
		}
		return buf
	}

	buildLint := func(name, out string) []string {
		var buf []string
		lines := strings.Split(out, "\n")
		for _, item := range lines {
			b := strings.Split(item, checkPatchSep)
			if len(b) >= checkPatchLen {
				num, _ := strconv.Atoi(strings.TrimSpace(b[1]))
				buf = append(buf, fmt.Sprintf("%s:%d:%s:%s",
					name,
					num,
					parseType(strings.TrimSpace(b[2])),
					strings.TrimSpace(strings.Join(b[3:], " "))))
			}
		}
		return buf
	}

	var buf []string

	for _, item := range files {
		opts := kl.cfg.Options
		opts = append(opts, "-f", filepath.Join(path, item))
		cmd := exec.Command(kl.linter, opts...)
		out, _ := cmd.CombinedOutput()
		buf = append(buf, buildLint(item, string(out))...)
	}

	return buf, nil
}
