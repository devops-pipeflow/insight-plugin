//go:build linux

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	ctx := context.Background()

	logger, err := initLogger(ctx, agentLevel)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, logger)
}

func TestInitDuration(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	duration, err := initDuration(ctx, logger, "")
	d, _ := time.ParseDuration(agentDuration)
	assert.Equal(t, nil, err)
	assert.Equal(t, d, duration)

	duration, err = initDuration(ctx, logger, "1s")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1*time.Second, duration)

	duration, err = initDuration(ctx, logger, "10m")
	assert.Equal(t, nil, err)
	assert.Equal(t, 10*time.Minute, duration)

	duration, err = initDuration(ctx, logger, "100h")
	assert.Equal(t, nil, err)
	assert.Equal(t, 100*time.Hour, duration)
}

func TestRunAgent(t *testing.T) {
	t.Skip("Skipping TestRunAgent.")
}

func TestFetchCpuStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchCpuStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchDiskStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchDiskStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchDockerStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	d, _ := time.ParseDuration(agentDuration)
	stat, _ := fetchDockerStat(ctx, logger, d)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchHostStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchHostStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchLoadStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchLoadStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchMemStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchMemStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchNetStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchNetStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchProcessStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchProcessStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}
