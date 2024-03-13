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

func timeDuration() func() {
	start := time.Now()

	return func() {
		d := time.Since(start)
		fmt.Printf("time duration = %v\n", d)
	}
}

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

	defer timeDuration()()

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchCpuStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchDiskStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchDiskStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchDockerStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	d, _ := time.ParseDuration(agentDuration)
	stat, _ := fetchDockerStat(ctx, logger, d)

	_, err := json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchHostStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchHostStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchLoadStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchLoadStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchMemStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchMemStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchNetStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchNetStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchProcessStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	d, _ := time.ParseDuration(agentDuration)
	stat, err := fetchProcessStat(ctx, logger, d)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}
