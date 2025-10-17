package sqlfunc

import (
	"context"
	"database/sql"
	"iter"
)

func queryRow[R any](
	ctx context.Context, query string, args ...any,
) (result sql.Null[R], err error) {
	executor, err := fromContext(ctx)
	if err != nil {
		return
	}
	// Execute query
	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	// Early return for no-rows
	if !rows.Next() {
		err = rows.Err()
		return
	}
	// Parse result
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	df, err := makeDestFunc[R](columns)
	if err != nil {
		return
	}
	dest := df(&result.V)
	if err = rows.Scan(dest...); err == nil {
		result.Valid = true
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
	defer func() {
		if err != nil {
			rows.Close()
		}
	}()
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	df, err := makeDestFunc[R](columns)
	if err == nil {
		results = func(yield func(R) bool) {
			defer rows.Close()
			for rows.Next() {
				var result R
				dest := df(&result)
				if err := rows.Scan(dest...); err != nil {
					continue
				}
				if !yield(result) {
					break
				}
			}
		}
	}
	return
}
