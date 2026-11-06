package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type ProfileModel struct {
	pool DBTX
}

func (m *ProfileModel) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func NewProfileModel(pool DBTX) *ProfileModel {
	return &ProfileModel{pool: pool}
}

func (p *ProfileModel) WithTx(tx pgx.Tx) *ProfileModel {
	return &ProfileModel{pool: tx}
}

func (p *ProfileModel) InsertDefaultProfile(ctx context.Context, profile *UserProfile) error {
	const query = `
		INSERT INTO profile.profiles (user_id, username, full_name, display_name)
		VALUES ($1, $2, $3, $4)
	`
	_, err := p.pool.Exec(ctx, query, profile.UserID, profile.Username, profile.FullName, profile.DisplayName)
	if err != nil {
		return translateErr("Insert Profile", err)
	}

	return nil
}

// translateErr converts Postgres errors into domain errors.
// Anything unrecognized is wrapped with context and returned as-is.
func translateErr(op string, err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			if pgErr.ConstraintName == "uq_profile_username" {
				return ErrUsernameTaken
			} else {
				return ErrProfileExists
			}
		case pgerrcode.CheckViolation:
			// data should have been validated already, so if we get here it's a bug upstream.
			return fmt.Errorf("%s: check constraint %q violated: %w", op, pgErr.ConstraintName, err)
		}
	}

	return nil
}
