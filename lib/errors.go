package lib

import (
	"errors"
)

var (
	ErrNotValidPassword = errors.New("not valid password")
	ErrNotValidEmail    = errors.New("not valid email string")
	ErrNoteValidPhone   = errors.New("not valid phone string")
	ErrUnexpectedFromDB = errors.New("unexpected database error")
	ErrUserAlreadyExist = errors.New("user already exist")
)
