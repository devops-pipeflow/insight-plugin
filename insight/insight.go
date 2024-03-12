package insight

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/proto"
	"github.com/devops-pipeflow/insight-plugin/sights"
)

const (
	routineNum = -1
)

type Insight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, *proto.BuildTrigger, *proto.CodeTrigger, *proto.NodeTrigger) (
		proto.BuildInfo, proto.CodeInfo, proto.NodeInfo, error)
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

func (i *insight) Run(ctx context.Context,
	buildTrigger *proto.BuildTrigger, codeTrigger *proto.CodeTrigger, nodeTrigger *proto.NodeTrigger) (
	proto.BuildInfo, proto.CodeInfo, proto.NodeInfo, error) {
	i.cfg.Logger.Debug("insight: Run")

	var (
		buildInfo proto.BuildInfo
		codeInfo  proto.CodeInfo
		nodeInfo  proto.NodeInfo
		err       error
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(routineNum)

	g.Go(func() error {
		buildInfo, err = i.cfg.BuildSight.Run(ctx, buildTrigger)
		if err != nil {
			return errors.Wrap(err, "failed to run buildsight")
		}
		return nil
	})

	g.Go(func() error {
		codeInfo, err = i.cfg.CodeSight.Run(ctx, codeTrigger)
		if err != nil {
			return errors.Wrap(err, "failed to run codesight")
		}
		return nil
	})

	g.Go(func() error {
		nodeInfo, err = i.cfg.NodeSight.Run(ctx, nodeTrigger)
		if err != nil {
			return errors.Wrap(err, "failed to run nodesight")
		}
		return nil
	})

	if err = g.Wait(); err != nil {
		return buildInfo, codeInfo, nodeInfo, errors.Wrap(err, "failed to wait routine")
	}

	return buildInfo, codeInfo, nodeInfo, nil
}
