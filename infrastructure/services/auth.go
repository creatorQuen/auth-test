package services

import (
	"auth-test/domain"
	"auth-test/dto"
	"auth-test/infrastructure/repository"
	"auth-test/lib"
	"auth-test/pkg/logging"
	"context"
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

type authUserService struct {
	repoAuthUser repository.AuthUserRepo
	log          *logging.Logger
	salt         string
}

func NewAuthUserService(repoAuthUser repository.AuthUserRepo, log *logging.Logger, salt string) *authUserService {
	if salt == "" {
		log.Fatal("salt is empty")
	}
	return &authUserService{
		repoAuthUser: repoAuthUser,
		log:          log,
	}
}

func (a *authUserService) CreateAuthUser(ctx context.Context, req dto.AuthUserRegisterReq) (index string, err error) {
	var authUser domain.AuthUser

	authUser.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if !isValidEmail(req.Email) {
		a.log.Error(lib.ErrNotValidEmail)
		return "", lib.ErrNotValidEmail
	}

	if !isValidPassword(req.Password) {
		a.log.Error(lib.ErrNotValidPassword)
		return "", lib.ErrNotValidPassword
	}

	dbUser, err := a.repoAuthUser.GetAuthUserByEmail(ctx, authUser.Email)
	if err != nil {
		a.log.Error("repoAuthUser.GetAuthUserByEmail: ", err)
		return "", err
	}
	if dbUser != nil {
		a.log.Error("repoAuthUser.GetAuthUserByEmail: ", lib.ErrUserAlreadyExist)
		return "", lib.ErrUserAlreadyExist
	}

	authUser.CreatedAt = time.Now().Format(lib.DbTLayout)
	authUser.PasswordHash = generatePasswordHash(req.Password, a.salt)
	authUser.Login = req.Login
	authUser.Phone = req.Phone

	return a.repoAuthUser.CreateUser(ctx, authUser)
}

func generatePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
