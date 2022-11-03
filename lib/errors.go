package lib

import (
	"errors"
)

var (
	ErrNotValidPassword = errors.New("not valid password")
	ErrNotValidEmail    = errors.New("not valid email string")
)
