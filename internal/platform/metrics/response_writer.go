package metrics

import "net/http"

type ResponseWriter struct {
	base					http.ResponseWriter
	headerWritten	bool
	Code					int
	WrittenLen		int
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

func NewResponseWriter(b http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		base: b,
		Code: http.StatusOK,
		WrittenLen: -1,
	}
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.base.WriteHeader(code)

	if !(code >= 100 && code <= 199) && !w.headerWritten {
		w.Code = code
		w.headerWritten = true
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.base.Header()
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	l, err := w.base.Write(body)
	if err != nil {
		return -1, err
	}

	w.WrittenLen = l

	return l, nil
}
