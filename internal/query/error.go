package query

import "errors"

var ErrQueryNotAllowed = errors.New("query: The received query is not allowed")
var ErrFieldNotWildcard = errors.New("query: The given index do not point to a wildcard field")
var ErrMultipleWildcard = errors.New("query: More than one wildcard char was given")
var ErrMultipleQueries = errors.New("query: More than one query was given")
var ErrQueryEndMarkerMissing = errors.New("query: End marker for query is missing, add a semi-colon")
var ErrNoQuery = errors.New("query: No query received")
