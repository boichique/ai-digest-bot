package sources

import (
	"context"
	"log"

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

	return err
}

func (r *Repository) GetUserSourcesByID(ctx context.Context, userID int) ([]string, error) {
	var source string

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
		if err := rows.Scan(&source); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (r *Repository) DeleteSourceByLink(ctx context.Context, source *Source) error {
	_, err := r.db.
		Exec(
			ctx,
			`DELETE FROM sources
			WHERE userid = $1 
			AND source = $2;`,
			source.UserID, source.Source)

	log.Print(err)
	return err
}
