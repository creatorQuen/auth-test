package lib

import (
	"errors"
)

var (
	ErrNotValidPassword      = errors.New("not valid password")
	ErrNotValidEmail         = errors.New("not valid email string")
	ErrNoteValidPhone        = errors.New("not valid phone string")
	ErrUnexpectedFromDB      = errors.New("unexpected database error")
	ErrUserEmailAlreadyExist = errors.New("email already exist")
	ErrUserLoginAlreadyExist = errors.New("login already exist")
)
