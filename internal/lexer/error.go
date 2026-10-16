package lexer

import "errors"

var ErrQueryNotAllowed = errors.New("query: The received query is not allowed")
var ErrMultipleWildcard = errors.New("query: More than one wildcard char was given")
var ErrMultipleQueries = errors.New("query: More than one query was given")
var ErrQueryEndMarkerMissing = errors.New("query: End marker for query is missing, add a semi-colon")
var ErrNoQuery = errors.New("query: No query received")

// Lexer
var ErrInvalidString = errors.New("String is not correctly closed, add end to string")
var ErrInvalidQuery = errors.New("Query is invalid formatted")
