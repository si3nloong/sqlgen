package sql

import (
	"context"
)

func SelectFrom[T Valuer[T], Ptr Scanner[T]](ctx context.Context, db DB) ([]T, error) {
	var v T

	stmt := acquireString()
	defer releaseString(stmt)
	stmt.WriteString("SELECT ")
	for i, col := range v.Columns() {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(dialect.Wrap(col))
	}
	stmt.WriteString(" FROM " + dialect.Wrap(v.Table()) + ";")

	rows, err := db.QueryContext(ctx, stmt.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]T, 0)
	for rows.Next() {
		var v T
		if err := rows.Scan(Ptr(&v).Addrs()...); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
