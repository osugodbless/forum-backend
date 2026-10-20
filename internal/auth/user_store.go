package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	db *pgxpool.Pool
}

func NewUserModel(db *pgxpool.Pool) *UserModel {
	return &UserModel{db: db}

}

func (u *UserModel) Create(ctx context.Context, user *User) error {
	return nil
}

func (u *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (u *UserModel) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return nil, nil
}

func (u *UserModel) Update(ctx context.Context, user *User) error {
	return nil
}
