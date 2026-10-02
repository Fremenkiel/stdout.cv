package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
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
	
	err := cache["user.html"].ExecuteTemplate(os.Stdout, "user.html", map[string]string{"Name": "Kevin"})
	if err != nil {
		fmt.Println("Execute Error:", err)
	}
}
