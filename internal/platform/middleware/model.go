package middleware

import "net/http"

type Middleware func(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request))
