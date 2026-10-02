package main

import (
	"fmt"
	"html/template"
	"os"
)

func main() {
	base := template.Must(template.New("").Parse(`{{ define "main" }}MAIN: {{ template "body" . }}{{ end }}`))
	
	tmpl := template.Must(base.Clone())
	// Emulate ParseFS by naming the parsed template "user.html"
	tmpl = template.Must(tmpl.New("user.html").Parse(`{{ define "body" }}BODY: {{ .Name }}{{ end }}{{ template "main" . }}`))
	
	err := tmpl.ExecuteTemplate(os.Stdout, "user.html", map[string]string{"Name": "Kevin"})
	if err != nil {
		fmt.Println("Error:", err)
	}
}
