package helpers

import "errors"

var (
	ErrUnsupportedQueryType = errors.New("unsupported query type")
	ErrRecordNotFound       = errors.New("requested record not found")
	ErrNotImplemented       = errors.New("handler not implemented for this query type")
)
