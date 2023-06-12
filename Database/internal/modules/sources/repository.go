package sources

import (
	"context"

	"digest_bot_database/internal/apperrors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSource(ctx context.Context, source *Source) error {
	sql, args, err := squirrel.
		Insert("sources").
		Columns("userid", "source").
		Values(source.UserID, source.Source).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (r *Repository) GetUsersIDList(ctx context.Context) ([]string, error) {
	sql, _, err := squirrel.
		Select("userid").
		From("sources").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.
		Query(ctx, sql)
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

func (r *Repository) GetSourcesByUserID(ctx context.Context, userID int) ([]string, error) {
	sql, args, err := squirrel.
		Select("source").
		From("sources").
		Where("userid = $1", userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.
		Query(ctx, sql, args...)
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

func (r *Repository) GetNewVideosForUserSources(ctx context.Context, source *Source) ([]string, error) {
	sql, args, err := squirrel.
		Select("new_vids").
		From("sources").
		Where("userid = $1 AND source = $2", source.UserID, source.Source).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&source.NewVids)
	return source.NewVids, err
}

// func (r *Repository) UpdateNewVidsForSourceByUserID(ctx context.Context, source *Source) error {
// 	sql, args, err := squirrel.
// 		Update("sources").
// 		Set("new_vids = $1", source.Source).
// 		Where("userid = $1", source.UserID).
// 		PlaceholderFormat(squirrel.Dollar).
// 		ToSql()
// 	if err != nil {
// 		return err
// 	}

// 	_, err = r.db.Exec(ctx, sql, args...)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *Repository) UpdateNewVidsForSource(ctx context.Context, newVids []string, source string) error {
	sql, args, err := squirrel.
		Update("sources").
		Set("new_vids = ARRAY_CAT(new_vids, $1)", newVids).
		Where("source = $1", source).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

// func (r *Repository) UpdateTodaysDigestForSourceByUserID(ctx context.Context, fullDigest string) error {
// 	sql, args, err := squirrel.
// 		Update("sources").
// 		Set("todays_digest = ARRAY_APPEND(new_vids, $1)", fullDigest).
// 		Where("userid = $1", source.UserID).
// 		PlaceholderFormat(squirrel.Dollar).
// 		ToSql()
// 	if err != nil {
// 		return err
// 	}

// 	_, err = r.db.Exec(ctx, sql, args...)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *Repository) PurgeNewVidsAndTodaysDigestColumns(ctx context.Context) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE sources
		SET new_vids = '{}',
		todays_digest = '';
		`)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteSourceByLink(ctx context.Context, source *Source) error {
	sql, args, err := squirrel.
		Delete("sources").
		Where("userid = $1 AND source = $2", source.UserID, source.Source).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	n, err := r.db.Exec(ctx, sql, args...)
	switch {
	case n.RowsAffected() == 0:
		return apperrors.NotFound("source", source)
	case err == nil:
		return apperrors.Internal(err)
	}

	return nil
}
