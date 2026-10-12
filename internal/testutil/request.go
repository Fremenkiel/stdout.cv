package testutil

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
)

func NewRequest(ctx context.Context, method, url string, cookies []*http.Cookie) *http.Request {
	buf := &bytes.Buffer{}

	request := httptest.NewRequestWithContext(
		ctx,
		method,
		url,
		buf,
		)

	var cookieSlice []string
	for i := range cookies {
		cookieSlice = append(cookieSlice, cookies[i].String())
	}
	request.Header["Cookie"] = cookieSlice

	return request
}
