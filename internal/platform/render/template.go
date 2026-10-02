package render

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/fremenkiel/stdout.cv/ui"
)

type TemplateRenderer struct {
	pageCache	map[string]*template.Template
	base			*template.Template
}

var _ Renderer = (*TemplateRenderer)(nil)

func NewTemplateRenderer() *TemplateRenderer {
	base := template.Must(template.New("").ParseFS(ui.Files, "html/layouts/*.html", "html/components/*.html"))

	cache := map[string]*template.Template{}
	pages, _ := fs.Glob(ui.Files, "html/pages/*.html")

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl := template.Must(base.Clone())
		tmpl = template.Must(tmpl.ParseFS(ui.Files, page))
		cache[name] = tmpl
	}

	return &TemplateRenderer{pageCache: cache, base: base}
}

func (r *TemplateRenderer) RenderPage(w http.ResponseWriter, name string, data any) error {
	filename := fmt.Sprintf("%s.html", name)

	temp, ok := r.pageCache[filename]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}

	return temp.ExecuteTemplate(w, filename, data)
}

func (r *TemplateRenderer) RenderFragment(w http.ResponseWriter, name string, data any) error {
	return r.base.ExecuteTemplate(w, name, data)
}
