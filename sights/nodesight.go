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

	artifactPath = "/devops-pipeflow/plugins/"

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

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(routineNum)

	g.Go(func() error {
		defer func(ns *nodesight, ctx context.Context, cfg *proto.SshConfig) {
			_ = ns.runClean(ctx, cfg)
		}(ns, ctx, &trigger.SshConfig)
		if err := ns.runDetect(ctx, &trigger.SshConfig); err != nil {
			return errors.Wrap(err, "failed to run detect")
		}
		health, err := ns.runHealth(ctx, &trigger.SshConfig)
		if err != nil {
			return errors.Wrap(err, "failed to run health")
		}
		stat, err := ns.runStat(ctx, &trigger.SshConfig)
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
		info.Error = err.Error()
		return info, errors.Wrap(err, "failed to wait routine")
	}

	return info, nil
}

func (ns *nodesight) runDetect(ctx context.Context, cfg *proto.SshConfig) error {
	ns.cfg.Logger.Debug("nodesight: runDetect")

	if err := ns.cfg.Ssh.Init(ctx, cfg); err != nil {
		return errors.Wrap(err, "failed to init ssh")
	}

	defer func() {
		_ = ns.cfg.Ssh.Deinit(ctx)
	}()

	cmds := []string{
		fmt.Sprintf("curl -s -u%s:%s -L %s -o %s",
			ns.cfg.Config.Spec.ArtifactConfig.User,
			ns.cfg.Config.Spec.ArtifactConfig.Pass,
			ns.cfg.Config.Spec.ArtifactConfig.Url+artifactPath+agentScript,
			agentPath+agentScript),
		fmt.Sprintf("cd %s; bash %s %s %s %s %s",
			agentPath,
			agentScript,
			ns.cfg.Config.Spec.ArtifactConfig.User,
			ns.cfg.Config.Spec.ArtifactConfig.Pass,
			ns.cfg.Config.Spec.ArtifactConfig.Url+artifactPath+agentExec,
			agentPath+agentExec),
	}

	out, err := ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return errors.Wrap(err, "failed to run ssh")
	}
	if out != "" {
		return errors.Wrap(errors.New(out), "failed to run detect")
	}

	return nil
}

func (ns *nodesight) runHealth(ctx context.Context, cfg *proto.SshConfig) (string, error) {
	ns.cfg.Logger.Debug("nodesight: runHealth")

	if err := ns.cfg.Ssh.Init(ctx, cfg); err != nil {
		return "", errors.Wrap(err, "failed to init ssh")
	}

	defer func() {
		_ = ns.cfg.Ssh.Deinit(ctx)
	}()

	cmds := []string{
		fmt.Sprintf("curl -s -u%s:%s -L %s -o %s",
			ns.cfg.Config.Spec.ArtifactConfig.User,
			ns.cfg.Config.Spec.ArtifactConfig.Pass,
			ns.cfg.Config.Spec.ArtifactConfig.Url+artifactPath+healthScript,
			healthPath+healthScript),
		fmt.Sprintf("cd %s; bash %s %s", healthPath, healthScript, healthSilent),
	}

	out, err := ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return "", errors.Wrap(err, "failed to run ssh")
	}

	return out, nil
}

func (ns *nodesight) runStat(ctx context.Context, cfg *proto.SshConfig) (*proto.NodeStat, error) {
	ns.cfg.Logger.Debug("nodesight: runStat")

	var stat proto.NodeStat

	if err := ns.cfg.Ssh.Init(ctx, cfg); err != nil {
		return &stat, errors.Wrap(err, "failed to init ssh")
	}

	defer func() {
		_ = ns.cfg.Ssh.Deinit(ctx)
	}()

	cmds := []string{
		fmt.Sprintf("%s %s %s",
			agentPath+agentExec,
			agentDurationTime+agentSep+ns.cfg.Config.Spec.NodeConfig.Duration,
			agentLogLevel+agentSep+"ERROR"),
	}

	out, err := ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return &stat, errors.Wrap(err, "failed to run ssh")
	}

	if err := json.Unmarshal([]byte(out), &stat); err != nil {
		return &stat, errors.Wrap(err, "failed to unmarshal json")
	}

	return &stat, nil
}

func (ns *nodesight) runClean(ctx context.Context, cfg *proto.SshConfig) error {
	ns.cfg.Logger.Debug("nodesight: runClean")

	if err := ns.cfg.Ssh.Init(ctx, cfg); err != nil {
		return errors.Wrap(err, "failed to init ssh")
	}

	defer func() {
		_ = ns.cfg.Ssh.Deinit(ctx)
	}()

	cmds := []string{
		fmt.Sprintf("rm -f %s", agentPath+agentScript),
		fmt.Sprintf("rm -f %s", healthPath+healthScript),
	}

	_, err := ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return errors.Wrap(err, "failed to run ssh")
	}

	return nil
}

func (ns *nodesight) runReport(_ context.Context, health string, stat *proto.NodeStat) (*proto.NodeReport, error) {
	ns.cfg.Logger.Debug("nodesight: runReport")

	var report proto.NodeReport

	// TBD: FIXME

	return &report, nil
}
