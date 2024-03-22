package sights

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/devops-pipeflow/insight-plugin/config"
	"github.com/devops-pipeflow/insight-plugin/gpt"
	"github.com/devops-pipeflow/insight-plugin/proto"
	"github.com/devops-pipeflow/insight-plugin/ssh"
)

const (
	agentDurationTime = "--duration-time"
	agentExec         = "agent"
	agentLogLevel     = "--log-level"
	agentPath         = "/tmp/"
	agentScript       = agentExec + ".sh"
	agentSep          = "="

	artifactPath = "zd-devops-nj-release-generic/devops-pipeflow/plugins"

	healthPath   = "/tmp/"
	healthScript = "healthcheck.sh"
	healthSilent = "--silent"

	routineNum = -1
)

type NodeSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, *proto.NodeTrigger) (proto.NodeInfo, error)
}

type NodeSightConfig struct {
	Config config.Config
	Logger hclog.Logger
	Gpt    gpt.Gpt
	Ssh    ssh.Ssh
}

type nodesight struct {
	cfg *NodeSightConfig
}

func NodeSightNew(_ context.Context, cfg *NodeSightConfig) NodeSight {
	return &nodesight{
		cfg: cfg,
	}
}

func DefaultNodeSightConfig() *NodeSightConfig {
	return &NodeSightConfig{}
}

func (ns *nodesight) Init(_ context.Context) error {
	ns.cfg.Logger.Debug("nodesight: Init")

	return nil
}

func (ns *nodesight) Deinit(_ context.Context) error {
	ns.cfg.Logger.Debug("nodesight: Deinit")

	return nil
}

func (ns *nodesight) Run(ctx context.Context, trigger *proto.NodeTrigger) (proto.NodeInfo, error) {
	ns.cfg.Logger.Debug("nodesight: Run")

	var info proto.NodeInfo

	if err := ns.cfg.Ssh.Init(ctx, &trigger.SshConfig); err != nil {
		return info, errors.Wrap(err, "failed to init ssh")
	}

	defer func() {
		_ = ns.cfg.Ssh.Deinit(ctx)
	}()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(routineNum)

	g.Go(func() error {
		defer func(ns *nodesight, ctx context.Context) {
			_ = ns.runClean(ctx)
		}(ns, ctx)
		if err := ns.runDetect(ctx); err != nil {
			return errors.Wrap(err, "failed to run detect")
		}
		health, err := ns.runHealth(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to run health")
		}
		stat, err := ns.runStat(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to run stat")
		}
		info.NodeStat = *stat
		report, err := ns.runReport(ctx, health, stat)
		if err != nil {
			return errors.Wrap(err, "failed to run report")
		}
		info.NodeReport = *report
		return nil
	})

	if err := g.Wait(); err != nil {
		return info, errors.Wrap(err, "failed to wait routine")
	}

	return info, nil
}

func (ns *nodesight) runDetect(ctx context.Context) error {
	ns.cfg.Logger.Debug("nodesight: runDetect")

	cmds := []string{
		fmt.Sprintf("curl -s -u%s:%s -L %s -o %s",
			ns.cfg.Config.Spec.ArtifactConfig.User,
			ns.cfg.Config.Spec.ArtifactConfig.Pass,
			ns.cfg.Config.Spec.ArtifactConfig.Url+"/"+artifactPath+"/"+agentScript,
			agentPath+agentScript),
		fmt.Sprintf("cd %s; bash %s %s %s %s %s %s %s",
			agentPath,
			agentScript,
			ns.cfg.Config.Spec.ArtifactConfig.User,
			ns.cfg.Config.Spec.ArtifactConfig.Pass,
			ns.cfg.Config.Spec.ArtifactConfig.Url,
			artifactPath,
			agentExec,
			agentPath+agentExec),
	}

	out, err := ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return errors.Wrap(errors.New(out), err.Error())
	}

	return nil
}

func (ns *nodesight) runHealth(ctx context.Context) (string, error) {
	ns.cfg.Logger.Debug("nodesight: runHealth")

	cmds := []string{
		fmt.Sprintf("curl -s -u%s:%s -L %s -o %s",
			ns.cfg.Config.Spec.ArtifactConfig.User,
			ns.cfg.Config.Spec.ArtifactConfig.Pass,
			ns.cfg.Config.Spec.ArtifactConfig.Url+"/"+artifactPath+"/"+healthScript,
			healthPath+healthScript),
	}

	out, err := ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return "", errors.Wrap(errors.New(out), err.Error())
	}

	cmds = []string{
		fmt.Sprintf("cd %s; bash %s %s", healthPath, healthScript, healthSilent),
	}

	out, err = ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return "", errors.Wrap(errors.New(out), err.Error())
	}

	return out, nil
}

func (ns *nodesight) runStat(ctx context.Context) (*proto.NodeStat, error) {
	ns.cfg.Logger.Debug("nodesight: runStat")

	var stat proto.NodeStat

	cmds := []string{
		fmt.Sprintf("%s %s %s",
			agentPath+agentExec,
			agentDurationTime+agentSep+ns.cfg.Config.Spec.NodeConfig.Duration,
			agentLogLevel+agentSep+"ERROR"),
	}

	out, err := ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return &stat, errors.Wrap(errors.New(out), err.Error())
	}

	if err := json.Unmarshal([]byte(out), &stat); err != nil {
		return &stat, errors.Wrap(err, "failed to unmarshal json")
	}

	return &stat, nil
}

func (ns *nodesight) runReport(_ context.Context, health string, stat *proto.NodeStat) (*proto.NodeReport, error) {
	ns.cfg.Logger.Debug("nodesight: runReport")

	var report proto.NodeReport

	// TBD: FIXME

	return &report, nil
}

func (ns *nodesight) runClean(ctx context.Context) error {
	ns.cfg.Logger.Debug("nodesight: runClean")

	cmds := []string{
		fmt.Sprintf("rm -f %s", agentPath+agentScript),
		fmt.Sprintf("rm -f %s", healthPath+healthScript),
	}

	out, err := ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return errors.Wrap(errors.New(out), err.Error())
	}

	return nil
}
