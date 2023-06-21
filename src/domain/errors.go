package domain

import "errors"

var (
	ErrInternalServer  = errors.New("internal server error")
	ErrInvalidArgument = errors.New("invalid argument")
)
