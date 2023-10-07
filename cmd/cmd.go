package cmd

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/gpt"
	"github.com/devops-pipeflow/insight-plugin/insight"
	"github.com/devops-pipeflow/insight-plugin/repo"
	"github.com/devops-pipeflow/insight-plugin/report"
	"github.com/devops-pipeflow/insight-plugin/review"
	"github.com/devops-pipeflow/insight-plugin/sights"
)

const (
	level = "INFO"
	name  = "insight"
)

var (
	app        = kingpin.New(name, "insight plugin").Version(config.Version + "-build-" + config.Build)
	configFile = app.Flag("config-file", "Config file (.yml)").Required().String()
	logLevel   = app.Flag("log-level", "Log level (DEBUG|INFO|WARN|ERROR)").Default(level).String()
)

func Run(ctx context.Context) error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger, err := initLogger(ctx, *logLevel)
	if err != nil {
		return errors.Wrap(err, "failed to init logger")
	}

	cfg, err := initConfig(ctx, logger, *configFile)
	if err != nil {
		return errors.Wrap(err, "failed to init config")
	}

	gt, err := initGpt(ctx, logger, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init gpt")
	}

	rp, err := initRepo(ctx, logger, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init repo")
	}

	rv, err := initReview(ctx, logger, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init review")
	}

	bs, err := initBuildSight(ctx, logger, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init buildsight")
	}

	cs, err := initCodeSight(ctx, logger, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init codesight")
	}

	gs, err := initGptSight(ctx, logger, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init gptsight")
	}

	rpt, err := initReport(ctx, logger, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init report")
	}

	i, err := initInsight(ctx, logger, cfg, gt, rp, rv, bs, cs, gs, rpt)
	if err != nil {
		return errors.Wrap(err, "failed to init insight")
	}

	if err := runInsight(ctx, logger, i); err != nil {
		return errors.Wrap(err, "failed to run insight")
	}

	return nil
}

func initLogger(_ context.Context, level string) (hclog.Logger, error) {
	return hclog.New(&hclog.LoggerOptions{
		Name:  name,
		Level: hclog.LevelFromString(level),
	}), nil
}

func initConfig(_ context.Context, logger hclog.Logger, name string) (*config.Config, error) {
	logger.Debug("cmd: initConfig")

	c := config.New()

	fi, err := os.Open(name)
	if err != nil {
		return c, errors.Wrap(err, "failed to open")
	}

	defer func() {
		_ = fi.Close()
	}()

	buf, _ := io.ReadAll(fi)

	if err := yaml.Unmarshal(buf, c); err != nil {
		return c, errors.Wrap(err, "failed to unmarshal")
	}

	return c, nil
}

func initGpt(ctx context.Context, logger hclog.Logger, cfg *config.Config) (gpt.Gpt, error) {
	logger.Debug("cmd: initGpt")

	c := gpt.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	return gpt.New(ctx, c), nil
}

func initRepo(ctx context.Context, logger hclog.Logger, cfg *config.Config) (repo.Repo, error) {
	logger.Debug("cmd: initRepo")

	c := repo.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	return repo.New(ctx, c), nil
}

func initReview(ctx context.Context, logger hclog.Logger, cfg *config.Config) (review.Review, error) {
	logger.Debug("cmd: initReview")

	c := review.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	return review.New(ctx, c), nil
}

func initBuildSight(ctx context.Context, logger hclog.Logger, cfg *config.Config) (sights.BuildSight, error) {
	logger.Debug("cmd: initBuildSight")

	c := sights.DefaultBuildSightConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	return sights.BuildSightNew(ctx, c), nil
}

func initCodeSight(ctx context.Context, logger hclog.Logger, cfg *config.Config) (sights.CodeSight, error) {
	logger.Debug("cmd: initCodeSight")

	c := sights.DefaultCodeSightConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	return sights.CodeSightNew(ctx, c), nil
}

func initGptSight(ctx context.Context, logger hclog.Logger, cfg *config.Config) (sights.GptSight, error) {
	logger.Debug("cmd: initGptSight")

	c := sights.DefaultGptSightConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	return sights.GptSightNew(ctx, c), nil
}

func initReport(ctx context.Context, logger hclog.Logger, cfg *config.Config) (report.Report, error) {
	logger.Debug("cmd: initReport")

	c := report.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	return report.New(ctx, c), nil
}

// nolint: lll
func initInsight(ctx context.Context, logger hclog.Logger, cfg *config.Config,
	gt gpt.Gpt, rp repo.Repo, rv review.Review,
	bs sights.BuildSight, cs sights.CodeSight, gs sights.GptSight,
	rpt report.Report) (insight.Insight, error) {
	logger.Debug("cmd: initInsight")

	c := insight.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	c.Gpt = gt
	c.Repo = rp
	c.Review = rv

	c.BuildSight = bs
	c.CodeSight = cs
	c.GptSight = gs

	c.Report = rpt

	return insight.New(ctx, c), nil
}

func runInsight(ctx context.Context, logger hclog.Logger, i insight.Insight) error {
	logger.Debug("cmd: runInsight")

	if err := i.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init")
	}

	s := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can"t be caught, so don't need add it
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	go func(c context.Context) {
		logger.Debug("cmd: runInsight: Run")
		_ = i.Run(c)
	}(ctx)

	go func(c context.Context, i insight.Insight, s chan os.Signal) {
		logger.Debug("cmd: runInsight: Deinit")
		<-s
		_ = i.Deinit(c)
		done <- true
	}(ctx, i, s)

	<-done

	return nil
}
