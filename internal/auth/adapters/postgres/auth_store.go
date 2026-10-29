package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/osugodbless/forum-backend.git/internal/auth/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	pool *pgxpool.Pool
}

func NewUserModel(pool *pgxpool.Pool) *UserModel {
	return &UserModel{pool: pool}
}

func (u *UserModel) Create(ctx context.Context, user *domain.User) error {

	const query = `
		INSERT INTO auth.users (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING created_at`

	err := u.pool.QueryRow(ctx, query, user.ID, user.Email, user.PasswordHash).Scan(&user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("%w: %w", domain.ErrUserAlreadyExists, err)
			default:
				return fmt.Errorf("unexpected error creating user: %w", err)
			}
		}
	}

	return nil
}

// func (u *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
// 	return nil, nil
// }

// func (u *UserModel) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
// 	return nil, nil
// }

// func (u *UserModel) Update(ctx context.Context, user *User) error {
// 	return nil
// }
