package repository

import (
	"auth-test/domain"
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
