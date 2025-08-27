package sqlfunc

import (
	"context"
	"database/sql"
)

// Executor interface declares required methods to execute query, that are
// implemented by *[sql.DB], *[sql.Conn] and *[sql.Tx].
type Executor interface {
	// ExecContext executes a query without returning any rows.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
