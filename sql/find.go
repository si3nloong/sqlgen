package sql

import (
	"context"
)

// FindOne is to find single record using primary key.
func FindOne[T KeyValuer[T], Ptr KeyValueScanner[T]](ctx context.Context, db DB, v Ptr) error {
	var t = (*v)
	pk, err := t.PK()
	if err != nil {
		return err
	}

	stmt := acquireString()
	defer releaseString(stmt)
	stmt.WriteString("SELECT ")
	for i, col := range t.Columns() {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(dialect.Wrap(col))
	}
	stmt.WriteString(" FROM " + dialect.Wrap(t.Table()) + " WHERE " + dialect.Wrap(t.PKName()) + " = " + dialect.Var(1) + " LIMIT 1;")

	return db.QueryRowContext(ctx, stmt.String(), pk).Scan(v.Addrs()...)
}
