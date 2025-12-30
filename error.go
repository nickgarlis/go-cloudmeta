package cloudmeta

import "errors"

var (
	ErrUnknownProvider = errors.New("unknown cloud provider")
	ErrNotFound        = errors.New("not found")
)
