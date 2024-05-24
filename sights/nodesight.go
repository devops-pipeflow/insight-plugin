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
	agentExec     = "agent"
	agentLogLevel = "--log-level"
	agentPath     = "/tmp/"
	agentScript   = agentExec + ".sh"
	agentSep      = "="

	artifactPath = "zd-devops-nj-release-generic/devops-pipeflow/plugins"

	healthPath   = "/tmp/"
	healthPlain  = "--plain"
	healthScript = "healthcheck.sh"
	healthSilent = "--silent"

	routineNum = -1
)

type NodeSight interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, *proto.NodeTrigger) (proto.NodeInfo, proto.MailInfo, error)
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

func (ns *nodesight) Run(ctx context.Context, trigger *proto.NodeTrigger) (proto.NodeInfo, proto.MailInfo, error) {
	ns.cfg.Logger.Debug("nodesight: Run")

	var nodeInfo proto.NodeInfo
	var mailInfo proto.MailInfo

	if err := ns.cfg.Ssh.Init(ctx, &trigger.SshConfig); err != nil {
		return nodeInfo, mailInfo, errors.Wrap(err, "failed to init ssh")
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
		_, err := ns.runHealth(ctx)
		if err != nil {
			nodeInfo.Error = err.Error()
			return errors.Wrap(err, "failed to run health")
		}
		stat, err := ns.runStat(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to run stat")
		}
		nodeInfo.NodeStat = *stat
		report, err := ns.runReport(ctx, stat)
		if err != nil {
			return errors.Wrap(err, "failed to run report")
		}
		nodeInfo.NodeReport = *report
		return nil
	})

	if err := g.Wait(); err != nil {
		return nodeInfo, mailInfo, errors.Wrap(err, "failed to wait routine")
	}

	return nodeInfo, mailInfo, nil
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
		fmt.Sprintf("cd %s; bash %s %s %s", healthPath, healthScript, healthPlain, healthSilent),
	}

	out, err = ns.cfg.Ssh.Run(ctx, cmds)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (ns *nodesight) runStat(ctx context.Context) (*proto.NodeStat, error) {
	ns.cfg.Logger.Debug("nodesight: runStat")

	var stat proto.NodeStat

	cmds := []string{
		fmt.Sprintf("%s %s", agentPath+agentExec, agentLogLevel+agentSep+"ERROR"),
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

func (ns *nodesight) runReport(_ context.Context, stat *proto.NodeStat) (*proto.NodeReport, error) {
	ns.cfg.Logger.Debug("nodesight: runReport")

	var report proto.NodeReport

	// TBD: FIXME
	// Split stat into cpuStat, diskStat, dockerStat, hostStat, loadStat, memStat, netStat, processStat

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
