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

	if !isValidPhone(req.Phone) {
		a.log.Error(lib.ErrNoteValidPhone)
		return "", lib.ErrNoteValidPhone
	}
	authUser.Phone = req.Phone

	dbUserByEmail, err := a.repoAuthUser.GetAuthUserByEmail(ctx, authUser.Email)
	if err != nil {
		a.log.Error("repoAuthUser.GetAuthUserByEmail: ", err)
		return "", err
	}
	if dbUserByEmail != nil {
		a.log.Error("repoAuthUser.GetAuthUserByEmail: ", lib.ErrUserEmailAlreadyExist)
		return "", lib.ErrUserEmailAlreadyExist
	}

	authUser.Login = strings.TrimSpace(req.Login)
	dbUserByLogin, err := a.repoAuthUser.GetAuthUserByLogin(ctx, authUser.Login)
	if err != nil {
		a.log.Error("repoAuthUser.GetAuthUserByLogin: ", err)
		return "", err
	}
	if dbUserByLogin != nil {
		a.log.Error("repoAuthUser.GetAuthUserByLogin: ", lib.ErrUserLoginAlreadyExist)
		return "", lib.ErrUserLoginAlreadyExist
	}

	authUser.CreatedAt = time.Now().Format(lib.DbTLayout)
	authUser.PasswordHash = generatePasswordHash(req.Password, a.salt)
	authUser.Login = req.Login

	return a.repoAuthUser.CreateUser(ctx, authUser)
}

func generatePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
