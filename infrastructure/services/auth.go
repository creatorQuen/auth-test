package services

import (
	"auth-test/domain"
	"auth-test/dto"
	"auth-test/infrastructure/repository"
	"auth-test/lib"
	"auth-test/pkg/logging"
	"context"
	"time"
)

type authUserService struct {
	repoAuthUser repository.AuthUserRepo
	log          *logging.Logger
}

func NewAuthUserService(repoAuthUser repository.AuthUserRepo, log *logging.Logger) *authUserService {
	return &authUserService{
		repoAuthUser: repoAuthUser,
		log:          log,
	}
}

func (a *authUserService) CreateAuthUser(ctx context.Context, req dto.AuthUserRegisterReq) (index string, err error) {
	var authUser domain.AuthUser

	authUser.CreatedAt = time.Now().Format(lib.DbTLayout)
	authUser.Email = req.Email
	authUser.PasswordHash = req.Password
	authUser.Login = req.Login
	authUser.Phone = req.Phone

	return a.repoAuthUser.CreateUser(ctx, authUser)
}
