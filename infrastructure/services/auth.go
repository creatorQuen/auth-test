package services

import (
	"auth-test/config"
	"auth-test/domain"
	"auth-test/dto"
	"auth-test/infrastructure/repository"
	"auth-test/lib"
	"auth-test/pkg/logging"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type authUserService struct {
	repoAuthUser    repository.AuthUserRepo
	log             *logging.Logger
	authServiceConf config.Config
}

func NewAuthUserService(repoAuthUser repository.AuthUserRepo, log *logging.Logger, authServiceConf *config.Config) *authUserService {
	if authServiceConf.ConfigAuthService.Salt == "" {
		log.Fatal("salt is empty")
	}
	if authServiceConf.ConfigAuthService.SignKey == "" {
		log.Fatal("sign key is empty")
	}
	if authServiceConf.ConfigAuthService.TokenTimeLeft == 0 {
		log.Fatal("token time left is empty")
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
	authUser.Login = req.Login

	authUser.PasswordHash = generatePasswordHash(req.Password, a.authServiceConf.ConfigAuthService.Salt)
	authUser.CreatedAt = time.Now().Format(lib.DbTLayout)

	return a.repoAuthUser.CreateUser(ctx, authUser)
}

func generatePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (a *authUserService) GenerateToken(ctx context.Context, req dto.AuthUserLoginReq) (string, error) {
	var passwordHashDB string
	var userID string

	if strings.Contains(req.LoginOrEmail, "@") {
		if !isValidEmail(req.LoginOrEmail) {
			a.log.Error(lib.ErrNotValidEmail)
			return "", lib.ErrNotValidEmail
		}

		dbAuthUser, err := a.repoAuthUser.GetAuthUserByEmail(ctx, req.LoginOrEmail)
		if err != nil {
			a.log.Error("repoAuthUser.GetAuthUserByEmail: ", err)
			return "", err
		}
		if dbAuthUser == nil {
			a.log.Error("repoAuthUser.GetAuthUserByEmail: ", lib.ErrUserNotExist)
			return "", lib.ErrUserNotExist
		}
		passwordHashDB = dbAuthUser.PasswordHash
		userID = dbAuthUser.Id

	} else {
		dbAuthUser, err := a.repoAuthUser.GetAuthUserByLogin(ctx, req.LoginOrEmail)
		if err != nil {
			a.log.Error("repoAuthUser.GetAuthUserByLogin: ", err)
			return "", err
		}
		if dbAuthUser == nil {
			a.log.Error("repoAuthUser.GetAuthUserByLogin: ", lib.ErrUserNotExist)
			return "", lib.ErrUserNotExist
		}

		passwordHashDB = dbAuthUser.PasswordHash
		userID = dbAuthUser.Id
	}

	hashInputPassword := generatePasswordHash(req.Password, a.authServiceConf.ConfigAuthService.Salt)

	if passwordHashDB != hashInputPassword {
		return "", lib.ErrPasswordNotEqual
	}

	tokenTL := a.authServiceConf.ConfigAuthService.TokenTimeLeft
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * (time.Duration(tokenTL))).Unix(),
	})

	return token.SignedString([]byte(a.authServiceConf.ConfigAuthService.SignKey))
}
