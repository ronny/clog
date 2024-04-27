package trace

import "errors"

var (
	ErrUnparseable  = errors.New("unparseable")
	ErrInvalidTrace = errors.New("invalid trace")
)
