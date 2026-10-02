package main

import (
	"fmt"
	"html/template"
	"os"
)

func main() {
	tmpl := template.New("test")
	// Using <script src=""></script> instead of <script src="" />
	_, err := tmpl.Parse(`{{ define "header" }}<script src=""></script>{{ end }}{{ template "header" . }}`)
	if err != nil {
		fmt.Println("Parse Error:", err)
		return
	}
	err = tmpl.Execute(os.Stdout, nil)
	if err != nil {
		fmt.Println("\nExecute Error:", err)
	}
}
