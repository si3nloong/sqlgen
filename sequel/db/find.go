package db

import (
	"context"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

// FindOne is to find single record using primary key.
func FindOne[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, db sequel.DB, v Ptr) error {
	var (
		pkName, _, pk = v.PK()
		dialect       = sequel.DefaultDialect()
		columns       = v.Columns()
		stmt          = strpool.AcquireString()
	)
	defer strpool.ReleaseString(stmt)

	stmt.WriteString("SELECT ")
	for i := range columns {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(dialect.Wrap(columns[i]))
	}
	stmt.WriteString(" FROM " + dialect.Wrap(v.TableName()) + " WHERE " + dialect.Wrap(pkName) + " = " + dialect.Var(1) + " LIMIT 1;")
	return db.QueryRowContext(ctx, stmt.String(), pk).Scan(v.Addrs()...)
}
