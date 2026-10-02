package main

import (
	"fmt"
	"net/http/httptest"
	"html/template"
	"io/fs"
	"path/filepath"
	"github.com/fremenkiel/stdout.cv/ui"
)

func main() {
	base := template.Must(template.New("").ParseFS(ui.Files, "html/layouts/*.html", "html/components/*.html"))
	
	cache := map[string]*template.Template{}
	pages, _ := fs.Glob(ui.Files, "html/pages/*.html")

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl := template.Must(base.Clone())
		tmpl = template.Must(tmpl.ParseFS(ui.Files, page))
		cache[name] = tmpl
	}
	
	w := httptest.NewRecorder()
	
	err := cache["user.html"].ExecuteTemplate(w, "user.html", map[string]string{"Name": "Kevin"})
	
	fmt.Printf("Status: %d\n", w.Result().StatusCode)
	fmt.Printf("Body length: %d\n", len(w.Body.Bytes()))
	fmt.Printf("Error: %v\n", err)
}
