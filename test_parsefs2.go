package main

import (
	"fmt"
	"io/fs"
	"github.com/fremenkiel/stdout.cv/ui"
)

func main() {
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	fmt.Println("Pages:", pages, "Err:", err)
}
