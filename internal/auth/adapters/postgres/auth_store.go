package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/osugodbless/forum-backend.git/internal/auth/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	pool *pgxpool.Pool
}

func NewUserModel(pool *pgxpool.Pool) *UserModel {
	return &UserModel{pool: pool}
}

// CreateUser inserts a new user into the database and returns nil if user is created.
// If the user already exists, it returns domain.ErrUserAlreadyExists.
// If any other error occurs, it returns a wrapped error with context.
func (m *UserModel) CreateUser(ctx context.Context, user *domain.User) error {

	const query = `
		INSERT INTO auth.users (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING created_at`

	err := m.pool.QueryRow(ctx, query, user.ID, user.Email, user.PasswordHash).Scan(&user.CreatedAt)
	if err != nil {
		return translateErr("Create User", err)
	}

	return nil
}

// GetUserByEmail retrieves a user from the database by their email.
// If the user does not exist, it returns domain.ErrUserNotFound.
// If any other error occurs, it returns a wrapped error with context.
func (m *UserModel) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
		FROM auth.users
		WHERE email = $1`

	return m.queryRow(ctx, "Get User By Email", query, email)

}

// GetUserByID retrieves a user from the database by their ID.
// It returns same errors as GetUserByEmail, but uses the user's ID instead of email for lookup.
func (m *UserModel) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
		FROM auth.users
		WHERE id = $1`

	return m.queryRow(ctx, "Get User By ID", query, id)
}

// UserExistsByEmail checks if a user exists in the database by their email.
// It returns true if the user exists, false if not, and an error if any issues occur during the query.
func (m *UserModel) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM auth.users
			WHERE email = $1
		)`

	var exists bool
	err := m.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, translateErr("User Exists By Email", err)
	}

	return exists, nil
}

// UserExistsByID checks if a user exists in the database by their ID.
// It returns true if the user exists, false if not, and an error if any issues occur during the query.
func (m *UserModel) UserExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM auth.users
			WHERE id = $1
		)`

	var exists bool
	err := m.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, translateErr("User Exists By ID", err)
	}

	return exists, nil
}

// UpdatePasswordHash updates the password hash for a user with the given ID.
func (m *UserModel) UpdatePasswordHash(ctx context.Context, id uuid.UUID, newHash string) error {
	const query = `
		UPDATE auth.users
		SET password_hash = $1
		WHERE id = $2`

	tag, err := m.pool.Exec(ctx, query, newHash, id)
	if err != nil {
		return translateErr("Update Password Hash", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// DeleteUser removes a user from the database by their ID. If the user does not exist, it returns domain.ErrUserNotFound.
func (m *UserModel) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM auth.users
		WHERE id = $1`

	tag, err := m.pool.Exec(ctx, query, id)
	if err != nil {
		return translateErr("Delete User", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// queryRow runs a single-row query and scans it straight into a domain.User.
func (m *UserModel) queryRow(ctx context.Context, op, query string, arg ...any) (*domain.User, error) {
	var user domain.User
	err := m.pool.QueryRow(ctx, query, arg...).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

// translateErr converts Postgres errors into domain errors.
// Anything unrecognized is wrapped with context and returned as-is.
func translateErr(op string, err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return domain.ErrUserAlreadyExists
		case pgerrcode.CheckViolation:
			// data should have been validated already, so if we get here it's a bug upstream.
			return fmt.Errorf("%s: check constraint %q violated: %w", op, pgErr.ConstraintName, err)
		}
	}

	return fmt.Errorf("%s: %w", op, err)
}
