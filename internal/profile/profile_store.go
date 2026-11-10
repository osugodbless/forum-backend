package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/osugodbless/forum-backend/internal/platform/database"
)

type ProfileModel struct {
	pool database.DBTX
}

func NewProfileModel(pool database.DBTX) *ProfileModel {
	return &ProfileModel{pool: pool}
}

func (m *ProfileModel) db(ctx context.Context) database.DBTX {
	return database.Executor(ctx, m.pool)
}

func (p *ProfileModel) InsertDefaultProfile(ctx context.Context, profile *UserProfile) error {
	const query = `
		INSERT INTO profile.profiles (user_id, username, full_name, display_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT ON CONSTRAINT uq_profile_username DO NOTHING
	`
	_, err := p.db(ctx).Exec(ctx, query, profile.UserID, profile.Username, profile.FullName, profile.DisplayName)
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

	return fmt.Errorf("%s: %w", op, err)
}
