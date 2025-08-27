package sqlfunc

import (
	"context"
	"iter"
	"strings"

	"github.com/deadblue/sqltmpl"
)

type (
	// QueryFunc executes a retrieving query, returns the first result.
	QueryFunc[P, R any] func(context.Context, P) (R, error)
	// QueryFunc executes a retrieving query without parameters, returns the
	// first result.
	QueryFunc0[R any] func(context.Context) (R, error)

	// QuerySeqFunc executes a retrieving query, returns an iterator for all
	// matched results.
	QuerySeqFunc[P, R any] func(context.Context, P) (iter.Seq[R], error)
	// QuerySeqFunc0 executes a retrieving query without parameters, returns
	// an iterator for all matched results.
	QuerySeqFunc0[R any] func(context.Context) (iter.Seq[R], error)

	// UpdateFunc executes an updating query, returns the number or rows
	// affected by updating query.
	UpdateFunc[P any] func(context.Context, P) (int64, error)
)

// MakeQueryFunc makes a QueryFunc from query template lines.
func MakeQueryFunc[P, R any](lines ...string) (f QueryFunc[P, R], err error) {
	preprocessResult[R]()
	tmpl, err := sqltmpl.Parse[P](lines...)
	if err == nil {
		f = func(ctx context.Context, params P) (result R, err error) {
			query, args, err := tmpl.Render(params)
			if err != nil {
				return
			}
			return queryRow[R](ctx, query, args...)
		}
	}
	return
}

// MakeQueryFunc0 makes a QueryFunc0 from query lines.
func MakeQueryFunc0[R any](lines ...string) (f QueryFunc0[R], err error) {
	preprocessResult[R]()
	query := strings.Join(lines, "\n")
	f = func(ctx context.Context) (result R, err error) {
		return queryRow[R](ctx, query)
	}
	return
}

// MakeQuerySeqFunc makes a QuerySeqFunc from query template lines.
func MakeQuerySeqFunc[P, R any](lines ...string) (f QuerySeqFunc[P, R], err error) {
	preprocessResult[R]()
	tmpl, err := sqltmpl.Parse[P](lines...)
	if err == nil {
		f = func(ctx context.Context, params P) (results iter.Seq[R], err error) {
			query, args, err := tmpl.Render(params)
			if err != nil {
				return
			}
			return queryRows[R](ctx, query, args...)
		}
	}
	return
}

// MakeQuerySeqFunc0 makes a QuerySeqFunc0 from query lines.
func MakeQuerySeqFunc0[R any](lines ...string) (f QuerySeqFunc0[R], err error) {
	preprocessResult[R]()
	query := strings.Join(lines, "\n")
	f = func(ctx context.Context) (results iter.Seq[R], err error) {
		return queryRows[R](ctx, query)
	}
	return
}

// MakeUpdateFunc makes an UpdateFunc from query template lines.
func MakeUpdateFunc[P any](lines ...string) (f UpdateFunc[P], err error) {
	tmpl, err := sqltmpl.Parse[P](lines...)
	if err != nil {
		return
	}
	f = func(ctx context.Context, params P) (affectedRows int64, err error) {
		executor, err := fromContext(ctx)
		if err != nil {
			return
		}
		query, args, err := tmpl.Render(params)
		if err != nil {
			return
		}
		ret, err := executor.ExecContext(ctx, query, args...)
		if err == nil {
			affectedRows, err = ret.RowsAffected()
		}
		return
	}
	return
}
