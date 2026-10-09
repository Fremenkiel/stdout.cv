package testutil

import (
	"bytes"
	"log"
	"net/http"
)

type Renderer struct {
	Writer *ResponseWriter	
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) RenderPage(w http.ResponseWriter, name string, data any) error {
	r.Writer =w.(*ResponseWriter) 

	var buf bytes.Buffer

	buf.WriteString("page:")
	buf.WriteString(name)

	_, err := r.Writer.Write(buf.Bytes())
		log.Printf("wrote header %d", r.Writer.status)
	return err
}

func (r *Renderer) RenderFragment(w http.ResponseWriter, name string, data any) error {
	r.Writer =w.(*ResponseWriter) 

	var buf bytes.Buffer

	buf.WriteString("fragment:")
	buf.WriteString(name)

	_, err := r.Writer.Write(buf.Bytes())
	return err
}
