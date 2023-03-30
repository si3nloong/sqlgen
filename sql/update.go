package sql

import (
	"context"
	"database/sql"
)

// UpdateOne is to update single record using primary key.
func UpdateOne[T KeyValuer[T]](ctx context.Context, db DB, v T) (sql.Result, error) {
	pk, err := v.PK()
	if err != nil {
		return nil, err
	}

	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)
	columns := v.Columns()
	values := v.Values()
	stmt.WriteQuery("UPDATE " + dialect.Wrap(v.Table()) + " SET ")
	noOfCols := len(columns)
	for i := 0; i < noOfCols; i++ {
		if i > 0 {
			stmt.WriteQuery(",")
		}
		stmt.WriteQuery(dialect.Wrap(columns[i])+" = "+dialect.Var(i+1), values[i])
	}

	// TODO: support `UUID()` etc
	switch vi := any(v).(type) {
	case Keyer:
		stmt.WriteQuery(" WHERE "+dialect.Wrap(vi.PKName())+" = "+dialect.Var(noOfCols+2)+";", pk)
	default:

	}

	return db.ExecContext(ctx, stmt.Query(), stmt.Args()...)
}
