package contextutils

import (
	"context"
)

//nolint:containedctx
type detachedContext struct {
	context.Context
	orig context.Context
}

// DetachedContext returns a new context detached from the lifetime
// of ctx, but which still returns the values of ctx.
//
// DetachedContext can be used to maintain the trace context required
// to correlate events, but where the operation is "fire-and-forget",
// and should not be affected by the deadline or cancellation of ctx.
func Detached(ctx context.Context) context.Context {
	return &detachedContext{Context: context.Background(), orig: ctx}
}

// Value returns c.orig.Value(key).
func (c *detachedContext) Value(key any) any {
	return c.orig.Value(key)
}
