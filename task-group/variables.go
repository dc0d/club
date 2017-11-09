package tgroup

import "github.com/dc0d/errgo/sentinel"

// errors
var (
	ErrTimeout = sentinel.Errorf("TIMEOUT")
)
