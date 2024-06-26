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
	"github.com/devops-pipeflow/insight-plugin/proto"
	"github.com/devops-pipeflow/insight-plugin/repo"
	"github.com/devops-pipeflow/insight-plugin/review"
	"github.com/devops-pipeflow/insight-plugin/sights"
	"github.com/devops-pipeflow/insight-plugin/ssh"
)

const (
	level = "INFO"
	name  = "insight"
)

var (
	app        = kingpin.New(name, "insight plugin")
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

	bs, cs, ns, err := initSights(ctx, logger, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init sights")
	}

	i, err := initInsight(ctx, logger, cfg, bs, cs, ns)
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

// nolint: lll
func initSights(ctx context.Context, logger hclog.Logger, cfg *config.Config) (sights.BuildSight, sights.CodeSight, sights.NodeSight, error) {
	buildSight := func(ctx context.Context, logger hclog.Logger, cfg *config.Config) sights.BuildSight {
		c := sights.DefaultBuildSightConfig()
		c.Config = *cfg
		c.Logger = logger
		g := gpt.DefaultConfig()
		g.Config = *cfg
		g.Logger = logger
		c.Gpt = gpt.New(ctx, g)
		r := repo.DefaultConfig()
		r.Config = *cfg
		r.Logger = logger
		c.Repo = repo.New(ctx, r)
		v := review.DefaultConfig()
		v.Config = *cfg
		v.Logger = logger
		c.Review = review.New(ctx, v)
		return sights.BuildSightNew(ctx, c)
	}

	codeSight := func(ctx context.Context, logger hclog.Logger, cfg *config.Config) sights.CodeSight {
		c := sights.DefaultCodeSightConfig()
		c.Config = *cfg
		c.Logger = logger
		g := gpt.DefaultConfig()
		g.Config = *cfg
		g.Logger = logger
		c.Gpt = gpt.New(ctx, g)
		r := repo.DefaultConfig()
		r.Config = *cfg
		r.Logger = logger
		c.Repo = repo.New(ctx, r)
		v := review.DefaultConfig()
		v.Config = *cfg
		v.Logger = logger
		c.Review = review.New(ctx, v)
		return sights.CodeSightNew(ctx, c)
	}

	nodeSight := func(ctx context.Context, logger hclog.Logger, cfg *config.Config) sights.NodeSight {
		c := sights.DefaultNodeSightConfig()
		c.Config = *cfg
		c.Logger = logger
		g := gpt.DefaultConfig()
		g.Config = *cfg
		g.Logger = logger
		c.Gpt = gpt.New(ctx, g)
		s := ssh.DefaultConfig()
		s.Config = *cfg
		s.Logger = logger
		c.Ssh = ssh.New(ctx, s)
		return sights.NodeSightNew(ctx, c)
	}

	logger.Debug("cmd: initSights")

	return buildSight(ctx, logger, cfg), codeSight(ctx, logger, cfg), nodeSight(ctx, logger, cfg), nil
}

func initInsight(ctx context.Context, logger hclog.Logger, cfg *config.Config,
	bs sights.BuildSight, cs sights.CodeSight, ns sights.NodeSight) (insight.Insight, error) {
	logger.Debug("cmd: initInsight")

	c := insight.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Config = *cfg
	c.Logger = logger

	c.BuildSight = bs
	c.CodeSight = cs
	c.NodeSight = ns

	return insight.New(ctx, c), nil
}

func runInsight(ctx context.Context, logger hclog.Logger, i insight.Insight) error {
	logger.Debug("cmd: runInsight")

	var buildTrigger proto.BuildTrigger
	var codeTrigger proto.CodeTrigger
	var nodeTrigger proto.NodeTrigger

	if err := i.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init")
	}

	s := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can"t be caught, so don't need add it
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	go func(ctx context.Context, buildTrigger *proto.BuildTrigger, codeTrigger *proto.CodeTrigger, nodeTrigger *proto.NodeTrigger) {
		logger.Debug("cmd: runInsight: Run")
		_, _, _, _, _ = i.Run(ctx, buildTrigger, codeTrigger, nodeTrigger)
	}(ctx, &buildTrigger, &codeTrigger, &nodeTrigger)

	go func(ctx context.Context, i insight.Insight, s chan os.Signal) {
		logger.Debug("cmd: runInsight: Deinit")
		<-s
		_ = i.Deinit(ctx)
		done <- true
	}(ctx, i, s)

	<-done

	return nil
}
