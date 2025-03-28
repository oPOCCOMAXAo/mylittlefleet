package hx

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey struct{}

type Context struct {
	IsHX   bool
	Target string
}

func MiddlewareInjectContext(ctx *gin.Context) {
	prev := ctx.Request.Context()

	value := Context{
		IsHX:   IsHX(ctx),
		Target: GetTarget(ctx),
	}

	ctx.Request = ctx.Request.WithContext(
		context.WithValue(prev, contextKey{}, value),
	)
}

func GetContext(ctx context.Context) Context {
	value, ok := ctx.Value(contextKey{}).(Context)
	if !ok {
		return Context{}
	}

	return value
}
