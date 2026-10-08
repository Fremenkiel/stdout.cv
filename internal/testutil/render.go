package testutil

import (
	"net/http"
)

type Renderer struct {
	pageRendered			map[string]any
	fragmentRendered	map[string]any
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) RenderPage(w http.ResponseWriter, name string, data any) error {
	r.pageRendered[name] = data

	return nil
}

func (r *Renderer) RenderFragment(w http.ResponseWriter, name string, data any) error {
	r.fragmentRendered[name] = data

	return nil
}
