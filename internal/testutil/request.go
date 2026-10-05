package testutil

import (
	"bytes"
	"context"
	"net/http"
)

func NewRequest(ctx context.Context, method, url string, cookies []*http.Cookie) *http.Request {
	buf := &bytes.Buffer{}

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		buf,
		)
	if err != nil {
		panic("unable to create request")
	}

	var cookieSlice []string
	for i := range cookies {
		cookieSlice = append(cookieSlice, cookies[i].String())
	}
	request.Header["Cookie"] = cookieSlice

	return request
}
