package gpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	requestTimeout = 100 * time.Second
)

type Gpt interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, string) (string, error)
}

type Config struct {
	Api    string
	Config config.Config
	Logger hclog.Logger
}

type Request struct {
	Content string `json:"content"`
}

type Response struct {
	Code uint   `json:"code"`
	Msg  string `json:"msg"`
	Ret  any    `json:"ret"`
}

type gpt struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Gpt {
	return &gpt{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (g *gpt) Init(_ context.Context) error {
	g.cfg.Logger.Debug("gpt: Init")

	return nil
}

func (g *gpt) Deinit(_ context.Context) error {
	g.cfg.Logger.Debug("gpt: Deinit")

	return nil
}

func (g *gpt) Run(ctx context.Context, content string) (string, error) {
	g.cfg.Logger.Debug("gpt: Run")

	return g.sendRequest(ctx, content)
}

func (g *gpt) sendRequest(_ context.Context, content string) (string, error) {
	g.cfg.Logger.Debug("gpt: sendRequest")

	var buf Response

	request := &Request{Content: content}
	marshal, err := json.Marshal(request)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal request")
	}

	req, _ := http.NewRequest("POST", g.cfg.Api, bytes.NewBuffer(marshal))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{
		Timeout: requestTimeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to send request")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return "", errors.Wrap(err, fmt.Sprintf("invalid response: %d", res.StatusCode))
	}

	err = json.NewDecoder(res.Body).Decode(&buf)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode response")
	}

	result, ok := buf.Ret.(string)
	if !ok {
		return "", errors.Wrap(err, "invalid return value")
	}

	return result, nil
}
