//go:build linux

package linters

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	artifactPath = "oxsecurity"
	artifactName = "megalinter-cupcake:latest"
)

type MegaLinter interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, string, []string) ([]string, error)
}

type MegaLinterConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type megalinter struct {
	cfg    *MegaLinterConfig
	client *client.Client
}

type artifactAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func MegaLinterNew(_ context.Context, cfg *MegaLinterConfig) MegaLinter {
	return &megalinter{
		cfg: cfg,
	}
}

func DefaultMegaLinterConfig() *MegaLinterConfig {
	return &MegaLinterConfig{}
}

func (ml *megalinter) Init(ctx context.Context) error {
	ml.cfg.Logger.Debug("megalinter: Init")

	var err error

	ml.client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Wrap(err, "failed to new client")
	}

	if err = ml.pullImage(ctx); err != nil {
		_ = ml.client.Close()
		return errors.Wrap(err, "failed to pull image")
	}

	return nil
}

func (ml *megalinter) Deinit(ctx context.Context) error {
	ml.cfg.Logger.Debug("megalinter: Deinit")

	defer func(client *client.Client) {
		if client != nil {
			_ = client.Close()
		}
	}(ml.client)

	if err := ml.removeImage(ctx); err != nil {
		return errors.Wrap(err, "failed to remove image")
	}

	return nil
}

func (ml *megalinter) Run(ctx context.Context, path string, files []string) ([]string, error) {
	ml.cfg.Logger.Debug("megalinter: Run")

	report, err := ml.runContainer(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run container")
	}

	buf, err := ml.parseReport(ctx, report)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run container")
	}

	return buf, nil
}

func (ml *megalinter) pullImage(ctx context.Context) error {
	ml.cfg.Logger.Debug("megalinter: pullImage")

	helper := func(v interface{}) string {
		var buf bytes.Buffer
		encoder := base64.NewEncoder(base64.StdEncoding, &buf)
		_ = json.NewEncoder(encoder).Encode(v)
		_ = encoder.Close()
		return buf.String()
	}

	auth := artifactAuth{
		Username: ml.cfg.Config.Spec.ArtifactConfig.User,
		Password: ml.cfg.Config.Spec.ArtifactConfig.Pass,
	}

	options := image.PullOptions{
		RegistryAuth: helper(auth),
	}

	out, err := ml.client.ImagePull(ctx, artifactPath+"/"+artifactName, options)
	if err != nil {
		return errors.Wrap(err, "failed to pull image")
	}

	_ = out.Close()

	return nil
}

func (ml *megalinter) removeImage(ctx context.Context) error {
	ml.cfg.Logger.Debug("megalinter: removeImage")

	options := image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	}

	if _, err := ml.client.ImageRemove(ctx, artifactPath+"/"+artifactName, options); err != nil {
		return errors.Wrap(err, "failed to remove image")
	}

	return nil
}

func (ml *megalinter) runContainer(ctx context.Context) ([]byte, error) {
	ml.cfg.Logger.Debug("megalinter: runContainer")

	var buf []byte

	resp, err := ml.client.ContainerCreate(ctx, &container.Config{}, &container.HostConfig{}, &network.NetworkingConfig{},
		nil, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create container")
	}

	defer func(ctx context.Context, ml *megalinter, id string) {
		opts := container.RemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   true,
			Force:         true,
		}
		_ = ml.client.ContainerRemove(ctx, id, opts)
	}(ctx, ml, resp.ID)

	if err = ml.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, errors.Wrap(err, "failed to start container")
	}

	// TBD: FIXME
	// Load output into buf

	return buf, nil
}

func (ml *megalinter) parseReport(_ context.Context, data []byte) ([]string, error) {
	ml.cfg.Logger.Debug("megalinter: parseReport")

	var buf []string

	// TBD: FIXME

	return buf, nil
}
