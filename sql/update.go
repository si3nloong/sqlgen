package sql

import (
	"context"
	"database/sql"
)

type SingleUpdater interface {
	UpdateQuery() (string, []any)
}

// UpdateOne is to update single record using primary key.
func UpdateOne[T KeyValuer[T]](ctx context.Context, db DB, v T) (sql.Result, error) {
	switch vi := any(v).(type) {
	case SingleUpdater:
		query, args := vi.UpdateQuery()
		return db.ExecContext(ctx, query, args...)
	}

	return nil, nil
	// idx, pk := v.PK()
	// stmt := acquireString()
	// defer releaseString(stmt)
	// columns, values := v.Columns(), v.Values()
	// stmt.WriteString("UPDATE " + dialect.Wrap(v.Table()) + " SET ")
	// noOfCols := len(columns)
	// for i := 0; i < noOfCols; i++ {
	// 	if i > 0 {
	// 		stmt.WriteByte(',')
	// 	}
	// 	stmt.WriteString(dialect.Wrap(columns[i]) + " = " + dialect.Var(i+1))
	// }

	// // TODO: support `UUID()` etc
	// switch vi := any(v).(type) {
	// case Keyer:
	// 	stmt.WriteString(" WHERE " + dialect.Wrap(vi.PKName()) + " = " + dialect.Var(noOfCols+1) + ";")
	// 	return db.ExecContext(ctx, stmt.String(), pk)
	// default:
	// 	return db.ExecContext(ctx, stmt.String(), values...)
	// }
}
