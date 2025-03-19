package ginutils

import "github.com/gin-gonic/gin"

type CtxValue[T any] string

func (v CtxValue[T]) GetOK(ctx *gin.Context) (T, bool) {
	var typ T

	value, ok := ctx.Get(string(v))
	if ok {
		typ, ok = value.(T)
	}

	return typ, ok
}

func (v CtxValue[T]) Get(ctx *gin.Context) T {
	value, _ := v.GetOK(ctx)

	return value
}

func (v CtxValue[T]) Set(ctx *gin.Context, value T) {
	ctx.Set(string(v), value)
}

func (v CtxValue[T]) Has(ctx *gin.Context) bool {
	_, has := v.GetOK(ctx)

	return has
}

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
