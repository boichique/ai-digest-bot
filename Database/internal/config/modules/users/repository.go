package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSource(ctx context.Context, user *User) error {
	_, err := r.db.
		Exec(
			ctx,
			`INSERT INTO sources (userID, source) 
			VALUES ($1, $2)`,
			user.ID, user.Source)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) GetUserSourcesByID(ctx context.Context, userID int) (*User, error) {
	return nil, fmt.Errorf("not implemented")
}
