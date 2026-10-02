package main

import (
	"fmt"
	"html/template"
	"github.com/fremenkiel/stdout.cv/ui"
)

func main() {
	_, err := template.New("").ParseFS(ui.Files, "html/layouts/*.html", "html/components/*.html")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Success")
	}
}
