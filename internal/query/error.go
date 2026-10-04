package query

import "errors"

var ErrQueryNotAllowed = errors.New("query: The received query is not allowed")
var ErrFieldNotWildcard = errors.New("query: The given index do not point to a wildcard field")
