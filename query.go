package sqlfunc

import (
	"context"
	"database/sql"
	"iter"
)

func queryRow[R any](
	ctx context.Context, query string, args ...any,
) (result R, err error) {
	executor, err := fromContext(ctx)
	if err != nil {
		return
	}
	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = sql.ErrNoRows
	} else {
		columns, _ := rows.Columns()
		dest := getResultDest(&result, columns)
		err = rows.Scan(dest...)
	}
	return
}

func queryRows[R any](
	ctx context.Context, query string, args ...any,
) (results iter.Seq[R], err error) {
	executor, err := fromContext(ctx)
	if err != nil {
		return
	}
	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	columns, err := rows.Columns()
	if err != nil {
		rows.Close()
		return
	}
	results = func(yield func(R) bool) {
		defer rows.Close()
		for rows.Next() {
			var result R
			dest := getResultDest(&result, columns)
			if err := rows.Scan(dest...); err != nil {
				continue
			}
			if !yield(result) {
				break
			}
		}
	}
	return
}
