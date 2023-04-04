package sql

import (
	"context"
	"database/sql"
	"strings"
)

// FindOne is to find single record using primary key.
func FindOne[T KeyValuer[T], Ptr KeyValueScanner[T]](ctx context.Context, db DB, v *T) error {
	var t = (*v)
	pk, err := t.PK()
	if err != nil {
		return err
	}

	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)
	stmt.WriteQuery("SELECT `" + strings.Join(t.Columns(), dialect.Wrap(",")) + "` FROM " + dialect.Wrap(t.Table()) + " ")
	stmt.WriteQuery("WHERE "+dialect.Wrap(t.PKName())+" = "+dialect.Var(1)+" LIMIT 1;", pk)

	rows, err := db.QueryContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	var result = Ptr(&t).Addrs()
	if err = rows.Scan(result...); err != nil {
		return err
	}
	return nil
}
