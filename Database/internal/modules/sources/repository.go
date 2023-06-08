package sources

import (
	"context"
	"log"

	"digest_bot_database/internal/apperrors"
	"digest_bot_database/internal/dbx"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSource(ctx context.Context, source *Source) error {
	_, err := r.db.
		Exec(
			ctx,
			`INSERT INTO sources (userid, source) 
			VALUES ($1, $2);`,
			source.UserID, source.Source)

	switch {
	case dbx.IsUniqueViolation(err, "source"):
		return apperrors.AlreadyExists("source", source)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (r *Repository) GetUsersIDList(ctx context.Context) ([]string, error) {
	rows, err := r.db.
		Query(
			ctx,
			`SELECT userid FROM sources`,
		)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, apperrors.Internal(err)
		}
		usersIDs = append(usersIDs, userID)
	}

	return usersIDs, nil
}

func (r *Repository) GetUserSourcesByID(ctx context.Context, userID int) ([]string, error) {
	rows, err := r.db.
		Query(
			ctx,
			`SELECT source FROM sources WHERE userid = $1;`,
			userID,
		)
	if err != nil {
		return nil, err
	}

	var sources []string
	for rows.Next() {
		var source string
		if err := rows.Scan(&source); err != nil {
			return nil, apperrors.Internal(err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (r *Repository) DeleteSourceByLink(ctx context.Context, source *Source) error {
	n, err := r.db.
		Exec(
			ctx,
			`DELETE FROM sources
			WHERE userid = $1 
			AND source = $2;`,
			source.UserID, source.Source)
	if err != nil {
		log.Print(err)
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("source", source)
	}

	return nil
}
