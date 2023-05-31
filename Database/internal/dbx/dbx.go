package dbx

import (
	"errors"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsUniqueViolation(err error, source string) bool {
	var perr *pgconn.PgError
	if errors.As(err, &perr) {
		return perr.Code == pgerrcode.UniqueViolation && strings.Contains(perr.ConstraintName, source)
	}

	return false
}

func IsNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
