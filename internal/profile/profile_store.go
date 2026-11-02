package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileModel struct {
	pool *pgxpool.Pool
}

func NewProfileModel(pool *pgxpool.Pool) *ProfileModel {
	return &ProfileModel{pool: pool}
}

func (p *ProfileModel) InsertProfile(ctx context.Context, profile *UserProfile) error {
	const query = `
		INSERT INTO profile.profiles (user_id, username, full_name, display_name)
		VALUES ($1, $2, $3, $4)
	`
	_, err := p.pool.Exec(ctx, query, profile.UserID, profile.Username, profile.FullName, profile.DisplayName)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if pgErr.ConstraintName == "uq_profile_username" {
					return fmt.Errorf("%w: %w", ErrUsernameTaken, err)
				} else {
					return fmt.Errorf("%w: %w", ErrProfileExists, err)
				}
			default:
				return fmt.Errorf("unexpected error creating user profile: %w", err)
			}
		}
	}

	return nil
}
