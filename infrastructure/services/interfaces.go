package services

import (
	"auth-test/dto"
	"context"
)

type AuthUserService interface {
	CreateAuthUser(ctx context.Context, req dto.AuthUserRegisterReq) (string, error)
}
