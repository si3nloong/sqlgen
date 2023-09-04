package sql

import (
	"context"
)

// FindOne is to find single record using primary key.
func FindOne[T KeyValuer[T], Ptr KeyValueScanner[T]](ctx context.Context, db DB, v Ptr) error {
	pkName, _, pk := v.PK()
	columns := v.Columns()
	stmt := acquireString()
	defer releaseString(stmt)

	stmt.WriteString("SELECT ")
	for i, col := range columns {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(dialect.Wrap(col))
	}
	stmt.WriteString(" FROM " + dialect.Wrap(v.Table()) + " WHERE " + dialect.Wrap(pkName) + " = " + dialect.Var(1) + " LIMIT 1;")

	return db.QueryRowContext(ctx, stmt.String(), pk).Scan(v.Addrs()...)
}
