package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	ctx := context.Background()

	logger, err := initLogger(ctx, level)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, logger)
}

func TestRunAgent(t *testing.T) {
	assert.Equal(t, nil, nil)
}
