package repository

import (
	"auth-test/domain"
	"context"
)

type AuthUserRepo interface {
	CreateUser(ctx context.Context, authUser domain.AuthUser) (string, error)
}