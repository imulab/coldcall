package codecall

import "errors"

var (
	ErrNoBody  = errors.New("body is nil")
	ErrBadBody = errors.New("body is invalid")
)
