package ginutils

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func RenderSSE(
	ctx *gin.Context,
	event string,
	component templ.Component,
) {
	var body bytes.Buffer

	_ = component.Render(ctx.Request.Context(), &body)
	lines := bytes.Split(body.Bytes(), []byte("\n"))

	_, _ = ctx.Writer.Write([]byte("event: "))
	_, _ = ctx.Writer.Write([]byte(event))
	_, _ = ctx.Writer.Write([]byte("\n"))

	for _, line := range lines {
		_, _ = ctx.Writer.Write([]byte("data: "))
		_, _ = ctx.Writer.Write(line)
		_, _ = ctx.Writer.Write([]byte("\n"))
	}

	_, _ = ctx.Writer.Write([]byte("\n"))
	ctx.Writer.Flush()
}
