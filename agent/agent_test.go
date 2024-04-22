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

func TestRunAgent(t *testing.T) {
	t.Skip("Skipping TestRunAgent.")
}

func TestFetchCpuStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	stat, err := fetchCpuStat(ctx, logger)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchDiskStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	stat, err := fetchDiskStat(ctx, logger)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchDockerStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	stat, _ := fetchDockerStat(ctx, logger)

	_, err := json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchHostStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	stat, err := fetchHostStat(ctx, logger)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchLoadStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	stat, err := fetchLoadStat(ctx, logger)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchMemStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	stat, err := fetchMemStat(ctx, logger)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchNetStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	stat, err := fetchNetStat(ctx, logger)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}

func TestFetchProcessStat(t *testing.T) {
	ctx := context.Background()
	logger, _ := initLogger(ctx, agentLevel)

	defer timeDuration()()

	stat, err := fetchProcessStat(ctx, logger)
	assert.Equal(t, nil, err)

	_, err = json.Marshal(stat)
	assert.Equal(t, nil, err)
}
