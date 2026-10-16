package schema

import "errors"

var ErrFieldNotWildcard = errors.New("query: The given index do not point to a wildcard field")
