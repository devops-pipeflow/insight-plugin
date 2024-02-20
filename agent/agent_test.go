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
	assert.Equal(t, nil, err)
	assert.Equal(t, agentDuration, duration)

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

	stat, err := fetchCpuStat(ctx, logger, agentDuration)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchDiskStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	stat, err := fetchDiskStat(ctx, logger, agentDuration)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchDockerStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	stat, _ := fetchDockerStat(ctx, logger, agentDuration)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchHostStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	stat, err := fetchHostStat(ctx, logger, agentDuration)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchLoadStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	stat, err := fetchLoadStat(ctx, logger, agentDuration)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchMemStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	stat, err := fetchMemStat(ctx, logger, agentDuration)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchNetStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	stat, err := fetchNetStat(ctx, logger, agentDuration)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}

func TestFetchProcessStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	stat, err := fetchProcessStat(ctx, logger, agentDuration)
	assert.Equal(t, nil, err)

	buf, _ := json.Marshal(stat)
	fmt.Println(string(buf))
}
