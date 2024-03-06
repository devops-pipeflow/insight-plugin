package insight

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/sights"
	pluginsInsight "github.com/devops-pipeflow/server/plugins/insight"
)

const (
	routineNum = -1
)

var (
	buildInfo pluginsInsight.BuildInfo
	codeInfo  pluginsInsight.CodeInfo
	nodeInfo  pluginsInsight.NodeInfo
)

type Insight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) (*pluginsInsight.BuildInfo, *pluginsInsight.CodeInfo, *pluginsInsight.NodeInfo, error)
}

type Config struct {
	Config     config.Config
	Logger     hclog.Logger
	BuildSight sights.BuildSight
	CodeSight  sights.CodeSight
	NodeSight  sights.NodeSight
}

type insight struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Insight {
	return &insight{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (i *insight) Init(ctx context.Context) error {
	i.cfg.Logger.Debug("insight: Init")

	if err := i.cfg.BuildSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init buildsight")
	}

	if err := i.cfg.CodeSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init codesight")
	}

	if err := i.cfg.NodeSight.Init(ctx); err != nil {
		return errors.Wrap(err, "failed to init nodesight")
	}

	return nil
}

func (i *insight) Deinit(ctx context.Context) error {
	i.cfg.Logger.Debug("insight: Deinit")

	_ = i.cfg.NodeSight.Deinit(ctx)
	_ = i.cfg.CodeSight.Deinit(ctx)
	_ = i.cfg.BuildSight.Deinit(ctx)

	return nil
}

func (i *insight) Run(ctx context.Context) (*pluginsInsight.BuildInfo, *pluginsInsight.CodeInfo, *pluginsInsight.NodeInfo, error) {
	i.cfg.Logger.Debug("insight: Run")

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(routineNum)

	g.Go(func() error {
		buildInfo, _ = i.cfg.BuildSight.Run(ctx)
		return nil
	})

	g.Go(func() error {
		codeInfo, _ = i.cfg.CodeSight.Run(ctx)
		return nil
	})

	g.Go(func() error {
		nodeInfo, _ = i.cfg.NodeSight.Run(ctx)
		return nil
	})

	_ = g.Wait()

	return &buildInfo, &codeInfo, &nodeInfo, nil
}
