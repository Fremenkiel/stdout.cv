package testutil

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	header				map[string][]string
	status				int
	wroteHeader,
	sentResponse,
	hasContent		bool
	buffer				bytes.Buffer
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header: make(map[string][]string),
		status: -1,
		wroteHeader: false,
		sentResponse: false,
		hasContent: false,
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.wroteHeader = true
}

func (w *ResponseWriter) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	w.hasContent = true
	w.sentResponse = true

	return w.buffer.Write(p)
}

func (w *ResponseWriter) Status() int {
	return w.status
}
