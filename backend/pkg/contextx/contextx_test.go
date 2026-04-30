package contextx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"friday/pkg/contextx"
)

// -----------------------------------------------------------------------
// TraceID
// -----------------------------------------------------------------------

func TestWithTraceID_roundtrip(t *testing.T) {
	ctx := contextx.WithTraceID(context.Background(), "abc-123")

	got, err := contextx.TraceIDFromContext(ctx)

	require.NoError(t, err)
	assert.Equal(t, contextx.TraceID("abc-123"), got)
	assert.Equal(t, "abc-123", got.String())
}

func TestTraceIDFromContext_missing(t *testing.T) {
	_, err := contextx.TraceIDFromContext(context.Background())

	require.Error(t, err)
	assert.ErrorIs(t, err, contextx.ErrNoValue)
}

// -----------------------------------------------------------------------
// Logger
// -----------------------------------------------------------------------

func TestWithLogger_roundtrip(t *testing.T) {
	logger := slog.Default()
	ctx := contextx.WithLogger(context.Background(), logger)

	got, err := contextx.LoggerFromContext(ctx)

	require.NoError(t, err)
	assert.Same(t, logger, got)
}

func TestLoggerFromContext_missing(t *testing.T) {
	_, err := contextx.LoggerFromContext(context.Background())

	require.Error(t, err)
	assert.ErrorIs(t, err, contextx.ErrNoValue)
}

func TestLoggerFromContextOrDefault_returnsDefault(t *testing.T) {
	got := contextx.LoggerFromContextOrDefault(context.Background())

	assert.NotNil(t, got)
}

func TestLoggerFromContextOrDefault_returnsStored(t *testing.T) {
	logger := slog.Default()
	ctx := contextx.WithLogger(context.Background(), logger)

	got := contextx.LoggerFromContextOrDefault(ctx)

	assert.Same(t, logger, got)
}

func TestEnrichLogger_addsAttrs(t *testing.T) {
	ctx := contextx.EnrichLogger(context.Background(), "key", "value")

	logger := contextx.LoggerFromContextOrDefault(ctx)

	assert.NotNil(t, logger)
}
