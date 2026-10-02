package render

import "net/http"

type Renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}
