package sql

import (
	"context"
	"database/sql"
	"strings"
)

// SelectOneFrom is to find single record using primary key.
func SelectOneFrom[T KeyValuer[T], Ptr KeyValueScanner[T]](ctx context.Context, db DB) (*T, error) {
	var t T
	pk, err := t.PK()
	if err != nil {
		return nil, err
	}

	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)
	stmt.WriteQuery("SELECT `" + strings.Join(t.Columns(), dialect.Wrap(",")) + "` FROM " + dialect.Wrap(t.Table()) + " ")
	stmt.WriteQuery("WHERE "+dialect.Wrap(t.PKName())+" = "+dialect.Var(1)+" LIMIT 1;", pk)

	rows, err := db.QueryContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var result = Ptr(&t).Addrs()
	if err = rows.Scan(result...); err != nil {
		return nil, err
	}
	return &t, nil
}

func SelectFrom[T Scanner[T]](ctx context.Context, db DB) error {
	return nil
}
