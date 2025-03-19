//nolint:varnamelen,containedctx,gochecknoglobals
package ginutils

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin/render"
	"github.com/pkg/errors"
)

var Default = &HTMLTemplRenderer{}

func NewTemplRender(
	prev render.HTMLRender,
) *HTMLTemplRenderer {
	return &HTMLTemplRenderer{
		FallbackHTMLRenderer: prev,
	}
}

type HTMLTemplRenderer struct {
	FallbackHTMLRenderer render.HTMLRender
}

func (r *HTMLTemplRenderer) Instance(s string, d any) render.Render {
	component, ok := d.(templ.Component)
	if !ok {
		if r.FallbackHTMLRenderer != nil {
			return r.FallbackHTMLRenderer.Instance(s, d)
		}
	}

	return &Renderer{
		Context:   context.Background(),
		Status:    -1,
		Component: component,
	}
}

type Renderer struct {
	Context   context.Context
	Status    int
	Component templ.Component
}

func (t Renderer) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)

	if t.Status != -1 {
		w.WriteHeader(t.Status)
	}

	if t.Component != nil {
		return errors.WithStack(t.Component.Render(t.Context, w))
	}

	return nil
}

func (t Renderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
