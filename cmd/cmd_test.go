package cmd

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/devops-pipeflow/insight-plugin/config"
)

func testInitConfig() *config.Config {
	cfg := config.New()

	fi, _ := os.Open("../test/config/config.yml")

	defer func() {
		_ = fi.Close()
	}()

	buf, _ := io.ReadAll(fi)
	_ = yaml.Unmarshal(buf, cfg)

	return cfg
}

func TestInitLogger(t *testing.T) {
	ctx := context.Background()

	logger, err := initLogger(ctx, level)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, logger)
}

func TestInitConfig(t *testing.T) {
	logger, _ := initLogger(context.Background(), level)
	ctx := context.Background()

	_, err := initConfig(ctx, logger, "invalid.yml")
	assert.NotEqual(t, nil, err)

	_, err = initConfig(ctx, logger, "../test/config/invalid.yml")
	assert.NotEqual(t, nil, err)

	_, err = initConfig(ctx, logger, "../test/config/config.yml")
	assert.Equal(t, nil, err)
}

func TestInitGpt(t *testing.T) {
	logger, _ := initLogger(context.Background(), level)
	cfg := testInitConfig()

	_, err := initGpt(context.Background(), logger, cfg)
	assert.Equal(t, nil, err)
}

func TestInitRepo(t *testing.T) {
	logger, _ := initLogger(context.Background(), level)
	cfg := testInitConfig()

	_, err := initRepo(context.Background(), logger, cfg)
	assert.Equal(t, nil, err)
}

func TestInitReview(t *testing.T) {
	logger, _ := initLogger(context.Background(), level)
	cfg := testInitConfig()

	_, err := initReview(context.Background(), logger, cfg)
	assert.Equal(t, nil, err)
}

// nolint: dogsled
func TestInitSights(t *testing.T) {
	logger, _ := initLogger(context.Background(), level)
	cfg := testInitConfig()

	_, _, _, err := initSights(context.Background(), logger, cfg)
	assert.Equal(t, nil, err)
}

func TestInitReport(t *testing.T) {
	logger, _ := initLogger(context.Background(), level)
	cfg := testInitConfig()

	_, err := initReport(context.Background(), logger, cfg)
	assert.Equal(t, nil, err)
}

func TestInitSsh(t *testing.T) {
	logger, _ := initLogger(context.Background(), level)
	cfg := testInitConfig()

	_, err := initSsh(context.Background(), logger, cfg)
	assert.Equal(t, nil, err)
}

func TestInitInsight(t *testing.T) {
	logger, _ := initLogger(context.Background(), level)
	cfg := testInitConfig()

	_, err := initInsight(context.Background(), logger, cfg, nil, nil, nil, nil, nil, nil, nil, nil)
	assert.Equal(t, nil, err)
}
