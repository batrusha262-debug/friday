package contextx

import (
	"context"
	"fmt"
)

type TraceID string

func (t TraceID) String() string {
	return string(t)
}

type contextKeyTraceID struct{}

func WithTraceID(ctx context.Context, traceID TraceID) context.Context {
	return context.WithValue(ctx, contextKeyTraceID{}, traceID)
}

func TraceIDFromContext(ctx context.Context) (TraceID, error) {
	traceID, ok := ctx.Value(contextKeyTraceID{}).(TraceID)
	if !ok {
		return "", fmt.Errorf("traceID: %w", ErrNoValue)
	}

	return traceID, nil
}
