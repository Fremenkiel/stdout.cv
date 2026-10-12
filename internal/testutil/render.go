package testutil

import (
	"bytes"
	"fmt"
	"net/http"
)

type Renderer struct {}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) RenderPage(w http.ResponseWriter, name string, data any) error {
	var buf bytes.Buffer

	buf.WriteString("page:")
	buf.WriteString(name)
	buf.WriteString(fmt.Sprint(data))

	_, err := w.Write(buf.Bytes())
	return err
}

func (r *Renderer) RenderFragment(w http.ResponseWriter, name string, data any) error {
	var buf bytes.Buffer

	buf.WriteString("fragment:")
	buf.WriteString(name)
	buf.WriteString(fmt.Sprint(data))

	_, err := w.Write(buf.Bytes())
	return err
}
