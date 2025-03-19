package ginutils

import "github.com/gin-gonic/gin"

type CtxValuePointer[T any] string

func (v CtxValuePointer[T]) GetOK(ctx *gin.Context) (*T, bool) {
	value, ok := ctx.Get(string(v))
	if !ok {
		return nil, false
	}

	typ, ok := value.(*T)
	if !ok {
		return nil, false
	}

	return typ, true
}

func (v CtxValuePointer[T]) Get(ctx *gin.Context) *T {
	value, _ := v.GetOK(ctx)

	return value
}

func (v CtxValuePointer[T]) Set(ctx *gin.Context, value *T) {
	ctx.Set(string(v), value)
}

func (v CtxValuePointer[T]) Has(ctx *gin.Context) bool {
	value, _ := v.GetOK(ctx)

	return value != nil
}
