package linters

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	conflictHead = "<<<<<<< HEAD"
	conflictTail = ">>>>>>> CHANGE"
	newLine      = "\n"

	descriptionMax = 80
	messageName    = "/COMMIT_MSG"
	messageSep     = "Change-Id"
	subjectMax     = 80
	subjectMin     = 25
)

var (
	conflictExcluded = []string{".apk", ".bin", ".so"}
	jsonIncluded     = []string{".json"}
	newlineIncluded  = []string{".te", "file_contexts"}
	xmlIncluded      = []string{".xml"}
)

type CommitLinter interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, string, []string) ([]string, error)
}

type CommitLinterConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type linterFunc func(context.Context, string, []string) ([]string, error)

type commitlinter struct {
	cfg    *CommitLinterConfig
	linter map[string]linterFunc
}

func CommitLinterNew(_ context.Context, cfg *CommitLinterConfig) CommitLinter {
	return &commitlinter{
		cfg: cfg,
	}
}

func DefaultCommitLinterConfig() *CommitLinterConfig {
	return &CommitLinterConfig{}
}

func (cl *commitlinter) Init(_ context.Context) error {
	cl.cfg.Logger.Debug("commitlinter: Init")

	cl.linter = map[string]linterFunc{
		"lintConflict": cl.lintConflict,
		"lintJson":     cl.lintJson,
		"lintMessage":  cl.lintMessage,
		"lintNewline":  cl.lintNewline,
		"lintXml":      cl.lintXml,
	}

	return nil
}

func (cl *commitlinter) Deinit(_ context.Context) error {
	cl.cfg.Logger.Debug("commitlinter: Deinit")

	return nil
}

func (cl *commitlinter) Run(ctx context.Context, path string, files []string) ([]string, error) {
	cl.cfg.Logger.Debug("commitlinter: Run")

	var buf []string
	var err error

	for _, linter := range cl.linter {
		var b []string
		b, err = linter(ctx, path, files)
		if err != nil {
			break
		}
		buf = append(buf, b...)
	}

	return buf, err
}

func (cl *commitlinter) lintConflict(_ context.Context, path string, files []string) ([]string, error) {
	cl.cfg.Logger.Debug("commitlinter: lintConflict")

	var buf []string

	for _, item := range files {
		suffix := filepath.Ext(item)
		if slices.IndexFunc(conflictExcluded, func(data string) bool {
			return data == suffix
		}) >= 0 {
			continue
		}

		name := filepath.Join(path, item)
		data, err := os.ReadFile(name)
		if err != nil {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", item, 0, "Error", err.Error()))
			continue
		}

		if strings.Contains(string(data), conflictHead) || strings.Contains(string(data), conflictTail) {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", item, 0, "Error", "Conflict character found"))
		}
	}

	return buf, nil
}

func (cl *commitlinter) lintJson(_ context.Context, path string, files []string) ([]string, error) {
	cl.cfg.Logger.Debug("commitlinter: lintJson")

	var buf []string

	for _, item := range files {
		suffix := filepath.Ext(item)
		if slices.IndexFunc(jsonIncluded, func(data string) bool {
			return data == suffix
		}) < 0 {
			continue
		}

		name := filepath.Join(path, item)
		data, err := os.ReadFile(name)
		if err != nil {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", item, 0, "Error", err.Error()))
			continue
		}

		var d interface{}

		err = json.Unmarshal(data, &d)
		if err != nil {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", item, 0, "Error", err.Error()))
		}
	}

	return buf, nil
}

// nolint:gocyclo
func (cl *commitlinter) lintMessage(_ context.Context, path string, files []string) ([]string, error) {
	cl.cfg.Logger.Debug("commitlinter: lintLength")

	loadMessage := func(name string) ([]string, error) {
		var buf []string
		file, err := os.Open(name)
		if err != nil {
			return buf, err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		reader := bufio.NewReader(file)
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			buf = append(buf, string(line))
		}
		return buf, nil
	}

	stripMessage := func(data []string) ([]string, error) {
		index := -1
		for i, d := range data {
			if strings.Contains(d, messageSep) {
				index = i + 1
				break
			}
		}
		if index > 0 {
			return data[:index], nil
		}
		return []string{}, nil
	}

	var buf []string

	for _, item := range files {
		name := filepath.Join(path, item)
		lines, err := loadMessage(name)
		if err != nil {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", messageName, 0, "Error", "Failed to load message"))
			continue
		}
		lines, err = stripMessage(lines)
		if err != nil {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", messageName, 0, "Error", "Failed to strip message"))
			continue
		}
		for index, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if index == 0 {
				if len(line) < subjectMin {
					buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", messageName, 0, "Error",
						fmt.Sprintf("Subject shorter than %d characters (found %d)", subjectMin, len(line))))
				} else if len(line) > subjectMax {
					buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", messageName, 0, "Error",
						fmt.Sprintf("Subject longer than %d characters (found %d)", subjectMax, len(line))))
				} else {
					// PASS
				}
			} else {
				if len(line) > descriptionMax {
					buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", messageName, 0, "Error",
						fmt.Sprintf("Description longer than %d characters (found %d)", descriptionMax, len(line))))
				} else {
					// PASS
				}
			}
		}
	}

	return buf, nil
}

func (cl *commitlinter) lintNewline(_ context.Context, path string, files []string) ([]string, error) {
	cl.cfg.Logger.Debug("commitlinter: lintNewline")

	var buf []string

	for _, item := range files {
		suffix := filepath.Ext(item)
		if slices.IndexFunc(newlineIncluded, func(data string) bool {
			if suffix != "" {
				return data == suffix
			} else {
				return data == item
			}
		}) < 0 {
			continue
		}

		name := filepath.Join(path, item)
		data, err := os.ReadFile(name)
		if err != nil {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", item, 0, "Error", err.Error()))
			continue
		}

		if !strings.HasSuffix(string(data), newLine) {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", item, 0, "Error", "No newline at end of file"))
		}
	}

	return buf, nil
}

func (cl *commitlinter) lintXml(_ context.Context, path string, files []string) ([]string, error) {
	cl.cfg.Logger.Debug("commitlinter: lintXml")

	var buf []string

	for _, item := range files {
		suffix := filepath.Ext(item)
		if slices.IndexFunc(xmlIncluded, func(data string) bool {
			return data == suffix
		}) < 0 {
			continue
		}

		name := filepath.Join(path, item)
		data, err := os.ReadFile(name)
		if err != nil {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", item, 0, "Error", err.Error()))
			continue
		}

		var d interface{}

		err = xml.Unmarshal(data, &d)
		if err != nil {
			buf = append(buf, fmt.Sprintf("%s:%d:%s:%s", item, 0, "Error", err.Error()))
		}
	}

	return buf, nil
}
