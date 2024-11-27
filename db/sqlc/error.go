package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// defines error handling for psql database interactions and handling of the pgx driver in Go
const (
	ForeignKeyViolation = "23503"
	UniqueViolations    = "23505"
)

// identify situations where a database query doesn't return any results
var ErrRecordNotFound = pgx.ErrNoRows

// custom error object that set code field to unique viloation,
// when unique constraint violations encountered during db operations
var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolations,
}

// takes an error (err) as input and assign it to *pgconn.PgError
// if its successfull, returns the error code with the psql error.
// this improves readability and maintainability of the code
func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
