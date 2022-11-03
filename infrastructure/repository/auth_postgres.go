package repository

import (
	"auth-test/domain"
	"auth-test/lib"
	"auth-test/pkg/logging"
	"context"
	"database/sql"
)

type authUserRepositoryDB struct {
	db  *sql.DB
	log *logging.Logger
}

func NewAuthUserRepositoryDb(dbClinet *sql.DB, log *logging.Logger) *authUserRepositoryDB {
	return &authUserRepositoryDB{
		db:  dbClinet,
		log: log,
	}
}

func (u *authUserRepositoryDB) CreateUser(ctx context.Context, authUser domain.AuthUser) (string, error) {
	var lastInsertID string
	err := u.db.QueryRowContext(ctx, "INSERT INTO auth_users(created_at, email, password_hash, login, phone) VALUES($1, $2, $3, $4, $5) returning id;",
		authUser.CreatedAt,
		authUser.Email,
		authUser.PasswordHash,
		authUser.Login,
		authUser.Phone).Scan(&lastInsertID)
	if err != nil {
		return "", err
	}

	return lastInsertID, nil
}

func (a *authUserRepositoryDB) GetAuthUserByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	query := `SELECT id, created_at, email, password_hash, login, phone FROM auth_users WHERE email=$1`

	var user domain.AuthUser
	err := a.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Email,
		&user.PasswordHash,
		&user.Login,
		&user.Phone,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			a.log.Error("error while scanning book " + err.Error())
			return nil, lib.ErrUnexpectedFromDB
		}
	}
	return &user, nil
}
